var SYNTINEL_URL = "http://syntinel.chenderson.org/api/v1";
var SYNTINEL_HEALTH = {
	SUCCESS_MIN : 87.5,
	WARN_MIN : 75
}

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
	var handler = function(request) {
		callback(request.response);
	};

	var request = new XMLHttpRequest();
	request.open("GET", url, true);
	request.withCredentials = true;
	request.setRequestHeader("Content-Type","application/json");
	request.send();
	request.onreadystatechange = function() {
		if(request.status === 200 && request.readyState == 4) {
			handler(request);
		} else if (request.readyState == 4) {
			console.log("Response code " + request.status);
			console.log(request.responseText);
		}
	};
}

function apiPost(url, body, callback) {
	var handler = function(request) {
		callback(request.response);
	};

	var request = new XMLHttpRequest();
	request.open("POST", url, true);
	request.withCredentials = true;
	request.setRequestHeader("Content-Type","application/json");
	request.send(JSON.stringify(body));
	request.onreadystatechange = function() {
		if(request.status === 201 && request.readyState == 4) {
			handler(request);
		} else if (request.readyState == 4) {
			console.log("Response code " + request.status);
			console.log(request.responseText);
		}
	};
}

function apiPatch(url, body, callback) {
	var handler = function(request) {
		callback(request.response);
	};

	var request = new XMLHttpRequest();
	request.open("PATCH", url, true);
	request.withCredentials = true;
	request.setRequestHeader("Content-Type","application/json");
	request.send(JSON.stringify(body));
	request.onreadystatechange = function() {
		if(request.status === 200 && request.readyState == 4) {
			handler(request);
		} else if (request.readyState == 4) {
			console.log("Response code " + request.status);
			console.log(request.responseText);
		}
	};
}

function apiDelete(url, body, callback) {
	var handler = function(request) {
		callback(request.response);
	};

	var request = new XMLHttpRequest();
	request.open("DELETE", url, true);
	request.withCredentials = true;
	request.setRequestHeader("Content-Type","application/json");
	request.send(JSON.stringify(body));
	request.onreadystatechange = function() {
		if(request.status === 200 && request.readyState == 4) {
			handler(request);
		} else if (request.readyState == 4) {
			console.log("Response code " + request.status);
			console.log(request.responseText);
		}
	};
}

String.prototype.replaceAll = function(search, replacement) {
  var target = this;
  return target.split(search).join(replacement);
};

function getRandomIntInclusive(min, max) {
  min = Math.ceil(min);
  max = Math.floor(max);
  return Math.floor(Math.random() * (max - min + 1)) + min;
}
