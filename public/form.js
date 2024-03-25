const { degrees, PDFDocument, rgb, StandardFonts } = PDFLib;

let existingPdfBytes;
const pngSize = 3200;

async function onFileChange(dom, e) {
	if (e.target.files[0]) {
		existingPdfBytes = await new Response(e.target.files[0]).arrayBuffer();
		reRenderPDF();
	} else {
		document.getElementById("pdf-editor").setAttribute("hidden", "true");
	}
}

async function reRenderPDF() {
	document.getElementById("qr-page").removeAttribute("aria-invalid");
	document.getElementById("qr-page-container").setAttribute("hidden", "true");
	const pngUrl = "/public/qr.png";
	const pngImageBytes = await fetch(pngUrl).then((res) => res.arrayBuffer());

	let qrWidth = document.getElementById("qr-location-width").value;
	let qrHeight = document.getElementById("qr-location-height").value;
	let qrSize = parseInt(document.getElementById("qr-size").value);
	let isAllPages = document.getElementById("all-page").checked;

	const pdfDoc = await PDFDocument.load(existingPdfBytes);

	const pngImage = await pdfDoc.embedPng(pngImageBytes);

	const pngDims = pngImage.scale((1 / pngSize) * qrSize);

	const pages = pdfDoc.getPages();

	if (isAllPages) {
		for (let i = 0; i < pages.length; i++) {
			const page = pdfDoc.getPages()[i];

			let widthLoc = (page.getWidth() - pngDims.width) * (qrWidth / 100);
			let heightLoc = (page.getHeight() - pngDims.height) * (qrHeight / 100);

			page.drawImage(pngImage, {
				x: widthLoc,
				y: heightLoc,
				width: pngDims.width,
				height: pngDims.height,
			});
		}
	} else {
		document.getElementById("qr-page-container").removeAttribute("hidden");
		let qrPage = parseInt(document.getElementById("qr-page").value);
		if (qrPage > pages.length || qrPage < 1) {
			document.getElementById("qr-page").setAttribute("aria-invalid", "true");
			return;
		}
		const page = pdfDoc.getPages()[qrPage - 1];

		let widthLoc = (page.getWidth() - pngDims.width) * (qrWidth / 100);
		let heightLoc = (page.getHeight() - pngDims.height) * (qrHeight / 100);

		page.drawImage(pngImage, {
			x: widthLoc,
			y: heightLoc,
			width: pngDims.width,
			height: pngDims.height,
		});
	}

	const pdfDocBytes = await pdfDoc.save();

	const editedSrc = URL.createObjectURL(
		new Blob([new Uint8Array(pdfDocBytes)], { type: "application/pdf" })
	);

	document.getElementById("pdf-editor").removeAttribute("hidden");
	document.getElementById("pdf").src = editedSrc;
}
