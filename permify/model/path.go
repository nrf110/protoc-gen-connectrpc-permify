package model

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

type fieldHolder struct {
	name  string
	field *protogen.Field
}

func (fh fieldHolder) VariableType(file *protogen.GeneratedFile) string {
	goType, pointer := variableType(file, fh.field)
	if pointer {
		return "*" + goType
	}
	return goType
}

func variableType(file *protogen.GeneratedFile, field *protogen.Field) (goType string, pointer bool) {
	if field.Desc.IsWeak() {
		return "struct{}", false
	}

	pointer = field.Desc.HasPresence()
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		goType = "bool"
	case protoreflect.EnumKind:
		goType = file.QualifiedGoIdent(field.Enum.GoIdent)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		goType = "int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		goType = "uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		goType = "int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		goType = "uint64"
	case protoreflect.FloatKind:
		goType = "float32"
	case protoreflect.DoubleKind:
		goType = "float64"
	case protoreflect.StringKind:
		goType = "string"
	case protoreflect.BytesKind:
		goType = "[]byte"
		pointer = false // rely on nullability of slices for presence
	case protoreflect.MessageKind, protoreflect.GroupKind:
		goType = "*" + file.QualifiedGoIdent(field.Message.GoIdent)
		pointer = false // pointer captured as part of the type
	}
	switch {
	case field.Desc.IsList():
		return "[]" + goType, false
	case field.Desc.IsMap():
		keyType, _ := variableType(file, field.Message.Fields[0])
		valType, _ := variableType(file, field.Message.Fields[1])
		return fmt.Sprintf("map[%v]%v", keyType, valType), false
	}
	return goType, pointer
}

type PathBuilder struct {
	file   *protogen.GeneratedFile
	fields []fieldHolder
	parent *PathBuilder
}

func NewPathBuilder(parent *PathBuilder) *PathBuilder {
	return &PathBuilder{
		file:   parent.file,
		parent: parent,
	}
}

func NewRootPathBuilder(root string, file *protogen.GeneratedFile) *PathBuilder {
	return &PathBuilder{
		file: file,
		fields: []fieldHolder{
			{
				name: root,
			},
		},
		parent: nil,
	}
}

func (node *PathBuilder) AddField(f *protogen.Field) *PathBuilder {
	return &PathBuilder{
		file: node.file,
		fields: append(node.fields, fieldHolder{
			name:  f.GoName,
			field: f,
		}),
		parent: node.parent,
	}
}

func (node *PathBuilder) Path() string {
	var sb strings.Builder
	lastIdx := len(node.fields) - 1
	for idx, f := range node.fields {
		sb.WriteString(f.name)
		if idx < lastIdx {
			sb.WriteString(".")
		}
	}
	return sb.String()
}

func (node *PathBuilder) VariableType() string {
	length := len(node.fields)
	if length > 0 {
		if holder := node.fields[length-1]; holder.field != nil {
			return holder.VariableType(node.file)
		}
	}
	return ""
}

func (node *PathBuilder) Build() *Path {
	return walk(node, nil)
}

func walk(currentNode *PathBuilder, path *Path) *Path {
	if currentNode.parent != nil {
		return walk(currentNode.parent, &Path{
			Path:         currentNode.Path(),
			VariableType: currentNode.VariableType(),
			Child:        path,
		})
	}
	return &Path{
		Path:         currentNode.Path(),
		VariableType: currentNode.VariableType(),
		Child:        path,
	}
}

type Path struct {
	Path         string
	VariableType string
	Child        *Path
}

func (path *Path) WithPrefix(prefix string) *Path {
	if path.Path == "" {
		path.Path = prefix
	} else {
		path.Path = fmt.Sprintf("%s.%s", prefix, path.Path)
	}
	return path
}
