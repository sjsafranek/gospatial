
function Utils() {
	/** 
	 * method:     randomColor()
	 * desciption: Generates and returns random hex color
	 * @returns    hex color code
	 */
	this.randomColor = function() {
		return '#'+Math.floor(Math.random()*16777215).toString(16);
	},

	/** 
	 * method:     uuid()
	 * desciption: Generates and returns randomly generated uuid string
	 * @returns    uuid string
	 */
	this.uuid = function() {
		function s4() {
			return Math.floor((1 + Math.random()) * 0x10000)
				.toString(16)
				.substring(1);
		}
		return s4() + s4() + '-' + s4() + '-' + s4() + '-' + s4() + '-' + s4() + s4() + s4();
	},

	/** 
	 * method:     parseUrl()
	 * desciption: parses url into pieces
	 * @returns    map of url parts
	 */
	this.parseUrl = function(url) {
		// new URL.searchParams.get('apikey')
		return new URL(url);
	}

}


// Source: Stacked Overflow
// http://stackoverflow.com/questions/1960473/unique-values-in-an-array
Array.prototype.getUnique = function(){
	var u = {}, a = [];
	for(var i = 0, l = this.length; i < l; ++i){
		if(u.hasOwnProperty(this[i])) {
			continue;
		}
		a.push(this[i]);
		u[this[i]] = 1;
	}
	return a;
}

Array.prototype.getMin = function() {
	var n = null;
	for (var i=0; i<this.length; i++) {
		if ("number" == typeof(this[i])) {
			if (null == n) {
				n = this[i];
			}
			if (n > this[i]) {
				n = this[i];
			}
		}
	}
	return n;
}

Array.prototype.getMax = function() {
	var n = null;
	for (var i=0; i<this.length; i++) {
		if ("number" == typeof(this[i])) {
			if (null == n) {
				n = this[i];
			}
			if (n < this[i]) {
				n = this[i];
			}
		}
	}
	return n;
}
