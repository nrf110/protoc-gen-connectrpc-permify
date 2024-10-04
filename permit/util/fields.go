package util

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"slices"
)

var IdKinds = []protoreflect.Kind{
	protoreflect.EnumKind,
	protoreflect.Int32Kind,
	protoreflect.Sint32Kind,
	protoreflect.Uint32Kind,
	protoreflect.Int64Kind,
	protoreflect.Sint64Kind,
	protoreflect.Uint64Kind,
	protoreflect.Sfixed32Kind,
	protoreflect.Fixed32Kind,
	protoreflect.Sfixed64Kind,
	protoreflect.Fixed64Kind,
	protoreflect.StringKind,
}

func IsIdKind(kind protoreflect.Kind) bool {
	return slices.Contains(IdKinds, kind)
}

func IsIdField(field *protogen.Field) bool {
	return IsIdKind(field.Desc.Kind())
}

func IsMessage(field *protogen.Field) bool {
	return field.Desc.Kind() == protoreflect.MessageKind
}

func GetMapFieldValue(field *protogen.Field) *protogen.Message {
	for _, field := range field.Message.Fields {
		if field.GoName == "Value" {
			return field.Message
		}
	}
	return nil
}

func IsMessageValueMap(field *protogen.Field) bool {
	return field.Desc.IsMap() && field.Desc.MapValue().Kind() == protoreflect.MessageKind
}
