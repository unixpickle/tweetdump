package main

import (
	"encoding/csv"
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/unixpickle/essentials"
)

func main() {
	var consumerKey, consumerSecret string
	var accessToken, accessSecret string
	var language string
	flag.StringVar(&consumerKey, "consumer-key", "", "consumer key")
	flag.StringVar(&consumerSecret, "consumer-secret", "", "consumer secret")
	flag.StringVar(&accessToken, "access-token", "", "access token")
	flag.StringVar(&accessSecret, "access-secret", "", "access secret")
	flag.StringVar(&language, "language", "en", "tweet language filter")
	flag.Parse()
	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	params := &twitter.StreamSampleParams{
		StallWarnings: twitter.Bool(true),
		Language:      []string{language},
	}
	stream, err := client.Streams.Sample(params)
	if err != nil {
		essentials.Die("streaming:", err)
	}

	writer := csv.NewWriter(os.Stdout)
	for msg := range stream.Messages {
		if tweet, ok := msg.(*twitter.Tweet); ok {
			parsed, err := time.Parse(time.RubyDate, tweet.CreatedAt)
			if err != nil {
				essentials.Die(err)
			}
			t := strconv.FormatInt(parsed.Unix(), 10)
			writer.Write([]string{tweet.IDStr, tweet.User.ScreenName, t, tweet.Text})
			writer.Flush()
		}
	}
}
