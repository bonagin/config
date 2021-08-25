# Credits API

Credits system for the [bulksms_client_django](https://github.com/dakab1/bulksms_client_django) system

---------------------------------------
  * [Features](#features)
  * [Requirements](#requirements)
  * [Installation](#installation)
  * [Usage](#usage)
  * [Note](#note)
  * [License](#license)

---------------------------------------

## Features
  * Configs stored in JSON format
  * Supports any terminal editor.
  * Config file is stored in the same directory as the exercutable. 
  * *.conf files can be updated externally.

## Requirements
  * Supports any Go version

---------------------------------------

## Installation
```bash
$ go get -u github.com/go-sql-driver/mysql
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
fmt.Println(config.Get("HELLO"))
```

### Note
The library only supports `string` variable ut will be changed in future to return an `interface`

## License
Go-MySQL-Driver is licensed under the [Mozilla Public License Version 2.0](https://raw.github.com/go-sql-driver/mysql/master/LICENSE)

Mozilla summarizes the license scope as follows:
> MPL: The copyleft applies to any files containing MPLed code.


That means:
  * You can **use** the **unchanged** source code both in private and commercially.
  * When distributing, you **must publish** the source code of any **changed files** licensed under the MPL 2.0 under a) the MPL 2.0 itself or b) a compatible license (e.g. GPL 3.0 or Apache License 2.0).
  * You **needn't publish** the source code of your library as long as the files licensed under the MPL 2.0 are **unchanged**.

Please read the [MPL 2.0 FAQ](https://www.mozilla.org/en-US/MPL/2.0/FAQ/) if you have further questions regarding the license.

You can read the full terms here: [LICENSE](https://github.com/bonagin/config/blob/master/LICENSE).
