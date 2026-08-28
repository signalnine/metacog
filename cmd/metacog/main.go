package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var Version = "6.11.0"
var StateSchemaVersion = 1

var rootCmd = &cobra.Command{
	Use:               "metacog",
	Short:             "Metacognitive compositional engine",
	SilenceUsage:      true,
	PersistentPreRunE: rejectBlankStringFlags,
}

func cobraTraverseRunHooksEnabled() bool { return cobra.EnableTraverseRunHooks }

// rejectBlankStringFlags fails any string flag that was set to whitespace
// only ("--name '  '"). Values are never altered: an explicit empty string
// (meditate --focus "") stays legal, and inner/leading whitespace in a real
// value (an excerpt fragment) is preserved verbatim.
func rejectBlankStringFlags(cmd *cobra.Command, args []string) error {
	var blank []string
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed && f.Value.Type() == "string" {
			v := f.Value.String()
			if v != "" && strings.TrimSpace(v) == "" {
				blank = append(blank, "--"+f.Name)
			}
		}
	})
	if len(blank) > 0 {
		return fmt.Errorf("%s must not be blank (whitespace only)", strings.Join(blank, ", "))
	}
	return nil
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
	// Run root's PersistentPreRunE even if a subcommand later defines its own
	// hook (cobra otherwise runs only the nearest one). main_test.go guards this.
	cobra.EnableTraverseRunHooks = true
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
