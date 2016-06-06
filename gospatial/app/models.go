package app

// Customer
type Customer struct {
	Apikey      string   `json:"apikey"`
	Datasources []string `json:"datasources"`
}

// MapData for template
type MapData struct {
	Datasource string
	Apikey     string
	Version    string
}
