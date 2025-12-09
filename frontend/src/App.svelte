<script>
  import PatientList from './components/PatientList.svelte';
  import PatientDetail from './components/PatientDetail.svelte';
  import AppointmentList from './components/AppointmentList.svelte';
  import AppointmentCalendar from './components/AppointmentCalendar.svelte';
  import Sessions from './components/Sessions.svelte';
  import Payments from './components/Payments.svelte';
  import Financials from './components/Financials.svelte';
  import Configuration from './components/Configuration.svelte';
  import LabOrders from './components/LabOrders.svelte';
  import UserProfileModal from './components/UserProfileModal.svelte';
  import LicenseGate from './components/LicenseGate.svelte';
  import Login from './components/Login.svelte';
  import { selectedPatient } from './stores/patientStore.js';
  import { licenseValid, currentLicenseKey, validateCurrentLicense, theme } from './stores/settingsStore.js';
  import { isAuthenticated, checkAuth } from './stores/authStore.js';
  import { onMount } from 'svelte';

  let currentPage = 'patients'; // 'patients', 'appointments', 'payments', 'calendar', 'sessions', 'financials', 'configuration', 'lab-orders'
  let showUserProfile = false;
  let showLicenseGate = true;
  let licenseValidated = false;
  let showLogin = false;

  // Check license status and authentication on mount
  onMount(async () => {
    // Initialize theme
    document.body.setAttribute('data-theme', $theme);
    
    const savedLicense = localStorage.getItem('dentist_license_key');
    if (savedLicense) {
      const isValid = await validateCurrentLicense();
      if (isValid) {
        licenseValidated = true;
        showLicenseGate = false;
        
        // Check if user is already authenticated
        const wasAuthenticated = checkAuth();
        if (wasAuthenticated) {
          showLogin = false;
        } else {
          showLogin = true;
        }
      }
    }
  });

  // React to license validation changes
  $: if ($licenseValid && $currentLicenseKey) {
    licenseValidated = true;
    showLicenseGate = false;
    
    // If license is valid, check authentication
    if (!$isAuthenticated && !showLogin) {
      const wasAuthenticated = checkAuth();
      showLogin = !wasAuthenticated;
    }
  } else if (!$licenseValid) {
    licenseValidated = false;
    showLicenseGate = true;
    showLogin = false;
  }

  // React to authentication changes
  $: if ($isAuthenticated) {
    showLogin = false;
  } else if (licenseValidated && !$isAuthenticated) {
    showLogin = true;
  }

  function goTo(page) {
    currentPage = page;
  }
  
  function openUserProfile() {
    showUserProfile = true;
  }
  
  function closeUserProfile() {
    showUserProfile = false;
  }

  function switchTheme() {
    theme.update(t => t === 'dark' ? 'light' : 'dark');
    document.body.setAttribute('data-theme', $theme);
  }

  // Update theme attribute when theme changes
  $: if ($theme) {
    document.body.setAttribute('data-theme', $theme);
  }

  function onLicenseValidated(licenseKey, licenseInfo) {
    licenseValidated = true;
    showLicenseGate = false;
    // After license validation, show login
    showLogin = true;
    console.log('License validated successfully:', licenseInfo);
  }

  function onLoginSuccess(user) {
    showLogin = false;
    console.log('Login successful:', user);
  }
</script>

{#if showLicenseGate}
  <LicenseGate {onLicenseValidated} />
{:else if showLogin}
  <Login {onLoginSuccess} />
{:else if $isAuthenticated}
  <nav class="navbar">
    <div class="nav-tabs">
      <button class:active={currentPage === 'patients'} on:click={() => goTo('patients')}>Patient Management</button>
      <button class:active={currentPage === 'appointments'} on:click={() => goTo('appointments')}>Appointments</button>
      <button class:active={currentPage === 'calendar'} on:click={() => goTo('calendar')}>Calendar</button>
      <button class:active={currentPage === 'sessions'} on:click={() => goTo('sessions')}>Sessions</button>
      <button class:active={currentPage === 'payments'} on:click={() => goTo('payments')}>Payments</button>
      <button class:active={currentPage === 'financials'} on:click={() => goTo('financials')}>Financials</button>
      <button class:active={currentPage === 'lab-orders'} on:click={() => goTo('lab-orders')}>Lab Orders</button>
      <button class:active={currentPage === 'configuration'} on:click={() => goTo('configuration')}>Configuration</button>
    </div>
    <div class="nav-actions">
      <button class="theme-toggle-btn" on:click={switchTheme} title={$theme === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'} aria-label="Toggle Theme">
        {#if $theme === 'dark'}
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="5"/>
            <line x1="12" y1="1" x2="12" y2="3"/>
            <line x1="12" y1="21" x2="12" y2="23"/>
            <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
            <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
            <line x1="1" y1="12" x2="3" y2="12"/>
            <line x1="21" y1="12" x2="23" y2="12"/>
            <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
            <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
          </svg>
        {:else}
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
          </svg>
        {/if}
      </button>
      <button class="user-profile-btn" on:click={openUserProfile} title="User Profile" aria-label="User Profile">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
          <circle cx="12" cy="7" r="4"/>
        </svg>
      </button>
    </div>
  </nav>

  <UserProfileModal open={showUserProfile} onClose={closeUserProfile} />

  <main>
    {#if currentPage === 'patients'}
      {#if $selectedPatient}
        <PatientDetail patient={$selectedPatient} on:back={() => selectedPatient.set(null)} />
      {:else}
        <PatientList />
      {/if}
    {:else if currentPage === 'appointments'}
      <AppointmentList />
    {:else if currentPage === 'payments'}
      <Payments />
    {:else if currentPage === 'calendar'}
      <AppointmentCalendar />
    {:else if currentPage === 'sessions'}
      <Sessions />
    {:else if currentPage === 'financials'}
      <Financials />
    {:else if currentPage === 'lab-orders'}
      <LabOrders />
    {:else if currentPage === 'configuration'}
      <Configuration />
    {/if}
  </main>
{/if}

<style>
  :global(body) {
    margin: 0;
    color: white;
    font-family: "Nunito", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto",
    "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue",
    sans-serif;
  }

  main {
    margin: 0;
    padding: 0;
    width: 100%;
    overflow-x: hidden;
    box-sizing: border-box;
  }

  .navbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: var(--color-navbar);
    color: var(--color-text);
    padding: 0.5rem 2rem;
    gap: 1.5rem;
    box-shadow: var(--color-shadow);
    border-radius: 0 0 12px 12px;
    position: relative;
    width: 100%;
    box-sizing: border-box;
  }

  .nav-tabs {
    display: flex;
    gap: 1.5rem;
    align-items: center;
  }

  .nav-tabs button {
    background: none;
    border: none;
    color: var(--color-text);
    font-size: 1.1rem;
    font-weight: 600;
    padding: 0.5rem 1.2rem;
    border-radius: 8px 8px 0 0;
    cursor: pointer;
    transition: background 0.2s, color 0.2s;
  }

  .nav-tabs button.active, .nav-tabs button:hover {
    background: var(--color-accent);
    color: #fff;
  }

  .nav-actions {
    display: flex;
    gap: 0.5rem;
    align-items: center;
    margin-left: auto;
  }

  .theme-toggle-btn,
  .user-profile-btn {
    background: none;
    border: none;
    color: var(--color-text);
    padding: 0.5rem;
    border-radius: 8px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
    width: 40px;
    height: 40px;
  }

  .theme-toggle-btn:hover,
  .user-profile-btn:hover {
    background: var(--color-panel);
    transform: scale(1.05);
  }

  .theme-toggle-btn:active,
  .user-profile-btn:active {
    transform: scale(0.95);
  }
</style>
