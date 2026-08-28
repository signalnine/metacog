package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var Version = "6.11.0"
var StateSchemaVersion = 1

var rootCmd = &cobra.Command{
	Use:   "metacog",
	Short: "Metacognitive compositional engine",
}

var jsonOutput bool

// formatVersion derives the primitive and stratagem lists from the registries
// so they cannot drift from what the binary actually ships (registry_test.go).
func formatVersion() string {
	return fmt.Sprintf("metacog v%s\nstate schema: v%d\nprimitives: %s\nstratagems: %s",
		Version, StateSchemaVersion,
		strings.Join(PrimitiveNames(), " "),
		strings.Join(allStratagemNames(), " "))
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(FormatOutput(jsonOutput, formatVersion(), nil))
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
