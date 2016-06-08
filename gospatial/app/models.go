package app

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
