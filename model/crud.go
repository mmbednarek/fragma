package model

import (
	"fmt"
)

type Storage interface {
	WriteObject(obj *Object) error
	ReadObject(typeName string, name string) (Object, error)
	ReadAllObjects(typeName string) ([]Object, error)
	RemoveObject(typeName string, name string) error
}

type Controller interface {
	OnDelete(typeName string, name string)
	OnUpdate(obj *Object)
	OnRead(obj *Object)
}

type CrudService[TStore Storage] struct {
	controllers []Controller
	storage     TStore
}

func NewCrudService[TStore Storage](storage TStore) CrudService[TStore] {
	return CrudService[TStore]{storage: storage}
}

func (s *CrudService[TStore]) AddController(listener Controller) {
	s.controllers = append(s.controllers, listener)
}

func (s *CrudService[TStore]) Update(obj *Object) error {
	if err := s.storage.WriteObject(obj); err != nil {
		return fmt.Errorf("s.storage.WriteObject: %w", err)
	}

	for _, listener := range s.controllers {
		listener.OnUpdate(obj)
	}
	return nil
}

func (s *CrudService[TStore]) Read(typeName string, name string) (Object, error) {
	obj, err := s.storage.ReadObject(typeName, name)
	if err != nil {
		return Object{}, fmt.Errorf("s.storage.ReadObject: %w", err)
	}

	for _, listener := range s.controllers {
		listener.OnRead(&obj)
	}

	return obj, nil
}

func (s *CrudService[TStore]) Delete(typeName string, name string) error {
	if err := s.storage.RemoveObject(typeName, name); err != nil {
		return fmt.Errorf("s.storage.RemoveObject: %w", err)
	}

	for _, cont := range s.controllers {
		cont.OnDelete(typeName, name)
	}

	return nil
}

func (s *CrudService[TStore]) ReadAll(typeName string) ([]Object, error) {
	obj, err := s.storage.ReadAllObjects(typeName)
	if err != nil {
		return nil, fmt.Errorf("s.storage.ReadAllObjects: %w", err)
	}

	for _, cont := range s.controllers {
		for i, ob := range obj {
			cont.OnRead(&ob)
			obj[i] = ob
		}
	}

	return obj, nil
}
