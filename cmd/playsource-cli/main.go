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
	command  = flag.String("cmd", "history", "Command to execute [queue, history, playing]")
	filename = flag.String("filename", "", "File to queue (only valid for queue command")
	songID   = flag.Int("songID", 0, "SongID of the song being queued (only valid for queue command)")
)

func queueSong(client playsource.PlaySourceClient) {
	resp, err := client.QueueSong(context.Background(), &playsource.QueueSongRequest{
		SongId:   int32(*songID),
		Filename: *filename,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Queued: %t", resp.Queued)
}

func playing(client playsource.PlaySourceClient) {
	resp, err := client.GetPlaying(context.Background(), &playsource.GetPlayingRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Playing: [%v] %v", resp.SongId, resp.Filename)
}
func history(client playsource.PlaySourceClient) {
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

		log.Printf("Song: [%v] %v", song.SongId, song.Filename)
	}
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", *host, *port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer conn.Close()

	c := playsource.NewPlaySourceClient(conn)

	switch *command {
	case "queue":
		queueSong(c)
		break
	case "history":
		history(c)
		break
	case "playing":
		playing(c)
		break
	default:
		log.Println("Unrecognized command:", *command)
	}
}
