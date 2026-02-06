package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "5.0.0"
var StateSchemaVersion = 1

var rootCmd = &cobra.Command{
	Use:   "metacog",
	Short: "Metacognitive compositional engine",
}

var jsonOutput bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		output := fmt.Sprintf("metacog v%s\nstate schema: v%d\nstratagems: pivot mirror stack anchor reset", Version, StateSchemaVersion)
		fmt.Println(FormatOutput(jsonOutput, output, nil))
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
