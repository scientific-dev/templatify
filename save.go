package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

var (
	templateDir string
	cwd         string
	config      Config
	manager     TemplateManager
)

func Save() {
	dirname := command.Dirname()
	config = command.ProcessConfig()
	cwd = command.Cwd()
	dirData := ReadDir(cwd)
	manager = command.CreateTemplateManager(dirname)
	parsedConfig, parseErr := json.Marshal(config)
	config.ProcessIgnoreMatches()

	if parseErr != nil {
		command.LogError("Unknown error: " + parseErr.Error())
	}

	flags := command.Flags()
	_, cleanDir := flags["clean"]

	if cleanDir {
		manager.CleanTemplateDir(config.Name)
		command.Println("[36mINFO[0m Successfully cleaned directory.")
	} else {
		manager.CheckTemplateDir(config.Name)
	}

	templateDir = path.Join(dirname, "templates", config.Name)
	saveTemplateDir("", dirData)
	manager.CreateAbsoluteFile(path.Join(templateDir, "templatify.lock.json"), parsedConfig, 0644)
	command.Println("[32mSUCCESS[0m Successfully copied template with name as \"" + config.Name + "\"")
}

func saveTemplateDir(pathname string, dirData []FolderAsset) {
	for _, asset := range dirData {
		if asset.IsDir {
			name := path.Join(pathname, asset.Name)
			if dirPath := path.Join(templateDir, name); !isIgnoredFile(name) {
				os.Mkdir(dirPath, asset.Mode())
				saveTemplateDir(name, asset.Assets)
			}
		} else {
			if filePath := path.Join(pathname, asset.Name); !isIgnoredFile(filePath) {
				saveTemplateFile(filePath, asset.Mode())
			}
		}
	}
}

func saveTemplateFile(pathname string, mode os.FileMode) {
	content, readErr := ioutil.ReadFile(path.Join(cwd, pathname))
	if readErr != nil {
		command.LogError("Failed reading a file: " + readErr.Error())
	}

	err := ioutil.WriteFile(path.Join(templateDir, pathname), content, mode)
	if err != nil {
		command.LogError("Error originated while storing the file in templates: " + err.Error())
	}
}

func Download() {
	if len(command.Args) < 1 {
		command.LogError("No repo url has been provided!")
	}

	command.CreateTemplateManager(command.Dirname()).SaveGithubRepo(command.Args[0])
}
