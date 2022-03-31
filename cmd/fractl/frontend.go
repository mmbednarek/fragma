package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mmbednarek/fragma/model"
	"github.com/mmbednarek/fragma/model/repo"
	"github.com/mmbednarek/fragma/pkg/protoutil"
	"github.com/mmbednarek/fragma/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

type Client interface {
	GetAll(api string, typeName string) ([]model.Object, error)
	GetObject(api string, typeName string, name string) (model.Object, error)
	WriteObject(obj model.Object) error
	DeleteObject(apiName string, typeName string, name string) error
}

type Frontend struct {
	Client Client
	Apis   repo.ApiRepository
}

func NewFrontend(client Client, apis repo.ApiRepository) *Frontend {
	return &Frontend{Client: client, Apis: apis}
}

func (f *Frontend) Mount() *cobra.Command {
	root := &cobra.Command{
		Use: "fractl",
	}

	getCmd := &cobra.Command{
		Use:  "get",
		Args: cobra.RangeArgs(1, 2),
		Run:  f.HandleGet,
	}
	root.AddCommand(getCmd)

	applyCmd := &cobra.Command{
		Use:  "apply",
		Args: cobra.NoArgs,
		Run:  f.HandleApply,
	}
	root.AddCommand(applyCmd)

	deleteCmd := &cobra.Command{
		Use:  "delete",
		Args: cobra.ExactArgs(2),
		Run:  f.HandleDelete,
	}
	root.AddCommand(deleteCmd)

	return root
}

func getApiAndTypeName(objectType string) (string, string) {
	if strings.Contains(objectType, ".") {
		index := strings.LastIndexByte(objectType, '.')
		return objectType[:index], strings.ToLower(objectType[index+1:])
	}

	return "fragma.core.v1", objectType
}

func (f *Frontend) HandleGet(cmd *cobra.Command, args []string) {
	api, typeName := getApiAndTypeName(args[0])

	det, ok := f.Apis.FindApi(api)
	if !ok {
		die("api not found: %s", api)
	}
	typeDep, err := det.GetObjectDetail(typeName)
	if err != nil {
		die("det.GetObjectDetail: %s", err)
	}

	if len(args) == 2 {
		objectName := args[1]
		obj, err := f.Client.GetObject(api, typeDep.SingularName, objectName)
		if err != nil {
			die("could not get objects: %s", err)
		}

		yamlMarshalled, err := yaml.Marshal(obj)
		if err != nil {
			die("could not marshall object: %s", err)
		}

		_, err = os.Stdout.Write(yamlMarshalled)
		if err != nil {
			die("could not write object: %s", err)
		}
		return
	}

	objs, err := f.Client.GetAll(api, typeDep.SingularName)
	if err != nil {
		die("could not get objects: %s", err)
	}

	table := util.NewTable()

	for _, obj := range objs {
		for _, field := range typeDep.HighlightedFields {
			table.Add(field, protoutil.ExtractValueByFieldName[string](obj.Spec, field))
		}
	}

	table.Print(os.Stdout)
}

func (f *Frontend) HandleApply(cmd *cobra.Command, args []string) {
	var meta model.JustMeta

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		die("could not read from stdin")
	}
	if err := yaml.Unmarshal(data, &meta); err != nil {
		die("yaml decode: %s", err)
	}

	api, typeName := getApiAndTypeName(meta.Kind)

	det, ok := f.Apis.FindApi(api)
	if !ok {
		die("api not found: %s", api)
	}
	typeDep, err := det.GetObjectDetail(typeName)
	if err != nil {
		die("det.GetObjectDetail: %s", err)
	}

	obj := model.Object{
		Spec: model.Spec{Message: typeDep.ProtoType.New().Interface()},
	}
	if err := yaml.Unmarshal(data, &obj); err != nil {
		die("yaml decode: %s", err)
	}

	if err := f.Client.WriteObject(obj); err != nil {
		die("error writing object: %s", err)
	}
}

func (f *Frontend) HandleDelete(cmd *cobra.Command, args []string) {
	api, typeName := getApiAndTypeName(args[0])

	if err := f.Client.DeleteObject(api, typeName, args[1]); err != nil {
		die("f.Client.WriteObject: %s", err)
	}
}

func die(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

type FlagErrChain struct {
	Flags    *pflag.FlagSet
	ErrorMap map[string]error
}

func NewFlagErrChain(flags *pflag.FlagSet) FlagErrChain {
	return FlagErrChain{
		Flags:    flags,
		ErrorMap: map[string]error{},
	}
}

func (c *FlagErrChain) GetInt(name string) *int {
	if !c.Flags.Changed(name) {
		return nil
	}
	value, err := c.Flags.GetInt(name)
	if err != nil {
		c.ErrorMap[name] = err
		return nil
	}
	return &value
}

func (c *FlagErrChain) GetString(name string) *string {
	if !c.Flags.Changed(name) {
		return nil
	}
	value, err := c.Flags.GetString(name)
	if err != nil {
		c.ErrorMap[name] = err
		return nil
	}
	return &value
}

func (c *FlagErrChain) GetStringSlice(name string) []string {
	value, err := c.Flags.GetStringSlice(name)
	if err != nil {
		c.ErrorMap[name] = err
		return nil
	}
	return value
}

func (c *FlagErrChain) GetBool(name string) bool {
	value, err := c.Flags.GetBool(name)
	if err != nil {
		c.ErrorMap[name] = err
		return false
	}
	return value
}

func (c *FlagErrChain) GetDuration(name string) *time.Duration {
	if !c.Flags.Changed(name) {
		return nil
	}

	value, err := c.Flags.GetString(name)
	if err != nil {
		c.ErrorMap[name] = err
		return nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		c.ErrorMap[name] = err
		return nil
	}

	return &parsed
}

func (c *FlagErrChain) Verify() error {
	if len(c.ErrorMap) == 0 {
		return nil
	}
	errText := "invalid flags:"
	for flag, err := range c.ErrorMap {
		errText += fmt.Sprintf(" %s (%v)", flag, err)
	}
	return errors.New(errText)
}
