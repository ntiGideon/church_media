document.addEventListener('DOMContentLoaded', function() {


    // Navbar scroll effect
    window.addEventListener('scroll', function() {
        const navbar = document.querySelector('.navbar');
        if (window.scrollY > 50) {
            navbar.classList.add('is-scrolled');
        } else {
            navbar.classList.remove('is-scrolled');
        }
    });

    // Initialize carousel
    let currentSlide = 0;
    const slides = document.querySelectorAll('.carousel-slide');
    const dots = document.querySelectorAll('.dot');
    const totalSlides = slides.length;

    // Function to show a specific slide
    function showSlide(n) {
        // Reset all slides and dots
        slides.forEach(slide => slide.classList.remove('active'));
        dots.forEach(dot => dot.classList.remove('active'));

        // Update current slide index
        currentSlide = (n + totalSlides) % totalSlides;

        // Activate current slide and dot
        slides[currentSlide].classList.add('active');
        dots[currentSlide].classList.add('active');
    }

    // Function to advance to next slide
    function nextSlide() {
        showSlide(currentSlide + 1);
    }

    // Auto-advance carousel every 5 seconds
    let slideInterval = setInterval(nextSlide, 5000);

    // Pause carousel on hover
    const heroCarousel = document.querySelector('.hero-carousel');
    if (heroCarousel) {
        heroCarousel.addEventListener('mouseenter', () => {
            clearInterval(slideInterval);
        });

        heroCarousel.addEventListener('mouseleave', () => {
            slideInterval = setInterval(nextSlide, 5000);
        });
    }

    // Dot navigation
    dots.forEach((dot, index) => {
        dot.addEventListener('click', () => {
            clearInterval(slideInterval);
            showSlide(index);
            slideInterval = setInterval(nextSlide, 5000);
        });
    });



    // Animation on scroll
    function animateOnScroll() {
        const elements = document.querySelectorAll('.slide-in-top, .slide-in-bottom, .hover-grow, .hover-zoom, .hover-tilt');

        elements.forEach(element => {
            const elementPosition = element.getBoundingClientRect().top;
            const screenPosition = window.innerHeight / 1.2;

            if (elementPosition < screenPosition) {
                element.style.opacity = '1';
            }
        });
    }

    window.addEventListener('scroll', animateOnScroll);
    animateOnScroll(); // Run once on page load
});

// Helper function for carousel navigation
function currentSlide(n) {
    const slides = document.querySelectorAll('.carousel-slide');
    const dots = document.querySelectorAll('.dot');

    slides.forEach(slide => slide.classList.remove('active'));
    dots.forEach(dot => dot.classList.remove('active'));

    slides[n - 1].classList.add('active');
    dots[n - 1].classList.add('active');
}


// Modern Carousel Functionality
class ModernCarousel {
    constructor(container) {
      this.container = container;
      this.slides = container.querySelectorAll('.carousel-slide');
      this.indicators = container.querySelectorAll('.indicator');
      this.prevBtn = container.querySelector('.prev');
      this.nextBtn = container.querySelector('.next');
      this.currentIndex = 0;
      this.interval = null;
      this.autoRotateDelay = 5000;

      this.init();
    }

    init() {
      // Show first slide
      this.showSlide(this.currentIndex);

      // Set up event listeners
      this.prevBtn.addEventListener('click', () => this.prevSlide());
      this.nextBtn.addEventListener('click', () => this.nextSlide());

      // Indicator clicks
      this.indicators.forEach((indicator, index) => {
        indicator.addEventListener('click', () => this.goToSlide(index));
      });

      // Start auto-rotation
      this.startAutoRotate();

      // Pause on hover
      this.container.addEventListener('mouseenter', () => this.stopAutoRotate());
      this.container.addEventListener('mouseleave', () => this.startAutoRotate());
    }

    showSlide(index) {
      // Hide all slides
      this.slides.forEach(slide => slide.classList.remove('active'));
      this.indicators.forEach(indicator => indicator.classList.remove('active'));

      // Show current slide
      this.slides[index].classList.add('active');
      this.indicators[index].classList.add('active');
      this.currentIndex = index;
    }

    nextSlide() {
      const nextIndex = (this.currentIndex + 1) % this.slides.length;
      this.showSlide(nextIndex);
      this.restartAutoRotate();
    }

    prevSlide() {
      const prevIndex = (this.currentIndex - 1 + this.slides.length) % this.slides.length;
      this.showSlide(prevIndex);
      this.restartAutoRotate();
    }

    goToSlide(index) {
      this.showSlide(index);
      this.restartAutoRotate();
    }

    startAutoRotate() {
      this.interval = setInterval(() => this.nextSlide(), this.autoRotateDelay);
    }

    stopAutoRotate() {
      clearInterval(this.interval);
    }

    restartAutoRotate() {
      this.stopAutoRotate();
      this.startAutoRotate();
    }
  }

  // Initialize the carousel when DOM is loaded
  document.addEventListener('DOMContentLoaded', () => {
    const carousels = document.querySelectorAll('.modern-carousel');
    carousels.forEach(carousel => new ModernCarousel(carousel));

    // Keep your existing navbar and other functionality here
    // ...
  });


// Accordion functionality
document.querySelectorAll('.accordion-header').forEach(header => {
    header.addEventListener('click', () => {
        const accordionItem = header.parentElement;
        const accordionContent = header.nextElementSibling;

        // Close all other items
        document.querySelectorAll('.accordion-item').forEach(item => {
            if (item !== accordionItem) {
                item.querySelector('.accordion-header').classList.remove('active');
                item.querySelector('.accordion-content').classList.remove('active');
            }
        });

        // Toggle current item
        header.classList.toggle('active');
        accordionContent.classList.toggle('active');
    });
});

// Initialize animations for ministries page
function animateMinistries() {
    const ministries = document.querySelectorAll('.ministry-card');

    ministries.forEach((ministry, index) => {
        setTimeout(() => {
            ministry.style.opacity = '1';
            ministry.style.transform = 'translateY(0)';
        }, index * 100);
    });
}


// Run when page loads
window.addEventListener('load', animateMinistries);


// Contact Form Submission
document.getElementById('contactFormClose').addEventListener('submit', function(e) {
    e.preventDefault();

    // Here you would typically send the form data to your server
    // For demonstration, we'll just show a success message
    alert('Thank you for your message! We will get back to you soon.');
    this.reset();
});

// Accordion functionality
document.querySelectorAll('.accordion-header').forEach(header => {
    header.addEventListener('click', () => {
        const accordionItem = header.parentElement;
        const accordionContent = header.nextElementSibling;

        // Close all other items
        document.querySelectorAll('.accordion-item').forEach(item => {
            if (item !== accordionItem) {
                item.querySelector('.accordion-header').classList.remove('active');
                item.querySelector('.accordion-content').classList.remove('active');
            }
        });

        // Toggle current item
        header.classList.toggle('active');
        accordionContent.classList.toggle('active');
    });
});



// Give Online Page Functionality
document.addEventListener('DOMContentLoaded', function() {
    // Recurring donation toggle
    const recurringCheckbox = document.getElementById('recurringCheckbox');
    const recurringOptions = document.getElementById('recurringOptions');

    if (recurringCheckbox) {
        recurringCheckbox.addEventListener('change', function() {
            if (this.checked) {
                recurringOptions.style.display = 'block';
            } else {
                recurringOptions.style.display = 'none';
            }
        });
    }

    // Payment method tabs
    const tabs = document.querySelectorAll('.tabs li');
    const tabContents = document.querySelectorAll('.tab-content');

    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            // Remove active class from all tabs
            tabs.forEach(t => t.classList.remove('is-active'));
            // Add active class to clicked tab
            tab.classList.add('is-active');

            // Hide all tab contents
            tabContents.forEach(content => {
                content.style.display = 'none';
            });

            // Show the selected tab content
            const tabId = tab.getAttribute('data-tab');
            document.getElementById(tabId).style.display = 'block';
        });
    });

    // Form submission
    const donationForm = document.getElementById('donationForm');
    if (donationForm) {
        donationForm.addEventListener('submit', function(e) {
            e.preventDefault();
            // In a real implementation, you would process the payment here
            alert('Thank you for your generous donation! A receipt has been sent to your email.');
            this.reset();
            // Hide recurring options if shown
            if (recurringOptions) recurringOptions.style.display = 'none';
        });
    }

    // Accordion functionality
    const accordionHeaders = document.querySelectorAll('.accordion-header');
    accordionHeaders.forEach(header => {
        header.addEventListener('click', () => {
            const accordionItem = header.parentElement;
            const accordionContent = header.nextElementSibling;

            // Close all other items
            document.querySelectorAll('.accordion-item').forEach(item => {
                if (item !== accordionItem) {
                    item.querySelector('.accordion-header').classList.remove('active');
                    item.querySelector('.accordion-content').classList.remove('active');
                }
            });

            // Toggle current item
            header.classList.toggle('active');
            accordionContent.classList.toggle('active');
        });
    });
});


//toaster plugins

document.querySelectorAll(".notification .delete")
            .forEach(function ($deleteButton) {

                // Step 2: Get the parent notification
                // of the delete button
                const parentNotification = $deleteButton.parentNode;

                // Add click event listener on delete
                // button and when the button get clicked
                // remove the parent notification
                $deleteButton.addEventListener('click', function () {
                    parentNotification.parentNode
                        .removeChild(parentNotification);
                });
            });