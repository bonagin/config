package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type NullValue struct {
	Value   string
	NotNull bool
}

var (
	Config                 map[string]NullValue
	ConfigFileName         NullValue
	ConfigFileParsed       bool
	gSecreteManagerEnabled bool
	gSecreteClient         *secretmanager.Client
	gSecreteProjectID      string
)

/*
EnableGSecreteManager enables Google Secret Manager as a third configuration option
projectID - Google Cloud Project ID where secrets are stored
*/
func EnableGSecreteManager(projectID string) error {
	if projectID == "" {
		return fmt.Errorf("projectID cannot be empty")
	}

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secret manager client: %v", err)
	}

	gSecreteClient = client
	gSecreteProjectID = projectID
	gSecreteManagerEnabled = true

	log.Println("Google Secret Manager enabled for project:", projectID)
	return nil
}

/*
DisableGSecreteManager disables Google Secret Manager integration
*/
func DisableGSecreteManager() {
	if gSecreteClient != nil {
		gSecreteClient.Close()
		gSecreteClient = nil
	}
	gSecreteManagerEnabled = false
	gSecreteProjectID = ""
	log.Println("Google Secret Manager disabled")
}

/*
getSecretFromGSM retrieves a secret from Google Secret Manager
secretName - name of the secret to retrieve
returns the secret value or empty string if not found
*/
func getSecretFromGSM(secretName string) string {
	if !gSecreteManagerEnabled || gSecreteClient == nil {
		return ""
	}

	ctx := context.Background()
	secretPath := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", gSecreteProjectID, secretName)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretPath,
	}

	result, err := gSecreteClient.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Printf("Failed to access secret %s: %v", secretName, err)
		return ""
	}

	return string(result.Payload.Data)
}

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
Priority order: 1. JSON config file, 2. Environment variables, 3. Google Secret Manager (if enabled)
*/
func Get(variable string) string {
	// First: Check if config file was parsed and variable exists in config
	if ConfigFileParsed && Config[variable].NotNull {
		return Config[variable].Value
	}

	// Second: Check environment variables
	env := os.Getenv(variable)
	if env != "" {
		return env
	}

	// Third: Check Google Secret Manager (if enabled)
	if gSecreteManagerEnabled {
		secret := getSecretFromGSM(variable)
		if secret != "" {
			return secret
		}
	}

	// If config file was parsed, show which file was checked
	if ConfigFileParsed {
		if gSecreteManagerEnabled {
			log.Fatal("'" + variable + "' not set as an environmental variable, in config file: '" + ConfigFileName.Value + "', or in Google Secret Manager")
		} else {
			log.Fatal("'" + variable + "' not set as an environmental variable or in config file: '" + ConfigFileName.Value + "'")
		}
	} else {
		// Config file not parsed, only checked env and possibly GSM
		if gSecreteManagerEnabled {
			log.Fatal("'" + variable + "' not set as an environmental variable or in Google Secret Manager")
		} else {
			log.Fatal("'" + variable + "' not set as an environmental variable")
		}
	}

	return ""
}

/*
IsGSecreteManagerEnabled returns whether Google Secret Manager is currently enabled
*/
func IsGSecreteManagerEnabled() bool {
	return gSecreteManagerEnabled
}

/*
GetGSecreteProjectID returns the current Google Cloud Project ID used for Secret Manager
*/
func GetGSecreteProjectID() string {
	return gSecreteProjectID
}

/*
CleanupGSecreteManager should be called before application exit to properly close the Secret Manager client
*/
func CleanupGSecreteManager() {
	if gSecreteClient != nil {
		gSecreteClient.Close()
		gSecreteClient = nil
	}
}
