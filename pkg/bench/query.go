package bench

import (
	"io"
	"log"
	"time"

	"github.com/crowdsoundsystem/go-client/pkg/crowdsound"
	"github.com/rcrowley/go-metrics"
	"golang.org/x/net/context"
)

func RunGetQueueBench(c crowdsound.CrowdSoundClient, amount int) (metrics.Histogram, error) {
	h := metrics.NewHistogram(metrics.NewUniformSample(amount))

	// To make the queue meaningful, we want there to actually
	// be a queue. Idealy, we could try various queue sizes, but
	// time is of the essence!
	postStream, err := c.PostSong(context.Background())
	if err != nil {
		return h, err
	}

	for song := range songGenerator(10, 0) {
		song.UserId = randomString(32)
		if err = postStream.Send(song); err != nil {
			return h, err
		}
	}

	if _, err = postStream.CloseAndRecv(); err != nil {
		return h, err
	}

	for i := 0; i < amount; i++ {
		start := time.Now()
		queueStream, err := c.GetQueue(context.Background(), &crowdsound.GetQueueRequest{})
		if err != nil {
			return h, err
		}

		for {
			_, err := queueStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
		}

		h.Update(time.Since(start).Nanoseconds())
	}

	return h, nil
}

func RunListTrendingArtists(c crowdsound.CrowdSoundClient, amount int) (metrics.Histogram, error) {
	h := metrics.NewHistogram(metrics.NewUniformSample(amount))

	// To make trending artists meaningful, add a bunch of artists.
	postStream, err := c.PostSong(context.Background())
	if err != nil {
		return h, err
	}

	for song := range songGenerator(100, 1.0) {
		song.UserId = randomString(32)
		if err = postStream.Send(song); err != nil {
			return h, err
		}
	}

	if _, err = postStream.CloseAndRecv(); err != nil {
		return h, err
	}

	for i := 0; i < amount; i++ {
		start := time.Now()
		listStream, err := c.ListTrendingArtists(context.Background(), &crowdsound.ListTrendingArtistsRequest{})
		if err != nil {
			return h, err
		}

		for {
			_, err := listStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
		}

		h.Update(time.Since(start).Nanoseconds())
	}

	return h, nil
}
