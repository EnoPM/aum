package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Release struct {
	Assets []Asset `json:"assets"`
	Id     int     `json:"id"`
	Path   string
}

type Asset struct {
	Name        string `json:"name"`
	DownloadUrl string `json:"browser_download_url"`
}

func GetLatestRelease(owner string, repo string) *Release {
	var url = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	reqClient := &http.Client{}
	res, err := reqClient.Do(request)
	if err != nil {
		panic(err)
	}
	var response Release
	defer closeBody(res.Body)
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		return nil
	}
	response.Path = fmt.Sprintf("%s/%s", owner, repo)
	return &response
}

func (release *Release) GetZipAsset() *Asset {
	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, ".zip") {
			return &asset
		}
	}
	return nil
}
