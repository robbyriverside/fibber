package fibber

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

// Global options and logger
var (
	logger  *zap.SugaredLogger
	once    sync.Once
	Options = &AppOptions{}
)

// AppOptions holds shared configuration for the app
type AppOptions struct {
	Verbose     bool
	AppName     string
	Version     string
	BuildTime   string
	Commit      string
	Environment string
}

// Core data types
type Gem struct {
	ID             string                 `json:"id"`
	Content        string                 `json:"content"`
	CreatedAt      time.Time              `json:"created_at"`
	StructuredData StructuredContribution `json:"structured_data"`
}

type StructuredContribution struct {
	GridContributions []GridContribution `json:"grid_contributions"`
}

type GridContribution struct {
	Beat  string `json:"beat"`
	Layer string `json:"layer"`
	Value string `json:"value"`
}

type Composition struct {
	ID        string                       `json:"id"`
	CreatedAt time.Time                    `json:"created_at"`
	Gems      []Gem                        `json:"gems"`
	Grid      map[string]map[string]string `json:"grid"`
}

var composition = Composition{
	ID:        "demo_composition",
	CreatedAt: time.Now(),
	Grid:      make(map[string]map[string]string),
}

func analyzeGem(content string) StructuredContribution {
	// Dummy AI behavior: hardcoded structure for example
	VLogf("Analyzing gem content: %q", content)
	return StructuredContribution{
		GridContributions: []GridContribution{
			{Beat: "midpoint", Layer: "emotional_arc", Value: "doubt and isolation"},
			{Beat: "midpoint", Layer: "theme", Value: "belonging vs rejection"},
		},
	}
}
