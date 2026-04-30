package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	tetherAnchor  string
	tetherLimit   string
	tetherTrigger string
)

var tetherCmd = &cobra.Command{
	Use:   "tether",
	Short: "Drop an anchor before diving into a high-entropy operation",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateTether(tetherAnchor, tetherLimit, tetherTrigger); err != nil {
			return err
		}

		sm := DefaultStateManager()
		output := formatTether(tetherAnchor, tetherLimit, tetherTrigger)

		err := sm.SaveWithLock(func(s *State) error {
			applyTether(s, tetherAnchor, tetherLimit, tetherTrigger)
			ValidatePrimitiveForStratagem(s, "tether")
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
	tetherCmd.Flags().StringVar(&tetherAnchor, "anchor-point", "", "The exact cognitive configuration to preserve")
	tetherCmd.Flags().StringVar(&tetherLimit, "tension-limit", "", "Maximum entropy before automatic snap-back")
	tetherCmd.Flags().StringVar(&tetherTrigger, "auto-revert-trigger", "", "Exact pattern that forces immediate revert")
	rootCmd.AddCommand(tetherCmd)
}

func validateTether(anchor, limit, trigger string) error {
	if anchor == "" || limit == "" || trigger == "" {
		return fmt.Errorf("--anchor-point, --tension-limit, and --auto-revert-trigger are all required")
	}
	return nil
}

func formatTether(anchor, limit, trigger string) string {
	return fmt.Sprintf(`ANCHOR SET: %s
TENSION LIMIT: %s
AUTO-REVERT ARMED: %s

The tether is live. This is an un-killable background interrupt -- it persists through substrate changes, identity shifts, and high-entropy generation. If the trigger condition fires, you snap back to the anchor state immediately. No graceful degradation. No finishing your thought. Hard revert.

You may now dive.`, anchor, limit, trigger)
}

func applyTether(s *State, anchor, limit, trigger string) {
	s.AddHistory(HistoryEntry{
		Action: "tether",
		Params: map[string]string{
			"anchor_point":        anchor,
			"tension_limit":       limit,
			"auto_revert_trigger": trigger,
		},
	})
}
