package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/src-d/go-git.v4"
)

type TemplateManager struct {
	Dirname string
}

func (self TemplateManager) Init() {
	templateDir := path.Join(self.Dirname, "templates")
	_, readErr := os.Stat(templateDir)

	if os.IsNotExist(readErr) {
		os.Mkdir(templateDir, 0755)
	} else if readErr != nil {
		command.LogError("Unknown error while reaching the templates directory: " + readErr.Error())
	}
}

func (self TemplateManager) CheckTemplateDir(name string) {
	self.CheckAbsoluteDir(path.Join(self.Dirname, "templates", name))
}

func (self TemplateManager) CheckAbsoluteDir(pathname string) {
	_, readErr := os.Stat(pathname)

	if os.IsNotExist(readErr) {
		if err := os.Mkdir(pathname, 0755); err != nil {
			command.LogError("Failed creating a directory: " + err.Error())
		}
	} else if readErr != nil {
		command.LogError("Unknown error while reaching the templates directory: " + readErr.Error())
	}
}

func (self TemplateManager) CleanTemplateDir(name string) {
	templateDir := path.Join(self.Dirname, "templates", name)
	if removeErr := os.RemoveAll(templateDir); removeErr != nil {
		command.LogError("Could not remove directory \"" + templateDir + "\": " + removeErr.Error())
	}

	if makeErr := os.Mkdir(templateDir, 0755); makeErr != nil {
		command.LogError("Could not create directory \"" + templateDir + "\": " + makeErr.Error())
	}
}

func (self TemplateManager) CreateAbsoluteFile(name string, content []byte, mode os.FileMode) {
	err := ioutil.WriteFile(name, content, mode)
	if err != nil {
		command.LogError("Error originated while storing the file in templates: " + err.Error())
	}
}

func (self TemplateManager) GetTemplateInfo(name string) Config {
	var config Config
	rawData, readErr := ioutil.ReadFile(path.Join(self.Dirname, "templates", name, "templatify.lock.json"))

	if readErr != nil {
		command.LogError("Could not read template \"" + name + "\"!")
	}

	if parseErr := json.Unmarshal(rawData, &config); parseErr != nil {
		command.LogError("Failed parsing lock file for the template: " + parseErr.Error())
	}

	return config
}

func (self TemplateManager) GetAllTemplates() []string {
	templateDirName := path.Join(self.Dirname, "templates")
	files, readErr := ioutil.ReadDir(templateDirName)

	if readErr != nil {
		command.LogError("Failed reading template files: " + readErr.Error())
	}

	return resolveTemplatesInDir(templateDirName, "", files)
}

func (self TemplateManager) SaveGithubRepo(repo string) {
	repo = GetGHBaseUrl(repo)
	templateDir = path.Join(self.Dirname, "templates", repo)
	flags := command.Flags()
	_, cleanTemplateDir := flags["clean"]

	if cleanTemplateDir {
		if err := os.RemoveAll(templateDir); err != nil {
			if os.IsPermission(err) {
				command.LogError("Access denied to clean directory before downloading it: " + err.Error())
			} else {
				command.LogError("Could not clean directory: " + err.Error())
			}
		}
	}

	if _, err := git.PlainClone(templateDir, false, &git.CloneOptions{URL: "https://github.com/" + repo}); err != nil {
		command.LogError("Could not clone repo \"" + repo + "\" from github: " + err.Error())
	}

	rawData, readErr := ioutil.ReadFile(path.Join(templateDir, "templatify.config.json"))

	if readErr != nil {
		command.LogError("Could not read template \"" + repo + "\"!")
	}

	if parseErr := json.Unmarshal(rawData, &config); parseErr != nil {
		command.LogError("Failed parsing config file for the repo: " + parseErr.Error())
	}

	for _, ignore := range config.Ignore {
		ignoreMatches = append(ignoreMatches, parseIgnoreMatch(ignore))
	}

	if command.Confirm("[90mCONFIRM[0m Perform template configuration? (y/n) ") {
		command.Println("[36mINFO[0m Performing termplate configuration.")
		checkRepoFolder("", ReadDir(path.Join(templateDir)))
	} else {
		command.Println("[36mINFO[0m Skipping termplate configuration.")
	}

	config.Name = repo
	data, parseErr := json.Marshal(config)

	if parseErr != nil {
		command.LogError("Failed marshalling json for lock file: " + parseErr.Error())
	}

	if err := ioutil.WriteFile(path.Join(templateDir, "templatify.lock.json"), data, 0644); err != nil {
		command.LogError("Failed writing the lock file: " + err.Error())
	}

	command.Println("[32mSUCCESS[0m Saved \"" + repo + "\" as a template.")
}

func checkRepoFolder(pathname string, assets []FolderAsset) {
	for _, asset := range assets {
		if asset.IsDir {
			if asset.Name != ".git" {
				pathname = path.Join(pathname, asset.Name)
				if isIgnoredFile(pathname) {
					if err := os.RemoveAll(path.Join(templateDir, pathname)); err != nil {
						command.LogWarn("Failed removing a ignored directory: " + err.Error())
					}
				} else {
					checkRepoFolder(pathname, asset.Assets)
				}
			}
		} else {
			checkRepoFile(path.Join(pathname, asset.Name))
		}
	}
}

func checkRepoFile(pathname string) {
	if isIgnoredFile(pathname) {
		if err := os.Remove(path.Join(templateDir, pathname)); err != nil {
			command.LogWarn("Failed removing a ignored file: " + err.Error())
		}
	}
}
