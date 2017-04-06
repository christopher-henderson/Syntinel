var SYNTINEL_URL = "http://syntinel.chenderson.org/api/v1";
var SYNTINEL_HEALTH = {
	SUCCESS_MIN : 87.5,
	WARN_MIN : 75
}
var SYNTINEL_ERRORREDIRECT = true;

function getQueryVariable(variable) {
	var query = window.location.search.substring(1);
	var vars = query.split("&");
	for (var i=0;i<vars.length;i++) {
		var pair = vars[i].split("=");
		if(pair[0] == variable){return pair[1];}
	}
}

function escapeNewLineChars(valueToEscape) {
	if (valueToEscape != null && valueToEscape != "") {
		valueToEscape = valueToEscape.replace(/\r/g, "\\r")
		return valueToEscape.replace(/\n/g, "\\n");
	} else {
		return valueToEscape;
	}
}

function apiGet(url, params, callback) {
	console.log("Requesting " + url);
	var handler = function(request) {
		callback(request.response);
	};

	var request = new XMLHttpRequest();
	request.open("GET", url, true);
	request.withCredentials = true;
	request.setRequestHeader("Content-Type","application/json");
	request.send();
	request.onreadystatechange = function() {
		if((request.status >= 200 && request.status < 300) && request.readyState == 4) { // Should be 200
			handler(request);
		} else if (request.readyState == 4) {
			console.log("Response code " + request.status);
			console.log(request.responseText);
			callback({"syntinelError" : true, "status" : request.status, "responseText" : request.responseText});
		}
	};
}

function apiPost(url, body, callback) {
	console.log("Requesting " + url);
	var handler = function(request) {
		callback(request.response);
	};

	var request = new XMLHttpRequest();
	request.open("POST", url, true);
	request.withCredentials = true;
	request.setRequestHeader("Content-Type","application/json");
	request.send(JSON.stringify(body));
	request.onreadystatechange = function() {
		if((request.status >= 200 && request.status < 300) && request.readyState == 4) { // Should be 201
			handler(request);
		} else if (request.readyState == 4) {
			console.log("Response code " + request.status);
			console.log(request.responseText);
			callback({"syntinelError" : true, "status" : request.status, "responseText" : request.responseText});
		}
	};
}

function apiPatch(url, body, callback) {
	console.log("Requesting " + url);
	var handler = function(request) {
		callback(request.response);
	};

	var request = new XMLHttpRequest();
	request.open("PATCH", url, true);
	request.withCredentials = true;
	request.setRequestHeader("Content-Type","application/json");
	request.send(JSON.stringify(body));
	request.onreadystatechange = function() {
		if((request.status >= 200 && request.status < 300) && request.readyState == 4) {
			handler(request);
		} else if (request.readyState == 4) {
			console.log("Response code " + request.status);
			console.log(request.responseText);
			callback({"syntinelError" : true, "status" : request.status, "responseText" : request.responseText});
		}
	};
}

function apiDelete(url, body, callback) {
	console.log("Requesting " + url);
	var handler = function(request) {
		callback(request.response);
	};

	var request = new XMLHttpRequest();
	request.open("DELETE", url, true);
	request.withCredentials = true;
	request.setRequestHeader("Content-Type","application/json");
	request.send(JSON.stringify(body));
	request.onreadystatechange = function() {
		if((request.status >= 200 && request.status < 300) && request.readyState == 4) {
			handler(request);
		} else if (request.readyState == 4) {
			console.log("Response code " + request.status);
			console.log(request.responseText);
			callback({"syntinelError" : true, "status" : request.status, "responseText" : request.responseText});
		}
	};
}

String.prototype.replaceAll = function(search, replacement) {
  var target = this;
  return target.split(search).join(replacement);
};

String.prototype.escape = function() {
    var tagsToReplace = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;'
    };
    return this.replace(/[&<>]/g, function(tag) {
        return tagsToReplace[tag] || tag;
    });
};

function getRandomIntInclusive(min, max) {
  min = Math.ceil(min);
  max = Math.floor(max);
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

function buildUrl(url, parameters) {
  var qs = "";
  for(var key in parameters) {
    var value = parameters[key];
    qs += encodeURIComponent(key) + "=" + encodeURIComponent(value) + "&";
  }
  if (qs.length > 0){
    qs = qs.substring(0, qs.length-1); //chop off last "&"
    url = url + "?" + qs;
  }
  return url;
}

function getTimestamp(date) {
  if(!date || typeof date != "object") {
    date = new Date();
  }
  
  var YYYY = date.getFullYear();
  var MM = date.getMonth()+1;
  var DD = date.getDate();
  var hh = date.getHours();
  var mm = date.getMinutes();
  
  if(DD < 10) {
    DD='0'+DD;
  } 
  if(MM < 10) {
    MM='0'+MM;
  } 
  if(hh < 10) {
    hh='0'+hh;
  } 
  if(mm < 10) {
    mm='0'+mm;
  } 
  
  return (YYYY+"-"+MM+"-"+DD+" @ "+hh+""+mm);
}