package util

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func GetBoolExtension(desc protoreflect.Descriptor, ext protoreflect.ExtensionType) bool {
	opts := desc.Options()
	if proto.HasExtension(opts, ext) {
		value, ok := proto.GetExtension(opts, ext).(bool)
		return ok && value
	}
	return false
}

func GetStringExtension(desc protoreflect.Descriptor, ext protoreflect.ExtensionType) (bool, string) {
	opts := desc.Options()
	if proto.HasExtension(opts, ext) {
		value, ok := proto.GetExtension(opts, ext).(string)
		return ok, value
	}
	return false, ""
}
