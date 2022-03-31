package model

import (
	"encoding/json"
	"errors"
	"fmt"

	core "github.com/mmbednarek/fragma/api/fragma/core/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	ErrObjectNotFound = errors.New("object not found")
)

type Spec struct {
	proto.Message
}

func (s Spec) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(s)
}

func (s Spec) UnmarshalJSON(data []byte) error {
	return protojson.Unmarshal(data, s)
}

func (s Spec) MarshalYAML() (interface{}, error) {
	jsonData, err := protojson.Marshal(s)
	if err != nil {
		return nil, err
	}

	dynObj := map[string]any{}
	if err := json.Unmarshal(jsonData, &dynObj); err != nil {
		return nil, err
	}

	return dynObj, nil
}

func (s Spec) UnmarshalYAML(unmarshal func(interface{}) error) error {
	dynObj := map[string]any{}
	if err := unmarshal(&dynObj); err != nil {
		return err
	}
	jsonData, err := json.Marshal(dynObj)
	if err != nil {
		return err
	}

	return protojson.Unmarshal(jsonData, &s)
}

type Metadata struct {
	Name        string            `json:"name" yaml:"name"`
	Labels      map[string]string `json:"labels" yaml:"labels"`
	Annotations map[string]string `json:"annotations" yaml:"annotations"`
}

type Object struct {
	Kind     string   `json:"kind" yaml:"kind"`
	Metadata Metadata `json:"metadata" yaml:"metadata"`
	Spec     Spec     `json:"spec" yaml:"spec"`
}

type JustMeta struct {
	Kind     string   `json:"kind" yaml:"kind"`
	Metadata Metadata `json:"metadata" yaml:"metadata"`
}

func (o Object) ToProto() (core.Object, error) {
	spec, err := anypb.New(o.Spec)
	if err != nil {
		return core.Object{}, fmt.Errorf("anypb.New: %w", err)
	}

	return core.Object{
		Kind: o.Kind,
		Metadata: &core.Metadata{
			Name:        o.Metadata.Name,
			Labels:      o.Metadata.Labels,
			Annotations: o.Metadata.Annotations,
		},
		Spec: spec,
	}, nil
}

func ObjectFromProto(object *core.Object) (Object, error) {
	message, err := object.Spec.UnmarshalNew()
	if err != nil {
		return Object{}, fmt.Errorf("object.Spec.UnmarshalNew: %w", err)
	}

	meta := Metadata{}
	if object.Metadata != nil {
		meta.Name = object.Metadata.Name
		meta.Labels = object.Metadata.Labels
		meta.Annotations = object.Metadata.Annotations
	}

	return Object{
		Kind:     object.Kind,
		Metadata: meta,
		Spec:     Spec{message},
	}, nil
}
