package lessons

import "fmt"

// Lesson represents a single structured touch-typing lesson.
type Lesson struct {
	ID          int
	Title       string
	Category    string
	FocusKeys   string
	Description string
	Prompt      string
}

var catalog = []Lesson{
	{
		ID:          1,
		Title:       "Home Row Anchors",
		Category:    "ASCII Basics",
		FocusKeys:   "f j d k Space",
		Description: "Find the tactile bumps on F and J with your index fingers.",
		Prompt:      "fff jjj fff jjj fjf jfj ddd kkk dkf jkd fjdk kdfj",
	},
	{
		ID:          2,
		Title:       "Full Home Row",
		Category:    "ASCII Basics",
		FocusKeys:   "a s d f g h j k l ;",
		Description: "Build muscle memory across all ten home-row positions.",
		Prompt:      "sad lad fall flask; glad dash half; ask a lad; all fall",
	},
	{
		ID:          3,
		Title:       "Index Finger Reach",
		Category:    "ASCII Basics",
		FocusKeys:   "r t y u v b n m",
		Description: "Stretch your index fingers up and down from F and J.",
		Prompt:      "run fun burn turn myth buy try; hurt my thumb by ruby",
	},
	{
		ID:          4,
		Title:       "Middle & Ring Fingers",
		Category:    "ASCII Basics",
		FocusKeys:   "e i c , w o x .",
		Description: "Coordinate middle and ring fingers across upper and lower rows.",
		Prompt:      "code, word, work, nice, ice, rice. exit the old box.",
	},
	{
		ID:          5,
		Title:       "Pinkies & Full Alphabet",
		Category:    "ASCII Basics",
		FocusKeys:   "q p z /",
		Description: "Train outer pinky columns and complete the lowercase alphabet.",
		Prompt:      "quiz zero park zip; quick lazy fox jumps over brown dog",
	},
	{
		ID:          6,
		Title:       "Shift & Punctuation",
		Category:    "Capitalization",
		FocusKeys:   "Shift ? ! ' -",
		Description: "Use opposite-hand Shift for capitals and everyday punctuation.",
		Prompt:      "Hello, World! It's a crisp-clear day. Are you ready?",
	},
	{
		ID:          7,
		Title:       "Numbers & Symbols",
		Category:    "Top Row",
		FocusKeys:   "1-0 @ # $ ( ) _ =",
		Description: "Reach the top number row and common developer symbols.",
		Prompt:      "v1.05: port=8080, #42 @user_99 ($100) [2026]",
	},
	{
		ID:          8,
		Title:       "English Full Sentences",
		Category:    "Fluency (ASCII)",
		FocusKeys:   "All ASCII keys",
		Description: "Type complete English sentences with steady rhythm and accuracy.",
		Prompt:      "The quick brown fox jumps over the lazy dog. Practice builds speed.",
	},
	{
		ID:          9,
		Title:       "Polish Diacritics (AltGr)",
		Category:    "Polish Layout",
		FocusKeys:   "AltGr + a c e l n o s x z",
		Description: "Hold Right Alt (AltGr) with your thumb to type Polish characters.",
		Prompt:      "żółć gęśl jaźń ćma łoś słońce księżyc źródło pączek",
	},
	{
		ID:          10,
		Title:       "Polish Full Sentences",
		Category:    "Polish Layout",
		FocusKeys:   "Full Polish QWERTY",
		Description: "Master Polish pangrams and natural sentences without looking down.",
		Prompt:      "Zażółć gęślą jaźń. Pchnąć w tę łódź jeża lub ośm skrzyń fig.",
	},
}

// All returns a copy of all available lessons in order.
func All() []Lesson {
	out := make([]Lesson, len(catalog))
	copy(out, catalog)
	return out
}

// Count returns the total number of lessons.
func Count() int {
	return len(catalog)
}

// Get returns the lesson matching id (1-based).
func Get(id int) (Lesson, error) {
	if id < 1 || id > len(catalog) {
		return Lesson{}, fmt.Errorf("lesson %d out of range (1-%d)", id, len(catalog))
	}
	return catalog[id-1], nil
}
