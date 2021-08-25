package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type NullValue struct {
	Value   string
	NotNull bool
}

var Config map[string]NullValue
var ConfigFileName NullValue
var Editor string

/*
	Initialize a new JSON config instance
	app - application name
*/
func NewConfig(app string) {
	args := os.Args
	print := false
	filename := "./" + app + ".conf"

	if (len(args) < 2) && (args[1] == "config") {
		if len(args) < 3 {
			log.Fatal("Error: missing flag" + app +
				"\n\t '" + args[0] + " config -h' for help")
		}

		print = configOptions(args[2], filename)

		if !print {
			os.Exit(0)
		}
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Config not set for applicatiion: '" + app +
			"\n\t '" + args[0] + " config -h' for help")
	}

	config, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Failed to read config file : ", err.Error())
	}

	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(config, &objmap)
	if err != nil {
		log.Fatal("Failed to process config file : ", err.Error())
	}

	if Config == nil {
		Config = make(map[string]NullValue)
	}

	for key, value := range objmap {
		val := string(*value)

		if len(val) > 0 && val[0] == '"' {
			val = val[1:]
		}
		if len(val) > 0 && val[len(val)-1] == '"' {
			val = val[:len(val)-1]
		}

		Config[key] = NullValue{Value: val, NotNull: true}
		if print {
			log.Println(key, "=", val)
		}
	}

	ConfigFileName = NullValue{Value: file.Name(), NotNull: true}

	log.Println("Config loaded")

	if print {
		os.Exit(0)
	}
}

/*
	Get a variable from JSON config instance
	variable - variable name
	return value
*/
func Get(variable string) string {
	if !Config[variable].NotNull {
		log.Fatal("'" + variable + "' not set in config file: '" + ConfigFileName.Value + "'")
	}

	return Config[variable].Value
}

func configOptions(flag, filename string) bool {

	switch flag {
	case "-e":
		if len(os.Args) < 4 {
			log.Fatal("\nEnter editor name eg. nano/vim")
		}

		Editor = os.Args[3]
		log.Println("\nEditor set to '" + Editor + "'")
	case "-w":
		if Editor == "" {
			Editor = "nano"
		}

		if !ConfigFileName.NotNull {
			cmd := exec.Command(Editor, filename)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Start()
			if err != nil {
				log.Fatal(err.Error())
			}
			err = cmd.Wait()
			if err != nil {
				log.Printf("Error while editing. Error: %v\n", err)
			} else {
				log.Printf("Successfully edited.")
			}
		}
	case "-r":
		return true
	default:
		log.Println("\nConfig Help: ")
		log.Println(" -e  set editor to use when editing the config file")
		log.Println(" -w  Write/Edit the config file")
		log.Println(" -r  Read the config")
	}

	return false
}
