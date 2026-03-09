package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/url"
	"sort"
	"strings"
)

type TelegramUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	PhotoURL  string `json:"photo_url"`
}

func ValidateTelegramInitData(initData string, botToken string) (*TelegramUser, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, err
	}

	hash := values.Get("hash")
	if hash == "" {
		return nil, errors.New("missing hash")
	}

	dataCheckString := buildDataCheckString(values)
	secret := sha256.Sum256([]byte(botToken))
	computed := computeHMAC(dataCheckString, secret[:])
	if !hmac.Equal([]byte(computed), []byte(hash)) {
		return nil, errors.New("invalid init data signature")
	}

	userJSON := values.Get("user")
	if userJSON == "" {
		return nil, errors.New("missing user payload")
	}

	var user TelegramUser
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func buildDataCheckString(values url.Values) string {
	pairs := make([]string, 0, len(values))
	for key := range values {
		if key == "hash" {
			continue
		}
		pairs = append(pairs, key+"="+values.Get(key))
	}
	sort.Strings(pairs)
	return strings.Join(pairs, "\n")
}

func computeHMAC(data string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
