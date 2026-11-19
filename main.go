package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"polkey/internal/config"
	"polkey/internal/tui"
)

func main() {
	cfgPath := flag.String("config", config.DefaultPath, "path to config.json")
	lessonID := flag.Int("lesson", 0, "start directly at lesson number (1-10)")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading %s: %v\n", *cfgPath, err)
		os.Exit(1)
	}

	m := tui.NewModel(*cfgPath, cfg, *lessonID)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running polkey: %v\n", err)
		os.Exit(1)
	}
}
