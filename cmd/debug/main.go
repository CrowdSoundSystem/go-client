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
	userID  = flag.String("user", "test_golang_user", "User ID when performing RPC calls")
	command = flag.String("cmd", "queue", "Command to execute")
	songs   = []*crowdsound.PostSongRequest{
		&crowdsound.PostSongRequest{Name: "Romeo", Artist: "Taylor Swift", Genre: "Country"},
		&crowdsound.PostSongRequest{Name: "Gay Fish", Artist: "Kanye West", Genre: "Rap"},
	}
)

func ping(client crowdsound.CrowdSoundClient) {
	_, err := client.Ping(context.Background(), &crowdsound.PingRequest{UserId: *userID})
	if err != nil {
		log.Fatalf("Error calling Ping(): %v", err)
	}
}

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
	defer stream.CloseSend()

	for _, song := range songs {
		song.UserId = *userID
		err := stream.Send(song)
		if err != nil {
			log.Fatalf("Error sending song: %v", err)
		}
	}
}

func vote(client crowdsound.CrowdSoundClient) {
	_, err := client.VoteSong(context.Background(), &crowdsound.VoteSongRequest{
		UserId: *userID,
		Name:   songs[0].Name,
		Artist: songs[0].Artist,
		Like:   true,
	})
	if err != nil {
		log.Fatalf("Error calling VoteSong(): %v", err)
	}

	_, err = client.VoteSong(context.Background(), &crowdsound.VoteSongRequest{
		UserId: *userID,
		Name:   songs[1].Name,
		Artist: songs[1].Artist,
		Like:   false,
	})
	if err != nil {
		log.Fatalf("Error calling VoteSong(): %v", err)
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
	case "ping":
		ping(c)
		break
	case "queue":
		printQueue(c)
		break
	case "post":
		postSongs(c)
		break
	case "vote":
		vote(c)
		break
	default:
		log.Println("Unrecognized command:", *command)
	}
}
