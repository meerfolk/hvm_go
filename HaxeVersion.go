package main

import "net/http"
import "os"
import "io"

import "compress/gzip"
import "path/filepath"
import "io/ioutil"
import "encoding/json"

import "runtime"
import "golang.org/x/net/html"
import "strings"
import "fmt"
import "archive/tar"

//import "compress/gzip"

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

func (h *HaxeVersion) download() error {
	var versionURL string
	sha, err := getCommit(h.version)
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
	resp, err := http.Get(versionURL)
	if err != nil {
		return err
	}

	file, err := os.Create(settings.Path + "temp.tar.gz")
	if err != nil {
		return err
	}
	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}
	if err := ungzip(settings.Path+"temp.tar.gz", settings.Path+"temp.tar"); err != nil {
		fmt.Println("ungzip err")
		return err
	}
	if err := untar(settings.Path+"temp.tar", settings.Path); err != nil {
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

func untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err := os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}
		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}

	}
	return nil
}
