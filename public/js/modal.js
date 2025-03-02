document.addEventListener("click", (event) => {
	document.querySelectorAll("dialog[id^='share-summary-popup-']").forEach(dialog => {
		const article = dialog.querySelector("article");
		if (
			dialog.open &&
			!article.contains(event.target)
		) {
			dialog.close();
		}
	});
});
