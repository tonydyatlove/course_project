// Package repository a
package repository

import (
	"github.com/tonydyatlove/leomessi/internal/model"

	"context"
)

// Repository transition to mongo or postgres db
type Repository interface {
	CreateGame(ctx context.Context, p *model.Person) (string, error)
	GetGameByID(ctx context.Context, idPerson string) (*model.Person, error)
	GetAllGames(ctx context.Context) ([]*model.Person, error)
	DeleteGame(ctx context.Context, id string) error
	UpdateGame(ctx context.Context, id string, per *model.Person) error
}