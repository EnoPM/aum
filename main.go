package main

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"regexp"
	"strings"
)

var AppdataFolderPath string
var AppFolderPath string
var ConfigFilePath string
var DownloadFolderPath string
var ModsFolderPath string

func init() {
	var DS = string(os.PathSeparator)
	var err error
	AppdataFolderPath, err = os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	AppFolderPath = fmt.Sprintf("%s%saum", AppdataFolderPath, DS)
	mkdirIfNoExists(AppFolderPath)
	ConfigFilePath = fmt.Sprintf("%s%sconfig.json", AppFolderPath, DS)
	DownloadFolderPath = fmt.Sprintf("%s%sdownloads", AppFolderPath, DS)
	mkdirIfNoExists(DownloadFolderPath)
	ModsFolderPath = fmt.Sprintf("%s%smods", AppFolderPath, DS)
	mkdirIfNoExists(ModsFolderPath)

	initCommands()
}

var HelpCommand = NewCommand("help", "h")
var InstallCommand = NewCommand("install", "i", "get")
var UpdateCommand = NewCommand("update", "u")
var RemoveCommand = NewCommand("remove", "r", "uninstall", "delete")
var VanillaCommand = NewCommand("vanilla", "clear", "prune")
var ResetCommand = NewCommand("reset")
var ConfigCommand = NewCommand("config", "c")

func initCommands() {
	InstallCommand.AddArgument("owner/repository", true, regexp.MustCompile("(?i)^([a-z\\d](?:[a-z\\d]|-((.+?)?:[a-z\\d])){0,38})\\/.+$"))
	InstallCommand.AddAllowedFlag("force")

	UpdateCommand.AddArgument("owner/repository", true, regexp.MustCompile("(?i)^([a-z\\d](?:[a-z\\d]|-((.+?)?:[a-z\\d])){0,38})\\/.+$"))
	UpdateCommand.AddAllowedFlag("force")

	RemoveCommand.AddArgument("owner/repository", true, regexp.MustCompile("(?i)^([a-z\\d](?:[a-z\\d]|-((.+?)?:[a-z\\d])){0,38})\\/.+$"))

	ConfigCommand.AddArgument("path", true, nil)
}

func main() {
	if !amAdmin() {
		fmt.Printf("%s %s", RedBold("[ERROR]"), Red("AmongUsMods require a terminal running as administrator. Please reopen your terminal with administrator permissions"))
		return
	}
	if HelpCommand.Match() || !IsCommand() {
		HelpCommand.ParseInput()
		err := HelpCommand.ParseArgs()
		if err != nil {
			HelpCommand.ShowError(err.Error())
			return
		}
		Help()
		return
	}
	if ConfigCommand.Match() {
		ConfigCommand.ParseInput()
		err := ConfigCommand.ParseArgs()
		if err != nil {
			ConfigCommand.ShowError(err.Error())
			return
		}
		SetGameFolderPath(ConfigCommand.GetArgValue("path"))
		return
	}

	if Mods.GameFolderPath == nil {
		color.Red("[ERROR] Aum needs to know the location of your Among Us game folder. You can fill in the path of this folder using the flag -gfp=\"Absolute Among Us Folder Path\"")
		return
	}

	if !dirExists(*Mods.GameFolderPath) {
		color.Red("[ERROR] Aum needs a valid Among Us game folder. You can fill in the path of this folder using the flag -gfp=\"Absolute Among Us Folder Path\"")
		return
	}

	if InstallCommand.Match() {
		InstallCommand.ParseInput()
		err := InstallCommand.ParseArgs()
		if err != nil {
			InstallCommand.ShowError(err.Error())
			return
		}
		parsed := strings.Split(InstallCommand.GetArgValue("owner/repository"), "/")
		owner := parsed[0]
		repo := parsed[1]
		mod := Mods.GetModByPath(InstallCommand.GetArgValue("owner/repository"))
		InstallMod(mod, owner, repo)
		return
	}

	if ResetCommand.Match() {
		ResetCommand.ParseInput()
		err := ResetCommand.ParseArgs()
		if err != nil {
			ResetCommand.ShowError(err.Error())
			return
		}
		Mods.Reset()
		return
	}

	if VanillaCommand.Match() {
		VanillaCommand.ParseInput()
		err := VanillaCommand.ParseArgs()
		if err != nil {
			VanillaCommand.ShowError(err.Error())
			return
		}
		RemoveCurrentMod(true)
		return
	}

	if UpdateCommand.Match() {
		UpdateCommand.ParseInput()
		err := UpdateCommand.ParseArgs()
		if err != nil {
			UpdateCommand.ShowError(err.Error())
			return
		}
		parsed := strings.Split(UpdateCommand.GetArgValue("owner/repository"), "/")
		owner := parsed[0]
		repo := parsed[1]
		mod := Mods.GetModByPath(UpdateCommand.GetArgValue("owner/repository"))
		UpdateMod(mod, owner, repo)
		return
	}

	if RemoveCommand.Match() {
		RemoveCommand.ParseInput()
		err := RemoveCommand.ParseArgs()
		if err != nil {
			RemoveCommand.ShowError(err.Error())
			return
		}
		parsed := strings.Split(RemoveCommand.GetArgValue("owner/repository"), "/")
		owner := parsed[0]
		repo := parsed[1]
		mod := Mods.GetModByPath(RemoveCommand.GetArgValue("owner/repository"))
		RemoveMod(mod, owner, repo)
		return
	}

	regex := regexp.MustCompile("(?i)^([a-z\\d](?:[a-z\\d]|-((.+?)?:[a-z\\d])){0,38})\\/.+$")
	if regex.MatchString(os.Args[1]) {
		ActivateMod(os.Args[1])
	} else {
		fmt.Printf("%s %s %s", RedBold("[ERROR]"), Red("Unexpected command"), Yellow(os.Args[1]))
	}
}

func Help() {
	response := fmt.Sprintf("%s %s%s%s         %s\n", BlueBold("$ aum"), Yellow("["), Green("owner/repository"), Yellow("]"), YellowItalic("Enable a installed Among Us mod in your game folder"))
	response += fmt.Sprintf("%s %s              %s\n", BlueBold("$ aum"), ConfigCommand.ToHelp(), YellowItalic("Setup your Among Us game folder path"))
	response += fmt.Sprintf("%s %s                       %s\n", BlueBold("$ aum"), HelpCommand.ToHelp(), YellowItalic("Display available commands syntax"))
	response += fmt.Sprintf("%s %s %s\n", BlueBold("$ aum"), InstallCommand.ToHelp(), YellowItalic("Install a Among Us mod from github"))
	response += fmt.Sprintf("%s %s  %s\n", BlueBold("$ aum"), UpdateCommand.ToHelp(), YellowItalic("Update a Among Us mod from github"))
	response += fmt.Sprintf("%s %s  %s\n", BlueBold("$ aum"), RemoveCommand.ToHelp(), YellowItalic("Remove a installed Among Us mod"))
	response += fmt.Sprintf("%s %s                    %s\n", BlueBold("$ aum"), VanillaCommand.ToHelp(), YellowItalic("Disable Among Us mod from your game folder"))
	response += fmt.Sprintf("%s %s                      %s\n", BlueBold("$ aum"), ResetCommand.ToHelp(), YellowItalic("Remove all installed Among Us mods"))
	fmt.Println(response)
}

func SetGameFolderPath(path string) {
	if !dirExists(path) {
		color.Red("[ERROR] Directory %s does not exist", path)
		return
	}
	Mods.SetGameFolderPath(strings.ReplaceAll(path, "\\", "/"))
	color.Green("[SUCCESS] Game folder path successfully saved!")
}

func RemoveMod(mod *ModType, owner string, repo string) {
	if mod == nil {
		color.Red("[ERROR] Among Us Mod %s/%s is not installed", owner, repo)
		return
	}
	Mods.Remove(mod)
	color.Green("[SUCCESS] Among Us Mod %s/%s successfully removed!", owner, repo)
}

func InstallMod(mod *ModType, owner string, repo string) {
	if mod != nil {
		ActivateMod(mod.Path)
		return
	}
	color.Blue("[INFO] Search for Among Us mod %s/%s...", owner, repo)
	release := GetLatestRelease(owner, repo)
	if release != nil {
		color.Blue("[INFO] Mod %s/%s found, downloading...", owner, repo)
		DownloadMod(release, owner, repo)
		Mods.Add(release)
		color.Green("[SUCCESS] Among Us Mod %s/%s successfully installed!", owner, repo)
		ActivateMod(fmt.Sprintf("%s/%s", owner, repo))
	} else {
		color.Red("[ERROR] Unable to find mod %s/%s!", owner, repo)
	}
}

func RemoveCurrentMod(showLog bool) {
	if Mods.CurrentMod != nil {
		Mods.UninstallMod(Mods.CurrentMod)
		Mods.SetCurrent(nil)
	}
	if showLog {
		color.Green("[SUCCESS] Among Us is now vanilla!")
	}
}

func UpdateMod(mod *ModType, owner string, repo string) {
	if mod == nil {
		color.Red("[ERROR] Among Us Mod %s/%s is not installed!", owner, repo)
		return
	}
	color.Blue("[INFO] Looking for a new release of mod %s/%s...", owner, repo)
	release := GetLatestRelease(owner, repo)
	if release.Id != mod.Id {
		color.Blue("[INFO] An update for Among Us Mod %s/%s was found, downloading...", owner, repo)
		Mods.Remove(mod)
		DownloadMod(release, owner, repo)
		Mods.Update(release)
		color.Green("[SUCCESS] Among Us Mod %s/%s successfully updated!", owner, repo)
		ActivateMod(mod.Path)
		return
	}
	color.Yellow("[WARNING] Among Us Mod %s/%s already up to date!", owner, repo)
}

func ActivateMod(modPath string) {
	RemoveCurrentMod(false)
	mod := Mods.GetModByPath(modPath)
	if mod != nil {
		Mods.SetCurrent(mod)
		color.Green("[SUCCESS] Among Us Mod %s successfully activated!", mod.Path)
		return
	}
	color.Red("[ERROR] Among Us Mod %s is not installed!", modPath)
}

func DownloadMod(release *Release, owner string, repo string) {
	if release != nil {
		asset := release.GetZipAsset()
		if asset != nil {
			DownloadFile(asset.DownloadUrl, DownloadFolderPath)
			modDirectory := fmt.Sprintf("%s/%s/%s", ModsFolderPath, owner, repo)
			mkdirIfNoExists(fmt.Sprintf("%s/%s", ModsFolderPath, owner))
			mkdirIfNoExists(modDirectory)
			color.Blue("[INFO] Extracting mod %s/%s...", owner, repo)
			handleError(UnzipSource(fmt.Sprintf("%s/%s", DownloadFolderPath, asset.Name), modDirectory))
			handleError(os.Remove(fmt.Sprintf("%s/%s", DownloadFolderPath, asset.Name)))
			return
		}
	}
	log.Fatal(fmt.Sprintf("Mod %s/%s not found", owner, repo))
}

func amAdmin() bool {
	file, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	closeFile(file)
	return true
}
