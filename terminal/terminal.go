package terminal

import (
	"io"
	"log"
	"os"

	"golang.org/x/term"
)

type Terminal struct {
	state *term.State
	term  *term.Terminal
}

func New(prompt string) Terminal {
	state, err := term.MakeRaw(0)
	if err != nil {
		log.Fatalln("setting stdin to raw:", err)
	}

	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	term := term.NewTerminal(screen, "")
	term.SetPrompt(prompt)

	return Terminal{state: state, term: term}
}

func (t *Terminal) Deactivate() {
	if err := term.Restore(0, t.state); err != nil {
		log.Println("warning, failed to restore terminal:", err)
	}
}

func (t *Terminal) Activate() {
	if state, err := term.MakeRaw(0); err != nil {
		log.Println("warning, failed to restore terminal:", err)
	} else {
		t.state = state
	}
}

func (t *Terminal) ReadLine() (string, error) {
	line, err := t.term.ReadLine()
	return line, err
}

func (t *Terminal) Print(s string) error {
	_, err := t.term.Write([]byte(s))
	return err
}
