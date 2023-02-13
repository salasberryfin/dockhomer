/*
Create a cobra CLI tool for interaction with Dockhomer

Only basic commands are implemented:
- Open an interactive shell to a container from image: dockhomer run <image>
- Run a single command in a container from image: dockhomer run command <image> <command>
*/
package main

import (
	"fmt"

	"github.com/salasberryfin/dockhomer/container"
	"github.com/spf13/cobra"
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "dockhomer",
		Short: "Dockhomer is a simple container creation tool",
	}
	var interactiveCmd = &cobra.Command{
		Use:   "run <image>",
		Short: "open an interactive shell to a container built from an image",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dockhomer run", args)
			cont := container.New(args[0], "/tmp/dockhomer")
			cont.OpenShell()
		},
	}
	var runCmd = &cobra.Command{
		Use:   "command <image> <cmd>",
		Short: "Run a command in a container built from an image",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dockhomer run with command", args)
			cont := container.New(args[0], "/tmp/dockhomer")
			cont.RunCmd(args[1:]...)
		},
	}

	rootCmd.AddCommand(interactiveCmd)
	interactiveCmd.AddCommand(runCmd)
	rootCmd.Execute()
}
