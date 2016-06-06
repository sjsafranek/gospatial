package app

// Customer structure for database
type Customer struct {
	Apikey      string   `json:"apikey"`
	Datasources []string `json:"datasources"`
}

// MapData for html templates
type MapData struct {
	Datasource string
	Apikey     string
	Version    string
}
