package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Settings struct {
	Path string
}

var stdPath string = os.Getenv("HOME") + "/.hvm/"
var params []string
var hvmPath string
var settings *Settings

func main() {
	params = os.Args[1:]
	settings = settingsResolver()
	key := ""
	val := ""
	if len(params) == 1 {
		key = params[0]
	}
	if len(params) == 2 {
		val = params[1]
	}
	paramsResolver(key, val)
}

func settingsResolver() (settings *Settings) {
	b, err := ioutil.ReadFile("./conf.json")
	if err != nil {
		fmt.Println("open file error : " + err.Error())
		newSettings := new(Settings)
		newSettings.Path = stdPath
		return newSettings
	}
	if err := json.Unmarshal(b, &settings); err != nil {
		fmt.Println("parse file error : " + err.Error())
	}
	return
}

func paramsResolver(key, val string) {
	switch key {
	case "download":
	case "setup":
		setup(val)
	}
}

func setup(path string) error {
	if path == "" {
		fmt.Println("current path : " + settings.Path)
		var newPath string
		fmt.Scanln(&newPath)
		if newPath == "" {
			path = settings.Path
		} else {
			path = newPath
		}
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	settings.Path = path
	b, err := json.Marshal(settings)
	fmt.Println("string settings " + string(b))
	if err != nil {
		fmt.Println("json error " + err.Error())
	}
	if err := ioutil.WriteFile("./conf.json", b, 0755); err != nil {
		fmt.Println("write file error : " + err.Error())
	}
	return nil
}
