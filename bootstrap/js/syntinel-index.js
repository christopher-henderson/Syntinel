function pageLoad() {
	apiGet(SYNTINEL_URL + "/project/all", "", function(res) {
		console.log(res);
	});

	var projects = "[{\"id\": 1,\"tests\": [1],\"name\": \"UltimateCode\"}]"
	projects = escapeNewLineChars(projects);
	projects = JSON.parse(projects);

	var tabAll = document.getElementById("projectsTab-all");
	for(var i = 0; i < projects.length; i++) {
		var project = projects[i];

		var tests = [];
		for(var j = 0; j < project.tests.length; j++) {
			// TODO API call to get test
		}

		var tab = "";
		tab += "<div aria-multiselectable=\"true\" class=\"panel-group\" id=\"accordion\" role=\"tablist\">";
		tab += "	<div class=\"panel panel-success\">";
		tab += "		<div class=\"panel-heading prev-emp activestate\" id=\"rc_branch\" role=\"tab\">";
		tab += "			<h4 class=\"panel-title\"><a aria-controls=\"collapseOne\" aria-expanded=\"true\" data-parent=\"#accordion\" data-toggle=\"collapse\" href=\"#collapseOne\">" + project.name + "</a>";
		tab += "			</h4>";
		tab += "		</div>";
		tab += "		<div aria-labelledby=\"rc_branch\" class=\"panel-collapse collapse\" id=\"collapseOne\" role=\"tabpanel\">";
		tab += "			<div class=\"container panel-body\">";
		tab += "				***** SOME TEXT *****";
		tab += "			</div>";
		tab += "		</div>";
		tab += "	</div>";
		tab += "</div>";

		tabAll.innerHTML += tab;
	}
}