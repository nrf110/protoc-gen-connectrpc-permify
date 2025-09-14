package main

import (
	"github.com/nrf110/protoc-gen-connectrpc-permify/permify/model"
	"github.com/nrf110/protoc-gen-connectrpc-permify/permify/util"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	util.InitLogger()
	defer util.LogFile.Close()

	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
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
	if len(file.Services) == 0 {
		return
	}

	// Reset the variable counter for each file to ensure deterministic output
	util.ResetVariableCounter()

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
