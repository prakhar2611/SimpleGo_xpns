package Utilities

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConectToMongo() (*mongo.Client, context.Context) {
	k := GetKaonf()
	username := k.String("AtlasConfig.userName")
	pass := k.String("AtlasConfig.password")
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.tb5xa.mongodb.net/?retryWrites=true&w=majority", username, pass)).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client, ctx
}
