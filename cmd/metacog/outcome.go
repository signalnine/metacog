package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func findLastCompletedStratagem(s *State) (string, int) {
	for i := len(s.History) - 1; i >= 0; i-- {
		h := s.History[i]
		if h.Action == "stratagem" && h.Params["event"] == "completed" {
			return h.Params["name"], i
		}
	}
	return "", -1
}

func hasOutcomeAfter(s *State, afterIdx int) bool {
	for i := afterIdx + 1; i < len(s.History); i++ {
		if s.History[i].Action == "outcome" {
			return true
		}
	}
	return false
}

func RecordOutcome(s *State, result, shift string) error {
	if result != "productive" && result != "unproductive" {
		return fmt.Errorf("result must be 'productive' or 'unproductive', got %q", result)
	}

	name, idx := findLastCompletedStratagem(s)
	if idx < 0 {
		return fmt.Errorf("no completed stratagem found in history")
	}

	if hasOutcomeAfter(s, idx) {
		return fmt.Errorf("outcome already recorded for this stratagem. Use --amend to update")
	}

	params := map[string]string{
		"result":    result,
		"stratagem": name,
	}
	if shift != "" {
		params["shift"] = shift
	}

	s.AddHistory(HistoryEntry{
		Action: "outcome",
		Params: params,
	})
	return nil
}

func AmendOutcome(s *State, result, shift string) error {
	if result != "productive" && result != "unproductive" {
		return fmt.Errorf("result must be 'productive' or 'unproductive', got %q", result)
	}

	// Find most recent outcome
	for i := len(s.History) - 1; i >= 0; i-- {
		if s.History[i].Action == "outcome" {
			s.History[i].Params["result"] = result
			if shift != "" {
				s.History[i].Params["shift"] = shift
			} else {
				delete(s.History[i].Params, "shift")
			}
			return nil
		}
	}
	return fmt.Errorf("no outcome to amend")
}

var outcomeResult string
var outcomeShift string
var outcomeAmend bool

var outcomeCmd = &cobra.Command{
	Use:   "outcome",
	Short: "Record stratagem effectiveness",
	RunE: func(cmd *cobra.Command, args []string) error {
		sm := DefaultStateManager()
		var output string
		err := sm.SaveWithLock(func(s *State) error {
			if outcomeAmend {
				err := AmendOutcome(s, outcomeResult, outcomeShift)
				if err != nil {
					return err
				}
				output = fmt.Sprintf("Outcome amended to %s.", outcomeResult)
				return nil
			}

			err := RecordOutcome(s, outcomeResult, outcomeShift)
			if err != nil {
				return err
			}
			name, _ := findLastCompletedStratagem(s)
			output = fmt.Sprintf("Outcome recorded: %s (%s).", outcomeResult, name)
			return nil
		})
		if err != nil {
			return err
		}
		fmt.Println(FormatOutput(jsonOutput, output, nil))
		return nil
	},
}

func init() {
	outcomeCmd.Flags().StringVar(&outcomeResult, "result", "", "Outcome: productive or unproductive (required)")
	outcomeCmd.Flags().StringVar(&outcomeShift, "shift", "", "Description of what changed (optional)")
	outcomeCmd.Flags().BoolVar(&outcomeAmend, "amend", false, "Update most recent outcome instead of creating new")
	outcomeCmd.MarkFlagRequired("result")
	rootCmd.AddCommand(outcomeCmd)
}
