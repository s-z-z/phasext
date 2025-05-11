package main

import (
	"github.com/suzi1037/pcmd/app/code-gen/internal"
)

func main() {

	dirs := []string{
		"app/demo/apis",
		//"test/pcmd",
	}

	for _, d := range dirs {
		println("generate: ", d)
		internal.DirLoopGenerate(d)
	}
}
