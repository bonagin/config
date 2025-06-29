# Config Library

Custom JSON config library fo go

---

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Note](#note)
- [License](#license)

---

## Features

- Configs stored in JSON format
- Supports any terminal editor.
- Config file is stored in the same directory as the exercutable.
- \*.conf files can be updated externally.
- **NEW**: Google Secret Manager integration as a third fallback option
- **NEW**: Three-tier configuration priority: JSON config → Environment variables → Google Secret Manager

## Requirements

- Supports Go 1.7 or later

---

## Installation

```bash
go get -u github.com/bonagin/config
```

Make sure [Git is installed](https://git-scm.com/downloads) on your machine and in your system's `PATH`.

## Usage

First create a config file in json format test.conf

```json
{
  "TEST": "Hello config"
}
```

```go
import (
  "github.com/bonagin/config"
)
// Initialize the config library using the application name as the argument
config.NewConfig("test")

// Call the Get method to read the config variable.
config_var := config.Get("HELLO")
```

Add the config flag and the path to the config file

```bash
$ /path/to/exercutable --config /path/to/configfile
```

### Google Secret Manager Integration (NEW)

The library now supports Google Secret Manager as a third fallback option. Configuration values are retrieved in the following priority order:

1. **JSON config file** (if provided via `--config` flag)
2. **Environment variables**
3. **Google Secret Manager** (if enabled)

#### Enabling Google Secret Manager

To enable Google Secret Manager integration, call `EnableGSecreteManager()` after initializing the config:

```go
import (
    "log"
    "os"
    "github.com/bonagin/config"
)

func main() {
    // Initialize the config system
    config.NewConfig()

    // Enable Google Secret Manager (optional)
    projectID := "your-google-cloud-project-id"
    err := config.EnableGSecreteManager(projectID)
    if err != nil {
        log.Printf("Failed to enable Google Secret Manager: %v", err)
        // Continue without GSM
    } else {
        // Ensure proper cleanup on exit
        defer config.CleanupGSecreteManager()
    }

    // Get configuration values - will check JSON config, env vars, then GSM
    dbURL := config.Get("DATABASE_URL")
    apiKey := config.Get("API_KEY")
}
```

#### Authentication

The library uses Google's default authentication methods. Ensure your application has access to Google Secret Manager by either:

- Running on Google Cloud (uses default service account)
- Setting `GOOGLE_APPLICATION_CREDENTIALS` environment variable to point to a service account key file
- Using `gcloud auth application-default login` for local development

#### Helper Functions

```go
// Check if Google Secret Manager is enabled
if config.IsGSecreteManagerEnabled() {
    log.Println("GSM is enabled")
}

// Get the current project ID
projectID := config.GetGSecreteProjectID()

// Disable Google Secret Manager
config.DisableGSecreteManager()

// Clean up resources (should be called before application exit)
config.CleanupGSecreteManager()
```

#### Secret Naming

Secrets in Google Secret Manager should have the same name as your configuration variable. For example:

- Configuration variable: `DATABASE_URL`
- Secret name in GSM: `DATABASE_URL`

### Note

The library only supports `string` variable ut will be changed in future to return an `interface`

## License

Go-MySQL-Driver is licensed under the [Mozilla Public License Version 2.0](https://raw.github.com/go-sql-driver/mysql/master/LICENSE)

Mozilla summarizes the license scope as follows:

> MPL: The copyleft applies to any files containing MPLed code.

That means:

- You can **use** the **unchanged** source code both in private and commercially.
- When distributing, you **must publish** the source code of any **changed files** licensed under the MPL 2.0 under a) the MPL 2.0 itself or b) a compatible license (e.g. GPL 3.0 or Apache License 2.0).
- You **needn't publish** the source code of your library as long as the files licensed under the MPL 2.0 are **unchanged**.

Please read the [MPL 2.0 FAQ](https://www.mozilla.org/en-US/MPL/2.0/FAQ/) if you have further questions regarding the license.

You can read the full terms here: [LICENSE](https://github.com/bonagin/config/blob/master/LICENSE).
