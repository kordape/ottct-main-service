package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAnalytics(t *testing.T) {

	tweets := []tweet{
		{
			Content: "a",
			Score:   Real,
		},
		{
			Content: "b",
			Score:   Real,
		},
		{
			Content: "c",
			Score:   Real,
		},
		{
			Content: "d",
			Score:   Fake,
		},
	}

	result := getAnalytics(tweets)
	assert.NotNil(t, result)
	assert.Equal(t, 4, result.Total)
	assert.Equal(t, float32(75), result.Authentic)
	assert.Equal(t, float32(25), result.Unauthentic)
}
