package systembolaget

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchWithCursor(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := NewClient(apiKey)

	cursor := client.SearchWithCursor(nil, FilterByQuery("Carlsberg"))

	items := 0
	for cursor.Next(context.TODO(), 0) {
		require.NoError(t, cursor.Error())

		product := cursor.At()

		// Assume there are more than 5 results for Carlsberg
		assert.NotNil(t, product)
		assert.Greater(t, len(product), 0)

		items++
		if items >= 5 {
			break
		}
	}

	assert.Equal(t, 5, items)
}
