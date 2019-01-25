// This package takes care of registering flags,subcommands and returns the
// command to the function who creates or holds the root command.
package commands

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	svcreate "neuron/cloudoperations/server/create"
	svdelete "neuron/cloudoperations/server/delete"
	svget "neuron/cloudoperations/server/get"
	svupdate "neuron/cloudoperations/server/update"
	err "neuron/error"
	"os"
	"strings"
)

var (
	servercmds map[string]*cobra.Command
	createsv   = svcreate.New()
	deletesv   = svdelete.New()
	updatesv   = svupdate.New()
	getsv      = svget.New()
)

// The function that helps in registering the subcommands with the respective main command.
// Make sure you call this, and this is the only way to register the subcommands.
func svRegister(name string, fn *cobra.Command) {
	if servercmds == nil {
		servercmds = make(map[string]*cobra.Command)
	}

	if servercmds[name] != nil {
		panic(fmt.Sprintf("Command %s is already registered", name))
	}
	servercmds[name] = fn
}

// The only way to create server command is to call this function and
// package commands will take care of calling this.
func getServCmds() *cobra.Command {

	// Creating "server" happens here.
	var cmdServer = &cobra.Command{
		Use:   "server [flags]",
		Short: "command to carry out server activities",
		Long:  `This will help user to create/update/get/delete server in cloud.`,
		RunE:  cc.echoServer,
	}
	registersvFlags("server", cmdServer)

	for cmdname, cmd := range servercmds {
		cmdServer.AddCommand(cmd)
		registersvFlags(cmdname, cmd)
	}
	return cmdServer
}

// Registering all the flags to the subcommands and command netwrok itself.
func registersvFlags(cmdname string, cmd *cobra.Command) {

	switch strings.ToLower(cmdname) {
	case "servercreate":
		cmd.Flags().StringVarP(&createsv.InstanceName, "name", "n", "", "give a name to your network.")
		cmd.Flags().Int64VarP(&createsv.Count, "count", "", 1, "specify the number of servers that has to be provisioned.")
		cmd.Flags().StringVarP(&createsv.ImageId, "imageid", "i", "", "ID of the base image from which the new server has to be provisioned.")
		cmd.Flags().StringVarP(&createsv.SubnetId, "subnetid", "s", "", "ID of the subnet in which servers has to be created.")
		cmd.Flags().StringVarP(&createsv.KeyName, "keyname", "k", "", "name of the kay-pair which has to be assigned to instances so it will be helpful while logging into it (works only with aws).")
		cmd.Flags().StringVarP(&createsv.Flavor, "flavor", "f", "", "flavor/configuration of the vm that has to be created. (checkout 'neuron flavor list' for the list of available flavors.)")
		cmd.Flags().StringVarP(&createsv.UserData, "userdata", "", "echo 'from neuron'", "if in case you need to execute certain scripts such as shell,ruby on the startup of server.? pass it from this flag.")
		cmd.Flags().BoolVarP(&createsv.AssignPubIp, "assignpublicip", "", false, "turnn this flag on if you need public ip for the machines which will be created.")
	case "serverdelete":
		cmd.Flags().StringVarP(&deletesv.VpcId, "vpcid", "v", "", "pass ID of vpc, from which servers has to be deleted")
	case "serverupdate":
		cmd.Flags().StringVarP(&updatesv.Action, "action", "", "", "action to be performed on the instances (supports start/stop).")
	case "serverget":
		cmd.Flags().StringSliceVarP(&getsv.VpcIds, "vpcids", "v", nil, "ID's of vpcs/vnets, pass comma separated value (if this flag is on which means you'll get servers in vpcs you mentioned)")
		cmd.Flags().StringSliceVarP(&getsv.SubnetIds, "subnetids", "s", nil, "ID's of subnets to filter the servers. pass comma separated value.")
	case "server":
		cmd.PersistentFlags().StringSliceVarP(&getsv.VpcIds, "instanceids", "", nil, "ID's of servers/instances, pass comma separated value.")
	}
}

func (cm *cliMeta) createServer(cmd *cobra.Command, args []string) error {
	if cm.CliSet == false {
		return err.CliNoStart()
	}
	createsv.Cloud = getCloud(cmd)
	createsv.Region = getRegion(cmd)
	createsv.Profile = getProfile(cmd)
	createsv.GetRaw = getGetRaw(cmd)
	server_response, ser_resp_err := createsv.CreateServer()
	if ser_resp_err != nil {
		return ser_resp_err
	} else {
		json_val, _ := json.MarshalIndent(server_response, "", " ")
		fmt.Fprintf(os.Stdout, "%v\n", string(json_val))
	}
	return nil
}

func (cm *cliMeta) deleteServer(cmd *cobra.Command, args []string) error {
	if cm.CliSet == false {
		return err.CliNoStart()
	}
	deletesv.Cloud = getCloud(cmd)
	deletesv.Region = getRegion(cmd)
	deletesv.Profile = getProfile(cmd)
	deletesv.GetRaw = getGetRaw(cmd)
	delete_sv_response, sv_err := deletesv.DeleteServer()
	if sv_err != nil {
		return sv_err
	} else {
		json_val, _ := json.MarshalIndent(delete_sv_response, "", " ")
		fmt.Fprintf(os.Stdout, "%v\n", string(json_val))
	}
	return nil
}

func (cm *cliMeta) getServer(cmd *cobra.Command, args []string) error {
	if cm.CliSet == false {
		return err.CliNoStart()
	}

	getsv.Cloud = getCloud(cmd)
	getsv.Region = getRegion(cmd)
	getsv.Profile = getProfile(cmd)
	getsv.GetRaw = getGetRaw(cmd)

	if isAll(cmd) {
		get_server_response, sv_get_err := getsv.GetAllServers()
		if sv_get_err != nil {
			return sv_get_err
		} else {
			json_val, _ := json.MarshalIndent(get_server_response, "", " ")
			fmt.Fprintf(os.Stdout, "%v\n", string(json_val))
		}
		return nil
	}

	get_server_response, sv_get_err := getsv.GetServersDetails()
	if sv_get_err != nil {
		return sv_get_err
	} else {
		json_val, _ := json.MarshalIndent(get_server_response, "", " ")
		fmt.Fprintf(os.Stdout, "%v\n", string(json_val))
	}
	return nil
}

func (cm *cliMeta) updateServer(cmd *cobra.Command, args []string) error {
	if cm.CliSet == false {
		return err.CliNoStart()
	}
	updatesv.Cloud = getCloud(cmd)
	updatesv.Region = getRegion(cmd)
	updatesv.Profile = getProfile(cmd)
	updatesv.GetRaw = getGetRaw(cmd)
	sv_update_response, sv_up_err := updatesv.UpdateServers()
	if sv_up_err != nil {
		return sv_up_err
	} else {
		json_val, _ := json.MarshalIndent(sv_update_response, "", " ")
		fmt.Fprintf(os.Stdout, "%v\n", string(json_val))
	}
	return nil
}

func (cm *cliMeta) echoServer(cmd *cobra.Command, args []string) error {
	if cm.CliSet == false {
		return err.CliNoStart()
	}
	printMessage()
	cmd.Usage()
	return nil
}