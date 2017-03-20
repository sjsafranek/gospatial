


var SwalPrettyPrint = Backbone.Model.extend({

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
			title: "<h3>" + this.get("title") + "</h3>",
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



var SwalPpError = SwalPrettyPrint.extend({

	width: 725,

	initialize: function(title, data) {
		SwalPrettyPrint.prototype.initialize.call(
			this, 
			title, 
			new Error(data).stack, 
			"error");
	}

});



var SwalPpSuccess = SwalPrettyPrint.extend({

	initialize: function(title, data) {
		SwalPrettyPrint.prototype.initialize.call(
			this, 
			title, 
			data, 
			"success");
	}

});



var SwalConfirm = Backbone.Model.extend({

	debug: false,

	defaults: {
		data: "",
		title: "Info",
		type: "info"
	},

	initialize: function(title, data, type, successCallback, cancelCallback) {
		this.set("type", type || "info");
		this.set("title", title || "Info");
		this.set("data", data);
		this.successCallback = successCallback;
		this.cancelCallback  = cancelCallback || function(dismiss) {};
		this.display();
	},
	
	text: function() {
		return JSON.stringify(this.toJSON());
	},

	display: function() {
		var self = this;

		var confirmColor = ""
		if ("info" == self.get("type")) {
			confirmColor = "#337ab7";
		}
		else if ("warning" == self.get("type")) {
			confirmColor = "#DD6B55";
		}

		swal({
			title: self.get("title"),
			text: self.get("data"),
			type: self.get("type"),
			showCancelButton: true,
			//confirmButtonColor: "#337ab7",
			confirmButtonColor: confirmColor,
			confirmButtonText: "Yes, pls!",
			cancelButtonText: "No, cancel pls!",
			showLoaderOnConfirm: true
		}).then(
			this.successCallback,
			this.cancelCallback
		);

		if (this.debug) { debugger; }
	}
	
});