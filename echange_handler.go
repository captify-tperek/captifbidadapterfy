package main

import (
	"context"
	"time"
)

type Exchange interface {
	sendBidRequest(ctx context.Context, bidRequest BidRequest, segments []Segment) (*BidResponse, error)
}

type ExchangesMultiplexer struct {
	exchanges          []Exchange
	LiveClassification *LiveClassification
}

func NewExchangesMultiplexer(exchanges []Exchange, classification LiveClassification) *ExchangesMultiplexer {
	return &ExchangesMultiplexer{exchanges: exchanges, LiveClassification: &classification}
}

func (e *ExchangesMultiplexer) sendToExchanges(request BidRequest) (*BidResponse, error) {
	// A channel to collect responses from ad exchanges
	responses := make(chan *BidResponse, len(e.exchanges)) // Assuming max 10 exchanges
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	//classify request
	segments, err := e.LiveClassification.classify(ctx, request)
	if err != nil {
		return nil, err
	}

	for _, exchange := range e.exchanges {
		_e := exchange
		go func(request BidRequest) {
			resp, err := _e.sendBidRequest(ctx, request, segments)
			if err == nil {
				responses <- resp
			} else {
				responses <- nil
			}
		}(request)
	}

	bidResponses := make([]*BidResponse, len(e.exchanges))
	for i := 0; i < len(e.exchanges); i++ {
		select {
		case resp := <-responses:
			bidResponses = append(bidResponses, resp)
		case <-ctx.Done():
			break
		}
	}
	// aggregate bid responses into one selecting best site bids for given seat
	return aggregateBestBids(request, bidResponses), nil

}

func aggregateBestBids(bidRequest BidRequest, responses []*BidResponse) *BidResponse {
	bestBids := make(map[string]Bid) // Key: ImpID, Value: Best Bid for that ImpID

	// Iterate through all responses
	for _, response := range responses {
		if response == nil {
			continue
		}
		for _, seatBid := range response.SeatBid {
			for _, bid := range seatBid.Bid {
				// Check if there's already a bid for the same ImpID and replace it if the current one is higher
				if bestBid, exists := bestBids[bid.ImpID]; !exists || bid.Price > bestBid.Price {
					bestBids[bid.ImpID] = bid
				}
			}
		}
	}

	// Create a new BidResponse with the best bids
	finalResponse := &BidResponse{
		ID:      bidRequest.ID,
		SeatBid: []SeatBid{},
		Cur:     "USD", // Assuming currency is the same for all bids or managed elsewhere
	}

	// Map to group the best bids by their Seat to mimic the original structure
	seatMap := make(map[string][]Bid)
	for _, bid := range bestBids {
		seatMap[bid.ImpID] = append(seatMap[bid.ImpID], bid)
	}

	// Populate the SeatBid slice
	for seat, bids := range seatMap {
		finalResponse.SeatBid = append(finalResponse.SeatBid, SeatBid{
			Seat: seat,
			Bid:  bids,
		})
	}

	return finalResponse
}
