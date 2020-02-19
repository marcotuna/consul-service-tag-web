package model

// Config ...
type Config struct {
	HTTP ConfigHTTP `json:"http"`
	Tag  ConfigTag  `json:"tag"`
}

// ConfigHTTP ...
type ConfigHTTP struct {
	Address string `json:"address"`
}

// ConfigTag ...
type ConfigTag struct {
	Filter string `json:"filter"`
}
