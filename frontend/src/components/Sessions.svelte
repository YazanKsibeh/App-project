<script>
  import { onMount } from 'svelte';
  import { sessions, sessionsLoading, sessionsError, currentPage, totalPages, sessionsPageSize, loadSessions, loadSession, sessionFilters, applyFilters, clearAllFilters, clearProcedureFilters, removeFilter } from '../stores/sessionStore.js';
  import { currentUser } from '../stores/authStore.js';
  import { patients, loadPatients } from '../stores/patientStore.js';
  import { procedures, loadProcedures } from '../stores/procedureStore.js';
  import { GetAllUsers } from '../../wailsjs/go/main/App.js';
  import { currentLicenseKey } from '../stores/settingsStore.js';
import { createSession } from '../stores/sessionStore.js';
import { loadInvoiceOverview } from '../stores/financialsStore.js';
import { refreshInvoices } from '../stores/invoiceListStore.js';
  import SessionDetail from './SessionDetail.svelte';
  import InvoiceConfirmationModal from './InvoiceConfirmationModal.svelte';
  import { getInvoiceBySession } from '../stores/invoiceStore.js';
  import { get } from 'svelte/store';

  let selectedSession = null;
  let invoiceModalSessionId = null;
  let showInvoiceModal = false;
  let selectedSection = 'list'; // 'list', 'new', 'filters'

  // Filter state
  let filterSectionExpanded = true; // Default expanded
  let filterPatientSearch = '';
  let selectedFilterPatient = null; // Store selected patient object
  let filterProcedureSearch = '';
  let filterStatus = null; // null, 'completed', 'in-progress'
  let filterDentistId = null;
  let filterDateFrom = ''; // Empty by default
  let filterDateTo = ''; // Empty by default
  let selectedProcedures = []; // Array of procedure IDs
  let dentists = [];
  let showProcedureDropdown = false;

  // New Session Form State
  let selectedPatient = null;
  let patientSearch = '';
  let sessionDate = '';
  let sessionStatus = 'completed';
  let sessionNotes = '';
  let sessionItems = [];
  let procedureSearch = '';
  let showProcedureDropdownForm = false;
  let isSaving = false;

  onMount(async () => {
    console.log('[Sessions] Component mounted, loading sessions...');
    await loadSessions(1);
    await loadPatients();
    await loadProcedures();
    await loadDentists();
    console.log('[Sessions] Sessions loaded, count:', $sessions.length);
    
    // Set default date/time to now
    const now = new Date();
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
    sessionDate = now.toISOString().slice(0, 16);
    
    // No default date range for filters - keep empty
  });

  // Load dentists/users
  async function loadDentists() {
    try {
      const licenseKey = get(currentLicenseKey) || localStorage.getItem('dentist_license_key') || '';
      if (licenseKey) {
        dentists = await GetAllUsers(licenseKey) || [];
      }
    } catch (err) {
      console.error('Failed to load dentists:', err);
      dentists = [];
    }
  }
  
  $: {
    console.log('[Sessions] Sessions store updated, count:', $sessions.length);
    console.log('[Sessions] Sessions data:', $sessions);
    console.log('[Sessions] Is array?', Array.isArray($sessions));
  }

  function selectSection(section) {
    selectedSection = section;
    // Reset form when switching away from new session
    if (section !== 'new') {
      resetNewSessionForm();
    }
  }

  function resetNewSessionForm() {
    selectedPatient = null;
    patientSearch = '';
    const now = new Date();
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
    sessionDate = now.toISOString().slice(0, 16);
    sessionStatus = 'completed';
    sessionNotes = '';
    sessionItems = [];
    procedureSearch = '';
    showProcedureDropdownForm = false;
  }

  async function handlePageChange(page) {
    await loadSessions(page);
  }

  async function handleRowClick(session) {
    const fullSession = await loadSession(session.id);
    if (fullSession) {
      selectedSession = fullSession;
    }
  }

  function handleCloseDetail() {
    selectedSession = null;
  }

  function formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', { 
      year: 'numeric', 
      month: 'short', 
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }

  function formatCurrency(amount) {
    return amount.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
  }

  function getStatusIcon(status) {
    return status === 'completed' ? '●' : '⏳';
  }

  function getStatusClass(status) {
    return status === 'completed' ? 'status-completed' : 'status-in-progress';
  }

  function handleCreateInvoice(event, session) {
    event.stopPropagation(); // Prevent row click
    invoiceModalSessionId = session.id;
    showInvoiceModal = true;
  }

  function handleInvoiceModalClose() {
    showInvoiceModal = false;
    invoiceModalSessionId = null;
  }

  async function handleInvoiceConfirmed(event) {
    const { invoice } = event.detail;
    // Update the session in the list with the new invoice number
    if (invoiceModalSessionId) {
      // Reload sessions to get updated invoice numbers
      await loadSessions($currentPage, get(sessionFilters));
    }
    // Keep Financials overview cards in sync
    loadInvoiceOverview();
    refreshInvoices();
    showInvoiceModal = false;
    invoiceModalSessionId = null;
  }

  function handleInvoiceNumberClick(event, session) {
    event.stopPropagation(); // Prevent row click
    // Future: Navigate to invoice details
    console.log('View invoice:', session.invoice_number);
  }

  // Filter Handlers
  function toggleFilterSection() {
    filterSectionExpanded = !filterSectionExpanded;
  }

  // Filter patient search
  $: filteredFilterPatients = $patients.filter(p => 
    !filterPatientSearch || 
    p.name.toLowerCase().includes(filterPatientSearch.toLowerCase()) ||
    p.phone.includes(filterPatientSearch)
  );

  // Filter procedure search
  $: filteredFilterProcedures = $procedures.filter(p =>
    !filterProcedureSearch ||
    p.name.toLowerCase().includes(filterProcedureSearch.toLowerCase())
  );

  function selectFilterPatient(patient) {
    selectedFilterPatient = patient;
    filterPatientSearch = patient.name; // Only show name in input
  }

  function selectFilterProcedure(procedure) {
    if (!selectedProcedures.includes(procedure.id)) {
      selectedProcedures = [...selectedProcedures, procedure.id];
    }
    filterProcedureSearch = '';
    showProcedureDropdownForm = false;
  }

  function removeSelectedProcedure(procedureId) {
    selectedProcedures = selectedProcedures.filter(id => id !== procedureId);
  }

  async function clearFilterInputs() {
    filterPatientSearch = '';
    selectedFilterPatient = null;
    filterProcedureSearch = '';
    filterStatus = null;
    filterDentistId = null;
    filterDateFrom = '';
    filterDateTo = '';
    selectedProcedures = [];
    showProcedureDropdown = false;
    // Auto-apply after clearing
    await applyFilterFilters();
  }

  async function applyFilterFilters() {
    // Use selected patient if available
    const patientId = selectedFilterPatient ? selectedFilterPatient.id : null;

    const filters = {
      patient_id: patientId,
      status: filterStatus,
      dentist_id: filterDentistId,
      date_from: filterDateFrom || null,
      date_to: filterDateTo || null,
      procedure_ids: selectedProcedures
    };

    await applyFilters(filters);
    
    // After filter is applied, clear patient input display (but keep selectedFilterPatient for reference)
    // The input should show placeholder when filter is active in store
    if (patientId !== null && $sessionFilters.patient_id === patientId) {
      filterPatientSearch = '';
      // Keep selectedFilterPatient - it's needed for state management and will be cleared when filter is removed
    }
  }

  // Remove a specific filter (wrapper to update local state)
  async function handleRemoveFilter(filterType, value = null) {
    // Update local state first
    switch (filterType) {
      case 'patient':
        filterPatientSearch = '';
        selectedFilterPatient = null;
        break;
      case 'status':
        filterStatus = null;
        // Preserve patient selection when removing status
        break;
      case 'dentist':
        filterDentistId = null;
        // Preserve patient selection when removing dentist
        break;
      case 'date':
        filterDateFrom = '';
        filterDateTo = '';
        // Preserve patient selection when removing date
        break;
      case 'procedure':
        if (value !== null) {
          selectedProcedures = selectedProcedures.filter(id => id !== value);
        }
        // Preserve patient selection when removing procedure
        break;
    }
    // Then call store function
    await removeFilter(filterType, value);
  }

  async function handleClearAllFilters() {
    clearFilterInputs();
    await clearAllFilters();
  }

  async function handleClearProcedureFilters() {
    selectedProcedures = [];
    await clearProcedureFilters();
  }

  // Get active filters for display (from store, not local state)
  $: activeFilters = {
    patient: $sessionFilters.patient_id ? $patients.find(p => p.id === $sessionFilters.patient_id) : null,
    status: $sessionFilters.status,
    dentist: $sessionFilters.dentist_id ? dentists.find(d => d.id === $sessionFilters.dentist_id) : null,
    date: $sessionFilters.date_from && $sessionFilters.date_to ? { from: $sessionFilters.date_from, to: $sessionFilters.date_to } : null,
    procedures: ($sessionFilters.procedure_ids || []).map(id => $procedures.find(p => p.id === id)).filter(Boolean)
  };

  $: hasActiveFilters = activeFilters.patient || activeFilters.status || activeFilters.dentist || 
    activeFilters.date || (activeFilters.procedures && activeFilters.procedures.length > 0);

  // Clear patient input display ONLY when filter is applied via Apply button
  // This reactive statement should NOT interfere with local state before Apply is pressed
  // The input will be cleared in applyFilterFilters() after the filter is applied

  // New Session Form Handlers
  $: filteredPatients = $patients.filter(p => 
    !patientSearch || 
    p.name.toLowerCase().includes(patientSearch.toLowerCase()) ||
    p.phone.includes(patientSearch)
  );

  $: filteredProcedures = $procedures.filter(p =>
    !procedureSearch ||
    p.name.toLowerCase().includes(procedureSearch.toLowerCase())
  );

  // Separate procedure search for form
  $: filteredFormProcedures = $procedures.filter(p =>
    !procedureSearch ||
    p.name.toLowerCase().includes(procedureSearch.toLowerCase())
  );

  $: totalAmount = sessionItems.reduce((sum, item) => sum + (item.amount || 0), 0);

  function handlePatientSelect(patient) {
    selectedPatient = patient;
    patientSearch = patient.name;
    showProcedureDropdownForm = false;
  }

  function handlePatientSearchInput(e) {
    patientSearch = e.target.value;
    if (!patientSearch) {
      selectedPatient = null;
    }
  }

  function handleAddProcedure() {
    showProcedureDropdownForm = true;
  }

  function handleSelectProcedure(procedure) {
    sessionItems = [...sessionItems, {
      procedure_id: procedure.id || null,
      item_name: procedure.name,
      amount: procedure.price
    }];
    procedureSearch = '';
    showProcedureDropdownForm = false;
  }

  function handleRemoveItem(index) {
    sessionItems = sessionItems.filter((_, i) => i !== index);
  }

  function handleUpdateItemAmount(index, value) {
    const amount = parseInt(value) || 0;
    sessionItems = sessionItems.map((item, i) => 
      i === index ? { ...item, amount } : item
    );
  }

  function handleUpdateItemName(index, value) {
    sessionItems = sessionItems.map((item, i) => 
      i === index ? { ...item, item_name: value } : item
    );
  }

  async function handleSaveSession() {
    if (!selectedPatient) {
      alert('Please select a patient');
      return;
    }

    if (!sessionDate) {
      alert('Please select a date and time');
      return;
    }

    if (sessionItems.length === 0) {
      alert('Please add at least one procedure');
      return;
    }

    isSaving = true;
    try {
      const user = get(currentUser);
      
      const sessionForm = {
        patient_id: selectedPatient.id,
        dentist_id: user?.id || 1,
        session_date: new Date(sessionDate).toISOString(),
        status: sessionStatus,
        notes: sessionNotes,
        items: sessionItems.map(item => ({
          procedure_id: item.procedure_id ? item.procedure_id : null,
          item_name: item.item_name,
          amount: item.amount
        }))
      };

      const success = await createSession(sessionForm);
      
      if (success) {
        resetNewSessionForm();
        selectedSection = 'list';
        await loadSessions(1);
      }
    } catch (error) {
      console.error('Error creating session:', error);
      alert('Failed to create session: ' + (error.message || 'Unknown error'));
    } finally {
      isSaving = false;
    }
  }
</script>

<div class="sessions">
  <div class="sessions-container">
    <!-- Left Sidebar Navigation -->
    <aside class="sidebar">
      <h2 class="sidebar-title">Sessions</h2>
      <nav class="sidebar-nav">
        <button 
          class="nav-item" 
          class:active={selectedSection === 'list'}
          on:click={() => selectSection('list')}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="8" y1="6" x2="21" y2="6"/>
            <line x1="8" y1="12" x2="21" y2="12"/>
            <line x1="8" y1="18" x2="21" y2="18"/>
            <line x1="3" y1="6" x2="3.01" y2="6"/>
            <line x1="3" y1="12" x2="3.01" y2="12"/>
            <line x1="3" y1="18" x2="3.01" y2="18"/>
          </svg>
          <span>Sessions List</span>
        </button>
        
        <button 
          class="nav-item" 
          class:active={selectedSection === 'new'}
          on:click={() => selectSection('new')}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          <span>New Session</span>
        </button>
        
        <button 
          class="nav-item" 
          class:active={selectedSection === 'filters'}
          on:click={() => selectSection('filters')}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/>
          </svg>
          <span>Filters & Reports</span>
        </button>
      </nav>
    </aside>

    <!-- Right Content Panel -->
    <main class="content-panel">
      <!-- Sessions List Section -->
      {#if selectedSection === 'list'}
        <div class="section-content">
          <div class="section-header">
            <div>
              <h1>Sessions List</h1>
              <p class="section-description">View and manage all patient sessions</p>
            </div>
            <button 
              class="filter-toggle-btn" 
              on:click={toggleFilterSection}
              title={filterSectionExpanded ? 'Collapse filters' : 'Expand filters'}
            >
              Filter {filterSectionExpanded ? '↓' : '→'}
            </button>
          </div>

          <!-- Filter Section -->
          {#if filterSectionExpanded}
            <div class="filter-section">
              <div class="filter-two-column">
                <!-- Left Column: Standard Filters -->
                <div class="filter-left-column">
                  <!-- Row 1: Patient and Dentist side by side -->
                  <div class="filter-row">
                    <div class="filter-group filter-group-half">
                      <label>Patient</label>
                      <div class="search-container">
                        <input
                          type="text"
                          class="form-input"
                          placeholder="Patient..."
                          value={selectedFilterPatient ? selectedFilterPatient.name : filterPatientSearch}
                          on:input={(e) => {
                            if (!$sessionFilters.patient_id) {
                              filterPatientSearch = e.target.value;
                              // Clear selectedFilterPatient if user starts typing
                              if (selectedFilterPatient && e.target.value !== selectedFilterPatient.name) {
                                selectedFilterPatient = null;
                              }
                            }
                          }}
                          disabled={$sessionFilters.patient_id !== null}
                        />
                        {#if filterPatientSearch && !$sessionFilters.patient_id && !selectedFilterPatient && filteredFilterPatients.length > 0}
                          <div class="dropdown">
                            {#each filteredFilterPatients.slice(0, 10) as patient}
                              <div class="dropdown-item" on:click={() => selectFilterPatient(patient)}>
                                {patient.name} {patient.phone ? `(${patient.phone})` : ''}
                              </div>
                            {/each}
                          </div>
                        {/if}
                      </div>
                    </div>

                    <div class="filter-group filter-group-half">
                      <label>Dentist</label>
                      <select class="form-input" bind:value={filterDentistId}>
                        <option value={null}>All Dentists...</option>
                        {#each dentists as dentist}
                          <option value={dentist.id}>{dentist.username}</option>
                        {/each}
                      </select>
                    </div>
                  </div>

                  <!-- Row 2: Status Filter (single line) -->
                  <div class="filter-row">
                    <div class="filter-group filter-group-full">
                      <label>Status</label>
                      <div class="status-radio-group">
                        <label class="radio-option">
                          <input
                            type="radio"
                            name="filterStatus"
                            value=""
                            checked={filterStatus === null}
                            on:change={() => filterStatus = null}
                          />
                          <span>All</span>
                        </label>
                        <label class="radio-option radio-completed">
                          <input
                            type="radio"
                            name="filterStatus"
                            value="completed"
                            checked={filterStatus === 'completed'}
                            on:change={() => filterStatus = 'completed'}
                          />
                          <span>Completed</span>
                        </label>
                        <label class="radio-option radio-in-progress">
                          <input
                            type="radio"
                            name="filterStatus"
                            value="in-progress"
                            checked={filterStatus === 'in-progress'}
                            on:change={() => filterStatus = 'in-progress'}
                          />
                          <span>In-progress</span>
                        </label>
                      </div>
                    </div>
                  </div>

                  <!-- Row 3: Date Range (From and To side by side) -->
                  <div class="filter-row">
                    <div class="filter-group filter-group-half">
                      <label>From:</label>
                      <input
                        type="date"
                        class="form-input"
                        bind:value={filterDateFrom}
                        max={filterDateTo}
                        placeholder="Date..."
                      />
                    </div>

                    <div class="filter-group filter-group-half">
                      <label>To:</label>
                      <input
                        type="date"
                        class="form-input"
                        bind:value={filterDateTo}
                        min={filterDateFrom}
                        placeholder="Date..."
                      />
                    </div>
                  </div>
                </div>

                <!-- Right Column: Procedure Selection -->
                <div class="filter-right-column">
                  <div class="filter-procedure-section">
                    <label class="section-label">Procedures</label>
                    <div class="procedure-search-container">
                      <input
                        type="text"
                        class="form-input procedure-search"
                        placeholder="Search procedures..."
                        bind:value={filterProcedureSearch}
                      />
                    </div>
                    <div class="procedures-grid-filter">
                      {#if filteredFilterProcedures.length > 0}
                        {#each filteredFilterProcedures as procedure}
                          <button
                            class="procedure-tag {selectedProcedures.includes(procedure.id) ? 'selected' : ''}"
                            on:click={async () => {
                              if (selectedProcedures.includes(procedure.id)) {
                                selectedProcedures = selectedProcedures.filter(id => id !== procedure.id);
                              } else {
                                selectedProcedures = [...selectedProcedures, procedure.id];
                              }
                              // Auto-apply after selection change
                              await applyFilterFilters();
                            }}
                            type="button"
                          >
                            {procedure.name}
                          </button>
                        {/each}
                      {:else}
                        <div class="empty-procedures">No procedures available</div>
                      {/if}
                    </div>
                  </div>
                </div>
              </div>

              <!-- Action Buttons (below filter section, right-aligned) -->
              <div class="filter-actions-bottom">
                <button class="btn btn-secondary" on:click={clearFilterInputs}>
                  Clear
                </button>
                <button class="btn btn-primary" on:click={applyFilterFilters}>
                  Apply
                </button>
              </div>
            </div>
          {/if}

          <!-- Active Filters Bar -->
          {#if hasActiveFilters}
            <div class="active-filters-bar">
              <div class="active-filters-main">
                {#if activeFilters.patient}
                  <span class="filter-tag">
                    Patient: {activeFilters.patient.name}
                    <button class="filter-tag-remove" on:click={() => handleRemoveFilter('patient')}>×</button>
                  </span>
                {/if}
                {#if activeFilters.status}
                  <span class="filter-tag">
                    Status: {activeFilters.status === 'completed' ? 'Completed' : 'In-progress'}
                    <button class="filter-tag-remove" on:click={() => handleRemoveFilter('status')}>×</button>
                  </span>
                {/if}
                {#if activeFilters.dentist}
                  <span class="filter-tag">
                    Dentist: {activeFilters.dentist.username}
                    <button class="filter-tag-remove" on:click={() => handleRemoveFilter('dentist')}>×</button>
                  </span>
                {/if}
                {#if activeFilters.date}
                  <span class="filter-tag">
                    Date: {new Date(activeFilters.date.from).toLocaleDateString('en-US', { month: 'short', year: 'numeric' })} - {new Date(activeFilters.date.to).toLocaleDateString('en-US', { month: 'short', year: 'numeric' })}
                    <button class="filter-tag-remove" on:click={() => handleRemoveFilter('date')}>×</button>
                  </span>
                {/if}
              </div>
              {#if activeFilters.procedures && activeFilters.procedures.length > 0}
                <div class="active-filters-procedures">
                  {#each activeFilters.procedures as procedure}
                    <span class="filter-tag">
                      {procedure.name}
                      <button class="filter-tag-remove" on:click={() => handleRemoveFilter('procedure', procedure.id)}>×</button>
                    </span>
                  {/each}
                </div>
              {/if}
              <div class="active-filters-actions">
                {#if activeFilters.procedures && activeFilters.procedures.length > 0}
                  <button class="btn-filter-action" on:click={handleClearProcedureFilters}>Clear Procedures</button>
                {/if}
                <button class="btn-filter-action" on:click={handleClearAllFilters}>Clear All Filters</button>
              </div>
            </div>
          {/if}

          {#if $sessionsLoading}
            <div class="loading-state">
              <div class="spinner"></div>
              <p>Loading sessions...</p>
            </div>
          {:else if $sessionsError}
            <div class="error-state">
              <p>{$sessionsError}</p>
            </div>
          {:else if $sessions.length === 0}
            <div class="empty-state">
              <p>No sessions found. Create your first session to get started.</p>
            </div>
          {:else}
            <div class="sessions-table-container">
              <table class="sessions-table">
                <thead>
                  <tr>
                    <th>Patient Name</th>
                    <th>Total Cost</th>
                    <th>Date</th>
                    <th>Status</th>
                    <th>Invoice</th>
                  </tr>
                </thead>
                <tbody>
                  {#if Array.isArray($sessions) && $sessions.length > 0}
                    {#each $sessions as session (session.id)}
                      <tr class="session-row" on:click={() => handleRowClick(session)}>
                        <td class="patient-name">{session.patient_name || 'Unknown'}</td>
                        <td class="total-cost">{formatCurrency(session.total_amount)} SYP</td>
                        <td class="date">{formatDate(session.session_date)}</td>
                        <td class="status">
                          <span class="status-badge {getStatusClass(session.status)}">
                            <span class="status-icon">{getStatusIcon(session.status)}</span>
                            <span class="status-text">{session.status === 'completed' ? 'Completed' : 'In-progress'}</span>
                          </span>
                        </td>
                        <td class="invoice" on:click|stopPropagation>
                          {#if session.invoice_number}
                            <button 
                              class="invoice-link" 
                              on:click={(e) => handleInvoiceNumberClick(e, session)}
                              title="View invoice"
                            >
                              {session.invoice_number}
                            </button>
                          {:else}
                            <button 
                              class="btn-create-invoice" 
                              on:click={(e) => handleCreateInvoice(e, session)}
                            >
                              Create Invoice
                            </button>
                          {/if}
                        </td>
                      </tr>
                    {/each}
                  {:else}
                    <tr>
                      <td colspan="5" style="text-align: center; padding: 2rem; color: var(--color-text); opacity: 0.6;">
                        {#if !Array.isArray($sessions)}
                          Error: Sessions data is not an array. Type: {typeof $sessions}
                        {:else}
                          No sessions found
                        {/if}
                      </td>
                    </tr>
                  {/if}
                </tbody>
              </table>
            </div>

            <div class="pagination">
              <button 
                class="page-btn" 
                disabled={$currentPage === 1}
                on:click={() => handlePageChange($currentPage - 1)}
              >
                Previous
              </button>
              <div class="page-meta">
                <span class="page-info">
                  Page {$currentPage} of {$totalPages}
                </span>
                <span class="page-size">
                  {$sessionsPageSize} per page
                </span>
              </div>
              <button 
                class="page-btn" 
                disabled={$currentPage >= $totalPages}
                on:click={() => handlePageChange($currentPage + 1)}
              >
                Next
              </button>
            </div>
          {/if}
        </div>
      {/if}

      <!-- New Session Section -->
      {#if selectedSection === 'new'}
        <div class="section-content">
          <div class="section-header">
            <h1>New Patient Session</h1>
          </div>

          <div class="new-session-form">
            <!-- Top Row: 3-Column Layout -->
            <div class="form-top-row">
              <!-- Patient Section -->
              <div class="form-section">
                <label class="section-label">Patient *</label>
                <div class="patient-search-container">
                  <svg class="search-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="11" cy="11" r="8"/>
                    <path d="m21 21-4.35-4.35"/>
                  </svg>
                  <input
                    type="text"
                    class="form-input search-input"
                    placeholder="Search by name or phone..."
                    value={patientSearch}
                    on:input={handlePatientSearchInput}
                    on:focus={() => showProcedureDropdownForm = false}
                  />
                  {#if patientSearch && !selectedPatient}
                    <div class="dropdown">
                      {#if filteredPatients.length > 0}
                        {#each filteredPatients.slice(0, 5) as patient}
                          <div 
                            class="dropdown-item"
                            on:click={() => handlePatientSelect(patient)}
                          >
                            <div class="patient-name">{patient.name}</div>
                            <div class="patient-phone">{patient.phone}</div>
                          </div>
                        {/each}
                      {:else}
                        <div class="dropdown-item empty">No patients found</div>
                      {/if}
                    </div>
                  {/if}
                </div>
                {#if selectedPatient}
                  <div class="selected-patient-badge">
                    <span>{selectedPatient.name} ({selectedPatient.phone})</span>
                  </div>
                {/if}
              </div>

              <!-- Date & Time Section -->
              <div class="form-section">
                <label class="section-label">Date & Time *</label>
                <div class="date-input-container">
                  <svg class="calendar-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                    <line x1="16" y1="2" x2="16" y2="6"/>
                    <line x1="8" y1="2" x2="8" y2="6"/>
                    <line x1="3" y1="10" x2="21" y2="10"/>
                  </svg>
                  <input
                    type="datetime-local"
                    class="form-input date-input"
                    bind:value={sessionDate}
                  />
                </div>
              </div>

              <!-- Status Section -->
              <div class="form-section">
                <label class="section-label">Status</label>
                <div class="status-buttons">
                  <button
                    class="status-btn"
                    class:active={sessionStatus === 'completed'}
                    on:click={() => sessionStatus = 'completed'}
                    type="button"
                  >
                    Completed
                  </button>
                  <button
                    class="status-btn"
                    class:active={sessionStatus === 'in-progress'}
                    on:click={() => sessionStatus = 'in-progress'}
                    type="button"
                  >
                    In-progress
                  </button>
                </div>
              </div>
            </div>

            <!-- Procedures Section - Two Column Layout -->
            <div class="form-section procedures-section">
              <div class="procedures-container">
                <!-- Available Procedures Column -->
                <div class="procedures-column available-procedures">
                  <label class="section-label">Procedures *</label>
                  <div class="procedure-search-container">
                    <input
                      type="text"
                      class="form-input procedure-search"
                      placeholder="Search procedures..."
                      bind:value={procedureSearch}
                    />
                  </div>
                  <div class="procedures-grid">
                    {#if filteredProcedures.length > 0}
                      {#each filteredProcedures as procedure}
                        <button
                          class="procedure-tag"
                          on:click={() => handleSelectProcedure(procedure)}
                          type="button"
                        >
                          {procedure.name}
                        </button>
                      {/each}
                    {:else}
                      <div class="empty-procedures">No procedures available</div>
                    {/if}
                  </div>
                </div>

                <!-- Selected Procedures Column -->
                <div class="procedures-column selected-procedures">
                  <label class="section-label">Selected Procedures</label>
                  {#if sessionItems.length > 0}
                    <div class="selected-procedures-list">
                      {#each sessionItems as item, index}
                        <div class="selected-procedure-item">
                          <div class="procedure-info">
                            <span class="procedure-name">{item.item_name}</span>
                            <span class="procedure-price">{formatCurrency(item.amount)} SYP</span>
                          </div>
                          <button
                            class="btn-remove-item"
                            on:click={() => handleRemoveItem(index)}
                            type="button"
                            title="Remove"
                          >
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                              <line x1="18" y1="6" x2="6" y2="18"/>
                              <line x1="6" y1="6" x2="18" y2="18"/>
                            </svg>
                          </button>
                        </div>
                      {/each}
                    </div>
                    <div class="total-display">
                      <div class="total-label">Total Amount</div>
                      <div class="total-value">{formatCurrency(totalAmount)} SYP</div>
                    </div>
                  {:else}
                    <div class="empty-selected">
                      <p>No procedures selected</p>
                      <p class="hint">Click on procedures from the left to add them</p>
                    </div>
                  {/if}
                </div>
              </div>
            </div>

            <!-- Session Notes Section -->
            <div class="form-section notes-section">
              <label class="section-label">Session Notes</label>
              <textarea
                class="form-textarea notes-textarea"
                rows="3"
                placeholder="Add details about the session..."
                bind:value={sessionNotes}
              ></textarea>
            </div>

            <!-- Final Amount Display and Action Buttons -->
            <div class="form-footer">
              <div class="final-amount-section">
                <div class="final-amount-label">Final Amount</div>
                <div class="final-amount-value">{formatCurrency(totalAmount)} SYP</div>
              </div>
              <div class="form-actions">
                <button class="btn btn-secondary" on:click={() => selectSection('list')} disabled={isSaving}>
                  Cancel
                </button>
                <button class="btn btn-primary" on:click={handleSaveSession} disabled={isSaving}>
                  {#if isSaving}
                    Saving...
                  {:else}
                    Save Session
                  {/if}
                </button>
              </div>
            </div>
          </div>
        </div>
      {/if}

      <!-- Filters & Reports Section -->
      {#if selectedSection === 'filters'}
        <div class="section-content">
          <div class="section-header">
            <h1>Filters & Reports</h1>
            <p class="section-description">Filter sessions and generate reports</p>
          </div>

          <div class="content-cards">
            <div class="info-card">
              <h3>Coming Soon</h3>
              <p class="card-description">Filters and reporting features will be available in a future update.</p>
            </div>
          </div>
        </div>
      {/if}
    </main>
  </div>
</div>

{#if selectedSession}
  <SessionDetail session={selectedSession} on:close={handleCloseDetail} />
{/if}

{#if showInvoiceModal && invoiceModalSessionId}
  <InvoiceConfirmationModal 
    sessionId={invoiceModalSessionId}
    bind:open={showInvoiceModal}
    on:close={handleInvoiceModalClose}
    on:confirmed={handleInvoiceConfirmed}
  />
{/if}

<style>
  .sessions {
    height: 100%;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .sessions-container {
    display: flex;
    flex: 1;
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
  .form-card,
  .info-card {
    background: var(--color-card);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: 1.5rem;
    box-shadow: var(--color-shadow);
  }

  .info-card h3 {
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--color-text);
    margin: 0 0 1rem 0;
  }

  .card-description {
    font-size: 0.875rem;
    color: var(--color-text);
    opacity: 0.7;
    margin: 0;
    line-height: 1.5;
  }

  /* New Session Form Styles */
  .new-session-form {
    background: var(--color-card);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: 2rem;
    box-shadow: var(--color-shadow);
  }

  .form-section {
    margin-bottom: 1.5rem;
  }

  /* Top Row: 3-Column Layout */
  .form-top-row {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    gap: 1.5rem;
    margin-bottom: 2rem;
  }

  @media (max-width: 768px) {
    .form-top-row {
      grid-template-columns: 1fr;
      gap: 1.5rem;
    }
  }

  .section-label {
    display: block;
    margin-bottom: 0.75rem;
    font-weight: 600;
    color: var(--color-text);
    font-size: 0.9375rem;
  }

  /* Patient Search */
  .patient-search-container {
    position: relative;
    display: flex;
    align-items: center;
  }

  .search-icon {
    position: absolute;
    left: 1rem;
    color: var(--color-text);
    opacity: 0.5;
    pointer-events: none;
    z-index: 1;
  }

  .search-input {
    width: 100%;
    padding: 0.875rem 1rem 0.875rem 2.75rem;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-size: 1rem;
    font-family: inherit;
    box-sizing: border-box;
  }

  .search-input::placeholder {
    color: var(--color-text);
    opacity: 0.5;
  }

  .search-input:focus {
    outline: none;
    border-color: var(--color-accent);
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  .selected-patient-badge {
    margin-top: 0.75rem;
    padding: 0.625rem 1rem;
    background: var(--color-panel);
    border-radius: 6px;
    display: inline-block;
    font-size: 0.875rem;
    color: var(--color-text);
  }

  /* Date Input */
  .date-input-container {
    position: relative;
    display: flex;
    align-items: center;
  }

  .calendar-icon {
    position: absolute;
    left: 1rem;
    color: var(--color-text);
    opacity: 0.5;
    pointer-events: none;
    z-index: 1;
  }

  .date-input {
    width: 100%;
    padding: 0.875rem 1rem 0.875rem 2.75rem;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-size: 1rem;
    font-family: inherit;
    box-sizing: border-box;
  }

  .date-input::-webkit-calendar-picker-indicator {
    opacity: 0;
    position: absolute;
    right: 0;
    width: 100%;
    height: 100%;
    cursor: pointer;
  }

  .date-input:focus {
    outline: none;
    border-color: var(--color-accent);
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  /* Status Buttons */
  .status-buttons {
    display: flex;
    gap: 0.75rem;
  }

  .status-btn {
    flex: 1;
    padding: 0.875rem 1.5rem;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-size: 0.9375rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .status-btn:hover {
    background: var(--color-border);
  }

  .status-btn.active {
    background: var(--color-accent);
    color: white;
    border-color: var(--color-accent);
  }

  /* Procedures Section */
  .procedures-section {
    margin-bottom: 2rem;
  }

  .procedures-container {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
    margin-top: 1rem;
  }

  @media (max-width: 768px) {
    .procedures-container {
      grid-template-columns: 1fr;
      gap: 1.5rem;
    }
  }

  .procedures-column {
    display: flex;
    flex-direction: column;
  }

  .available-procedures {
    background: var(--color-bg);
    padding: 1.5rem;
    border-radius: 8px;
    border: 1px solid var(--color-border);
  }

  .selected-procedures {
    background: var(--color-panel);
    padding: 1.5rem;
    border-radius: 8px;
    border: 1px solid var(--color-border);
  }

  .procedure-search-container {
    margin-bottom: 1rem;
  }

  .procedure-search {
    width: 100%;
    padding: 0.75rem;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-size: 0.9375rem;
    font-family: inherit;
    box-sizing: border-box;
  }

  .procedure-search::placeholder {
    color: var(--color-text);
    opacity: 0.5;
  }

  .procedure-search:focus {
    outline: none;
    border-color: var(--color-accent);
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  .procedures-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    max-height: 400px;
    overflow-y: auto;
  }

  .procedure-tag {
    padding: 0.625rem 1rem;
    background: var(--color-card);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 6px;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    white-space: nowrap;
  }

  .procedure-tag:hover:not(:disabled) {
    background: var(--color-panel);
    border-color: var(--color-accent);
    transform: translateY(-1px);
  }

  .procedure-tag:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .procedure-tag.selected {
    background: var(--color-accent);
    color: white;
    border-color: var(--color-accent);
  }

  .empty-procedures {
    padding: 2rem;
    text-align: center;
    color: var(--color-text);
    opacity: 0.6;
    font-size: 0.875rem;
  }

  /* Selected Procedures */
  .selected-procedures-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    margin-bottom: 1.5rem;
    max-height: 300px;
    overflow-y: auto;
  }

  .selected-procedure-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.875rem;
    background: var(--color-card);
    border-radius: 6px;
    border: 1px solid var(--color-border);
  }

  .procedure-info {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    flex: 1;
  }

  .procedure-info .procedure-name {
    font-weight: 500;
    color: var(--color-text);
    font-size: 0.9375rem;
  }

  .procedure-info .procedure-price {
    font-weight: 600;
    color: var(--color-accent);
    font-size: 0.875rem;
  }

  .btn-remove-item {
    background: transparent;
    border: none;
    color: var(--color-text);
    opacity: 0.6;
    cursor: pointer;
    padding: 0.25rem;
    border-radius: 4px;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .btn-remove-item:hover {
    opacity: 1;
    background: rgba(239, 68, 68, 0.1);
    color: #ef4444;
  }

  .total-display {
    padding-top: 1rem;
    border-top: 1px solid var(--color-border);
  }

  .total-label {
    font-size: 0.875rem;
    color: var(--color-text);
    opacity: 0.7;
    margin-bottom: 0.5rem;
  }

  .total-value {
    font-size: 1.25rem;
    font-weight: 700;
    color: var(--color-accent);
  }

  .empty-selected {
    padding: 2rem;
    text-align: center;
    color: var(--color-text);
    opacity: 0.6;
  }

  .empty-selected p {
    margin: 0.5rem 0;
    font-size: 0.875rem;
  }

  .empty-selected .hint {
    font-size: 0.8125rem;
    opacity: 0.5;
  }

  /* Session Notes */
  .notes-section {
    margin-bottom: 1.5rem;
  }

  .notes-textarea {
    width: 100%;
    padding: 0.875rem;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-size: 1rem;
    font-family: inherit;
    box-sizing: border-box;
    resize: vertical;
  }

  .notes-textarea::placeholder {
    color: var(--color-text);
    opacity: 0.5;
  }

  .notes-textarea:focus {
    outline: none;
    border-color: var(--color-accent);
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  /* Form Footer */
  .form-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 2rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--color-border);
    margin-top: 1rem;
  }

  @media (max-width: 768px) {
    .form-footer {
      flex-direction: column;
      align-items: stretch;
      gap: 1.5rem;
    }
  }

  /* Final Amount Display */
  .final-amount-section {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .final-amount-label {
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--color-text);
    opacity: 0.7;
  }

  .final-amount-value {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--color-accent);
  }

  /* Dropdown Styles */
  .dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: var(--color-card);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    margin-top: 0.25rem;
    max-height: 200px;
    overflow-y: auto;
    z-index: 10;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  .dropdown-item {
    padding: 0.75rem;
    cursor: pointer;
    border-bottom: 1px solid var(--color-border);
    transition: background 0.2s;
  }

  .dropdown-item:last-child {
    border-bottom: none;
  }

  .dropdown-item:hover {
    background: var(--color-panel);
  }

  .dropdown-item.empty {
    cursor: default;
    color: var(--color-text);
    opacity: 0.6;
  }

  .patient-name {
    font-weight: 500;
    color: var(--color-text);
  }

  .patient-phone {
    font-size: 0.875rem;
    color: var(--color-text);
    opacity: 0.7;
    margin-top: 0.25rem;
  }

  .form-actions {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
  }

  @media (max-width: 768px) {
    .form-actions {
      width: 100%;
    }
    
    .form-actions .btn {
      flex: 1;
    }
  }

  .btn {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .btn-primary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
  }

  .btn-primary:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  }

  .btn-secondary {
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
  }

  .btn-secondary:hover:not(:disabled) {
    background: var(--color-border);
  }

  /* Sessions Table Styles */
  .loading-state,
  .error-state,
  .empty-state {
    text-align: center;
    padding: 3rem;
    color: var(--color-text);
  }

  .spinner {
    width: 40px;
    height: 40px;
    border: 3px solid var(--color-border);
    border-top: 3px solid var(--color-accent);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin: 0 auto 1rem;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  .sessions-table-container {
    background: var(--color-card);
    border-radius: 12px;
    border: 1px solid var(--color-border);
    overflow: hidden;
  }

  .sessions-table {
    width: 100%;
    border-collapse: collapse;
  }

  .sessions-table thead {
    background: var(--color-panel);
    border-bottom: 2px solid var(--color-border);
  }

  .sessions-table th {
    padding: 1rem;
    text-align: left;
    font-weight: 600;
    color: var(--color-text);
    font-size: 0.9rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .sessions-table tbody tr {
    border-bottom: 1px solid var(--color-border);
    cursor: pointer;
    transition: background 0.2s ease;
  }

  .sessions-table tbody tr:hover {
    background: var(--color-panel);
  }

  .sessions-table tbody tr:last-child {
    border-bottom: none;
  }

  .sessions-table td {
    padding: 1rem;
    color: var(--color-text);
  }

  .patient-name {
    font-weight: 500;
  }

  .total-cost {
    font-weight: 600;
    color: var(--color-accent);
  }

  .status-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.375rem 0.75rem;
    border-radius: 6px;
    font-size: 0.875rem;
    font-weight: 500;
  }

  .status-completed {
    background: rgba(34, 197, 94, 0.1);
    color: #22c55e;
  }

  .status-in-progress {
    background: rgba(251, 191, 36, 0.1);
    color: #fbbf24;
  }

  .status-icon {
    font-size: 0.75rem;
  }

  .pagination {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 1rem;
    margin-top: 2rem;
    flex-wrap: wrap;
  }

  .page-btn {
    padding: 0.5rem 1rem;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .page-btn:hover:not(:disabled) {
    background: var(--color-border);
  }

  .page-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .page-info {
    color: var(--color-text);
    font-weight: 500;
  }

  .page-meta {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-wrap: wrap;
  }

  .page-size {
    color: var(--color-text);
    opacity: 0.7;
    font-size: 0.9rem;
  }

  .invoice {
    text-align: center;
  }

  .btn-create-invoice {
    padding: 0.5rem 1rem;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    white-space: nowrap;
  }

  .btn-create-invoice:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  }

  .invoice-link {
    background: none;
    border: none;
    color: var(--color-accent);
    font-weight: 600;
    font-size: 0.875rem;
    cursor: pointer;
    text-decoration: underline;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    transition: all 0.2s ease;
  }

  .invoice-link:hover {
    background: var(--color-panel);
    text-decoration: none;
  }

  /* Filter Styles */
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 1.5rem;
  }

  .filter-toggle-btn {
    padding: 0.5rem 1rem;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 6px;
    cursor: pointer;
    font-size: 0.875rem;
    transition: all 0.2s ease;
  }

  .filter-toggle-btn:hover {
    background: var(--color-border);
  }

  .filter-section {
    background: var(--color-card);
    border: 1px solid var(--color-border);
    border-radius: 10px;
    padding: 1rem;
    margin-bottom: 1.5rem;
    animation: slideDown 0.3s ease;
  }

  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .filter-two-column {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
    align-items: stretch; /* Make both columns same height */
  }

  /* Match left column height for right column */
  .filter-left-column {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .filter-right-column {
    display: flex;
    flex-direction: column;
    align-self: stretch;
    position: relative;
    padding-left: 1.5rem;
    height: 100%; /* Match left column height */
    min-height: 0; /* Allow height constraint */
  }

  .filter-right-column::before {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 1px;
    background: var(--color-border);
    opacity: 0.5;
  }

  @media (max-width: 768px) {
    .filter-two-column {
      grid-template-columns: 1fr;
    }
  }


  .filter-row {
    display: flex;
    gap: 1rem;
    align-items: flex-end;
    position: relative;
    padding-bottom: 1rem;
  }

  .filter-row:not(:last-child)::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: var(--color-border);
    opacity: 0.5;
  }

  .filter-group-half {
    flex: 1 1 50%;
  }

  .filter-group-full {
    flex: 1 1 100%;
  }

  .filter-procedure-section {
    display: flex;
    flex-direction: column;
    height: 100%;
    min-height: 0; /* Allow height constraint from parent */
    overflow: hidden; /* Prevent section from exceeding parent */
  }

  .procedures-grid-filter {
    display: grid;
    grid-template-columns: 1fr 1fr; /* Two-column layout like New Session modal */
    gap: 0.75rem; /* Match New Session modal gap */
    flex: 1 1 0; /* Take remaining space, start from 0 */
    min-height: 0; /* Critical: Allow flexbox to constrain height */
    max-height: 100%; /* Don't exceed parent */
    overflow-y: auto; /* Enable scrolling when content exceeds height */
    overflow-x: hidden;
    align-content: start; /* Align items to top */
  }

  .filter-procedure-section .section-label {
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--color-text);
    margin-bottom: 0.5rem;
    flex-shrink: 0; /* Don't shrink label */
  }

  .procedure-search-container {
    margin-bottom: 0.75rem;
    flex-shrink: 0; /* Don't shrink search container */
  }


  /* Filter procedure tags - match New Session modal exactly */
  .procedures-grid-filter button.procedure-tag {
    /* Complete button reset */
    all: unset;
    box-sizing: border-box;
    display: inline-block;
    
    /* Apply exact styling from New Session modal */
    padding: 0.625rem 1rem;
    background: var(--color-card);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 6px;
    font-size: 0.875rem;
    font-weight: 500;
    font-family: inherit;
    cursor: pointer;
    transition: all 0.2s ease;
    white-space: nowrap;
    width: 100%; /* Fill grid cell */
    text-align: left; /* Left-align text in tags */
    outline: none;
    margin: 0;
  }

  .procedures-grid-filter .procedure-tag:hover:not(:disabled) {
    background: var(--color-panel);
    border-color: var(--color-accent);
    transform: translateY(-1px);
  }

  .procedures-grid-filter .procedure-tag.selected {
    background: var(--color-accent);
    color: white;
    border-color: var(--color-accent);
  }

  .procedures-grid-filter .empty-procedures {
    grid-column: 1 / -1; /* Span both columns */
    padding: 2rem;
    text-align: center;
    color: var(--color-text);
    opacity: 0.6;
    font-size: 0.875rem;
  }

  .filter-actions-bottom {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid var(--color-border);
  }

  .filter-group {
    display: flex;
    flex-direction: column;
    gap: 0.375rem;
  }

  .filter-group label {
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--color-text);
    min-width: 0; /* Ensure labels don't affect alignment */
    width: 100%;
  }

  /* Ensure all inputs align perfectly */
  .filter-group .form-input,
  .filter-group .search-container,
  .filter-group .status-radio-group {
    width: 100%;
  }


  .search-container {
    position: relative;
  }

  .search-container input:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }

  .status-radio-group {
    display: flex;
    gap: 0.75rem;
    align-items: center;
  }

  .radio-option {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    cursor: pointer;
    font-size: 0.875rem;
    color: var(--color-text);
    padding: 0.375rem 0.75rem;
    border-radius: 6px;
    transition: all 0.2s ease;
  }

  .radio-option input[type="radio"] {
    margin: 0;
    cursor: pointer;
  }

  .radio-option span {
    user-select: none;
  }

  .radio-option:has(input[type="radio"]:checked) {
    background: var(--color-panel);
    font-weight: 500;
  }

  .radio-completed:has(input[type="radio"]:checked) {
    background: rgba(76, 175, 80, 0.15);
    color: #2e7d32;
  }

  .radio-in-progress:has(input[type="radio"]:checked) {
    background: rgba(255, 152, 0, 0.15);
    color: #f57c00;
  }

  .clear-selection-btn {
    position: absolute;
    right: 0.5rem;
    top: 50%;
    transform: translateY(-50%);
    background: none;
    border: none;
    color: var(--color-text);
    cursor: pointer;
    font-size: 1.25rem;
    line-height: 1;
    padding: 0;
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    transition: all 0.2s ease;
    opacity: 0.6;
  }

  .clear-selection-btn:hover {
    opacity: 1;
    background: var(--color-border);
  }


  .dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: var(--color-card);
    border: 1px solid var(--color-border);
    border-radius: 6px;
    margin-top: 0.25rem;
    max-height: 200px;
    overflow-y: auto;
    z-index: 100;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  .dropdown-item {
    padding: 0.75rem 1rem;
    cursor: pointer;
    transition: background 0.2s ease;
    border-bottom: 1px solid var(--color-border);
  }

  .dropdown-item:last-child {
    border-bottom: none;
  }

  .dropdown-item:hover {
    background: var(--color-panel);
  }

  /* Active Filters Bar */
  .active-filters-bar {
    background: var(--color-card);
    border: 1px solid var(--color-border);
    border-radius: 10px;
    padding: 0.75rem;
    margin-bottom: 1.5rem;
  }

  .active-filters-main,
  .active-filters-procedures {
    display: flex;
    flex-wrap: wrap;
    gap: 0.375rem;
    margin-bottom: 0.5rem;
  }

  .active-filters-procedures {
    padding-top: 0.5rem;
    border-top: 1px solid var(--color-border);
  }

  .filter-tag {
    display: inline-flex;
    align-items: center;
    gap: 0.375rem;
    padding: 0.25rem 0.625rem;
    background: var(--color-panel);
    border: 1px solid var(--color-border);
    border-radius: 6px;
    font-size: 0.8125rem;
    color: var(--color-text);
  }

  .filter-tag-remove {
    background: none;
    border: none;
    color: var(--color-text);
    cursor: pointer;
    font-size: 1.125rem;
    line-height: 1;
    padding: 0;
    width: 18px;
    height: 18px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    transition: all 0.2s ease;
  }

  .filter-tag-remove:hover {
    background: var(--color-border);
    color: var(--color-danger);
  }

  .active-filters-actions {
    display: flex;
    gap: 0.75rem;
    padding-top: 0.5rem;
    border-top: 1px solid var(--color-border);
  }

  .btn-filter-action {
    padding: 0.375rem 0.75rem;
    background: var(--color-panel);
    border: 1px solid var(--color-border);
    border-radius: 6px;
    color: var(--color-accent);
    cursor: pointer;
    font-size: 0.875rem;
    font-weight: 500;
    transition: all 0.2s ease;
  }

  .btn-filter-action:hover {
    background: var(--color-accent);
    color: #fff;
    border-color: var(--color-accent);
  }
</style>
