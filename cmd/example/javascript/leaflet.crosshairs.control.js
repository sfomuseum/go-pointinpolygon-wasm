// version v1.0.0

L.Control.Crosshairs = L.Control.extend({
    _map: null,
    options: {
        position: 'topright',
	coordinates: 'latlon',
    },

    onAdd: function(map) {
	
	this._map = map;

	var self = this;
	this.draw_crosshairs();
	
	window.onresize = function(){
	    self.draw_crosshairs();
	};

	var draw_coords = function(){
	    var pos = map.getCenter();
	    var lat = pos['lat'];
	    var lon = pos['lng'];	  
	    
	    var zoom = map.getZoom();
	    
	    var ll;

	    if (self.options.coordinates == "lonlat"){
		ll = lon.toFixed(6) + ", " + lat.toFixed(6) + " #" + zoom.toFixed(2);
	    } else {	    
		ll = lat.toFixed(6) + ", " + lon.toFixed(6) + " #" + zoom.toFixed(2);
	    }
	    
	    self.div.innerText = ll;	    
	};
	
	map.on("moveend", function(){
	    draw_coords();
	});
	
	this.div = L.DomUtil.create('div','leaflet-crosshairs-container leaflet-bar');

	// This is necessary in order to copy/paste
	L.DomEvent.disableClickPropagation(this.div);
	
	draw_coords();
	return this.div;	
    },
    
    'draw_crosshairs': function(){
	
	var map_el = this._map.getContainer();
	
	var container = map_el.getBoundingClientRect();
	
	var height = container.height;
	var width = container.width;
	
	var crosshair_y = (height / 2) - 8;
	var crosshair_x = (width / 2);
	
	// http://www.sveinbjorn.org/dataurls_css
	
	var data_url = '"data:image/gif;base64,R0lGODlhEwATAKEBAAAAAP///////////' + 
		       'yH5BAEKAAIALAAAAAATABMAAAIwlI+pGhALRXuRuWopPnOj7hngEpRm6Z' + 
		       'ymAbTuC7eiitJlNHr5tmN99cNdQpIhsVIAADs="';
	
	var style = [];
	style.push("position:absolute");
	style.push("top:" + crosshair_y + "px");
	style.push("height:19px");
	style.push("width:19px");
	style.push("left:" + crosshair_x + "px");
	style.push("margin-left:-8px;");
	style.push("display:block");
	style.push("background-position: center center");
	style.push("background-repeat: no-repeat");
	style.push("background: url(" + data_url + ")");
	style.push("z-index:10000");
	
	style = style.join(";");
	
	var crosshairs = document.getElementById("crosshairs");

	if (! crosshairs){
	    crosshairs = document.createElement("div");
	    crosshairs.setAttribute("id", "crosshairs");
	    map_el.appendChild(crosshairs);
	}
	
	crosshairs.style.cssText = style;
	return true;
    },
    
    
});
