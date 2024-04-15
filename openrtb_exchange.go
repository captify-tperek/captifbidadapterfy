package main

import (
	"context"
	"encoding/json"
	"github.com/bsm/openrtb/v3"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/samber/lo"
)

type OpenRTBExchange struct {
	creatives []BannerConfig
}

func (e *OpenRTBExchange) sendBidRequest(ctx context.Context, bidRequest BidRequest, segments []Segment) (*BidResponse, error) {
	// Convert BidRequest to openrtb.BidRequest

	openRTBBidRequest := transformToOpenRTB(&bidRequest, segments)
	marshal, err := json.Marshal(openRTBBidRequest)
	if err != nil {
		return nil, err
	}

	log.Infof("transformed request: %s", string(marshal))

	openRtbBidResponse := e.generateResponse(openRTBBidRequest)
	marshal, err = json.Marshal(openRtbBidResponse)
	if err != nil {
		return nil, err
	}

	log.Infof("received response: %s", string(marshal))
	return transformOpenRTBResp(openRtbBidResponse), nil

}

func transformOpenRTBResp(openRTBResponse *openrtb.BidResponse) *BidResponse {
	// Create a new instance of BidResponse
	bidResponse := &BidResponse{
		ID:      openRTBResponse.ID,
		SeatBid: transformSeatBids(openRTBResponse.SeatBids),
		Cur:     "USD",
	}
	return bidResponse
}

func transformSeatBids(bids []openrtb.SeatBid) []SeatBid {
	var seatBids []SeatBid
	for _, bid := range bids {
		seatBid := SeatBid{
			Bid:  transformBids(bid.Bids),
			Seat: bid.Seat,
		}
		seatBids = append(seatBids, seatBid)
	}
	return seatBids

}

func transformBids(bids []openrtb.Bid) []Bid {
	var transformedBids []Bid
	for _, bid := range bids {
		transformedBid := Bid{
			ID:      bid.ID,
			ImpID:   bid.ImpID,
			Price:   bid.Price,
			Adm:     bid.AdMarkup,
			Crid:    bid.CreativeID,
			W:       bid.Width,
			H:       bid.Height,
			Adomain: bid.AdvDomains,
			Iurl:    bid.ImageURL,
			Nurl:    bid.NoticeURL,
		}
		transformedBids = append(transformedBids, transformedBid)
	}
	return transformedBids

}

func (e *OpenRTBExchange) generateResponse(openRTBBidRequest *openrtb.BidRequest) *openrtb.BidResponse {
	// find matching banners in map
	bidResponse := openrtb.BidResponse{ID: openRTBBidRequest.ID}

	for _, imp := range openRTBBidRequest.Impressions {
		seatBid := openrtb.SeatBid{}
		for _, format := range imp.Banner.Formats {
			matchingBanner := e.findMatchingBanner(format.Height, format.Width, openRTBBidRequest.Site.Content.Data)

			if matchingBanner != nil {
				seatBid.Bids = append(seatBid.Bids, openrtb.Bid{
					ID:         uuid.New().String(),
					ImpID:      imp.ID,
					Price:      1.0,
					AdMarkup:   matchingBanner.AdMarkup,
					CreativeID: matchingBanner.CreativeID,
					Width:      matchingBanner.Width,
					Height:     matchingBanner.Height,
					AdvDomains: matchingBanner.AdvDomains,
					ImageURL:   matchingBanner.ImageURL,
					NoticeURL:  matchingBanner.NoticeURL,
				})
			}
		}
		bidResponse.SeatBids = append(bidResponse.SeatBids, seatBid)

	}
	return &bidResponse
}

func (e *OpenRTBExchange) findMatchingBanner(height int, width int, data []openrtb.Data) *BannerConfig {
	captifyData, found := lo.Find(data, func(d openrtb.Data) bool {
		return d.ID == "captify"
	})
	if found {
		segments := lo.Map(captifyData.Segment, func(s openrtb.Segment, _ int) Segment {
			return Segment{ID: s.ID, Name: s.Name}
		})
		for _, b := range e.creatives {
			if b.Height == height && b.Width == width && len(lo.Intersect(b.Segments, segments)) > 0 {
				return &b
			}
		}

	}
	return nil
}

func transformToOpenRTB(myBidRequest *BidRequest, segments []Segment) *openrtb.BidRequest {
	// Create a new instance of OpenRTB's BidRequest
	openRTBBidRequest := &openrtb.BidRequest{
		ID:          myBidRequest.ID,
		Impressions: transformImpressions(myBidRequest.Imp),
		Site: &openrtb.Site{
			Inventory: transformInventory(myBidRequest.Site.Domain, segments),
			Page:      myBidRequest.Site.Page,
		},
		Device: &openrtb.Device{
			UA: myBidRequest.Device.UA,
			IP: myBidRequest.Device.IP,
		},
		User: &openrtb.User{
			ID: myBidRequest.User.ID,
		},
	}

	return openRTBBidRequest
}

func transformInventory(domain string, segments []Segment) openrtb.Inventory {
	return openrtb.Inventory{
		Domain: domain,
		Content: &openrtb.Content{
			Data: transformSegments(segments),
		},
	}
}

func transformImpressions(imps []Imp) []openrtb.Impression {
	var openRTBImps []openrtb.Impression
	for _, imp := range imps {
		openRTBImp := openrtb.Impression{
			ID: imp.ID,
			Banner: &openrtb.Banner{
				Formats: transformFormats(imp.Banner.Format),
			},
			BidFloor: imp.BidFloor,
		}
		openRTBImps = append(openRTBImps, openRTBImp)
	}
	return openRTBImps
}

func transformFormats(formats []Format) []openrtb.Format {
	var openRTBFormats []openrtb.Format
	for _, format := range formats {
		openRTBFormat := openrtb.Format{
			Width:  format.W,
			Height: format.H,
		}
		openRTBFormats = append(openRTBFormats, openRTBFormat)
	}
	return openRTBFormats
}

func transformSegments(segments []Segment) []openrtb.Data {
	return []openrtb.Data{
		{
			ID:   "captify",
			Name: "captify",
			Segment: lo.Map(segments, func(s Segment, _ int) openrtb.Segment {
				return openrtb.Segment{
					ID:   s.ID,
					Name: s.Name,
				}
			}),
		},
	}
}
