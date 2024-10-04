package model

import (
	"fmt"
	permitv1 "github.com/nrf110/connectrpc-permit/gen/nrf110/permit/v1"
	"github.com/nrf110/protoc-gen-connectrpc-permit/permit/util"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

type Resource struct {
	GoName       string
	Type         string
	Path         *Path
	IdPath       *Path
	TenantIdPath *Path
}

func NewResource(plugin *protogen.Plugin, pb *protogen.Message) *Resource {
	util.Log.Println("finding resource")
	return findResourcePath(plugin, pb, NewRootPathBuilder("req"))
}

func findResourcePath(plugin *protogen.Plugin, pb *protogen.Message, path *PathBuilder) *Resource {
	options := pb.Desc.Options()
	if proto.HasExtension(options, permitv1.E_ResourceType) {
		util.Log.Println("found resource type in", pb.GoIdent.String())
		resourceType := proto.GetExtension(options, permitv1.E_ResourceType).(string)
		return &Resource{
			GoName:       pb.GoIdent.GoName,
			Type:         resourceType,
			Path:         path.Build(),
			IdPath:       findIdPath(plugin, pb, NewRootPathBuilder("resource")),
			TenantIdPath: findTenantPath(pb, NewRootPathBuilder("resource")),
		}
	}

	for _, field := range pb.Fields {
		util.Log.Println("checking field", field.GoName)
		if !field.Desc.HasOptionalKeyword() {
			if util.IsMessageValueMap(field) {
				util.Log.Println(field.GoName, "is a map of messages")
				result := findResourcePath(plugin, util.GetMapFieldValue(field), NewPathBuilder(path.AddField(field)))
				if result != nil {
					return result
				}
			}

			if util.IsMessage(field) {
				util.Log.Println(field.GoName, "is a message")
				if field.Desc.IsList() {
					util.Log.Println(field.GoName, "is a repeated message")
					if result := findResourcePath(plugin, field.Message, NewPathBuilder(path.AddField(field))); result != nil {
						return result
					}
				}

				if result := findResourcePath(plugin, field.Message, path.AddField(field)); result != nil {
					return result
				}
			}
		}
	}

	return nil
}

func findIdPath(plugin *protogen.Plugin, pb *protogen.Message, path *PathBuilder) *Path {
	for _, field := range pb.Fields {
		if util.GetBoolExtension(field.Desc, permitv1.E_ResourceId) {
			if field.Desc.IsList() {
				if !util.IsIdField(field) {
					plugin.Error(fmt.Errorf("%s must be a string, enum, or integer type or a list or map of those types", field.GoName))
				}
				return NewPathBuilder(path.AddField(field)).Build()
			}

			if field.Desc.IsMap() {
				if !util.IsIdKind(field.Desc.MapValue().Kind()) {
					plugin.Error(fmt.Errorf("%s must be a string, enum, or integer type or a list or map of those types", field.GoName))
				}
				return NewPathBuilder(path.AddField(field)).Build()
			}
			return path.AddField(field).Build()
		}

		if util.IsMessageValueMap(field) {
			if result := findIdPath(plugin, util.GetMapFieldValue(field), NewPathBuilder(path.AddField(field))); result != nil {
				return result
			}
		}

		if util.IsMessage(field) {
			if field.Desc.IsList() {
				if result := findIdPath(plugin, field.Message, NewPathBuilder(path.AddField(field))); result != nil {
					return result
				}
			}

			if result := findIdPath(plugin, field.Message, path.AddField(field)); result != nil {
				return result
			}
		}
	}

	return nil
}

func findTenantPath(pb *protogen.Message, path *PathBuilder) *Path {
	for _, field := range pb.Fields {
		if field.Desc.IsList() || field.Desc.IsMap() {
			continue
		}

		if util.GetBoolExtension(field.Desc, permitv1.E_TenantId) {
			return path.AddField(field).Build()
		}

		if field.Desc.Kind() == protoreflect.MessageKind {
			if result := findTenantPath(field.Message, NewPathBuilder(path).AddField(field)); result != nil {
				return result
			}
		}
	}
	return nil
}

func (resource *Resource) Generate(file *protogen.GeneratedFile, nestingLevel int) {
	resource.checksFromResources(file, resource.Path, nestingLevel, make(map[string]bool))
}

func (resource *Resource) checksFromResources(file *protogen.GeneratedFile, remainingPath *Path, nestingLevel int, usedLoopVars map[string]bool) {
	if remainingPath.Child != nil {
		if remainingPath.Child.Child != nil || resource.IdPath != nil || resource.TenantIdPath != nil {
			varName := loopVar(usedLoopVars)
			file.P(util.Indent(nestingLevel), "for _, ", varName, " := range ", remainingPath.Path, " {")
			resource.checksFromResources(file, remainingPath.Child.WithPrefix(varName), nestingLevel+1, usedLoopVars)
		} else {
			file.P(util.Indent(nestingLevel), "for range ", remainingPath.Path, " {")
			resource.checksFromResources(file, remainingPath.Child, nestingLevel+1, usedLoopVars)
		}
		file.P(util.Indent(nestingLevel), "}")
	} else {
		if resource.IdPath != nil {
			file.P(util.Indent(nestingLevel), "resource := ", remainingPath.Path)
		}
		resource.checksFromIds(file, resource.IdPath, nestingLevel, usedLoopVars)
	}
}

func (resource *Resource) checksFromIds(file *protogen.GeneratedFile, remainingPath *Path, nestingLevel int, usedLoopVars map[string]bool) {
	if remainingPath != nil && remainingPath.Child != nil {
		varName := loopVar(usedLoopVars)
		file.P(util.Indent(nestingLevel), "for _, ", varName, " := range ", remainingPath.Path, "{")
		resource.checksFromIds(file, remainingPath.Child.WithPrefix(varName), nestingLevel+1, usedLoopVars)
		file.P(util.Indent(nestingLevel), "}")
	} else {
		file.P(util.Indent(nestingLevel), `tenantId := "default"`)
		if resource.TenantIdPath != nil {
			var sb strings.Builder
			file.P(util.Indent(nestingLevel), "if ", resource.checkTenantId(&sb, resource.TenantIdPath, ""), "{")
			file.P(util.Indent(nestingLevel+1), "tenantId = ", resource.tenantIdPath())
			file.P(util.Indent(nestingLevel), "}")
		}
		file.P(util.Indent(nestingLevel), "check := connectrpc_permit.Check {")
		file.P(util.Indent(nestingLevel+1), "Action:   action,")
		file.P(util.Indent(nestingLevel+1), "Resource: connectrpc_permit.Resource {")
		file.P(util.Indent(nestingLevel+2), `Type:   "`, resource.Type, `",`)
		if remainingPath != nil {
			file.P(util.Indent(nestingLevel+2), "Key:    ", remainingPath.Path, ",")
		}
		file.P(util.Indent(nestingLevel+2), "Tenant: tenantId,")
		file.P(util.Indent(nestingLevel+1), "},")
		file.P(util.Indent(nestingLevel), "}")
		file.P(util.Indent(nestingLevel), "checks = append(checks, check)")
	}
}

func (resource *Resource) checkTenantId(sb *strings.Builder, remainingPath *Path, checkedPath string) string {
	var cumulativePath string
	if checkedPath == "" {
		cumulativePath = remainingPath.Path
	} else {
		cumulativePath = fmt.Sprintf("%s.%s", checkedPath, remainingPath.Path)
	}

	sb.WriteString(cumulativePath)
	if remainingPath.Child != nil {
		sb.WriteString(" != nil &&")
		return resource.checkTenantId(sb, remainingPath.Child, cumulativePath)
	} else {
		sb.WriteString(` != ""`)
	}
	return sb.String()
}

func (resource *Resource) tenantIdPath() string {
	var sb strings.Builder
	currentPath := resource.TenantIdPath
	for currentPath != nil {
		sb.WriteString(currentPath.Path)
		if currentPath.Child != nil {
			sb.WriteString(".")
		}
		currentPath = currentPath.Child
	}
	return sb.String()
}

func loopVar(inUse map[string]bool) string {
	for {
		name := util.VariableName()
		if exists, _ := inUse[name]; !exists {
			inUse[name] = true
			return name
		}
	}
}
