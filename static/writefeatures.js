
/* INITIATE MAP OBJECT */

    function initMap(div, datasources) {

    // CREATE MAP OBJ
        var map = L.map('map',{maxZoom: 22 });

    // Prepare baselayers
        osm = L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png',{ 
            attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a> | F.I.N.D.'
        }).addTo(map);
        var baseMaps = {};
        baseMaps["OSM"] = osm;
        baseMaps["Topographic"] = L.tileLayer("http://services.arcgisonline.com/ArcGIS/rest/services/World_Topo_Map/MapServer/tile/{z}/{y}/{x}.png",{maxNativeZoom:22});
        baseMaps["Streets"] = L.tileLayer("http://services.arcgisonline.com/ArcGIS/rest/services/World_Street_Map/MapServer/tile/{z}/{y}/{x}.png",{maxNativeZoom:22});
        baseMaps["Imagery"] = L.tileLayer("http://services.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}.png",{maxNativeZoom:22});

    // BUILDINGS
        var osmb = new OSMBuildings(map)
           .date(new Date(2015, 5, 15, 17, 30))
           .load()
           .click(function(id) {
                console.log('feature id clicked:', id);
           }
        );

        var overlayMaps = { Buildings: osmb };

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


    // LAYER CONTROL
        map.layerControl = function() {
            for (var _i=0; _i < this.datasources.length; _i++) {
                var obj = document.createElement('option');
                obj.value = this.datasources[_i];
                obj.text = this.datasources[_i];
                $('#layers').append(obj);
            }
        }

    // ATTRIBUTE LEGEND
        var legend = L.control({position: 'bottomright'});
        legend.onAdd = function (map) {
            var div = L.DomUtil.create('div', 'info legend');
            div.innerHTML = "<h4>Attributes</h4><div id='attributes'>Hover over features</div>";
            return div;
        };
        legend.addTo(map);

    // GEOJSON LAYERS
        var geojsonLayers = L.control({position: 'topright'});
        geojsonLayers.onAdd = function (map) {
            var div = L.DomUtil.create('div', 'info legend');
            div.innerHTML = '';
            div.innerHTML += '<i class="fa fa-search-plus" id="zoom" style="padding-left:5px; margin-right:0px;"></i><select name="basemaps" id="layers"></select>';
            return div;
        };
        geojsonLayers.addTo(map);

        L.control.layers(baseMaps, overlayMaps, {position: 'topright'}).addTo(map);

    // LAUNCH MAP OBJ
        map.launch = function() {
            try {
                this.setView([0,0], 1);
                this.layerControl();
                $(document).ready(function(){ 
                    $('select').on('change', function(){ 
                        map.getLayer($('#layers').val());
                    });
                });
                $(document).ready(function(){ 
                    $('#zoom').on('click', function(){ 
                        map.fitBounds(
                            map.featureLayers[$('#layers').val()].getBounds()
                        );
                    });
                });
                map.getLayer($('#layers').val());
            }
            catch(err) { console.log(err); }
        }
    // START MAP OBJ
        map.launch();
    // RETURN ENHANCED MAP OBJ
        return map;

    }



