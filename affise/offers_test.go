package affise

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestClient() *Client {
	panic("not implemented")
}

func TestOffersClient_GetByID(t *testing.T) {
	tests := []struct {
		id   int
		want *Offer
	}{
		{1, nil},
	}

	client := newTestClient()
	for _, test := range tests {
		offer, _, err := client.Offers.GetByID(context.Background(), test.id)
		require.NoError(t, err)
		require.Equal(t, test.want, offer)
	}
}
