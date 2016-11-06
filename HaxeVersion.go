package main

import "fmt"
import "net/http"
import "os"
import "io"
import "archive/zip"
import "path/filepath"

type HaxeVersion struct {
	version    string
	isDownload bool
	path       string
}

func (h *HaxeVersion) download() {
	haxeDownloadURL := "https://haxe.org/website-content/downloads/" + h.version + "/downloads" + "/haxe-" + h.version + "-win.zip"
	resp, getErr := http.Get(haxeDownloadURL)
	if getErr != nil {
		fmt.Println(getErr)
	}
	dirErr := os.MkdirAll(h.path, 0755)
	if dirErr != nil {
		fmt.Println(dirErr)
	}
	file, crtErr := os.Create(h.path + "temp.zip")
	if crtErr != nil {
		fmt.Println(crtErr)
	}
	io.Copy(file, resp.Body)
	resp.Body.Close()
	file.Close()
	unzip(h.path, "temp.zip")
    if err := os.Remove(h.path + "temp.zip") ; err != nil {
        fmt.Println(err)
    }
    
}

func unzip(path, fileName string) {
	reader, err := zip.OpenReader(path + fileName)
	defer reader.Close()
    if err != nil {
		fmt.Println(err)
	}
	for _, file := range reader.File {
		newPath := filepath.Join(path, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(newPath, file.Mode())
			continue
		}
		fileReader, err := file.Open()
		if err != nil {
			fmt.Println("readerErr")
            fmt.Println(err)
		}
		defer fileReader.Close()
		targetFile, err := os.OpenFile(newPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fmt.Println("openErr")
            fmt.Println(err)
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			fmt.Println("copyErr")
            fmt.Println(err)
		}
	}
}
