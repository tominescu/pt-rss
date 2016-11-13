package config

import "encoding/json"

type Config struct {
	Timeout     int    `json:"timeout"`
	SettingsDir string `json:"settings_dir"`
	Sites       []Site `json:"sites"`
}

type Site struct {
	Name        string `json:"name"`
	Rss         string `json:"rss"`
	DownloadDir string `json:"download_dir"`
	Interval    int    `json:"interval"`
}

func NewConfig(b []byte) (*Config, error) {
	c := Config{}
	if err := json.Unmarshal(b, &c); err != nil {
		return &c, err
	}
	return &c, nil
}
