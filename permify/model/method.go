package model

import (
	"fmt"

	permifyv1 "github.com/nrf110/connectrpc-permify/gen/nrf110/permify/v1"
	connectpermify "github.com/nrf110/connectrpc-permify/pkg"
	"github.com/nrf110/protoc-gen-connectrpc-permify/permify/util"
	"google.golang.org/protobuf/compiler/protogen"
)

type Method struct {
	file        *protogen.GeneratedFile
	IsPublic    bool
	Action      string
	RequestType string
	Resource    *Resource
	CheckType   connectpermify.CheckType
}

func NewMethod(plugin *protogen.Plugin, file *protogen.GeneratedFile, pb *protogen.Method) *Method {
	isPublic := util.GetBoolExtension(pb.Desc, permifyv1.E_Public)
	hasAction, action := util.GetStringExtension(pb.Desc, permifyv1.E_Action)

	if !isPublic && !hasAction {
		plugin.Error(fmt.Errorf("method %s in service %s must specify a permission action", pb.GoName, pb.Parent.GoName))
	}

	resource := NewResource(plugin, file, pb.Input)
	if !isPublic && resource == nil {
		plugin.Error(fmt.Errorf("method %s in service %s must specify a resource", pb.GoName, pb.Parent.GoName))
	}

	method := Method{
		file:        file,
		IsPublic:    isPublic,
		Action:      action,
		RequestType: pb.Input.GoIdent.GoName,
		Resource:    resource,
	}

	return &method
}

func (method *Method) Generate() {
	method.file.P("func (req *", method.RequestType, ") GetChecks() connectpermify.CheckConfig {")
	if method.IsPublic {
		method.generatePublic()
	} else {
		method.generateChecks()
	}
	method.file.P("}")
}

func (method *Method) generatePublic() {
	file := method.file
	file.P(util.Indent(1), "return connectpermify.CheckConfig {")
	file.P(util.Indent(2), "Type:   connectpermify.PUBLIC,")
	file.P(util.Indent(2), `Checks: []connectpermify.Check{},`)
	file.P(util.Indent(1), "}")
}

func (method *Method) generateChecks() {
	file := method.file
	file.P(util.Indent(1), `action := "`, method.Action, `"`)
	file.P(util.Indent(1), "var checks []connectpermify.Check")
	if method.Resource != nil {
		method.Resource.Generate(1)
	}

	file.P(util.Indent(1), "return connectpermify.CheckConfig {")
	file.P(util.Indent(2), "Type:   connectpermify.SINGLE,")
	file.P(util.Indent(2), "Checks: checks,")
	file.P(util.Indent(1), "}")
}
