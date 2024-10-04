package model

import (
	"google.golang.org/protobuf/compiler/protogen"
)

type Service struct {
	Methods []*Method
}

func NewService(plugin *protogen.Plugin, pb *protogen.Service) *Service {
	var methods []*Method
	for _, method := range pb.Methods {
		methods = append(methods, NewMethod(plugin, method))
	}
	return &Service{Methods: methods}
}

func (service *Service) Generate(file *protogen.GeneratedFile) {
	for _, method := range service.Methods {
		method.Generate(file)
		file.P()
	}
}
