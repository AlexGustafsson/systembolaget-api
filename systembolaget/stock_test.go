package systembolaget

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthenticatedClient_GetStockStatus(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	status, err := testClient.GetStockStatus(context.TODO(), "0102", "507849")
	require.NoError(t, err)
	fmt.Printf("%+v\n", status)
}
