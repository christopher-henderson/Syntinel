function pageLoad() {
	var projectID = getQueryVariable("project");

	var add = document.getElementById("button-test-add");
	add.addEventListener('click', function() {
		// Open the modal
		$("#modal-add").modal();
	});

}