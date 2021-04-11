package main

import (
	"os"
	"path"
)

func Remove() {
	if len(command.Args) < 1 {
		command.LogError("No template name has been provided!")
	}

	if command.Confirm("[90mCONFIRM[0m Are you sure? (y/n) ") {
		if removeErr := os.RemoveAll(path.Join(command.Dirname(), "templates", command.Args[0])); removeErr != nil {
			if os.IsPermission(removeErr) {
				command.LogError("Access denied to remove directories. Run the program as administrator. " + removeErr.Error())
			}

			command.LogError("Could not remove directory \"" + templateDir + "\": " + removeErr.Error())
		}

		command.Println("[32mSUCCESS[0m Deleted template successfully.")
	} else {
		command.LogError("Failed deleting the template.")
	}
}

func RemoveAll() {
	dirname := command.Dirname()
	if command.Confirm("[90mCONFIRM[0m Are you sure? (y/n) ") {
		for _, template := range command.CreateTemplateManager(dirname).GetAllTemplates() {
			err := os.RemoveAll(path.Join(dirname, "templates", template))
			if err != nil {
				if os.IsPermission(err) {
					command.LogError("Access denied to remove directories. Run the program as administrator. " + err.Error())
				}
				command.LogError("Failed removing a directory: " + err.Error())
			}
		}

		command.Println("[32mSUCCESS[0m Removed all templates.")
	} else {
		command.LogError("Failed deleting the templates.")
	}
}
