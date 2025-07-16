# TUI Typing Game

A simple terminal-based typing game implemented in Go using the x/term package.

## Features

- Random English text typing practice
- Typing speed (WPM) measurement
- Accuracy measurement
- Error count tracking
- Optimized layout with center-screen positioning
- Safe termination with signal handling
- Stabilized display with buffer control

## Requirements

- Go 1.24 or higher
- ANSI-compatible terminal emulator

## Installation

```bash
git clone https://github.com/kenta-takeuchi/tui-typing-game.git
cd tui-typing-game
go mod tidy
```

## Running the Game

```bash
go run main.go
```

## How to Play

1. When the game starts, text to type will appear in the center of the screen
2. Type the displayed text accurately
3. After completing each text, results (time, accuracy, speed) will be displayed
4. The next text will appear automatically
5. Press `Ctrl+C` to exit (safe termination will be executed)

## Controls

- Regular key input: Text input
- Backspace: Delete one character
- Ctrl+C: Exit game (executes safe termination)

## Implementation Features

- Screen control using ANSI escape sequences
- Dynamic layout adjustment based on terminal size
- Signal handling (SIGINT, SIGTERM) for safe termination
- Display stabilization through buffer control
- Appropriate cursor show/hide control