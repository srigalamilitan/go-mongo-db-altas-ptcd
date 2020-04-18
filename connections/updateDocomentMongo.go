package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	client, err:=mongo.NewClient(options.Client().ApplyURI("mongodb+srv://username:password@cluster0-1avcl.mongodb.net/test?retryWrites=true&w=majority"))
	if err!=nil{
		log.Fatal(err);
	}
	ctx,_:=context.WithTimeout(context.Background(),10*time.Second)
	err=client.Connect(ctx)
	if err!=nil{
		log.Fatal(err);
	}
	defer client.Disconnect(ctx)
	database:=client.Database("quickstart")
	podcastsCollection:=database.Collection("podcasts")
	id,_:= primitive.ObjectIDFromHex("5e9a97a0f8978b891fa17b50")
	result,err:= podcastsCollection.UpdateOne(ctx,
		bson.M{"_id":id},
		bson.D{
			{"$set",bson.D{{"author","Anak Agung Krisna Putra"}}},
			},
		)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("update %v Documents!\n",result.ModifiedCount)
	result,err=podcastsCollection.UpdateMany(
		ctx,
		bson.M{"title":"The Polygot Developer Podcast"},
		bson.D{
			{"$set",bson.D{{"author","Anak Agung Maldiva Gandhi"}}},
		},
		)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("update %v Documents!\n",result.ModifiedCount)
	result,err=podcastsCollection.ReplaceOne(
		ctx,
		bson.M{"_id":id},
		bson.M{
			"title": "Rerajahan Ngeleak",
			"author":"Nyen Kaden",
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("update %v Documents!\n",result.ModifiedCount)
}

