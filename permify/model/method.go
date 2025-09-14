package model

import (
	"fmt"

	permifyv1 "github.com/nrf110/connectrpc-permify/gen/nrf110/permify/v1"
	"github.com/nrf110/protoc-gen-connectrpc-permify/permify/util"
	"google.golang.org/protobuf/compiler/protogen"
)

type Method struct {
	file        *protogen.GeneratedFile
	IsPublic    bool
	Permission  string
	RequestType string
	Resource    *Resource
}

func NewMethod(plugin *protogen.Plugin, file *protogen.GeneratedFile, pb *protogen.Method) *Method {
	isPublic := util.GetBoolExtension(pb.Desc, permifyv1.E_Public)
	hasPermission, permission := util.GetStringExtension(pb.Desc, permifyv1.E_Permission)

	if !isPublic && !hasPermission {
		plugin.Error(fmt.Errorf("method %s in service %s must specify a permission", pb.GoName, pb.Parent.GoName))
	}

	resource := NewResource(plugin, file, pb.Input)
	if !isPublic && resource == nil {
		plugin.Error(fmt.Errorf("method %s in service %s must specify a resource", pb.GoName, pb.Parent.GoName))
	}

	method := Method{
		file:        file,
		IsPublic:    isPublic,
		Permission:  permission,
		RequestType: pb.Input.GoIdent.GoName,
		Resource:    resource,
	}

	return &method
}

func (method *Method) Generate() {
	method.file.P("func (req *", method.RequestType, ") GetChecks() pkg.CheckConfig {")
	if method.IsPublic {
		method.generatePublic()
	} else {
		method.generateChecks()
	}
	method.file.P("}")
}

func (method *Method) generatePublic() {
	file := method.file
	file.P(util.Indent(1), "return pkg.CheckConfig {")
	file.P(util.Indent(2), "Type:   pkg.PUBLIC,")
	file.P(util.Indent(2), `Checks: []pkg.Check{},`)
	file.P(util.Indent(1), "}")
}

func (method *Method) generateChecks() {
	file := method.file
	file.P(util.Indent(1), `permission := "`, method.Permission, `"`)
	file.P(util.Indent(1), "var checks []pkg.Check")
	if method.Resource != nil {
		method.Resource.Generate(1)
	}

	file.P(util.Indent(1), "return pkg.CheckConfig {")
	file.P(util.Indent(2), "Type:   pkg.SINGLE,")
	file.P(util.Indent(2), "Checks: checks,")
	file.P(util.Indent(1), "}")
}
