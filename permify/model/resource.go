package model

import (
	"fmt"
	"strings"

	permifyv1 "github.com/nrf110/connectrpc-permify/gen/nrf110/permify/v1"
	"github.com/nrf110/protoc-gen-connectrpc-permify/permify/util"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Resource struct {
	file           *protogen.GeneratedFile
	GoName         string
	Type           string
	Path           *Path
	IdPath         *Path
	TenantIdPath   *Path
	AttributePaths map[string]*Path
}

func NewResource(plugin *protogen.Plugin, file *protogen.GeneratedFile, pb *protogen.Message) *Resource {
	util.Log.Println("finding resource")
	return findResourcePath(plugin, file, pb, NewRootPathBuilder("req", file))
}

func (resource *Resource) Generate(nestingLevel int) {
	resource.checksFromResources(resource.Path, nestingLevel, make(map[string]bool))
}

func findResourcePath(plugin *protogen.Plugin, file *protogen.GeneratedFile, pb *protogen.Message, path *PathBuilder) *Resource {
	options := pb.Desc.Options()
	if proto.HasExtension(options, permifyv1.E_ResourceType) {
		util.Log.Println("found resource type in", pb.GoIdent.String())
		resourceType := proto.GetExtension(options, permifyv1.E_ResourceType).(string)
		return &Resource{
			file:           file,
			GoName:         pb.GoIdent.GoName,
			Type:           resourceType,
			Path:           path.Build(),
			IdPath:         findIdPath(plugin, pb, NewRootPathBuilder("resource", file)),
			TenantIdPath:   findTenantPath(pb, NewRootPathBuilder("resource", file)),
			AttributePaths: findAttributes(plugin, pb, NewRootPathBuilder("resource", file), make(map[string]*Path)),
		}
	}

	for _, field := range pb.Fields {
		util.Log.Println("checking field", field.GoName)
		if !field.Desc.HasOptionalKeyword() {
			if util.IsMessageValueMap(field) {
				util.Log.Println(field.GoName, "is a map of messages")
				result := findResourcePath(plugin, file, util.GetMapFieldValue(field), NewPathBuilder(path.AddField(field)))
				if result != nil {
					return result
				}
			}

			if util.IsMessage(field) {
				util.Log.Println(field.GoName, "is a message")
				if field.Desc.IsList() {
					util.Log.Println(field.GoName, "is a repeated message")
					if result := findResourcePath(plugin, file, field.Message, NewPathBuilder(path.AddField(field))); result != nil {
						return result
					}
				}

				if result := findResourcePath(plugin, file, field.Message, path.AddField(field)); result != nil {
					return result
				}
			}
		}
	}

	return nil
}

func findIdPath(plugin *protogen.Plugin, pb *protogen.Message, path *PathBuilder) *Path {
	for _, field := range pb.Fields {
		if util.GetBoolExtension(field.Desc, permifyv1.E_ResourceId) {
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

func findAttributes(plugin *protogen.Plugin, pb *protogen.Message, path *PathBuilder, accum map[string]*Path) map[string]*Path {
	for _, field := range pb.Fields {
		if found, attributeName := util.GetStringExtension(field.Desc, permifyv1.E_AttributeName); found {
			accum[attributeName] = path.AddField(field).Build()
			continue
		}

		if util.IsMessageValueMap(field) {
			findAttributes(plugin, util.GetMapFieldValue(field), NewPathBuilder(path.AddField(field)), accum)
			continue
		}

		if util.IsMessage(field) {
			if field.Desc.IsList() {
				findAttributes(plugin, field.Message, NewPathBuilder(path.AddField(field)), accum)
			} else {
				findAttributes(plugin, field.Message, path.AddField(field), accum)
			}
		}
	}
	return accum
}

func findTenantPath(pb *protogen.Message, path *PathBuilder) *Path {
	for _, field := range pb.Fields {
		if field.Desc.IsList() || field.Desc.IsMap() {
			continue
		}

		if util.GetBoolExtension(field.Desc, permifyv1.E_TenantId) {
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

func (resource *Resource) checksFromResources(remainingPath *Path, nestingLevel int, usedLoopVars map[string]bool) {
	file := resource.file
	if remainingPath.Child != nil {
		if remainingPath.Child.Child != nil || resource.IdPath != nil || resource.TenantIdPath != nil {
			varName := loopVar(usedLoopVars)
			file.P(util.Indent(nestingLevel), "for _, ", varName, " := range ", remainingPath.Path, " {")
			resource.checksFromResources(remainingPath.Child.WithPrefix(varName), nestingLevel+1, usedLoopVars)
		} else {
			file.P(util.Indent(nestingLevel), "for range ", remainingPath.Path, " {")
			resource.checksFromResources(remainingPath.Child, nestingLevel+1, usedLoopVars)
		}
		file.P(util.Indent(nestingLevel), "}")
	} else {
		if resource.IdPath != nil || resource.TenantIdPath != nil || len(resource.AttributePaths) > 0 {
			file.P(util.Indent(nestingLevel), "resource := ", remainingPath.Path)
		}

		var idPath string
		if resource.IdPath != nil {
			idPath = resource.renderId(resource.IdPath, nestingLevel, usedLoopVars)
		}
		file.P(util.Indent(nestingLevel), `tenantId := "default"`)
		if resource.TenantIdPath != nil {
			resource.renderTenantId(resource.TenantIdPath, nestingLevel, usedLoopVars)
		}
		resource.renderAttributes(nestingLevel, usedLoopVars)
		resource.renderCheck(nestingLevel, idPath)
		file.P(util.Indent(nestingLevel), "checks = append(checks, check)")
	}
}

func (resource *Resource) renderId(remainingPath *Path, nestingLevel int, usedLoopVars map[string]bool) string {
	file := resource.file

	if remainingPath != nil && remainingPath.Child != nil {
		varName := loopVar(usedLoopVars)
		file.P(util.Indent(nestingLevel), "for _, ", varName, " := range ", remainingPath.Path, "{")
		result := resource.renderId(remainingPath.Child.WithPrefix(varName), nestingLevel+1, usedLoopVars)
		file.P(util.Indent(nestingLevel), "}")
		return result
	} else {
		return remainingPath.Path
	}
}

func (resource *Resource) renderTenantId(remainingPath *Path, nestingLevel int, usedLoopVars map[string]bool) {
	var (
		sb   strings.Builder
		file = resource.file
	)
	file.P(util.Indent(nestingLevel), "if ", resource.checkTenantId(&sb, resource.TenantIdPath, ""), "{")
	file.P(util.Indent(nestingLevel+1), "tenantId = ", resource.tenantIdPath())
	file.P(util.Indent(nestingLevel), "}")
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

func (resource *Resource) renderAttributes(nestingLevel int, usedLoopVars map[string]bool) {
	file := resource.file

	file.P(util.Indent(nestingLevel), "attributes := make(map[string]any)")
	for name, path := range resource.AttributePaths {
		resource.renderAttribute(name, path, nestingLevel, usedLoopVars)
	}
}

func (resource *Resource) renderAttribute(name string, remainingPath *Path, nestingLevel int, usedLoopVars map[string]bool) {
	file := resource.file

	if remainingPath.Child != nil {
		varName := loopVar(usedLoopVars)
		file.P(util.Indent(nestingLevel), "for _, ", varName, " := range ", remainingPath.Path, "{")
		resource.renderAttribute(name, remainingPath.Child.WithPrefix(varName), nestingLevel+1, usedLoopVars)
		file.P(util.Indent(nestingLevel), "}")
	} else {
		file.P(util.Indent(nestingLevel+1), `attributes["`, name, `"] = `, remainingPath.Path)
	}
}

func (resource *Resource) renderCheck(nestingLevel int, idPath string) {
	file := resource.file

	file.P(util.Indent(nestingLevel), "check := pkg.Check {")
	file.P(util.Indent(nestingLevel+1), "TenantID:     tenantId,")
	file.P(util.Indent(nestingLevel+1), "Permission:   permission,")
	file.P(util.Indent(nestingLevel+1), "Entity: &pkg.Resource {")
	file.P(util.Indent(nestingLevel+2), `Type:       "`, resource.Type, `",`)
	if idPath != "" {
		file.P(util.Indent(nestingLevel+2), "ID:        ", idPath, ",")
	}
	file.P(util.Indent(nestingLevel+2), "Attributes: attributes,")
	file.P(util.Indent(nestingLevel+1), "},")
	file.P(util.Indent(nestingLevel), "}")
}

func loopVar(inUse map[string]bool) string {
	for {
		name := util.VariableName()
		if exists := inUse[name]; !exists {
			inUse[name] = true
			return name
		}
	}
}
