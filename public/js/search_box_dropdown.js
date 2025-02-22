const options = ["VAS", "VGS", "VDHG", "DHHF", "IVV", "NDQ"];
const searchBox = document.getElementById("ticker-search-box");
const dropdownList = document.getElementById("ticker-search-dropdown-list");
let selectedIndex = -1; // Track selected index

function updateDropdown() {
	const value = searchBox.value.toLowerCase();
	dropdownList.innerHTML = "";
	selectedIndex = -1; // Reset selection
	if (value) {
		const filteredOptions = options.filter(opt => opt.toLowerCase().includes(value));
		filteredOptions.forEach((opt, index) => {
			const div = document.createElement("div");
			div.textContent = opt;
			div.setAttribute("data-index", index);
			div.addEventListener("click", () => selectOption(opt));
			dropdownList.appendChild(div);
		});
		dropdownList.style.display = filteredOptions.length ? "block" : "none";
	} else {
		dropdownList.style.display = "none";
	}
}

function selectOption(value) {
	searchBox.value = value;
	dropdownList.style.display = "none";
}

searchBox.addEventListener("input", updateDropdown);

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
			selectOption(items[selectedIndex].textContent)
		};
	} else if (e.key === "Tab") {
		dropdownList.style.display = "none"; // Hide dropdown on Tab
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

document.addEventListener("click", (e) => {
	if (!searchBox.contains(e.target) && !dropdownList.contains(e.target)) {
		dropdownList.style.display = "none";
	}
});
