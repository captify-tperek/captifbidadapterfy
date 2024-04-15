package main

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
)

type Server struct {
	ExchangesMultiplexer *ExchangesMultiplexer
}

func main() {

	config, err := ReadConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Allow all domains, or specify allowed domains
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	var exchanges []Exchange
	for _, exchangeConfig := range config.Exchanges {
		exchange := &OpenRTBExchange{}
		for _, bannerConfig := range exchangeConfig.Banners {
			exchange.creatives = append(exchange.creatives, bannerConfig)
		}
		exchanges = append(exchanges, exchange)
	}

	liveClassification := LiveClassification{config.LiveClassification.Urls}

	server := &Server{
		ExchangesMultiplexer: NewExchangesMultiplexer(
			exchanges,
			liveClassification,
		)}

	e.GET("/status", helloHandler)
	e.Static("/assets", "assets")
	e.POST("/bid", server.handleBidRequest)
	e.POST("/bidwon", server.hanldleBidWon)

	e.Logger.Fatal(e.Start(":8080"))
}

func helloHandler(c echo.Context) error {
	return c.String(http.StatusOK, "ok!")
}

func (s *Server) handleBidRequest(c echo.Context) error {
	bidRequest := new(BidRequest)
	if err := c.Bind(bidRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	bidResponse, err := s.ExchangesMultiplexer.sendToExchanges(*bidRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	} else {
		return c.JSONPretty(http.StatusOK, bidResponse, " ")
	}
}

func (s *Server) hanldleBidWon(c echo.Context) error {
	var msg map[string]any
	b, err := io.ReadAll(c.Request().Body)
	str := string(b)
	log.Infof("received bid won: %v", str)
	if err != nil {
		log.Error("error unmarshalling bid won message: %v", err)
		return c.String(http.StatusBadRequest, "Invalid input")
	}
	err = json.Unmarshal(b, &msg)
	if err != nil {
		log.Error("error unmarshalling bid won message: %v", err)
		return c.String(http.StatusBadRequest, "Invalid input")
	}
	log.Infof("registering bid won: %v", msg)
	return c.String(http.StatusOK, "ok!")

}
