package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type NullValue struct {
	Value   string
	NotNull bool
}

var (
	Config           map[string]NullValue
	ConfigFileName   NullValue
	ConfigFileParsed bool
)

/*
Initialize a new JSON config instance
*/
func NewConfig() {
	args := os.Args
	print := false
	filename := ""
	ConfigFileParsed = false

	for i, arg := range args {
		if arg == "--config" {
			if len(args) > (i + 1) {
				filename = args[i+1]
			}
		}
		if arg == "--print_config" {
			if len(args) > (i + 1) {
				print = true
			}
		}
	}

	if filename != "" {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal("Failed to open config file ", err.Error())
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

		ConfigFileParsed = true
	} else {
		log.Println("\nConfig tag not set\n\t '" + os.Args[0] + "' --config /path/to/json/configfile'")
	}

	// log.Println("Config loaded")

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
	if !ConfigFileParsed {
		return os.Getenv(variable)
	}

	if !Config[variable].NotNull {
		env := os.Getenv(variable)

		if env != "" {
			return env
		}

		log.Fatal("'" + variable + "' not set as an enviromental variable or in config file: '" + ConfigFileName.Value + "'")
	}

	return Config[variable].Value
}
