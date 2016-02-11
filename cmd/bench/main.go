package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/crowdsoundsystem/go-client/pkg/crowdsound"
)

var (
	host        = flag.String("hostname", "localhost", "Host of the service")
	port        = flag.Int("port", 50051, "Port of the service")
	iterations  = flag.Int("interations", 100, "Number of iterations per benchmark")
	songsToPost = flag.Int("songsToPost", 1000, "Number of songs to post for PostSongs bench")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", *host, *port),
		grpc.WithInsecure(),
		grpc.WithTimeout(3*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := crowdsound.NewCrowdSoundClient(conn)

	for _, b := range benchmarks {
		benchName, results, err := b(client)
		if err != nil {
			log.Fatal("Failed to run benchmark:", err)
		}

		log.Printf("Benchmark results: %v", benchName)
		log.Printf("--------------------------------")
		log.Printf("Min: \t%v μs", results.Min()/1000)
		log.Printf("Max: \t%v μs", results.Max()/1000)
		log.Printf("Mean: \t%v μs", results.Mean()/1000)
		log.Printf("Std: \t%v μs", results.StdDev()/1000)
		log.Print("")
	}
}
