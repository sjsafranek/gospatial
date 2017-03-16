
var SwalJSON = Backbone.Model.extend({

	debug: false,
	width: 500,

	defaults: {
		data: "",
		title: "Success",
		type: "success"
	},

	initialize: function(title, data, type) {
		this.set("type", type || "success");
		this.set("title", title || "Success");
		this.set("data", data);
		this.display();
	},
	
	text: function() {
		return JSON.stringify(this.toJSON());
	},

	// http://stackoverflow.com/questions/4810841/how-can-i-pretty-print-json-using-javascript
	syntaxHighlight: function(json) {
		if (typeof json != 'string') {
			 json = JSON.stringify(json, undefined, 2);
		}
		json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
		return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
			var cls = 'number';
			if (/^"/.test(match)) {
				if (/:$/.test(match)) {
					cls = 'key';
				} else {
					cls = 'string';
				}
			} else if (/true|false/.test(match)) {
				cls = 'boolean';
			} else if (/null/.test(match)) {
				cls = 'null';
			}
			return '<span class="' + cls + '">' + match + '</span>';
		});
	},

	display: function() {
		var self = this;
		swal({
			title: this.get("title"),
			html: (function() {
					return	 "<pre class='well'>"
							+	"<small>"
							+		self.syntaxHighlight(
										self.get("data")
									)
							+	"</small>"
							+ "</pre>";

				  })(),
			width: this.width || 500,
			type: this.get("type")
		});
		if (this.debug) { debugger; }
	}
	
});



var SwalError = SwalJSON.extend({

	width: 725,

	initialize: function(title, data) {
		SwalJSON.prototype.initialize.call(
			this, 
			title, 
			new Error(data).stack, 
			"error");
	}

});


