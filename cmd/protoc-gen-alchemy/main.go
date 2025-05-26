package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/wjiec/alchemy/cmd/protoc-gen-alchemy/internal/gengo"
)

func main() {
	protogen.Options{ParamFunc: flag.CommandLine.Set}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = gengo.SupportedFeatures
		for _, file := range gen.Files {
			if !file.Generate {
				continue
			}

			if err := gengo.GenerateFile(gen, file); err != nil {
				return err
			}
		}
		return nil
	})
}
