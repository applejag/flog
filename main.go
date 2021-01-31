package main

import (
	"os"

	"github.com/jilleJr/flog/pkg/logparser"
)

func main() {
	p := logparser.NewIOParser(os.Stdin)
	printer := NewConsolePrinter(&p)

	for printer.Next() {
	}
}
