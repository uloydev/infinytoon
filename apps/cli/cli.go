package main

import (
	"infinitoon.dev/infinitoon/apps/cli/cmd"
	appctx "infinitoon.dev/infinitoon/pkg/context"
)

func main() {
	cmd.Execute(appctx.NewAppContext())
}
