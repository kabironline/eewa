package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kabironline/monke/lexer"
	"github.com/kabironline/monke/tokens"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		for tok := l.NextToken(); tok.Type != tokens.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
