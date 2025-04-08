 // document.addEventListener('DOMContentLoaded', function(){
 //            // Sample message data
 //            const messages = {
 //                1: {
 //                    sender: "John Smith",
 //                    email: "john.smith@example.com",
 //                    subject: "Prayer request for my family",
 //                    date: "June 15, 2023 at 10:30 AM",
 //                    content: "Dear Pastors,\n\nMy wife has been diagnosed with a serious illness and we would appreciate your prayers. Our children are also struggling with this news. Please keep our family in your prayers during this difficult time.\n\nThank you,\nJohn Smith",
 //                    tags: ["Prayer", "Unread"],
 //                    responded: false
 //                },
 //                2: {
 //                    sender: "Sarah Johnson",
 //                    email: "sarah.j@example.com",
 //                    subject: "Question about baptism",
 //                    date: "June 14, 2023 at 3:45 PM",
 //                    content: "Hello,\n\nI'm interested in being baptized but have some questions about the process. What are the requirements? Do I need to take classes first? When is the next baptism service scheduled?\n\nThanks,\nSarah",
 //                    tags: ["Question", "Responded"],
 //                    responded: true
 //                },
 //                3: {
 //                    sender: "Michael Brown",
 //                    email: "michael.b@example.com",
 //                    subject: "Volunteer opportunity inquiry",
 //                    date: "June 14, 2023 at 11:20 AM",
 //                    content: "Good morning,\n\nI'd like to get involved in serving at the church. I have experience with sound systems and would love to help with the tech team. Who should I contact about volunteering?\n\nBlessings,\nMichael",
 //                    tags: ["Volunteer", "Unread"],
 //                    responded: false
 //                },
 //                4: {
 //                    sender: "Emily Davis",
 //                    email: "emily.d@example.com",
 //                    subject: "Need information about VBS",
 //                    date: "June 13, 2023 at 4:15 PM",
 //                    content: "Hi there,\n\nI have two children who would like to attend Vacation Bible School this summer. Can you send me information about dates, times, and how to register? Also, is there a cost?\n\nThank you,\nEmily Davis",
 //                    tags: ["Question", "Unread"],
 //                    responded: false
 //                },
 //                5: {
 //                    sender: "Robert Wilson",
 //                    email: "robert.w@example.com",
 //                    subject: "Building fund donation question",
 //                    date: "June 12, 2023 at 9:10 AM",
 //                    content: "Dear Church Office,\n\nI'd like to make a donation to the building fund. Can this be done online or should I send a check? Also, is the donation tax-deductible?\n\nSincerely,\nRobert Wilson",
 //                    tags: ["Question", "Responded"],
 //                    responded: true
 //                }
 //            };
 //
 //            // Filter messages
 //            const filterTabs = document.querySelectorAll('.tabs li');
 //            filterTabs.forEach(tab => {
 //                tab.addEventListener('click', function() {
 //                    const filter = this.getAttribute('data-filter');
 //
 //                    // Update active tab
 //                    filterTabs.forEach(t => t.classList.remove('is-active'));
 //                    this.classList.add('is-active');
 //
 //                    // Filter messages
 //                    const messageCards = document.querySelectorAll('.message-card');
 //                    messageCards.forEach(card => {
 //                        if (filter === 'all') {
 //                            card.style.display = 'block';
 //                        } else if (filter === 'unread') {
 //                            if (card.classList.contains('unread')) {
 //                                card.style.display = 'block';
 //                            } else {
 //                                card.style.display = 'none';
 //                            }
 //                        } else if (filter === 'responded') {
 //                            if (card.classList.contains('responded')) {
 //                                card.style.display = 'block';
 //                            } else {
 //                                card.style.display = 'none';
 //                            }
 //                        } else if (filter === 'prayer') {
 //                            const tags = card.querySelector('.tags').textContent;
 //                            if (tags.includes('Prayer')) {
 //                                card.style.display = 'block';
 //                            } else {
 //                                card.style.display = 'none';
 //                            }
 //                        }
 //                    });
 //                });
 //            });
 //
 //            // Message selection
 //            const messageCards = document.querySelectorAll('.message-card');
 //            messageCards.forEach(card => {
 //                card.addEventListener('click', function() {
 //                    const messageId = this.getAttribute('data-id');
 //                    const message = messages[messageId];
 //
 //                    // Update message view
 //                    document.querySelector('.message-sender').textContent = message.sender;
 //                    document.querySelector('.message-email').textContent = message.email;
 //                    document.querySelector('.message-subject').textContent = message.subject;
 //                    document.querySelector('.message-date').textContent = message.date;
 //                    document.querySelector('.message-content').innerHTML = message.content.replace(/\n/g, '<br>');
 //
 //                    // Update reply fields
 //                    document.querySelector('.reply-email').value = message.email;
 //                    document.querySelector('.reply-subject').value = `Re: ${message.subject}`;
 //
 //                    // Show/hide elements
 //                    document.querySelector('.empty-state').style.display = 'none';
 //                    document.querySelector('.message-view').style.display = 'block';
 //                    document.querySelector('.reply-box').style.display = 'none';
 //
 //                    if (message.responded) {
 //                        document.querySelector('.previous-responses').style.display = 'block';
 //                    } else {
 //                        document.querySelector('.previous-responses').style.display = 'none';
 //                    }
 //
 //                    // Highlight selected card
 //                    messageCards.forEach(c => c.classList.remove('has-background-light'));
 //                    this.classList.add('has-background-light');
 //                });
 //            });
 //
 //            // Reply functionality
 //            document.querySelector('.reply-button').addEventListener('click', function() {
 //                document.querySelector('.reply-box').style.display = 'block';
 //            });
 //
 //            document.querySelector('.cancel-reply').addEventListener('click', function() {
 //                document.querySelector('.reply-box').style.display = 'none';
 //            });
 //
 //            document.querySelector('.send-reply').addEventListener('click', function() {
 //                const email = document.querySelector('.reply-email').value;
 //                const subject = document.querySelector('.reply-subject').value;
 //                const message = document.querySelector('.reply-message').value;
 //
 //                // In a real implementation, you would send the email here
 //                alert(`Reply sent to ${email} with subject: ${subject}`);
 //
 //                // Mark as responded
 //                const activeCard = document.querySelector('.message-card.has-background-light');
 //                if (activeCard) {
 //                    activeCard.classList.remove('unread');
 //                    activeCard.classList.add('responded');
 //
 //                    // Update tags
 //                    const tagsContainer = activeCard.querySelector('.tags');
 //                    tagsContainer.innerHTML = `
 //                        <span class="tag is-success is-light">Responded</span>
 //                        ${tagsContainer.innerHTML.includes('Prayer') ? '<span class="tag is-primary is-light">Prayer</span>' : ''}
 //                    `;
 //
 //                    // Show previous responses
 //                    document.querySelector('.previous-responses').style.display = 'block';
 //                }
 //
 //                // Clear and hide reply box
 //                document.querySelector('.reply-message').value = '';
 //                document.querySelector('.reply-box').style.display = 'none';
 //            });
 //
 //            // Mark as read
 //            document.querySelector('.mark-as-read').addEventListener('click', function() {
 //                const activeCard = document.querySelector('.message-card.has-background-light');
 //                if (activeCard && activeCard.classList.contains('unread')) {
 //                    activeCard.classList.remove('unread');
 //
 //                    // Update tags
 //                    const tagsContainer = activeCard.querySelector('.tags');
 //                    const tagsHtml = tagsContainer.innerHTML.replace('Unread', '').replace('tag is-warning is-light', '');
 //                    tagsContainer.innerHTML = tagsHtml;
 //                }
 //            });
 //        });

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