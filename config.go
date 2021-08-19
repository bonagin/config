package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var Config map[string]string

/*
	Initialize a new JSON config instance
	app - application name
*/
func NewConfig(app string) {
	//env := os.Getenv("CONFIG")
	xmlfile, _ := os.Open("/usr/config/" + app + "/config")
	config, err := ioutil.ReadAll(xmlfile)

	if err != nil {
		panic(("Failed to load config : " + err.Error()))
	}

	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(config, &objmap)

	if Config == nil {
		Config = make(map[string]string)
	}

	for key, value := range objmap {
		val := string(*value)

		if len(val) > 0 && val[0] == '"' {
			val = val[1:]
		}
		if len(val) > 0 && val[len(val)-1] == '"' {
			val = val[:len(val)-1]
		}

		Config[key] = val
		//log.Printf("%s:%s\n", key, val)
	}

	log.Println("Config loaded")
}

/*
	Get a variable from JSON config instance
	variable - variable name
	return value
*/
func Get(variable string) string {
	return Config[variable]
}
