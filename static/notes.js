




/*
        map.heatLayers = {};

        map.createHeatLayer = function(data) {
            var pts = [];
            for (_i=0; _i < data.length; _i++) {
                var pt = [parseFloat(data[_i][0]),parseFloat(data[_i][1])];
                pts.push(pt);
            }
            var featureLayer = L.heatLayer(pts);
            //featureLayer.setOptions({radius: 15, blur: 20, maxZoom: 25});
            //featureLayer.setOptions({maxZoom: 12});
            return featureLayer;
        }

        map.getHeatLayer = function() {
            $.get(
                "http://" + baseURL + "/api/v1/layer/heat", 
                {'uuid':$('#layers').val()},
                function(data,status){
                    if (status == 'success') {
                        var layer = map.updateHeatLayers(data);
                    }
                    else {
                        console.log(status);
                    }
                },
                "json"
            );
        }
        map.updateHeatLayers = function(data) {
            for (var _i in this.heatLayers){    // Remove old featurelayers
                if (this.hasLayer(this.heatLayers[_i])) {
                    this.removeLayer(this.heatLayers[_i]);
                }
            }
            this.heatLayers = {};    // Clear old featurelayers
            try {
                this.heatLayers[data.uuid] = this.createHeatLayer(data.geodata);    // Create new featurelayers
                this.heatLayers[data.uuid].addTo(this);    // Apply new featurelayers to map
            }
            catch(err) { console.log(err); }
        }



                    "baselayers": {
                        "Topographic": "http://services.arcgisonline.com/ArcGIS/rest/services/World_Topo_Map/MapServer/tile/{z}/{y}/{x}.png",
                        "Streets": "http://services.arcgisonline.com/ArcGIS/rest/services/World_Street_Map/MapServer/tile/{z}/{y}/{x}.png",
                        "Oceans": "http://services.arcgisonline.com/arcgis/rest/services/Ocean/World_Ocean_Base/MapServer/tile/{z}/{y}/{x}.png",
                        "NationalGeographic": "http://services.arcgisonline.com/ArcGIS/rest/services/NatGeo_World_Map/MapServer/tile/{z}/{y}/{x}.png",
                        "Gray": "http://services.arcgisonline.com/ArcGIS/rest/services/Canvas/World_Light_Gray_Base/MapServer/tile/{z}/{y}/{x}.png",
                        "DarkGray": "http://services.arcgisonline.com/ArcGIS/rest/services/Canvas/World_Dark_Gray_Base/MapServer/tile/{z}/{y}/{x}.png",
                        "Imagery": "http://services.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}.png",
                        "ShadedRelief": "http://services.arcgisonline.com/ArcGIS/rest/services/World_Shaded_Relief/MapServer/tile/{z}/{y}/{x}.png",
                        "Terrain": "http://services.arcgisonline.com/ArcGIS/rest/services/World_Terrain_Base/MapServer/tile/{z}/{y}/{x}.png",
                        "World at Night": "https://tiles2.arcgis.com/tiles/P3ePLMYs2RVChkJx/arcgis/rest/services/Earth_at_Night_WM/MapServer/tile/{z}/{y}/{x}"
                    },
                    "overlays": {
                        "None": null,
                        "World Transportation":"https://services.arcgisonline.com/ArcGIS/rest/services/Reference/World_Transportation/MapServer/tile/{z}/{y}/{x}",
                        "LandScan World 2010 Population":"https://utility.arcgis.com/usrsvcs/rest/services/5864fa4352d34c859bd5fa4c0344f500/MapServer/tile/{z}/{y}/{x}"
                    //    ,
                    //    "tilestache": config.tilestache
                    },
                    "labels": {
                        "None": null,
                        "OceansLabels": "http://services.arcgisonline.com/arcgis/rest/services/Ocean/World_Ocean_Reference/MapServer/tile/{z}/{y}/{x}.png",
                        "GrayLabels": "http://services.arcgisonline.com/ArcGIS/rest/services/Canvas/World_Light_Gray_Reference/MapServer/tile/{z}/{y}/{x}.png",
                        "DarkGrayLabels": "http://services.arcgisonline.com/ArcGIS/rest/services/Canvas/World_Dark_Gray_Reference/MapServer/tile/{z}/{y}/{x}.png",
                        "ImageryLabels": "http://services.arcgisonline.com/ArcGIS/rest/services/Reference/World_Boundaries_and_Places/MapServer/tile/{z}/{y}/{x}.png",
                        "ShadedReliefLabels": "http://services.arcgisonline.com/ArcGIS/rest/services/Reference/World_Boundaries_and_Places_Alternate/MapServer/tile/{z}/{y}/{x}.png",
                        "TerrainLabels": "http://services.arcgisonline.com/ArcGIS/rest/services/Reference/World_Reference_Overlay/MapServer/tile/{z}/{y}/{x}.png"
                    }

*/
