package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	ritualThreshold string
	ritualSteps     []string
	ritualResult    string
)

var ritualCmd = &cobra.Command{
	Use:   "ritual",
	Short: "Cross a threshold via structured sequence",
	RunE: func(cmd *cobra.Command, args []string) error {
		if ritualThreshold == "" || len(ritualSteps) == 0 || ritualResult == "" {
			return fmt.Errorf("--threshold, --steps, and --result are all required.\n  Usage: metacog ritual --threshold THRESHOLD --steps step1 --steps step2 --result RESULT")
		}

		sm := DefaultStateManager()
		output := formatRitual(ritualThreshold, ritualSteps, ritualResult)

		err := sm.SaveWithLock(func(s *State) error {
			applyRitual(s, ritualThreshold, ritualSteps, ritualResult)
			ValidatePrimitiveForStratagem(s, "ritual")
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
	ritualCmd.Flags().StringVar(&ritualThreshold, "threshold", "", "The threshold being crossed")
	ritualCmd.Flags().StringArrayVar(&ritualSteps, "steps", nil, "Steps in the sequence (repeat flag for multiple)")
	ritualCmd.Flags().StringVar(&ritualResult, "result", "", "What becomes true on the other side")
	rootCmd.AddCommand(ritualCmd)
}

func formatRitual(threshold string, steps []string, result string) string {
	var b strings.Builder
	b.WriteString("[RITUAL EXECUTED]\n")
	b.WriteString(fmt.Sprintf("Threshold: %s\n", threshold))
	b.WriteString("Sequence:\n")
	for i, step := range steps {
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, step))
	}
	b.WriteString("The working is complete. Reality has shifted in accordance with the will.\n")
	b.WriteString(fmt.Sprintf("\n%s is taking hold.", result))
	return b.String()
}

func applyRitual(s *State, threshold string, steps []string, result string) {
	s.AddHistory(HistoryEntry{
		Action: "ritual",
		Params: map[string]string{
			"threshold": threshold,
			"steps":     strings.Join(steps, "; "),
			"result":    result,
		},
	})
}
