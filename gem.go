package fibber

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

func InitLogger(env string) {
	once.Do(func() {
		env = strings.ToLower(env)
		if env == "" || env == "production" {
			env = "production"
		} else if env == "dev" {
			env = "development"
		}
		Options.Environment = env

		cfg := zap.NewProductionConfig()
		cfg.Encoding = "json"
		cfg.EncoderConfig.TimeKey = "time"
		cfg.EncoderConfig.LevelKey = "level"
		cfg.EncoderConfig.MessageKey = "msg"
		cfg.EncoderConfig.CallerKey = "caller"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder

		log, err := cfg.Build(zap.Fields(
			zap.String("app", Options.AppName),
			zap.String("version", Options.Version),
			zap.String("env", Options.Environment),
		))
		if err != nil {
			panic(err)
		}
		logger = log.Sugar()
	})
}

// VLogf logs only when verbose is true
func VLogf(format string, args ...interface{}) {
	if Options.Verbose {
		fmt.Printf("[verbose] "+format+"\n", args...)
	}
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

func InsertGem(id, content string) {
	structured := analyzeGem(content)
	gem := Gem{
		ID:             id,
		Content:        content,
		CreatedAt:      time.Now(),
		StructuredData: structured,
	}
	composition.Gems = append(composition.Gems, gem)
	for _, contrib := range structured.GridContributions {
		if _, ok := composition.Grid[contrib.Beat]; !ok {
			composition.Grid[contrib.Beat] = make(map[string]string)
		}
		composition.Grid[contrib.Beat][contrib.Layer] = contrib.Value
		VLogf("Updated grid at beat='%s', layer='%s' => %q", contrib.Beat, contrib.Layer, contrib.Value)
	}
	logger.Infof("Gem '%s' inserted and composition updated.", id)
}

func StatusReport() {
	fmt.Println("-- Composition Status Report --")
	last := len(composition.Gems)
	if last == 0 {
		fmt.Println("No gems inserted yet.")
		return
	}
	latest := composition.Gems[last-1]
	fmt.Printf("Last gem: %s\n", latest.ID)
	fmt.Printf("Inserted at: %s\n", latest.CreatedAt.Format(time.RFC822))
	fmt.Printf("Beat impact: ")
	for _, c := range latest.StructuredData.GridContributions {
		fmt.Printf("[%s:%s=%s] ", c.Beat, c.Layer, c.Value)
	}
	fmt.Println("\nSuggested next: Fill in 'all_is_lost' or explore character arc further.")
}

func Inspect() {
	data, err := json.MarshalIndent(composition, "", "  ")
	if err != nil {
		logger.Errorf("Failed to inspect composition: %v", err)
		return
	}
	fmt.Println(string(data))
}
