<script>
  import { onMount } from 'svelte';

  let selectedSection = 'labs'; // 'labs', 'new-order', 'orders-list', 'tracking'

  onMount(() => {
    selectedSection = 'labs';
  });

  function selectSection(section) {
    selectedSection = section;
  }
</script>

<div class="lab-orders">
  <div class="lab-orders-container">
    <!-- Left Sidebar Navigation -->
    <aside class="sidebar">
      <h2 class="sidebar-title">Lab Orders</h2>
      <nav class="sidebar-nav">
        <button 
          class="nav-item" 
          class:active={selectedSection === 'labs'}
          on:click={() => selectSection('labs')}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M9.5 2A2.5 2.5 0 0 1 12 4.5v15a2.5 2.5 0 0 1-4.96.44L2.5 19.5"/>
            <path d="M14.5 2A2.5 2.5 0 0 0 12 4.5v15a2.5 2.5 0 0 0 4.96.44L21.5 19.5"/>
          </svg>
          <span>Labs Management</span>
        </button>
        
        <button 
          class="nav-item" 
          class:active={selectedSection === 'new-order'}
          on:click={() => selectSection('new-order')}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          <span>New Order</span>
        </button>
        
        <button 
          class="nav-item" 
          class:active={selectedSection === 'orders-list'}
          on:click={() => selectSection('orders-list')}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="8" y1="6" x2="21" y2="6"/>
            <line x1="8" y1="12" x2="21" y2="12"/>
            <line x1="8" y1="18" x2="21" y2="18"/>
            <line x1="3" y1="6" x2="3.01" y2="6"/>
            <line x1="3" y1="12" x2="3.01" y2="12"/>
            <line x1="3" y1="18" x2="3.01" y2="18"/>
          </svg>
          <span>Orders List</span>
        </button>
        
        <button 
          class="nav-item" 
          class:active={selectedSection === 'tracking'}
          on:click={() => selectSection('tracking')}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"/>
            <polyline points="12 6 12 12 16 14"/>
          </svg>
          <span>Order Tracking</span>
        </button>
      </nav>
    </aside>

    <!-- Right Content Panel -->
    <main class="content-panel">
      <!-- Labs Management Section -->
      {#if selectedSection === 'labs'}
        <div class="section-content">
          <div class="section-header">
            <h1>Lab Orders</h1>
            <p class="section-description">This section is under development.</p>
          </div>

          <div class="content-cards">
            <div class="action-card">
              <h3>Labs Management</h3>
              <p class="card-description">
                Features coming soon:
              </p>
              <ul class="feature-list">
                <li>Labs Management</li>
              </ul>
            </div>
          </div>
        </div>
      {/if}

      <!-- New Order Section -->
      {#if selectedSection === 'new-order'}
        <div class="section-content">
          <div class="section-header">
            <h1>Lab Orders</h1>
            <p class="section-description">This section is under development.</p>
          </div>

          <div class="content-cards">
            <div class="action-card">
              <h3>New Order</h3>
              <p class="card-description">
                Features coming soon:
              </p>
              <ul class="feature-list">
                <li>New Order Creation</li>
              </ul>
            </div>
          </div>
        </div>
      {/if}

      <!-- Orders List Section -->
      {#if selectedSection === 'orders-list'}
        <div class="section-content">
          <div class="section-header">
            <h1>Lab Orders</h1>
            <p class="section-description">This section is under development.</p>
          </div>

          <div class="content-cards">
            <div class="action-card">
              <h3>Orders List</h3>
              <p class="card-description">
                Features coming soon:
              </p>
              <ul class="feature-list">
                <li>Orders List & Tracking</li>
              </ul>
            </div>
          </div>
        </div>
      {/if}

      <!-- Order Tracking Section -->
      {#if selectedSection === 'tracking'}
        <div class="section-content">
          <div class="section-header">
            <h1>Lab Orders</h1>
            <p class="section-description">This section is under development.</p>
          </div>

          <div class="content-cards">
            <div class="action-card">
              <h3>Order Tracking</h3>
              <p class="card-description">
                Features coming soon:
              </p>
              <ul class="feature-list">
                <li>Order Status Monitoring</li>
              </ul>
            </div>
          </div>
        </div>
      {/if}
    </main>
  </div>
</div>

<style>
  .lab-orders {
    height: calc(100vh - 80px);
    background: var(--color-bg);
    color: var(--color-text);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .lab-orders-container {
    display: flex;
    height: 100%;
    overflow: hidden;
  }

  /* Sidebar Styles */
  .sidebar {
    width: 240px;
    background: var(--color-card);
    border-right: 1px solid var(--color-border);
    display: flex;
    flex-direction: column;
    padding: 1.5rem 0;
    overflow-y: auto;
  }

  .sidebar-title {
    font-size: 1.25rem;
    font-weight: 700;
    color: var(--color-text);
    padding: 0 1.5rem;
    margin: 0 0 1.5rem 0;
  }

  .sidebar-nav {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    padding: 0 0.75rem;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.625rem 1rem;
    border: none;
    background: transparent;
    color: var(--color-text);
    font-size: 0.9375rem;
    font-weight: 500;
    cursor: pointer;
    border-radius: 8px;
    transition: all 0.2s ease;
    text-align: left;
  }

  .nav-item:hover {
    background: var(--color-panel);
  }

  .nav-item.active {
    background: var(--color-accent);
    color: white;
  }

  .nav-item.danger.active {
    background: var(--color-danger);
    color: white;
  }

  .nav-item svg {
    flex-shrink: 0;
  }

  /* Content Panel Styles */
  .content-panel {
    flex: 1;
    overflow-y: auto;
    padding: 2rem;
    background: var(--color-bg);
  }

  .section-content {
    max-width: 900px;
    animation: fadeIn 0.3s ease;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .section-header {
    margin-bottom: 2rem;
  }

  .section-header h1 {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--color-text);
    margin: 0 0 0.5rem 0;
  }

  .section-description {
    font-size: 0.9375rem;
    color: var(--color-text);
    opacity: 0.7;
    margin: 0;
  }

  .content-cards {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }

  /* Card Styles */
  .action-card {
    background: var(--color-card);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: 1.5rem;
    box-shadow: var(--color-shadow);
  }

  .action-card h3 {
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--color-text);
    margin: 0 0 1rem 0;
  }

  .card-description {
    font-size: 0.875rem;
    color: var(--color-text);
    opacity: 0.7;
    margin: 0 0 1rem 0;
    line-height: 1.5;
  }

  .feature-list {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .feature-list li {
    padding: 0.5rem 0;
    color: var(--color-text);
    opacity: 0.8;
    font-size: 0.875rem;
    position: relative;
    padding-left: 1.5rem;
  }

  .feature-list li::before {
    content: "•";
    position: absolute;
    left: 0;
    color: var(--color-accent);
    font-weight: bold;
    font-size: 1.2rem;
  }
</style>

