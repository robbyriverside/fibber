package fibber

import (
	"encoding/json"
	"fmt"
	"time"
)

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
