package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	milkywayapp "github.com/milkyway-labs/milkyway/v7/app"
	"github.com/milkyway-labs/milkyway/v7/cmd/milkywayd/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, milkywayapp.EnvPrefix, milkywayapp.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
