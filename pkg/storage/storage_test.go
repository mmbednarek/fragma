package storage

import (
	"testing"

	v1 "github.com/mmbednarek/fragma/api/fragma/core/v1"
	"github.com/mmbednarek/fragma/model"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func Test_makeKey(t *testing.T) {
	app := v1.Application{
		Name: "Name",
		Path: "Path",
	}

	obj := model.Object{
		Version: "v1",
		Name:    "App",
		Kind:    "Application",
		Labels:  map[string]string{},
		Spec:    &app,
	}

	protoObj, err := obj.ToProto()
	require.NoError(t, err)

	key := makeKey(&protoObj)
	require.Equal(t, key, []byte("type.googleapis.com/fragma.core.v1.Application/App"))
}

func TestStorage_WriteObject(t *testing.T) {
	store, err := NewStorage("/tmp/storetest")
	require.NoError(t, err)

	app := v1.Application{
		Name: "Name",
		Path: "Path",
	}

	obj := model.Object{
		Version: "v1",
		Name:    "App",
		Kind:    "Application",
		Labels:  map[string]string{"option": "true"},
		Spec:    &app,
	}

	err = store.WriteObject(&obj)
	require.NoError(t, err)

	dbObj, err := store.ReadObject("fragma.core.v1.Application", "App")
	require.NoError(t, err)

	require.Equal(t, obj.Version, dbObj.Version)
	require.Equal(t, obj.Name, dbObj.Name)
	require.Equal(t, obj.Kind, dbObj.Kind)
	require.Equal(t, obj.Labels, dbObj.Labels)
	require.True(t, proto.Equal(&app, dbObj.Spec))
}
