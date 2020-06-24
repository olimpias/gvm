# gvm

[![GVM](https://circleci.com/gh/olimpias/gvm.svg?style=svg)](<https://app.circleci.com/pipelines/github/olimpias/gvm>)


gvm is a short version of `Go Version Manager`. It allows you to manage your installed go version and change its versions
according to your requirements. 

Supports OSx, linux and windows. It requires to preinstalled go(for now...) and `GOPATH` to be
set in environmental variables. In addition to that, the user privileges must be set to run the program. Otherwise, there will
be permission issues to change and replace folders.

## Installation

Download release version from [releases](https://github.com/olimpias/gvm/releases) according to your operating system. For now, releases are done for OSX and windows 64.

Extract the executable from tar/zip. Then enjoy it!

**Note:** If you are trying to install go with this file, you have to set bash configuration by yourself. Also, **gvm** requires `GOROOT` environmental variable.

Checkout example for [setting up environment](#setting-up-env-for-v001-from-scratch-in-macos)

## Command Usages
gvm provides 5 type of commands that you can apply; `help`, `list`, `use`, `download` and `del`

### List Command
Lists the possible go versions that you have downloaded to your local machine. To download a specific version checkout [Download Command](#download-command)

Command usage: `gvm list`

Example usage: 

`gvm list`

```
go1.10.1
go1.13.1
go1.13.8
go1.14.2
go1.14.3
```

### Download Command
Downloads the version that you inputted the command. It downloads tar file into your `$HOME/.gvm` path.

`gvm dl <go-version>`

Example Usage:

`gvm dl 1.11.1`

```
118.43 MiB / 118.43 MiB [-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------] 100.00% 23.81 MiB p/s 5s
```

### Use Command
Sets go version that you inputted into the command. It unzip the downloaded file and moves files into your `GOROOT` path. It may take sometime. In future versions, progressbar is going to be added.
To able to use a go version, you need to download it before using `gvm dl <go-version>` [Download Command](#download-command)

Command Usage: `gvm use <go-version>`

Example Usage:

`gvm use 1.11.1`

### Del Command
Deletes the version of go file that you inputted as a version from your `$HOME/.gvm`.

Command usage:`gvm del <go-version>`

### Help Command
To print help

Command usage: `gvm help`

```
gvm is a go version controller
Commands:
list  list the possible downloaded versions that ready to use.
dl    downloads the version that you specify to your machine.
use   uses the version that specify as an input. It has to be downloaded first using dl command.
del   deletes the version that you specify as an input
```

## Setting Up Env For v0.0.1 From Scratch in MacOS

Be sure if go is installed. Otherwise it is hardly recommended to install for alpha version.

Then, check if `goroot` is set with  `echo $GOROOT`.

If it returns empty try using `which go`, this will return go path.

Set env with `export GOROOT=PATH`(put result from `which go` command, you should exclude /bin/go in path). For default installation it is most likely `/usr/local/go`. It is recommanded to add it to bash profile
otherwise you need to set `GOROOT` into environmental variable all the time to use **gvm**.
Example result `/usr/local/go/bin/go`, you need to use `/usr/local/go` as `GOROOT`.

Check the permission for users. Most likely for default installation it is assigned to root user and your current user does not have access for that directory to edit.

Go to `cd $GOROOT/..` path and use `sudo chmod -R 777 go` or `sudo chmod -R 755 go`.

Click [link](https://github.com/olimpias/gvm/releases/download/v0.0.1/gvm_0.0.1.darwin-amd64.tar.gz) to download tar file.

Use `tar xvf gvm_0.0.1.darwin-amd64.tar.gz` to extract tar file and it will extract `gvm` executable. It is ready for use

You need to use it with `./` at the beginning. Example: `./gvm dl 1.14.4`

## TODOs
- [X] Add progressbar for unzipping
- [ ] Add CI/CD for testing windows and linux over circleci
- [ ] Add executable as a downloadable so that it could be usable through homebrew or other example platforms.
- [ ] Provide initial installation of go(without preinstalled go) with executable file.
