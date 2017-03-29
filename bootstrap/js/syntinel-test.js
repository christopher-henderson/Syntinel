function pageLoad() {
	var project = {};
	var test = {};
	var r = [];

	var projectID = getQueryVariable("project");
	var testID = getQueryVariable("test");

	var populatePage = function() {
		// Page header
		document.getElementById("header-test-name").innerHTML = test.name + " <small>Syntinel Test</small>";
		document.getElementById("breadcrumb-project-name").innerHTML = "<i class=\"fa fa-file\"></i> <a href=\"project.html?project=" + projectID + "\">" + project.name + "</a>";
		document.getElementById("breadcrumb-test-name").innerHTML = "<i class=\"fa fa-file\"></i> " + test.name;

		// Setting - Name
		document.getElementById("setting-project").value = projectID;

		// Setting - Name
		document.getElementById("setting-testName").value = test.name;

		// Setting - ID
		document.getElementById("setting-testID").value = test.id;

		// Setting - Script
		document.getElementById("setting-testScript").value = test.script;

		// Setting - Docker
		document.getElementById("setting-testDocker").value = test.dockerfile;

		// Setting - Environment Variables
		var envs = test.environmentVariables;
		var envStr = "";
		for(var i = 0; i < envs.length; i++) {
			var env = envs[i].split("=");
			envStr += "<tr><td>" + env[0] + "</td><td>" + env[1] + "</td><td>";
			envStr += "<button type=\"button\" class=\"btn btn-xs btn-info\">Edit</button>";
			envStr += "<button type=\"button\" class=\"btn btn-xs btn-danger\">Remove</button>";
			envStr += "</td></tr>";
		}
		document.getElementById("setting-environmentVariables").innerHTML = envStr;

		// Run histories
		var testRuns = document.getElementById("table-test-runs-body");
		testRuns.innerHTML = "";
		for(var i = 0; i < r.length; i++) {
			var run = r[i];
			var runRow = "";

			var runStatus;
			if(!run.successful || run.successful == null) {
				runStatus = "Running";
			} else if(run.successful == true) {
				runStatus = "Successful";
			} else {
				runStatus = "Failed";
			}

			runRow += "<tr class=\"" + (runStatus == "Successful" ? "success" : (runStatus == "Running" ? "warning" : "danger")) + "\">";
			runRow += "	<td>" + run.id + "</td>";
			runRow += "	<td>" + runStatus + "</td>";
			runRow += "</tr>";

			testRuns.innerHTML += runRow;
		}
	}

	// Make all the calls
	apiGet(SYNTINEL_URL + "/project/" + projectID, "", function(res) {
		project = res;
		project = escapeNewLineChars(project);
		project = JSON.parse(project);

		// Now get the test
		apiGet(SYNTINEL_URL + "/test/" + testID, "", function(res) {
			test = res;
			test = escapeNewLineChars(test);
			test = JSON.parse(test);

			// Get test runs
			apiGet(SYNTINEL_URL + "/testrun/all?test=" + testID, "", function(res) {
				var runs = res;
				runs = escapeNewLineChars(runs);
				runs = JSON.parse(runs);

				for(var i = 0; i < runs.length; i++) {
					var run = runs[i];
					r.push(run);
				}

				populatePage();
			});
		});
	});

	$('#table-test-runs').find('tr').click(function() {
		var index = ($(this).index());
		var row = document.getElementById("table-test-runs-body").childNodes[index];
		var id = row.childNodes[1].innerHTML;
		window.location = "run.html?project="+projectID+"&test="+testID+"&run="+id;
	});
}