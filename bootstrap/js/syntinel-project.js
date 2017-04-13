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
			"dockerfile" : 
				"# Defaults (Fill in as necessary)\n" +
				"FROM docker.io/centos\n" +
				"MAINTAINER Allstate Engineer\n" +
				"RUN yum install -y go git wget\n" +
				"# These must be the last two lines in your Dockerfile:\n" +
				"COPY script.sh $HOME/script.sh\n" +
				"CMD chmod +x script.sh && ./script.sh",
			"script" : "# Required - Input script here",
			"project" : Number(projectID),
			"interval": null
		};

		apiPost(SYNTINEL_URL + "/test/", postBody, function(res) {
			if(res.syntinelError && SYNTINEL_ERRORREDIRECT) {
				var qs = {};
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
		breadcrumbProject.innerHTML = "<i class=\"fa fa-sitemap\"></i> " + project.name + " <button id=\"button-project-delete\" type=\"button\" class=\"btn btn-xs btn-danger\">Delete Project</button>";

		var projectTests = document.getElementById("table-project-tests-body");
		projectTests.innerHTML = "";

		// Sort tests
		t.sort(function(a,b) {return (a.id > b.id) ? 1 : ((b.id > a.id) ? -1 : 0);} );
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

		document.getElementById("button-project-delete").addEventListener('click', function() {
			apiDelete(SYNTINEL_URL + "/project/" + projectID, {}, function(res) {
					if(res.syntinelError && SYNTINEL_ERRORREDIRECT) {
						var qs = {};
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
				window.location = "index.html";
			});
		});


		$('#table-project-tests').find('tr').click(function() {
			var index = ($(this).index());
			var row = document.getElementById("table-project-tests-body").childNodes[index];
			var id = row.childNodes[1].innerHTML;
			window.location = "test.html?project="+projectID+"&test="+id;
		});
	}

	// Make all the calls
	apiGet(SYNTINEL_URL + "/project/" + projectID, "", function(res) {
		if(res.syntinelError && SYNTINEL_ERRORREDIRECT) {
			var qs = {};
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
		if(project.tests.length > 0) {
			for(var j = 0; j < project.tests.length; j++) {
					apiGet(SYNTINEL_URL + "/test/" + project.tests[j], null, function(res) {
					if(res.syntinelError && SYNTINEL_ERRORREDIRECT) {
						var qs = {};
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
		} else {
			populatePage();
		}
	});
}
