package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	err=client.Ping(ctx,readpref.Primary())
	if err!=nil{
		log.Fatal(err)
	}
	log.Print("Done")
	databases, err:=client.ListDatabaseNames(ctx,bson.M{})
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println(databases);
	quickStartDatabase:=client.Database("quickstart")
	podcastCollection :=quickStartDatabase.Collection("podcasts")
	episodeCollection :=quickStartDatabase.Collection("episodes")

	podCastResult,err:=podcastCollection.InsertOne(ctx,bson.D{
		{"title","The Polygot Developer Podcast"},
		{"author","Krisna Putra"},
		{"tags",bson.A{"Development","programming","coding"}},
	})
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println(podCastResult.InsertedID)
	episodeResult,err:=episodeCollection.InsertMany(ctx,[]interface{}{
		bson.D{
			{"podcast",podCastResult.InsertedID},
			{"title","Episode #1"},
			{"description","This is The first Episode"},
			{"duration",25},
		},
		bson.D{
			{"podcast",podCastResult.InsertedID},
			{"title","Episode #2"},
			{"description","This is The second Episode"},
			{"duration",30},
		},
	})
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println(episodeResult.InsertedIDs)
}
