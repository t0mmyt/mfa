package storage

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"time"
)

type Key struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

func NewKey(id, url string) *Key {
	return &Key{
		Id:  id,
		Url: url,
	}
}

func (k *Key) Totp() (string, error) {
	key, err := otp.NewKeyFromURL(k.Url)
	if err != nil {
		return "", err
	}
	code, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		return "", err
	}
	return code, nil
}
