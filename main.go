package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/iamthemuffinman/logsip"
)

var log = logsip.New()

func main() {
	// Config options.
	var sourceElb string
	var destElb string
	flag.StringVar(&sourceElb, "source", "", "Source ELB")
	flag.StringVar(&destElb, "dest", "", "Destination ELB")
	flag.Parse()

	// Halt if parameters are not specified.
	if sourceElb == "" {
		log.Fatal("No source ELB specified, swapELB -h to help")
	}
	if destElb == "" {
		log.Fatal("No destination ELB specified, swapELB -h to help")
	}

	// Create a new aws session.
	stsSvc := session.New()

	// Get new temporary STS credentials for assumed role.
	getCreds := stscreds.NewCredentials(stsSvc, "<roleARNhere>")

	// Open a new elb session with the aws-sdk and pass in temporary sts credentials.
	svc := elb.New(session.New(&aws.Config{Region: aws.String("us-east-1"), Credentials: getCreds}))

	// Define parameters to pass to DescribeInstanceHealth
	params := &elb.DescribeInstanceHealthInput{LoadBalancerName: aws.String(sourceElb), Instances: []*elb.Instance{}}

	result, err := svc.DescribeInstanceHealth(params)
	if err != nil {
		log.Fatal("Failed to describe ELB: \n", err.Error())
	}

	// See instances currently registered with sourceElb.
	fmt.Printf("Currently registered instances with %s: \n", sourceElb)
	fmt.Println(result)

	// Loop through registered instances to get instance ids and register them with the destElb.
	for _, instances := range result.InstanceStates {

		id := aws.StringValue(instances.InstanceId)
		_, err := svc.RegisterInstancesWithLoadBalancer(&elb.RegisterInstancesWithLoadBalancerInput{Instances: []*elb.Instance{{InstanceId: aws.String(id)}}, LoadBalancerName: aws.String(destElb)})
		if err != nil {
			log.Fatal("Failed to register instance", err.Error())
		}

	}
	// Sleep for 20 seconds to allow instance registration in the destElb.
	fmt.Println("Going to sleep for 20 seconds before checking instance registration status ... ")
	time.Sleep(20 * time.Second)

	// Define parameters to pass to DescribeInstanceHealth for destElb
	paramsDest := &elb.DescribeInstanceHealthInput{LoadBalancerName: aws.String(destElb), Instances: []*elb.Instance{}}

	resultDest, err := svc.DescribeInstanceHealth(paramsDest)
	if err != nil {
		log.Fatal("Failed to describe ELBs", err.Error())
	}
	// Loop through Instance states for destELB and check if state is InService.
	for _, instances := range resultDest.InstanceStates {
		id := aws.StringValue(instances.InstanceId)
		state := aws.StringValue(instances.State)
		if state != "InService" {
			log.Fatalf("%s is not registered successfully with the load balancer", id)
		} else {
			fmt.Printf("%s registered successfully with the load balancer.\n", id)
		}
	}

}
