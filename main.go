package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
)

func main() {

	config, err := ReadConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:                             []string{"*"}, // Allow all domains, or specify allowed domains
		AllowMethods:                             []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders:                             []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		UnsafeWildcardOriginWithAllowCredentials: true,
		AllowCredentials:                         true,
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

	server := NewServer(NewExchangesMultiplexer(
		exchanges,
		liveClassification,
	))

	e.GET("/status", helloHandler)
	e.Static("/static", "static")
	e.POST("/bid", server.handleBidRequest)
	e.GET("/bidwon", server.handleBidWon)
	e.GET("/stats", server.handleStats)

	e.Logger.Fatal(e.Start(":8080"))
}

func helloHandler(c echo.Context) error {
	return c.String(http.StatusOK, "ok!")
}
