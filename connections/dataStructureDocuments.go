package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type Podcast struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `bson:"title,omitempty"`
	Author string             `bson:"author,omitempty"`
	Tags   []string           `bson:"tags,omitempty"`
}
type Episode struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Podcast     primitive.ObjectID `bson:"podcast,omitempty"`
	Title       string             `bson:"title,omitempty"`
	Description string             `bson:"description,omitempty"`
	Duration    int32              `bson:"duration,omitempty"`
}


func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://username:password@cluster0-1avcl.mongodb.net/test?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Done")
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
	quickStartDatabase := client.Database("quickstart")
	podcastCollection := quickStartDatabase.Collection("podcasts")
	episodeCollection := quickStartDatabase.Collection("episodes")

	mongoPodcast :=Podcast{
		Title:  "The Monggo Db",
		Author: "Monggo",
		Tags:   []string{"monggo","coba mongo"},
	}
	insertResult,err:= podcastCollection.InsertOne(ctx,mongoPodcast)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertResult)
	cursor, err := episodeCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Println(err)
	}
	//var episodes []bson.M
	//if err=cursor.All(ctx,&episodes);err!=nil{
	//	log.Fatal(err);
	//}
	//for _,episode:=range episodes{
	//	fmt.Println(episode)
	//}
	defer cursor.Close(ctx)
	/*
		find all data using batching
	*/
	fmt.Println("Find All Data Using Batching")
	for cursor.Next(ctx) {
		var episode Episode
		if err = cursor.Decode(&episode); err != nil {
			log.Fatal(err)
		}
		fmt.Println(episode)
	}
	/*
		find one data , top one data because not using filter
	*/
	fmt.Println("Find One From PodcastCollection")
	var podcast Podcast
	if err = podcastCollection.FindOne(ctx, bson.M{}).Decode(&podcast); err != nil {
		log.Fatal(err)
	}
	fmt.Println(podcast)

	/*
		filter data Episode using duration equals 25
	*/
	filterCurson, err := episodeCollection.Find(ctx, bson.M{"duration": 25})
	if err != nil {
		log.Fatal(err)
	}
	var episodesFilter []Episode
	if err = filterCurson.All(ctx, &episodesFilter); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Episode Filter by : duration equals 25")
	fmt.Println(episodesFilter)

	fmt.Println("Sorting Data ")
	opts := options.Find()
	opts.SetSort(bson.D{{"duration", -1,}})
	sortCursor, err := episodeCollection.Find(ctx, bson.D{
		{"duration", bson.D{
			{"$gt", 24},
		},},
	}, opts)
	var episodeSorted []Episode
	if err = sortCursor.All(ctx, &episodeSorted); err != nil {
		log.Fatal(err)
	}
	fmt.Println(episodeSorted)

}
