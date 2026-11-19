package lessons

import (
	"testing"
	"unicode/utf8"
)

func TestLessonsIntegrity(t *testing.T) {
	all := All()
	if len(all) != Count() || len(all) < 8 {
		t.Fatalf("unexpected lesson count: %d", len(all))
	}

	for i, l := range all {
		wantID := i + 1
		if l.ID != wantID {
			t.Errorf("lesson %d has ID %d, want %d", i, l.ID, wantID)
		}
		if l.Title == "" || l.Prompt == "" || l.Description == "" {
			t.Errorf("lesson %d has empty required fields", l.ID)
		}
		if !utf8.ValidString(l.Prompt) {
			t.Errorf("lesson %d prompt is not valid UTF-8", l.ID)
		}
		got, err := Get(l.ID)
		if err != nil {
			t.Errorf("Get(%d) returned error: %v", l.ID, err)
		}
		if got.Prompt != l.Prompt {
			t.Errorf("Get(%d) mismatch", l.ID)
		}
	}

	if _, err := Get(0); err == nil {
		t.Error("expected error for lesson 0")
	}
	if _, err := Get(Count() + 1); err == nil {
		t.Error("expected error for out-of-range lesson")
	}
}
