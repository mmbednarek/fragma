package model

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ObjectDetail struct {
	Version           string
	SingularName      string
	PluralName        string
	FullName          string
	ProtoType         protoreflect.MessageType
	HighlightedFields []string
}
