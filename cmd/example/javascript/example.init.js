var map;

window.addEventListener("load", function load(event){

    // var sfo_sw = [ 37.6099139, -122.407007 ];
    // var sfo_ne = [ 37.624803, -122.370057 ];

    var sfo_sw = [ 37.601617, -122.408061 ];
    var sfo_ne = [ 37.640167, -122.354907 ];
    
    var result_el = document.getElementById("result");

    var features = {};
    var feature_layers = [];
    
    var fetch_places = function(places){

	var count = places.length;

	for (var i=0; i < count; i++){
	    fetch_place(places[i]["wof:id"]);
	}
    };

    var fetch_place = function(id){

	var url = "https://static.sfomuseum.org/geojson/" + id;
	console.log(url);

	if (features[id]){
	    add_feature(features[id]);
	    return;
	}
	
	fetch(url)
	    .then((rsp) => rsp.json())
	    .then((data) => {
		features[data["wof:id"]] = data;
		add_feature(data);
	    })
	    .catch(err => {
		console.log(err);
	    });
	

    };

    var add_feature = function(f){

	var args = {
	    fillColor: "#ff7800",
	    color: "#000",
	    weight: 1,
	    opacity: 1,
	    fillOpacity: 0.2
	};
	
	var l = L.geoJSON(f, args);
	l.addTo(map);

	feature_layers.push(l);
    };
    
    var do_pointinpolygon = function(){
	
	var pos = map.getCenter();
	console.log(pos);
	
	result_el.innerHTML = "";
	result_el.innerText = "Perform point-in-polygon query for " + pos.lat + "," + pos.lng;
	
	var req = {
	    "longitude": pos.lng,
	    "latitude": pos.lat,
	};
	
	var str_req = JSON.stringify(req);
	
	sfomuseum_pointinpolygon(str_req).then((rsp) => {

	    var count_layers = feature_layers.length;

	    for (var i=0; i < count_layers; i++){
		map.removeLayer(feature_layers[i]);
	    }

	    feature_layers = [];
	    
	    var data = JSON.parse(rsp);
	    var str_rsp = JSON.stringify(data, "", " ");
	    
	    var pre = document.createElement("pre");
	    pre.innerText = str_rsp;
	    
	    result_el.innerHTML = "";
	    result_el.appendChild(pre);

	    fetch_places(data["places"]);
	    
	}).catch((err) => {

	    result_el.innerHTML = "";
	    result_el.innerText = "Point-in-polygon query failed, " + err;
	});
	
    };
    
    var setup_map = function(){

	map = L.map('map');

	// map.setMaxBounds([ sfo_sw, sfo_ne]);
	// map.setMinZoom(16);
	
	map.fitBounds([ sfo_sw, sfo_ne ]);

	L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
	    maxZoom: 19,
	    attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
	}).addTo(map);

	var cr = new L.Control.Crosshairs({
	    coordinates: 'latlon',
	});

	cr.addTo(map);
		
	map.on("moveend", function(){
	    do_pointinpolygon();	    
	});

	do_pointinpolygon();	    	
    };
    
    sfomuseum.wasm.fetch("wasm/sfomuseum_pointinpolygon.wasm").then(rsp => {
	
	result_el.innerHTML = "";
	
	setup_map();
    }).catch((err) => {
	
	console.log("Failed to load point-in-polygon	data, "	+ err);
	
	result_el.innerText = "Failed to load point-in-polygon data, " + err;
    });
    
    
});
