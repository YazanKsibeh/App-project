<script>
  import { onMount } from 'svelte';
  import {
    filteredDentalLabs,
    dentalLabsLoading,
    dentalLabsError,
    dentalLabsSuccess,
    dentalLabSearch,
    dentalLabsCurrentPage,
    dentalLabsTotalPages,
    loadDentalLabsPaginated,
    getDentalLab,
    createDentalLab,
    updateDentalLab,
    deleteDentalLab
  } from '../stores/dentalLabStore.js';

  let selectedSection = 'labs'; // 'labs', 'new-order', 'orders-list', 'tracking'

  // Lab modal state
  let showLabModal = false;
  let labModalMode = 'create'; // 'create' or 'edit'
  let labForm = {
    id: null,
    name: '',
    contact_person: '',
    phone_primary: '',
    phone_secondary: '',
    email: '',
    specialties: '',
    is_active: true,
    notes: ''
  };
  let labFormErrors = {};
  let labFormLoading = false;
  let showOptionalFields = false;
  let labToEdit = null;

  // Lab details modal state
  let showLabDetailModal = false;
  let selectedLab = null;

  // Delete confirmation state
  let showDeleteConfirm = false;
  let labToDelete = null;

  onMount(() => {
    selectedSection = 'labs';
    loadDentalLabsPaginated(1);
  });

  function selectSection(section) {
    selectedSection = section;
    if (section === 'labs') {
      loadDentalLabsPaginated(1);
    }
  }

  // Phone formatting (same as AddPatientModal)
  function formatPhone(event, field) {
    let value = event.target.value.replace(/\D/g, '');
    if (value.length > 10) {
      value = value.slice(0, 10);
    }
    labForm[field] = value;
    // Clear error when user types
    if (labFormErrors[field]) {
      labFormErrors[field] = '';
    }
  }

  // Validation
  function validateLabForm() {
    labFormErrors = {};
    
    if (!labForm.name.trim()) {
      labFormErrors.name = 'Name is required';
    }
    
    if (!labForm.contact_person.trim()) {
      labFormErrors.contact_person = 'Contact person is required';
    }
    
    if (!labForm.phone_primary.trim()) {
      labFormErrors.phone_primary = 'Phone primary is required';
    } else if (!/^\d{10}$/.test(labForm.phone_primary.replace(/\D/g, ''))) {
      labFormErrors.phone_primary = 'Please enter a valid 10-digit phone number';
    }
    
    if (labForm.phone_secondary && !/^\d{10}$/.test(labForm.phone_secondary.replace(/\D/g, ''))) {
      labFormErrors.phone_secondary = 'Please enter a valid 10-digit phone number';
    }
    
    return Object.keys(labFormErrors).length === 0;
  }

  // Open create lab modal
  function openCreateLabModal() {
    labModalMode = 'create';
    labForm = {
      id: null,
      name: '',
      contact_person: '',
      phone_primary: '',
      phone_secondary: '',
      email: '',
      specialties: '',
      is_active: true,
      notes: ''
    };
    labFormErrors = {};
    showOptionalFields = false;
    labToEdit = null;
    showLabModal = true;
  }

  // Open edit lab modal
  function openEditLabModal(lab) {
    labModalMode = 'edit';
    labToEdit = lab;
    labForm = {
      id: lab.id,
      name: lab.name || '',
      contact_person: lab.contact_person || '',
      phone_primary: lab.phone_primary || '',
      phone_secondary: lab.phone_secondary || '',
      email: lab.email || '',
      specialties: lab.specialties || '',
      is_active: lab.is_active !== undefined ? lab.is_active : true,
      notes: lab.notes || ''
    };
    labFormErrors = {};
    showOptionalFields = false;
    showLabModal = true;
  }

  // Handle lab save
  async function handleLabSave() {
    if (!validateLabForm()) {
      return;
    }
    
    labFormLoading = true;
    try {
      const formData = {
        name: labForm.name.trim(),
        contact_person: labForm.contact_person.trim(),
        phone_primary: labForm.phone_primary,
        phone_secondary: labForm.phone_secondary || '',
        email: labForm.email.trim() || '',
        specialties: labForm.specialties.trim() || '',
        is_active: labForm.is_active,
        notes: labForm.notes.trim() || ''
      };

      if (labModalMode === 'create') {
        await createDentalLab(formData);
      } else {
        await updateDentalLab(labForm.id, formData);
      }
      
      showLabModal = false;
      labForm = {
        id: null,
        name: '',
        contact_person: '',
        phone_primary: '',
        phone_secondary: '',
        email: '',
        specialties: '',
        is_active: true,
        notes: ''
      };
      labFormErrors = {};
      showOptionalFields = false;
    } catch (error) {
      console.error('Error saving lab:', error);
    } finally {
      labFormLoading = false;
    }
  }

  // Open lab details modal
  async function openLabDetailModal(lab) {
    try {
      const fullLab = await getDentalLab(lab.id);
      selectedLab = fullLab;
      showLabDetailModal = true;
    } catch (error) {
      console.error('Error loading lab details:', error);
    }
  }

  // Close lab details modal
  function closeLabDetailModal() {
    showLabDetailModal = false;
    selectedLab = null;
  }

  // Open edit from details modal
  function editFromDetails() {
    if (selectedLab) {
      closeLabDetailModal();
      openEditLabModal(selectedLab);
    }
  }

  // Confirm delete
  function confirmDeleteLab(lab) {
    labToDelete = lab;
    showDeleteConfirm = true;
  }

  // Handle delete
  async function handleDeleteLab() {
    if (!labToDelete) return;
    const success = await deleteDentalLab(labToDelete.id);
    if (success) {
      showDeleteConfirm = false;
      labToDelete = null;
    }
  }

  // Pagination
  function previousLabs() {
    const current = $dentalLabsCurrentPage;
    if (current > 1) {
      loadDentalLabsPaginated(current - 1);
    }
  }

  function nextLabs() {
    const current = $dentalLabsCurrentPage;
    const total = $dentalLabsTotalPages;
    if (current < total) {
      loadDentalLabsPaginated(current + 1);
    }
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
            <h1>Labs Management</h1>
            <p class="section-description">Manage your dental lab vendors.</p>
          </div>

          <div class="content-cards">
            <div class="action-card">
              <div class="labs-toolbar">
                <div class="search-group">
                  <input
                    type="text"
                    class="form-input search-input"
                    placeholder="Search labs..."
                    bind:value={$dentalLabSearch}
                  />
                </div>
                <button class="btn btn-primary" on:click={openCreateLabModal}>
                  ➕ Create New Lab
                </button>
              </div>

              {#if $dentalLabsError}
                <div class="alert alert-error">
                  {$dentalLabsError}
                </div>
              {/if}
              {#if $dentalLabsSuccess}
                <div class="alert alert-success">
                  {$dentalLabsSuccess}
                </div>
              {/if}

              <div class="labs-list">
                {#if $dentalLabsLoading}
                  <div class="loading-state">
                    <div class="spinner"></div>
                    <p>Loading labs...</p>
                  </div>
                {:else if $filteredDentalLabs.length === 0}
                  <div class="empty-state">
                    <p>No labs found. Add your first lab to get started.</p>
                  </div>
                {:else}
                  <div class="sessions-table-container">
                    <table class="sessions-table">
                      <thead>
                        <tr>
                          <th>Code</th>
                          <th>Name</th>
                          <th>Contact Person</th>
                          <th>Phone Primary</th>
                          <th>Phone Secondary</th>
                          <th class="actions-col">Actions</th>
                        </tr>
                      </thead>
                      <tbody>
                        {#each $filteredDentalLabs as lab}
                          <tr class="lab-row" on:click={() => openLabDetailModal(lab)}>
                            <td class="lab-code">{lab.code || '-'}</td>
                            <td class="lab-name">{lab.name}</td>
                            <td>{lab.contact_person || '-'}</td>
                            <td>{lab.phone_primary || '-'}</td>
                            <td>{lab.phone_secondary || '-'}</td>
                            <td class="actions-col" on:click|stopPropagation>
                              <button class="icon-btn" on:click={() => openEditLabModal(lab)} title="Edit">
                                ✏️
                              </button>
                              <button class="icon-btn danger" on:click={() => confirmDeleteLab(lab)} title="Delete">
                                🗑️
                              </button>
                            </td>
                          </tr>
                        {/each}
                      </tbody>
                    </table>
                  </div>
                  
                  {#if $dentalLabsTotalPages > 1}
                    <div class="pagination">
                      <button 
                        class="page-btn" 
                        disabled={$dentalLabsCurrentPage === 1 || $dentalLabsLoading} 
                        on:click={previousLabs}
                      >
                        Previous
                      </button>
                      <span class="page-info">
                        Page {$dentalLabsCurrentPage} of {$dentalLabsTotalPages}
                      </span>
                      <button 
                        class="page-btn" 
                        disabled={$dentalLabsCurrentPage >= $dentalLabsTotalPages || $dentalLabsLoading} 
                        on:click={nextLabs}
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

<!-- Create/Edit Lab Modal -->
{#if showLabModal}
  <div 
    class="modal-overlay" 
    role="button"
    tabindex="0"
    on:click={() => showLabModal = false}
    on:keydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        showLabModal = false;
      }
    }}
  >
    <div class="modal-content" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header">
        <h2>{labModalMode === 'create' ? 'Create New Dental Lab' : 'Edit Dental Lab'}</h2>
        <button class="close-btn" on:click={() => showLabModal = false} title="Close">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
      
      <form on:submit|preventDefault={handleLabSave} class="lab-form">
        <!-- Required Fields Section -->
        <div class="form-section">
          <div class="form-group">
            <label for="labName">Name *</label>
            <input
              id="labName"
              type="text"
              bind:value={labForm.name}
              class="form-input {labFormErrors.name ? 'error' : ''}"
              placeholder="Enter lab name"
              required
            />
            {#if labFormErrors.name}
              <span class="error-message">{labFormErrors.name}</span>
            {/if}
          </div>
          
          <div class="form-group">
            <label for="contactPerson">Contact Person *</label>
            <input
              id="contactPerson"
              type="text"
              bind:value={labForm.contact_person}
              class="form-input {labFormErrors.contact_person ? 'error' : ''}"
              placeholder="Enter contact person name"
              required
            />
            {#if labFormErrors.contact_person}
              <span class="error-message">{labFormErrors.contact_person}</span>
            {/if}
          </div>
          
          <div class="form-group">
            <label for="phonePrimary">Phone Primary *</label>
            <input
              id="phonePrimary"
              type="tel"
              bind:value={labForm.phone_primary}
              on:input={(e) => formatPhone(e, 'phone_primary')}
              class="form-input {labFormErrors.phone_primary ? 'error' : ''}"
              placeholder="0941414122"
              maxlength="10"
              required
            />
            {#if labFormErrors.phone_primary}
              <span class="error-message">{labFormErrors.phone_primary}</span>
            {/if}
          </div>
          
          <div class="form-group">
            <label for="labActive">Active</label>
            <div class="toggle-group">
              <div class="toggle-item">
                <label for="labActive" class="toggle-label">
                  <span>{labForm.is_active ? 'Active' : 'Inactive'}</span>
                  <div class="toggle-switch">
                    <input
                      type="checkbox"
                      id="labActive"
                      bind:checked={labForm.is_active}
                      class="toggle-input"
                      disabled={labFormLoading}
                    />
                    <span class="toggle-slider"></span>
                  </div>
                </label>
              </div>
            </div>
          </div>
        </div>

        <!-- Optional Fields Section (Collapsible) -->
        <div class="optional-fields-section">
          <button 
            type="button" 
            class="optional-fields-header"
            on:click={() => showOptionalFields = !showOptionalFields}
            aria-expanded={showOptionalFields}
          >
            <span class="section-title">Show Optional Fields {showOptionalFields ? '▼' : '▶'}</span>
          </button>
          
          {#if showOptionalFields}
            <div class="optional-fields-content">
              <div class="form-group">
                <label for="phoneSecondary">Phone Secondary</label>
                <input
                  id="phoneSecondary"
                  type="tel"
                  bind:value={labForm.phone_secondary}
                  on:input={(e) => formatPhone(e, 'phone_secondary')}
                  class="form-input {labFormErrors.phone_secondary ? 'error' : ''}"
                  placeholder="0941414122"
                  maxlength="10"
                />
                {#if labFormErrors.phone_secondary}
                  <span class="error-message">{labFormErrors.phone_secondary}</span>
                {/if}
              </div>
              
              <div class="form-group">
                <label for="labEmail">Email</label>
                <input
                  id="labEmail"
                  type="email"
                  bind:value={labForm.email}
                  class="form-input"
                  placeholder="info@lab.com"
                />
              </div>
              
              <div class="form-group">
                <label for="labSpecialties">Specialties</label>
                <textarea
                  id="labSpecialties"
                  bind:value={labForm.specialties}
                  class="form-textarea"
                  placeholder="Crowns, Bridges, Implants"
                  rows="3"
                ></textarea>
              </div>
              
              <div class="form-group">
                <label for="labNotes">Notes</label>
                <textarea
                  id="labNotes"
                  bind:value={labForm.notes}
                  class="form-textarea"
                  placeholder="Additional notes about the lab"
                  rows="3"
                ></textarea>
              </div>
            </div>
          {/if}
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn-secondary" on:click={() => showLabModal = false} disabled={labFormLoading}>
            Cancel
          </button>
          <button type="submit" class="btn btn-primary" disabled={labFormLoading}>
            {#if labFormLoading}
              Saving...
            {:else if labModalMode === 'create'}
              Create Lab
            {:else}
              Save Changes
            {/if}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- Lab Details Modal -->
{#if showLabDetailModal && selectedLab}
  <div 
    class="modal-overlay" 
    role="button"
    tabindex="0"
    on:click={closeLabDetailModal}
    on:keydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        closeLabDetailModal();
      }
    }}
  >
    <div class="modal-content" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header">
        <h2>Lab Details: {selectedLab.name} ({selectedLab.code})</h2>
        <button class="close-btn" on:click={closeLabDetailModal} title="Close">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
      
      <div class="lab-details-content">
        <div class="detail-row">
          <span class="detail-label">Code</span>
          <span class="detail-value">{selectedLab.code || '-'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Name</span>
          <span class="detail-value">{selectedLab.name || '-'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Contact Person</span>
          <span class="detail-value">{selectedLab.contact_person || '-'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Phone Primary</span>
          <span class="detail-value">{selectedLab.phone_primary || '-'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Phone Secondary</span>
          <span class="detail-value">{selectedLab.phone_secondary || '-'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Email</span>
          <span class="detail-value">{selectedLab.email || '-'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Specialties</span>
          <span class="detail-value">{selectedLab.specialties || '-'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Active</span>
          <span class="detail-value">
            {#if selectedLab.is_active}
              ● Active
            {:else}
              ○ Inactive
            {/if}
          </span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Notes</span>
          <span class="detail-value">{selectedLab.notes || '-'}</span>
        </div>
      </div>

      <div class="form-actions">
        <button type="button" class="btn btn-secondary" on:click={closeLabDetailModal}>
          Close
        </button>
        <button type="button" class="btn btn-primary" on:click={editFromDetails}>
          Edit
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Delete Confirmation Modal -->
{#if showDeleteConfirm}
  <div 
    class="modal-overlay" 
    role="button"
    tabindex="0"
    on:click={() => showDeleteConfirm = false}
    on:keydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        showDeleteConfirm = false;
      }
    }}
  >
    <div class="confirmation-modal" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header danger-header">
        <h3>Delete Lab</h3>
      </div>
      <div class="modal-content">
        <p>Are you sure you want to delete the lab <strong>{labToDelete?.name}</strong>? This action cannot be undone.</p>
      </div>
      <div class="modal-actions">
        <button class="btn btn-danger" on:click={handleDeleteLab}>Delete</button>
        <button class="btn btn-secondary" on:click={() => { showDeleteConfirm = false; labToDelete = null; }}>Cancel</button>
      </div>
    </div>
  </div>
{/if}

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

  /* Labs Management Styles */
  .labs-toolbar {
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
    color: #e53935;
  }

  .loading-state,
  .empty-state {
    padding: 2rem 1rem;
    text-align: center;
    color: var(--color-text);
    opacity: 0.75;
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

  /* Sessions Table Styles (matching Sessions.svelte) */
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

  .lab-row {
    cursor: pointer;
  }

  .lab-code {
    font-weight: 600;
    color: var(--color-accent);
  }

  .lab-name {
    font-weight: 500;
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
    background: #e53935;
    color: #fff;
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
    max-width: 600px;
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

  .modal-header h2 {
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

  .modal-header .close-btn svg {
    width: 18px;
    height: 18px;
  }

  .danger-header {
    background: linear-gradient(135deg, #ff6b6b, #ee5a52) !important;
  }

  /* Form Styles */
  .lab-form {
    padding: 1.5rem;
    overflow-y: auto;
    flex: 1;
    background: var(--color-card);
    color: var(--color-text);
  }

  .form-section {
    margin-bottom: 1.5rem;
  }

  .form-group {
    margin-bottom: 1.5rem;
    width: 100%;
    box-sizing: border-box;
  }

  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-text);
  }

  .form-input {
    width: 100%;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 0.6rem 1rem;
    font-size: 1rem;
    margin-bottom: 0.3rem;
    transition: border-color 0.2s, box-shadow 0.2s;
    box-sizing: border-box;
  }

  .form-input:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  .form-input.error {
    border-color: #e53935;
  }

  .form-input.error:focus {
    box-shadow: 0 0 0 3px rgba(229, 57, 53, 0.1);
  }

  .form-input:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .form-textarea {
    width: 100%;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 0.6rem 1rem;
    font-size: 1rem;
    font-family: inherit;
    resize: vertical;
    min-height: 80px;
    box-sizing: border-box;
    margin-bottom: 0.3rem;
  }

  .form-textarea:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  .error-message {
    color: #e53935;
    font-size: 0.875rem;
    margin-top: 0.25rem;
    display: block;
  }

  /* Optional Fields Section */
  .optional-fields-section {
    margin-top: 1.5rem;
    margin-bottom: 1rem;
  }

  .optional-fields-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    background: var(--color-panel);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 0.75rem 1rem;
    cursor: pointer;
    transition: all 0.2s ease;
    color: var(--color-text);
    font-size: 1rem;
    font-weight: 500;
    font-family: inherit;
    margin-bottom: 0;
  }

  .optional-fields-header:hover {
    background: var(--color-border);
    border-color: var(--color-accent);
  }

  .section-title {
    color: var(--color-text);
  }

  .optional-fields-content {
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid var(--color-border);
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

  /* Toggle Switches */
  .toggle-group {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-bottom: 1.5rem;
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
    font-size: 1rem;
  }

  .toggle-switch {
    position: relative;
    width: 50px;
    height: 26px;
    flex-shrink: 0;
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
    border-radius: 26px;
  }

  .toggle-slider:before {
    position: absolute;
    content: "";
    height: 20px;
    width: 20px;
    left: 3px;
    bottom: 3px;
    background-color: white;
    transition: 0.3s;
    border-radius: 50%;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }

  .toggle-input:checked + .toggle-slider {
    background-color: #667eea;
  }

  .toggle-input:checked + .toggle-slider:before {
    transform: translateX(24px);
  }

  .toggle-input:disabled + .toggle-slider {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .form-actions {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
    margin-top: 2rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--color-border);
  }

  /* Lab Details Modal */
  .lab-details-content {
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .detail-row {
    display: flex;
    justify-content: space-between;
    padding: 1rem;
    background: var(--color-panel);
    border-radius: 8px;
    border: 1px solid var(--color-border);
  }

  .detail-label {
    font-weight: 600;
    color: var(--color-text);
    opacity: 0.8;
  }

  .detail-value {
    color: var(--color-text);
    text-align: right;
  }

  /* Confirmation Modal */
  .confirmation-modal {
    background: var(--color-card);
    border-radius: 12px;
    width: 90%;
    max-width: 500px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  }

  .confirmation-modal .modal-content {
    padding: 1.5rem;
    color: var(--color-text);
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
</style>

