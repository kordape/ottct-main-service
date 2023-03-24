package sqs

import "time"

type FakeNewsEvent struct {
	TweetContent   string    `json:"tweetContent"`
	EntityID       string    `json:"entityId"`
	TweetTimestamp time.Time `json:"tweetTimestamp"`
}
