package main

import (
	"path"
	"strconv"
	"strings"
)

var listSpaces int = -3

func Info() {

	if len(command.Args) < 1 {
		command.LogError("No template name has been provided!")
	}

	template := command.CreateTemplateManager(command.Dirname()).GetTemplateInfo(strings.Join(command.Args, " "))
	description := template.Description
	preScripts := strings.Join(config.PreScripts, " && ")
	ignore := strings.Join(config.Ignore, ", ")

	if len(description) == 0 {
		description = "No description."
	}

	if len(preScripts) == 0 {
		preScripts = "No pre scripts."
	}

	if len(ignore) == 0 {
		ignore = "No ignored files."
	}

	command.Println("Template information of [1m" + command.Args[0] + "[0m\n")
	command.Println("[90m-[0m [1mName:[0m             " + template.Name)
	command.Println("[90m-[0m [1mDescription:[0m      " + description)
	command.Println("[90m-[0m [1mPre-Scripts:[0m      " + preScripts)
	command.Println("[90m-[0m [1mIgnored files:[0m    " + ignore)

	maxLength := 0

	for name := range template.Scripts {
		l := len(name)
		if maxLength < l {
			maxLength = l
		}
	}

	if maxLength > 0 {
		command.Println("[90m-[0m [1mScripts:[0m")
		for name, script := range template.Scripts {
			command.Println("    [90m-[0m [1m" + name + ":[0m" + strings.Repeat(" ", (maxLength-len(name))+3) + script)
		}
	}

}

func All() {
	command.Print("[1mAll the templates saved.[0m\n\n")
	templates := command.CreateTemplateManager(command.Dirname()).GetAllTemplates()

	if len(templates) == 0 {
		command.Println("No templates have been saved.")
	} else {
		for i, template := range templates {
			command.Println("[1m" + strconv.Itoa(i+1) + ".[0m " + template[1:])
		}
	}

	command.Println("\nYou can use [90mtemplatify info <template-name>[0m to show the information!")
}

func List() {
	if len(command.Args) < 1 {
		command.LogError("No template name has been provided!")
	}

	dirData := ReadDir(path.Join(command.Dirname(), "templates", command.Args[0]))
	lockExists := false

	for _, file := range dirData {
		if file.Name == "templatify.lock.json" {
			lockExists = true
		}
	}

	if !lockExists {
		command.LogError("Unknown template \"" + command.Args[0] + "\".")
	}

	command.Print("Template structure for [1m" + command.Args[0] + "[0m\n\n")
	pushList(dirData)
}

func pushList(assets []FolderAsset) {
	listSpaces += 3
	for _, asset := range assets {
		if asset.IsDir {
			command.Println(strings.Repeat(" ", listSpaces) + "[90m-[0m " + asset.Name + "/")
			pushList(asset.Assets)
		} else {
			command.Println(strings.Repeat(" ", listSpaces) + "[90m-[0m " + asset.Name)
		}
	}
}

func Help() {
	help := []string{
		"[1mTemplatify Commands[0m",
		"[90m[arg] means optional argument and <arg> means required argument.[0m\n",
		"help           To show this help page",
		"all            Shows all the templates saved",
		"removeall      Removes all the templates",
		"init           Creates a default templatify.config.json",
		"pre-scripts    Runs the pre scripts saved in the lock",
		"save [90m[name][0m    Save the current working directory as a template!",
		"       [90m--clean[0m                 Total rewrite over the old template else will just merge",
		"use [90m<name>[0m     Use a template",
		"       [90m--custom-path=<path>[0m    Set a custom path to use the template else will use the path as current working directory",
		"       [90m--no-git[0m                Will not clone .git folder",
		"       [90m--remove-lock[0m           Will remove templatify's lock file",
		"       [90m--disable-pre-scripts[0m   Will not run pre scripts",
		"get [90m<repo>[0m     Download a repo as a template",
		"info [90m<name>[0m    Shows info of a template",
		"list [90m<name>[0m    Returns the folder map of the template",
		"remove [90m<name>[0m  Delete a template",
		"exec [90m<name>[0m    Execute a script by its name",
		"test           Executes the test script if it exists",
		"\nYou can read the guide in github [94;4mhttps://github.com/scientific-guy/templatify[0m.",
	}

	for _, str := range help {
		command.Println(str)
	}
}
