# Templatify

A cli to create local templates and copy templates from github which is saved within your pc and used easily!

> This project is currently under development. Incase if you want to help create a pr or an issue. ~~The code might look unorganized~~.

![Example](/assets/example.gif)

## Quick guide

### Version

You can view up the version by

```sh
> templatify --version
Current templatify version: 1.0.0
```

### Saving templates

You can save the current working directory as a template!

```sh
> templatify save TestTemplate
SUCCESS: Successfully copied template with name as "TestTemplate"
```

So here `TestTemplate` is the name of the template to save. If none provided will save with the name of directory. 
> Name of the template should not contain spaces!

You can also use configuration file like `templatify.config.json`! An example of the config file is

```json
{
    "name": "TestTemplate",
    "description": "This is my test template",
    "delimiter": "%",
    "parseFiles": "*.js",
    "preScripts": [
        "echo \"I have been ran by templatify!\""
    ],
    "ignore": [
        "node_modules",
        "package-lock.json",
        "test/*.js"
    ]
}
```

There are more configuration can be done with that! Here is are list of fields. **All are optional fields**.

- **name** - Name of the template to save.
- **description** - Not useful. But will help you to remember what the template is for.
- **preScripts** - Array of subprocess scripts to execute while creating the template.
- **ignore** - Array of files to ignore when saving as a template.
- **delimiter** - The delimiter for parsing files, by default it will be as `%` which will be used as `%{key}`. If you set it as `$` then it would be `${key}` to be parsed.
- **parseFiles** - A glob string for the files to be parsed. This will be required to enable parsing else it would not. Remember if the template consists of any binary files, it might throw error while parsing files so use a perfect glob string for it.

> You can directly use `templatify init` to create a default config file!

Normally this template overwrite over the old changes if exists. If you want to completely delete and save as a fresh template you have to use `--clean` flag

```sh
> templatify save TestTemplate --clean
INFO: Successfully cleaned directory.
SUCCESS: Successfully copied template with name as "TestTemplate"
```

You can also use templates stored in github

```sh
> templatify get repo
CONFIRM Perform template configuration? (y/n) y
INFO Performing termplate configuration.
SUCCESS Saved "repo" as a template.
```

So here repo should be something like `username/reponame`. But any github repo cannot be a templatify template it must have a `templatify.config.json` file in the root of it! Currently there is no support to import from branches. While downloading the repo as a template it might ask you to allow template configuration, this will remove ignore files. And also it would have the `.git` folder with it but if you want to prevent it while using the template you can use the `--no-git` flag while using it. 

> Using the `name` field in the config file in the repo is useless because when storing the template it would be `username/reponame`. You can still keep the name field.

> Templates are stored in `templates` directory in the directory where templatify exists

### Using templates

You can use a template into a directory like this

```sh
> templatify use TestTemplate
INFO Copying template "TestTemplate" to "E:\Projects\test".
INFO Cloned files.
PRE-SCRIPT echo "I have been ran by templatify!"
I have been ran by templatify!
INFO Ran all preScripts.
SUCCESS Finished in 0s
```

In any case you want to store it in a custom path within the current working directory lets say you want to store it in `my-template` subdir you can do something like this

```sh
> templatify use TestTemplate --custom-path=my-template
INFO Copying template "TestTemplate" to "E:\Projects\test/my-template".
INFO Cloned files.
PRE-SCRIPT: echo "I have been ran by templatify!"
I have been ran by templatify!
INFO Ran all preScripts.
SUCCESS Finished in 0s
```

> There might be problems with forward slashes and backward slashes on file paths which will not matter alot because they are only for representation.

While using the template you might see a additional file named `templatify.lock.json` which is a lock file saving all configurations. You can remove it though by using `--remove-lock` flag.

```sh
> templatify use TestTemplate --custom-path=my-template --remove-lock
INFO Copying template "TestTemplate" to "E:\Projects\test/my-template".
INFO Cloned files.
PRE-SCRIPT echo "I have been ran by templatify!"
I have been ran by templatify!
INFO Ran all preScripts.
SUCCESS Finished in 0s
```

When saving any template with the `.git` directory with it might confuse with the current git in the path to use the template in that case you can use the `--no-git` flag. This might be useful to use a template from github and not use the git.

### Get all templates

You can view what and all templates you have stored by

```sh
> templatify all
All the templates saved.

1. TestTemplate

You can use `templatify info <template-name>` to show the information!
```

### Get paticular template

You can view the information about a template

```sh
> templatify info TestTemplate
Template information of "TestTemplate"

- Name:             TestTemplate
- Description:      This is my test template
- Pre-Scripts:      echo "I have been ran by templatify!"
- Ignored files:    node_modules, package-lock.json, test/*.js
```

To view the files of a template you can do something like

```sh
> templatify list TestTemplate
Template structure for "TestTemplate"

- templatify.lock.json
- package.json
- index.js
- templatify.config.json
- test/
    - main_test.go
    - main.go
```

> You will be viewing a `templatify.lock.json` file in the template structure. You can remove it while using the template using `--remove-lock`.

### Removing templates

If you want to remove template you have to do:

```sh
> templatify remove TestTemplate
CONFIRM Are you sure? (y/n) y
SUCCESS Deleted template successfully.
```

You can also remove all templates

```sh
> templatify removeall
CONFIRM Are you sure? (y/n) y
SUCCESS Deleted all templates successfully.
```