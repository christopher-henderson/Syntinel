function pageLoad() {
	var status = getQueryVariable("status");
	var project = getQueryVariable("project");
	var test = getQueryVariable("test");
	var run = getQueryVariable("run");
	var reason = getQueryVariable("reason");

	if(status) {
		var statusStr = ""
		if(Number(status) >= 100 && Number(status) <= 900) {
			statusStr = "Status code: " + status;
		} else {
			statusStr = status;
		}
		document.getElementById("error-header").innerHTML = "<h3>Uh oh! Something went wrong (" + statusStr + "):</h3>";
	}

	var body = document.getElementById("error-body");
	body.innerHTML = "<p>You can try to:</p>";

	if(run) {
		body.innerHTML += "- <a href=\"run.html?project="+project+"&test="+test+"&run="+run+"\">Reload the <strong>run</strong> page</a><br>";
	}
	if(test) {
		body.innerHTML += "- <a href=\"test.html?project="+project+"&test="+test+"\">Reload the <strong>test</strong> page</a><br>";
	}
	if(project) {
		body.innerHTML += "- <a href=\"project.html?project="+project+"\">Reload the <strong>project</strong> page</a><br>";
	}

	body.innerHTML += "- <a href=\"index.html\">Reload the <strong>dashboard</strong> page</a><br>";

	if(reason && reason.length > 0) {
		body.innerHTML += "<br><br><h4>Here's what else we know:</h4>";
		body.innerHTML += "<p>" + decodeURI(reason).escape() + "</p>";
	}
}