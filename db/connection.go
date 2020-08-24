package db

import (
	"bufio"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strings"
)

var (
	MongoCN       = ConnectDB()
	clientOptions *options.ClientOptions
)

func ConnectDB() *mongo.Client {
	opts()
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err.Error())
		return client
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err.Error())
		return client
	}

	log.Println("Database connection successful")
	return client
}

func opts() {
	var uri string
	f, err := os.Open("./props.txt")
	if err != nil {
		log.Fatalf("connection: %v", err)
	}
	input := bufio.NewScanner(f)
	for input.Scan() {
		lines := strings.SplitN(input.Text(), "=", 2)
		if lines[0] == "db" {
			uri = strings.Trim(lines[1], " \"")
		}
	}
	clientOptions = options.Client().ApplyURI(uri)
	_ = f.Close()
}
