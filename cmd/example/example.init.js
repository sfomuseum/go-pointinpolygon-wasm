window.addEventListener("load", function load(event){

    sfomuseum.wasm.fetch("sfomuseum_pointinpolygon.wasm").then(rsp => {
	console.log("PIP");
    }).catch((err) => {
	console.log("SAD", err);
    });
    
    
});
