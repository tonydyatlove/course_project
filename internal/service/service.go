// Package service i
package service

import (
	"github.com/tonydyatlove/leomessi/internal/model"
	"github.com/tonydyatlove/leomessi/internal/repository"

	"context"
)

// Service s
type Service struct {
	jwtKey []byte
	rps    repository.Repository
}

// NewService create new service connection
func NewService(pool repository.Repository, jwtKey []byte) *Service {
	return &Service{rps: pool, jwtKey: jwtKey}
}


func (se *Service) GetGame(ctx context.Context, id string) (*model.Person, error) {
	return se.rps.GetGameByID(ctx, id)
}


func (se *Service) GetAllGames(ctx context.Context) ([]*model.Person, error) {
	return se.rps.GetAllGames(ctx)
}


func (se *Service) DeleteGame(ctx context.Context, id string) error {
	return se.rps.DeleteGame(ctx, id)
}


func (se *Service) UpdateGame(ctx context.Context, id string, game *model.Person) error {
	return se.rps.UpdateGame(ctx, id, game)
}