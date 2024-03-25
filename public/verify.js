async function verifyFile(e) {
	e.preventDefault();

	const submitButton = document.getElementById("submit");
	submitButton.setAttribute("aria-busy", "true");
	const formData = new FormData();
	formData.append("id", document.getElementById("id").value);
	formData.append("file", document.getElementById("file").files[0]);

	let response = await fetch(window.location.href.split("?")[0], {
		method: "POST",
		body: formData,
	});

	const resultField = document.getElementById("result");
	submitButton.setAttribute("hidden", true);
	submitButton.removeAttribute("aria-busy");
	resultField.removeAttribute("hidden");
	if (response.status == 200) {
		resultField.innerHTML = "<p> File Integrity Verified! </p>";
		resultField.className = "pico-background-jade-500";
	} else {
		resultField.innerHTML = "<p> File Integrity Compromised! </p>";
		resultField.className = "pico-background-red-500";
	}
}

function showButton(e) {
	const submitButton = document.getElementById("submit");
	const resultField = document.getElementById("result");
	resultField.setAttribute("hidden", true);
	submitButton.removeAttribute("hidden");
}
