# gomake

## Redesign of GNU Make (Makefile) oriented towards project management and command uniformisation among all your projects.

When working in an microservice environment with multiple languages and/or frameworks it is sometimes hard to switch between command sets.
This is why I have created gomake to help better manage running commands between projects:

## Examples: 

Go project
```yml
targets:
    run: go run /path/to/package
```

Ruby on rails project: 
```yml
targets:
    run: rails -s -b 0.0.0.0 -p 3000
```

PHP: 
```yml
targets:
    run: php -S localhost:8000
```

#### For all three projects you need only to run `gomake run` in order to start the service or server.


## Installation

If you have golang installed (and GOBIN is added to your path): 

```
go get github.com/efimovalex/gomake
```

### Linux:

```
wget https://github.com/efimovalex/gomake/releases/download/v1.0.0/gomake_1.0.0_Linux_x86_64.tar.gz
tar -C /usr/local -xzf gomake_1.0.0_Linux_x86_64.tar.gz
```

### OS X
```
brew tap efimovalex/gomake
brew install efimovalex/gomake/gomake
```

### Microsoft Windows: 

Download the exe file for your system from the [latest release](https://github.com/efimovalex/gomake/releases/latest) and unarchive it to your `C:\Windows` folder in order to be accesible from everywhere in the system.
You can create an alias to remove de `.exe` ending by running in PowerShell:

```
Set-Alias -Name gomake -Value C:\Windows\gomake.exe
```


## Project settings file: makefile.yml
You can provide you own file with the `-file=path/to/makefile.yml`

### CLI 
Point to the CLI you want to use: `bash`,`sh` and for windows `cmd`, `powershell` and `bash` (`bash.exe` if you have wsl)

### VARS: target variables
### Define variables that are replaced troughout your tagets defined in the yml file

```yml
cli: bash
vars:
    buildCmd: go build -ldflags "-X main.BuildName=${package_name} -X main.BuildVersion=${version}" 
    package_name: pkg
    package: pkg
    GOARCH: amd64
    version: "1.0.0"
targets:
    build_darwin: env GOOS=darwin GOARCH=${GOARCH} ${buildCmd} -o ${package_name}_v${version}_darwin_amd64 ${package}
    build_windows: env GOOS=windows GOARCH=${GOARCH} ${buildCmd} -o ${package_name}_v${version}_windows_amd64.exe ${package}
    build_linux: env GOOS=linux GOARCH=${GOARCH} ${buildCmd} -o ${package_name}_v${version}_linux_amd64 ${package}
```

## ENV: Environment variables
### Define env variables that are loaded only for the purpose of your project isolated from your current environment

```yml
cli: bash
vars: 
    package_name: pkg
    package: pkg
    version: "1.0.0"
    buildCmd: go build -ldflags "-X main.BuildName=${package_name} -X main.BuildVersion=${version}" 
env:
    GOARCH: amd64
targets:
    build_darwin: env GOOS=darwin ${buildCmd} -o ${package_name}_v${version}_darwin_amd64 ${package}
    build_windows: env GOOS=windows ${buildCmd} -o ${package_name}_v${version}_windows_amd64.exe ${package}
    build_linux: env GOOS=linux ${buildCmd} -o ${package_name}_v${version}_linux_amd64 ${package}
```

## Targets:
### Commands that are run when the target name is provided

```yml
targets:
    build_darwin: env GOOS=darwin ${buildCmd} -o ${package_name}_v${version}_darwin_amd64 ${package}
    build_windows: env GOOS=windows ${buildCmd} -o ${package_name}_v${version}_windows_amd64.exe ${package}
    build_linux: env GOOS=linux ${buildCmd} -o ${package_name}_v${version}_linux_amd64 ${package}
```

```sh
gomake build_darwin
```

If the target is not found/defined you will get an error message saying: `no target found for {target}`

### Multi OS
If you are working in an multi os team, you could define a separate file for windows users and provide that file on windows machines which should contain the same targets but defined as windows commands. 

```powershell
gomake -f makefile_win.yml
```
