
/* INITIATE MAP OBJECT */

    function initMap(div, datasources) {

    // CREATE MAP OBJ
        var map = L.map('map',{maxZoom: 22 });


    // STORE CONFIG
        map.datasources = datasources;

    // CREATE FEATURE LAYERS
        map.featureLayers = {};

        map.getLayer = function(datasource) {
            data = Utils.getRequest("/api/v1/layer/" + datasource);
            map.updateFeatureLayers(data);
        }

        map.updateFeatureLayers = function(data) {
            for (var _i in this.featureLayers){    // Remove old featurelayers
                if (this.hasLayer(this.featureLayers[_i])) {
                    this.removeLayer(this.featureLayers[_i]);
                }
            }
            try {
                this.featureLayers[$('#layers').val()] = this.createFeatureLayer(data);    // Create new featurelayers
                this.featureLayers[$('#layers').val()].addTo(this);    // Apply new featurelayers to map
            }
            catch(err) { console.log(err); }
        }

        map.createFeatureLayer = function(data) {
            var featureLayer = L.geoJson(data, {
                style: {
                    "weight": 2, 
                    "color": "#000", 
                    "fillOpacity": 0.25,
                },
                pointToLayer: function(feature, latlng) {
                    return L.circleMarker(latlng, {
                        radius: 4,
                        weight: 1,
                        fillOpacity: 0.25,
                        color: '#000'
                    });
                },
                onEachFeature: function (feature, layer) {
                    // layer.bindPopup(
                    //     "<button onclick=map.editfeature(" + JSON.stringify(feature) + ")>Edit</button>"
                    // );
                    layer.on({
                        mouseover: function(feature){
                            var properties = feature.target.feature.properties;
                            var results = "";
                            for (var item in properties) {
                                results += item + ": " + properties[item] + "<br>";
                            }
                            $("#attributes")[0].innerHTML = results;
                        },
                        mouseout: function(){
                            $("#attributes")[0].innerHTML = "Hover over features";
                        }
                    });
                }
            });
            return featureLayer;
        }

        map.getFeature = function(datasource, k){
            results = Utils.getRequest("/api/v1/layer/" + datasource + "/feature/" + k);
            return results;
        }

        // map.editfeature = function(data) {
        //     console.log(data);
        //     var layer = L.geoJson(data);
        //     layer.k = false;
        //     layer.addTo(map);
        //     layer.on('click', function(e){
        //         var layer = e.layer;
        //         layer.editing.enable();
        //         if(!layer.k){
        //             layer.k = layer.feature.properties.k;
        //             layer.on('click', function(e){
        //                 var save = confirm("save changes?");
        //                 if(save) {
        //                     $.ajax({
        //                         crossDomain: true,
        //                         dataType: 'jsonp',
        //                         async: false,
        //                         method: "PUT",
        //                         headers: {"X-HTTP-Method-Override": "PUT"},
        //                         data: {
        //                             'uuid': $('#layers').val(),
        //                             'k': this.k,
        //                             'geom': JSON.stringify([[this._latlng.lng,this._latlng.lat]])                   
        //                         },
        //                         url: '/api/v1/layer/feature',
        //                         success: function (data) {
        //                             try {
        //                                 results = data;
        //                                 alert(results.message);
        //                             }
        //                             catch(err){  console.log('Error:', err);  }
        //                         },
        //                         error: function(xhr,errmsg,err) {
        //                             console.log(xhr.status,xhr.responseText,errmsg,err);
        //                             // console.log(xhr);
        //                         }
        //                     });
        //                     map.removeLayer(layer);
        //                     map.getLayer();
        //                 }
        //             });
        //         }
        //     });
        // }



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



    // ATTRIBUTE LEGEND
        featureAttributesControl = L.control({position: 'bottomright'});
        featureAttributesControl.onAdd = function (map) {
            var div = L.DomUtil.create('div', 'info legend');
            div.innerHTML = "<h4>Attributes</h4><div id='attributes'>Hover over features</div>";
            return div;
        };
        featureAttributesControl.addTo(map);
        map.preventPropogation(featureAttributesControl);

    // GEOJSON LAYERS
        geojsonLayerControl = L.control({position: 'topright'});
        geojsonLayerControl.onAdd = function (map) {
            var div = L.DomUtil.create('div', 'info legend');
            div.innerHTML = '';
            div.innerHTML += '<i class="fa fa-search-plus" id="zoom" style="padding-left:5px; margin-right:0px;"></i><select name="basemaps" id="layers"></select>';
            return div;
        };
        geojsonLayerControl.addTo(map);
        map.preventPropogation(geojsonLayerControl);

    // LAYER CONTROL
        for (var _i=0; _i < map.datasources.length; _i++) {
            var obj = document.createElement('option');
            obj.value = map.datasources[_i];
            obj.text = map.datasources[_i];
            $('#layers').append(obj);
        }
    // Zoom to current layer
        $('#zoom').on('click', function(){ 
            map.fitBounds(
                map.featureLayers[$('#layers').val()].getBounds()
            );
        });


    // Prepare baselayers
        osm = L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png',{ 
            attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>'
        }).addTo(map);
        var baseMaps = {
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
        L.control.layers(baseMaps, overlayMaps, {position: 'topright'}).addTo(map);



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
                map.getLayer($('#layers').val());
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
                div.innerHTML = "<h4>Feature Properties</h4>";
                div.innerHTML += "<a href='#' id='add_property'><i class='fa fa-plus' style='padding-left:5px; margin-right:0px;'></i>Add Field</a><br>";
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

        var logo = L.control({position : 'topleft'});
        logo.onAdd = function () {
            this._div = L.DomUtil.create('div', 'logo');
            // this._div.innerHTML = "<div><img class='img-logo-compass' src='/images/compass.png' alt='logo'></div>"
            this._div.innerHTML = "<div><img class='img-logo-hypercube' src='/images/HyperCube2.png' alt='logo'></div>"
            return this._div;
        };
        logo.addTo(map);


    // LAUNCH MAP OBJ
        map.launch = function() {
            try {
                this.setView([0,0], 1);
                $(document).ready(function(){ 
                    $('select').on('change', function(){ 
                        map.getLayer($('#layers').val());
                    });
                });
                $(document).ready(function(){ 
                    map.fitBounds(
                        map.featureLayers[$('#layers').val()].getBounds()
                    );
                });
                map.getLayer($('#layers').val());
                map.enableDrawing();
            }
            catch(err) { console.log(err); }
        }
    // START MAP OBJ
        map.launch();
    // RETURN ENHANCED MAP OBJ
        return map;

    }



