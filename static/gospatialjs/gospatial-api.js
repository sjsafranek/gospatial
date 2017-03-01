
/*
var GoSpatialApiHandler = Backbone.Model.extend({
	
	urlRoot: '/api/v1',
	
	defaults: {
		server: null,
		apikey: null
	},

	initialize: function(){
		console.log(this.get("server"),this.get("apikey"));
	}

});

var API = new GoSpatialApiHandler({server:"test",apikey:"1234", id:"1"});
*/


	function GoSpatialApi(apikey, server) {

		this.apikey = apikey;
		this.server = server || "";
		this.ajaxActive = 0;
		this.debug = false;
		// this.url = new URL("http://localhost:8888/map?apikey=iNx1xvBPDrZb#");

		this.getCustomer = function(callback) {
			var self = this;
			this.GET(this.server + "/api/v1/customer" + "?apikey=" + self.apikey, function(error, result){
				return callback(error, result);
			});
		}

		this.getLayer = function(datasource, callback) {
			var self = this;
			this.GET(this.server + "/api/v1/layer/" + datasource + "?apikey=" + self.apikey, function(error, result){
				return callback(error, result);
			});
		}

		this.submitFeature = function(datasource, feature, callback) {
			var self = this;
			if ("object" == typeof(feature)) {
				feature = JSON.stringify(feature);
			}
			this.POST(
				this.server + '/api/v1/layer/' + datasource + '/feature?apikey=' + self.apikey,
				feature,
				function(error, result) {
					callback(error, result);
				}
			)
		}
		
		this.editFeature = function(datasource, feature, callback) {
			var self = this;
			var geo_id = ""+feature.properties.geo_id;
			if ("object" == typeof(feature)) {
				feature = JSON.stringify(feature);
			}
			this.PUT(
				this.server + '/api/v1/layer/' + datasource + '/feature/' + geo_id + '?apikey=' + self.apikey,
				feature,
				function(error, result) {
					callback(error, result);
				}
			)
		}

		this.activeRequests = function() {
			console.log(new Date().toISOString(), "[DEBUG] API {requests:", this.ajaxActive, "}");
		}

		this._getAjaxObject = function(url, type, data, opts, callback) {
			var self = this;
			
			if(!url || !type) {
        		return false;
			}

			var ajaxObject = {
				crossDomain: true,
				url: url,
				type: type,
				dataType: opts.dataType || "json",
				beforeSend: function() {
					self.ajaxActive++;
					if (self.debug) { self.activeRequests(); }
				},
				complete: function() {
					self.ajaxActive--;
					if (self.debug) { self.activeRequests(); }
				},
				success: function (data) {
					return callback(null, data);
				},
				error: function(xhr,errmsg,err) {
					console.log(xhr.status,xhr.responseText,errmsg,err);
					var message = xhr.status + " " + xhr.responseText;
					return callback(new Error(message));
				}
			}

			if (data) {
				ajaxObject.data = data;
			}


		    ajaxObject.cache = typeof opts.cache != "undefined" ? opts.cache : false;
		    //ajaxObject.dataType = opts.dataType || "json";
		    ajaxObject.async = typeof opts.async != "undefined" ? opts.async : false;
		    
		    ajaxObject.timeout = 30000 // sets timeout to 3 seconds
		    
		    if(type.toLowerCase() == "post"){
		        ajaxObject.contentType = opts.contentType || "application/json";
		    }

		    var headers;
		    if(opts.headers){
		        headers = opts.headers;
		        delete opts.headers;
		    }
		    if(headers){
		        ajaxObject.headers = headers;
		    }

			return ajaxObject;

		}

		this.GET = function(route, callback) {
			var ajaxObj = this._getAjaxObject(route, "GET", null, {}, callback);
			$.ajax(ajaxObj);
		}

		this.POST = function(route, data, callback) {
			var ajaxObj = this._getAjaxObject(route, "POST", data, {}, callback);
			$.ajax(ajaxObj);
		}

		this.PUT = function(route, data, callback) {
			var ajaxObj = this._getAjaxObject(route, "PUT", data, {}, callback);
			console.log(ajaxObj.data);
			$.ajax(ajaxObj);
		}

	}


/*



Db4iotTurnstileApi.prototype.ajaxWithOpts = function(url, type, data, opts, onError, onSuccess){
    var self = this;
\

    var ajaxObject = opts;
    
    if(type.toLowerCase() == "get"){
        if ("object" == typeof(data)) {
			var queryString = [];
			for(var i in data){
				queryString.push(i + "=" + data[i]);
			}
			if (-1 == url.indexOf("?")) {
				url += "?" + queryString.join("&");
			} else {
				url += queryString.join("&");
			}
        //data = null;
		}
    }
    
    
    ajaxObject.success = onSuccess;
    ajaxObject.error = function(response){
		console.log(new Date().toISOString(), "[ERROR]:", "AJAX", response);
        var message = response.responseJSON || response.responseText;
        if(message && message.error_message){
            message = message.error_message;
        }else{
            message = "Unkown Error";
        }
        var code = response.status || 500;
        var error = {
            "message": message,
            "code": code
        }
    
        onError(error);
    };
    
    return $.ajax(ajaxObject);
};




 */
