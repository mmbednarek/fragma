package protoutil

import (
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func getProtoValue[T any](value protoreflect.Value) T {
	return value.Interface().(T)
}

func extractValueByFields[T any](msg protoreflect.Message, field []string) T {
	if len(field) == 1 {
		var res T
		msg.Range(func(desc protoreflect.FieldDescriptor, value protoreflect.Value) bool {
			if string(desc.Name()) == field[0] {
				res = getProtoValue[T](value)
				return false
			}
			return true
		})
		return res
	}

	var res protoreflect.Message
	msg.Range(func(desc protoreflect.FieldDescriptor, value protoreflect.Value) bool {
		if string(desc.Name()) == field[0] {
			res = getProtoValue[protoreflect.Message](value)
			return false
		}
		return true
	})

	return extractValueByFields[T](res, field[1:])
}

func ExtractValueByFieldName[T any](msg proto.Message, field string) T {
	fields := strings.Split(field, ".")
	return extractValueByFields[T](msg.ProtoReflect(), fields)
}
