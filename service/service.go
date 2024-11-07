package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func GetNoOfAppTierEC2() (int, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return -1, nil
	}

	client := ec2.NewFromConfig(cfg)
	result, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []ec2types.Filter{
			{
				Name: aws.String("image-id"),
				Values: []string{
					os.Getenv("AWS_AMI_ID"),
				},
			},
			{
				Name: aws.String("instance-state-name"),
				Values: []string{
					"running",
					"pending",
				},
			},
		},
	})

	noOfAppTierInstaces := len(result.Reservations)

	return noOfAppTierInstaces, err
}

func CreateAppTierEC2(cnt int) error {
	log.Printf("Creating app-tier-instance-%v ec2 instance", cnt)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil
	}

	client := ec2.NewFromConfig(cfg)
	if _, err = client.RunInstances(context.TODO(), &ec2.RunInstancesInput{
		MaxCount:     aws.Int32(1),
		MinCount:     aws.Int32(1),
		ImageId:      aws.String(os.Getenv("AWS_AMI_ID")),
		InstanceType: ec2types.InstanceTypeT2Micro,
		TagSpecifications: []ec2types.TagSpecification{{
			ResourceType: ec2types.ResourceTypeInstance,
			Tags: []ec2types.Tag{{
				Key:   aws.String("Name"),
				Value: aws.String(fmt.Sprintf("app-tier-instance-%v", cnt)),
			}}},
		},
		SecurityGroupIds: []string{os.Getenv("AWS_SECURITY_GROUP_ID")},
		UserData:         aws.String(os.Getenv("AWS_USER_DATA")),
	}); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetNoOfMessagesInRequestQueue() (int, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return -1, err
	}

	client := sqs.NewFromConfig(cfg)
	result, err := client.GetQueueAttributes(context.TODO(), &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(os.Getenv("AWS_REQ_URL")),
		AttributeNames: []sqstypes.QueueAttributeName{
			sqstypes.QueueAttributeNameApproximateNumberOfMessages,
		},
	})
	if err != nil {
		return -1, err
	}

	noOfMessages, err := strconv.Atoi(result.Attributes[string(sqstypes.QueueAttributeNameApproximateNumberOfMessages)])

	return noOfMessages, err
}
