package service

import "github.com/golang-jwt/jwt/v5"

func issueToken(userID, uri, key string) (string, error) {
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"uri": uri,
	})
	return jwt.SignedString([]byte(key))
}
