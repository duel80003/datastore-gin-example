package common

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

var hmacSecret = os.Getenv("HMAC_SECRET")

func stripBearerPrefixFromTokenString(tok string) (string, error) {
	// Should be a bearer token
	if len(tok) > 5 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}

	return tok, nil
}

var AuthorizationHeaderExtractor = &request.PostExtractionFilter{
	request.HeaderExtractor{"Authorization"},
	stripBearerPrefixFromTokenString,
}

var TokenExtractor = &request.MultiExtractor{
	AuthorizationHeaderExtractor,
	request.ArgumentExtractor{"access_token"},
}

func CreateFakeToken() (string, error) {
	ttl := 24 * time.Hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"exp": time.Now().UTC().Add(ttl).Unix()})

	hmacSampleSecret := []byte(hmacSecret)
	tokenString, err := token.SignedString(hmacSampleSecret)
	return tokenString, err
}

func CreateToken(name string, account string) (string, error) {
	ttl := 24 * time.Hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userName": name,
		"account":  account,
		"exp":      time.Now().UTC().Add(ttl).Unix(),
	})

	hmacSampleSecret := []byte(hmacSecret)
	tokenString, err := token.SignedString(hmacSampleSecret)
	return tokenString, err
}

func ParseToken(jwtToken string) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		hmacSampleSecret := []byte(hmacSecret)
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["exp"])
	} else {
		fmt.Println(err)
	}
}

func AuthenticationToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := request.ParseFromRequest(c.Request, TokenExtractor, func(token *jwt.Token) (interface{}, error) {
			b := ([]byte(hmacSecret))
			return b, nil
		})
		if err != nil {
			LogError("validation token error: ", err)
		}
		if !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		return
	}
}
