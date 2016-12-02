package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"local"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

import "golang.org/x/net/html"

//HaxeVersion comment
type HaxeVersion struct {
	version string
}

//TagObject test
type TagObject struct {
	Sha string
	URL string
}

//TagInfo test
type TagInfo struct {
	Ref    string
	URL    string
	Object TagObject
}

func (h *HaxeVersion) set(path string) error {
	err := myUtils.RecursiveCopy(filepath.Join(settings.Path, h.version), path)
	if err != nil {
		return err
	}
	return nil
}

func (h *HaxeVersion) download() error {
	var versionURL string
	sha, err := getCommit(h.version)
	fmt.Println("Commit SHA " + sha)
	if err != nil {
		return err
	}
	hrefs, err := getHrefs()
	if err != nil {
		return err
	}
	for _, h := range hrefs {
		if strings.Index(h, sha) != (-1) {
			versionURL = h
		}
	}
	fmt.Println("Downloading from " + versionURL)
	resp, err := http.Get(versionURL)
	if err != nil {
		return err
	}
	ttgz := filepath.Join(settings.Path, "temp.tar.gz")
	file, err := os.Create(ttgz)
	if err != nil {
		return err
	}
	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}
	fmt.Println("GUnzip temp.tar.gz")
	tt := filepath.Join(settings.Path, "temp.tar")
	if err := ungzip(ttgz, tt); err != nil {
		fmt.Println("ungzip err")
		return err
	}
	file.Close()
	fmt.Println("Remove temp.tar.gz")
	if err := os.Remove(ttgz); err != nil {
		return err
	}
	fmt.Println("Untar temp.tar")
	dirName, err := untar(tt, settings.Path)
	if err != nil {
		return err
	}
	fmt.Println("Remove temp.tar")
	if err := os.Remove(tt); err != nil {
		return err
	}
	fmt.Println("Rename path")
	oldPath := filepath.Join(settings.Path, dirName)
	newPath := filepath.Join(settings.Path, h.version)
	p, _ := os.Stat(newPath)
	if p != nil {
		if err := os.RemoveAll(newPath); err != nil {
			return err
		}
	}
	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}
	return nil
}

func getHrefs() ([]string, error) {
	var osName string
	var hrefs []string
	switch runtime.GOOS {
	case "windows":
		osName = "windows"
	}
	haxeDownloadURL := "http://hxbuilds.s3-website-us-east-1.amazonaws.com/builds/haxe/" + osName + "/"
	resp, err := http.Get(haxeDownloadURL)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return hrefs, nil
		case html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if isAnchor {
				for _, a := range t.Attr {
					if a.Key == "href" {
						hrefs = append(hrefs, a.Val)
					}
				}
			}
		}

	}
}

func getCommit(version string) (sha string, err error) {
	haxeTagsURL := "https://api.github.com/repos/HaxeFoundation/haxe/git/refs/tags/"
	tagJSON, err := http.Get(haxeTagsURL + version)
	defer tagJSON.Body.Close()
	if err != nil {
		return "", err
	}
	tag := new(TagInfo)
	b, err := ioutil.ReadAll(tagJSON.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(b, &tag); err != nil {
		return "", err
	}
	sha = tag.Object.Sha[:7]
	return sha, nil
}

func ungzip(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()
	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()
	target = filepath.Join(target, archive.Name)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = io.Copy(writer, archive)
	return nil
}

func untar(tarball, target string) (string, error) {
	var basePath = ""
	reader, err := os.Open(tarball)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if basePath == "" {
				basePath = info.Name()
			}
			if err := os.MkdirAll(path, info.Mode()); err != nil {
				return "", err
			}
			continue
		}
		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return "", err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return "", err
		}

	}
	return basePath, nil
}
