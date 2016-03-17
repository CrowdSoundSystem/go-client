package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/crowdsoundsystem/go-client/pkg/bench"
	"github.com/crowdsoundsystem/go-client/pkg/crowdsound"

	"google.golang.org/grpc"
)

var (
	host        = flag.String("hostname", "localhost", "Host of the service")
	port        = flag.Int("port", 50051, "Port of the service")
	numClients  = flag.Int("numClients", 100, "Number of clients to use")
	songsToPost = flag.Int("songsToPost", 10, "Number of songs to post for PostSongs bench")
)

func main() {
	flag.Parse()

	// Run the PostSongs bench with a bunch of clients
	var wg sync.WaitGroup
	wg.Add(*numClients)
	for i := 0; i < *numClients; i++ {
		go func() {
			defer wg.Done()

			conn, err := grpc.Dial(fmt.Sprintf("%v:%v", *host, *port),
				grpc.WithInsecure(),
				grpc.WithTimeout(5*time.Second),
			)
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			client := crowdsound.NewCrowdSoundClient(conn)
			_, err = bench.RunPostSongsBench(client, 1, *songsToPost)
			if err != nil {
				log.Fatal(err)
			}

			// Now just periodically get the queue
			log.Println("Just pollin things")
			for {
				wait := time.Duration(10+rand.Intn(10)) * time.Second
				time.Sleep(wait)

				var count int
				stream, err := client.GetQueue(context.Background(), &crowdsound.GetQueueRequest{})
				for {
					_, err = stream.Recv()
					if err == io.EOF {
						break
					} else if err != nil {
						log.Fatal(err)
					}

					count++
				}

				fmt.Printf("Received queue of %v items", count)
			}

		}()
	}

	wg.Wait()
}
