package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/crowdsoundsystem/go-client/pkg/crowdsound"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

var (
	host    = flag.String("host", "localhost", "Hostname of the service")
	port    = flag.Int("port", 50051, "Port of the service")
	userID  = flag.String("user", "test_golang_user", "User ID when performing RPC calls")
	command = flag.String("cmd", "queue", "Command to execute [ping, queue, post, vote, voteSkip, meta]")
	songs   = []*crowdsound.PostSongRequest{
		&crowdsound.PostSongRequest{Name: "Shivers", Artist: []string{"Armin van Buuren"}, Genre: "Trance"},
		&crowdsound.PostSongRequest{Name: "Games (Standerwick Remix)", Artist: []string{"John O'Callaghan"}, Genre: "Trance"},
		// This will yield search results, but NOT be able to play, because T-swift thinks spotify will steal all of her business!
		// Get it together, T-Swift. You could have been part of something big!
		&crowdsound.PostSongRequest{Name: "Love Story", Artist: []string{"Taylor Swift"}, Genre: "Country"},
		&crowdsound.PostSongRequest{Name: "Never Cry Again", Artist: []string{"Dash Berlin"}, Genre: "Trance"},
		// This guy, on the other hand, is a legit gangster. If he thought people were stealing
		// from him, he'd just shoot them up. Yet, Spotify has yet to be shot...
		&crowdsound.PostSongRequest{Name: "What You Know", Artist: []string{"T.I."}, Genre: "Phils Genre"},
		&crowdsound.PostSongRequest{Name: "Test song", Genre: "Phils Genre"},
	}
)

func ping(client crowdsound.CrowdSoundClient) {
	_, err := client.Ping(context.Background(), &crowdsound.PingRequest{UserId: *userID})
	if err != nil {
		log.Fatalf("Error calling Ping(): %v", err)
	}
}

func printMeta(client crowdsound.CrowdSoundClient) {
	resp, err := client.GetSessionData(context.Background(), &crowdsound.GetSessionDataRequest{})
	if err != nil {
		log.Fatalf("Error calling GetSessionData(): %v", err)
	}

	log.Println("Session Data:", resp)
}

func printPlaying(client crowdsound.CrowdSoundClient) {
	resp, err := client.GetPlaying(context.Background(), &crowdsound.GetPlayingRequest{})
	if err != nil {
		log.Fatalf("Error calling GetPlaying(): %v", err)
	}

	log.Println("Playing:", resp)
}

func printQueue(client crowdsound.CrowdSoundClient) {
	stream, err := client.GetQueue(context.Background(), &crowdsound.GetQueueRequest{})
	if err != nil {
		log.Fatalf("Error calling GetQueue(): %v", err)
	}

	var count int

	for {
		song, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error retrieving song: %v", err)
		}

		if song.IsBuffered {
			if song.IsPlaying {
				log.Printf("Song: [%v] %v - %v (buffered) (playing)", song.Genre, song.Artist, song.Name)
			} else {
				log.Printf("Song: [%v] %v - %v (buffered)", song.Genre, song.Artist, song.Name)
			}
		} else {
			log.Printf("Song: [%v] %v - %v", song.Genre, song.Artist, song.Name)
		}
		count++
	}

	log.Println("Queue size:", count)
}

func postSongs(client crowdsound.CrowdSoundClient) {
	stream, err := client.PostSong(context.Background())
	if err != nil {
		log.Fatalf("Error calling PostSong(): %v", err)
	}
	defer stream.CloseAndRecv()

	for _, song := range songs {
		log.Println("Posting song:", song)
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
		Name:   songs[4].Name,
		Artist: songs[4].Artist[0],
		Like:   true,
	})
	if err != nil {
		log.Fatalf("Error calling VoteSong(): %v", err)
	}
}

func voteSkip(client crowdsound.CrowdSoundClient) {
	_, err := client.VoteSkip(context.Background(), &crowdsound.VoteSkipRequest{
		UserId: *userID,
	})
	if err != nil {
		log.Fatalf("Error calling VoteSkip(): %v", err)
	}

	log.Println("Voted for skip")
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", *host, *port),
		grpc.WithInsecure(),
		grpc.WithTimeout(10*time.Second),
	)

	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer conn.Close()

	c := crowdsound.NewCrowdSoundClient(conn)

	switch *command {
	case "ping":
		ping(c)
		break
	case "playing":
		printPlaying(c)
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
	case "voteSkip":
		voteSkip(c)
	case "meta":
		printMeta(c)
		break
	default:
		log.Println("Unrecognized command:", *command)
	}
}
