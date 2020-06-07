# gvm
gvm is a short version of `Go Version Manager`. It allows you to manage your installed go version and change its versions
according to your requirements. 

Supports OSx, linux and windows. It requires to preinstalled go(for now...) and `GOPATH` to be
set in environmental variables. In addition to that, the user privileges must be set to run the program. Otherwise, there will
be permission issues to change and replace folders.

##Installation

It is going to be fill in... in progress

##Command Usages
gvm provides 5 type of commands that you can apply; `help`, `list`, `use`, `download` and `del`

###List Command
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
###Download Command
Downloads the version that you inputted the command. It downloads tar file into your `$HOME/.gvm` path.

`gvm dl <go-version>`

Example Usage:

`gvm dl 1.11.1`

```
118.43 MiB / 118.43 MiB [-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------] 100.00% 23.81 MiB p/s 5s
```

###Use Command
Sets go version that you inputted into the command. It unzip the downloaded file and moves files into your `GOROOT` path. It may take sometime. In future versions, progressbar is going to be added.
To able to use a go version, you need to download it before using `gvm dl <go-version>` [Download Command](#download-command)

Command Usage: `gvm use <go-version>`

Example Usage:

`gvm use 1.11.1`

###Del Command

Deletes the version of go file that you inputted as a version from your `$HOME/.gvm`.

Command usage:`gvm del <go-version>`


###Help Command
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

## TODOs
* Add progressbar for unzipping
* Add CI/CD for testing windows and linux over circleci
* Add executable as a downloadable so that it could be usable through homebrew or other example platforms.
* Provide initial installation of go(without preinstalled go) with executable file.
