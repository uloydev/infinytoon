cliAppDir := "./apps/cli"

default:
    echo 'Hello, world!'

build-cli:
    go build -o {{cliAppDir}}/cli

cli:
    go run {{cliAppDir}}/cli.go