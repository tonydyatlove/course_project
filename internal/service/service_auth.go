// Package service ...
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/tonydyatlove/leomessi/internal/model"
	"github.com/tonydyatlove/leomessi/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)


var (
	AccessTokenWorkTime  = time.Now().Add(time.Minute * 5).Unix()
	RefreshTokenWorkTime = time.Now().Add(time.Hour * 3).Unix()
)


func (se *Service) Authentication(ctx context.Context, id, password string) (accessToken, refreshToken string, err error) {
	authUser, err := se.rps.GetGameByID(ctx, id)
	if err != nil {
		return "", "", err
	}
	incoming := []byte(password)
	existing := []byte(authGame.Password)
	err = bcrypt.CompareHashAndPassword(existing, incoming) // check passwords
	if err != nil {
		return "", "", err
	}
	authGame.Password = password
	accessToken, refreshToken, err = se.CreateJWT(ctx, se.rps, authGame)
	if err != nil {
		return "", "", err
	}
	return
}

// CreateUser create new user, add him to db
func (se *Service) CreateGame(ctx context.Context, p *model.Person) (string, error) {
	hashedPassword, err := hashingPassword(p.Password)
	if err != nil {
		return "", err
	}
	p.Password = hashedPassword
	return se.rps.CreateGame(ctx, p)
}

// RefreshTokens refresh tokens
func (se *Service) RefreshTokens(ctx context.Context, refreshTokenStr string) (newRefreshToken, newAccessToken string, err error) { // refresh our tokens
	refreshToken, err := jwt.Parse(refreshTokenStr, func(t *jwt.Token) (interface{}, error) {
		return se.jwtKey, nil
	}) // parse it into string format
	if err != nil {
		log.Errorf("service: can't parse refresh token - %e", err)
		return "", "", err
	}
	if !refreshToken.Valid {
		return "", "", fmt.Errorf("service: expired refresh token")
	}
	claims := refreshToken.Claims.(jwt.MapClaims)
	gameUUID := claims["jti"]
	if gameUUID == "" {
		return "", "", fmt.Errorf("service: error while parsing claims, ID couldnt be empty")
	}
	person, err := se.rps.SelectByIDAuth(ctx, gameUUID.(string))
	if err != nil {
		return "", "", fmt.Errorf("service: token refresh failed - %e", err)
	}
	if refreshTokenStr != person.RefreshToken {
		return "", "", fmt.Errorf("service: invalid refresh token")
	}

	return se.CreateJWT(ctx, se.rps, &person)
}

// CreateJWT create jwt tokens
func (se *Service) CreateJWT(ctx context.Context, rps repository.Repository, person *model.Person) (accessTokenStr, refreshTokenStr string, err error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)            // encrypt access token by SigningMethodHS256 method
	claimsA := accessToken.Claims.(jwt.MapClaims)             // fill access-token`s claims
	claimsA["exp"] = AccessTokenWorkTime                      // work time
	claimsA["teamsname"] = person.Teams                         // payload
	accessTokenStr, err = accessToken.SignedString(se.jwtKey) // convert token to string format
	if err != nil {
		log.Errorf("service: can't generate access token - %v", err)
		return "", "", err
	}
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	claimsR := refreshToken.Claims.(jwt.MapClaims)
	claimsR["teamsname"] = person.Teams
	claimsR["exp"] = RefreshTokenWorkTime
	claimsR["jti"] = person.ID
	refreshTokenStr, err = refreshToken.SignedString(se.jwtKey)
	if err != nil {
		log.Errorf("service: can't generate access token - %v", err)
		return "", "", err
	}
	err = rps.UpdateAuth(ctx, person.ID, refreshTokenStr)
	if err != nil {
		log.Errorf("service: can't generate access token - %v", err)
		return "", "", err
	}
	return
}

// UpdateUserAuth update auth user, add token
func (se *Service) UpdateGameAuth(ctx context.Context, id, refreshToken string) error {
	return se.rps.UpdateAuth(ctx, id, refreshToken)
}

// Verify verify access token
func (se *Service) Verify(accessTokenString string) error {
	accessToken, err := jwt.Parse(accessTokenString, func(t *jwt.Token) (interface{}, error) {
		return se.jwtKey, nil
	})
	if err != nil {
		log.Errorf("service: can't parse refresh token - ", err)
		return err
	}
	if !accessToken.Valid {
		return fmt.Errorf("service: expired refresh token")
	}
	return nil
}

// hashingPassword _
func hashingPassword(password string) (string, error) {
	if len(password) < 5 || len(password) > 30 {
		return "", fmt.Errorf("password is too short or too long")
	}
	bytesPassword := []byte(password)
	hashedBytesPassword, err := bcrypt.GenerateFromPassword(bytesPassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hashPassword := string(hashedBytesPassword)
	return hashPassword, nil
}