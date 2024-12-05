package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/voidwyrm-2/ferrite/lexer"
	"github.com/voidwyrm-2/ferrite/runtime/interpreter"
	"github.com/voidwyrm-2/ferrite/terminal"
	"github.com/voidwyrm-2/ferrite/utils"
)

//go:embed version.txt
var version string

func ferriteRepl(useTermx bool) {
	i := interpreter.New()

	if useTermx {
		term := terminal.New("> ")
		defer term.Deactivate()
		term.Print("type '!exit' to end the repl\n")

		for {
			line, err := term.ReadLine()
			if err == io.EOF {
				break
			}

			if strings.TrimSpace(line) == "!exit" {
				break
			}

			l := lexer.New(line)
			tokens, err := l.Lex()
			if err != nil {
				term.Print(err.Error() + "\n")
				continue
			}

			if err := i.InterpretTokens(tokens); err != nil {
				if err.Error() != "BYE" {
					term.Print(err.Error() + "\n")
				}
			}
		}
	} else {
		fmt.Println("type '!exit' to end the repl")
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			scanner.Scan()
			input := scanner.Text()
			if strings.TrimSpace(input) == "!exit" {
				break
			}

			l := lexer.New(input)
			tokens, err := l.Lex()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			if err := i.InterpretTokens(tokens); err != nil {
				if err.Error() != "BYE" {
					fmt.Println(err.Error())
				}
			}
		}
	}
}

func main() {
	version = strings.TrimSpace(version)

	parser := argparse.NewParser("ferrite", "The Ferrite interpreter")

	showVersion := parser.Flag("v", "version", &argparse.Options{Required: false, Default: false, Help: "Shows the current Pasta version"})
	run := parser.NewCommand("run", "Run a Ferrite('.fer') file")
	rfile := run.String("f", "file", &argparse.Options{Required: true, Help: "The Ferrite file to run"})

	repl := parser.NewCommand("repl", "Run the Ferrite REPL")

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	if *showVersion {
		fmt.Println(version)
		return
	} else if run.Happened() {
		if content, err := utils.ReadFile(*rfile); err != nil {
			fmt.Println(err.Error())
			return
		} else {
			i := interpreter.New()
			if err := i.Interpret(content); err != nil {
				fmt.Println(err.Error())
			}
		}
	} else if repl.Happened() {
		ferriteRepl(runtime.GOOS != "wasip1")
	}
}
