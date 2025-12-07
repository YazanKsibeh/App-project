<script>
  // No JavaScript needed for this static footer
</script>

<footer class="company-footer">
  <div class="footer-content">
    <div class="footer-section">
      <div class="footer-brand">
        <img
          src="/assets/images/novatech-logo.png"
          alt="NovaTech Labs Logo"
          class="footer-logo"
        />
        <span class="brand-name">NovaTech Labs</span>
      </div>
      <p class="footer-description">
        NovaTech Labs – Engineering Ideas Into Reality.
      </p>
      <p class="footer-description">
        0941-414-122
      </p>
      <div class="social-icons">
        <a href="https://github.com/nova-tech-labs-solutions" class="social-link" aria-label="GitHub" target="_blank" rel="noopener noreferrer">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"/>
          </svg>
        </a>
        <a href="https://www.linkedin.com/company/nova-tech-labs/" class="social-link" aria-label="LinkedIn" target="_blank" rel="noopener noreferrer">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z"/>
            <rect x="2" y="9" width="4" height="12"/>
            <circle cx="4" cy="4" r="2"/>
          </svg>
        </a>
        <a href="https://www.facebook.com/profile.php?id=61579677733198" class="social-link" aria-label="Facebook" target="_blank" rel="noopener noreferrer">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M18 2h-3a5 5 0 0 0-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 0 1 1-1h3z"/>
          </svg>
        </a>
        <a href="https://x.com/Novatech_Labs" class="social-link" aria-label="X (Twitter)" target="_blank" rel="noopener noreferrer">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M18 6L6 18"/>
            <path d="M6 6l12 12"/>
          </svg>
        </a>
      </div>
    </div>

    <div class="footer-section">
      <h4>Products</h4>
      <ul>
        <li><a href="https://novatech-labs.com/#products" target="_blank" rel="noopener noreferrer">DentistApp Pro</a></li>
        <li><a href="https://novatech-labs.com/#products" target="_blank" rel="noopener noreferrer">LabConnect</a></li>
        <!-- <li><a href="https://novatech-labs.com/#products" target="_blank" rel="noopener noreferrer">Enterprise ERP (Soon)</a></li> -->
        <li><a href="https://novatech-labs.com/#products" target="_blank" rel="noopener noreferrer">Analytics Pro (Soon)</a></li>
      </ul>
    </div>

    <div class="footer-section">
      <h4>Company</h4>
      <ul>
        <li><a href="https://novatech-labs.com/#about" target="_blank" rel="noopener noreferrer">About Us</a></li>
        <li><a href="https://novatech-labs.com/#contact" target="_blank" rel="noopener noreferrer">Contact</a></li>
        <li><a href="https://novatech-labs.com/privacy" target="_blank" rel="noopener noreferrer">Privacy Policy</a></li>
        <li><a href="https://novatech-labs.com/terms" target="_blank" rel="noopener noreferrer">Terms of Service</a></li>
      </ul>
    </div>

    <div class="footer-section">
      <h4>Support</h4>
      <ul>
        <li><a href="https://novatech-labs.com/docs" target="_blank" rel="noopener noreferrer">Documentation</a></li>
        <li><a href="https://novatech-labs.com/help" target="_blank" rel="noopener noreferrer">Help Center</a></li>
        <li><a href="mailto:novatech.labs@outlook.com">Contact Support</a></li>
      </ul>
    </div>
  </div>

  <div class="footer-bottom">
    <p>&copy; 2025 NovaTech Labs. All rights reserved.</p>
    <p>Engineering Ideas Into Reality</p>
  </div>
</footer>

<style>
  .company-footer {
    background: var(--color-bg-footer, #000000);
    color: var(--color-text-inverted, #ffffff);
    padding: 3rem 0 1rem;
    margin-top: 0;
    border-top: 1px solid var(--color-border-primary, #2f3336);
    width: 100%;
    box-sizing: border-box;
    overflow-x: hidden;
  }

  .footer-content {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 2rem;
    margin-bottom: 2rem;
    max-width: 1200px;
    margin-left: auto;
    margin-right: auto;
    padding: 0 2rem;
    width: 100%;
  }

  .footer-brand {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  .footer-logo {
    width: 40px;
    height: 40px;
    object-fit: contain;
    border-radius: 8px;
  }

  .brand-name {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--color-text-inverted, #ffffff);
  }

  .footer-description {
    color: var(--color-text-secondary, #71767b);
    margin-bottom: 1.5rem;
    line-height: 1.6;
  }

  .social-icons {
    display: flex;
    gap: 1rem;
  }

  .social-link {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 44px;
    height: 44px;
    border-radius: 8px;
    background: rgba(255, 255, 255, 0.1);
    color: rgba(255, 255, 255, 0.8);
    text-decoration: none;
    transition: all 0.15s ease-in-out;
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);
  }

  .social-link:hover {
    background: rgba(255, 255, 255, 0.2);
    color: white;
    transform: translateY(-2px);
  }

  .footer-section h4 {
    font-weight: 600;
    margin-bottom: 1rem;
    color: white;
  }

  .footer-section ul {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .footer-section ul li {
    margin-bottom: 0.5rem;
  }

  .footer-section ul li a {
    color: var(--color-text-secondary, #71767b);
    text-decoration: none;
    transition: color 0.15s ease-in-out;
  }

  .footer-section ul li a:hover {
    color: white;
  }

  .footer-bottom {
    border-top: 1px solid var(--color-border-primary, #2f3336);
    padding-top: 2rem;
    text-align: center;
    color: var(--color-text-secondary, #71767b);
    max-width: 1200px;
    margin-left: auto;
    margin-right: auto;
    padding-left: 2rem;
    padding-right: 2rem;
  }

  .footer-bottom p {
    margin-bottom: 0.5rem;
  }

  /* Responsive Design */
  @media (max-width: 768px) {
    .footer-content {
      grid-template-columns: 1fr;
      gap: 1.5rem;
      padding: 0 1rem;
    }
    
    .footer-bottom {
      padding-left: 1rem;
      padding-right: 1rem;
    }
  }
</style>
