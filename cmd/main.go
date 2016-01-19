package main

import (
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/crowdsoundsystem/go-client/pkg/crowdsound"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

var (
	host  = flag.String("host", "localhost", "Hostname of the service")
	port  = flag.Int("port", 50051, "Port of the service")
	songs = []*crowdsound.PostSongRequest{
		&crowdsound.PostSongRequest{Name: "Romeo", Artist: "Taylor Swift", Genre: "Country"},
		&crowdsound.PostSongRequest{Name: "Gay Fish", Artist: "Kanye West", Genre: "Rap"},
	}
)

func printSongs(client crowdsound.CrowdSoundClient) {
	flag.Parse()

	log.Println("Retrieving songs...")

	stream, err := client.GetQueue(context.Background(), &crowdsound.GetQueueRequest{})
	if err != nil {
		log.Fatalf("Error calling ListSongs(): %v", err)
	}

	for {
		song, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error retrieving song: %v", err)
		}

		log.Printf("Song: [%v] %v - %v", song.Genre, song.Artist, song.Name)
	}
}

func postSongs(client crowdsound.CrowdSoundClient) {
	log.Println("Posting songs...")

	stream, err := client.PostSong(context.Background())
	if err != nil {
		log.Fatalf("Error calling PostSong(): %v", err)
	}

	for _, song := range songs {
		err := stream.Send(song)
		if err != nil {
			log.Fatalf("Error sending song: %v", err)
		}
	}
}

func main() {
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", *host, *port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer conn.Close()

	c := crowdsound.NewCrowdSoundClient(conn)

	postSongs(c)
	printSongs(c)
}
