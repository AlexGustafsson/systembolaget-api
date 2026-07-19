package systembolaget

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetStockStatus(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := NewClient(apiKey)

	status, err := client.GetStockStatus(context.TODO(), "0102", "507849")
	require.NoError(t, err)
	fmt.Printf("%+v\n", status)
}
