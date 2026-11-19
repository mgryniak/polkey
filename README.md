# polkey

A terminal-based touch-typing tutor written in Go using [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

Designed for learning blind typing from scratch—starting on ASCII home-row anchors and progressing through full English sentences, Polish diacritics (`AltGr`), and Polish pangrams.

## Features

- **Live Visual Keyboard**: Displays a full QWERTY keyboard in the console with:
  - **Cyan** highlight on the next expected key
  - **Yellow** highlight on required modifiers (`SHIFT` or `AltGr` for Polish characters)
  - **Green / Red** flash on the key you just pressed
  - **Tactile Home Row Markers** (`F*` and `J*`) and per-key finger hints
- **Bottom Command Input Line**: Shows your typed text in real time at the bottom of the terminal alongside live WPM, accuracy, and error counts.
- **10 Progressive Lessons**: Structured from home-row muscle memory to full sentences in both English and Polish.
- **Automatic Progression**: Passing a lesson within the allowed error limit and minimum WPM automatically saves your progress and advances to the next lesson.
- **Persistent Configuration (`config.json`)**: Stores `last_level`, `max_errors`, and `min_wpm`, editable directly in the file or via the in-app Settings screen (`s`).

## Installation & Usage

### Prerequisites

- Go `1.24+` (or let Go's toolchain automatically download the version specified in `go.mod`)

### Run from source

```bash
go run .
```

### Build a binary

```bash
go build -o polkey .
./polkey
```

On Windows:

```powershell
go build -o polkey.exe .
.\polkey.exe
```

### CLI Flags

- `--lesson <N>`: Jump directly into lesson `N` (`1`–`10`) on startup.
- `--config <path>`: Use a custom path for `config.json` (defaults to `config.json` in the current directory).

### Controls

- **Menu**: Enter a lesson number (`1`–`10`) and press `Enter` (or press `Enter` on an empty prompt to resume `last_level`), `s` for Settings, `q` to quit.
- **Training**: Type the displayed text. `Ctrl+R` restarts the current lesson; `ESC` returns to the menu.
- **Result**: Automatically advances to the next lesson after 3 seconds on pass (or press `Enter` immediately, `r` to retry, `m` for menu).

## Lesson Curriculum

| # | Title | Category | Focus Keys |
|---|-------|----------|------------|
| 1 | Home Row Anchors | ASCII Basics | `f j d k Space` |
| 2 | Full Home Row | ASCII Basics | `a s d f g h j k l ;` |
| 3 | Index Finger Reach | ASCII Basics | `r t y u v b n m` |
| 4 | Middle & Ring Fingers | ASCII Basics | `e i c , w o x .` |
| 5 | Pinkies & Full Alphabet | ASCII Basics | `q p z /` |
| 6 | Shift & Punctuation | Capitalization | `Shift ? ! ' -` |
| 7 | Numbers & Symbols | Top Row | `1-0 @ # $ ( ) _ =` |
| 8 | English Full Sentences | Fluency (ASCII) | All ASCII keys |
| 9 | Polish Diacritics (AltGr) | Polish Layout | `AltGr + a c e l n o s x z` |
| 10 | Polish Full Sentences | Polish Layout | Full Polish QWERTY |

## Configuration (`config.json`)

If `config.json` does not exist, `polkey` creates it automatically with default values:

```json
{
  "last_level": 1,
  "max_errors": 5,
  "min_wpm": 20
}
```

- `last_level`: The default lesson selected on startup; updated automatically when you pass a lesson.
- `max_errors`: Maximum number of mistakes allowed to pass a lesson.
- `min_wpm`: Minimum typing speed (Words Per Minute) required to pass a lesson.

## Development & Testing

```bash
go vet ./...
go test -v ./...
```

## License

This project is licensed under the [MIT License](LICENSE).
