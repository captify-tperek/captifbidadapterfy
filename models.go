package main

// BidRequest represents the top-level structure of the bid request JSON.
type BidRequest struct {
	ID     string `json:"id"`
	Imp    []Imp  `json:"imp"`
	Site   Site   `json:"site"`
	Device Device `json:"device"`
	User   User   `json:"user"`
}

// Imp represents an impression object in the bid request.
type Imp struct {
	ID       string  `json:"id"`
	Banner   Banner  `json:"banner"`
	BidFloor float64 `json:"bidfloor"`
}

// Banner represents the banner object in the impression.
type Banner struct {
	Format []Format `json:"format"`
}

// Format represents one of the possible formats the banner can have.
type Format struct {
	W int `json:"w"` // Width of the banner
	H int `json:"h"` // Height of the banner
}

// Site represents website information in the bid request.
type Site struct {
	Domain string `json:"domain"`
	Page   string `json:"page"`
}

// Device represents the device information from which the request originated.
type Device struct {
	UA string `json:"ua"` // User Agent
	IP string `json:"ip"` // IP address
}

// User represents a user object in the bid request.
type User struct {
	ID string `json:"id"`
}

// BidResponse represents the top-level structure of the bid response JSON.
type BidResponse struct {
	ID      string    `json:"id"`
	SeatBid []SeatBid `json:"seatbid"`
	Cur     string    `json:"cur"`
}

// SeatBid represents a seatbid object in the bid response.
type SeatBid struct {
	Bid  []Bid  `json:"bid"`
	Seat string `json:"seat"`
}

// Bid represents a single bid object within a seatbid.
type Bid struct {
	ID      string   `json:"id"`
	ImpID   string   `json:"impid"`
	Price   float64  `json:"price"`
	Adm     string   `json:"adm"`
	Crid    string   `json:"crid"`
	W       int      `json:"w"`
	H       int      `json:"h"`
	Adomain []string `json:"adomain"`
	Iurl    string   `json:"iurl"`
	Nurl    string   `json:"nurl"`
}
