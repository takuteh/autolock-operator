package controller

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

func configHash(config interface{}) (string, error) {
	data, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)

	return fmt.Sprintf("%x", hash), nil
}
