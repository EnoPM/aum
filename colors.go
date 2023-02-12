package main

import "github.com/fatih/color"

func RedBold(text string) string {
	c := color.New(color.Bold, color.FgRed)
	return c.Sprint(text)
}

func Green(text string) string {
	return color.GreenString(text)
}

func Red(text string) string {
	return color.RedString(text)
}

func Yellow(text string) string {
	return color.YellowString(text)
}

func BlueBold(text string) string {
	c := color.New(color.Bold, color.FgBlue)
	return c.Sprint(text)
}

func YellowItalic(text string) string {
	c := color.New(color.Italic, color.FgYellow)
	return c.Sprint(text)
}
