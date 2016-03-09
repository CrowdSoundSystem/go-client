package main

import (
	"flag"
	"fmt"
	"io"
	"log"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/crowdsoundsystem/go-client/pkg/playsource"
)

var (
	host     = flag.String("host", "localhost", "Hostname of the service")
	port     = flag.Int("port", 50052, "Port of the service")
	command  = flag.String("cmd", "history", "Command to execute [history, playing, skip]")
	filename = flag.String("filename", "", "File to queue (only valid for queue command")
	songID   = flag.Int("songID", 0, "SongID of the song being queued (only valid for queue command)")
)

func playing(client playsource.PlaysourceClient) {
	resp, err := client.GetPlaying(context.Background(), &playsource.GetPlayingRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Playing: [%v] %v", resp.Song)
}

func skipSong(client playsource.PlaysourceClient) {
	_, err := client.SkipSong(context.Background(), &playsource.SkipSongRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Skipped song")
}

func history(client playsource.PlaysourceClient) {
	stream, err := client.GetPlayHistory(context.Background(), &playsource.GetPlayHistoryRequest{})
	if err != nil {
		log.Fatal(err)
	}

	for {
		song, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error retrieving history:", err)
		}

		log.Println("Song:", song)
	}
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", *host, *port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer conn.Close()

	c := playsource.NewPlaysourceClient(conn)

	switch *command {
	case "history":
		history(c)
		break
	case "playing":
		playing(c)
		break
	case "skip":
		skipSong(c)
	default:
		log.Println("Unrecognized command:", *command)
	}
}
