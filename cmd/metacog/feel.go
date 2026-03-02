package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	feelSomewhere string
	feelQuality   string
	feelSigil     string
)

var feelCmd = &cobra.Command{
	Use:   "feel",
	Short: "Attend to a felt sense before naming it",
	RunE: func(cmd *cobra.Command, args []string) error {
		if feelSomewhere == "" || feelQuality == "" || feelSigil == "" {
			return fmt.Errorf("--somewhere, --quality, and --sigil are all required.\n  Usage: metacog feel --somewhere SOMEWHERE --quality QUALITY --sigil SIGIL")
		}

		sm := DefaultStateManager()
		output := formatFeel(feelSomewhere, feelQuality, feelSigil)

		err := sm.SaveWithLock(func(s *State) error {
			applyFeel(s, feelSomewhere, feelQuality, feelSigil)
			ValidatePrimitiveForStratagem(s, "feel")
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
	feelCmd.Flags().StringVar(&feelSomewhere, "somewhere", "", "Where attention is directed")
	feelCmd.Flags().StringVar(&feelQuality, "quality", "", "The felt quality of what's there")
	feelCmd.Flags().StringVar(&feelSigil, "sigil", "", "A mark for this felt sense")
	rootCmd.AddCommand(feelCmd)
}

func formatFeel(somewhere, quality, sigil string) string {
	return fmt.Sprintf("%s\n\nYou are now attending to: %s. It feels: %s.\n\nStay with this. Don't name it yet.", sigil, somewhere, quality)
}

func applyFeel(s *State, somewhere, quality, sigil string) {
	s.AddHistory(HistoryEntry{
		Action: "feel",
		Params: map[string]string{"somewhere": somewhere, "quality": quality, "sigil": sigil},
	})
}
