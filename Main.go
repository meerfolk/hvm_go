package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Settings struct {
	IsSetup  bool
	Path     string
	HaxePath string
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
		newSettings.IsSetup = false
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
		if err := showList(); err != nil {
			panic(err)
		}
	case "set":
		if err := set(val); err != nil {
			panic(err)
		}
	}
}

func showList() error {
	fmt.Println("version's list :")
	files, err := getList()
	if err != nil {
		return err
	}
	for _, file := range files {
		fmt.Println(file)
	}
	return nil
}

func set(version string) error {
	if version == "" {
		fmt.Println("You didn't write the version")
		return nil
	}
	isDownloaded, err := isVersionDownloaded(version)
	if err != nil {
		return err
	}
	if isDownloaded {
		setVersion(version)
	}
	return nil
}

func setVersion(version string) error {
	haxePath := os.Getenv("HAXEPATH")
	if haxePath == "" {
		// fmt.Println("Please give haxe path")
		// fmt.Scanln(&haxePath)
		// err := os.MkdirAll(haxePath, 0755)
		// if err != nil {
		// 	return err
		// }
		// os.Setenv("HAXEPATH", haxePath)
		fmt.Println(os.Getenv("PATH"))
		//os.Setenv("PATH",os.Getenv("PATH") + + haxePath)
	}
	settings.HaxePath = haxePath
	hx := HaxeVersion{version}
	hx.set(haxePath)
	return nil
}

func isVersionDownloaded(version string) (bool, error) {
	files, err := getList()
	if err != nil {
		return false, err
	}
	for _, file := range files {
		if file == version {
			return true, nil
		}
	}
	return false, nil
}

func getList() ([]string, error) {
	var fileList []string
	files, err := ioutil.ReadDir(settings.Path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		fileList = append(fileList, file.Name())
	}
	return fileList, nil
}

func download(version string) error {
	fmt.Println("Download version " + version)
	if version == "" {
		fmt.Println("You didn't write the version")
		return nil
	}
	if !settings.IsSetup {
		setup("")
	}
	hx := HaxeVersion{version}
	if err := hx.download(); err != nil {
		return err
	}
	return nil
}

func setup(path string) error {
	if path == "" {
		fmt.Println("Default path for Haxe versions : " + settings.Path)
		var newPath string
		fmt.Println("Enter new path for Haxe versions :")
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
	settings.IsSetup = true
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
