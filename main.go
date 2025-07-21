package main

import (
	"github.com/nrf110/protoc-gen-connectrpc-permify/permify/model"
	"github.com/nrf110/protoc-gen-connectrpc-permify/permify/util"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	util.InitLogger()
	defer util.LogFile.Close()

	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		for _, f := range plugin.Files {
			if !f.Generate {
				continue
			}
			buildModel(plugin, f)
		}
		return nil
	})
}

func buildModel(plugin *protogen.Plugin, file *protogen.File) {
	filename := file.GeneratedFilenamePrefix + "_permit.pb.go"
	gen := plugin.NewGeneratedFile(filename, file.GoImportPath)
	gen.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "permifyv1",
		GoImportPath: "github.com/nrf110/connectrpc-permify/pkg",
	})

	gen.P("package " + file.GoPackageName)
	gen.P("")

	for _, service := range file.Services {
		svc := model.NewService(plugin, gen, service)
		svc.Generate()
	}
}
