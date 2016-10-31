package app

import "github.com/paulmach/go.geojson"

// Customer structure for database
type Customer struct {
	Apikey      string      `json:"apikey"`
	Datasources []string    `json:"datasources"`
	TileLayers  []TileLayer `json:"tilelayers"`
	// TileLayers  map[string]string  `json:"tilelayers"`
}

type TileLayer struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

// MapData for html templates
type PageViewData struct {
	Apikey  string
	Version string
}

type TcpData struct {
	Apikey      string                     `json:"apikey"`
	Datasources []string                   `json:"datasources"`
	Datasource  string                     `json:"datasource"`
	Layer       *geojson.FeatureCollection `json:"layer"`
	Feature     *geojson.Feature           `json:"feature"`
}

type TcpMessage struct {
	Authkey    string  `json:"authkey"`
	Apikey     string  `json:"apikey"`
	Method     string  `json:"method"`
	Data       TcpData `json:"data"`
	Datasource string  `json:"datasource"`
}
