package asgmanager

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/morfien101/asg-health/ec2metadatareader"
	"os"
)

func awsSession() (*session.Session, error) {
	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		guess, err := ec2metadatareader.Region()
		if err != nil {
			return nil, err
		}
		region = guess
	}
	return session.NewSession(&aws.Config{Region: aws.String(region)})
}
func asgSession(session *session.Session) *autoscaling.AutoScaling {
	return autoscaling.New(session)
}
func newAWSSession() (*autoscaling.AutoScaling, error) {
	basicSession, err := awsSession()
	if err != nil {
		return nil, err
	}
	asgSession := asgSession(basicSession)
	return asgSession, nil
}

func SetUnhealthy(instanceID string) error {
	asgSession, err := newAWSSession()
	if err != nil {
		return err
	}

	currentHealth, err := getCurrentHealth(asgSession, instanceID)
	if currentHealth == "Unhealthy" {
		return nil
	}
	return setUnhealthy(asgSession, instanceID)
}

// IsInService return ok bool, current state string and an error
func IsInService(instanceID string) (bool, string, error) {
	asgSession, err := newAWSSession()
	if err != nil {
		return false, "Not_Available", err
	}

	details, err := getInstanceDetails(asgSession, instanceID)
	if err != nil {
		return false, "Not_Available", err
	}

	if *details.LifecycleState == "InService" {
		return true, *details.LifecycleState, nil
	}
	return false, *details.LifecycleState, nil
}

func getCurrentHealth(asgSession *autoscaling.AutoScaling, instanceID string) (string, error) {
	instanceDetails, err := getInstanceDetails(asgSession, instanceID)
	if err != nil {
		return "", err
	}
	return *instanceDetails.LifecycleState, nil
}

func getInstanceDetails(asgSession *autoscaling.AutoScaling, instanceID string) (*autoscaling.InstanceDetails, error) {
	input := &autoscaling.DescribeAutoScalingInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}
	output, err := asgSession.DescribeAutoScalingInstances(input)
	if err != nil {
		return nil, err
	}
	if len(output.AutoScalingInstances) == 0 {
		return nil, fmt.Errorf("No health status found for instance %s", instanceID)
	}
	return output.AutoScalingInstances[0], nil
}

func setUnhealthy(asgSession *autoscaling.AutoScaling, instanceID string) error {
	_, err := asgSession.SetInstanceHealth(&autoscaling.SetInstanceHealthInput{
		HealthStatus: aws.String("Unhealthy"),
		InstanceId:   aws.String(instanceID),
	})
	if err != nil {
		return err
	}
	return nil
}
