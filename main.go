package main

import (
	"fmt"
	"os"

	"github.com/jilleJr/flog/pkg/logparser"
)

func main() {
	p := logparser.NewIOParser(os.Stdin)

	for p.Scan() {
		log := p.ParsedLog()
		if log.Level > logparser.LevelInformation {
			fmt.Println(log.String)
		}
	}
}
