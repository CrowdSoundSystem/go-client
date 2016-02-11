package main

import (
	"github.com/crowdsoundsystem/go-client/pkg/bench"
	"github.com/crowdsoundsystem/go-client/pkg/crowdsound"
	"github.com/rcrowley/go-metrics"
)

type Benchmark func(client crowdsound.CrowdSoundClient) (name string, results metrics.Histogram, err error)

var benchmarks = []Benchmark{
	func(c crowdsound.CrowdSoundClient) (name string, results metrics.Histogram, err error) {
		name = "Get Session Data"
		results, err = bench.RunGetSessionDataBench(c, *iterations)
		return
	},
	func(c crowdsound.CrowdSoundClient) (name string, results metrics.Histogram, err error) {
		name = "Get Queue"
		results, err = bench.RunGetQueueBench(c, *iterations)
		return
	},
	func(c crowdsound.CrowdSoundClient) (name string, results metrics.Histogram, err error) {
		name = "List Trending Artists"
		results, err = bench.RunListTrendingArtists(c, *iterations)
		return
	},
	func(c crowdsound.CrowdSoundClient) (name string, results metrics.Histogram, err error) {
		name = "Post Songs"
		results, err = bench.RunPostSongsBench(c, *iterations/(*songsToPost / *iterations), *songsToPost)
		return
	},
	func(c crowdsound.CrowdSoundClient) (name string, results metrics.Histogram, err error) {
		name = "Vote Song"
		results, err = bench.RunVoteSongsBench(c, *iterations)
		return
	},
}
