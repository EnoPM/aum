package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type ConfigType struct {
	GameFolderPath *string   `json:"game_folder_path"`
	CurrentMod     *ModType  `json:"current_mod"`
	ModsList       []ModType `json:"mods"`
}

type ModType struct {
	Id   int    `json:"id"`
	Path string `json:"path"`
}

var Mods = ConfigType{CurrentMod: nil, ModsList: make([]ModType, 0)}

func init() {
	if fileExists(ConfigFilePath) {
		Mods.Read()
	} else {
		Mods.Save()
	}
}

func (mod *ModType) GetAbsolutePath() string {
	return fmt.Sprintf("%s/%s", ModsFolderPath, mod.Path)
}

func (mods *ConfigType) Read() {
	file, err := os.Open(ConfigFilePath)
	if err != nil {
		panic(err)
	}
	defer closeFile(file)
	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, mods)
	if err != nil {
		panic(err)
	}
}

func (mods *ConfigType) Save() {
	content, err := json.Marshal(mods)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(ConfigFilePath, content, 0755)
	if err != nil {
		panic(err)
	}
}

func (mods *ConfigType) Add(release *Release) {
	asset := release.GetZipAsset()
	if asset != nil {
		mods.ModsList = append(mods.ModsList, ModType{Id: release.Id, Path: release.Path})
		mods.Save()
	}
}

func (mods *ConfigType) Update(release *Release) {
	for key, mod := range mods.ModsList {
		if mod.Path == release.Path {
			mods.ModsList[key].Id = release.Id
			mods.Save()
			return
		}
	}
}

func (mods *ConfigType) Reset() {
	for _, mod := range mods.ModsList {
		mods.Remove(&mod)
	}
}

func (mods *ConfigType) Remove(mod *ModType) {
	if mods.CurrentMod != nil && mod.Id == mods.CurrentMod.Id {
		mods.UninstallMod(mod)
		mods.SetCurrent(nil)
	}
	if fileOrDirExists(mod.GetAbsolutePath()) {
		err := os.RemoveAll(mod.GetAbsolutePath())
		if err != nil {
			panic(err)
		}
	}
	for key, currentMod := range mods.ModsList {
		if currentMod.Id == mod.Id {
			copy(mods.ModsList[key:], mods.ModsList[key+1:])
			mods.ModsList[len(mods.ModsList)-1] = ModType{}
			mods.ModsList = mods.ModsList[:len(mods.ModsList)-1]
			mods.Save()
			break
		}
	}
}

func (mods *ConfigType) SetCurrent(mod *ModType) {
	if mod != nil {
		mods.InstallMod(mod)
	}
	mods.CurrentMod = mod
	mods.Save()
}

func (mods *ConfigType) SetGameFolderPath(folderPath string) {
	mods.GameFolderPath = &folderPath
	mods.Save()
}

func (mods *ConfigType) GetModById(id int) *ModType {
	for _, mod := range mods.ModsList {
		if mod.Id == id {
			return &mod
		}
	}
	return nil
}

func (mods *ConfigType) GetModByPath(path string) *ModType {
	for _, mod := range mods.ModsList {
		if mod.Path == path {
			return &mod
		}
	}
	return nil
}

func (mods *ConfigType) InstallMod(mod *ModType) {
	if mods.GameFolderPath == nil {
		panic("Please set your game folder path (aum -gfp=\"C:\\Program Files (x86)\\Steam\\steamapps\\common\\Among Us\")")
	}
	files, err := os.ReadDir(mod.GetAbsolutePath())
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s/%s", ModsFolderPath, mod.Path, file.Name())
		symLinkPath := fmt.Sprintf("%s/%s", *mods.GameFolderPath, file.Name())
		err = os.Symlink(filePath, symLinkPath)
		if err != nil {
			panic(err)
		}
	}
}

func (mods *ConfigType) UninstallMod(mod *ModType) {
	if mods.GameFolderPath == nil {
		panic("Please set your game folder path (aum -gfp=\"C:\\Program Files (x86)\\Steam\\steamapps\\common\\Among Us\")")
	}
	files, err := os.ReadDir(mod.GetAbsolutePath())
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		symLinkPath := fmt.Sprintf("%s/%s", *mods.GameFolderPath, file.Name())
		if fileOrDirExists(symLinkPath) {
			err = os.Remove(symLinkPath)
			if err != nil {
				panic(err)
			}
		}
	}
}
