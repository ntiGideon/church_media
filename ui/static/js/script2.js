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
 document.querySelector('.file-input').addEventListener('change', function(e) {
    const file = e.target.files[0];
    const preview = document.getElementById('photo-preview');
    const previewContainer = document.getElementById('photo-preview-container');

    if (file) {
        const reader = new FileReader();
        reader.onload = function(e) {
            preview.src = e.target.result;
            previewContainer.style.display = 'flex';
        }
        reader.readAsDataURL(file);
        document.querySelector('.file-name').textContent = file.name;
    } else {
        previewContainer.style.display = 'none';
        document.querySelector('.file-name').textContent = 'No file selected';
    }
});

// Function to show modal with member details
function showSuccessModal(memberId) {
  document.getElementById('successMemberId').textContent = memberId;
  document.getElementById('successModal').classList.add('is-active');
}

// Function to close modal
function closeModal() {
  document.getElementById('successModal').classList.remove('is-active');
}

// Function to handle printing (placeholder)
function printMemberCard() {
  // Implement your print functionality here
  alert('Print functionality would generate ID card for member');
}

// Close modal when clicking background
document.querySelector('.modal-background').addEventListener('click', closeModal);



// Close modal when clicking background or close button
document.querySelectorAll('.modal-background, .modal-card-head .delete').forEach(el => {
    el.addEventListener('click', function() {
        document.querySelector('.modal').classList.remove('is-active');
    });
});


function confirmDelete() {
  return confirm('Are you sure you want to delete this member?');
}





// Delete Confirmation Modal
let currentDeleteUrl = '';

function showDeleteModal(memberId, memberName) {
  const modal = document.getElementById('deleteModal');
  const confirmationText = document.getElementById('deleteConfirmationText');
  const deleteButton = document.getElementById('confirmDeleteButton');

  // Set the confirmation message
  confirmationText.textContent = `Are you sure you want to delete ${memberName}?`;

  // Set the delete URL
  currentDeleteUrl = `/members/${memberId}/delete`;

  // Show the modal
  modal.classList.add('is-active');

  // Remove any existing click handlers to prevent duplicates
  deleteButton.replaceWith(deleteButton.cloneNode(true));

  // Add new click handler
  document.getElementById('confirmDeleteButton').addEventListener('click', function() {
    window.location.href = currentDeleteUrl;
  });
}

function closeDeleteModal() {
  document.getElementById('deleteModal').classList.remove('is-active');
}

// Close modal when clicking background
document.querySelector('.modal-background').addEventListener('click', closeDeleteModal);

// Initialize all delete buttons
document.addEventListener('DOMContentLoaded', function() {
  document.querySelectorAll('[data-delete-member]').forEach(button => {
    button.addEventListener('click', function(e) {
      e.preventDefault();
      const memberId = this.getAttribute('data-member-id');
      const memberName = this.getAttribute('data-member-name');
      showDeleteModal(memberId, memberName);
    });
  });
});



document.addEventListener('DOMContentLoaded', function() {
  // Membership Growth Chart
  const growthCtx = document.getElementById('growthChart').getContext('2d');
  const growthChart = new Chart(growthCtx, {
    type: 'line',
    data: {
      labels: {{ .ChartData.Months}},
      datasets: [{
        label: 'New Members',
        data: {{.ChartData.NewMembers}},
        borderColor: '#485fc7',
        backgroundColor: 'rgba(72, 95, 199, 0.1)',
        tension: 0.3,
        fill: true
      }, {
        label: 'Total Members',
        data: {{.ChartData.TotalMembers}},
        borderColor: '#00d1b2',
        backgroundColor: 'rgba(0, 209, 178, 0.1)',
        tension: 0.3,
        fill: true
      }]
    },
    options: {
      responsive: true,
      plugins: {
        legend: {
          position: 'top',
        }
      },
      scales: {
        y: {
          beginAtZero: true
        }
      }
    }
  });

  // Age Distribution Chart
  const ageCtx = document.getElementById('ageChart').getContext('2d');
  const ageChart = new Chart(ageCtx, {
    type: 'doughnut',
    data: {
      labels: ['0-18', '19-30', '31-45', '46-60', '60+'],
      datasets: [{
        data: {{.ChartData.AgeGroups}},
        backgroundColor: [
          '#36a2eb',
          '#ff6384',
          '#ffcd56',
          '#4bc0c0',
          '#9966ff'
        ]
      }]
    },
    options: {
      responsive: true,
      plugins: {
        legend: {
          position: 'right',
        }
      }
    }
  });

  const regionCtx = document.getElementById('regionChart').getContext('2d');
  const regionChart = new Chart(regionCtx, {
    type: 'bar',
    data: {
      labels: {ChartData.Regions}},
      datasets: [{
        label: 'Members',
        data: {{.ChartData.RegionCounts}},
        backgroundColor: 'rgba(72, 95, 199, 0.7)'
      }]
    },
    options: {
      responsive: true,
      scales: {
        y: {
          beginAtZero: true
        }
      }
    }
  });

  // Gender Distribution Chart
  const genderCtx = document.getElementById('genderChart').getContext('2d');
  const genderChart = new Chart(genderCtx, {
    type: 'pie',
    data: {
      labels: ['Male', 'Female', 'Other'],
      datasets: [{
        data: {{.ChartData.GenderDistribution}},
        backgroundColor: [
          '#36a2eb',
          '#ff6384',
          '#ffcd56'
        ]
      }]
    },
    options: {
      responsive: true,
      plugins: {
        legend: {
          position: 'right',
        }
      }
    }
  });

  // Time filter change
  document.getElementById('timeFilter').addEventListener('change', function() {
    // In a real app, you would reload data based on the filter
    console.log('Time filter changed to:', this.value);
  });
});




document.addEventListener('DOMContentLoaded', function() {
    // Calculate total attendance
    const malesInput = document.querySelector('input[name="males"]');
    const femalesInput = document.querySelector('input[name="females"]');
    const totalInput = document.querySelector('input[name="total"]');

    function calculateTotal() {
        const males = parseInt(malesInput.value) || 0;
        const females = parseInt(femalesInput.value) || 0;
        totalInput.value = males + females;
    }

    malesInput.addEventListener('input', calculateTotal);
    femalesInput.addEventListener('input', calculateTotal);

    // Set default date to today
    const today = new Date().toISOString().split('T')[0];
    document.querySelector('input[type="date"]').value = today;
});

// Sorting functionality
document.querySelectorAll('.is-sortable').forEach(header => {
    header.addEventListener('click', () => {
        const sortField = header.dataset.sort;
        const currentUrl = new URL(window.location.href);

        if (currentUrl.searchParams.get('sort') === sortField) {
            // Toggle sort direction if same field clicked
            currentUrl.searchParams.set('dir',
                currentUrl.searchParams.get('dir') === 'asc' ? 'desc' : 'asc');
        } else {
            // New sort field, default to ascending
            currentUrl.searchParams.set('sort', sortField);
            currentUrl.searchParams.set('dir', 'asc');
        }

        window.location.href = currentUrl.toString();
    });
});

// Filter functionality
document.getElementById('service-type-filter').addEventListener('change', applyFilters);
document.getElementById('date-filter').addEventListener('change', applyFilters);
document.getElementById('search-input').addEventListener('input', applyFilters);
document.getElementById('clear-filters').addEventListener('click', clearFilters);

function applyFilters() {
    const serviceType = document.getElementById('service-type-filter').value;
    const date = document.getElementById('date-filter').value;
    const search = document.getElementById('search-input').value.toLowerCase();

    const rows = document.querySelectorAll('#records-table tbody tr');

    rows.forEach(row => {
        const rowServiceType = row.querySelector('td:nth-child(2) span').textContent.toLowerCase();
        const rowDate = row.querySelector('td:nth-child(1)').textContent;
        const rowText = row.textContent.toLowerCase();

        const typeMatch = !serviceType ||
            (serviceType === 'first_service' && rowServiceType === 'first') ||
            (serviceType === 'second_service' && rowServiceType === 'second') ||
            (serviceType === 'special_service' && rowServiceType === 'special');

        const dateMatch = !date || rowDate.includes(new Date(date).toLocaleDateString('en-US', {
            month: 'short',
            day: '2-digit',
            year: 'numeric'
        }).replace(',', ''));

        const searchMatch = !search || rowText.includes(search);

        row.style.display = (typeMatch && dateMatch && searchMatch) ? '' : 'none';
    });
}

function clearFilters() {
    document.getElementById('service-type-filter').value = '';
    document.getElementById('date-filter').value = '';
    document.getElementById('search-input').value = '';
    applyFilters();
}

// Delete modal functionality
let recordToDelete = null;

function confirmDelete(id) {
    recordToDelete = id;
    document.getElementById('delete-modal').classList.add('is-active');
}

function confirmEdit(id) {
    recordToDelete = id;
    document.getElementById('edit-modal').classList.add('is-active');
}

function closeModal() {
    document.getElementById('delete-modal').classList.remove('is-active');
}

document.getElementById('confirm-delete').addEventListener('click', () => {
    if (recordToDelete) {
        fetch(`/service-record/${recordToDelete}/delete`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            }
        })
        .then(response => {
            if (response.ok) {
                window.location.reload();
            }
        });
    }
});





// Initialize charts with data from backend
document.addEventListener('DOMContentLoaded', function() {
    // Growth Chart
    const growthCtx = document.getElementById('growthChart').getContext('2d');
    new Chart(growthCtx, {
        type: 'line',
        data: {
            labels: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"],
            datasets: [{
                label: 'Membership Growth',
                data: [34,56,22,76,53,67,26],
                backgroundColor: 'rgba(54, 162, 235, 0.2)',
                borderColor: 'rgba(54, 162, 235, 1)',
                borderWidth: 2,
                tension: 0.3
            }]
        },
        options: {
            responsive: true,
            plugins: {
                legend: { display: false }
            }
        }
    });

    // Age Distribution Chart
    const ageCtx = document.getElementById('ageChart').getContext('2d');
    new Chart(ageCtx, {
        type: 'pie',
        data: {
            labels: {{.ChartData.AgeGroups }},
            datasets: [{
                data: {{.ChartData.AgeCounts }},
                backgroundColor: [
                    '#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF'
                ]
            }]
        }
    });

    // Gender Chart
    const genderCtx = document.getElementById('genderChart').getContext('2d');
    new Chart(genderCtx, {
        type: 'doughnut',
        data: {
            labels: ['Male', 'Female'],
            datasets: [{
                data: [{{.Stats.MaleMembers}}, {{.Stats.FemaleMembers}}],
                backgroundColor: ['#36A2EB', '#FF6384']
            }]
        },
        options: {
            cutout: '70%',
            plugins: { legend: { position: 'bottom' } }
        }
    });

    // Region Chart
    const regionCtx = document.getElementById('regionChart').getContext('2d');
    new Chart(regionCtx, {
        type: 'bar',
        data: {
            labels: {{.ChartData.Regions }},
            datasets: [{
                label: 'Members',
                data: {{.ChartData.RegionCounts }},
                backgroundColor: '#00D1B2'
            }]
        },
        options: {
            indexAxis: 'y',
            plugins: { legend: { display: false } }
        }
    });
});