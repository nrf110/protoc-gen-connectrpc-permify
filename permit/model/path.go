package model

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

type PathBuilder struct {
	fields []string
	parent *PathBuilder
}

func NewPathBuilder(parent *PathBuilder) *PathBuilder {
	return &PathBuilder{
		parent: parent,
	}
}

func NewRootPathBuilder(root string) *PathBuilder {
	return &PathBuilder{
		parent: nil,
		fields: []string{root},
	}
}

func (node *PathBuilder) AddField(field *protogen.Field) *PathBuilder {
	return &PathBuilder{
		fields: append(node.fields, field.GoName),
		parent: node.parent,
	}
}

func (node *PathBuilder) Path() string {
	return strings.Join(node.fields, ".")
}

func (node *PathBuilder) Build() *Path {
	return walk(node, nil)
}

func walk(currentNode *PathBuilder, path *Path) *Path {
	if currentNode.parent != nil {
		return walk(currentNode.parent, &Path{
			Path:  currentNode.Path(),
			Child: path,
		})
	}
	return &Path{
		Path:  currentNode.Path(),
		Child: path,
	}
}

type Path struct {
	Path  string
	Child *Path
}

func (path *Path) WithPrefix(prefix string) *Path {
	if path.Path == "" {
		path.Path = prefix
	} else {
		path.Path = fmt.Sprintf("%s.%s", prefix, path.Path)
	}
	return path
}
