package model

import (
	"encoding/json"
	"fmt"
	"testing"

	core_v1 "github.com/mmbednarek/fragma/api/fragma/core/v1"
	"github.com/stretchr/testify/require"
)

func Test_ObjectMarshall(t *testing.T) {
	app := core_v1.Application{
		Name: "TestApp",
		Path: "/usr/bin/test",
	}

	obj := Object{
		Version: "core_v1",
		Name:    "test",
		Kind:    "fragma.core_v1.core.Application",
		Labels:  map[string]string{"label": "something"},
		Spec:    Spec{&app},
	}

	bytes, err := json.Marshal(obj)
	require.NoError(t, err)

	fmt.Println(string(bytes))

	obj2 := Object{
		Spec: Spec{&core_v1.Application{}},
	}
	require.NoError(t, json.Unmarshal(bytes, &obj2))

	fmt.Println(obj2)
}
