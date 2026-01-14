package service

import "github.com/golang-jwt/jwt/v5"

func issueToken(userID, uri, addr, key string) (string, error) {
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"iss": "http://" + addr,
		"uri": uri,
	})
	return jwt.SignedString([]byte(key))
}
