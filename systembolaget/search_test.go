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

	cursor := client.SearchWithCursor(nil, FilterByCategory("Alkoholfritt", "Öl", ""))

	results, err := client.Search(context.TODO(), nil, FilterByCategory("Alkoholfritt", "Öl", ""))
	require.NoError(t, err)
	totalItems := results.Metadata.FullAssortmentDocumentCount

	yieldedItems := 0
	for cursor.Next(context.TODO(), 0) {
		require.NoError(t, cursor.Error())

		product := cursor.At()

		assert.NotNil(t, product)
		assert.Greater(t, len(product), 0)

		yieldedItems++
	}

	// Ensure that we retrieved all products
	assert.Equal(t, totalItems, yieldedItems)
	// Ensure that there were multiple pages processed
	assert.Greater(t, yieldedItems, 30)
}
