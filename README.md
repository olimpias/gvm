# gvm

[![GVM](https://circleci.com/gh/olimpias/gvm.svg?style=svg)](<https://app.circleci.com/pipelines/github/olimpias/gvm>)


gvm is a short name of `Go Version Manager`. It allows you to manage your go versions and change versions
according to your requirements. It is completely written in Go and provides cross-platform usability

gvm supports Windows, Linux, Freebsd and MacOS. It requires to preinstalled go for windows and `GOROOT` to be
set in environmental variables for windows. In addition to that, the user privileges must be set to run the program for all OSs. Otherwise, there will
be permission issues while removing and replacing the folders/files. [Check Troubleshooting section](#Troubleshooting).

## Installation

**Linux**:
```shell script
#Example for amd64 arch
# Linux Example (assumes ~/bin is in PATH).
curl -o gvm.tar.gz -O https://github.com/olimpias/gvm/releases/download/v0.1.0/gvm.linux.amd64.tar.gz
tar -C ~/bin -xzf gvm.tar.gz
chmod +x ~/bin/gvm
gvm use 1.14.4
```

Supported Archs For Linux : `386`, `amd64`, `arm`, `arm64`, `mips`, `mipsle`, `mips64`, `mips64le`, `ppc64`, `ppc64le`, `s390x`

**MacOS**: 
```shell script
# Example for amd64 arch
curl -o gvm.tar.gz -O https://github.com/olimpias/gvm/releases/download/v0.1.0/gvm.darwin.amd64.tar.gz
tar -C /usr/local/bin -xzf gvm.tar.gz
chmod +x /usr/local/bin/gvm
gvm use 1.14.4
```

Supported Archs For MacOs : `amd64`

**Windows**:
```shell script
# Example for amd64 arch
curl -o gvm.zip -O https://github.com/olimpias/gvm/releases/download/v0.1.0/gvm.windows.amd64.exe.zip
unzip gvm.zip
gvm use 1.14.4
```

Supported Archs For MacOs : `386`, `amd64`

**FreeBSD**:
```shell script
#Example for amd64 arch
# Freebsd Example (assumes ~/bin is in PATH).
curl -o gvm.tar.gz -O https://github.com/olimpias/gvm/releases/download/v0.1.0/gvm.freebsd.amd64.tar.gz
tar -C ~/bin -xzf gvm.tar.gz
chmod +x ~/bin/gvm
gvm use 1.14.4
```

Supported Archs For Freebsd : `386`, `amd64`, `arm`, `arm64`


To support use other architecture, you just need to change **amd64** in the link for `curl` command.

**Note:** If you are trying to install go for the first time, it will not work for Windows OS. It requires preinstalled go for windows.

Download release version from [releases](https://github.com/olimpias/gvm/releases) according to your operating system and architecture.

Extract the executable from tar/zip. Then enjoy it!

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
If you have not downloaded the version before, `use` command will trigger download of it as well. Then it will set the go version.

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

## Troubleshooting

When Go is installed by default installer, the created folder could have access restrictions. If it is the case please use
following command to bypass it.

Mostly Go is installed to `/usr/local/go` path by default installer. We need to grant access for `/usr/local/go` file path.

```shell script
sudo chmod -R 755 /usr/local/go
```