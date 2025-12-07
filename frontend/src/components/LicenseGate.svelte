<script>
  import { onMount } from 'svelte';
  import { ValidateLicense } from '../../wailsjs/go/main/App.js';

  export let onLicenseValidated = () => {};

  let licenseKey = '';
  let error = '';
  let loading = false;
  let licenseInfo = null;

  // Check for existing license in localStorage
  onMount(() => {
    const savedLicense = localStorage.getItem('dentist_license_key');
    if (savedLicense) {
      licenseKey = savedLicense;
      validateLicense();
    }
  });

  async function validateLicense() {
    if (!licenseKey.trim()) {
      error = 'Please enter a license key';
      return;
    }

    loading = true;
    error = '';

    try {
      licenseInfo = await ValidateLicense(licenseKey);
      
      if (licenseInfo.is_valid) {
        // Save valid license to localStorage
        localStorage.setItem('dentist_license_key', licenseKey);
        localStorage.setItem('dentist_license_expiry', licenseInfo.expiry_date);
        
        // Notify parent component that license is valid
        onLicenseValidated(licenseKey, licenseInfo);
      } else {
        error = licenseInfo.message || 'Invalid license key';
        // Remove invalid license from localStorage
        localStorage.removeItem('dentist_license_key');
        localStorage.removeItem('dentist_license_expiry');
      }
    } catch (err) {
      error = `Validation error: ${err.message}`;
      console.error('License validation error:', err);
    } finally {
      loading = false;
    }
  }

  function handleKeyPress(event) {
    if (event.key === 'Enter') {
      validateLicense();
    }
  }
</script>

<div class="license-gate">
  <div class="license-container">
    <div class="logo-section">
      <h1>🦷 DentistApp</h1>
      <p class="subtitle">Professional Dental Practice Management</p>
    </div>

    <div class="license-form">
      <h2>License Activation Required</h2>
      <p class="description">
        Please enter your license key to access the application. 
        If you don't have a license key, contact your software provider.
      </p>

      <div class="input-group">
        <label for="license-input">License Key:</label>
        <input
          id="license-input"
          type="text"
          bind:value={licenseKey}
          on:keypress={handleKeyPress}
          placeholder="Enter your license key..."
          disabled={loading}
          class="license-input"
        />
      </div>

      {#if error}
        <div class="error-message">
          ❌ {error}
        </div>
      {/if}

      {#if licenseInfo && !licenseInfo.is_valid}
        <div class="license-details">
          <p><strong>Expiry Date:</strong> {licenseInfo.expiry_date}</p>
          <p><strong>Status:</strong> <span class="expired">Expired</span></p>
        </div>
      {/if}

      <button
        on:click={validateLicense}
        disabled={loading || !licenseKey.trim()}
        class="validate-button"
      >
        {#if loading}
          🔄 Validating...
        {:else}
          🔑 Activate License
        {/if}
      </button>

      <div class="help-section">
        <h3>Need Help?</h3>
        <ul>
          <li>Contact your software provider for a new license key</li>
          <li>Ensure your license key is entered exactly as provided</li>
          <li>Check that your license hasn't expired</li>
        </ul>
      </div>
    </div>
  </div>
</div>

<style>
  .license-gate {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    z-index: 10000;
  }

  .license-container {
    background: white;
    border-radius: 16px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
    max-width: 500px;
    width: 100%;
    overflow: hidden;
  }

  .logo-section {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    text-align: center;
    padding: 3rem 2rem;
  }

  .logo-section h1 {
    margin: 0;
    font-size: 2.5rem;
    font-weight: 700;
  }

  .subtitle {
    margin: 0.5rem 0 0 0;
    font-size: 1.1rem;
    opacity: 0.9;
  }

  .license-form {
    padding: 2.5rem;
  }

  .license-form h2 {
    margin: 0 0 1rem 0;
    color: #333;
    font-size: 1.5rem;
    text-align: center;
  }

  .description {
    color: #666;
    line-height: 1.6;
    margin-bottom: 2rem;
    text-align: center;
  }

  .input-group {
    margin-bottom: 1.5rem;
  }

  .input-group label {
    display: block;
    margin-bottom: 0.5rem;
    color: #333;
    font-weight: 600;
  }

  .license-input {
    width: 100%;
    padding: 1rem;
    border: 2px solid #e1e5e9;
    border-radius: 8px;
    font-size: 1rem;
    font-family: 'Courier New', monospace;
    transition: border-color 0.2s;
    box-sizing: border-box;
  }

  .license-input:focus {
    outline: none;
    border-color: #667eea;
  }

  .license-input:disabled {
    background-color: #f5f5f5;
    cursor: not-allowed;
  }

  .error-message {
    background: #fee;
    color: #c53030;
    padding: 1rem;
    border-radius: 8px;
    margin-bottom: 1.5rem;
    text-align: center;
    border: 1px solid #fed7d7;
  }

  .license-details {
    background: #fafafa;
    padding: 1rem;
    border-radius: 8px;
    margin-bottom: 1.5rem;
    border: 1px solid #e1e5e9;
  }

  .expired {
    color: #c53030;
    font-weight: 600;
  }

  .validate-button {
    width: 100%;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    padding: 1rem 2rem;
    border-radius: 8px;
    font-size: 1.1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    margin-bottom: 2rem;
  }

  .validate-button:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 10px 30px rgba(102, 126, 234, 0.3);
  }

  .validate-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none;
  }

  .help-section {
    border-top: 1px solid #e1e5e9;
    padding-top: 1.5rem;
  }

  .help-section h3 {
    margin: 0 0 1rem 0;
    color: #333;
    font-size: 1.1rem;
  }

  .help-section ul {
    margin: 0;
    padding-left: 1.2rem;
    color: #666;
    line-height: 1.6;
  }

  .help-section li {
    margin-bottom: 0.5rem;
  }
</style>