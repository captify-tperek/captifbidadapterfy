package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
)

type Server struct {
	ExchangesMultiplexer *ExchangesMultiplexer
	creativeStats        map[string]float64
}

func NewServer(multiplexer *ExchangesMultiplexer) *Server {
	return &Server{
		ExchangesMultiplexer: multiplexer,
		creativeStats:        make(map[string]float64),
	}

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

func (s *Server) handleBidWon(c echo.Context) error {
	creativeId := c.QueryParam("creativeId")
	if creativeId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input, missing creativeId")
	}
	winningCPMstr := c.QueryParam("cpm")
	if winningCPMstr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input, missing cpm")
	}
	winningCPM, err := strconv.ParseFloat(winningCPMstr, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input, unable to parse cpm")
	}

	log.Infof("registering bid won with creative: %s, cpm: %d",
		creativeId,
		winningCPM)

	if _, ok := s.creativeStats[creativeId]; !ok {
		s.creativeStats[creativeId] = winningCPM / 1000.0
	} else {
		s.creativeStats[creativeId] += winningCPM / 1000.0
	}

	return c.String(http.StatusOK, "ok!")
}

func (s *Server) handleStats(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, s.creativeStats, " ")
}
