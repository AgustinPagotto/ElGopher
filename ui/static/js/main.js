// Enhanced navigation active state
var navLinks = document.querySelectorAll("nav a");
var currentPath = window.location.pathname;

for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i];
	var href = link.getAttribute('href');
	
	if (href === currentPath || (href !== '/' && currentPath.startsWith(href))) {
		link.classList.add("live");
		link.setAttribute('aria-current', 'page');
		break;
	}
}

// HTMX event listeners for better UX
document.body.addEventListener('htmx:beforeRequest', function(evt) {
	// Disable submit buttons during requests
	var button = evt.detail.elt.querySelector('button[type="submit"]');
	if (button) {
		button.disabled = true;
	}
});

document.body.addEventListener('htmx:afterRequest', function(evt) {
	// Re-enable submit buttons after requests
	var button = evt.detail.elt.querySelector('button[type="submit"]');
	if (button) {
		button.disabled = false;
	}
	
	// Show success feedback for successful requests
	if (evt.detail.successful && evt.detail.xhr.status === 200) {
		// Could add toast notification here
	}
});

document.body.addEventListener('htmx:responseError', function(evt) {
	// Handle error responses
	console.error('Request failed:', evt.detail);
	var errorDiv = document.getElementById('form-errors');
	if (errorDiv) {
		errorDiv.innerHTML = '<p role="alert" style="color: var(--pico-color-red-600);">An error occurred. Please try again.</p>';
	}
});
