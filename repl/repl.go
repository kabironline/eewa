package repl

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/kabironline/monke/evaluator"
	"github.com/kabironline/monke/lexer"
	"github.com/kabironline/monke/object"
	"github.com/kabironline/monke/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		startTime := time.Now()
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
		io.WriteString(out, fmt.Sprintf("Operation took %s\n", time.Now().Sub(startTime).String()))
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! Someone made an Oopsie!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
