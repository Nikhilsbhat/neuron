package awsinstance

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	aws "neuron/cloud/aws/interface"
)

type CommonComputeInput struct {
	InstanceIds []string
	VpcId       string
	GetRaw      bool
}

type CommonComputeResponse struct {
	InstanceIds []string                     `json:"InstanceIds,omitempty"`
	GetInstRaw  *ec2.DescribeInstancesOutput `json:"DescribeInstRaw,omitempty"`
}

func (des *CommonComputeInput) SearchInstance(con aws.EstablishConnectionInput) (bool, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return false, sesserr
	}
	if des.InstanceIds != nil {
		result, err := ec2.DescribeInstance(
			&aws.DescribeComputeInput{
				InstanceIds: des.InstanceIds,
			},
		)

		if err != nil {
			return false, err
		}
		if len(result.Reservations) == len(des.InstanceIds) {
			return true, nil
		}
		return false, fmt.Errorf("Failed to fetch the data of all instances you entered, found multiple entries of same values")
	}
	return false, fmt.Errorf("You provided empty struct to retrive the data, this is not acceptable")
}

func (des *CommonComputeInput) GetInstanceFromVpc(con aws.EstablishConnectionInput) (CommonComputeResponse, error) {

	//get the relative sessions before proceeding further
	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return CommonComputeResponse{}, sesserr
	}
	if des.VpcId == "" {
		return CommonComputeResponse{}, fmt.Errorf("You provided empty struct to retrive the data, this is not acceptable")
	}

	result, err := ec2.DescribeInstance(
		&aws.DescribeComputeInput{
			Filters: aws.Filters{
				Name:  "vpc-id",
				Value: []string{des.VpcId},
			},
		},
	)

	if err != nil {
		return CommonComputeResponse{}, err
	}

	if des.GetRaw == true {
		return CommonComputeResponse{GetInstRaw: result}, nil
	}

	instances := make([]string, 0)
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instances = append(instances, *instance.InstanceId)
		}
	}

	return CommonComputeResponse{InstanceIds: instances}, nil
}
