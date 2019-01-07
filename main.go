package main

import (
	log "neuron/logger"
	err "neuron/error"
	"neuron/app"
	cli "neuron/cli"
)

//This function is responsible for starting the application.
func main() {

	neuerr := neuron.StartNeuron()
	if neuerr != nil {
                log.Error(neuerr)
		log.Error(err.FailStartError())
	}
	cli.CliMain()
}
