
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
			this.POST(
				this.server + '/api/v1/layer/' + datasource + '/feature?apikey=' + self.apikey,
				feature,
				function(error, result) {
					callback(error, result);
				}
			)
		}
		
		this._getAjaxObject = function(route, callback) {
			var self = this;
			return {
				crossDomain: true,
				url: route,
				dataType: 'JSON',
				beforeSend: function() {
					self.ajaxActive++;
				},
				complete: function() {
					self.ajaxActive--;
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
		}

		this.GET = function(route, callback) {
			var ajaxObj = this._getAjaxObject(route, callback);
			ajaxObject.type = "GET";
			$.ajax(ajaxObj);
		}

		this.POST = function(route, data, callback) {
			var ajaxObj = this._getAjaxObject(route, callback);
			ajaxObject.type = "POST";
			ajaxObject.data = data;
			$.ajax(ajaxObject);
		}

	}


/*



Db4iotTurnstileApi.prototype.ajaxWithOpts = function(url, type, data, opts, onError, onSuccess){
    var self = this;
    
    if (this.debug) {
		console.log(new Date().toISOString(), "[DEBUG] API {requests:", self.ajaxActive, "}");
	}

    if(!url || !type)
        return false;
    
    var headers;
    if(opts.headers){
        headers = opts.headers;
        delete opts.headers;
    }
    
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
    
    ajaxObject.url = url;
    ajaxObject.cache = typeof opts.cache != "undefined" ? opts.cache : false;
    ajaxObject.type = type;
    ajaxObject.dataType = opts.dataType || "json";
    ajaxObject.async = typeof opts.async != "undefined" ? opts.async : false;
    
    ajaxObject.timeout = 30000 // sets timeout to 3 seconds
    
    if(type.toLowerCase() == "post"){
        ajaxObject.contentType = opts.contentType || "application/json";
    }
    
    if(headers){
        ajaxObject.headers = headers;
    }
    
    if(data){
        ajaxObject.data = typeof data != "string" ? JSON.stringify(data) : data;
    } else if ("" == data) {
		ajaxObject.data = data;
	}

    ajaxObject.beforeSend = function() {
		self.ajaxActive++;
	}
    
    ajaxObject.complete = function() {
		self.ajaxActive--;
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
