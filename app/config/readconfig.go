package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	neuerr "github.com/nikhilsbhat/neuron/error"
	log "github.com/nikhilsbhat/neuron/logger"
	"io/ioutil"
	"os"
)

func findConfig(pathtofile string) (AppConfig, error) {

	if pathtofile != "" {
		if _, dir_neuerr := os.Stat(pathtofile); os.IsNotExist(dir_neuerr) {
			return AppConfig{}, neuerr.NoFileFoundError()
		} else {
			config_data, confneuerr := readConfig(pathtofile)
			if confneuerr != nil {
				return AppConfig{}, confneuerr
			}
			return config_data, nil
		}
	}

	if _, dir_neuerr := os.Stat("/var/lib/neuron/neuron.json"); os.IsNotExist(dir_neuerr) {

		log.Info("+++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		log.Info("")
		log.Info("You did not provide the configuration file hence setting configurations to default")

		decoder := json.NewDecoder(bytes.NewReader([]byte(`{"port": "80","logfile": "neuron.log","loglocation": "/var/log/neuron/", "home": "/var/lib/neuron"}`)))
		var config_data AppConfig
		if decodneuerr := decoder.Decode(&config_data); decodneuerr != nil {
			log.Error(neuerr.JsonDecodeError())
			log.Error("Please provide us valid file")
			log.Error("Hence quitting installation...")
			return AppConfig{}, neuerr.JsonDecodeError()
		}
		return config_data, nil
	} else {

		log.Info("+++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		log.Info("")
		log.Info("Found configuration file, configuring application as per the entries")
		config_data, confneuerr := readConfig("/var/lib/neuron/neuron.json")
		if confneuerr != nil {
			return AppConfig{}, confneuerr
		}
		return config_data, nil
	}
}

func readConfig(pathtofile string) (AppConfig, error) {
	conf_file, confneuerr := ioutil.ReadFile(pathtofile)
	if confneuerr != nil {
		log.Error(neuerr.InvalidConfig())
		return AppConfig{}, confneuerr
	}

	var confdata AppConfig
	if decodneuerr := JsonDecode([]byte(conf_file), &confdata); decodneuerr != nil {
		log.Error(neuerr.JsonDecodeError())
		log.Error("Hence quitting installation...")
		return AppConfig{}, decodneuerr
	}
	return confdata, nil
}

// below functions are placed here temporarily, once respective package is created will be moved to it.
func JsonDecode(data []byte, i interface{}) error {
	err := json.Unmarshal(data, i)
	if err != nil {

		switch err.(type) {
		case *json.UnmarshalTypeError:
			return unknownTypeError(data, err)
		case *json.SyntaxError:
			return syntaxError(data, err)
		}
	}

	return nil
}

func syntaxError(data []byte, err error) error {
	syntaxErr, ok := err.(*json.SyntaxError)
	if !ok {
		return err
	}

	newline := []byte{'\x0a'}

	start := bytes.LastIndex(data[:syntaxErr.Offset], newline) + 1
	end := len(data)
	if index := bytes.Index(data[start:], newline); index >= 0 {
		end = start + index
	}

	line := bytes.Count(data[:start], newline) + 1

	err = fmt.Errorf("error occurred at line %d, %s\n%s",
		line, syntaxErr, data[start:end])
	return err
}

func unknownTypeError(data []byte, err error) error {
	unknownTypeErr, ok := err.(*json.UnmarshalTypeError)
	if !ok {
		return err
	}

	newline := []byte{'\x0a'}

	fmt.Println(bytes.LastIndex(data[:unknownTypeErr.Offset], newline))
	start := bytes.LastIndex(data[:unknownTypeErr.Offset], newline) + 1
	end := len(data)
	if index := bytes.Index(data[start:], newline); index >= 0 {
		end = start + index
	}

	line := bytes.Count(data[:start], newline) + 1

	err = fmt.Errorf("error occurred at line %d, %s\n%s\nThe data type you entered for the value is wrong",
		line, unknownTypeErr, data[start:end])
	return err
}
