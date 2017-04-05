function pageLoad() {
	var project = {};
	var t = [];

	var projectID = getQueryVariable("project");

	document.getElementById("button-test-add").addEventListener('click', function() {
		// Open the modal
		$("#modal-add").modal();
	});

	var createProject = document.getElementById("modal-add-test").addEventListener('click', function() {
		var postBody = {
			"name" : document.getElementById("modal-add-test-name").value,
			"environmentVariables" : null,
			"dockerfile" : "Required",
			"script" : "Required",
			"interval": null
		};

		apiPost(SYNTINEL_URL + "/test/", postBody, function(res) {
			if(res.error && SYNTINEL_ERRORREDIRECT) {
				if(res.responseText && res.responseText.length > 0) {
					qs.reason = res.responseText;
				}
				if(res.status) {
					qs.status = res.status;
				}
				
				qs.project = projectID;

				window.location = buildUrl("error.html", qs);
				return;
			}

			var test = res;
			test = escapeNewLineChars(test);
			test = JSON.parse(test);

			window.location = "test.html?project="+projectID+"&test="+test.id;
		});
	});

	var populatePage = function() {
		var pageHeader = document.getElementById("header-project-name");
		pageHeader.innerHTML = project.name + " <small>Syntinel Project</small>";

		var breadcrumbProject = document.getElementById("breadcrumb-project-name");
		breadcrumbProject.innerHTML = "<i class=\"fa fa-sitemap\"></i> " + project.name;

		var projectTests = document.getElementById("table-project-tests-body");
		projectTests.innerHTML = "";
		for(var i = 0; i < t.length; i++) {
			var test = t[i];
			var testRow = "";

			testRow += "<tr class=\"" + (test.health >= SYNTINEL_HEALTH.SUCCESS_MIN ? "success" : (test.health > SYNTINEL_HEALTH.WARN_MIN ? "warning" : "danger")) + "\">";
			testRow += "	<td>" + test.id + "</td>";
			testRow += "	<td>" + test.name + "</td>";
			testRow += "	<td>" + test.health + "%</td>";
			testRow += "</tr>";

			projectTests.innerHTML += testRow;
		}
	}

	// Make all the calls
	apiGet(SYNTINEL_URL + "/project/" + projectID, "", function(res) {
		if(res.error && SYNTINEL_ERRORREDIRECT) {
			if(res.responseText && res.responseText.length > 0) {
				qs.reason = res.responseText;
			}
			if(res.status) {
				qs.status = res.status;
			}
			
			qs.project = projectID;

			window.location = buildUrl("error.html", qs);
			return;
		}

		project = res;
		project = escapeNewLineChars(project);
		project = JSON.parse(project);

		var count = 0;
		for(var j = 0; j < project.tests.length; j++) {
				apiGet("/test/" + projectID, "", function(res) {
				if(res.error && SYNTINEL_ERRORREDIRECT) {
					if(res.responseText && res.responseText.length > 0) {
						qs.reason = res.responseText;
					}
					if(res.status) {
						qs.status = res.status;
					}
					
					qs.project = projectID;

					window.location = buildUrl("error.html", qs);
					return;
				}

				t.push(JSON.parse(escapeNewLineChars(res)));
				count++;
				if(count == project.tests.length) {
					populatePage();
				}
			});
		}
	});

	$('#table-project-tests').find('tr').click(function() {
		var index = ($(this).index());
		var row = document.getElementById("table-project-tests-body").childNodes[index];
		var id = row.childNodes[1].innerHTML;
		window.location = "test.html?project="+projectID+"&test="+id;
	});
}
