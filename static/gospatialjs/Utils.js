
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
		return new URL(url);
	}

}
