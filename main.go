package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/kabironline/monke/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf(`███████╗   ███████╗   ██╗    ██╗    █████╗ 
██╔════╝   ██╔════╝   ██║    ██║   ██╔══██╗
█████╗     █████╗     ██║ █╗ ██║   ███████║
██╔══╝     ██╔══╝     ██║███╗██║   ██╔══██║
███████╗██╗███████╗██╗╚███╔███╔╝██╗██║  ██║
╚══════╝╚═╝╚══════╝╚═╝ ╚══╝╚══╝ ╚═╝╚═╝  ╚═╝
                                           `)
	fmt.Printf("\nWelcome %s, To the Expression Evaluator on Web Assembly\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
