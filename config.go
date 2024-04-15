package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Segment struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

type LiveClassificationConfig struct {
	Urls map[string][]Segment `yaml:"urls"`
}

type BannerConfig struct {
	Width      int       `yaml:"width"`
	Height     int       `yaml:"height"`
	Segments   []Segment `yaml:"segments"`
	Price      float64   `yaml:"price"`
	AdMarkup   string    `yaml:"ad_markup"`
	CreativeID string    `yaml:"creative_id"`
	AdvDomains []string  `yaml:"adv_domains"`
	ImageURL   string    `yaml:"image_url"`
	NoticeURL  string    `yaml:"notice_url"`
}

type ExchangeConfig struct {
	Name    string         `yaml:"name"`
	Banners []BannerConfig `yaml:"banners"`
}

type Config struct {
	Exchanges          []ExchangeConfig         `yaml:"exchanges"`
	LiveClassification LiveClassificationConfig `yaml:"live_classification"`
}

func ReadConfig() (*Config, error) {
	f, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
