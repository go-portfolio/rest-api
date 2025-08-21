package jwt

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-portfolio/rest-api/internal/config"
	"log"
)

var jwtKey []byte

func init() {
	// Загружаем конфиг
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Инициализация jwtKey из секрета в конфиге
	jwtKey = []byte(cfg.Jwt.JwtSecretKey)
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(username string) (string, error) {
	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "go-portfolio",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, jwt.NewValidationError("could not parse claims", jwt.ValidationErrorClaimsInvalid)
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, jwt.NewValidationError("token is expired", jwt.ValidationErrorExpired)
	}

	return claims, nil
}
