// Package server p
package server

import (
	"github.com/tonydyatlove/leomessi/internal/model"
	"github.com/tonydyatlove/leomessi/internal/service"
	pb "github.com/tonydyatlove/leomessi/proto"

	"context"
)

// Server struct
type Server struct {
	pb.UnimplementedCRUDServer
	se *service.Service
}

// NewServer create new server connection
func NewServer(serv *service.Service) *Server {
	return &Server{se: serv}
}


func (s *Server) GetGame(ctx context.Context, request *pb.GetGameRequest) (*pb.GetGameResponse, error) {
	accessToken := request.GetAccessToken()
	if err := s.se.Verify(accessToken); err != nil {
		return nil, err
	}
	idPerson := request.GetId()
	personDB, err := s.se.GetGame(ctx, idPerson)
	if err != nil {
		return nil, err
	}
	personProto := &pb.GetGameResponse{
		Person: &pb.Person{
			Id:       personDB.ID,
			Teams:     personDB.Teams,
			xG:      personDB.xG,
			Score:    personDB.Score,
			MVP: personDB.MVP,
		},
	}
	return personProto, nil
}

// GetAllUsers get all users from db
func (s *Server) GetAllGames(ctx context.Context, _ *pb.GetAllGamesRequest) (*pb.GetAllGamesResponse, error) {
	persons, err := s.se.GetAllGames(ctx)
	if err != nil {
		return nil, err
	}
	var list []*pb.Person
	for _, person := range persons {
		personProto := new(pb.Person)
		personProto.Id = person.ID
		personProto.Teams = person.Teams
		personProto.Score = person.Score
		personProto.xG = person.xG
		list = append(list, personProto)
	}
	return &pb.GetAllGamesResponse{Persons: list}, nil
}

// DeleteUser delete user by id
func (s *Server) DeleteGame(ctx context.Context, request *pb.DeleteGameRequest) (*pb.Response, error) {
	idUser := request.GetId()
	err := s.se.DeleteGame(ctx, idGame)
	if err != nil {
		return nil, err
	}
	return new(pb.Response), nil
}

// UpdateUser update user with new parameters
func (s *Server) UpdateGame(ctx context.Context, request *pb.UpdateGameRequest) (*pb.Response, error) {
	accessToken := request.GetAccessToken()
	if err := s.se.Verify(accessToken); err != nil {
		return nil, err
	}
	user := &model.Person{
		Teams:  request.Person.Teams,
		xG: request.Person.xG,
		Score:   request.Person.Score,
	}
	idGame := request.GetId()
	err := s.se.UpdateGame(ctx, idGame, game)
	if err != nil {
		return nil, err
	}
	return new(pb.Response), nil
}