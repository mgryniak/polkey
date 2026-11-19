package tui

import (
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"polkey/internal/config"
	"polkey/internal/lessons"
)

func TestMenuStartsSelectedLesson(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "config.json")
	cfg := config.Config{LastLevel: 2, MaxErrors: 5, MinWPM: 10}
	m := NewModel(cfgPath, cfg, 0)

	// Type '3' and Enter
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	m = next.(Model)
	next, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = next.(Model)

	if m.state != StateTraining {
		t.Fatalf("expected StateTraining, got %v", m.state)
	}
	if m.lesson.ID != 3 {
		t.Fatalf("expected lesson 3, got %d", m.lesson.ID)
	}
}

func TestTrainingPassAndAutoAdvance(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "config.json")
	cfg := config.Config{LastLevel: 1, MaxErrors: 2, MinWPM: 5}
	m := NewModel(cfgPath, cfg, 1)

	// Simulate startTime 5 seconds ago so WPM is well above MinWPM=5.
	m.started = true
	m.startTime = time.Now().Add(-5 * time.Second)

	l, _ := lessons.Get(1)
	runes := []rune(l.Prompt)

	// Inject one mistake first.
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	m = next.(Model)
	if m.errors != 1 {
		t.Fatalf("expected 1 error, got %d", m.errors)
	}

	// Type the entire prompt accurately.
	for _, r := range runes {
		if r == ' ' {
			next, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
		} else {
			next, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		}
		m = next.(Model)
	}

	if m.state != StateResult {
		t.Fatalf("expected StateResult, got %v", m.state)
	}
	if !m.passed {
		t.Fatalf("expected lesson to pass (errors=%d, wpm=%.1f)", m.errors, m.finalWPM)
	}

	savedCfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("failed to load saved config: %v", err)
	}
	if savedCfg.LastLevel != 2 {
		t.Fatalf("expected LastLevel=2 in config.json, got %d", savedCfg.LastLevel)
	}

	// Simulate 3 auto-advance ticks.
	for i := 0; i < 3; i++ {
		next, _ = m.Update(autoAdvanceTickMsg{seq: m.resultSeq})
		m = next.(Model)
	}
	if m.state != StateTraining || m.lesson.ID != 2 {
		t.Fatalf("expected auto-advance to Lesson 2, got state=%v lesson=%d", m.state, m.lesson.ID)
	}
}

func TestTrainingFailTooManyErrors(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "config.json")
	cfg := config.Config{LastLevel: 1, MaxErrors: 0, MinWPM: 5}
	m := NewModel(cfgPath, cfg, 1)

	m.started = true
	m.startTime = time.Now().Add(-5 * time.Second)

	// One mistake exceeds MaxErrors=0.
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	m = next.(Model)

	l, _ := lessons.Get(1)
	for _, r := range []rune(l.Prompt) {
		if r == ' ' {
			next, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
		} else {
			next, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		}
		m = next.(Model)
	}

	if m.state != StateResult || m.passed {
		t.Fatalf("expected failed StateResult, got state=%v passed=%v", m.state, m.passed)
	}
}
