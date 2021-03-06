package awsnetwork

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/nikhilsbhat/neuron/cloud/aws/interface"
	common "github.com/nikhilsbhat/neuron/cloud/aws/operations/common"
	"strings"
)

// VpcResponse is a struct that will be the response type of almost all the VPC related activities under cloud/operations.
type VpcResponse struct {
	Name              string                           `json:"Name,omitempty"`
	Type              string                           `json:"Type,omitempty"`
	VpcId             string                           `json:"VpcId,omitempty"`
	IgwId             string                           `json:"IgwId,omitempty"`
	SecGroupIds       []string                         `json:"SecGroupId,omitempty"`
	IsDefault         bool                             `json:"IsDefault,omitempty"`
	State             string                           `json:"State,omitempty"`
	CreateVpcRaw      *ec2.CreateVpcOutput             `json:"CreateVpcRaw,omitempty"`
	GetVpcRaw         *ec2.DescribeVpcsOutput          `json:"GetVpcRaw,omitempty"`
	CreateIgwRaw      *ec2.CreateInternetGatewayOutput `json:"CreateIgwRaw,omitempty,omitempty"`
	CreateSecurityRaw *ec2.CreateSecurityGroupOutput   `json:"CreateSecRaw,omitempty,omitempty"`
}

// CreateVpc is a customized method for vpc creation, if one needs plain vpc creation then he/she has to call the GOD, interface which talks to cloud.
func (vpc *NetworkCreateInput) CreateVpc(con neuronaws.EstablishConnectionInput) (VpcResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return VpcResponse{}, seserr
	}
	// I am gathering inputs since create vpc needs it
	vpc_result, vpc_err := ec2.CreateVpc(
		&neuronaws.CreateNetworkInput{
			Cidr:    vpc.VpcCidr,
			Tenancy: "default",
		})

	// handling the error if it throws while vpc is under creation process
	if vpc_err != nil {
		return VpcResponse{}, vpc_err
	}

	// I will program wait until vpc become available
	wait_err := ec2.WaitTillVpcAvailable(
		&neuronaws.DescribeNetworkInput{
			Filters: neuronaws.Filters{
				Name: "vpc-id", Value: []string{*vpc_result.Vpc.VpcId},
			},
		},
	)
	if wait_err != nil {
		return VpcResponse{}, wait_err
	}

	// I will pass name to create_tags to set a name to the vpc
	vpctagin := common.Tag{*vpc_result.Vpc.VpcId, "Name", vpc.Name}
	vpctag, tag_err := vpctagin.CreateTags(con)
	if tag_err != nil {
		return VpcResponse{}, tag_err
	}

	// I will make the decision whether we need public network or private, based on the input I receive
	netcomp := new(NetworkComponentInput)
	netcomp.Name = vpc.Name
	netcomp.VpcIds = []string{*vpc_result.Vpc.VpcId}
	netcomp.GetRaw = vpc.GetRaw
	vpcresponse := new(VpcResponse)

	if (strings.ToLower(vpc.Type) == "public") || (strings.ToLower(vpc.Type) == "") {
		ig, ig_err := netcomp.CreateIgw(con)
		if ig_err != nil {
			return VpcResponse{}, ig_err
		}
		if vpc.GetRaw != true {
			vpcresponse.IgwId = ig.IgwIds[0]
		} else {
			vpcresponse.CreateIgwRaw = ig.CreateIgwRaw
		}
	} else if strings.ToLower(vpc.Type) == "private" {
		vpcresponse.IgwId = ""
	} else {
		return VpcResponse{}, fmt.Errorf("You provided unknown network type. There are two possibility, either we do not support this type else you would have misspelled")
	}

	// I will initialize data required to create security group and pass it to respective person to create one
	netcomp.Ports = vpc.Ports
	sec, sec_err := netcomp.CreateSecurityGroup(con)
	if sec_err != nil {
		return VpcResponse{}, sec_err
	}

	if vpc.GetRaw == true {
		vpcresponse.CreateSecurityRaw = sec.CreateSecurityRaw
		vpcresponse.CreateVpcRaw = vpc_result
		return *vpcresponse, nil
	}

	vpcresponse.SecGroupIds = sec.SecGroupIds
	vpcresponse.Name = vpctag
	vpcresponse.VpcId = *vpc_result.Vpc.VpcId
	vpcresponse.Type = vpc.Type
	return *vpcresponse, nil
}

// DeleteVpc is a customized method for vpc deletion, if one needs plain vpc deletion then he/she has to call the GOD, interface which talks to cloud.
func (vpc *DeleteNetworkInput) DeleteVpc(con neuronaws.EstablishConnectionInput) error {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return seserr
	}

	err := ec2.DeleteVpc(
		&neuronaws.DescribeNetworkInput{
			VpcIds: vpc.VpcIds,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// GetVpcs is a customized method for fetching details of all vpc for a given region, if one needs plain get subnet then he/she has to call the GOD, interface which talks to cloud.
func (v *GetNetworksInput) GetVpcs(con neuronaws.EstablishConnectionInput) (NetworkResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return NetworkResponse{}, seserr
	}
	response, err := ec2.DescribeAllVpc(
		&neuronaws.DescribeNetworkInput{
			VpcIds: v.VpcIds,
		},
	)
	if err != nil {
		return NetworkResponse{}, err
	}

	if v.GetRaw == true {
		return NetworkResponse{GetVpcsRaw: response}, nil
	}

	vpcs := make([]VpcResponse, 0)
	for _, vpc := range response.Vpcs {
		vpcs = append(vpcs, VpcResponse{Name: *vpc.Tags[0].Value, VpcId: *vpc.VpcId, State: *vpc.State, IsDefault: *vpc.IsDefault})
	}
	return NetworkResponse{Vpcs: vpcs}, nil
}

// FindVpcs is a customized method which sends back the response to the caller about the existence of vpc asked for.
func (v *GetNetworksInput) FindVpcs(con neuronaws.EstablishConnectionInput) (bool, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return false, seserr
	}
	response, err := ec2.DescribeVpc(
		&neuronaws.DescribeNetworkInput{
			VpcIds: v.VpcIds,
		},
	)
	if err != nil {
		return false, err
	}
	if len(response.Vpcs) != 0 {
		return true, nil
	}
	return false, fmt.Errorf("Could not find the VPC's you asked for")
}
