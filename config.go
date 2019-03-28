package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var Config map[string]string

func NewConfig() {
	//env := os.Getenv("CONFIG")
	jsondoc, _ := os.Open("/usr/config/config")
	config, err := ioutil.ReadAll(jsondoc)

	if err != nil {
		panic(("Failed to load config new : " + err.Error()))
	}

	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(config, &objmap)

	if Config == nil {
		Config = make(map[string]string)
	}

	log.Println("\n========================Config============================:\n")
	for key, value := range objmap {
		val := string(*value)

		if len(val) > 0 && val[0] == '"' {
			val = val[1:]
		}
		if len(val) > 0 && val[len(val)-1] == '"' {
			val = val[:len(val)-1]
		}

		Config[key] = val
		log.Printf("%s:%s\n", key, val)
	}
	log.Println("\n==========================END==============================:\n")
}
