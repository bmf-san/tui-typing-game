package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/term"
)

var words = []string{
	"Hello World",
	"Go Programming",
	"Type this text",
	"Practice makes perfect",
	"Keep coding",
}

type Game struct {
	oldState  *term.State
	word      string
	input     string
	startTime time.Time
	mistakes  int
	width     int
	height    int
	done      chan struct{}
}

func newGame() (*Game, error) {
	// Ensure stdout is synced
	os.Stdout.Sync()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("failed to set terminal to raw mode: %v", err)
	}

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil, fmt.Errorf("failed to get terminal size: %v", err)
	}

	return &Game{
		oldState:  oldState,
		word:      words[rand.Intn(len(words))],
		startTime: time.Now(),
		width:     width,
		height:    height,
		done:      make(chan struct{}),
	}, nil
}

func (g *Game) cleanup() {
	// Clear screen and move cursor to top
	fmt.Print("\033[H\033[2J")
	fmt.Print("\033[H")
	fmt.Print("\033[?25h") // Show cursor
	term.Restore(int(os.Stdin.Fd()), g.oldState)
	fmt.Print("\n") // Add a newline after restoring
	os.Stdout.Sync()
}

func (g *Game) centerText(text string, row int) {
	padding := (g.width - len(text)) / 2
	if padding < 0 {
		padding = 0
	}
	fmt.Printf("\033[%d;%dH%s\033[K", row, padding+1, text)
	os.Stdout.Sync()
}

func (g *Game) clearLine(row int) {
	fmt.Printf("\033[%d;1H\033[K", row)
	os.Stdout.Sync()
}

func (g *Game) render() {
	// Hide cursor and clear screen
	fmt.Print("\033[?25l")
	fmt.Print("\033[H\033[2J")
	os.Stdout.Sync()

	// Clear all lines we're going to use
	for i := 1; i <= 7; i++ {
		g.clearLine(i)
	}

	// Title
	g.centerText("=== Typing Game ===", 1)

	// Target text
	targetText := fmt.Sprintf("Type this: %s", g.word)
	g.centerText(targetText, 3)

	// User input
	inputText := fmt.Sprintf("Your input: %s", g.input)
	g.centerText(inputText, 5)

	// Mistakes
	mistakesText := fmt.Sprintf("Mistakes: %d", g.mistakes)
	g.centerText(mistakesText, 7)

	// Move cursor to input position
	fmt.Printf("\033[5;%dH", (g.width-len("Your input: ")+len(g.input))/2+1)
	os.Stdout.Sync()
}

func (g *Game) showResults(duration time.Duration, accuracy, wpm float64) {
	// Clear screen
	fmt.Print("\033[H\033[2J")
	os.Stdout.Sync()

	// Clear all lines we're going to use
	for i := 1; i <= 7; i++ {
		g.clearLine(i)
	}

	g.centerText("=== Results ===", 1)

	timeText := fmt.Sprintf("Time: %.2f seconds", duration.Seconds())
	g.centerText(timeText, 3)

	accuracyText := fmt.Sprintf("Accuracy: %.1f%%", accuracy)
	g.centerText(accuracyText, 4)

	speedText := fmt.Sprintf("Speed: %.1f WPM", wpm)
	g.centerText(speedText, 5)

	g.centerText("Press any key to continue, Ctrl+C to exit...", 7)
	os.Stdout.Sync()
}

func (g *Game) handleInput(b byte) bool {
	if b == 3 { // Ctrl+C
		close(g.done)
		return false
	}

	if b == 127 { // Backspace
		if len(g.input) > 0 {
			g.input = g.input[:len(g.input)-1]
		}
		return true
	}

	if len(g.input) < len(g.word) {
		g.input += string(b)

		// Check for mistakes
		if g.input[len(g.input)-1] != g.word[len(g.input)-1] {
			g.mistakes++
		}

		// Check if word is complete
		if len(g.input) == len(g.word) {
			duration := time.Since(g.startTime)
			accuracy := 100.0 * (float64(len(g.word)-g.mistakes) / float64(len(g.word)))
			wpm := float64(len(g.word)) / duration.Minutes() / 5.0

			g.showResults(duration, accuracy, wpm)

			// Reset for next word
			g.word = words[rand.Intn(len(words))]
			g.input = ""
			g.mistakes = 0
			g.startTime = time.Now()
		}
	}

	return true
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game, err := newGame()
	if err != nil {
		fmt.Printf("Error initializing game: %v\n", err)
		return
	}
	defer game.cleanup()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Handle signals in a separate goroutine
	go func() {
		<-sigChan
		close(game.done)
	}()

	var b [1]byte
	for {
		select {
		case <-game.done:
			return
		default:
			game.render()

			_, err := os.Stdin.Read(b[:])
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				return
			}

			if !game.handleInput(b[0]) {
				return
			}
		}
	}
}
