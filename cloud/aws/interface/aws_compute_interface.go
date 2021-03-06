package neuronaws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"reflect"
)

type CreateServerInput struct {
	ImageId        string   `json:"imageid,omitempty"`
	InstanceType   string   `json:"instancetype,omitempty"`
	KeyName        string   `json:"keyname,omitempty"`
	MaxCount       int64    `json:"maxcount,omitempty"`
	MinCount       int64    `json:"mincount,omitempty"`
	SubnetId       string   `json:"subnetid,omitempty"`
	SecurityGroups []string `json:"securitygroup,omitempty"`
	UserData       string   `json:"userdata,omitempty"`
	AssignPubIp    bool     `json:"assignpubip,omitempty"`
}

type DescribeComputeInput struct {
	InstanceIds []string `json:"instanceids,omitempty"`
	ImageIds    []string `json:"imageids,omitempty"`
	Filters     Filters  `json:"filters,omitempty"`
}

type UpdateComputeInput struct {
	InstanceIds []string `json:"instanceids,omitempty"`
	Force       string   `json:"force,omitempty"`
}

type DeleteComputeInput struct {
	ImageId     string   `json:"ImageId,omitempty"`
	SnapshotId  string   `json:"SnapshotId,omitempty"`
	InstanceIds []string `json:"InstanceIds,omitempty"`
}

type ImageCreateInput struct {
	Description string `json:"Description,omitempty"`
	ServerName  string `json:"ServerName,omitempty"`
	InstanceId  string `json:"InstanceId,omitempty"`
}

func (sess *EstablishedSession) CreateInstance(ins *CreateServerInput) (*ec2.Reservation, error) {

	if sess.Ec2 != nil {
		if (ins.ImageId != "") || (ins.InstanceType != "") || (ins.KeyName != "") || (ins.MinCount != 0) || (ins.MaxCount != 0) || (ins.UserData != "") || (ins.SubnetId != "") || (ins.SecurityGroups != nil) {
			// support for custom ebs mapping will be rolled out soon
			create_server_input := &ec2.RunInstancesInput{
				ImageId:      aws.String(ins.ImageId),
				InstanceType: aws.String(ins.InstanceType),
				KeyName:      aws.String(ins.KeyName),
				MaxCount:     aws.Int64(ins.MaxCount),
				MinCount:     aws.Int64(ins.MinCount),
				UserData:     aws.String(ins.UserData),
				NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{{
					AssociatePublicIpAddress: aws.Bool(ins.AssignPubIp),
					DeviceIndex:              aws.Int64(0),
					DeleteOnTermination:      aws.Bool(true),
					SubnetId:                 aws.String(ins.SubnetId),
					Groups:                   aws.StringSlice(ins.SecurityGroups),
				}},
			}
			server_create_result, err := (sess.Ec2).RunInstances(create_server_input)
			// handling the error if it throws while subnet is under creation process
			if err != nil {
				return nil, err
			}
			return server_create_result, nil
		}
		return nil, fmt.Errorf("You provided empty/wrong details to CreateInstance, this is not acceptable")
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")

}

func (sess *EstablishedSession) DescribeInstance(des *DescribeComputeInput) (*ec2.DescribeInstancesOutput, error) {

	if sess.Ec2 != nil {
		if des.InstanceIds != nil {
			input := &ec2.DescribeInstancesInput{
				InstanceIds: aws.StringSlice(des.InstanceIds),
			}
			result, err := (sess.Ec2).DescribeInstances(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}

		if reflect.DeepEqual(des.Filters, Filters{}) {
			return nil, fmt.Errorf("You provided empty struct to DescribeInstance, this is not acceptable")
		}
		if (des.Filters.Name == "") || (des.Filters.Value == nil) {
			return nil, fmt.Errorf("You chose Filters to fetch server details and did not provided required value for Filters.")
		}
		input := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{Name: aws.String(des.Filters.Name),
					Values: aws.StringSlice(des.Filters.Value),
				},
			},
		}
		result, err := (sess.Ec2).DescribeInstances(input)

		if err != nil {
			return nil, err
		}
		return result, nil

	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

func (sess *EstablishedSession) DescribeAllInstances(des *DescribeComputeInput) (*ec2.DescribeInstancesOutput, error) {

	if sess.Ec2 != nil {
		input := &ec2.DescribeInstancesInput{}
		result, err := (sess.Ec2).DescribeInstances(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

func (sess *EstablishedSession) DeleteInstance(d *DeleteComputeInput) (*ec2.TerminateInstancesOutput, error) {

	if sess.Ec2 != nil {
		if d.InstanceIds != nil {
			terminate_instance_input := &ec2.TerminateInstancesInput{
				InstanceIds: aws.StringSlice(d.InstanceIds),
			}
			_, err := (sess.Ec2).TerminateInstances(terminate_instance_input)

			if err != nil {
				return nil, err
			}
			return nil, nil
		}
		return nil, fmt.Errorf("You provided empty struct to DeleteInstance, this is not acceptable")
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

func (sess *EstablishedSession) StartInstances(s *UpdateComputeInput) (*ec2.StartInstancesOutput, error) {

	if sess.Ec2 != nil {
		if s.InstanceIds != nil {
			input := &ec2.StartInstancesInput{
				InstanceIds: aws.StringSlice(s.InstanceIds),
			}
			result, err := (sess.Ec2).StartInstances(input)

			if err != nil {
				return nil, err
			}
			return result, nil

		}
		return nil, fmt.Errorf("You provided empty struct to StartInstances, this is not acceptable")
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

func (sess *EstablishedSession) StopInstances(s *UpdateComputeInput) (*ec2.StopInstancesOutput, error) {

	if sess.Ec2 != nil {
		if s.InstanceIds != nil {
			input := &ec2.StopInstancesInput{
				InstanceIds: aws.StringSlice(s.InstanceIds),
			}
			result, err := (sess.Ec2).StopInstances(input)

			if err != nil {
				return nil, err
			}
			return result, nil

		}
		return nil, fmt.Errorf("You provided empty struct to StopInstances, this is not acceptable")
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

// I will be responsible for capturing the image of the server when I am called
func (sess *EstablishedSession) CreateImage(img *ImageCreateInput) (*ec2.CreateImageOutput, error) {

	if sess.Ec2 != nil {
		if (img.ServerName != "") || (img.InstanceId != "") || (img.Description != "") {
			input := &ec2.CreateImageInput{
				Description: aws.String(img.Description),
				InstanceId:  aws.String(img.InstanceId),
				Name:        aws.String(img.ServerName),
			}
			result, err := (sess.Ec2).CreateImage(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf("You provided empty struct to CreateImage, this is not acceptable")
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

//Below function along with DeleteSnapshot has to be used to delete an image.
func (sess *EstablishedSession) DeregisterImage(img *DeleteComputeInput) error {

	if sess.Ec2 != nil {
		if img.ImageId != "" {
			// deregistering image will be done by below code
			input := &ec2.DeregisterImageInput{ImageId: aws.String(img.ImageId)}
			_, err := (sess.Ec2).DeregisterImage(input)

			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("You provided empty struct to DeregisterImage, this is not acceptable")
	}
	return fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

func (sess *EstablishedSession) DeleteSnapshot(img *DeleteComputeInput) error {

	if sess.Ec2 != nil {
		if img.SnapshotId != "" {
			// Deletion of snapshot will addressed by below code
			input := &ec2.DeleteSnapshotInput{SnapshotId: aws.String(img.SnapshotId)}
			_, err := (sess.Ec2).DeleteSnapshot(input)

			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("You provided empty struct to DeleteSnapshot, this is not acceptable")
	}
	return fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

func (sess *EstablishedSession) DescribeImages(img *DescribeComputeInput) (*ec2.DescribeImagesOutput, error) {

	if sess.Ec2 != nil {
		if img.ImageIds != nil {
			// desribing image to check if image exists
			search_image_input := &ec2.DescribeImagesInput{
				ImageIds: aws.StringSlice(img.ImageIds),
			}
			result, err := (sess.Ec2).DescribeImages(search_image_input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf("You provided empty struct to DescribeImages, this is not acceptable")
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

func (sess *EstablishedSession) DescribeAllImages(img *DescribeComputeInput) (*ec2.DescribeImagesOutput, error) {

	if sess.Ec2 != nil {
		// desribing image to check if image exists
		input := &ec2.DescribeImagesInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{Name: aws.String("is-public"),
					Values: aws.StringSlice([]string{"false"}),
				},
			},
		}
		result, err := (sess.Ec2).DescribeImages(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

func (sess *EstablishedSession) WaitTillInstanceAvailable(d *DescribeComputeInput) error {

	if sess.Ec2 != nil {
		if d.InstanceIds != nil {
			input := &ec2.DescribeInstancesInput{
				InstanceIds: aws.StringSlice(d.InstanceIds),
			}
			err := (sess.Ec2).WaitUntilInstanceRunning(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("You provided empty struct to WaitTillInstanceAvailable, this is not acceptable")
	}
	return fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

func (sess *EstablishedSession) WaitTillInstanceRunning(d *DescribeComputeInput) error {

	if sess.Ec2 != nil {
		if d.InstanceIds != nil {
			input := &ec2.DescribeInstancesInput{
				InstanceIds: aws.StringSlice(d.InstanceIds),
			}
			err := (sess.Ec2).WaitUntilInstanceRunning(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("You provided empty struct to WaitTillInstanceRunning, this is not acceptable")
	}
	return fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

func (sess *EstablishedSession) WaitTillInstanceStopped(d *DescribeComputeInput) error {

	if sess.Ec2 != nil {
		if d.InstanceIds != nil {
			input := &ec2.DescribeInstancesInput{
				InstanceIds: aws.StringSlice(d.InstanceIds),
			}
			err := (sess.Ec2).WaitUntilInstanceStopped(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("You provided empty struct to WaitTillInstanceStopped, this is not acceptable")
	}
	return fmt.Errorf("Did not get session to perform action, cannot proceed further")
}

func (sess *EstablishedSession) WaitTillInstanceTerminated(d *DescribeComputeInput) error {

	if sess.Ec2 != nil {
		if d.InstanceIds != nil {
			input := &ec2.DescribeInstancesInput{
				InstanceIds: aws.StringSlice(d.InstanceIds),
			}
			err := (sess.Ec2).WaitUntilInstanceTerminated(input)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("You provided empty struct to WaitTillInstanceTerminated, this is not acceptable")
	}
	return fmt.Errorf("Did not get session to perform action, cannot proceed further")
}
