package model

import (
	"google.golang.org/protobuf/compiler/protogen"
)

type Service struct {
	file    *protogen.GeneratedFile
	Methods []*Method
}

func NewService(plugin *protogen.Plugin, file *protogen.GeneratedFile, pb *protogen.Service) *Service {
	var methods []*Method
	for _, method := range pb.Methods {
		methods = append(methods, NewMethod(plugin, file, method))
	}
	return &Service{
		file:    file,
		Methods: methods,
	}
}

func (service *Service) Generate() {
	for _, method := range service.Methods {
		method.Generate()
		service.file.P()
	}
}
