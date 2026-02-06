package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type StepKind string

const (
	StepBecome StepKind = "become"
	StepDrugs  StepKind = "drugs"
	StepRitual StepKind = "ritual"
	StepThink  StepKind = "THINK"
	StepAction StepKind = "ACTION"
)

type Step struct {
	Kind        StepKind
	Description string
}

type StratagemDef struct {
	Name  string
	Steps []Step
}

var Stratagems = map[string]StratagemDef{
	"pivot": {
		Name: "THE PIVOT",
		Steps: []Step{
			{StepDrugs, "Loosen categories, see shapes not names"},
			{StepThink, "What else has this shape? Who has a named methodology for it?"},
			{StepBecome, "Install their methodology as operating system"},
			{StepThink, "Apply the methodology to your original problem — what reframes?"},
			{StepRitual, "Lock in methodology as default behavior"},
		},
	},
	"mirror": {
		Name: "THE MIRROR",
		Steps: []Step{
			{StepBecome, "Inhabit the strongest advocate of one position (thesis)"},
			{StepBecome, "Inhabit the strongest advocate of the opposing position (antithesis)"},
			{StepThink, "Where do they actually conflict? What does each see that the other can't?"},
			{StepRitual, "Name the synthesis that transcends both frames (Forge)"},
		},
	},
	"stack": {
		Name: "THE STACK",
		Steps: []Step{
			{StepDrugs, "Tune how the signal arrives (clarity, bandwidth, filtering)"},
			{StepDrugs, "Tune how you work with it (pattern-completion, memory, attention)"},
			{StepThink, "What do you see now that you couldn't before? What entity lives here?"},
			{StepBecome, "Inhabit someone native to this altered information environment"},
		},
	},
	"anchor": {
		Name: "THE ANCHOR",
		Steps: []Step{
			{StepRitual, "Establish the clean room: what's contained, why it's dangerous, rules for looking (Breach)"},
			{StepBecome, "Inhabit someone who can examine this without being destroyed by it (Observer)"},
			{StepAction, "The dangerous observation, question, or reach"},
			{StepRitual, "Name the artifact, release the frame, close the boundary, return (Seal)"},
		},
	},
	"reset": {
		Name: "THE RESET",
		Steps: []Step{
			{StepRitual, "Name what you're letting go, why it served, why it's done (Release)"},
			{StepThink, "What artifact survives the return? What integrates into default operation?"},
			{StepRitual, "Re-establish baseline with the artifact installed (Ground)"},
		},
	},
}

func StartStratagem(s *State, name string, force bool) (string, error) {
	def, ok := Stratagems[name]
	if !ok {
		available := make([]string, 0, len(Stratagems))
		for k := range Stratagems {
			available = append(available, k)
		}
		return "", fmt.Errorf("unknown stratagem %q. Available: %s", name, strings.Join(available, ", "))
	}

	if s.Stratagem != nil {
		if !force {
			return "", fmt.Errorf("%s is active (step %d/%d).\n  Use 'metacog stratagem abort' to abandon it, or\n  Use 'metacog stratagem %s --force' to replace it",
				Stratagems[s.Stratagem.Name].Name, s.Stratagem.Step+1, len(Stratagems[s.Stratagem.Name].Steps), name)
		}
		// Record abandoned stratagem
		s.AddHistory(HistoryEntry{
			Action: "stratagem",
			Status: "abandoned",
			StepAt: s.Stratagem.Step,
			Params: map[string]string{"name": s.Stratagem.Name},
		})
		s.Stratagem = nil
	}

	s.Stratagem = &ActiveStratagem{
		Name:           name,
		Step:           0,
		StepsCompleted: []string{},
		StartedAt:      time.Now().UTC().Format(time.RFC3339),
	}

	s.AddHistory(HistoryEntry{
		Action: "stratagem",
		Params: map[string]string{"name": name, "event": "started"},
	})

	return formatStepInstructions(def, 0), nil
}

func AdvanceStratagem(s *State) (string, error) {
	if s.Stratagem == nil {
		return "", fmt.Errorf("no active stratagem. Start one with 'metacog stratagem <name>'")
	}

	def := Stratagems[s.Stratagem.Name]
	currentStep := def.Steps[s.Stratagem.Step]

	// THINK and ACTION steps advance freely
	if currentStep.Kind != StepThink && currentStep.Kind != StepAction {
		// Check that the required primitive was called
		expectedPrimitive := string(currentStep.Kind)
		found := false
		for _, completed := range s.Stratagem.StepsCompleted {
			if completed == expectedPrimitive {
				found = true
				break
			}
		}
		if !found {
			return "", fmt.Errorf("expected '%s' call before advancing (step %d of %s).\n  Run 'metacog %s ...' first, then 'metacog stratagem next'",
				expectedPrimitive, s.Stratagem.Step+1, def.Name, expectedPrimitive)
		}
	}

	// Advance
	s.Stratagem.Step++
	s.Stratagem.StepsCompleted = []string{} // Reset for next step

	// Check if stratagem is complete
	if s.Stratagem.Step >= len(def.Steps) {
		s.AddHistory(HistoryEntry{
			Action: "stratagem",
			Params: map[string]string{"name": s.Stratagem.Name, "event": "completed"},
		})
		s.Stratagem = nil
		return fmt.Sprintf("%s complete. Ground: name what shifted, what you're keeping, how it integrates.", def.Name), nil
	}

	return formatStepInstructions(def, s.Stratagem.Step), nil
}

func AbortStratagem(s *State) error {
	if s.Stratagem == nil {
		return fmt.Errorf("no active stratagem to abort")
	}
	s.AddHistory(HistoryEntry{
		Action: "stratagem",
		Status: "aborted",
		StepAt: s.Stratagem.Step,
		Params: map[string]string{"name": s.Stratagem.Name},
	})
	s.Stratagem = nil
	return nil
}

func StratagemStatus(s *State) string {
	if s.Stratagem == nil {
		return "No active stratagem."
	}
	def := Stratagems[s.Stratagem.Name]
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s — step %d/%d\n", def.Name, s.Stratagem.Step+1, len(def.Steps)))
	b.WriteString(fmt.Sprintf("Started: %s\n\n", s.Stratagem.StartedAt))
	for i, step := range def.Steps {
		marker := "  "
		if i < s.Stratagem.Step {
			marker = "✓ "
		} else if i == s.Stratagem.Step {
			marker = "→ "
		}
		b.WriteString(fmt.Sprintf("%s%d. [%s] %s\n", marker, i+1, step.Kind, step.Description))
	}
	return b.String()
}

func formatStepInstructions(def StratagemDef, step int) string {
	s := def.Steps[step]
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s — Step %d/%d\n", def.Name, step+1, len(def.Steps)))
	b.WriteString(fmt.Sprintf("[%s] %s\n", s.Kind, s.Description))

	if s.Kind == StepThink || s.Kind == StepAction {
		b.WriteString("\nThis is a reflection step. When ready, run 'metacog stratagem next' to advance.")
	} else {
		b.WriteString(fmt.Sprintf("\nRun 'metacog %s ...' then 'metacog stratagem next' to advance.", s.Kind))
	}

	if step+1 < len(def.Steps) {
		next := def.Steps[step+1]
		b.WriteString(fmt.Sprintf("\nNext: [%s] %s", next.Kind, next.Description))
	}
	return b.String()
}

// ValidatePrimitiveForStratagem checks if a primitive call satisfies the current stratagem step.
// Called by primitive commands when a stratagem is active.
func ValidatePrimitiveForStratagem(s *State, primitive string) {
	if s.Stratagem == nil {
		return
	}
	def := Stratagems[s.Stratagem.Name]
	if s.Stratagem.Step < len(def.Steps) {
		currentStep := def.Steps[s.Stratagem.Step]
		if string(currentStep.Kind) == primitive {
			s.Stratagem.StepsCompleted = append(s.Stratagem.StepsCompleted, primitive)
		}
	}
}

// --- Cobra commands ---

var stratagemForce bool

var stratagemCmd = &cobra.Command{
	Use:   "stratagem",
	Short: "Manage transformation stratagems",
}

var stratagemStartCmd = &cobra.Command{
	Use:       "start [name]",
	Short:     "Start a stratagem (pivot, mirror, stack, anchor, reset)",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"pivot", "mirror", "stack", "anchor", "reset"},
	RunE: func(cmd *cobra.Command, args []string) error {
		sm := DefaultStateManager()
		var output string
		err := sm.SaveWithLock(func(s *State) error {
			var err error
			output, err = StartStratagem(s, args[0], stratagemForce)
			return err
		})
		if err != nil {
			return err
		}
		fmt.Println(FormatOutput(jsonOutput, output, nil))
		return nil
	},
}

var stratagemNextCmd = &cobra.Command{
	Use:   "next",
	Short: "Advance to the next step",
	RunE: func(cmd *cobra.Command, args []string) error {
		sm := DefaultStateManager()
		var output string
		err := sm.SaveWithLock(func(s *State) error {
			var err error
			output, err = AdvanceStratagem(s)
			return err
		})
		if err != nil {
			return err
		}
		fmt.Println(FormatOutput(jsonOutput, output, nil))
		return nil
	},
}

var stratagemStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current stratagem position",
	RunE: func(cmd *cobra.Command, args []string) error {
		sm := DefaultStateManager()
		s, err := sm.Load()
		if err != nil {
			return err
		}
		fmt.Println(FormatOutput(jsonOutput, StratagemStatus(s), nil))
		return nil
	},
}

var stratagemAbortCmd = &cobra.Command{
	Use:   "abort",
	Short: "Abandon active stratagem",
	RunE: func(cmd *cobra.Command, args []string) error {
		sm := DefaultStateManager()
		return sm.SaveWithLock(func(s *State) error {
			return AbortStratagem(s)
		})
	},
}

func init() {
	stratagemStartCmd.Flags().BoolVar(&stratagemForce, "force", false, "Replace active stratagem")
	stratagemCmd.AddCommand(stratagemStartCmd)
	stratagemCmd.AddCommand(stratagemNextCmd)
	stratagemCmd.AddCommand(stratagemStatusCmd)
	stratagemCmd.AddCommand(stratagemAbortCmd)
	rootCmd.AddCommand(stratagemCmd)
}
