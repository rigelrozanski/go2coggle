package main

import (
	"fmt"
	"os"

	"github.com/rigelrozanski/go2coggle/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
