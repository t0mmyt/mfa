package storage

import (
	"encoding/json"
)
import bolt "go.etcd.io/bbolt"

type Bbolt struct {
	db *bolt.DB
}

const BucketName = "default"

func Open(path string) (*Bbolt, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			_, err := tx.CreateBucket([]byte(BucketName))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return &Bbolt{db: db}, nil
}

func (b Bbolt) Close() error {
	return b.db.Close()
}

func (b Bbolt) Import(in []byte) error {
	var keys []Key
	err := json.Unmarshal(in, &keys)
	if err != nil {
		return err
	}
	for _, key := range keys {
		err := b.Set(key.Id, key.Url)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b Bbolt) Export() (string, error) {
	keys, err := b.getAllKeys()
	if err != nil {
		return "", err
	}
	j, err := json.Marshal(&keys)
	return string(j), nil
}

func (b Bbolt) Get(id string) *Key {
	var url []byte
	_ = b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		url = b.Get([]byte(id))
		return nil
	})
	if url == nil {
		return nil
	}
	return NewKey(id, string(url))
}

func (b Bbolt) Set(id, url string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		return b.Put([]byte(id), []byte(url))
	})
}

func (b Bbolt) List() ([]string, error) {
	keys, err := b.getAllKeys()
	if err != nil {
		return nil, err
	}
	keyIds := make([]string, len(keys))
	for i, k := range keys {
		keyIds[i] = k.Id
	}
	return keyIds, nil
}

func (b Bbolt) getAllKeys() ([]*Key, error) {
	var keys []*Key
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		return b.ForEach(func(k, v []byte) error {
			keys = append(keys, NewKey(string(k), string(v)))
			return nil
		})
	})
	return keys, err
}
