package keyboard

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
)

// Hand indicates which hand operates a key (used for opposite-hand Shift hints).
type Hand int

const (
	HandLeft Hand = iota
	HandRight
	HandThumb
)

// KeyDef describes a single physical key on a QWERTY/Polish-Programmers keyboard.
type KeyDef struct {
	ID        string
	Label     string
	Width     int
	Base      rune
	Shifted   rune
	AltGr     rune
	Hand      Hand
	Finger    string
	HomeBump  bool
	IsSpace   bool
	IsSpecial bool
}

// KeyCombo describes the physical key(s) needed to produce a rune.
type KeyCombo struct {
	KeyID      string
	NeedShift  bool
	ShiftKeyID string
	NeedAltGr  bool
	FingerHint string
}

// VisualState controls which keys are highlighted on the rendered keyboard.
type VisualState struct {
	ExpectedRune rune
	PressedRune  rune
	PressedOK    bool
	ShowPressed  bool
}

var rows = [][]KeyDef{
	{
		{ID: "`", Label: "`", Width: 4, Base: '`', Shifted: '~', Hand: HandLeft, Finger: "Left Pinky"},
		{ID: "1", Label: "1", Width: 4, Base: '1', Shifted: '!', Hand: HandLeft, Finger: "Left Pinky"},
		{ID: "2", Label: "2", Width: 4, Base: '2', Shifted: '@', Hand: HandLeft, Finger: "Left Ring"},
		{ID: "3", Label: "3", Width: 4, Base: '3', Shifted: '#', Hand: HandLeft, Finger: "Left Middle"},
		{ID: "4", Label: "4", Width: 4, Base: '4', Shifted: '$', Hand: HandLeft, Finger: "Left Index"},
		{ID: "5", Label: "5", Width: 4, Base: '5', Shifted: '%', Hand: HandLeft, Finger: "Left Index"},
		{ID: "6", Label: "6", Width: 4, Base: '6', Shifted: '^', Hand: HandRight, Finger: "Right Index"},
		{ID: "7", Label: "7", Width: 4, Base: '7', Shifted: '&', Hand: HandRight, Finger: "Right Index"},
		{ID: "8", Label: "8", Width: 4, Base: '8', Shifted: '*', Hand: HandRight, Finger: "Right Middle"},
		{ID: "9", Label: "9", Width: 4, Base: '9', Shifted: '(', Hand: HandRight, Finger: "Right Ring"},
		{ID: "0", Label: "0", Width: 4, Base: '0', Shifted: ')', Hand: HandRight, Finger: "Right Pinky"},
		{ID: "-", Label: "-", Width: 4, Base: '-', Shifted: '_', Hand: HandRight, Finger: "Right Pinky"},
		{ID: "=", Label: "=", Width: 4, Base: '=', Shifted: '+', Hand: HandRight, Finger: "Right Pinky"},
		{ID: "bksp", Label: "BKSP", Width: 8, Hand: HandRight, Finger: "Right Pinky", IsSpecial: true},
	},
	{
		{ID: "tab", Label: "TAB", Width: 6, Hand: HandLeft, Finger: "Left Pinky", IsSpecial: true},
		{ID: "q", Label: "Q", Width: 4, Base: 'q', Shifted: 'Q', Hand: HandLeft, Finger: "Left Pinky"},
		{ID: "w", Label: "W", Width: 4, Base: 'w', Shifted: 'W', Hand: HandLeft, Finger: "Left Ring"},
		{ID: "e", Label: "E", Width: 4, Base: 'e', Shifted: 'E', AltGr: 'ę', Hand: HandLeft, Finger: "Left Middle"},
		{ID: "r", Label: "R", Width: 4, Base: 'r', Shifted: 'R', Hand: HandLeft, Finger: "Left Index"},
		{ID: "t", Label: "T", Width: 4, Base: 't', Shifted: 'T', Hand: HandLeft, Finger: "Left Index"},
		{ID: "y", Label: "Y", Width: 4, Base: 'y', Shifted: 'Y', Hand: HandRight, Finger: "Right Index"},
		{ID: "u", Label: "U", Width: 4, Base: 'u', Shifted: 'U', Hand: HandRight, Finger: "Right Index"},
		{ID: "i", Label: "I", Width: 4, Base: 'i', Shifted: 'I', Hand: HandRight, Finger: "Right Middle"},
		{ID: "o", Label: "O", Width: 4, Base: 'o', Shifted: 'O', AltGr: 'ó', Hand: HandRight, Finger: "Right Ring"},
		{ID: "p", Label: "P", Width: 4, Base: 'p', Shifted: 'P', Hand: HandRight, Finger: "Right Pinky"},
		{ID: "[", Label: "[", Width: 4, Base: '[', Shifted: '{', Hand: HandRight, Finger: "Right Pinky"},
		{ID: "]", Label: "]", Width: 4, Base: ']', Shifted: '}', Hand: HandRight, Finger: "Right Pinky"},
		{ID: "\\", Label: "\\", Width: 6, Base: '\\', Shifted: '|', Hand: HandRight, Finger: "Right Pinky"},
	},
	{
		{ID: "caps", Label: "CAPS", Width: 7, Hand: HandLeft, Finger: "Left Pinky", IsSpecial: true},
		{ID: "a", Label: "A", Width: 4, Base: 'a', Shifted: 'A', AltGr: 'ą', Hand: HandLeft, Finger: "Left Pinky"},
		{ID: "s", Label: "S", Width: 4, Base: 's', Shifted: 'S', AltGr: 'ś', Hand: HandLeft, Finger: "Left Ring"},
		{ID: "d", Label: "D", Width: 4, Base: 'd', Shifted: 'D', Hand: HandLeft, Finger: "Left Middle"},
		{ID: "f", Label: "F*", Width: 4, Base: 'f', Shifted: 'F', Hand: HandLeft, Finger: "Left Index", HomeBump: true},
		{ID: "g", Label: "G", Width: 4, Base: 'g', Shifted: 'G', Hand: HandLeft, Finger: "Left Index"},
		{ID: "h", Label: "H", Width: 4, Base: 'h', Shifted: 'H', Hand: HandRight, Finger: "Right Index"},
		{ID: "j", Label: "J*", Width: 4, Base: 'j', Shifted: 'J', Hand: HandRight, Finger: "Right Index", HomeBump: true},
		{ID: "k", Label: "K", Width: 4, Base: 'k', Shifted: 'K', Hand: HandRight, Finger: "Right Middle"},
		{ID: "l", Label: "L", Width: 4, Base: 'l', Shifted: 'L', AltGr: 'ł', Hand: HandRight, Finger: "Right Ring"},
		{ID: ";", Label: ";", Width: 4, Base: ';', Shifted: ':', Hand: HandRight, Finger: "Right Pinky"},
		{ID: "'", Label: "'", Width: 4, Base: '\'', Shifted: '"', Hand: HandRight, Finger: "Right Pinky"},
		{ID: "enter", Label: "ENTER", Width: 9, Hand: HandRight, Finger: "Right Pinky", IsSpecial: true},
	},
	{
		{ID: "lshift", Label: "SHIFT", Width: 9, Hand: HandLeft, Finger: "Left Pinky", IsSpecial: true},
		{ID: "z", Label: "Z", Width: 4, Base: 'z', Shifted: 'Z', AltGr: 'ż', Hand: HandLeft, Finger: "Left Pinky"},
		{ID: "x", Label: "X", Width: 4, Base: 'x', Shifted: 'X', AltGr: 'ź', Hand: HandLeft, Finger: "Left Ring"},
		{ID: "c", Label: "C", Width: 4, Base: 'c', Shifted: 'C', AltGr: 'ć', Hand: HandLeft, Finger: "Left Middle"},
		{ID: "v", Label: "V", Width: 4, Base: 'v', Shifted: 'V', Hand: HandLeft, Finger: "Left Index"},
		{ID: "b", Label: "B", Width: 4, Base: 'b', Shifted: 'B', Hand: HandLeft, Finger: "Left Index"},
		{ID: "n", Label: "N", Width: 4, Base: 'n', Shifted: 'N', AltGr: 'ń', Hand: HandRight, Finger: "Right Index"},
		{ID: "m", Label: "M", Width: 4, Base: 'm', Shifted: 'M', Hand: HandRight, Finger: "Right Index"},
		{ID: ",", Label: ",", Width: 4, Base: ',', Shifted: '<', Hand: HandRight, Finger: "Right Middle"},
		{ID: ".", Label: ".", Width: 4, Base: '.', Shifted: '>', Hand: HandRight, Finger: "Right Ring"},
		{ID: "/", Label: "/", Width: 4, Base: '/', Shifted: '?', Hand: HandRight, Finger: "Right Pinky"},
		{ID: "rshift", Label: "SHIFT", Width: 11, Hand: HandRight, Finger: "Right Pinky", IsSpecial: true},
	},
	{
		{ID: "lctrl", Label: "CTRL", Width: 7, Hand: HandLeft, Finger: "Left Pinky", IsSpecial: true},
		{ID: "lalt", Label: "ALT", Width: 6, Hand: HandLeft, Finger: "Thumb", IsSpecial: true},
		{ID: "space", Label: "SPACE", Width: 32, Base: ' ', Shifted: ' ', Hand: HandThumb, Finger: "Thumb", IsSpace: true},
		{ID: "altgr", Label: "AltGr", Width: 8, Hand: HandRight, Finger: "Right Thumb", IsSpecial: true},
		{ID: "rctrl", Label: "CTRL", Width: 7, Hand: HandRight, Finger: "Right Pinky", IsSpecial: true},
	},
}

// ResolveRune maps any character to the physical key combination required to type it.
func ResolveRune(r rune) (KeyCombo, bool) {
	if r == 0 {
		return KeyCombo{}, false
	}

	for _, row := range rows {
		for _, k := range row {
			if k.Base != 0 && k.Base == r {
				return KeyCombo{
					KeyID:      k.ID,
					FingerHint: k.Finger,
				}, true
			}
			if k.Shifted != 0 && k.Shifted == r {
				shiftID := "rshift"
				shiftLabel := "Right Shift"
				if k.Hand == HandRight {
					shiftID = "lshift"
					shiftLabel = "Left Shift"
				}
				return KeyCombo{
					KeyID:      k.ID,
					NeedShift:  true,
					ShiftKeyID: shiftID,
					FingerHint: fmt.Sprintf("%s + %s", shiftLabel, k.Finger),
				}, true
			}
			if k.AltGr != 0 && (k.AltGr == r || unicode.ToUpper(k.AltGr) == r) {
				isUpper := unicode.IsUpper(r)
				shiftID := ""
				hint := fmt.Sprintf("AltGr (Right Alt) + %s", k.Finger)
				if isUpper {
					shiftID = "rshift"
					if k.Hand == HandRight {
						shiftID = "lshift"
					}
					hint = fmt.Sprintf("Shift + AltGr + %s", k.Finger)
				}
				return KeyCombo{
					KeyID:      k.ID,
					NeedShift:  isUpper,
					ShiftKeyID: shiftID,
					NeedAltGr:  true,
					FingerHint: hint,
				}, true
			}
		}
	}

	return KeyCombo{}, false
}

var (
	styleNormal = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("236"))

	styleSpecial = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Background(lipgloss.Color("235"))

	styleHomeBump = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.Color("238")).
			Bold(true)

	styleExpected = lipgloss.NewStyle().
			Foreground(lipgloss.Color("16")).
			Background(lipgloss.Color("45")).
			Bold(true)

	styleModifier = lipgloss.NewStyle().
			Foreground(lipgloss.Color("16")).
			Background(lipgloss.Color("220")).
			Bold(true)

	stylePressedOK = lipgloss.NewStyle().
			Foreground(lipgloss.Color("16")).
			Background(lipgloss.Color("42")).
			Bold(true)

	stylePressedErr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.Color("196")).
			Bold(true)

	styleFrame = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)
)

// Render draws the full visual keyboard with expected and pressed key highlights.
func Render(state VisualState) string {
	expectedCombo, hasExpected := ResolveRune(state.ExpectedRune)
	var pressedCombo KeyCombo
	var hasPressed bool
	if state.ShowPressed && state.PressedRune != 0 {
		pressedCombo, hasPressed = ResolveRune(state.PressedRune)
	}

	var renderedRows []string
	for _, row := range rows {
		var cells []string
		for _, k := range row {
			st := styleNormal
			if k.IsSpecial {
				st = styleSpecial
			} else if k.HomeBump {
				st = styleHomeBump
			}

			// Highlight expected key and its modifiers.
			if hasExpected {
				if k.ID == expectedCombo.KeyID {
					st = styleExpected
				} else if (expectedCombo.NeedShift && k.ID == expectedCombo.ShiftKeyID) ||
					(expectedCombo.NeedAltGr && k.ID == "altgr") {
					st = styleModifier
				}
			}

			// Highlight last pressed key (green if correct, red if mistake).
			if hasPressed && k.ID == pressedCombo.KeyID {
				if state.PressedOK {
					st = stylePressedOK
				} else {
					st = stylePressedErr
				}
			}

			cells = append(cells, st.Width(k.Width).Align(lipgloss.Center).Render(k.Label))
		}
		renderedRows = append(renderedRows, strings.Join(cells, " "))
	}

	legend := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Render("Legend: [Cyan] Next Key  [Yellow] Shift/AltGr  [Green] Correct  [Red] Error  [*]=Home Bump")

	body := strings.Join(renderedRows, "\n") + "\n" + legend
	return styleFrame.Render(body)
}
