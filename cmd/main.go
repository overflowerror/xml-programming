package main

import (
	"flag"
	"io"
	"os"
	"xml-programming/internal/analysis"
	"xml-programming/internal/parser"
	"xml-programming/internal/vm"
)

func main() {
	flag.Parse()
	filename := flag.Arg(0)
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	program, err := parser.Parse(content)
	if err != nil {
		panic(err)
	}

	//json, _ := json.MarshalIndent(program, "", "  ")
	//fmt.Println(string(json))

	err = analysis.StaticAnalysis(program)
	if err != nil {
		panic(err)
	}

	vm.Run(program)
}