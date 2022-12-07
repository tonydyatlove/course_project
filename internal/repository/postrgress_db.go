package repository

import (
	"github.com/tonydyatlove/leomessi/internal/model"

	"context"

	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

// PRepository p
type PRepository struct {
	Pool *pgxpool.Pool
}


func (p *PRepository) CreateGame(ctx context.Context, person *model.Person) (string, error) {
	newID := uuid.New().String()
	if Game.Score < 0 || Game.Score > 20 {
		return "", fmt.Errorf("database error with create game: age less then 0 or more then 20")
	}
	_, err := p.Pool.Exec(ctx, "insert into persons(id,teams,xG,score,mvp) values($1,$2,$3,$4,$5)",
		newID, &game.Teams, &game.xG, &game.Score, &game.MVP)
	if err != nil {
		log.Errorf("database error with create game: %v", err)
		return "", err
	}
	return newID, nil
}


func (p *PRepository) GetGameByID(ctx context.Context, idPerson string) (*model.Person, error) {
	u := model.Person{}
	err := p.Pool.QueryRow(ctx, "select id,teams,xG,score,mvp from persons where id=$1", idPerson).Scan(
		&u.ID, &uTeams, &u.xG, &u.Score, &u.MVP)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &model.Person{}, fmt.Errorf("game with this id doesnt exist: %v", err)
		}
		log.Errorf("database error, select by id: %v", err)
		return &model.Person{}, err
	}
	return &u, nil
}


func (p *PRepository) GetAllGames(ctx context.Context) ([]*model.Person, error) {
	var persons []*model.Person
	rows, err := p.Pool.Query(ctx, "select id,teams,xG,score from persons")
	if err != nil {
		log.Errorf("database error with select all games, %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		per := model.Person{}
		err = rows.Scan(&per.ID, &per.Teams, &per.xG, &per.Score)
		if err != nil {
			log.Errorf("database error with select all games, %v", err)
			return nil, err
		}
		persons = append(persons, &per)
	}

	return persons, nil
}

func (p *PRepository) DeleteGame(ctx context.Context, id string) error {
	a, err := p.Pool.Exec(ctx, "delete from persons where id=$1", id)
	if a.RowsAffected() == 0 {
		return fmt.Errorf("game with this id doesnt exist")
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("game with this id doesnt exist: %v", err)
		}
		log.Errorf("error with delete game %v", err)
		return err
	}
	return nil
}


func (p *PRepository) UpdateGame(ctx context.Context, id string, per *model.Person) error {
	a, err := p.Pool.Exec(ctx, "update persons set teams=$1,xG=$2,score=$3 where id=$4", &per.Teams, &per.xG, &per.Score, id)
	if per.Score < 0 || per.Score > 20 {
		return fmt.Errorf("database error with create game: score less then 0 or more then 20")
	}
	if a.RowsAffected() == 0 {
		return fmt.Errorf("game with this id doesnt exist")
	}
	if err != nil {
		log.Errorf("error with update game %v", err)
		return err
	}
	return nil
}
