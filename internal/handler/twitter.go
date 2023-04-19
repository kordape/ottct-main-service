package handler

import (
	"context"
	"errors"
	"fmt"
	"math"
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
	entityStorage    EntityStorage
}

func NewTwitterManager(validator *validator.Validate, fetcher twitter.TweetsFetcher, classifier predictor.FakeNewsClassifier, entityStorage EntityStorage) (*TwitterManager, error) {

	m := TwitterManager{
		fetcher:          fetcher,
		classifier:       classifier,
		requestValidator: validator,
		entityStorage:    entityStorage,
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

	if m.entityStorage == nil {
		return errors.New("entity storage is nil")
	}

	return nil
}

type Classification int

const (
	Real Classification = 1
	Fake Classification = 0
)

type tweet struct {
	Content string
	Score   Classification
}

func (m *TwitterManager) GetTweets(ctx context.Context, request api.GetTweetsRequest, log *logrus.Entry) (api.GetTweetsResponse, error) {
	err := m.validate()

	if err != nil {
		return api.GetTweetsResponse{}, fmt.Errorf("[TwitterManager] manager validation error: %w", err)
	}

	err = m.requestValidator.Struct(request)
	if err != nil {
		log.WithError(err).Error("[TwitterManager] Invalid GetTweetsRequest request")
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

	entityId, err := m.entityStorage.GetEntity(request.EntityID)
	if err != nil {
		log.WithError(err).Error("[TwitterManager] Failed to get entity by id")
		return api.GetTweetsResponse{}, fmt.Errorf("[TwitterManager] storage error: %w", err)
	}

	fetchRequest := twitter.FetchTweetsRequest{
		MaxResults: maxResults,
		EntityID:   entityId.TwitterId,
		StartTime:  from,
		EndTime:    to,
	}

	resp, err := m.fetcher.FetchTweets(ctx, logger.New("debug"), fetchRequest)

	if err != nil {
		return api.GetTweetsResponse{}, fmt.Errorf("error while fetching tweets: %w", err)
	}

	if len(resp) == 0 {
		return api.GetTweetsResponse{}, nil
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

	tweets := toDomainTweets(resp, classifyResp.Classification)

	return api.GetTweetsResponse{
		Result: getAnalytics(tweets),
	}, nil
}

func toDomainTweets(tweets twitter.FetchTweetsResponse, classifications []predictor.Classification) []tweet {
	dts := make([]tweet, len(tweets))

	for i, t := range tweets {
		dts[i] = tweet{
			Content: t.Text,
			Score:   Classification(classifications[i]),
		}
	}

	return dts
}

func getAnalytics(tweets []tweet) api.Analytics {
	a := api.Analytics{
		Total: len(tweets),
	}

	var authentic int
	var unauthentic int
	for _, t := range tweets {
		if t.Score == Real {
			authentic++
		} else {
			unauthentic++
		}
	}

	a.Authentic = roundFloat((float64(authentic)/float64(a.Total))*100, 1)
	a.Unauthentic = roundFloat((float64(unauthentic)/float64(a.Total))*100, 1)

	return a
}

func roundFloat(val float64, precision uint) float32 {
	ratio := math.Pow(10, float64(precision))
	return float32(math.Round(val*ratio) / ratio)
}
