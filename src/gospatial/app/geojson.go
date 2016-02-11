package app

type PolyGeom struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

type LineGeom struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

type PointGeom struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Feature struct {
	Type       string                 `json:"type"`
	Geometry   interface{}            `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type Geojson struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

func NewFeature() Feature {
	feat := Feature{Type: "Feature"}
	return feat
}

func NewGeojson() Geojson {
	geojs := Geojson{Type: "FeatureCollection"}
	return geojs
}
