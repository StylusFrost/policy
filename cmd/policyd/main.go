package main

import (
	"os"

	"github.com/StylusFrost/policy/app"
	"github.com/StylusFrost/policy/cmd/policyd/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
