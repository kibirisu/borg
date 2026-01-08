package service

import (
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

func issueToken(userID int32, username, addr, key string) (string, error) {
	id := strconv.Itoa(int(userID))
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  id,
		"iss":  "http://" + addr,
		"name": username,
	})
	return jwt.SignedString([]byte(key))
}
