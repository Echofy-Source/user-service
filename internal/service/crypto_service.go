package service

import (
	"fmt"
	"time"

	"github.com/Echofy-Source/user-service/internal/model"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type CryptoService interface {
	GenerateTokens(username string) (*model.TokenModel, error)
	RenewTokens(refreshTokenStr string) (*model.TokenModel, error)
	HashPassword(password string) (string, error)
}

type CryptoServiceImpl struct {
	secretKey string
}

// GenerateTokens implements CryptoService.
func (c *CryptoServiceImpl) GenerateTokens(username string) (*model.TokenModel, error) {
	accessTokenExpTime := time.Now().Add(time.Minute * 15).Unix()
	refreshTokenExpTime := time.Now().Add(time.Hour * 24).Unix()

	// Add username to the accessToken
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"type":     "access",
		"exp":      accessTokenExpTime,
	})

	// Add username to the refreshToken
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"type":     "refresh",
		"exp":      refreshTokenExpTime,
	})

	// Sign the tokens
	accessTokenStr, err := accessToken.SignedString([]byte(c.secretKey))
	if err != nil {
		return nil, err
	}

	refreshTokenStr, err := refreshToken.SignedString([]byte(c.secretKey))
	if err != nil {
		return nil, err
	}

	return &model.TokenModel{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}, nil
}

// RenewTokens implements CryptoService.
func (c *CryptoServiceImpl) RenewTokens(refreshTokenStr string) (*model.TokenModel, error) {
	// Check if the refreshToken is valid
	refreshToken, err := jwt.Parse(refreshTokenStr, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is valid
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(c.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok || !refreshToken.Valid {
		return nil, err
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, fmt.Errorf("invalid token type")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username")
	}

	// Generate new tokens
	return c.GenerateTokens(username)
}

// HashPassword implements CryptoService.
func (c *CryptoServiceImpl) HashPassword(password string) (string, error) {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bcryptPassword), nil
}

// Ensure CryptoServiceImpl implements CryptoService.
var _ CryptoService = &CryptoServiceImpl{}
