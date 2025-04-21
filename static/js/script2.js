function applyServiceTypeFilter(type) {
    updateUrlParam('type', type);
}

function applyDateFilter(date) {
    updateUrlParam('date', date);
}

function applySearchFilter(query) {
    updateUrlParam('search', query);
}

function updateUrlParam(param, value) {
    const url = new URL(window.location.href);
    url.searchParams.set('page', '1');
    if (value) {
        url.searchParams.set(param, value);
    } else {
        url.searchParams.delete(param);
    }
    window.location.href = url.toString();
}