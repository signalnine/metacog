package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func FormatStatus(s *State) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Session: %s\n\n", s.SessionID))

	if s.Identity != nil {
		b.WriteString(fmt.Sprintf("Identity: %s\n  Lens: %s\n  Environment: %s\n\n", s.Identity.Name, s.Identity.Lens, s.Identity.Env))
	} else {
		b.WriteString("Identity: (none)\n\n")
	}

	if s.Substrate != nil {
		b.WriteString(fmt.Sprintf("Substrate: %s\n  Method: %s\n  Qualia: %s\n\n", s.Substrate.Substance, s.Substrate.Method, s.Substrate.Qualia))
	} else {
		b.WriteString("Substrate: (none)\n\n")
	}

	if s.Stratagem != nil {
		b.WriteString(StratagemStatus(s))
	} else {
		b.WriteString("Stratagem: (none)\n")
	}

	return b.String()
}

func FormatHistory(s *State) string {
	if len(s.History) == 0 {
		return "No history."
	}
	var b strings.Builder
	for i, h := range s.History {
		status := ""
		if h.Status != "" {
			status = fmt.Sprintf(" [%s]", h.Status)
		}
		b.WriteString(fmt.Sprintf("%d. [%s] %s%s", i+1, h.Timestamp, h.Action, status))
		if len(h.Params) > 0 {
			parts := make([]string, 0, len(h.Params))
			for k, v := range h.Params {
				parts = append(parts, fmt.Sprintf("%s=%s", k, v))
			}
			b.WriteString(fmt.Sprintf(" (%s)", strings.Join(parts, ", ")))
		}
		b.WriteString("\n")
	}
	return b.String()
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current state",
	RunE: func(cmd *cobra.Command, args []string) error {
		sm := DefaultStateManager()
		s, err := sm.Load()
		if err != nil {
			return err
		}
		fmt.Println(FormatOutput(jsonOutput, FormatStatus(s), nil))
		return nil
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Clear all state",
	RunE: func(cmd *cobra.Command, args []string) error {
		sm := DefaultStateManager()
		return sm.SaveWithLock(func(s *State) error {
			fresh := NewState()
			*s = *fresh
			return nil
		})
	},
}

var historyFull bool

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show transformation history",
	RunE: func(cmd *cobra.Command, args []string) error {
		sm := DefaultStateManager()
		s, err := sm.Load()
		if err != nil {
			return err
		}
		fmt.Println(FormatOutput(jsonOutput, FormatHistory(s), nil))
		return nil
	},
}

var repairCmd = &cobra.Command{
	Use:   "repair",
	Short: "Validate and fix corrupted state file",
	RunE: func(cmd *cobra.Command, args []string) error {
		sm := DefaultStateManager()
		err := sm.Repair()
		if err != nil {
			return err
		}
		fmt.Println("State file repaired.")
		return nil
	},
}

func init() {
	historyCmd.Flags().BoolVar(&historyFull, "full", false, "Show full history from log file")
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(historyCmd)
	rootCmd.AddCommand(repairCmd)
}
