<script>
  import { onMount } from 'svelte';
  import { account, licenseValid, setLicense, licenseValidationStatus, validateCurrentLicense } from '../stores/settingsStore.js';
  import { deleteAllPatients } from '../stores/patientStore.js';
  import { isAdmin } from '../stores/authStore.js';
  import UserManagement from './UserManagement.svelte';
  import {
    filteredProcedures,
    procedures,
    proceduresLoading,
    proceduresError,
    proceduresSuccess,
    procedureSearch,
    proceduresCurrentPage,
    proceduresTotalPages,
    loadProcedures,
    loadProceduresPaginated,
    createProcedure,
    updateProcedure,
    deleteProcedure
  } from '../stores/procedureStore.js';
  import {
    filteredWorkTypes,
    workTypesLoading,
    workTypesError,
    workTypesSuccess,
    workTypeSearch,
    workTypesCurrentPage,
    workTypesTotalPages,
    loadWorkTypesPaginated,
    createWorkType,
    updateWorkType,
    deleteWorkType
  } from '../stores/workTypeStore.js';
  import {
    filteredColorShades,
    colorShadesLoading,
    colorShadesError,
    colorShadesSuccess,
    colorShadeSearch,
    colorShadesCurrentPage,
    colorShadesTotalPages,
    loadColorShadesPaginated,
    createColorShade,
    updateColorShade,
    deleteColorShade
  } from '../stores/colorShadeStore.js';

  let selectedSection = 'license'; // 'license', 'users', 'procedures', 'work-types', 'color-shades', 'danger'
  let showLicenseInput = false;
  let newKey = '';
  let validatingLicense = false;
  let showDeleteAllConfirm = false;
  let deletingAllPatients = false;
  let showUserManagement = false;

  // Procedure modal state
  let showProcedureModal = false;
  let procedureModalMode = 'create';
  let procedureForm = { id: null, name: '', priceDisplay: '' };
  let procedureFormError = '';
  let procedureFormLoading = false;
  let showProcedureDeleteConfirm = false;
  let procedureToDelete = null;

  // Work Type modal state
  let showWorkTypeModal = false;
  let workTypeModalMode = 'create';
  let workTypeForm = { id: null, name: '', description: '' };
  let workTypeFormError = '';
  let workTypeFormLoading = false;
  let showWorkTypeDeleteConfirm = false;
  let workTypeToDelete = null;

  // Color Shade modal state
  let showColorShadeModal = false;
  let colorShadeModalMode = 'create';
  let colorShadeForm = { id: null, name: '', description: '', hex_color: '#F1ECE4', is_active: true };
  let colorShadeFormError = '';
  let colorShadeFormLoading = false;
  let showColorShadeDeleteConfirm = false;
  let colorShadeToDelete = null;

  onMount(() => {
    selectedSection = 'license';
    // Don't load procedures on mount - load when procedures section is selected
  });

  // Reactive statement to log when procedures section is rendered
  // $: if (selectedSection === 'procedures') {
  //   console.log('[Configuration] Procedures section is being rendered');
  //   console.log('[Configuration] proceduresLoading:', $proceduresLoading);
  //   console.log('[Configuration] proceduresError:', $proceduresError);
  //   console.log('[Configuration] filteredProcedures count:', $filteredProcedures?.length || 0);
  //   console.log('[Configuration] procedures store value:', $filteredProcedures);
  // }

  function selectSection(section) {
    // console.log('[Configuration] selectSection called with section:', section);
    // console.log('[Configuration] Current selectedSection before change:', selectedSection);
    // console.log('[Configuration] isAdmin() check:', isAdmin());
    
    selectedSection = section;
    // console.log('[Configuration] selectedSection updated to:', selectedSection);
    
    if (section !== 'license') {
      showLicenseInput = false;
    }
    if (section !== 'users') {
      showUserManagement = false;
    }
    if (section !== 'procedures') {
      showProcedureModal = false;
      showProcedureDeleteConfirm = false;
      procedureFormError = '';
    }
    if (section !== 'work-types') {
      showWorkTypeModal = false;
      showWorkTypeDeleteConfirm = false;
      workTypeFormError = '';
    }
    if (section !== 'color-shades') {
      showColorShadeModal = false;
      showColorShadeDeleteConfirm = false;
      colorShadeFormError = '';
    }
    
    if (section === 'procedures') {
      // Load paginated procedures when section is selected
      loadProceduresPaginated(1).then(() => {
        // console.log('[Configuration] loadProceduresPaginated() completed');
      }).catch((err) => {
        console.error('[Configuration] loadProceduresPaginated() failed:', err);
      });
    } else if (section === 'work-types') {
      loadWorkTypesPaginated(1).catch((err) => {
        console.error('[Configuration] loadWorkTypesPaginated() failed:', err);
      });
    } else if (section === 'color-shades') {
      loadColorShadesPaginated(1).catch((err) => {
        console.error('[Configuration] loadColorShadesPaginated() failed:', err);
      });
    }
    
    // console.log('[Configuration] selectSection completed');
  }

  async function handleLicenseSave() {
    if (!newKey.trim()) {
      return;
    }
    
    validatingLicense = true;
    try {
      const isValid = await setLicense(newKey);
      if (isValid) {
        showLicenseInput = false;
        newKey = '';
      }
    } catch (error) {
      console.error('License save error:', error);
    } finally {
      validatingLicense = false;
    }
  }

  async function refreshLicense() {
    await validateCurrentLicense();
  }

  async function handleDeleteAllPatients() {
    deletingAllPatients = true;
    try {
      await deleteAllPatients();
      showDeleteAllConfirm = false;
    } catch (error) {
      console.error('Error deleting all patients:', error);
      alert('Failed to delete all patients. Please try again.');
    } finally {
      deletingAllPatients = false;
    }
  }

  function formatPrice(price) {
    return price.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
  }

  function openCreateProcedureModal() {
    procedureModalMode = 'create';
    procedureForm = { id: null, name: '', priceDisplay: '' };
    procedureFormError = '';
    showProcedureModal = true;
  }

  function openEditProcedureModal(proc) {
    procedureModalMode = 'edit';
    procedureForm = { id: proc.id, name: proc.name, priceDisplay: formatPrice(proc.price) };
    procedureFormError = '';
    showProcedureModal = true;
  }

  async function handleProcedureSave() {
    console.log('[Configuration] handleProcedureSave called');
    console.log('[Configuration] procedureForm:', procedureForm);
    procedureFormError = '';
    
    // Validate procedure name
    const trimmedName = procedureForm.name.trim();
    if (!trimmedName) {
      procedureFormError = 'Procedure name is required';
      console.log('[Configuration] Validation failed: Name is required');
      return;
    }
    
    // Check for duplicate procedure name (case-insensitive)
    const existingProcedures = $procedures || [];
    const duplicateProcedure = existingProcedures.find(proc => 
      proc.name.toLowerCase() === trimmedName.toLowerCase() && 
      proc.id !== procedureForm.id // Exclude current procedure when editing
    );
    if (duplicateProcedure) {
      procedureFormError = 'A procedure with this name already exists';
      console.log('[Configuration] Validation failed: Duplicate procedure name');
      return;
    }
    
    // Validate price - must be a positive integer
    const priceString = String(procedureForm.priceDisplay || '').replace(/,/g, '');
    const priceValue = parseInt(priceString, 10);
    console.log('[Configuration] Parsed price value:', priceValue);
    
    // Check if price is a valid number
    if (isNaN(priceValue) || priceString === '') {
      procedureFormError = 'Price must be a valid number';
      console.log('[Configuration] Validation failed: Price is not a valid number');
      return;
    }
    
    // Check if price is a positive integer
    if (priceValue <= 0) {
      procedureFormError = 'Price must be a positive integer greater than zero';
      console.log('[Configuration] Validation failed: Price must be greater than zero');
      return;
    }
    
    // Check if price is actually an integer (not a decimal)
    if (parseFloat(priceString) !== priceValue) {
      procedureFormError = 'Price must be a whole number (no decimals)';
      console.log('[Configuration] Validation failed: Price must be an integer');
      return;
    }
    procedureFormLoading = true;
    const payloadPrice = priceValue;
    console.log('[Configuration] Payload to send:', {
      name: procedureForm.name.trim(),
      price: payloadPrice,
      mode: procedureModalMode
    });
    let success = false;
    try {
      if (procedureModalMode === 'create') {
        console.log('[Configuration] Creating procedure...');
        success = await createProcedure({
          name: procedureForm.name.trim(),
          price: payloadPrice
        });
        console.log('[Configuration] createProcedure returned:', success);
      } else {
        console.log('[Configuration] Updating procedure...');
        success = await updateProcedure({
          id: procedureForm.id,
          name: procedureForm.name.trim(),
          price: payloadPrice
        });
        console.log('[Configuration] updateProcedure returned:', success);
      }
    } catch (err) {
      console.error('[Configuration] Error in handleProcedureSave:', err);
      procedureFormError = err.message || 'Failed to save procedure';
      success = false;
    }
    procedureFormLoading = false;
    if (success) {
      console.log('[Configuration] Procedure saved successfully, closing modal');
      showProcedureModal = false;
      procedureForm = { id: null, name: '', priceDisplay: '' };
      procedureFormError = '';
    } else {
      // Check if there's an error from the store
      if ($proceduresError) {
        procedureFormError = $proceduresError;
      }
      console.log('[Configuration] Procedure save failed, error:', procedureFormError || $proceduresError);
    }
  }

  function confirmDeleteProcedure(proc) {
    procedureToDelete = proc;
    showProcedureDeleteConfirm = true;
  }

  async function handleDeleteProcedure() {
    if (!procedureToDelete) return;
    const success = await deleteProcedure(procedureToDelete.id);
    if (success) {
      showProcedureDeleteConfirm = false;
      procedureToDelete = null;
    }
  }

  // Pagination functions
  function previousProcedures() {
    const current = $proceduresCurrentPage;
    if (current > 1) {
      loadProceduresPaginated(current - 1);
    }
  }

  function nextProcedures() {
    const current = $proceduresCurrentPage;
    const total = $proceduresTotalPages;
    if (current < total) {
      loadProceduresPaginated(current + 1);
    }
  }

  // Work Type handlers
  function openCreateWorkTypeModal() {
    workTypeModalMode = 'create';
    workTypeForm = { id: null, name: '', description: '' };
    workTypeFormError = '';
    showWorkTypeModal = true;
  }

  function openEditWorkTypeModal(workType) {
    workTypeModalMode = 'edit';
    workTypeForm = { id: workType.id, name: workType.name, description: workType.description || '' };
    workTypeFormError = '';
    showWorkTypeModal = true;
  }

  async function handleWorkTypeSave() {
    workTypeFormError = '';
    
    const trimmedName = workTypeForm.name.trim();
    if (!trimmedName) {
      workTypeFormError = 'Work type name is required';
      return;
    }
    
    workTypeFormLoading = true;
    try {
      if (workTypeModalMode === 'create') {
        await createWorkType({
          name: trimmedName,
          description: workTypeForm.description.trim() || ''
        });
      } else {
        await updateWorkType(workTypeForm.id, {
          name: trimmedName,
          description: workTypeForm.description.trim() || ''
        });
      }
      showWorkTypeModal = false;
      workTypeForm = { id: null, name: '', description: '' };
      workTypeFormError = '';
    } catch (err) {
      workTypeFormError = err.message || 'Failed to save work type';
    } finally {
      workTypeFormLoading = false;
    }
  }

  function confirmDeleteWorkType(workType) {
    workTypeToDelete = workType;
    showWorkTypeDeleteConfirm = true;
  }

  async function handleDeleteWorkType() {
    if (!workTypeToDelete) return;
    const success = await deleteWorkType(workTypeToDelete.id);
    if (success) {
      showWorkTypeDeleteConfirm = false;
      workTypeToDelete = null;
    }
  }

  function previousWorkTypes() {
    const current = $workTypesCurrentPage;
    if (current > 1) {
      loadWorkTypesPaginated(current - 1);
    }
  }

  function nextWorkTypes() {
    const current = $workTypesCurrentPage;
    const total = $workTypesTotalPages;
    if (current < total) {
      loadWorkTypesPaginated(current + 1);
    }
  }

  // Color Shade handlers
  function openCreateColorShadeModal() {
    colorShadeModalMode = 'create';
    colorShadeForm = { id: null, name: '', description: '', hex_color: '#F1ECE4', is_active: true };
    colorShadeFormError = '';
    showColorShadeModal = true;
  }

  function openEditColorShadeModal(shade) {
    colorShadeModalMode = 'edit';
    colorShadeForm = { 
      id: shade.id, 
      name: shade.name, 
      description: shade.description || '', 
      hex_color: shade.hex_color || '#F1ECE4',
      is_active: shade.is_active !== undefined ? shade.is_active : true
    };
    colorShadeFormError = '';
    showColorShadeModal = true;
  }

  async function handleColorShadeSave() {
    colorShadeFormError = '';
    
    const trimmedName = colorShadeForm.name.trim();
    if (!trimmedName) {
      colorShadeFormError = 'Color shade name is required';
      return;
    }
    
    colorShadeFormLoading = true;
    try {
      if (colorShadeModalMode === 'create') {
        await createColorShade({
          name: trimmedName,
          description: colorShadeForm.description.trim() || '',
          hex_color: colorShadeForm.hex_color || '#F1ECE4',
          is_active: colorShadeForm.is_active
        });
      } else {
        await updateColorShade(colorShadeForm.id, {
          name: trimmedName,
          description: colorShadeForm.description.trim() || '',
          hex_color: colorShadeForm.hex_color || '#F1ECE4',
          is_active: colorShadeForm.is_active
        });
      }
      showColorShadeModal = false;
      colorShadeForm = { id: null, name: '', description: '', hex_color: '#F1ECE4', is_active: true };
      colorShadeFormError = '';
    } catch (err) {
      colorShadeFormError = err.message || 'Failed to save color shade';
    } finally {
      colorShadeFormLoading = false;
    }
  }

  function confirmDeleteColorShade(shade) {
    colorShadeToDelete = shade;
    showColorShadeDeleteConfirm = true;
  }

  async function handleDeleteColorShade() {
    if (!colorShadeToDelete) return;
    const success = await deleteColorShade(colorShadeToDelete.id);
    if (success) {
      showColorShadeDeleteConfirm = false;
      colorShadeToDelete = null;
    }
  }

  function previousColorShades() {
    const current = $colorShadesCurrentPage;
    if (current > 1) {
      loadColorShadesPaginated(current - 1);
    }
  }

  function nextColorShades() {
    const current = $colorShadesCurrentPage;
    const total = $colorShadesTotalPages;
    if (current < total) {
      loadColorShadesPaginated(current + 1);
    }
  }
</script>

<div class="configuration">
  <div class="config-container">
    <!-- Left Sidebar Navigation -->
    <aside class="sidebar">
      <h2 class="sidebar-title">Configuration</h2>
      <nav class="sidebar-nav">
        <button 
          class="nav-item" 
          class:active={selectedSection === 'license'}
          on:click={() => selectSection('license')}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
            <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
          </svg>
          <span>License & Account</span>
        </button>
        
        {#if isAdmin()}
        <button 
          class="nav-item" 
          class:active={selectedSection === 'users'}
          on:click={() => selectSection('users')}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
          <span>User Management</span>
        </button>
        {/if}
        
        <button 
          class="nav-item" 
          class:active={selectedSection === 'procedures'}
          on:click={() => {
            // console.log('[Configuration] Dental Procedures button clicked');
            selectSection('procedures');
          }}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
          <span>Dental Procedures</span>
        </button>
        
        <button 
          class="nav-item" 
          class:active={selectedSection === 'work-types'}
          on:click={() => selectSection('work-types')}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
          </svg>
          <span>Work Types</span>
        </button>
        
        <button 
          class="nav-item" 
          class:active={selectedSection === 'color-shades'}
          on:click={() => selectSection('color-shades')}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="13.5" cy="6.5" r=".5" fill="currentColor"/>
            <circle cx="17.5" cy="10.5" r=".5" fill="currentColor"/>
            <circle cx="8.5" cy="7.5" r=".5" fill="currentColor"/>
            <circle cx="6.5" cy="12.5" r=".5" fill="currentColor"/>
            <path d="M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10c.926 0 1.648-.746 1.648-1.688 0-.437-.18-.835-.437-1.125-.29-.289-.438-.652-.438-1.125a1.64 1.64 0 0 1 1.668-1.668h1.996c3.051 0 5.555-2.503 5.555-5.554C21.965 6.012 17.461 2 12 2z"/>
          </svg>
          <span>Color Shades</span>
        </button>
        
        <button 
          class="nav-item danger" 
          class:active={selectedSection === 'danger'}
          on:click={() => selectSection('danger')}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
            <line x1="12" y1="9" x2="12" y2="13"/>
            <line x1="12" y1="17" x2="12.01" y2="17"/>
          </svg>
          <span>Reset Application Data</span>
        </button>
      </nav>
    </aside>

    <!-- Right Content Panel -->
    <main class="content-panel">
      <!-- License & Account Section -->
      {#if selectedSection === 'license'}
        <div class="section-content">
          <div class="section-header">
            <h1>License & Account</h1>
            <p class="section-description">Manage your license key and account information</p>
          </div>

          <div class="content-cards">
            <!-- Account Information Card -->
            <div class="info-card">
              <h3>Account Information</h3>
              <div class="info-list">
                <div class="info-item">
                  <span class="info-label">Email:</span>
                  <span class="info-value">{$account.email}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">License Key:</span>
                  <span class="info-value">{$account.licenseKey ? `${$account.licenseKey.substring(0, 20)}...` : 'Not set'}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">License Expiry:</span>
                  <span class="info-value">{$account.licenseExpiry || 'Unknown'}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">Status:</span>
                  <span class="info-value">
                    {#if $licenseValidationStatus.isChecking}
                      <span class="status-badge checking">Checking...</span>
                    {:else if $licenseValid}
                      <span class="status-badge valid">✓ Valid</span>
                    {:else}
                      <span class="status-badge expired">✗ Invalid/Expired</span>
                    {/if}
                  </span>
                </div>
                {#if $licenseValidationStatus.message}
                  <div class="info-item message-item">
                    <span class="info-label">Message:</span>
                    <span class="info-value">{$licenseValidationStatus.message}</span>
                  </div>
                {/if}
              </div>
            </div>

            <!-- License Actions Card -->
            <div class="action-card">
              <h3>License Actions</h3>
              <div class="action-buttons">
                <button 
                  class="btn btn-secondary" 
                  on:click={refreshLicense} 
                  disabled={$licenseValidationStatus.isChecking}
                >
                  {#if $licenseValidationStatus.isChecking}
                    🔄 Checking...
                  {:else}
                    🔍 Check License
                  {/if}
                </button>
                <button 
                  class="btn btn-primary" 
                  on:click={() => showLicenseInput = !showLicenseInput}
                >
                  {#if showLicenseInput}❌ Cancel{:else}✏️ Update License{/if}
                </button>
              </div>
              
              {#if showLicenseInput}
                <div class="license-form">
                  <input 
                    type="text" 
                    placeholder="Enter new license key..." 
                    bind:value={newKey}
                    disabled={validatingLicense}
                    class="form-input"
                  />
                  <button 
                    class="btn btn-primary btn-small" 
                    on:click={handleLicenseSave}
                    disabled={validatingLicense || !newKey.trim()}
                  >
                    {#if validatingLicense}🔄 Validating...{:else}💾 Save & Validate{/if}
                  </button>
                </div>
              {/if}
            </div>
          </div>
        </div>
      {/if}

      <!-- User Management Section -->
      {#if selectedSection === 'users'}
        <div class="section-content">
          <div class="section-header">
            <h1>User Management</h1>
            <p class="section-description">Create and manage user accounts</p>
          </div>

          <div class="content-cards">
            <div class="action-card">
              <h3>Manage Users</h3>
              <p class="card-description">View, create, and manage user accounts in the system.</p>
              <button class="btn btn-primary" on:click={() => showUserManagement = true}>
                👥 Open User Management
              </button>
            </div>
          </div>
        </div>
      {/if}

      <!-- Dental Procedures Section -->
      {#if selectedSection === 'procedures'}
        <div class="section-content">
          <div class="section-header">
            <h1>Dental Procedures</h1>
            <p class="section-description">Manage your clinic's available procedures and pricing.</p>
          </div>

          <div class="content-cards">
            <div class="action-card">
              <div class="procedures-toolbar">
                <div class="search-group">
                  <input
                    type="text"
                    class="form-input search-input"
                    placeholder="Search procedures..."
                    bind:value={$procedureSearch}
                  />
                </div>
                <button class="btn btn-primary" on:click={openCreateProcedureModal}>
                  ➕ Create New Procedure
                </button>
              </div>

              {#if $proceduresError}
                <div class="alert alert-error">
                  {$proceduresError}
                </div>
              {/if}
              {#if $proceduresSuccess}
                <div class="alert alert-success">
                  {$proceduresSuccess}
                </div>
              {/if}

              <div class="procedure-list">
                {#if $proceduresLoading}
                  <div class="loading-state">
                    <div class="spinner"></div>
                    <p>Loading procedures...</p>
                  </div>
                {:else if $filteredProcedures.length === 0}
                  <div class="empty-state">
                    <p>No procedures found.</p>
                  </div>
                {:else}
                  <table>
                    <thead>
                      <tr>
                        <th>Procedure Name</th>
                        <th>Price</th>
                        <th class="actions-col">Actions</th>
                      </tr>
                    </thead>
                    <tbody>
                      {#each $filteredProcedures as procedure}
                        <tr>
                          <td>{procedure.name}</td>
                          <td>{formatPrice(procedure.price)} SYP</td>
                          <td class="actions-col">
                            <button class="icon-btn" on:click={() => openEditProcedureModal(procedure)} title="Edit">
                              ✏️
                            </button>
                            <button class="icon-btn danger" on:click={() => confirmDeleteProcedure(procedure)} title="Delete">
                              🗑️
                            </button>
                          </td>
                        </tr>
                      {/each}
                    </tbody>
                  </table>
                  
                  {#if $proceduresTotalPages > 1}
                    <div class="pagination">
                      <button 
                        class="page-btn" 
                        disabled={$proceduresCurrentPage === 1 || $proceduresLoading} 
                        on:click={previousProcedures}
                      >
                        Previous
                      </button>
                      <span class="page-info">
                        Page {$proceduresCurrentPage} of {$proceduresTotalPages}
                      </span>
                      <button 
                        class="page-btn" 
                        disabled={$proceduresCurrentPage >= $proceduresTotalPages || $proceduresLoading} 
                        on:click={nextProcedures}
                      >
                        Next
                      </button>
                    </div>
                  {/if}
                {/if}
              </div>
            </div>
          </div>
        </div>
      {/if}

      <!-- Work Types Section -->
      {#if selectedSection === 'work-types'}
        <div class="section-content">
          <div class="section-header">
            <h1>Work Types</h1>
            <p class="section-description">Define your lab work types for orders.</p>
          </div>

          <div class="content-cards">
            <div class="action-card">
              <div class="procedures-toolbar">
                <div class="search-group">
                  <input
                    type="text"
                    class="form-input search-input"
                    placeholder="Search work types..."
                    bind:value={$workTypeSearch}
                  />
                </div>
                <button class="btn btn-primary" on:click={openCreateWorkTypeModal}>
                  ➕ Create New Work Type
                </button>
              </div>

              {#if $workTypesError}
                <div class="alert alert-error">
                  {$workTypesError}
                </div>
              {/if}
              {#if $workTypesSuccess}
                <div class="alert alert-success">
                  {$workTypesSuccess}
                </div>
              {/if}

              <div class="procedure-list">
                {#if $workTypesLoading}
                  <div class="loading-state">
                    <div class="spinner"></div>
                    <p>Loading work types...</p>
                  </div>
                {:else if $filteredWorkTypes.length === 0}
                  <div class="empty-state">
                    <p>No work types defined yet. Create your first work type.</p>
                  </div>
                {:else}
                  <table>
                    <thead>
                      <tr>
                        <th>Name</th>
                        <th>Description</th>
                        <th class="actions-col">Actions</th>
                      </tr>
                    </thead>
                    <tbody>
                      {#each $filteredWorkTypes as workType}
                        <tr>
                          <td>{workType.name}</td>
                          <td>{workType.description || '-'}</td>
                          <td class="actions-col">
                            <button class="icon-btn" on:click={() => openEditWorkTypeModal(workType)} title="Edit">
                              ✏️
                            </button>
                            <button class="icon-btn danger" on:click={() => confirmDeleteWorkType(workType)} title="Delete">
                              🗑️
                            </button>
                          </td>
                        </tr>
                      {/each}
                    </tbody>
                  </table>
                  
                  {#if $workTypesTotalPages > 1}
                    <div class="pagination">
                      <button 
                        class="page-btn" 
                        disabled={$workTypesCurrentPage === 1 || $workTypesLoading} 
                        on:click={previousWorkTypes}
                      >
                        Previous
                      </button>
                      <span class="page-info">
                        Page {$workTypesCurrentPage} of {$workTypesTotalPages}
                      </span>
                      <button 
                        class="page-btn" 
                        disabled={$workTypesCurrentPage >= $workTypesTotalPages || $workTypesLoading} 
                        on:click={nextWorkTypes}
                      >
                        Next
                      </button>
                    </div>
                  {/if}
                {/if}
              </div>
            </div>
          </div>
        </div>
      {/if}

      <!-- Color Shades Section -->
      {#if selectedSection === 'color-shades'}
        <div class="section-content">
          <div class="section-header">
            <h1>Color Shades</h1>
            <p class="section-description">Define tooth color shades for lab orders.</p>
          </div>

          <div class="content-cards">
            <div class="action-card">
              <div class="procedures-toolbar">
                <div class="search-group">
                  <input
                    type="text"
                    class="form-input search-input"
                    placeholder="Search color shades..."
                    bind:value={$colorShadeSearch}
                  />
                </div>
                <button class="btn btn-primary" on:click={openCreateColorShadeModal}>
                  ➕ Create New Color Shade
                </button>
              </div>

              {#if $colorShadesError}
                <div class="alert alert-error">
                  {$colorShadesError}
                </div>
              {/if}
              {#if $colorShadesSuccess}
                <div class="alert alert-success">
                  {$colorShadesSuccess}
                </div>
              {/if}

              <div class="procedure-list">
                {#if $colorShadesLoading}
                  <div class="loading-state">
                    <div class="spinner"></div>
                    <p>Loading color shades...</p>
                  </div>
                {:else if $filteredColorShades.length === 0}
                  <div class="empty-state">
                    <p>No color shades defined yet. Create your first color shade.</p>
                  </div>
                {:else}
                  <table>
                    <thead>
                      <tr>
                        <th>Name</th>
                        <th>Description</th>
                        <th>Color Sample</th>
                        <th class="actions-col">Actions</th>
                      </tr>
                    </thead>
                    <tbody>
                      {#each $filteredColorShades as shade}
                        <tr>
                          <td>{shade.name}</td>
                          <td>{shade.description || '-'}</td>
                          <td>
                            <div 
                              class="color-sample" 
                              style="background-color: {shade.hex_color || '#F1ECE4'};"
                              title="{shade.hex_color || '#F1ECE4'}"
                            ></div>
                          </td>
                          <td class="actions-col">
                            <button class="icon-btn" on:click={() => openEditColorShadeModal(shade)} title="Edit">
                              ✏️
                            </button>
                            <button class="icon-btn danger" on:click={() => confirmDeleteColorShade(shade)} title="Delete">
                              🗑️
                            </button>
                          </td>
                        </tr>
                      {/each}
                    </tbody>
                  </table>
                  
                  {#if $colorShadesTotalPages > 1}
                    <div class="pagination">
                      <button 
                        class="page-btn" 
                        disabled={$colorShadesCurrentPage === 1 || $colorShadesLoading} 
                        on:click={previousColorShades}
                      >
                        Previous
                      </button>
                      <span class="page-info">
                        Page {$colorShadesCurrentPage} of {$colorShadesTotalPages}
                      </span>
                      <button 
                        class="page-btn" 
                        disabled={$colorShadesCurrentPage >= $colorShadesTotalPages || $colorShadesLoading} 
                        on:click={nextColorShades}
                      >
                        Next
                      </button>
                    </div>
                  {/if}
                {/if}
              </div>
            </div>
          </div>
        </div>
      {/if}

      <!-- Danger Zone Section -->
      {#if selectedSection === 'danger'}
        <div class="section-content">
          <div class="section-header">
            <h1>Reset Application Data</h1>
            <p class="section-description danger-text">These actions cannot be undone. Please be careful!</p>
          </div>

          <div class="content-cards">
            <div class="danger-card">
              <h3>Delete All Data</h3>
              <p class="card-description">Permanently delete all patient records, appointments, payments, and folder data. This action cannot be undone.</p>
              <button class="btn btn-danger" on:click={() => showDeleteAllConfirm = true}>
                🗑️ Delete All Data
              </button>
            </div>
          </div>
        </div>
      {/if}
    </main>
  </div>
</div>

<!-- User Management Modal -->
{#if showUserManagement}
  <div 
    class="modal-overlay" 
    role="button"
    tabindex="0"
    on:click={() => showUserManagement = false}
    on:keydown={(e) => (e.key === 'Enter' || e.key === ' ') && (showUserManagement = false)}
  >
    <div class="modal-content" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header">
        <h3>User Management</h3>
        <button class="close-btn" on:click={() => showUserManagement = false}>×</button>
      </div>
      <div class="modal-body">
        <UserManagement />
      </div>
    </div>
  </div>
{/if}

<!-- Delete All Patients Confirmation Modal -->
{#if showDeleteAllConfirm}
  <div 
    class="modal-overlay" 
    role="button"
    tabindex="0"
    on:click={() => showDeleteAllConfirm = false}
    on:keydown={(e) => (e.key === 'Enter' || e.key === ' ') && (showDeleteAllConfirm = false)}
  >
    <div class="confirmation-modal" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header danger-header">
        <h3>⚠️ Delete All Patients</h3>
      </div>
      <div class="modal-content">
        <p><strong>This action will permanently delete:</strong></p>
        <ul class="deletion-list">
          <li>All patient records</li>
          <li>All patient appointments</li>
          <li>All patient payment records</li>
          <li>All patient folder data</li>
        </ul>
        <p class="final-warning"><strong>This action cannot be undone!</strong></p>
      </div>
      <div class="modal-actions">
        <button 
          class="btn btn-danger" 
          on:click={handleDeleteAllPatients}
          disabled={deletingAllPatients}
        >
          {#if deletingAllPatients}🔄 Deleting...{:else}✅ Yes, Delete All{/if}
        </button>
        <button 
          class="btn btn-secondary" 
          on:click={() => showDeleteAllConfirm = false}
          disabled={deletingAllPatients}
        >
          ❌ Cancel
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Procedure Modal -->
{#if showProcedureModal}
  <div 
    class="modal-overlay" 
    role="button"
    tabindex="0"
    on:click={() => showProcedureModal = false}
    on:keydown={(e) => (e.key === 'Enter' || e.key === ' ') && (showProcedureModal = false)}
  >
    <div class="modal-content small" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header">
        <h3>{procedureModalMode === 'create' ? 'Create Procedure' : 'Edit Procedure'}</h3>
        <button class="close-btn" on:click={() => showProcedureModal = false}>×</button>
      </div>
      <div class="modal-body">
        {#if procedureFormError}
          <div class="alert alert-error">{procedureFormError}</div>
        {/if}
        <div class="form-group">
          <label for="procedureName">Procedure Name</label>
          <input 
            type="text" 
            class="form-input" 
            placeholder="Enter procedure name"
            id="procedureName"
            bind:value={procedureForm.name}
            disabled={procedureFormLoading}
            on:keydown={(e) => e.key === 'Enter' && !procedureFormLoading && handleProcedureSave()}
          />
        </div>
        <div class="form-group">
          <label for="procedurePrice">Price (SYP)</label>
          <input 
            type="number" 
            class="form-input" 
            min="0" 
            step="1"
            placeholder="300000"
            id="procedurePrice"
            bind:value={procedureForm.priceDisplay}
            disabled={procedureFormLoading}
            on:keydown={(e) => e.key === 'Enter' && !procedureFormLoading && handleProcedureSave()}
          />
        </div>
      </div>
      <div class="modal-actions">
        <button class="btn btn-secondary" on:click={() => showProcedureModal = false} disabled={procedureFormLoading}>Cancel</button>
        <button class="btn btn-primary" on:click={handleProcedureSave} disabled={procedureFormLoading}>
          {#if procedureFormLoading}
            Saving...
          {:else if procedureModalMode === 'create'}
            Create Procedure
          {:else}
            Save Changes
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Delete Procedure Confirmation -->
{#if showProcedureDeleteConfirm}
  <div 
    class="modal-overlay" 
    role="button"
    tabindex="0"
    on:click={() => showProcedureDeleteConfirm = false}
    on:keydown={(e) => (e.key === 'Enter' || e.key === ' ') && (showProcedureDeleteConfirm = false)}
  >
    <div class="confirmation-modal" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header danger-header">
        <h3>Delete Procedure</h3>
      </div>
      <div class="modal-content">
        <p>Are you sure you want to delete the procedure <strong>{procedureToDelete?.name}</strong>? This action cannot be undone.</p>
      </div>
      <div class="modal-actions">
        <button class="btn btn-danger" on:click={handleDeleteProcedure}>Delete</button>
        <button class="btn btn-secondary" on:click={() => { showProcedureDeleteConfirm = false; procedureToDelete = null; }}>Cancel</button>
      </div>
    </div>
  </div>
{/if}

<!-- Work Type Modal -->
{#if showWorkTypeModal}
  <div 
    class="modal-overlay" 
    role="button"
    tabindex="0"
    on:click={() => showWorkTypeModal = false}
    on:keydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        showWorkTypeModal = false;
      }
    }}
  >
    <div class="modal-content small" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header">
        <h3>{workTypeModalMode === 'create' ? 'Create New Work Type' : 'Edit Work Type'}</h3>
        <button class="close-btn" on:click={() => showWorkTypeModal = false}>×</button>
      </div>
      <div class="modal-body">
        {#if workTypeFormError}
          <div class="alert alert-error">{workTypeFormError}</div>
        {/if}
        <div class="form-group">
          <label for="workTypeName">Name *</label>
          <input 
            type="text" 
            class="form-input" 
            placeholder="Enter work type name"
            id="workTypeName"
            bind:value={workTypeForm.name}
            disabled={workTypeFormLoading}
            on:keydown={(e) => e.key === 'Enter' && !workTypeFormLoading && handleWorkTypeSave()}
          />
        </div>
        <div class="form-group">
          <label for="workTypeDescription">Description</label>
          <textarea 
            class="form-textarea" 
            placeholder="Enter description (optional)"
            id="workTypeDescription"
            bind:value={workTypeForm.description}
            disabled={workTypeFormLoading}
            rows="3"
            on:keydown={(e) => {
              if (e.key === 'Enter' && e.ctrlKey && !workTypeFormLoading) {
                handleWorkTypeSave();
              }
            }}
          ></textarea>
        </div>
      </div>
      <div class="modal-actions">
        <button class="btn btn-secondary" on:click={() => showWorkTypeModal = false} disabled={workTypeFormLoading}>Cancel</button>
        <button class="btn btn-primary" on:click={handleWorkTypeSave} disabled={workTypeFormLoading}>
          {#if workTypeFormLoading}
            Saving...
          {:else if workTypeModalMode === 'create'}
            Create Work Type
          {:else}
            Save Changes
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Delete Work Type Confirmation -->
{#if showWorkTypeDeleteConfirm}
  <div 
    class="modal-overlay" 
    role="button"
    tabindex="0"
    on:click={() => showWorkTypeDeleteConfirm = false}
    on:keydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        showWorkTypeDeleteConfirm = false;
      }
    }}
  >
    <div class="confirmation-modal" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header danger-header">
        <h3>Delete Work Type</h3>
      </div>
      <div class="modal-content">
        <p>Are you sure you want to delete the work type <strong>{workTypeToDelete?.name}</strong>? This action cannot be undone.</p>
      </div>
      <div class="modal-actions">
        <button class="btn btn-danger" on:click={handleDeleteWorkType}>Delete</button>
        <button class="btn btn-secondary" on:click={() => { showWorkTypeDeleteConfirm = false; workTypeToDelete = null; }}>Cancel</button>
      </div>
    </div>
  </div>
{/if}

<!-- Color Shade Modal -->
{#if showColorShadeModal}
  <div 
    class="modal-overlay" 
    role="button"
    tabindex="0"
    on:click={() => showColorShadeModal = false}
    on:keydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        showColorShadeModal = false;
      }
    }}
  >
    <div class="modal-content small" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header">
        <h3>{colorShadeModalMode === 'create' ? 'Create New Color Shade' : 'Edit Color Shade'}</h3>
        <button class="close-btn" on:click={() => showColorShadeModal = false}>×</button>
      </div>
      <div class="modal-body">
        {#if colorShadeFormError}
          <div class="alert alert-error">{colorShadeFormError}</div>
        {/if}
        <div class="form-group">
          <label for="colorShadeName">Name *</label>
          <input 
            type="text" 
            class="form-input" 
            placeholder="Enter color shade name"
            id="colorShadeName"
            bind:value={colorShadeForm.name}
            disabled={colorShadeFormLoading}
            on:keydown={(e) => e.key === 'Enter' && !colorShadeFormLoading && handleColorShadeSave()}
          />
        </div>
        <div class="form-group">
          <label for="colorShadeDescription">Description</label>
          <textarea 
            class="form-textarea" 
            placeholder="Enter description (optional)"
            id="colorShadeDescription"
            bind:value={colorShadeForm.description}
            disabled={colorShadeFormLoading}
            rows="3"
            on:keydown={(e) => {
              if (e.key === 'Enter' && e.ctrlKey && !colorShadeFormLoading) {
                handleColorShadeSave();
              }
            }}
          ></textarea>
        </div>
        <div class="form-group">
          <label for="colorShadeColor">Color</label>
          <input 
            type="color" 
            class="form-input color-input" 
            id="colorShadeColor"
            bind:value={colorShadeForm.hex_color}
            disabled={colorShadeFormLoading}
          />
        </div>
        <div class="form-group">
          <label for="colorShadeActive">Active</label>
          <div class="toggle-group">
            <div class="toggle-item">
              <label for="colorShadeActive" class="toggle-label">
                <span>{colorShadeForm.is_active ? 'Active' : 'Inactive'}</span>
                <div class="toggle-switch">
                  <input
                    type="checkbox"
                    id="colorShadeActive"
                    bind:checked={colorShadeForm.is_active}
                    class="toggle-input"
                    disabled={colorShadeFormLoading}
                  />
                  <span class="toggle-slider"></span>
                </div>
              </label>
            </div>
          </div>
        </div>
      </div>
      <div class="modal-actions">
        <button class="btn btn-secondary" on:click={() => showColorShadeModal = false} disabled={colorShadeFormLoading}>Cancel</button>
        <button class="btn btn-primary" on:click={handleColorShadeSave} disabled={colorShadeFormLoading}>
          {#if colorShadeFormLoading}
            Saving...
          {:else if colorShadeModalMode === 'create'}
            Create Color Shade
          {:else}
            Save Changes
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Delete Color Shade Confirmation -->
{#if showColorShadeDeleteConfirm}
  <div 
    class="modal-overlay" 
    role="button"
    tabindex="0"
    on:click={() => showColorShadeDeleteConfirm = false}
    on:keydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        showColorShadeDeleteConfirm = false;
      }
    }}
  >
    <div class="confirmation-modal" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header danger-header">
        <h3>Delete Color Shade</h3>
      </div>
      <div class="modal-content">
        <p>Are you sure you want to delete the color shade <strong>{colorShadeToDelete?.name}</strong>? This action cannot be undone.</p>
      </div>
      <div class="modal-actions">
        <button class="btn btn-danger" on:click={handleDeleteColorShade}>Delete</button>
        <button class="btn btn-secondary" on:click={() => { showColorShadeDeleteConfirm = false; colorShadeToDelete = null; }}>Cancel</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .configuration {
    height: calc(100vh - 80px);
    background: var(--color-bg);
    color: var(--color-text);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .config-container {
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

  .danger-text {
    color: var(--color-danger);
  }

  .content-cards {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }

  /* Card Styles */
  .info-card,
  .action-card,
  .danger-card {
    background: var(--color-card);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: 1.5rem;
    box-shadow: var(--color-shadow);
  }

  .info-card h3,
  .action-card h3,
  .danger-card h3 {
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

  .info-list {
    display: flex;
    flex-direction: column;
    gap: 0.875rem;
  }

  .info-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.625rem 0;
    border-bottom: 1px solid var(--color-border);
  }

  .info-item:last-child {
    border-bottom: none;
  }

  .message-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.375rem;
  }

  .info-label {
    font-weight: 600;
    color: var(--color-text);
    opacity: 0.8;
    font-size: 0.875rem;
  }

  .info-value {
    color: var(--color-text);
    font-family: 'Courier New', monospace;
    font-size: 0.875rem;
    text-align: right;
  }

  .message-item .info-value {
    text-align: left;
  }

  .status-badge {
    display: inline-block;
    padding: 0.25rem 0.625rem;
    border-radius: 6px;
    font-size: 0.8125rem;
    font-weight: 600;
  }

  .status-badge.valid {
    background: rgba(76, 175, 80, 0.15);
    color: #4caf50;
  }

  .status-badge.expired {
    background: rgba(229, 115, 115, 0.15);
    color: var(--color-danger);
  }

  .status-badge.checking {
    background: rgba(255, 152, 0, 0.15);
    color: #ff9800;
  }

  /* Button Styles - Compact */
  .action-buttons {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
  }

  .btn {
    padding: 0.5rem 1rem;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
    border: none;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.375rem;
  }

  .btn-primary {
    background: var(--color-accent-gradient);
    color: white;
  }

  .btn-primary:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  }

  .btn-secondary {
    background: var(--color-panel);
    color: var(--color-accent);
    border: 1px solid var(--color-accent);
  }

  .btn-secondary:hover:not(:disabled) {
    background: var(--color-accent);
    color: white;
  }

  .btn-danger {
    background: linear-gradient(135deg, #ff6b6b, #ee5a52);
    color: white;
  }

  .btn-danger:hover:not(:disabled) {
    background: linear-gradient(135deg, #ff5252, #d32f2f);
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(255, 107, 107, 0.3);
  }

  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none !important;
  }

  .btn-small {
    padding: 0.375rem 0.75rem;
    font-size: 0.8125rem;
  }

  /* Form Styles */
  .license-form {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid var(--color-border);
  }

  .form-input {
    padding: 0.625rem;
    border-radius: 8px;
    border: 1px solid var(--color-border);
    font-size: 0.875rem;
    font-family: 'Courier New', monospace;
    background: var(--color-card);
    color: var(--color-text);
    transition: border-color 0.2s;
  }

  .form-input:focus {
    outline: none;
    border-color: var(--color-accent);
  }

  .form-input:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .danger-card {
    border-color: var(--color-danger);
  }

  .danger-card h3 {
    color: var(--color-danger);
  }

  /* Modal Styles */
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 3000;
    animation: fadeIn 0.2s ease;
  }

  .modal-content {
    background: var(--color-card);
    border-radius: 12px;
    width: 90%;
    max-width: 800px;
    max-height: 90vh;
    overflow: hidden;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
    display: flex;
    flex-direction: column;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.25rem 1.5rem;
    border-bottom: 1px solid var(--color-border);
    background: var(--color-accent-gradient);
    color: white;
  }

  .modal-header h3 {
    margin: 0;
    font-size: 1.125rem;
    color: white;
  }

  .modal-header .close-btn {
    background: rgba(255, 255, 255, 0.2);
    color: white;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    padding: 0;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    transition: background 0.2s;
    line-height: 1;
  }

  .modal-header .close-btn:hover {
    background: rgba(255, 255, 255, 0.3);
  }

  .modal-body {
    padding: 1.5rem;
    overflow-y: auto;
    flex: 1;
    background: var(--color-card);
    color: var(--color-text);
  }

  .confirmation-modal {
    background: var(--color-card);
    border-radius: 12px;
    width: 90%;
    max-width: 500px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  }

  .danger-header {
    background: linear-gradient(135deg, #ff6b6b, #ee5a52) !important;
  }

  .confirmation-modal .modal-content {
    padding: 1.5rem;
    color: var(--color-text);
  }

  .deletion-list {
    margin: 1rem 0;
    padding-left: 1.5rem;
  }

  .deletion-list li {
    margin: 0.5rem 0;
    color: var(--color-danger);
    font-weight: 500;
  }

  .final-warning {
    color: var(--color-danger);
    font-size: 1rem;
    text-align: center;
    margin-top: 1rem;
    font-weight: 600;
  }

  .modal-actions {
    display: flex;
    gap: 0.75rem;
    padding: 1.25rem 1.5rem;
    border-top: 1px solid var(--color-border);
  }

  .modal-actions .btn {
    flex: 1;
  }

  /* Procedure management styles */
  .procedures-toolbar {
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1.25rem;
  }

  .search-group {
    flex: 1;
    min-width: 220px;
    max-width: 320px;
  }

  .search-input {
    width: 100%;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    border: 1px solid var(--color-border);
    border-radius: 10px;
    overflow: hidden;
  }

  th, td {
    padding: 0.75rem 1rem;
    border-bottom: 1px solid var(--color-border);
    text-align: left;
  }

  th {
    background: var(--color-panel);
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--color-text);
  }

  tr:last-child td {
    border-bottom: none;
  }

  tr:hover td {
    background: rgba(0, 0, 0, 0.02);
  }

  .actions-col {
    text-align: right;
    width: 140px;
  }

  .icon-btn {
    border: none;
    background: var(--color-panel);
    border-radius: 8px;
    padding: 0.35rem 0.65rem;
    cursor: pointer;
    font-size: 0.95rem;
    transition: all 0.2s ease;
    margin-left: 0.35rem;
  }

  .icon-btn:hover {
    background: var(--color-accent);
    color: #fff;
  }

  .icon-btn.danger:hover {
    background: var(--color-danger);
    color: #fff;
  }

  .loading-state,
  .empty-state {
    padding: 2rem 1rem;
    text-align: center;
    color: var(--color-text);
    opacity: 0.75;
  }

  .pagination {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 1rem;
    margin-top: 2rem;
    padding-top: 1rem;
  }

  .page-btn {
    padding: 0.5rem 1rem;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 6px;
    cursor: pointer;
    font-size: 0.875rem;
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
    font-size: 0.875rem;
  }

  .spinner {
    width: 32px;
    height: 32px;
    border: 3px solid rgba(255, 255, 255, 0.3);
    border-top-color: var(--color-accent);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin: 0 auto 1rem;
  }

  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }

  .alert {
    padding: 0.75rem 1rem;
    border-radius: 10px;
    font-size: 0.9rem;
    margin-bottom: 1rem;
  }

  .alert-success {
    background: rgba(76, 175, 80, 0.15);
    color: #2e7d32;
  }

  .alert-error {
    background: rgba(229, 115, 115, 0.15);
    color: var(--color-danger);
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
    margin-bottom: 1rem;
  }

  .form-group label {
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--color-text);
  }

  .modal-content.small {
    max-width: 420px;
  }

  .form-textarea {
    width: 100%;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 0.6rem 1rem;
    font-size: 0.875rem;
    font-family: inherit;
    resize: vertical;
    min-height: 80px;
    box-sizing: border-box;
    transition: border-color 0.2s;
  }

  .form-textarea:focus {
    outline: none;
    border-color: var(--color-accent);
  }

  .form-textarea:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .color-input {
    height: 40px;
    cursor: pointer;
    padding: 0.25rem;
  }

  .color-sample {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    border: 2px solid var(--color-border);
    display: inline-block;
    cursor: pointer;
    transition: transform 0.2s ease;
  }

  .color-sample:hover {
    transform: scale(1.1);
  }

  .toggle-group {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .toggle-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .toggle-label {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    cursor: pointer;
    color: var(--color-text);
    font-size: 0.9rem;
  }

  .toggle-switch {
    position: relative;
    width: 44px;
    height: 24px;
  }

  .toggle-input {
    opacity: 0;
    width: 0;
    height: 0;
  }

  .toggle-slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: var(--color-border);
    transition: 0.3s;
    border-radius: 24px;
  }

  .toggle-slider:before {
    position: absolute;
    content: "";
    height: 18px;
    width: 18px;
    left: 3px;
    bottom: 3px;
    background-color: white;
    transition: 0.3s;
    border-radius: 50%;
  }

  .toggle-input:checked + .toggle-slider {
    background-color: var(--color-accent);
  }

  .toggle-input:checked + .toggle-slider:before {
    transform: translateX(20px);
  }

  .toggle-input:disabled + .toggle-slider {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>

