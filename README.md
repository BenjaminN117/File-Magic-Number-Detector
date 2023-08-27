# MacOS File Magic Number Detection

![Static Badge](https://img.shields.io/badge/Golang-=>1.19.4-blue)
![Static Badge](https://img.shields.io/badge/License-BSD2-yellow)
![Static Badge](https://img.shields.io/badge/MacOS-%2399ffcc)

## Description

Observes the directory chosen for files that their advertised file types do not match their true file type based on the magic number. This then sends a notification to the user alerting them of the file. Files that can not be identified get recorded in a log file.

## Prerequisites

- Golang >= 1.19.4

## Dependancies

- github.com/gen2brain/beeep v0.0.0
- github.com/go-toast/toast v0.0.0
- github.com/godbus/dbus/v5 v5.1.0
- github.com/nu7hatch/gouuid v0.0.0
- github.com/tadvi/systray v0.0.0
- golang.org/x/exp v0.0.0
- golang.org/x/sys v0.11.0 

## Install

```
git clone https://github.com/BenjaminN117/File-Magic-Number-Detector.git
```

Build binaries

```
go build src/main.go
```

## Usage

```
-filepath string
        Please enter a target directory (default "./")
  -logger string
        Please specify the log file location (default "./")
```

Run Script
```
go run src/main.go -filepath ~/Downloads
```