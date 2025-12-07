<script>
  import { createEventDispatcher, onMount } from 'svelte';
  import { createSession } from '../stores/sessionStore.js';
  import { patients, loadPatients } from '../stores/patientStore.js';
  import { procedures, loadProcedures } from '../stores/procedureStore.js';
  import { currentUser } from '../stores/authStore.js';
  import { get } from 'svelte/store';

  const dispatch = createEventDispatcher();

  let selectedPatient = null;
  let patientSearch = '';
  let sessionDate = '';
  let sessionStatus = 'completed';
  let sessionNotes = '';
  let sessionItems = [];
  let procedureSearch = '';
  let showProcedureDropdown = false;
  let isSaving = false;

  onMount(async () => {
    await loadPatients();
    await loadProcedures();
    
    // Set default date/time to now
    const now = new Date();
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
    sessionDate = now.toISOString().slice(0, 16);
  });

  $: filteredPatients = $patients.filter(p => 
    !patientSearch || 
    p.name.toLowerCase().includes(patientSearch.toLowerCase()) ||
    p.phone.includes(patientSearch)
  );

  $: filteredProcedures = $procedures.filter(p =>
    !procedureSearch ||
    p.name.toLowerCase().includes(procedureSearch.toLowerCase())
  );

  $: totalAmount = sessionItems.reduce((sum, item) => sum + (item.amount || 0), 0);

  function handlePatientSelect(patient) {
    selectedPatient = patient;
    patientSearch = patient.name;
    showProcedureDropdown = false;
  }

  function handlePatientSearchInput(e) {
    patientSearch = e.target.value;
    if (!patientSearch) {
      selectedPatient = null;
    }
  }

  function handleAddProcedure() {
    showProcedureDropdown = true;
  }

  function handleSelectProcedure(procedure) {
    sessionItems = [...sessionItems, {
      procedure_id: procedure.id || null,
      item_name: procedure.name,
      amount: procedure.price
    }];
    procedureSearch = '';
    showProcedureDropdown = false;
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

  function handleCancel() {
    console.log('[NewSessionPanel] handleCancel called');
    dispatch('close');
  }

  async function handleSave() {
    console.log('[NewSessionPanel] handleSave called');
    console.log('[NewSessionPanel] selectedPatient:', selectedPatient);
    console.log('[NewSessionPanel] sessionDate:', sessionDate);
    console.log('[NewSessionPanel] sessionItems:', sessionItems);
    
    if (!selectedPatient) {
      console.log('[NewSessionPanel] Validation failed: No patient selected');
      alert('Please select a patient');
      return;
    }

    if (!sessionDate) {
      console.log('[NewSessionPanel] Validation failed: No date selected');
      alert('Please select a date and time');
      return;
    }

    if (sessionItems.length === 0) {
      console.log('[NewSessionPanel] Validation failed: No procedures added');
      alert('Please add at least one procedure');
      return;
    }

    isSaving = true;
    try {
      const user = get(currentUser);
      console.log('[NewSessionPanel] Current user:', user);
      
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

      console.log('[NewSessionPanel] Session form to submit:', sessionForm);
      const success = await createSession(sessionForm);
      console.log('[NewSessionPanel] createSession returned:', success);
      
      if (success) {
        console.log('[NewSessionPanel] Dispatching sessionCreated event');
        dispatch('sessionCreated');
      } else {
        console.log('[NewSessionPanel] Session creation failed');
      }
    } catch (error) {
      console.error('[NewSessionPanel] Error creating session:', error);
      console.error('[NewSessionPanel] Error stack:', error.stack);
      alert('Failed to create session: ' + (error.message || 'Unknown error'));
    } finally {
      isSaving = false;
      console.log('[NewSessionPanel] handleSave completed');
    }
  }

  function formatCurrency(amount) {
    return amount.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
  }
</script>

<div class="panel-overlay" on:click={handleCancel}>
  <div class="panel-content" on:click|stopPropagation>
    <div class="panel-header">
      <h2>New Session</h2>
      <button class="close-btn" on:click={handleCancel} title="Close">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18"/>
          <line x1="6" y1="6" x2="18" y2="18"/>
        </svg>
      </button>
    </div>

    <div class="panel-body">
      <div class="form-group">
        <label>Patient *</label>
        <div class="patient-search-container">
          <input
            type="text"
            class="form-input"
            placeholder="Search by name or phone..."
            value={patientSearch}
            on:input={handlePatientSearchInput}
            on:focus={() => showProcedureDropdown = false}
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
          <div class="selected-patient">
            <span class="selected-label">Selected:</span>
            <span class="selected-value">{selectedPatient.name} ({selectedPatient.phone})</span>
          </div>
        {/if}
      </div>

      <div class="form-group">
        <label>Date & Time *</label>
        <input
          type="datetime-local"
          class="form-input"
          bind:value={sessionDate}
        />
      </div>

      <div class="form-group">
        <label>Status *</label>
        <div class="radio-group">
          <label class="radio-label">
            <input
              type="radio"
              value="completed"
              bind:group={sessionStatus}
            />
            <span>Completed</span>
          </label>
          <label class="radio-label">
            <input
              type="radio"
              value="in-progress"
              bind:group={sessionStatus}
            />
            <span>In-progress</span>
          </label>
        </div>
      </div>

      <div class="form-group">
        <label>Procedures *</label>
        <button 
          class="btn-add-procedure" 
          on:click={handleAddProcedure}
          type="button"
        >
          + Add Procedure
        </button>

        {#if showProcedureDropdown}
          <div class="procedure-dropdown-container">
            <input
              type="text"
              class="form-input"
              placeholder="Search procedures..."
              bind:value={procedureSearch}
              on:focus
            />
            <div class="dropdown">
              {#if filteredProcedures.length > 0}
                {#each filteredProcedures.slice(0, 10) as procedure}
                  <div
                    class="dropdown-item"
                    on:click={() => handleSelectProcedure(procedure)}
                  >
                    <span class="procedure-name">{procedure.name}</span>
                    <span class="procedure-price">{formatCurrency(procedure.price)} SYP</span>
                  </div>
                {/each}
              {:else}
                <div class="dropdown-item empty">No procedures found</div>
              {/if}
            </div>
          </div>
        {/if}

        {#if sessionItems.length > 0}
          <div class="items-list">
            {#each sessionItems as item, index}
              <div class="item-row">
                <input
                  type="text"
                  class="form-input item-name"
                  placeholder="Procedure name"
                  value={item.item_name}
                  on:input={(e) => handleUpdateItemName(index, e.target.value)}
                />
                <input
                  type="number"
                  class="form-input item-amount"
                  placeholder="Amount (SYP)"
                  value={item.amount}
                  on:input={(e) => handleUpdateItemAmount(index, e.target.value)}
                />
                <button
                  class="btn-remove"
                  on:click={() => handleRemoveItem(index)}
                  type="button"
                >
                  ×
                </button>
              </div>
            {/each}
          </div>
        {/if}
      </div>

      <div class="form-group">
        <label>Total Amount</label>
        <div class="total-amount">{formatCurrency(totalAmount)} SYP</div>
      </div>

      <div class="form-group">
        <label>Notes</label>
        <textarea
          class="form-textarea"
          rows="4"
          placeholder="Session notes..."
          bind:value={sessionNotes}
        ></textarea>
      </div>
    </div>

    <div class="panel-footer">
      <button class="btn btn-secondary" on:click={handleCancel} disabled={isSaving}>
        Cancel
      </button>
      <button class="btn btn-primary" on:click={handleSave} disabled={isSaving}>
        {#if isSaving}
          Saving...
        {:else}
          Save Session
        {/if}
      </button>
    </div>
  </div>
</div>

<style>
  .panel-overlay {
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 1000;
    display: flex;
    justify-content: flex-end;
  }

  .panel-content {
    background: var(--color-card);
    color: var(--color-text);
    width: 500px;
    max-width: 90vw;
    height: 100vh;
    display: flex;
    flex-direction: column;
    box-shadow: -4px 0 20px rgba(0, 0, 0, 0.3);
    border-left: 1px solid var(--color-border);
  }

  .panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1.5rem;
    border-bottom: 1px solid var(--color-border);
  }

  .panel-header h2 {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 600;
  }

  .close-btn {
    background: none;
    border: none;
    color: var(--color-text);
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 4px;
    transition: background 0.2s;
  }

  .close-btn:hover {
    background: var(--color-panel);
  }

  .close-btn svg {
    width: 20px;
    height: 20px;
  }

  .panel-body {
    flex: 1;
    overflow-y: auto;
    padding: 1.5rem;
  }

  .form-group {
    margin-bottom: 1.5rem;
  }

  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: var(--color-text);
  }

  .form-input,
  .form-textarea {
    width: 100%;
    padding: 0.75rem;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-size: 1rem;
    font-family: inherit;
    box-sizing: border-box;
  }

  .form-input:focus,
  .form-textarea:focus {
    outline: none;
    border-color: var(--color-accent);
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  .form-textarea {
    resize: vertical;
  }

  .patient-search-container {
    position: relative;
  }

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

  .selected-patient {
    margin-top: 0.5rem;
    padding: 0.75rem;
    background: var(--color-panel);
    border-radius: 8px;
    display: flex;
    gap: 0.5rem;
  }

  .selected-label {
    font-weight: 500;
    color: var(--color-text);
    opacity: 0.7;
  }

  .selected-value {
    font-weight: 500;
    color: var(--color-text);
  }

  .radio-group {
    display: flex;
    gap: 1.5rem;
  }

  .radio-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
  }

  .radio-label input[type="radio"] {
    cursor: pointer;
  }

  .btn-add-procedure {
    padding: 0.75rem 1rem;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    cursor: pointer;
    font-size: 1rem;
    margin-bottom: 1rem;
    transition: all 0.2s;
  }

  .btn-add-procedure:hover {
    background: var(--color-border);
  }

  .procedure-dropdown-container {
    position: relative;
    margin-bottom: 1rem;
  }

  .procedure-name {
    flex: 1;
    color: var(--color-text);
  }

  .procedure-price {
    font-weight: 600;
    color: var(--color-accent);
  }

  .items-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    margin-top: 1rem;
  }

  .item-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .item-name {
    flex: 1;
  }

  .item-amount {
    width: 150px;
  }

  .btn-remove {
    background: #ef4444;
    color: white;
    border: none;
    border-radius: 4px;
    width: 32px;
    height: 32px;
    cursor: pointer;
    font-size: 1.25rem;
    line-height: 1;
    flex-shrink: 0;
  }

  .btn-remove:hover {
    opacity: 0.8;
  }

  .total-amount {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--color-accent);
    padding: 0.75rem;
    background: var(--color-panel);
    border-radius: 8px;
    text-align: center;
  }

  .panel-footer {
    padding: 1.5rem;
    border-top: 1px solid var(--color-border);
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
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
</style>

