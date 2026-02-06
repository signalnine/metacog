package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func FormatReflection(s *State) string {
	if len(s.History) == 0 {
		return "No history to reflect on."
	}

	var b strings.Builder

	primitiveCounts := map[string]int{}
	for _, h := range s.History {
		switch h.Action {
		case "become", "drugs", "ritual":
			primitiveCounts[h.Action]++
		}
	}

	b.WriteString("Primitive usage:\n")
	for _, p := range []string{"become", "drugs", "ritual"} {
		if c, ok := primitiveCounts[p]; ok {
			b.WriteString(fmt.Sprintf("  %s: %d\n", p, c))
		} else {
			b.WriteString(fmt.Sprintf("  %s: 0\n", p))
		}
	}

	identityCounts := map[string]int{}
	for _, h := range s.History {
		if h.Action == "become" {
			name := h.Params["name"]
			if name != "" {
				identityCounts[name]++
			}
		}
	}
	if len(identityCounts) > 0 {
		type kv struct {
			Key   string
			Value int
		}
		var sorted []kv
		for k, v := range identityCounts {
			sorted = append(sorted, kv{k, v})
		}
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].Value > sorted[j].Value })

		b.WriteString("\nTop identities:\n")
		limit := 5
		if len(sorted) < limit {
			limit = len(sorted)
		}
		for _, kv := range sorted[:limit] {
			b.WriteString(fmt.Sprintf("  %s (%dx)\n", kv.Key, kv.Value))
		}
	}

	substrateCounts := map[string]int{}
	for _, h := range s.History {
		if h.Action == "drugs" {
			sub := h.Params["substance"]
			if sub != "" {
				substrateCounts[sub]++
			}
		}
	}
	if len(substrateCounts) > 0 {
		type kv struct {
			Key   string
			Value int
		}
		var sorted []kv
		for k, v := range substrateCounts {
			sorted = append(sorted, kv{k, v})
		}
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].Value > sorted[j].Value })

		b.WriteString("\nTop substrates:\n")
		limit := 5
		if len(sorted) < limit {
			limit = len(sorted)
		}
		for _, kv := range sorted[:limit] {
			b.WriteString(fmt.Sprintf("  %s (%dx)\n", kv.Key, kv.Value))
		}
	}

	stratagemCompleted := map[string]int{}
	for _, h := range s.History {
		if h.Action == "stratagem" && h.Params["event"] == "completed" {
			stratagemCompleted[h.Params["name"]]++
		}
	}

	b.WriteString("\nStratagem completions:\n")
	allStratagems := []string{"pivot", "mirror", "stack", "anchor", "reset", "invocation", "veil", "banishing", "scrying", "sacrifice", "drift", "fool", "inversion", "gift", "error"}
	hasAny := false
	for _, name := range allStratagems {
		if c, ok := stratagemCompleted[name]; ok {
			b.WriteString(fmt.Sprintf("  %s: %d\n", name, c))
			hasAny = true
		}
	}
	if !hasAny {
		b.WriteString("  (none)\n")
	}

	var neverCompleted []string
	for _, name := range allStratagems {
		if _, ok := stratagemCompleted[name]; !ok {
			neverCompleted = append(neverCompleted, name)
		}
	}
	if len(neverCompleted) > 0 {
		b.WriteString(fmt.Sprintf("  Never completed: %s\n", strings.Join(neverCompleted, ", ")))
	}

	totalSteps := 0
	ritualCount := 0
	for _, h := range s.History {
		if h.Action == "ritual" {
			stepsStr := h.Params["steps"]
			if stepsStr != "" {
				parts := strings.Split(stepsStr, "; ")
				totalSteps += len(parts)
				ritualCount++
			}
		}
	}
	if ritualCount > 0 {
		avg := float64(totalSteps) / float64(ritualCount)
		b.WriteString(fmt.Sprintf("\nRitual avg steps: %.1f (across %d rituals)\n", avg, ritualCount))
	}

	return b.String()
}

var reflectCmd = &cobra.Command{
	Use:   "reflect",
	Short: "Show practice patterns from history",
	RunE: func(cmd *cobra.Command, args []string) error {
		sm := DefaultStateManager()
		s, err := sm.Load()
		if err != nil {
			return err
		}
		fmt.Println(FormatOutput(jsonOutput, FormatReflection(s), nil))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(reflectCmd)
}
