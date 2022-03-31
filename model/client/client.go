package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mmbednarek/fragma/model"
	"github.com/mmbednarek/fragma/model/repo"
	"github.com/valyala/fasthttp"
)

type Client struct {
	host     string
	insecure bool
}

func NewClient(host string) *Client {
	return &Client{
		host:     host,
		insecure: true,
	}
}

func (c *Client) protocolPrefix() string {
	if c.insecure {
		return "http"
	}
	return "https"
}

func getObjDetailByName(api string, name string) (model.ObjectDetail, error) {
	repository := repo.GetStandardRepository()
	apiDet, ok := repository.FindApi(api)
	if !ok {
		return model.ObjectDetail{}, errors.New("invalid api name")
	}

	object, err := apiDet.GetObjectDetail(name)
	if err != nil {
		return model.ObjectDetail{}, errors.New("invalid object name")
	}
	return object, nil
}

type JsonSliceIterator struct {
	Data []byte
	At   int
}

func (j *JsonSliceIterator) Next() []byte {
	counter := 0
	for i, b := range j.Data[j.At:] {
		switch b {
		case '{':
			counter++
		case '}':
			counter--
			if counter == 0 {
				slice := j.Data[j.At:(j.At + i + 1)]
				j.At += i + 1
				return slice
			}
		}
	}
	return nil
}

func ReadJsonSlice(data []byte, start int) ([]byte, int) {
	counter := 0
	for i, b := range data[start:] {
		switch b {
		case '{':
			counter++
		case '}':
			counter--
			if counter == 0 {
				return data[:i+1], start + i + 1
			}
		}
	}
	return data, len(data)
}

func splitApiAndTypeName(full string) (string, string) {
	idx := strings.LastIndexByte(full, '.')
	return full[:idx], strings.ToLower(full[idx+1:])
}

func (c *Client) GetAll(api string, typeName string) ([]model.Object, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(fmt.Sprintf("%s://%s/apis/%s/%s", c.protocolPrefix(), c.host, api, typeName))
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		return nil, fmt.Errorf("fasthttp.Do: %w", err)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode())
	}

	objDetail, err := getObjDetailByName(api, typeName)
	if err != nil {
		return nil, fmt.Errorf("getObjDetailByName: %w", err)
	}

	body := resp.Body()
	var result []model.Object

	jsonSlices := JsonSliceIterator{
		Data: body,
		At:   0,
	}
	for {
		slice := jsonSlices.Next()
		if slice == nil {
			break
		}

		obj := model.Object{
			Spec: model.Spec{Message: objDetail.ProtoType.New().Interface()},
		}
		if err := json.Unmarshal(slice, &obj); err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w", err)
		}

		result = append(result, obj)
	}

	return result, nil
}

func (c *Client) GetObject(api string, typeName string, name string) (model.Object, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	objDetail, err := getObjDetailByName(api, typeName)
	if err != nil {
		return model.Object{}, fmt.Errorf("getObjDetailByName: %w", err)
	}

	req.SetRequestURI(fmt.Sprintf("%s://%s/apis/%s/%s/%s", c.protocolPrefix(), c.host, api, objDetail.PluralName, name))
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		return model.Object{}, fmt.Errorf("fasthttp.Do: %w", err)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		return model.Object{}, fmt.Errorf("invalid status code: %d", resp.StatusCode())
	}

	obj := model.Object{
		Spec: model.Spec{Message: objDetail.ProtoType.New().Interface()},
	}
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return model.Object{}, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return obj, nil
}

func (c *Client) WriteObject(obj model.Object) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	apiName, typeName := splitApiAndTypeName(obj.Kind)
	objDetail, err := getObjDetailByName(apiName, typeName)
	if err != nil {
		return fmt.Errorf("getObjDetailByName: %w", err)
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	req.Header.SetMethod(fasthttp.MethodPut)
	req.SetBody(data)
	req.SetRequestURI(fmt.Sprintf("%s://%s/apis/%s/%s", c.protocolPrefix(), c.host, apiName, objDetail.SingularName))
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		return fmt.Errorf("fasthttp.Do: %w", err)
	}
	if resp.StatusCode() != fasthttp.StatusCreated {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode())
	}

	return nil
}

func (c *Client) DeleteObject(apiName string, typeName string, name string) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	objDetail, err := getObjDetailByName(apiName, typeName)
	if err != nil {
		return fmt.Errorf("getObjDetailByName: %w", err)
	}

	req.Header.SetMethod(fasthttp.MethodDelete)
	req.SetRequestURI(fmt.Sprintf("%s://%s/apis/%s/%s/%s", c.protocolPrefix(), c.host, apiName, objDetail.PluralName, name))
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		return fmt.Errorf("fasthttp.Do: %w", err)
	}
	if resp.StatusCode() != fasthttp.StatusNoContent {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode())
	}

	return nil
}
