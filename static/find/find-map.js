
    var map,
        find,
        findDraw,
        datasources;

    function initialize(div, datasources) {

        map = L.map('map',{maxZoom: 22 });

        // find = L.find(datasources);
        // find.addTo(map);
        findDraw = L.find.draw(datasources);
        findDraw.addTo(map);

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

        var osmb = new OSMBuildings(map)
           .date(new Date(2015, 5, 15, 17, 30))
           .load()
           .click(function(id) {
                console.log('feature id clicked:', id);
           }
        );

        var overlayMaps = { Buildings: osmb };

        L.control.layers(map._baseMaps, overlayMaps, {position: 'topright'}).addTo(map);

    }
