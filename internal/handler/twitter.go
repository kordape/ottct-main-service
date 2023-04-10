package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kordape/ottct-main-service/pkg/api"
	"github.com/kordape/ottct-main-service/pkg/logger"
	"github.com/kordape/ottct-poller-service/pkg/predictor"
	"github.com/kordape/ottct-poller-service/pkg/twitter"
	"github.com/sirupsen/logrus"
)

const (
	defaultMaxResults        = 50
	defaultTimeWindowInHours = 24
)

type TwitterManager struct {
	requestValidator *validator.Validate
	fetcher          twitter.TweetsFetcher
	classifier       predictor.FakeNewsClassifier
}

func NewTwitterManager(validator *validator.Validate, fetcher twitter.TweetsFetcher, classifier predictor.FakeNewsClassifier) (*TwitterManager, error) {

	m := TwitterManager{
		fetcher:          fetcher,
		classifier:       classifier,
		requestValidator: validator,
	}

	err := m.validate()

	if err != nil {
		return &m, fmt.Errorf("[TwitterManager] validation error: %w", err)
	}

	return &m, nil
}

func (m TwitterManager) validate() error {
	if m.fetcher == nil {
		return errors.New("fetcher is nil")
	}

	if m.fetcher == nil {
		return errors.New("classifier is nil")
	}

	if m.requestValidator == nil {
		return errors.New("request validator is nil")
	}

	return nil
}

func (m *TwitterManager) GetTweets(ctx context.Context, request api.GetTweetsRequest, log *logrus.Entry) (api.GetTweetsResponse, error) {
	err := m.validate()

	if err != nil {
		return api.GetTweetsResponse{}, fmt.Errorf("[TwitterManager] manager validation error: %w", err)
	}

	err = m.requestValidator.Struct(request)
	if err != nil {
		log.Error(fmt.Errorf("[TwitterManager] Invalid GetTweetsRequest request: %w", err))
		return api.GetTweetsResponse{}, ErrInvalidRequest
	}

	maxResults := defaultMaxResults
	if request.MaxResults > 0 {
		maxResults = request.MaxResults
	}

	to := time.Now()
	if !request.To.IsZero() {
		to = request.To
	}

	from := to.Add(-time.Hour * defaultTimeWindowInHours)
	if !request.From.IsZero() {
		from = request.From
		if from.After(to) {
			return api.GetTweetsResponse{}, ErrInvalidRequest
		}
	}

	fetchRequest := twitter.FetchTweetsRequest{
		MaxResults: maxResults,
		EntityID:   request.EntityID,
		StartTime:  from,
		EndTime:    to,
	}

	resp, err := m.fetcher.FetchTweets(ctx, logger.New("debug"), fetchRequest)

	if err != nil {
		return api.GetTweetsResponse{}, fmt.Errorf("error while fetching tweets: %w", err)
	}

	classifyRequest := make([]string, len(resp))
	for i, c := range resp {
		classifyRequest[i] = c.Text
	}
	classifyResp, err := m.classifier.Classify(ctx, classifyRequest)
	if err != nil {
		return api.GetTweetsResponse{}, fmt.Errorf("error while classifying tweets: %w", err)
	}

	if len(classifyResp.Classification) != len(resp) {
		return api.GetTweetsResponse{}, fmt.Errorf("error in mismatched number of classifications: %w", err)
	}

	tweets := make([]api.Tweet, len(resp))

	for i, t := range resp {
		tweets[i] = api.Tweet{
			ID:            t.ID,
			Content:       t.Text,
			CreatedAt:     t.CreatedAt,
			RealnessScore: float32(classifyResp.Classification[i]),
		}
	}

	return api.GetTweetsResponse{
		Tweets: tweets,
	}, nil
}
