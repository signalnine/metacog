package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"

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

var inspirePoolName string
var inspireList bool

var inspireCmd = &cobra.Command{
	Use:   "inspire",
	Short: "Draw a random stance from the pool",
	RunE: func(cmd *cobra.Command, args []string) error {
		pools, err := LoadStancePools()
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
	rootCmd.AddCommand(inspireCmd)
}
