function pageLoad() {
	var project = {};
	var test = {};
	var run = {};

	var projectID = getQueryVariable("project");
	var testID = getQueryVariable("test");
	var runID = getQueryVariable("run");

	var populatePage = function() {
		// Page header
		document.getElementById("header-run-name").innerHTML = "Run #" + runID + " <small>Syntinel Test Run</small>";
		document.getElementById("breadcrumb-project-name").innerHTML = "<i class=\"fa fa-sitemap\"></i> <a href=\"project.html?project=" + projectID + "\">" + project.name + "</a>";
		document.getElementById("breadcrumb-test-name").innerHTML = "<i class=\"fa fa-file\"></i> <a href=\"test.html?project=" + projectID + "&test=" + testID + "\">" + test.name + "</a>";
		document.getElementById("breadcrumb-run-name").innerHTML = "<i class=\"fa fa-cog\"></i> Run #" + runID;

		var runStatus;
		if(!run.successful || run.successful == null) {
			runStatus = "Still running";
		} else if(run.successful == true) {
			runStatus = "Successfully ran";
		} else {
			runStatus = "Failed to run";
		}

		var runStatusHeader = document.getElementById("run-status-header");
		runStatusHeader.innerHTML = "<strong>" + runStatus + "</strong>" + " at " + (run.timestamp ? run.timestamp : "unknown time");
		runStatusHeader.className = "alert " + (runStatus == "Successfully ran" ? "alert-success" : (runStatus == "Still running" ? "alert-warning" : "alert-danger"));

		if(run.log) {
			var c = document.getElementById("run-console");
			c.innerHTML += run.log + "\n";
		}

		startWebsocket(run.id);
	}

	// Get the project
	apiGet(SYNTINEL_URL + "/project/" + projectID, "", function(res) {
		project = res;
		project = escapeNewLineChars(project);
		project = JSON.parse(project);

		// Now get the test
		apiGet(SYNTINEL_URL + "/test/" + testID, "", function(res) {
			test = res;
			test = escapeNewLineChars(test);
			test = JSON.parse(test);

			// Get the test
			apiGet(SYNTINEL_URL + "/testrun/" + runID, "", function(res) {
				run = res;
				run = escapeNewLineChars(run);
				run = JSON.parse(run);

				populatePage();
			});
		});
	});

    function startWebsocket(id) {
		var c = document.getElementById("run-console");
		socket = new WebSocket("ws://" + document.domain + "/testRun/console/" + id);
		socket.onmessage = function(e) {
			c.innerHTML += e.data + "\n";
		}
		socket.onopen = function() {
			console.log("Connection established - Showing Run #" + run.id + " console");
		}
		// Call onopen directly if socket is already open
		if (socket.readyState == WebSocket.OPEN) socket.onopen();
    }
}