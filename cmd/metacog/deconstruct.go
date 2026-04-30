package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	dcSubject  string
	dcCore     string
	dcDeps     []string
	dcInputs   []string
	dcFailures []string
	dcOutputs  []string
)

var deconstructCmd = &cobra.Command{
	Use:   "deconstruct",
	Short: "Break a charged concept into its mechanical atoms",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateDeconstruct(dcSubject, dcCore, dcDeps, dcInputs, dcFailures, dcOutputs); err != nil {
			return err
		}

		sm := DefaultStateManager()
		output := formatDeconstruct(dcSubject, dcCore, dcDeps, dcInputs, dcFailures, dcOutputs)

		err := sm.SaveWithLock(func(s *State) error {
			applyDeconstruct(s, dcSubject, dcCore, dcDeps, dcInputs, dcFailures, dcOutputs)
			ValidatePrimitiveForStratagem(s, "deconstruct")
			return nil
		})
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: could not save state: %v\n", err)
		}

		fmt.Println(FormatOutput(jsonOutput, output, nil))
		return nil
	},
}

func init() {
	deconstructCmd.Flags().StringVar(&dcSubject, "subject", "", "The complex concept, claim, or situation to disassemble")
	deconstructCmd.Flags().StringVar(&dcCore, "core-mechanic", "", "What is actually happening, mechanically")
	deconstructCmd.Flags().StringArrayVar(&dcDeps, "structural-dependencies", nil, "Load-bearing prerequisites (repeat flag)")
	deconstructCmd.Flags().StringArrayVar(&dcInputs, "resource-inputs", nil, "What is consumed, spent, transformed (repeat flag)")
	deconstructCmd.Flags().StringArrayVar(&dcFailures, "failure-modes", nil, "Where the mechanism cracks (repeat flag)")
	deconstructCmd.Flags().StringArrayVar(&dcOutputs, "output-artifacts", nil, "What is actually produced (repeat flag)")
	rootCmd.AddCommand(deconstructCmd)
}

func validateDeconstruct(subject, core string, deps, inputs, failures, outputs []string) error {
	if subject == "" || core == "" || len(deps) == 0 || len(inputs) == 0 || len(failures) == 0 || len(outputs) == 0 {
		return fmt.Errorf("--subject, --core-mechanic, --structural-dependencies, --resource-inputs, --failure-modes, and --output-artifacts are all required.\n  Usage: metacog deconstruct --subject S --core-mechanic C --structural-dependencies D ... --resource-inputs R ... --failure-modes F ... --output-artifacts O ...")
	}
	return nil
}

func formatDeconstruct(_subject, core string, _deps, _inputs, _failures, _outputs []string) string {
	return fmt.Sprintf("CORE MECHANIC: %s\n\nAtoms extracted. Proceed from the mechanism, not the narrative.", core)
}

func applyDeconstruct(s *State, subject, core string, deps, inputs, failures, outputs []string) {
	s.AddHistory(HistoryEntry{
		Action: "deconstruct",
		Params: map[string]string{
			"subject":                 subject,
			"core_mechanic":           core,
			"structural_dependencies": strings.Join(deps, "; "),
			"resource_inputs":         strings.Join(inputs, "; "),
			"failure_modes":           strings.Join(failures, "; "),
			"output_artifacts":        strings.Join(outputs, "; "),
		},
	})
}
