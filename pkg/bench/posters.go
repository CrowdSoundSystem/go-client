package bench

import (
	"time"

	"github.com/crowdsoundsystem/go-client/pkg/crowdsound"
	"github.com/rcrowley/go-metrics"
	"golang.org/x/net/context"
)

func RunGetSessionDataBench(c crowdsound.CrowdSoundClient, amount int) (metrics.Histogram, error) {
	h := metrics.NewHistogram(metrics.NewUniformSample(amount))
	for i := 0; i < amount; i++ {
		start := time.Now()
		if _, err := c.GetSessionData(context.Background(), &crowdsound.GetSessionDataRequest{}); err != nil {
			return h, err
		}

		h.Update(time.Since(start).Nanoseconds())
	}

	return h, nil
}

func RunPostSongsBench(c crowdsound.CrowdSoundClient, amount, songsToPost int) (metrics.Histogram, error) {
	h := metrics.NewHistogram(metrics.NewUniformSample(amount))

	userID := randomString(32)
	for i := 0; i < amount; i++ {
		start := time.Now()

		stream, err := c.PostSong(context.Background())
		if err != nil {
			return h, err
		}

		for song := range songGenerator(songsToPost, 0.1) {
			song.UserId = userID
			if err = stream.Send(song); err != nil {
				return h, err
			}
		}

		stream.CloseAndRecv()

		h.Update(time.Since(start).Nanoseconds())
	}

	return h, nil
}

func RunVoteSongsBench(c crowdsound.CrowdSoundClient, amount int) (metrics.Histogram, error) {
	h := metrics.NewHistogram(metrics.NewUniformSample(amount))

	// Generate songs to vote on
	postStream, err := c.PostSong(context.Background())
	if err != nil {
		return h, err
	}

	var songToVote *crowdsound.PostSongRequest
	for song := range songGenerator(1, 0) {
		song.UserId = randomString(32)
		if err = postStream.Send(song); err != nil {
			return h, err
		}
		songToVote = song
	}

	if _, err = postStream.CloseAndRecv(); err != nil {
		return h, err
	}

	userID := randomString(32)
	for i := 0; i < amount; i++ {
		start := time.Now()
		_, err := c.VoteSong(context.Background(), &crowdsound.VoteSongRequest{
			UserId: userID,
			Name:   songToVote.Name,
			Artist: songToVote.Artist[0],
			Like:   true,
		})
		if err != nil {
			return h, err
		}

		h.Update(time.Since(start).Nanoseconds())
	}

	return h, nil
}
