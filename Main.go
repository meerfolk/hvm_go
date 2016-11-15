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
		key = params[0]
		val = params[1]
	}
	paramsResolver(key, val)
}

func settingsResolver() (settings *Settings) {
	b, err := ioutil.ReadFile("./conf.json")
	if err != nil {
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
		if err := download(val); err != nil {
			panic(err)
		}
	case "setup":
		if err := setup(val); err != nil {
			panic(err)
		}
	case "list":
		if err := getList(); err != nil {
			panic(err)
		}
	}
}

func getList() error {
	files, err := ioutil.ReadDir(settings.Path)
	if err != nil {
		return err
	}
	fmt.Println("version's list :")
	for _, file := range files {
		fmt.Println(file.Name())
	}
	return nil
}

func download(version string) error {
	fmt.Println("download version " + version)
	hx := HaxeVersion{version}
	if err := hx.download(); err != nil {
		return err
	}
	return nil
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
		return err
	}
	if err := ioutil.WriteFile("./conf.json", b, 0755); err != nil {
		return err
	}
	return nil
}
