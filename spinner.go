package bget

import (
	"fmt"
	"os"
	"time"

	"github.com/tj/go-spin"
)

func ShowSpinnerWhile(message string) func() {
	spinner := spin.New()
	spinner.Set(spin.Box1)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				// reprint new spinner state
				fmt.Fprintf(os.Stderr, "\r%s %s", spinner.Next(), message)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	return func() {
		done <- true
		// remove spinner
		fmt.Fprint(os.Stderr, "\n\033[1A\033[K")
	}
}
