function updateRowNumbers() {
    const rows = document.querySelectorAll('#share-table-body tr[id^="table-row-"]');
    rows.forEach((row, index) => {
        const th = row.querySelector('th');
        th.textContent = index + 1;
    });
}

const targetNode = document.getElementById('share-table-body');
const config = { childList: true, subtree: true };

const callback = function(mutationsList) {
    for (let mutation of mutationsList) {
        if (mutation.type === 'childList') {
            observer.disconnect();
            updateRowNumbers();
            observer.observe(targetNode, config);
        }
    }
};

const observer = new MutationObserver(callback);
observer.observe(targetNode, config);

updateRowNumbers();
