package db

import (
	"context"
	"curse-count/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func InsertPerson(person models.Person) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := MongoCN.Database("curse-count").Collection("person")

	result, err := col.InsertOne(ctx, person)
	if err != nil {
		log.Printf("insert person: %v", err)
		return "", err
	}

	id, _ := result.InsertedID.(primitive.ObjectID)

	return id.Hex(), nil
}

func GetAll() ([]*models.Person, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := MongoCN.Database("curse-count").Collection("person")

	var results []*models.Person
	cursor, err := col.Find(ctx, bson.D{}) // Para incluir todos los documentos se usa un documento vacío como filtro

	if err != nil {
		log.Printf("getall person: %v", err)
		return results, err
	}

	for cursor.Next(ctx) {
		var person models.Person
		err := cursor.Decode(&person)
		if err != nil {
			log.Printf("getall person: %v", err)
			return results, err
		}

		results = append(results, &person)
	}

	if err = cursor.Err(); err != nil {
		log.Printf("getall person: %v", err)
		return results, err
	}

	_ = cursor.Close(ctx)

	return results, nil
}

func Count(id string, qt int) (*models.Person, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var p *models.Person

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("count person: %v", err)
		return p, false
	}

	gteValue := 1
	if qt == 1 {
		gteValue = 0
	}

	col := MongoCN.Database("curse-count").Collection("person")
	filter := bson.D{
		{"_id", _id},
		{"curses", bson.D{{"$gte", gteValue}},
		},
	}
	update := bson.D{{"$inc", bson.D{{"curses", qt}}}}
	after := options.After
	opts := options.FindOneAndUpdateOptions{ReturnDocument: &after}
	result := col.FindOneAndUpdate(ctx, filter, update, &opts)

	if result.Err() != nil {
		log.Printf("count person: %v", result.Err())
		return p, false
	}

	err = result.Decode(&p)
	if err != nil {
		log.Printf("count person: %v", err)
		return p, false
	}

	return p, true
}
