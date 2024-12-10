package utils

import (
	"errors"
	constants "go-server/internal/const"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type CustomClaim struct {
	Id string `json:"id"`
	jwt.RegisteredClaims
}

func LoadEnv() {
	if err := godotenv.Load("config.env"); err != nil {
		log.Fatal("Error Loading Env Files")
	}
}

func CreateJWT(userId string) (string, error) {
	claimsData := CustomClaim{
		Id: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 10)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Test",
		},
	}
	getSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsData)
	return token.SignedString([]byte(getSecret))
}

func ParseClaimFromJWT(r *http.Request) (*jwt.Token, *CustomClaim) {
	bearerText := r.Header.Get("authorization")
	tokenArr := strings.Split(bearerText, " ")

	var token string = ""
	if len(tokenArr) > 1 {
		token = tokenArr[1]
	} else {
		return &jwt.Token{Valid: false}, nil
	}
	secret := os.Getenv("JWT_SECRET")

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}
	claims := &CustomClaim{}
	jwtToken, _ := jwt.ParseWithClaims(token, claims, keyFunc)
	return jwtToken, claims
}

func VerifyJWT(r *http.Request) (*CustomClaim, error) {
	x, claim := ParseClaimFromJWT(r)
	if x.Valid {
		expTime, _ := x.Claims.GetExpirationTime()
		if expTime.Unix() < time.Now().Unix() {
			return nil, errors.New("jwt expired")
		} else {
			return claim, nil
		}
	} else {
		return nil, errors.New("invalid jwt token")
	}

}

func GetClaimFromCtx(r *http.Request) CustomClaim {
	val := r.Context().Value(constants.Claim).(*CustomClaim)
	return *val
}
