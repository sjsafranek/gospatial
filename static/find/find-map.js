
    var map,
        find,
        datasources;

/* INITIATE MAP OBJECT */

    function initialize(div, datasources) {

        map = L.map('map',{maxZoom: 22 });

        // find = L.find(datasources);
        // find.addTo(map);
        findDraw = L.find.draw(datasources);
        findDraw.addTo(map);

    // PREVENT EVENT PROPOGATION TO MAP FOR LEAFLET CONTROL ELEMENTS
        map.preventPropogation = function(obj) {
            // http://gis.stackexchange.com/questions/104507/disable-panning-dragging-on-leaflet-map-for-div-within-map
            // Disable dragging when user's cursor enters the element
            obj.getContainer().addEventListener('mouseover', function () {
                map.dragging.disable();
            });
            // Re-enable dragging when user's cursor leaves the element
            obj.getContainer().addEventListener('mouseout', function () {
                map.dragging.enable();
            });
        }


    // Prepare baselayers
        osm = L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png',{ 
            attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>'
        });
        osm.addTo(map);
        map._baseMaps = {
            "OSM": osm,
            "Topographic": L.tileLayer("http://services.arcgisonline.com/ArcGIS/rest/services/World_Topo_Map/MapServer/tile/{z}/{y}/{x}.png",{maxNativeZoom:22}),
            "Streets": L.tileLayer("http://services.arcgisonline.com/ArcGIS/rest/services/World_Street_Map/MapServer/tile/{z}/{y}/{x}.png",{maxNativeZoom:22}),
            "Imagery": L.tileLayer("http://services.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}.png",{maxNativeZoom:22})
        };

    // BUILDINGS
        var osmb = new OSMBuildings(map)
           .date(new Date(2015, 5, 15, 17, 30))
           .load()
           .click(function(id) {
                console.log('feature id clicked:', id);
           }
        );

        var overlayMaps = { Buildings: osmb };

    // Baselayers
        L.control.layers(map._baseMaps, overlayMaps, {position: 'topright'}).addTo(map);


    // Drawing
        map.enableDrawing = function() {
            if (this.drawing) {
                return;
            }
            this.drawing = {};
            // Add controls
            var measureControl = new L.Control.Measure();
            measureControl.addTo(this);
            L.control.locate().addTo(this);
            new L.Control.GeoSearch({
                provider: new L.GeoSearch.Provider.OpenStreetMap()
            }).addTo(this);

            // Draw Features
            this.drawing.drawnItems = L.featureGroup().addTo(this);
            this.addControl(new L.Control.Draw({
                draw: { circle: false },
                edit: { featureGroup: this.drawing.drawnItems }
            }));

            this.drawing._feature_types = {
                "marker": "Point",
                "polygon": "Polygon",
                "rectangle": "Polygon",
                "polyline": "LineString"
            };

            this.drawing.sendFeature = function(id) {
                var results;
                var feature = this.drawnItems._layers[id];
                var payload = {
                    "geometry": {
                        "type": this._feature_types[feature.layerType],
                        "coordinates": []
                    },
                    "properties": this._getProperties()
                }
                if (payload.geometry.type == "Point") {
                    payload.geometry.coordinates.push(feature._latlng.lng);
                    payload.geometry.coordinates.push(feature._latlng.lat);
                } else if (payload.geometry.type == "LineString") {
                    for (var i = 0; i < feature._latlngs.length; i++) {
                        payload.geometry.coordinates.push([feature._latlngs[i].lng,feature._latlngs[i].lat])
                    }
                } else if (payload.geometry.type == "Polygon") {
                    payload.geometry.coordinates.push([]);
                    for (var i = 0; i < feature._latlngs.length; i++) {
                        payload.geometry.coordinates[0].push([feature._latlngs[i].lng,feature._latlngs[i].lat])
                    }
                } else {
                    alert("Unknown feature type!")
                    return results
                }
                results = Utils.postRequest(
                    '/api/v1/layer/' + $('#layers').val() + '/feature',
                    JSON.stringify(payload)
                );
                console.log(this.$super);
                find.getLayer($('#layers').val());
                map.removeLayer(this.drawnItems._layers[id]);
                $("#properties .attr").val("");
                return results;
            }

            var popup = L.popup();
            function onMapClick(e) {
                if (e.target.editing._enabled) {  console.log('editing enabled')  }
                else {
                    popup
                        .setLatLng(e.latlng)
                        .setContent("<div class='button' value='Submit Feature' onClick='map.drawing.sendFeature(" + e.target._leaflet_id + ")'><h4>Submit Feature</h4><div>")
                        .openOn(map);
                }
            }

            this.on('draw:created', function(event) {
                var layer = event.layer;
                layer.on('click', onMapClick);
                layer.options.color='blue';
                layer.layerType = event.layerType;
                this.drawing.drawnItems.addLayer(layer);
            });


            featurePropertiesControl = L.control({position: 'bottomleft'});
            featurePropertiesControl.onAdd = function () {
                var div = L.DomUtil.create('div', 'info legend properties_form');
                div.innerHTML = "<div>";
                div.innerHTML += "<strong>Feature Properties </strong>";
                // div.innerHTML += "<a href='#' id='add_property'><i class='fa fa-plus'></i>Add Field</a>";
                div.innerHTML += "<a href='#' id='add_property'>[Add Field]</a>";
                div.innerHTML += "</div>";
                div.innerHTML += "<div id='properties'>";
                div.innerHTML += "</div>";
                return div;
            };
            featurePropertiesControl.addTo(this);
            this.preventPropogation(featurePropertiesControl);

            $("#add_property").on("click", function() {
                $("#properties").append("<input type='text' class='field' placeholder='field'><input type='text' class='attr' placeholder='attribute'><br>");
            });

            this.drawing._getProperties = function() {
                var prop = {};
                var fields = $("#properties .field");
                var attrs = $("#properties .attr");
                for (var _i=0; _i < fields.length; _i++) {
                    prop[fields[_i].value] = attrs[_i].value;
                }
                return prop;
            }
        }


        // map.setView([0,0], 1);
        map.enableDrawing();

    }



