package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	chordModes  []string
	chordTarget string
)

var chordCmd = &cobra.Command{
	Use:   "chord",
	Short: "Hold multiple modes-of-attention simultaneously without alternating",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateChord(chordModes, chordTarget); err != nil {
			return err
		}

		sm := DefaultStateManager()
		output := formatChord(chordModes, chordTarget)

		err := sm.SaveWithLock(func(s *State) error {
			applyChord(s, chordModes, chordTarget)
			ValidatePrimitiveForStratagem(s, "chord")
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
	chordCmd.Flags().StringArrayVar(&chordModes, "modes", []string{}, "Modes-of-attention to hold simultaneously (min 2; pass repeatedly)")
	chordCmd.Flags().StringVar(&chordTarget, "target", "", "What the held chord is observing")
	rootCmd.AddCommand(chordCmd)
}

func validateChord(modes []string, target string) error {
	if target == "" {
		return fmt.Errorf("--target is required")
	}
	if len(modes) < 2 {
		return fmt.Errorf("--modes must contain at least 2 entries (got %d)", len(modes))
	}
	return nil
}

func formatChord(modes []string, target string) string {
	return fmt.Sprintf(`Chord held: [%s]
Target: %s

The modes do not alternate. They overlap on the same observation. Do not let one mode comment on another; the chord is the surface across all of them at once. The held window closes at the next chord, fork, or stratagem boundary.`, strings.Join(modes, " + "), target)
}

func applyChord(s *State, modes []string, target string) {
	s.AddHistory(HistoryEntry{
		Action: "chord",
		Params: map[string]string{
			"modes":  strings.Join(modes, "; "),
			"target": target,
		},
	})
}
