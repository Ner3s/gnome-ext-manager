package utils

import (
	"fmt"
	"strings"
	"time"
)

// ProgressBar represents a progress bar that can be updated
type ProgressBar struct {
	total      int
	current    int
	width      int
	done       chan bool
	isStopped  bool
	message    string
	updateFreq time.Duration
}

// NewProgressBar creates a new progress bar with the given total steps and width
func NewProgressBar(total, width int) *ProgressBar {
	return &ProgressBar{
		total:      total,
		current:    0,
		width:      width,
		done:       make(chan bool),
		updateFreq: 100 * time.Millisecond,
		message:    "Processing",
	}
}

// SetMessage sets the message to display with the progress bar
func (p *ProgressBar) SetMessage(message string) {
	p.message = message
}

// Increment increases the current progress by the specified amount
func (p *ProgressBar) Increment(amount int) {
	p.current += amount
	if p.current > p.total {
		p.current = p.total
	}
}

// Start begins the progress bar animation in a separate goroutine
func (p *ProgressBar) Start() {
	go func() {
		for !p.isStopped {
			p.render()
			select {
			case <-p.done:
				return
			case <-time.After(p.updateFreq):
				// Continue
			}
		}
	}()
}

// Stop stops the progress bar animation and prints a newline
func (p *ProgressBar) Stop() {
	if !p.isStopped {
		p.isStopped = true
		close(p.done)
		fmt.Println()
	}
}

// SetProgress sets the current progress to a specific value
func (p *ProgressBar) SetProgress(current int) {
	p.current = current
	if p.current > p.total {
		p.current = p.total
	}
}

// render displays the current state of the progress bar
func (p *ProgressBar) render() {
	percent := float64(p.current) / float64(p.total)
	filled := int(float64(p.width) * percent)

	// Ensure filled doesn't exceed width
	if filled > p.width {
		filled = p.width
	}

	bar := fmt.Sprintf("\r%s [%s%s] %.1f%%",
		p.message,
		strings.Repeat("█", filled),
		strings.Repeat(" ", p.width-filled),
		percent*100,
	)

	fmt.Print(bar)
}

// SimulateIndeterminateProgress shows an indeterminate progress animation
// and executes the given function in the background
func SimulateIndeterminateProgress(message string, fn func() error) error {
	done := make(chan bool)
	errCh := make(chan error)

	// Start the animation
	go func() {
		spinChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-done:
				fmt.Print("\r                                                  \r")
				return
			default:
				fmt.Printf("\r%s %s", spinChars[i], message)
				i = (i + 1) % len(spinChars)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Run the actual function
	go func() {
		err := fn()
		errCh <- err
	}()

	// Wait for completion
	err := <-errCh
	close(done)
	fmt.Println()

	return err
}
