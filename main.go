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

	// // Set env variable to use the correct aws profile from ./aws/credentials
	// awsProfile := "multiops"
	// // Set AWS_PROFILE env variable on OS
	// err := os.Setenv("AWS_PROFILE", awsProfile)
	// if err != nil {
	// 	log.Fatal("Failed to set AWS_PROFILE environment variable")
	// }
	//
	// svc := s3.New(session.New())
	// result, err := svc.ListBuckets(&s3.ListBucketsInput{})
	// if err != nil {
	// 	log.Fatal("Failed to list buckets!", err)
	// 	return
	// }
	//
	// fmt.Println("Buckets:")
	// for _, bucket := range result.Buckets {
	// 	fmt.Printf("%s : %s\n", aws.StringValue(bucket.Name), bucket.CreationDate)
	// }
	awsProfile := "multiops"
	sourceElb := "plinko-admin-api-production"
	destElb := "plinko-admin-api-internal"

	// Set AWS_PROFILE env variable on OS.
	err := os.Setenv("AWS_PROFILE", awsProfile)
	if err != nil {
		log.Fatal("Failed to set AWS_PROFILE environment variable", err)
	}
	// Open new elb session with aws-sdk
	svc := elb.New(session.New())

	params := &elb.DescribeInstanceHealthInput{LoadBalancerName: aws.String(sourceElb), Instances: []*elb.Instance{}}

	result, err := svc.DescribeInstanceHealth(params)
	if err != nil {
		log.Fatal("Failed to describe ELBs", err)
	}
	fmt.Println(result)

	for _, instances := range result.InstanceStates {

		id := aws.StringValue(instances.InstanceId)
		resp, err := svc.RegisterInstancesWithLoadBalancer(&elb.RegisterInstancesWithLoadBalancerInput{Instances: []*elb.Instance{{InstanceId: aws.String(id)}}, LoadBalancerName: aws.String(destElb)})
		if err != nil {
			log.Fatal("Failed to register instance", err.Error())
		}
		fmt.Println(resp)
	}

}
