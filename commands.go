package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Command struct {
	Name         string
	Aliases      []string
	Arguments    []CommandArgument
	AllowedFlags []string
	InputArgs    []string
	InputFlags   []string
	ArgsValues   []ArgsValue
}

type CommandArgument struct {
	Name     string
	Regex    *regexp.Regexp
	Required bool
}

type ArgsValue struct {
	Name     string
	Value    string
	HasError bool
}

func NewCommand(name string, aliases ...string) *Command {
	return &Command{
		Name:       name,
		Aliases:    aliases,
		Arguments:  make([]CommandArgument, 0),
		InputArgs:  make([]string, 0),
		InputFlags: make([]string, 0),
	}
}

func (a *CommandArgument) Stringify(hasError bool) string {
	var name string
	var startChar string
	var endChar string
	if hasError {
		name = Red(a.Name)
	} else {
		name = Green(a.Name)
	}
	if a.Required {
		startChar = Yellow("[")
		endChar = Yellow("]")
	} else {
		startChar = Yellow("(")
		endChar = Yellow(")")
	}
	return fmt.Sprintf("%s%s%s", startChar, name, endChar)
}

func (c *Command) AddArgument(name string, required bool, regex *regexp.Regexp) {
	c.Arguments = append(c.Arguments, CommandArgument{
		Name:     name,
		Regex:    regex,
		Required: required,
	})
}

func (c *Command) AddAllowedFlag(name string) {
	c.AllowedFlags = append(c.AllowedFlags, name)
}

func (c *Command) ParseInput() {
	for index, arg := range os.Args {
		if index < 2 {
			continue
		}
		if strings.HasPrefix(arg, "-") {
			c.InputFlags = append(c.InputFlags, strings.Replace(arg, "-", "", 1))
		} else if strings.HasPrefix(arg, "--") {
			c.InputFlags = append(c.InputFlags, strings.Replace(arg, "--", "", 1))
		} else {
			c.InputArgs = append(c.InputArgs, arg)
		}
	}
}

func (c *Command) ParseArgs() error {
	if len(c.InputArgs) > len(c.Arguments) {
		return errors.New("too many arguments")
	}
	for index, argument := range c.Arguments {
		value := ""
		if len(c.InputArgs) >= index+1 {
			value = c.InputArgs[index]
		}
		var hasError = false
		if argument.Regex != nil {
			hasError = !argument.Regex.MatchString(value)
		}
		c.ArgsValues = append(c.ArgsValues, ArgsValue{
			Name:     argument.Name,
			Value:    value,
			HasError: hasError,
		})
	}
	hasError := false
	for _, arg := range c.ArgsValues {
		hasError = arg.HasError
	}
	if hasError {
		errorMessage := fmt.Sprintf("%s $ %s %s ", RedBold("Expected"), Yellow("aum"), Yellow(c.Name))
		errorMessage2 := fmt.Sprintf("%s $ %s %s ", RedBold("             Got"), Yellow("aum"), Yellow(c.Name))
		for index, arg := range c.ArgsValues {
			argument := c.Arguments[index]
			errorMessage += fmt.Sprintf("%s ", argument.Stringify(arg.HasError))
			if arg.HasError {
				errorMessage2 += fmt.Sprintf("%s ", Red(arg.Value))
			} else {
				errorMessage2 += fmt.Sprintf("%s ", Green(arg.Value))
			}
		}
		return errors.New(fmt.Sprintf("%s\n%s", errorMessage, errorMessage2))
	}
	return nil
}

func (c *Command) HasFlag(name string) bool {
	for _, flag := range c.InputFlags {
		if flag == name {
			return true
		}
	}
	return false
}

func (c *Command) GetArgValue(name string) string {
	for _, arg := range c.ArgsValues {
		if arg.Name == name {
			return arg.Value
		}
	}
	return ""
}

func (c *Command) Match() bool {
	if os.Args[1] == c.Name {
		return true
	}
	for _, alias := range c.Aliases {
		if os.Args[1] == alias {
			return true
		}
	}
	return false
}

func (c *Command) ShowError(message string) {
	_, _ = fmt.Printf("%s %s", RedBold("[ERROR]"), message)
}

func (c *Command) ToHelp() string {
	help := fmt.Sprintf("%s", BlueBold(c.Name))
	for _, argument := range c.Arguments {
		help += fmt.Sprintf(" %s", argument.Stringify(false))
	}
	return help
}

func IsCommand() bool {
	return len(os.Args) > 1
}
