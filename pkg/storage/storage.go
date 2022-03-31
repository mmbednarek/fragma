package storage

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
	core "github.com/mmbednarek/fragma/api/fragma/core/v1"
	"github.com/mmbednarek/fragma/model"
	"google.golang.org/protobuf/proto"
)

type Storage struct {
	db *badger.DB
}

func NewStorage(path string) (Storage, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return Storage{}, fmt.Errorf("badger.Open: %w", err)
	}

	return Storage{
		db: db,
	}, nil
}

func (s Storage) WriteObject(obj *model.Object) error {
	protoObj, err := obj.ToProto()
	if err != nil {
		return fmt.Errorf("obj.ToProto: %w", err)
	}

	bytes, err := proto.Marshal(&protoObj)
	if err != nil {
		return fmt.Errorf("proto.Marshal: %w", err)
	}

	trans := s.db.NewTransaction(true)
	defer trans.Discard()

	key := makeKey(&protoObj)
	if key == nil {
		return fmt.Errorf("invalid metadata")
	}

	if err := trans.Set(key, bytes); err != nil {
		return fmt.Errorf("trans.Set: %w", err)
	}

	if err := trans.Commit(); err != nil {
		return fmt.Errorf("trans.Commit: %w", err)
	}
	return nil
}

func (s Storage) ReadObject(typeName string, name string) (model.Object, error) {
	var result []byte

	err := s.db.View(func(txn *badger.Txn) error {
		value, err := txn.Get(makeKeyWithTypeUrl("type.googleapis.com/"+typeName, name))
		if err != nil {
			return fmt.Errorf("trans.Get: %w", err)
		}

		err = value.Value(func(val []byte) error {
			result = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return fmt.Errorf("value.Value: %w", err)
		}

		return nil
	})
	if err != nil {
		return model.Object{}, fmt.Errorf("s.db.View: %w", err)
	}

	protoObj := core.Object{}
	if err := proto.Unmarshal(result, &protoObj); err != nil {
		return model.Object{}, fmt.Errorf("proto.Unmarshal: %w", err)
	}

	obj, err := model.ObjectFromProto(&protoObj)
	if err != nil {
		return model.Object{}, fmt.Errorf("model.ObjectFromProto: %w", err)
	}

	return obj, nil
}

func (s Storage) ReadAllObjects(typeName string) ([]model.Object, error) {
	var result []model.Object
	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			Prefix: []byte("type.googleapis.com/" + typeName + "/"),
		})
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				protoObj := core.Object{}
				if err := proto.Unmarshal(val, &protoObj); err != nil {
					return fmt.Errorf("proto.Unmarshal: %w", err)
				}

				obj, err := model.ObjectFromProto(&protoObj)
				if err != nil {
					return fmt.Errorf("model.ObjectFromProto: %w", err)
				}

				result = append(result, obj)
				return nil
			})
			if err != nil {
				return fmt.Errorf("item.Value: %s", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("s.db.View: %s", err)
	}

	return result, nil
}

func (s Storage) RemoveObject(typeName string, name string) error {
	tx := s.db.NewTransaction(true)
	defer tx.Discard()

	if err := tx.Delete(makeKeyWithTypeUrl("type.googleapis.com/"+typeName, name)); err != nil {
		return fmt.Errorf("tx.Delete: %s", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("tx.Commit: %s", err)
	}

	return nil
}

func makeKey(obj *core.Object) []byte {
	if obj.Metadata == nil {
		return nil
	}
	return makeKeyWithTypeUrl(obj.Spec.TypeUrl, obj.Metadata.Name)
}

func makeKeyWithTypeUrl(typeUrl string, name string) []byte {
	key := make([]byte, len(typeUrl)+1+len(name))
	copy(key, []byte(typeUrl))
	key[len(typeUrl)] = '/'
	copy(key[len(typeUrl)+1:], []byte(name))
	return key
}
