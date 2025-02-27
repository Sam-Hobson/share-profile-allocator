const searchBox = document.getElementById("ticker-search-box");
const dropdownList = document.getElementById("ticker-search-dropdown-list");
document.addEventListener("htmx:afterSwap", function(event) {
	if (event.detail.target.id === "ticker-search-dropdown-list") {
		dropdownList.style.display = "block";
	}
});

let selectedIndex = -1; // Track selected index

function closeDropdown() {
	dropdownList.style.display = "none";
	selectedIndex = -1;
}

function selectOption(value) {
	const item = value.querySelector("span");
	searchBox.value = item.textContent;
	closeDropdown();
}

searchBox.addEventListener("keydown", (e) => {
	const items = dropdownList.querySelectorAll("div");
	if (items.length === 0) return;

	if (e.key === "ArrowDown") {
		e.preventDefault();
		selectedIndex = (selectedIndex + 1) % items.length;
		highlightItem(items);
	} else if (e.key === "ArrowUp") {
		e.preventDefault();
		if (selectedIndex === -1) {
			selectedIndex = items.length - 1;
		} else {
			selectedIndex = (selectedIndex - 1 + items.length) % items.length;
		}
		highlightItem(items);
	} else if (e.key === "Enter") {
		e.preventDefault();
		if (selectedIndex >= 0) {
			selectOption(items[selectedIndex])
		};
	} else if (e.key === "Tab") {
		closeDropdown();
	}
});

function highlightItem(items) {
	items.forEach(item => item.classList.remove("active"));
	if (selectedIndex >= 0) {
		const selectedItem = items[selectedIndex];
		selectedItem.classList.add("active");

		// Auto-scroll selected item into view
		dropdownList.scrollTop = selectedItem.offsetTop - dropdownList.offsetTop;
	}
}

dropdownList.addEventListener("click", (e) => {
	dropdownList.querySelectorAll("div").forEach(item => {
		if (item.contains(e.target)) {
			selectOption(item);
			return;
		}
	})
})

document.addEventListener("click", (e) => {
	if (!searchBox.contains(e.target) && !dropdownList.contains(e.target)) {
		closeDropdown();
	}
});
