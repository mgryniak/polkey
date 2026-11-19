package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"polkey/internal/config"
	"polkey/internal/keyboard"
	"polkey/internal/lessons"
)

// ScreenState represents the active view of the application.
type ScreenState int

const (
	StateMenu ScreenState = iota
	StateTraining
	StateResult
	StateSettings
)

type clearFlashMsg struct{ seq int }
type autoAdvanceTickMsg struct{ seq int }

// Model is the Bubble Tea model for the typing tutor.
type Model struct {
	cfgPath string
	cfg     config.Config

	state      ScreenState
	menuInput  string
	menuErr    string
	width      int
	height     int

	// Training state
	lesson        lessons.Lesson
	targetRunes   []rune
	typedRunes    []rune
	errors        int
	started       bool
	startTime     time.Time
	endTime       time.Time
	lastPressed   rune
	lastPressedOK bool
	showPressed   bool
	flashSeq      int
	statusMsg     string

	// Result state
	finalWPM      float64
	finalAccuracy float64
	passed        bool
	autoCountdown int
	resultSeq     int

	// Settings state
	settingsField int // 0 = MaxErrors, 1 = MinWPM
	settingsBuf   string
	settingsErr   string
}

// NewModel initializes the TUI model. If initialLesson > 0, it starts directly
// in training mode for that lesson.
func NewModel(cfgPath string, cfg config.Config, initialLesson int) Model {
	m := Model{
		cfgPath: cfgPath,
		cfg:     cfg,
		state:   StateMenu,
	}
	if initialLesson > 0 {
		if l, err := lessons.Get(initialLesson); err == nil {
			m.startLesson(l)
		} else {
			m.menuErr = err.Error()
		}
	}
	return m
}

func (m *Model) startLesson(l lessons.Lesson) {
	m.state = StateTraining
	m.lesson = l
	m.targetRunes = []rune(l.Prompt)
	m.typedRunes = nil
	m.errors = 0
	m.started = false
	m.startTime = time.Time{}
	m.endTime = time.Time{}
	m.lastPressed = 0
	m.lastPressedOK = false
	m.showPressed = false
	m.statusMsg = ""
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case clearFlashMsg:
		if msg.seq == m.flashSeq {
			m.showPressed = false
		}
		return m, nil

	case autoAdvanceTickMsg:
		if m.state != StateResult || msg.seq != m.resultSeq || !m.passed || m.lesson.ID >= lessons.Count() {
			return m, nil
		}
		m.autoCountdown--
		if m.autoCountdown <= 0 {
			next, err := lessons.Get(m.lesson.ID + 1)
			if err == nil {
				m.startLesson(next)
			}
			return m, nil
		}
		return m, scheduleAutoAdvanceTick(m.resultSeq)

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		switch m.state {
		case StateMenu:
			return m.updateMenu(msg)
		case StateTraining:
			return m.updateTraining(msg)
		case StateResult:
			return m.updateResult(msg)
		case StateSettings:
			return m.updateSettings(msg)
		}
	}

	return m, nil
}

func (m Model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		raw := strings.TrimSpace(m.menuInput)
		m.menuInput = ""
		m.menuErr = ""

		if strings.EqualFold(raw, "q") || strings.EqualFold(raw, "quit") {
			return m, tea.Quit
		}
		if strings.EqualFold(raw, "s") {
			m.state = StateSettings
			m.settingsField = 0
			m.settingsBuf = strconv.Itoa(m.cfg.MaxErrors)
			m.settingsErr = ""
			return m, nil
		}

		level := m.cfg.LastLevel
		if raw != "" {
			n, err := strconv.Atoi(raw)
			if err != nil {
				m.menuErr = fmt.Sprintf("Invalid input %q: enter a lesson number (1-%d)", raw, lessons.Count())
				return m, nil
			}
			level = n
		}

		l, err := lessons.Get(level)
		if err != nil {
			m.menuErr = err.Error()
			return m, nil
		}
		m.startLesson(l)
		return m, nil

	case tea.KeyBackspace:
		if len(m.menuInput) > 0 {
			_, size := utf8.DecodeLastRuneInString(m.menuInput)
			m.menuInput = m.menuInput[:len(m.menuInput)-size]
		}
		return m, nil

	case tea.KeyRunes:
		m.menuInput += string(msg.Runes)
		m.menuErr = ""
		return m, nil
	}

	return m, nil
}

func (m Model) updateTraining(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.state = StateMenu
		m.menuInput = ""
		m.menuErr = ""
		return m, nil
	case tea.KeyCtrlR:
		m.startLesson(m.lesson)
		return m, nil
	}

	var pressed rune
	switch msg.Type {
	case tea.KeySpace:
		pressed = ' '
	case tea.KeyRunes:
		if len(msg.Runes) == 0 {
			return m, nil
		}
		pressed = msg.Runes[0]
	default:
		return m, nil
	}

	if !m.started {
		m.started = true
		m.startTime = time.Now()
	}

	idx := len(m.typedRunes)
	if idx >= len(m.targetRunes) {
		return m, nil
	}

	expected := m.targetRunes[idx]
	m.lastPressed = pressed
	m.showPressed = true
	m.flashSeq++
	flashCmd := scheduleClearFlash(m.flashSeq)

	if pressed == expected {
		m.typedRunes = append(m.typedRunes, pressed)
		m.lastPressedOK = true
		m.statusMsg = ""

		if len(m.typedRunes) == len(m.targetRunes) {
			m.endTime = time.Now()
			return m.finishLesson()
		}
	} else {
		m.errors++
		m.lastPressedOK = false
		m.statusMsg = fmt.Sprintf("Typed %q, expected %q", string(pressed), string(expected))
	}

	return m, flashCmd
}

func (m Model) finishLesson() (tea.Model, tea.Cmd) {
	m.finalWPM = m.computeWPM(m.endTime)
	m.finalAccuracy = m.computeAccuracy()
	m.passed = m.errors <= m.cfg.MaxErrors && m.finalWPM >= float64(m.cfg.MinWPM)
	m.state = StateResult
	m.resultSeq++

	if m.passed && m.lesson.ID < lessons.Count() {
		nextLevel := m.lesson.ID + 1
		if nextLevel > m.cfg.LastLevel {
			m.cfg.LastLevel = nextLevel
		} else {
			m.cfg.LastLevel = nextLevel
		}
		_ = config.Save(m.cfgPath, m.cfg)
		m.autoCountdown = 3
		return m, scheduleAutoAdvanceTick(m.resultSeq)
	}

	return m, nil
}

func (m Model) updateResult(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		if m.passed && m.lesson.ID < lessons.Count() {
			if next, err := lessons.Get(m.lesson.ID + 1); err == nil {
				m.startLesson(next)
				return m, nil
			}
		}
		m.startLesson(m.lesson)
		return m, nil

	case tea.KeyEsc:
		m.state = StateMenu
		return m, nil

	case tea.KeyRunes:
		if len(msg.Runes) == 0 {
			return m, nil
		}
		switch msg.Runes[0] {
		case 'r', 'R':
			m.startLesson(m.lesson)
			return m, nil
		case 'm', 'M':
			m.state = StateMenu
			return m, nil
		case 'q', 'Q':
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) updateSettings(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.state = StateMenu
		return m, nil

	case tea.KeyTab:
		if m.settingsField == 0 {
			if v, err := strconv.Atoi(strings.TrimSpace(m.settingsBuf)); err == nil && v >= 0 {
				m.cfg.MaxErrors = v
			}
			m.settingsField = 1
			m.settingsBuf = strconv.Itoa(m.cfg.MinWPM)
		} else {
			if v, err := strconv.Atoi(strings.TrimSpace(m.settingsBuf)); err == nil && v >= 1 {
				m.cfg.MinWPM = v
			}
			m.settingsField = 0
			m.settingsBuf = strconv.Itoa(m.cfg.MaxErrors)
		}
		m.settingsErr = ""
		return m, nil

	case tea.KeyEnter:
		v, err := strconv.Atoi(strings.TrimSpace(m.settingsBuf))
		if err != nil {
			m.settingsErr = "Please enter a valid whole number."
			return m, nil
		}
		if m.settingsField == 0 {
			if v < 0 {
				m.settingsErr = "Max errors cannot be negative."
				return m, nil
			}
			m.cfg.MaxErrors = v
			m.settingsField = 1
			m.settingsBuf = strconv.Itoa(m.cfg.MinWPM)
			m.settingsErr = ""
			return m, nil
		}

		if v < 1 {
			m.settingsErr = "Min WPM must be at least 1."
			return m, nil
		}
		m.cfg.MinWPM = v
		if err := config.Save(m.cfgPath, m.cfg); err != nil {
			m.settingsErr = fmt.Sprintf("Failed to save config: %v", err)
			return m, nil
		}
		m.state = StateMenu
		return m, nil

	case tea.KeyBackspace:
		if len(m.settingsBuf) > 0 {
			m.settingsBuf = m.settingsBuf[:len(m.settingsBuf)-1]
		}
		return m, nil

	case tea.KeyRunes:
		for _, r := range msg.Runes {
			if r >= '0' && r <= '9' {
				m.settingsBuf += string(r)
			}
		}
		m.settingsErr = ""
		return m, nil
	}

	return m, nil
}

func (m Model) computeWPM(now time.Time) float64 {
	if !m.started || len(m.typedRunes) == 0 {
		return 0
	}
	mins := now.Sub(m.startTime).Minutes()
	if mins < 0.005 {
		mins = 0.005
	}
	return (float64(len(m.typedRunes)) / 5.0) / mins
}

func (m Model) computeAccuracy() float64 {
	total := len(m.typedRunes) + m.errors
	if total == 0 {
		return 100
	}
	return (float64(len(m.typedRunes)) / float64(total)) * 100
}

func scheduleClearFlash(seq int) tea.Cmd {
	return tea.Tick(350*time.Millisecond, func(time.Time) tea.Msg {
		return clearFlashMsg{seq: seq}
	})
}

func scheduleAutoAdvanceTick(seq int) tea.Cmd {
	return tea.Tick(1*time.Second, func(time.Time) tea.Msg {
		return autoAdvanceTickMsg{seq: seq}
	})
}

// View implements tea.Model.
func (m Model) View() string {
	switch m.state {
	case StateMenu:
		return m.viewMenu()
	case StateTraining:
		return m.viewTraining()
	case StateResult:
		return m.viewResult()
	case StateSettings:
		return m.viewSettings()
	default:
		return ""
	}
}

var (
	styleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("45"))

	styleSubtitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	styleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1)

	stylePromptTyped = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42"))

	stylePromptCursor = lipgloss.NewStyle().
				Foreground(lipgloss.Color("16")).
				Background(lipgloss.Color("45")).
				Bold(true)

	stylePromptRest = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	styleCmdBar = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	stylePass = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	styleFail = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)
)

func (m Model) viewMenu() string {
	var b strings.Builder
	b.WriteString(styleTitle.Render("POLKEY — Console Touch Typing Tutor"))
	b.WriteString("\n")
	b.WriteString(styleSubtitle.Render("Learn QWERTY & Polish diacritics without looking at the keyboard."))
	b.WriteString("\n\n")

	for _, l := range lessons.All() {
		marker := "  "
		if l.ID == m.cfg.LastLevel {
			marker = styleTitle.Render("->")
		}
		line := fmt.Sprintf("%s %2d. %-26s [%s] (%s)",
			marker, l.ID, l.Title, l.Category, l.FocusKeys)
		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	cfgLine := fmt.Sprintf("Config (%s): Last Level = %d | Max Errors = %d | Min WPM = %d",
		m.cfgPath, m.cfg.LastLevel, m.cfg.MaxErrors, m.cfg.MinWPM)
	b.WriteString(styleSubtitle.Render(cfgLine))
	b.WriteString("\n")

	if m.menuErr != "" {
		b.WriteString(styleFail.Render("Error: " + m.menuErr))
		b.WriteString("\n")
	}

	prompt := fmt.Sprintf("Enter lesson number (1-%d) [default: %d], 's' settings, 'q' quit: > %s_",
		lessons.Count(), m.cfg.LastLevel, m.menuInput)
	b.WriteString("\n")
	b.WriteString(styleCmdBar.Render(prompt))
	return b.String()
}

func (m Model) viewTraining() string {
	var b strings.Builder

	wpm := 0.0
	if m.started {
		wpm = m.computeWPM(time.Now())
	}
	acc := m.computeAccuracy()

	header := fmt.Sprintf("Lesson %d/%d: %s   |   WPM: %.0f (min %d)   |   Acc: %.0f%%   |   Errors: %d/%d   |   %d/%d",
		m.lesson.ID, lessons.Count(), m.lesson.Title,
		wpm, m.cfg.MinWPM, acc, m.errors, m.cfg.MaxErrors,
		len(m.typedRunes), len(m.targetRunes))
	b.WriteString(styleTitle.Render(header))
	b.WriteString("\n")
	b.WriteString(styleSubtitle.Render(m.lesson.Description + "  (Focus: " + m.lesson.FocusKeys + ")"))
	b.WriteString("\n\n")

	// Render target prompt with typed/cursor/remaining styling.
	idx := len(m.typedRunes)
	var promptView strings.Builder
	if idx > 0 {
		promptView.WriteString(stylePromptTyped.Render(string(m.targetRunes[:idx])))
	}
	var expectedRune rune
	if idx < len(m.targetRunes) {
		expectedRune = m.targetRunes[idx]
		cur := string(expectedRune)
		if expectedRune == ' ' {
			cur = "␣"
		}
		promptView.WriteString(stylePromptCursor.Render(cur))
		if idx+1 < len(m.targetRunes) {
			promptView.WriteString(stylePromptRest.Render(string(m.targetRunes[idx+1:])))
		}
	}
	b.WriteString(styleBox.Render(promptView.String()))
	b.WriteString("\n\n")

	// Visual keyboard.
	b.WriteString(keyboard.Render(keyboard.VisualState{
		ExpectedRune: expectedRune,
		PressedRune:  m.lastPressed,
		PressedOK:    m.lastPressedOK,
		ShowPressed:  m.showPressed,
	}))
	b.WriteString("\n")

	// Finger hint + status line.
	hint := "Done!"
	if combo, ok := keyboard.ResolveRune(expectedRune); ok {
		charLabel := fmt.Sprintf("%q", string(expectedRune))
		if expectedRune == ' ' {
			charLabel = "SPACE"
		}
		hint = fmt.Sprintf("Next: %s -> %s", charLabel, combo.FingerHint)
	}
	if m.statusMsg != "" {
		b.WriteString(styleFail.Render(m.statusMsg) + "   |   " + styleSubtitle.Render(hint))
	} else {
		b.WriteString(styleSubtitle.Render(hint + "   |   [Ctrl+R] Restart   [ESC] Menu"))
	}
	b.WriteString("\n")

	// Bottom command line showing typed text.
	cmdText := fmt.Sprintf("> %s_", string(m.typedRunes))
	b.WriteString(styleCmdBar.Render(cmdText))

	return b.String()
}

func (m Model) viewResult() string {
	var b strings.Builder
	b.WriteString(styleTitle.Render(fmt.Sprintf("Lesson %d Complete: %s", m.lesson.ID, m.lesson.Title)))
	b.WriteString("\n\n")

	stats := fmt.Sprintf("Speed:    %.1f WPM (target >= %d WPM)\nAccuracy: %.1f%%\nErrors:   %d (max allowed: %d)",
		m.finalWPM, m.cfg.MinWPM, m.finalAccuracy, m.errors, m.cfg.MaxErrors)
	b.WriteString(styleBox.Render(stats))
	b.WriteString("\n\n")

	if m.passed {
		b.WriteString(stylePass.Render("PASSED! Great job."))
		b.WriteString("\n")
		if m.lesson.ID < lessons.Count() {
			b.WriteString(fmt.Sprintf("Advancing automatically to Lesson %d in %ds...\n", m.lesson.ID+1, m.autoCountdown))
			b.WriteString(styleSubtitle.Render("[Enter] Next Lesson Now   [r] Retry   [m] Menu   [q] Quit"))
		} else {
			b.WriteString(stylePass.Render("Congratulations! You completed all lessons!"))
			b.WriteString("\n")
			b.WriteString(styleSubtitle.Render("[r] Retry   [m] Menu   [q] Quit"))
		}
	} else {
		b.WriteString(styleFail.Render("NOT PASSED YET — keep practicing!"))
		b.WriteString("\n")
		if m.errors > m.cfg.MaxErrors {
			b.WriteString(fmt.Sprintf("- Too many errors (%d > %d)\n", m.errors, m.cfg.MaxErrors))
		}
		if m.finalWPM < float64(m.cfg.MinWPM) {
			b.WriteString(fmt.Sprintf("- Speed too low (%.1f < %d WPM)\n", m.finalWPM, m.cfg.MinWPM))
		}
		b.WriteString("\n")
		b.WriteString(styleSubtitle.Render("[Enter / r] Retry Lesson   [m] Menu   [q] Quit"))
	}

	return b.String()
}

func (m Model) viewSettings() string {
	var b strings.Builder
	b.WriteString(styleTitle.Render("Settings (saved to " + m.cfgPath + ")"))
	b.WriteString("\n\n")

	f0 := strconv.Itoa(m.cfg.MaxErrors)
	f1 := strconv.Itoa(m.cfg.MinWPM)
	if m.settingsField == 0 {
		f0 = m.settingsBuf + "_"
	} else {
		f1 = m.settingsBuf + "_"
	}

	b.WriteString(fmt.Sprintf("1. Max Allowed Errors: %s\n", f0))
	b.WriteString(fmt.Sprintf("2. Min Passing WPM:    %s\n", f1))
	b.WriteString("\n")
	if m.settingsErr != "" {
		b.WriteString(styleFail.Render(m.settingsErr) + "\n")
	}
	b.WriteString(styleSubtitle.Render("[Tab] Switch Field   [Enter] Save / Next   [ESC] Cancel"))
	return b.String()
}
