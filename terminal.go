package main

import (
	"fmt"
	"github.com/gosuri/uilive"
)

type Terminal struct {
	StaticText string
	writer     *uilive.Writer
}

func NewTerminal() *Terminal {
	return &Terminal{
		writer: uilive.New(),
	}
}

func (t *Terminal) String(text string) {
	_, err := fmt.Fprintf(t.writer, "%s %s", t.StaticText, text)
	if err != nil {
		panic(err)
	}
}

func (t *Terminal) Start() {
	t.writer.Start()
}

func (t *Terminal) Stop() {
	t.writer.Stop()
}
