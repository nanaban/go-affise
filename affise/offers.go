package affise

import (
	"context"
	"fmt"
)

type OffersClient struct {
	client *Client
}

// Certificate defines the schema of an Offer.
type Offer struct {
	ID      int    `json:"id"`
	OfferID string `json:"offer_id"`
	//todo other fields
}

// OfferGetResponse defines the schema of the response when retrieving a single Offer.
type OfferGetResponse struct {
	Status int    `json:"status"` //todo handle status
	Offer  *Offer `json:"offer"`
}

// GetByID gets full information of an offer
func (c *OffersClient) GetByID(ctx context.Context, id int) (*Offer, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/offer/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body OfferGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	return body.Offer, resp, nil
}

// OfferListOpts specifies options for listing Offers.
type OfferListOpts struct {
	Q   string   `json:"q"`   //Search by title and id
	IDs []string `json:"ids"` //Search by string offer ID
	OS  []string `json:"os"`  //OS
	//todo other fields
}

// OfferListResponse defines the schema of the response when listing Offers
type OfferListResponse struct {
	Status int      `json:"status"`
	Offers []*Offer `json:"offers"`
}

// List gets a list of offers
func (c *OffersClient) List(ctx context.Context, opts *OfferListOpts) ([]*Offer, *Response, error) {
	req, err := c.client.NewRequestQueryParams(ctx, "GET", "/offers", opts, nil)
	if err != nil {
		return nil, nil, err
	}

	var body OfferListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	return body.Offers, resp, nil
}
