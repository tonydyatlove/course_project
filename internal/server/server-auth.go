// Package server package server
package server

import (
	"context"
	"fmt"

	"github.com/tonydyatlove/leomessi/internal/model"
	pb "github.com/tonydyatlove/leomessi/proto"
	"golang.org/x/crypto/bcrypt"
)

// Authentication log-in
func (s *Server) Authentification(ctx context.Context, request *pb.AuthentificationRequest) (*pb.AuthentificationResponse, error) {
	idGame := request.GetId()
	password := request.GetPassword()
	accessToken, refreshToken, err := s.se.Authentification(ctx, idGame, password)
	if err != nil {
		return nil, err
	}
	return &pb.AuthenticationResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Registration sign-up
func (s *Server) Registration(ctx context.Context, request *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
	hashPassword, err := hashingPassword(request.Password)
	if err != nil {
		return nil, fmt.Errorf("server: error while hashing password, %e", err)
	}
	request.Password = hashPassword
	p := model.Person{
		Teams:     request.Teams,
		xG:         request.xG,
		Score:      request.Score,
		MVP: request.MVP,
	}
	newID, err := s.se.CreateGame(ctx, &p)
	if err != nil {
		return nil, err
	}
	return &pb.RegistrationResponse{Id: newID}, nil
}


func (s *Server) RefreshMyTokens(ctx context.Context, refreshTokenString *pb.RefreshTokensRequest) (*pb.RefreshTokensResponse, error) { // refresh our tokens
	refreshTokenStr := refreshTokenString.GetRefreshToken()
	newRefreshToken, newAccessToken, err := s.se.RefreshTokens(ctx, refreshTokenStr)
	if err != nil {
		return nil, err
	}
	return &pb.RefreshTokensResponse{
		RefreshToken: newRefreshToken,
		AccessToken:  newAccessToken,
	}, nil
}


func (s *Server) Logout(ctx context.Context, request *pb.LogoutRequest) (*pb.Response, error) {
	err := s.se.Verify(request.AccessToken)
	if err != nil {
		return nil, err
	}
	idUser := request.GetId()
	err = s.se.UpdateGameAuth(ctx, idGame, "")
	if err != nil {
		return nil, err
	}
	return new(pb.Response), nil
}

// hashingPassword _
func hashingPassword(password string) (string, error) {
	bytesPassword := []byte(password)
	hashedBytesPassword, err := bcrypt.GenerateFromPassword(bytesPassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hashPassword := string(hashedBytesPassword)
	return hashPassword, nil
}