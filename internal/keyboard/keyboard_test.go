package keyboard

import (
	"strings"
	"testing"
)

func TestResolveRune(t *testing.T) {
	tests := []struct {
		r          rune
		wantKey    string
		wantShift  bool
		wantShiftK string
		wantAltGr  bool
	}{
		{'f', "f", false, "", false},
		{' ', "space", false, "", false},
		{'A', "a", true, "rshift", false},
		{'J', "j", true, "lshift", false},
		{'?', "/", true, "lshift", false},
		{'ą', "a", false, "", true},
		{'ł', "l", false, "", true},
		{'Ż', "z", true, "rshift", true},
		{'Ś', "s", true, "rshift", true},
	}

	for _, tc := range tests {
		got, ok := ResolveRune(tc.r)
		if !ok {
			t.Fatalf("ResolveRune(%q) returned ok=false", tc.r)
		}
		if got.KeyID != tc.wantKey || got.NeedShift != tc.wantShift || got.ShiftKeyID != tc.wantShiftK || got.NeedAltGr != tc.wantAltGr {
			t.Errorf("ResolveRune(%q) = %+v, want key=%s shift=%v(%s) altgr=%v",
				tc.r, got, tc.wantKey, tc.wantShift, tc.wantShiftK, tc.wantAltGr)
		}
	}
}

func TestRenderContainsExpectedRows(t *testing.T) {
	out := Render(VisualState{
		ExpectedRune: 'ą',
		PressedRune:  'f',
		PressedOK:    true,
		ShowPressed:  true,
	})
	for _, token := range []string{"F*", "J*", "SPACE", "AltGr", "ENTER"} {
		if !strings.Contains(out, token) {
			t.Errorf("Render output missing %q", token)
		}
	}
}
