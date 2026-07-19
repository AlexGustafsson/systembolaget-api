package systembolaget

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticatedClient_GetStores(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	stores, err := testClient.GetStores(context.TODO())
	require.NoError(t, err)
	fmt.Printf("%+v\n", stores)

	for _, store := range stores {
		if assert.NotNil(t, store) {
			assert.NotEmpty(t, store.SiteID)
		}
	}
}
