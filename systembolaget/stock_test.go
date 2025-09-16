package systembolaget

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStock(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := NewClient(apiKey)

	stock, err := client.StockBalance(context.TODO(), "1208", "507811")
	require.NoError(t, err)
	fmt.Printf("%+v\n", stock)
	assert.NotEmpty(t, stock)
}
