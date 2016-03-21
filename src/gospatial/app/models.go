package app

/*=======================================*/
// Customer
/*=======================================*/
type Customer struct {
	Apikey      string   `json:"apikey"`
	Datasources []string `json:"datasources"`
}

/*=======================================*/
// GeoJSON
/*=======================================*/
// type PolyGeom struct {
// 	Type        string      `json:"type"`
// 	Coordinates [][]float64 `json:"coordinates"`
// }

// type LineGeom struct {
// 	Type        string      `json:"type"`
// 	Coordinates [][]float64 `json:"coordinates"`
// }

// type PointGeom struct {
// 	Type        string    `json:"type"`
// 	Coordinates []float64 `json:"coordinates"`
// }

// type Feature struct {
// 	Type       string                 `json:"type"`
// 	Geometry   interface{}            `json:"geometry"`
// 	Properties map[string]interface{} `json:"properties"`
// }

// type Geojson struct {
// 	Type     string    `json:"type"`
// 	Features []Feature `json:"features"`
// }

// func NewFeature() Feature {
// 	feat := Feature{Type: "Feature"}
// 	return feat
// }

// func NewGeojson() Geojson {
// 	geojs := Geojson{Type: "FeatureCollection"}
// 	return geojs
// }

/*=======================================*/
// Layer
/*=======================================*/
// type Layer struct {
// 	Datasource string  `json:"datasource"`
// 	Geojson    Geojson `json:"geojson"`
// }

// func (lyr *Layer) Save() error {
// 	DB.insertLayer(lyr.Datasource, lyr.Geojson)
// 	return nil
// }

/*=======================================*/
// MapData for template
/*=======================================*/
type MapData struct {
	Datasource string
	Apikey     string
}
