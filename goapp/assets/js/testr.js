function reduce() {
	var requests = $("#filter").val().toLowerCase();
	var queries = requests.split(" ")	
	$("tbody").find("tr").each(function (idx, r) {
		var labels = $(r).data("label").toLowerCase();	
						
		var ok = queries.every(function (q) {
			return labels.indexOf(q) >= 0
		});
		if (!ok) {			
			$(r).addClass("hidden");
		} else {			
			$(r).removeClass("hidden");
		}
	});
}

function filterOn(f) {
	$("#filter").val(f);
	reduce();
}

function showOutput(e) {
	var o = $("#output");
	o.find("h4").html($(e).data("test"));
	o.find("code").html($(e).data("output"));
	o.modal().show();
}