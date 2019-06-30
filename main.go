// oops2core: convert Linux kernel crash report to a core file
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"opensource.go.fig.lu/oops2core/internal/elfcore32"
	"opensource.go.fig.lu/oops2core/internal/parser"
)

func handleFatalError() {
	if e := recover(); e != nil {
		fmt.Fprintln(os.Stderr, "Fatal error occured: ", e)
		os.Exit(1)
	}
}

func main() {
	defer handleFatalError()
	rawtext, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	ci, err := parser.ParseCrash(string(rawtext))
	if err != nil {
		panic(err)
	}
	_, err = os.Stdout.Write(elfcore32.NewElfcore(ci).Bytes())
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(elfcore32.NewElfcore(ci).Bytes())
}
