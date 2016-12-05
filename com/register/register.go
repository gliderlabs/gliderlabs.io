package register

import (
	"bytes"
	"fmt"

	"github.com/gliderlabs/gliderlabs.io/com/auth0"
	"github.com/lytics/base62"
	"github.com/mitchellh/mapstructure"
	"github.com/satori/go.uuid"
)

type User struct {
	Name        string
	Nickname    string
	Email       string
	AppMetadata AppMetadata `mapstructure:"app_metadata"`
	ID          string      `mapstructure:"user_id"`
}

type AppMetadata struct {
	Registered map[string]RegisteredProject
}

type RegisteredProject struct {
	Key    string
	Charge string
}

func LookupUser(uid string) (User, error) {
	data, err := auth0.Client().User(uid)
	if err != nil {
		return User{}, err
	}
	var user User
	err = mapstructure.Decode(data, &user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func GenerateKey(lookupKey string) string {
	id := uuid.NewV4()
	rawKey := append([]byte(fmt.Sprintf("1:%s:", lookupKey)), id.Bytes()...)
	return base62.StdEncoding.EncodeToString(rawKey)
}

func ExtractLookupKey(key string) string {
	data, err := base62.StdEncoding.DecodeString(key)
	if err != nil {
		return ""
	}
	parts := bytes.Split(data, []byte(":"))
	if len(parts) < 3 {
		return ""
	}
	return string(parts[1])
}
