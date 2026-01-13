package service

import "github.com/golang-jwt/jwt/v5"

func issueToken(userID, username, addr, key string) (string, error) {
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"iss":  "http://" + addr,
		"name": username,
	})
	return jwt.SignedString([]byte(key))
}
