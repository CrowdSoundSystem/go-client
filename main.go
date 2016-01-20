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
	host    = flag.String("host", "localhost", "Hostname of the service")
	port    = flag.Int("port", 50051, "Port of the service")
	command = flag.String("cmd", "queue", "Command to execute")
	songs   = []*crowdsound.PostSongRequest{
		&crowdsound.PostSongRequest{Name: "Romeo", Artist: "Taylor Swift", Genre: "Country"},
		&crowdsound.PostSongRequest{Name: "Gay Fish", Artist: "Kanye West", Genre: "Rap"},
	}
)

func printQueue(client crowdsound.CrowdSoundClient) {
	stream, err := client.GetQueue(context.Background(), &crowdsound.GetQueueRequest{})
	if err != nil {
		log.Fatalf("Error calling GetQueue(): %v", err)
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
	flag.Parse()

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", *host, *port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer conn.Close()

	c := crowdsound.NewCrowdSoundClient(conn)

	switch *command {
	case "queue":
		printQueue(c)
		break
	case "post":
		postSongs(c)
		break
	default:
		log.Println("Unrecognized command:", *command)
	}
}
