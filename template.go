package main

const initFile = `cli: bash
vars: 
    current_dir: ~/
    app_name: app
env: 
    LOCATION: ./

default: build
    
targets:
    # build command for a go project
    build:  go build -o ${app_name}

    # runs tests with coverage
    test:  go test ./... -coverprofile=coverage.out

    # displays test coverage as a web page
    cov: go tool cover -html=coverage.out

    # runs the build binary/exe
    run:  ./${app_name} 

    # Runs different golang checks
    check: |
        go fmt $(go list ./... | grep -v /vendor/)
        go vet $(go list ./... | grep -v /vendor/)
        golint $(go list ./... | grep -v /vendor/)
        misspell $(go list ./... | grep -v /vendor/)
        gocyclo -top=10 .

    # Removes unwanted files
    clean: |
        rm -rf $(app_name) tests.log artifacts coverage.out
        go clean
        find . -name '*.test' -delete

    # starts required docker containers if not inside a docker env
    up: |
        if [ -f /.dockerenv ]; then
            docker-compose up -d
        fi

    # brings up the shell from your dockerr env
    shell: |
        if [ -f /.dockerenv ]; then
            docker exec -i -t $(app_name) /bin/bash
        fi
    
    # Destroys the docker stuff
    docker_clean: |
        if [ "$(docker ps -q -f name=$(app_name))" ]; then \
            docker rm -v -f $(docker ps -q -f name=$(app_name)); \
        fi
        docker network prune --force

`
