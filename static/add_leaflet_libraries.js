
function add_leaflet_libraries(mapObj) {

	// Measure
    var measureControl = new L.Control.Measure();
	measureControl.addTo(mapObj);

    // Locate control
    L.control.locate().addTo(mapObj);

    // Geosearch
    new L.Control.GeoSearch({
        provider: new L.GeoSearch.Provider.OpenStreetMap()
    }).addTo(mapObj);

}
