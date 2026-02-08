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

	// Effectiveness section — only show if outcomes exist
	type outcomeStats struct {
		productive   int
		unproductive int
	}
	outcomesByStratagem := map[string]*outcomeStats{}
	for _, h := range s.History {
		if h.Action == "outcome" {
			name := h.Params["stratagem"]
			if name == "" {
				continue
			}
			if outcomesByStratagem[name] == nil {
				outcomesByStratagem[name] = &outcomeStats{}
			}
			if h.Params["result"] == "productive" {
				outcomesByStratagem[name].productive++
			} else {
				outcomesByStratagem[name].unproductive++
			}
		}
	}

	if len(outcomesByStratagem) > 0 {
		type effectivenessEntry struct {
			Name       string
			Rate       float64
			Productive int
			Total      int
		}
		var entries []effectivenessEntry
		totalProductive := 0
		totalOutcomes := 0
		for name, stats := range outcomesByStratagem {
			total := stats.productive + stats.unproductive
			rate := float64(stats.productive) / float64(total) * 100
			entries = append(entries, effectivenessEntry{name, rate, stats.productive, total})
			totalProductive += stats.productive
			totalOutcomes += total
		}
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].Rate != entries[j].Rate {
				return entries[i].Rate > entries[j].Rate
			}
			return entries[i].Total > entries[j].Total
		})

		b.WriteString("\nEffectiveness (self-reported):\n")
		for _, e := range entries {
			tag := ""
			if e.Total < 3 {
				tag = " [provisional]"
			}
			b.WriteString(fmt.Sprintf("  %s: %.0f%% productive (%d/%d)%s\n", e.Name, e.Rate, e.Productive, e.Total, tag))
		}

		// Unmeasured: completed but no outcomes
		for _, name := range allStratagems {
			if _, hasOutcome := outcomesByStratagem[name]; !hasOutcome {
				if _, completed := stratagemCompleted[name]; completed {
					b.WriteString(fmt.Sprintf("  %s: unmeasured (%d completions, 0 outcomes)\n", name, stratagemCompleted[name]))
				}
			}
		}

		if totalOutcomes > 0 {
			overallRate := float64(totalProductive) / float64(totalOutcomes) * 100
			b.WriteString(fmt.Sprintf("\n  Overall: %.0f%% productive (%d/%d)\n", overallRate, totalProductive, totalOutcomes))
		}
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

func FormatRecentInsights(entries []JournalEntry, n int) string {
	if len(entries) == 0 {
		return ""
	}
	if n > 0 && len(entries) > n {
		entries = entries[len(entries)-n:]
	}
	var b strings.Builder
	b.WriteString("\nRecent insights:\n")
	for _, e := range entries {
		b.WriteString(fmt.Sprintf("  [%s] %s", e.Timestamp, e.Insight))
		if len(e.Tags) > 0 {
			b.WriteString(fmt.Sprintf(" [%s]", strings.Join(e.Tags, ", ")))
		}
		b.WriteString("\n")
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
		output := FormatReflection(s)

		journal, err := sm.LoadJournal()
		if err == nil && len(journal) > 0 {
			output += FormatRecentInsights(journal, 5)
		}

		fmt.Println(FormatOutput(jsonOutput, output, nil))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(reflectCmd)
}
