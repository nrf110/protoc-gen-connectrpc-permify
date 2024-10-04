package model

import (
	"fmt"
	permitv1 "github.com/nrf110/connectrpc-permit/gen/nrf110/permit/v1"
	"github.com/nrf110/protoc-gen-connectrpc-permit/permit/util"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

type Method struct {
	IsPublic      bool
	Action        string
	RequestType   string
	Resource      *Resource
	BulkCheckMode permitv1.BulkCheckMode
}

func NewMethod(plugin *protogen.Plugin, pb *protogen.Method) *Method {
	isPublic := util.GetBoolExtension(pb.Desc, permitv1.E_Public)
	hasAction, action := util.GetStringExtension(pb.Desc, permitv1.E_Action)

	if !isPublic && !hasAction {
		plugin.Error(fmt.Errorf("method %s in service %s must specify a permission action", pb.GoName, pb.Parent.GoName))
	}

	resource := NewResource(plugin, pb.Input)
	if !isPublic && resource == nil {
		plugin.Error(fmt.Errorf("method %s in service %s must specify a resource", pb.GoName, pb.Parent.GoName))
	}

	method := Method{
		IsPublic:    isPublic,
		Action:      action,
		RequestType: pb.Input.GoIdent.GoName,
		Resource:    resource,
	}

	if resource.Path.Child != nil || resource.IdPath != nil && resource.IdPath.Child != nil {
		opts := pb.Desc.Options()
		if proto.HasExtension(opts, permitv1.E_BulkCheckMode) {
			if mode, hasMode := proto.GetExtension(opts, permitv1.E_BulkCheckMode).(permitv1.BulkCheckMode); hasMode {
				method.BulkCheckMode = mode
			}
		}
	}

	return &method
}

func (method *Method) Generate(file *protogen.GeneratedFile) {
	file.P("func (req *", method.RequestType, ") GetChecks() connectrpc_permit.CheckConfig {")
	if method.IsPublic {
		method.generatePublic(file)
	} else {
		method.generateChecks(file)
	}
	file.P("}")
}

func (method *Method) generatePublic(file *protogen.GeneratedFile) {
	file.P(util.Indent(1), "return connectrpc_permit.CheckConfig {")
	file.P(util.Indent(2), "Type:   connectrpc_permit.PUBLIC,")
	method.generateBulkCheckMode(file)
	file.P(util.Indent(2), `Checks: []connectrpc_permit.Check{},`)
	file.P(util.Indent(1), "}")
}

func (method *Method) generateChecks(file *protogen.GeneratedFile) {
	file.P(util.Indent(1), `action := "`, method.Action, `"`)
	file.P(util.Indent(1), "var checks []connectrpc_permit.Check")
	if method.Resource != nil {
		method.Resource.Generate(file, 1)
	}
	file.P(util.Indent(1), "checkType := connectrpc_permit.SINGLE")
	file.P(util.Indent(1), "if len(checks) > 1 {")
	file.P(util.Indent(2), "checkType = connectrpc_permit.BULK")
	file.P(util.Indent(1), "}")
	file.P(util.Indent(1), "return connectrpc_permit.CheckConfig {")
	file.P(util.Indent(2), "Type:   checkType,")
	method.generateBulkCheckMode(file)
	file.P(util.Indent(2), "Checks: checks,")
	file.P(util.Indent(1), "}")
}

func (method *Method) generateBulkCheckMode(file *protogen.GeneratedFile) {
	value := "connectrpc_permit.AllOf"
	if method.BulkCheckMode == permitv1.BulkCheckMode_any_of {
		value = "connectrpc_permit.AnyOf"
	}
	file.P(util.Indent(2), "Mode:   ", value, ",")
}
