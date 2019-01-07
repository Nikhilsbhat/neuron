// The package which makes the tool cloud agnostic with respect to creation of network.
// The decision will be made here to route the request to respective package based on input.
package networkCreate

// The struct which impliments method CreateNetwork.
type NetworkCreateInput struct {
	// The name for the Network that has to be created.
	Name string `json:"Name"`

	// The CIDR block which will be used to create VPC and this
	// contains info that how many IP should be present in the network
	// so decide that in prior before calling this.
	VpcCidr string `json:"VpcCidr"`

	// The CIDR for the subnet that has to be created in the VPC.
	// Pass an array of CIDR's and neuron will take care of creating
	// appropriate number of subnets and attaching to created VPC.
	SubCidr []string `json:"SubCidr"`

	// The type of the network that has to be created, public or private.
	// Accordingly IGW will be created and attached.
	Type string `json:"Type"`

	// The ports that has to be opened for the network,
	// if not passed, by default 22 will be made open so that
	// one can access machines that will be created inside the created network.
	Ports []string `json:"Ports"`

	// Pass the cloud in which the resource has to be created. usage: "aws","azure" etc.
	Cloud string `json:"Cloud"`

	// Along with cloud, pass region in which resource has to be created.
	Region string `json:"Region"`

	// Passing the profile is important, because this will help in fetching the the credentials
	// of cloud stored along with user details.
	Profile string `json:"Profile"`

	// Use this option if in case you need unfiltered output from cloud.
	GetRaw bool `json:"GetRaw"`
}

//Nothing much from this file. This file contains only the structs for network/create
