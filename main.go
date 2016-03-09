package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/iamthemuffinman/logsip"
)

var log = logsip.New(os.Stdout)

func main() {
	// Set config options.
	awsProfile := "multiops"
	sourceElb := "plinko-admin-api-production"
	destElb := "plinko-admin-api-internal"

	// Set AWS_PROFILE env variable on OS.
	err := os.Setenv("AWS_PROFILE", awsProfile)
	if err != nil {
		log.Fatal("Failed to set AWS_PROFILE environment variable", err)
	}
	// Open a new elb session with the aws-sdk.
	svc := elb.New(session.New())

	// Define parameters to pass to DescribeInstanceHealth
	params := &elb.DescribeInstanceHealthInput{LoadBalancerName: aws.String(sourceElb), Instances: []*elb.Instance{}}

	result, err := svc.DescribeInstanceHealth(params)
	if err != nil {
		log.Fatal("Failed to describe ELBs", err)
	}
	// See instances currently registered with sourceElb.
	fmt.Println(result)

	// Loop through registered instances to get instance ids and register them with the destElb.
	for _, instances := range result.InstanceStates {

		id := aws.StringValue(instances.InstanceId)
		resp, err := svc.RegisterInstancesWithLoadBalancer(&elb.RegisterInstancesWithLoadBalancerInput{Instances: []*elb.Instance{{InstanceId: aws.String(id)}}, LoadBalancerName: aws.String(destElb)})
		if err != nil {
			log.Fatal("Failed to register instance", err.Error())
		}
		fmt.Println(resp)
	}

}
