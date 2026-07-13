package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TelegramUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	PhotoURL  string `json:"photo_url"`
	// StartParam is the startapp deep-link payload, taken from the signed
	// initData so referral attribution cannot be spoofed by the client.
	StartParam string `json:"-"`
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

	authDateStr := values.Get("auth_date")
	if authDateStr == "" {
		return nil, errors.New("missing auth_date")
	}
	authDateUnix, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		return nil, errors.New("invalid auth_date")
	}
	const maxAuthAge = 30 * time.Minute
	authTime := time.Unix(authDateUnix, 0)
	if time.Since(authTime) > maxAuthAge || authTime.After(time.Now().Add(1*time.Minute)) {
		return nil, errors.New("init data expired")
	}

	dataCheckString := buildDataCheckString(values)
	secret := computeTelegramWebAppSecret(botToken)
	computed := computeHMAC(dataCheckString, secret)
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
	user.StartParam = values.Get("start_param")

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

func computeTelegramWebAppSecret(botToken string) []byte {
	h := hmac.New(sha256.New, []byte("WebAppData"))
	h.Write([]byte(botToken))
	return h.Sum(nil)
}
