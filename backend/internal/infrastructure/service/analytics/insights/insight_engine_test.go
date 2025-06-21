package insights

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRecommendationsEmptyData(t *testing.T) {
	ie := NewInsightEngine()
	recs := ie.GenerateRecommendations(nil, nil, nil)
	assert.NotEmpty(t, recs)
	for _, r := range recs {
		assert.NotEmpty(t, r.Title)
	}
}

func TestGetRecommendationStrings(t *testing.T) {
	ie := NewInsightEngine()
	msgs := ie.GetRecommendationStrings(nil, nil, nil)
	assert.NotEmpty(t, msgs)
	for _, m := range msgs {
		assert.NotEmpty(t, m)
	}
}
