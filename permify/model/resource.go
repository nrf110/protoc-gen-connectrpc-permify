package model

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	permifyv1 "github.com/nrf110/connectrpc-permify/gen/nrf110/permify/v1"
	"github.com/nrf110/protoc-gen-connectrpc-permify/permify/util"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
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
	resource.checksFromResources(resource.Path, nestingLevel)
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
			IdPath:         findPath(plugin, pb, permifyv1.E_ResourceId, NewRootPathBuilder("resource", file)),
			TenantIdPath:   findPath(plugin, pb, permifyv1.E_TenantId, NewRootPathBuilder("resource", file)),
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

func findPath(plugin *protogen.Plugin, pb *protogen.Message, extension *protoimpl.ExtensionInfo, path *PathBuilder) *Path {
	for _, field := range pb.Fields {
		if util.GetBoolExtension(field.Desc, extension) {
			if !util.IsIdField(field) {
				plugin.Error(fmt.Errorf("%s must be a string or integer type", field.GoName))
			}
			return path.AddField(field).Build()
		}

		if util.IsMessage(field) {
			if field.Desc.IsList() || field.Desc.IsMap() {
				continue
			}

			if result := findPath(plugin, field.Message, extension, NewPathBuilder(path).AddField(field)); result != nil {
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

func (resource *Resource) checksFromResources(remainingPath *Path, nestingLevel int) {
	file := resource.file
	if remainingPath.Child != nil {
		if remainingPath.Child.Child != nil || resource.IdPath != nil || resource.TenantIdPath != nil {
			varName := util.VariableName()
			file.P(util.Indent(nestingLevel), "for _, ", varName, " := range ", remainingPath.Path, " {")
			resource.checksFromResources(remainingPath.Child.WithPrefix(varName), nestingLevel+1)
		} else {
			file.P(util.Indent(nestingLevel), "for range ", remainingPath.Path, " {")
			resource.checksFromResources(remainingPath.Child, nestingLevel+1)
		}
		file.P(util.Indent(nestingLevel), "}")
	} else {
		if resource.IdPath != nil || resource.TenantIdPath != nil || len(resource.AttributePaths) > 0 {
			file.P(util.Indent(nestingLevel), "resource := ", remainingPath.Path)
		}

		var idPath string
		file.P(util.Indent(nestingLevel), `var id string`)
		if resource.IdPath != nil {
			util.Log.Printf("rendering id path for %v", resource.IdPath)
			resource.renderIdPath(resource.IdPath, nestingLevel, "id")
		}
		file.P(util.Indent(nestingLevel), `tenantId := "default"`)
		if resource.TenantIdPath != nil {
			resource.renderIdPath(resource.TenantIdPath, nestingLevel, "tenantId")
		}
		resource.renderAttributes(nestingLevel)
		resource.renderCheck(nestingLevel, idPath)
		file.P(util.Indent(nestingLevel), "checks = append(checks, check)")
	}
}

func (resource *Resource) renderIdPath(path *Path, nestingLevel int, varName string) {
	var (
		sb   strings.Builder
		file = resource.file
	)
	file.P(util.Indent(nestingLevel), "if ", resource.renderNilChecks(&sb, path, ""), "{")
	file.P(util.Indent(nestingLevel+1), varName, " = ", path)
	file.P(util.Indent(nestingLevel), "}")
}

func (resource *Resource) renderNilChecks(sb *strings.Builder, remainingPath *Path, checkedPath string) string {
	var cumulativePath string
	if checkedPath == "" {
		cumulativePath = remainingPath.Path
	} else {
		cumulativePath = fmt.Sprintf("%s.%s", checkedPath, remainingPath.Path)
	}

	sb.WriteString(cumulativePath)
	if remainingPath.Child != nil {
		sb.WriteString(" != nil &&")
		return resource.renderNilChecks(sb, remainingPath.Child, cumulativePath)
	} else {
		sb.WriteString(` != ""`)
	}
	return sb.String()
}

func (resource *Resource) renderAttributes(nestingLevel int) {
	file := resource.file

	file.P(util.Indent(nestingLevel), "attributes := make(map[string]any)")
	keys := maps.Keys(resource.AttributePaths)
	sortedKeys := slices.Sorted(keys)
	for _, name := range sortedKeys {
		resource.renderAttribute(name, resource.AttributePaths[name], nestingLevel, false)
	}
}

func (resource *Resource) renderAttribute(name string, remainingPath *Path, nestingLevel int, isNested bool) {
	file := resource.file

	if remainingPath.Child != nil {
		// We have a nested collection, need to collect values into a slice
		varName := util.VariableName()
		if !isNested {
			// First time entering a collection, initialize the slice
			file.P(util.Indent(nestingLevel), "var ", name, "Values []any")
		}
		file.P(util.Indent(nestingLevel), "for _, ", varName, " := range ", remainingPath.Path, "{")
		resource.renderAttribute(name, remainingPath.Child.WithPrefix(varName), nestingLevel+1, true)
		file.P(util.Indent(nestingLevel), "}")
		if !isNested {
			// After collecting all values, assign to attributes
			file.P(util.Indent(nestingLevel), `if len(`, name, `Values) > 0 {`)
			file.P(util.Indent(nestingLevel+1), `attributes["`, name, `"] = `, name, `Values`)
			file.P(util.Indent(nestingLevel), "}")
		}
	} else {
		if isNested {
			// We're inside a collection, append to the slice
			file.P(util.Indent(nestingLevel), name, `Values = append(`, name, `Values, `, remainingPath.Path, `)`)
		} else {
			// Direct assignment for non-nested attributes
			file.P(util.Indent(nestingLevel), `attributes["`, name, `"] = `, remainingPath.Path)
		}
	}
}

func (resource *Resource) renderCheck(nestingLevel int, idPath string) {
	file := resource.file

	file.P(util.Indent(nestingLevel), "check := pkg.Check {")
	file.P(util.Indent(nestingLevel+1), "TenantID:     tenantId,")
	file.P(util.Indent(nestingLevel+1), "Permission:   permission,")
	file.P(util.Indent(nestingLevel+1), "Entity: &pkg.Resource {")
	file.P(util.Indent(nestingLevel+2), `Type:       "`, resource.Type, `",`)
	file.P(util.Indent(nestingLevel+2), `ID:         id,`)
	file.P(util.Indent(nestingLevel+2), `Attributes: attributes,`)
	file.P(util.Indent(nestingLevel+1), "},")
	file.P(util.Indent(nestingLevel), "}")
}
