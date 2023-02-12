package main

import (
	"archive/zip"
	"io"
	"os"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func closeBody(Body io.ReadCloser) {
	err := Body.Close()
	if err != nil {
		panic(err)
	}
}

func closeFile(out *os.File) {
	err := out.Close()
	if err != nil {
		panic(err)
	}
}

func closeZipReader(reader *zip.ReadCloser) {
	err := reader.Close()
	if err != nil {
		panic(err)
	}
}
