package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

//go:embed stances/*.json
var stancesFS embed.FS

type Stance struct {
	Who   string `json:"who"`
	Where string `json:"where"`
	Lens  string `json:"lens"`
}

type StancePool struct {
	Name    string
	Stances []Stance
}

func LoadStancePools() (map[string]StancePool, error) {
	pools := make(map[string]StancePool)

	entries, err := stancesFS.ReadDir("stances")
	if err != nil {
		return nil, fmt.Errorf("reading stances directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		data, err := stancesFS.ReadFile("stances/" + entry.Name())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not read %s: %v\n", entry.Name(), err)
			continue
		}

		var stances []Stance
		if err := json.Unmarshal(data, &stances); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not parse %s: %v\n", entry.Name(), err)
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".json")
		pools[name] = StancePool{
			Name:    name,
			Stances: stances,
		}
	}

	return pools, nil
}

func ListPoolNames(pools map[string]StancePool) []string {
	names := make([]string, 0, len(pools))
	for name := range pools {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func RandomStance(pools map[string]StancePool, poolName string) (*Stance, string, error) {
	if poolName != "" {
		pool, ok := pools[poolName]
		if !ok {
			return nil, "", fmt.Errorf("unknown pool %q. Use --list to see available pools", poolName)
		}
		s := pool.Stances[rand.Intn(len(pool.Stances))]
		return &s, poolName, nil
	}

	// Pick a random pool, then a random stance from it
	names := ListPoolNames(pools)
	if len(names) == 0 {
		return nil, "", fmt.Errorf("no stance pools loaded")
	}
	chosen := names[rand.Intn(len(names))]
	pool := pools[chosen]
	s := pool.Stances[rand.Intn(len(pool.Stances))]
	return &s, chosen, nil
}

type PersonalStance struct {
	Who       string `json:"who"`
	Where     string `json:"where"`
	Lens      string `json:"lens"`
	Substance string `json:"substance,omitempty"`
	Method    string `json:"method,omitempty"`
	Qualia    string `json:"qualia,omitempty"`
}

func SavePersonalStance(metacogDir string, s *State) error {
	if s.Identity == nil {
		return fmt.Errorf("no identity set. Use 'metacog become' first")
	}

	stance := PersonalStance{
		Who:   s.Identity.Name,
		Where: s.Identity.Env,
		Lens:  s.Identity.Lens,
	}
	if s.Substrate != nil {
		stance.Substance = s.Substrate.Substance
		stance.Method = s.Substrate.Method
		stance.Qualia = s.Substrate.Qualia
	}

	stancesDir := filepath.Join(metacogDir, "stances")
	os.MkdirAll(stancesDir, 0755)
	poolPath := filepath.Join(stancesDir, "personal.json")
	lockPath := filepath.Join(stancesDir, ".personal.lock")

	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("cannot open personal pool lock: %w", err)
	}
	defer func() {
		syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)
		lockFile.Close()
	}()
	if err := syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX); err != nil {
		return fmt.Errorf("cannot acquire personal pool lock: %w", err)
	}

	var stances []PersonalStance
	if data, err := os.ReadFile(poolPath); err == nil {
		json.Unmarshal(data, &stances)
	}

	for _, existing := range stances {
		if existing.Who == stance.Who && existing.Where == stance.Where && existing.Lens == stance.Lens {
			return nil
		}
	}

	stances = append(stances, stance)

	data, err := json.MarshalIndent(stances, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal stances: %w", err)
	}

	return os.WriteFile(poolPath, data, 0644)
}

func LoadStancePoolsWithPersonal(metacogDir string) (map[string]StancePool, error) {
	pools, err := LoadStancePools()
	if err != nil {
		return nil, err
	}

	poolPath := filepath.Join(metacogDir, "stances", "personal.json")
	if data, err := os.ReadFile(poolPath); err == nil {
		var personalStances []PersonalStance
		if err := json.Unmarshal(data, &personalStances); err == nil && len(personalStances) > 0 {
			stances := make([]Stance, len(personalStances))
			for i, ps := range personalStances {
				stances[i] = Stance{Who: ps.Who, Where: ps.Where, Lens: ps.Lens}
			}
			pools["personal"] = StancePool{Name: "personal", Stances: stances}
		}
	}

	return pools, nil
}

var inspirePoolName string
var inspireList bool
var inspireSave bool

var inspireCmd = &cobra.Command{
	Use:   "inspire",
	Short: "Draw a random stance from the pool",
	RunE: func(cmd *cobra.Command, args []string) error {
		sm := DefaultStateManager()

		if inspireSave {
			s, err := sm.Load()
			if err != nil {
				return err
			}
			err = SavePersonalStance(sm.dir, s)
			if err != nil {
				return err
			}
			output := fmt.Sprintf("Saved current identity as personal stance: %s", s.Identity.Name)
			fmt.Println(FormatOutput(jsonOutput, output, nil))
			return nil
		}

		pools, err := LoadStancePoolsWithPersonal(sm.dir)
		if err != nil {
			return err
		}

		if inspireList {
			names := ListPoolNames(pools)
			output := fmt.Sprintf("%d pools:\n%s", len(names), strings.Join(names, "\n"))
			fmt.Println(FormatOutput(jsonOutput, output, nil))
			return nil
		}

		stance, pool, err := RandomStance(pools, inspirePoolName)
		if err != nil {
			return err
		}

		output := fmt.Sprintf("[%s]\nWho: %s\nWhere: %s\nLens: %s", pool, stance.Who, stance.Where, stance.Lens)
		fmt.Println(FormatOutput(jsonOutput, output, nil))
		return nil
	},
}

func init() {
	inspireCmd.Flags().StringVar(&inspirePoolName, "pool", "", "Draw from a specific pool")
	inspireCmd.Flags().BoolVar(&inspireList, "list", false, "List available stance pools")
	inspireCmd.Flags().BoolVar(&inspireSave, "save", false, "Save current identity as a personal stance")
	rootCmd.AddCommand(inspireCmd)
}
