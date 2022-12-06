// Package repository : file contains operations with MongoDB
package repository

import (
	"context"
	"fmt"

	"github.com/tonydyatlove/leomessi/internal/model"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type MRepository struct {
	Pool *mongo.Client
}

func (m *MRepository) CreateGame(ctx context.Context, person *model.Person) (string, error) {
	if game.Score < 0 || game.Score > 20 {
		return "", fmt.Errorf("mongo repository: error with create, score must be more then 0 and less then 20")
	}
	newID := uuid.New().String()
	collection := m.Pool.Database("person").Collection("person")
	_, err := collection.InsertOne(ctx, bson.D{
		{Key: "id", Value: newID},
		{Key: "teams", Value: game.Teams},
		{Key: "xG", Value: game.xG},
		{Key: "score", Value: game.Score},
		{Key: "mvp", Value: game.MVP},
	
	})
	if err != nil {
		return "", fmt.Errorf("mongo: unable to create new game: %v", err)
	}
	return newID, nil
}


func (m *MRepository) UpdateGame(ctx context.Context, id string, person *model.Person) error {
	if game.Score < 0 || game.Score > 20 {
		return fmt.Errorf("mongo repository: error with create, game`s score must be more then 0 and less then 20")
	}
	collection := m.Pool.Database("game").Collection("game")
	_, err := collection.UpdateOne(ctx, bson.D{primitive.E{Key: "id", Value: id}}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "teams", Value: game.Teams},
		{Key: "xG", Value: game.xG},
		{Key: "score", Value: game.Score},
	}}})
	if err != nil {
		return fmt.Errorf("mongo: unable to update game %v", err)
	}
	return nil
}

// GetAllUsers take all users from db
func (m *MRepository) GetAllGames(ctx context.Context) ([]*model.Person, error) {
	var users []*model.Person
	collection := m.Pool.Database("person").Collection("person")
	c, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("mongo: unable to select all games %v", err)
	}
	for c.Next(ctx) {
		game := model.Person{}
		err := c.Decode(&game)
		if err != nil {
			return games, err
		}
		users = append(games, &game)
	}
	return games, nil
}

// DeleteUser user from db
func (m *MRepository) DeleteGame(ctx context.Context, id string) error {
	collection := m.Pool.Database("person").Collection("person")
	_, err := collection.DeleteOne(ctx, bson.D{primitive.E{Key: "id", Value: id}})
	if err != nil {
		return fmt.Errorf("mongo: unable to delete game, %v", err)
	}
	return nil
}

// GetUserByID select exist user from db by his id
func (m *MRepository) GetGameByID(ctx context.Context, id string) (*model.Person, error) {
	user := model.Person{}
	collection := m.Pool.Database("person").Collection("person")
	err := collection.FindOne(ctx, bson.D{primitive.E{Key: "id", Value: id}}).Decode(&user)
	if err != nil {
		return &game, err
	}
	return &game, nil
}
