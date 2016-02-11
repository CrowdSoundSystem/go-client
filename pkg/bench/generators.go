package bench

import (
	"math/rand"
	"time"

	"github.com/crowdsoundsystem/go-client/pkg/crowdsound"
)

var (
	genres = []string{"dance", "crackle", "pop", "drop-it-low", "and other", "things"}
)

func songGenerator(amount int, duplicateRate float64) <-chan *crowdsound.PostSongRequest {
	out := make(chan *crowdsound.PostSongRequest, amount)

	// Let's arbitrarily make 10% of songs non-unique
	duplicateSong := &crowdsound.PostSongRequest{
		Name:   "Random",
		Artist: "Creator",
		Genre:  "Transcending",
	}

	go func() {
		for i := 0; i < amount; i++ {
			if rand.Float64() <= duplicateRate {
				out <- duplicateSong
			} else {
				out <- &crowdsound.PostSongRequest{
					Name:   randomString(10),
					Artist: randomString(5),
					Genre:  genres[i%len(genres)],
				}
			}
		}

		close(out)
	}()

	return out
}

func randomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

	rand.Seed(time.Now().UTC().UnixNano())
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}
