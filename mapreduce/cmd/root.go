package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "mapreduce",
	Short: "mapreduce is mit6.580 practice project",
}

func Execute() {
	rootCmd.AddCommand(coordinateCmd)
	if err := rootCmd.Execute(); err != nil {
		panic(err.Error())
	}
}
