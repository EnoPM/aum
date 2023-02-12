package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
	"path"
)

func DownloadFile(url string, dest string) {
	fileName := path.Base(url)
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer closeBody(resp.Body)
	f, _ := os.OpenFile(fmt.Sprintf("%s/%s", dest, fileName), os.O_CREATE|os.O_WRONLY, 0644)
	defer closeFile(f)

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		fileName,
	)
	_, err := io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil {
		panic(err)
	}
}
