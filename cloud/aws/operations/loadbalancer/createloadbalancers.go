package loadbalancer

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	aws "neuron/cloud/aws/interface"
	network "neuron/cloud/aws/operations/network"
	"strings"
)

type LoadBalanceCreateInput struct {
	//optional parameter; If you provide the name to the loadbalancer well and good, else we will name it with a default one.
	Name string `json:"Name"`

	//optional parameter; The Id of vpc in which the loadbalancer has to be created. Use this only if you don't want to pass subnets directly.
	//once this option is used we automatically fetch the random subnets from this network.
	VpcId string `json:"VpcId"`

	//optional parameter; The Ids of subnet in which the loadbalancer has to be created.
	SubnetIds []string `json:"SubnetIds"`

	//optional parameter; The names of availability zones to which loadbalancers has to be tagged. Either this or subnets has to be passed, passing both does't work
	AvailabilityZones []string `json:"AvailabilityZones"`

	//optional parameter; The Ids of secutiry group to be attached to loadbalancer.
	//If not mentioned, default security group of VPC will be attached.
	SecurityGroupIds []string `json:"SubnetIds"`

	//optional parameter; This field is to select the catageory of loadbalancer ex: internal, internet-facing. If not mentioned internet-facing will be created by default.
	Scheme string `json:"Scheme"`
	//mandatory parameter; The type of loadbalancer required ex: classic, application etc.
	Type string `json:"Type"`

	//required only if the LB protocol is HTTPS else can be initiazed with dummy value
	SslCert   string `json:"SslCert"`
	SslPolicy string `json:"SslPolicy"`

	//mandatory parameter; The port of the loabalancer. ex: 8080, 80 etc.
	LbPort   int64 `json:"LbPort"`
	InstPort int64 `json:"InstPort"`

	//mandatory parameter; The protocol of loadbalancer. ex: HTTPS, HTTP.
	Lbproto   string `json:"Lbproto"`
	Instproto string `json:"Instproto"`

	//optional parameter; The http code. ex: 200, 404 etc.
	HttpCode   string `json:"HttpCode"`
	HealthPath string `json:"HealthPath"`

	//optional parameter; Select Ip address type ex: ipv4, ipv6. If nothing is passed ipv4 is considered by default.
	IpAddressType string `json:"IpAddressType"`

	//optional parameter; Only when you need unfiltered result from cloud, enable this field by setting it to true. By default it is set to false.
	GetRaw bool `json:"GetRaw"`
}

type LoadBalanceResponse struct {
	Name                   string                           `json:"Name,omitempty""`
	Type                   string                           `json:"Type,omitempty""`
	LbDns                  string                           `json:"LbDns,omitempty""`
	LbArn                  string                           `json:"LbArn,omitempty"`
	LbArns                 []string                         `json:"LbArns,omitempty"`
	TargetArn              interface{}                      `json:"TargetArn,omitempty"`
	ListnerArn             interface{}                      `json:"ListnerArn,omitempty"`
	Createdon              string                           `json:"Createdon,omitempty"`
	VpcId                  string                           `json:"VpcId,omitempty"`
	Scheme                 string                           `json:"Scheme,omitempty"`
	DefaultResponse        interface{}                      `json:"DefaultResponse,omitempty"`
	LbDeleteStatus         string                           `json:"LbDeleteStatus,omitempty"`
	ClassicLb              []LoadBalanceResponse            `json:"ClassicLb,omitempty"`
	ApplicationLb          []LoadBalanceResponse            `json:"ApplicationLb,omitempty"`
	CreateClassicLbRaw     *elb.CreateLoadBalancerOutput    `json:"CreateClassicLbRaw,omitempty"`
	GetClassicLbsRaw       *elb.DescribeLoadBalancersOutput `json:"GetClassicLbsRaw,omitempty"`
	GetClassicLbRaw        *elb.LoadBalancerDescription     `json:"GetClassicLbRaw,omitempty"`
	CreateApplicationLbRaw ApplicationLbRaw                 `json:"CreateApplicationLbRaw,omitempty"`
	GetApplicationLbRaw    ApplicationLbRaw                 `json:"GetApplicationLbRaw,omitempty"`
}

type ApplicationLbRaw struct {
	CreateApplicationLbRaw *elbv2.CreateLoadBalancerOutput    `json:"CreateApplicationLbRaw,omitempty"`
	GetApplicationLbsRaw   *elbv2.DescribeLoadBalancersOutput `json:"GetApplicationLbsRaw,omitempty"`
	GetApplicationLbRaw    *elbv2.LoadBalancer                `json:"GetApplicationLbRaw,omitempty"`
	CreateTargetGroupRaw   *elbv2.CreateTargetGroupOutput     `json:"CreateTargetGroupRaw,omitempty"`
	GetTargetGroupRaw      *elbv2.DescribeTargetGroupsOutput  `json:"GetTargetGroupRaw,omitempty"`
	CreateListnersRaw      *elbv2.CreateListenerOutput        `json:"CreateListnersRaw,omitempty"`
	GetListnersRaw         *elbv2.DescribeListenersOutput     `json:"GetListnersRaw,omitempty"`
}

func (load *LoadBalanceCreateInput) CreateLoadBalancer(con aws.EstablishConnectionInput) (LoadBalanceResponse, error) {

	// creating LB according to the input ex: application/classic
	//get the relative sessions before proceeding further
	elb, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return LoadBalanceResponse{}, sesserr
	}

	lbin := new(aws.LoadBalanceCreateInput)
	//giving name to the loadbalancer which wil be created
	lbin.Name = load.Name
	// collecting subnet details
	if load.SubnetIds != nil {
		lbin.Subnets = load.SubnetIds
	} else {
		subnets_in := network.GetNetworksInput{VpcIds: []string{load.VpcId}}
		subnets_result, suberr := subnets_in.GetSubnetsFromVpc(con)
		if suberr != nil {
			return LoadBalanceResponse{}, suberr
		}
		for _, subnet := range subnets_result.Subnets {
			lbin.Subnets = append(lbin.Subnets, subnet.Id)
		}
	}

	// fetching security group so that I can attach it to the loabalancer which I am about to create
	if load.SecurityGroupIds != nil {
		lbin.SecurityGroups = load.SecurityGroupIds
	} else {
		sec_input := network.NetworkComponentInput{VpcIds: []string{load.VpcId}}
		sec_result, err := sec_input.GetSecFromVpc(con)
		if err != nil {
			return LoadBalanceResponse{}, err
		}
		lbin.SecurityGroups = sec_result.SecGroupIds
	}
	// creating load balancer

	// selecting scheme
	if load.Scheme == "external" {
		lbin.Scheme = "internet-facing"
	} else if load.Scheme == "internal" {
		lbin.Scheme = "internal"
	} else {
		lbin.Scheme = "internet-facing"
	}

	//setting availability zones
	if load.AvailabilityZones != nil {
		lbin.AvailabilityZones = load.AvailabilityZones
	}

	switch strings.ToLower(load.Type) {
	case "classic":

		lbin.InstPort = load.InstPort
		lbin.Instproto = load.Instproto
		lbin.LbPort = load.LbPort
		lbin.Lbproto = load.Lbproto
		lbin.SslCert = load.SslCert
		lb_create_response, err := elb.CreateClassicLb(*lbin)

		if err != nil {
			return LoadBalanceResponse{}, err
		}

		response := new(LoadBalanceResponse)
		if load.GetRaw == true {
			response.CreateClassicLbRaw = lb_create_response
			return *response, nil
		}

		response.Name = load.Name
		response.Type = load.Type
		response.LbDns = *lb_create_response.DNSName
		return *response, nil

	case "application":

		if load.IpAddressType == "" {
			lbin.IpAddressType = "ipv4"
		} else {
			lbin.IpAddressType = load.IpAddressType
		}
		// creating load balancer logic
		lb_create_response, lberr := elb.CreateApplicationLb(*lbin)
		if lberr != nil {
			return LoadBalanceResponse{}, lberr
		}

		lbin.Name = load.Name + "-target"
		lbin.LbPort = load.LbPort
		lbin.Lbproto = load.Lbproto
		lbin.VpcId = load.VpcId
		lbin.Instproto = load.Instproto
		lbin.InstPort = load.InstPort
		lbin.HealthPath = load.HealthPath
		lbin.HttpCode = load.HttpCode
		// creating target group
		target_group_response, tarerr := elb.CreateTargetGroups(lbin)
		if tarerr != nil {
			return LoadBalanceResponse{}, tarerr
		}

		lbin.SslCert = load.SslCert
		lbin.TargetArn = *target_group_response.TargetGroups[0].TargetGroupArn
		lbin.LbArn = *lb_create_response.LoadBalancers[0].LoadBalancerArn
		lbin.LbPort = load.LbPort
		lbin.Lbproto = load.Lbproto
		lbin.SslPolicy = load.SslPolicy
		listner_create_response, liserr := elb.CreateApplicationListners(lbin)
		if liserr != nil {
			return LoadBalanceResponse{}, liserr
		}

		response := new(LoadBalanceResponse)

		if load.GetRaw == true {
			response.CreateApplicationLbRaw.CreateApplicationLbRaw = lb_create_response
			response.CreateApplicationLbRaw.CreateTargetGroupRaw = target_group_response
			response.CreateApplicationLbRaw.CreateListnersRaw = listner_create_response
			return *response, nil
		}

		response.Name = load.Name
		response.Type = load.Type
		response.LbDns = *lb_create_response.LoadBalancers[0].DNSName
		response.LbArn = *lb_create_response.LoadBalancers[0].LoadBalancerArn
		response.TargetArn = *target_group_response.TargetGroups[0].TargetGroupArn
		response.ListnerArn = *listner_create_response.Listeners[0].ListenerArn
		return *response, nil

	default:
		return LoadBalanceResponse{}, fmt.Errorf("You provided unknown loadbalancer type, enter a valid LB type")
	}
}