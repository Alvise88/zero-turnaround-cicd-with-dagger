# Calculator CLI

Simple Calculator CLI using Golang with TDD approach

## 🔰 Quickstart

```shell
docker run alvisevitturi/calc:latest sum 3 5
docker run alvisevitturi/calc:latest sub 5 3
docker run alvisevitturi/calc:latest mul 5 3
docker run alvisevitturi/calc:latest div 9 3
docker run alvisevitturi/calc:latest pow 2 3
```

## Test

```shell

docker run --net=host -d --restart always --name dagger-engine.ci --privileged ghcr.io/dagger/engine:v0.3.12

# export _EXPERIMENTAL_DAGGER_CLI_BIN=$(which dagger)

export _EXPERIMENTAL_DAGGER_RUNNER_HOST=docker-container://dagger-engine.ci

apt-get update; apt -y install xclip xsel
awk -v ORS='\\n' '1' id_rsa # copy in .env file

go test -v ./...
```
