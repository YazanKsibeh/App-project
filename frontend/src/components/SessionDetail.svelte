<script>
  import { createEventDispatcher, onMount } from 'svelte';
  import { updateSession, deleteSession, loadSession } from '../stores/sessionStore.js';
  import { getInvoiceBySession } from '../stores/invoiceStore.js';
  import InvoiceConfirmationModal from './InvoiceConfirmationModal.svelte';

  export let session;

  const dispatch = createEventDispatcher();

  let isEditing = false;
  let editedSession = { ...session };
  let editedItems = session.items ? [...session.items] : [];
  let showDeleteConfirm = false;
  let isSaving = false;
  let isDeleting = false;
  let hasInvoice = false;
  let invoiceStatus = null; // Store invoice status: 'issued', 'paid', 'cancelled', or null
  let showInvoiceModal = false;
  let invoiceSuccess = null;

  function formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', { 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }

  function formatCurrency(amount) {
    return amount.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
  }

  function getStatusClass(status) {
    return status === 'completed' ? 'status-completed' : 'status-in-progress';
  }

  function getStatusText(status) {
    return status === 'completed' ? 'Completed' : 'In-progress';
  }

  function calculateTotal() {
    return editedItems.reduce((sum, item) => sum + item.amount, 0);
  }

  function handleEdit() {
    isEditing = true;
    editedSession = { ...session };
    editedItems = session.items ? session.items.map(item => ({
      procedure_id: item.procedure_id,
      item_name: item.item_name,
      amount: item.amount
    })) : [];
  }

  function handleCancel() {
    isEditing = false;
    editedSession = { ...session };
    editedItems = session.items ? [...session.items] : [];
  }

  async function handleSave() {
    isSaving = true;
    try {
      const itemsForm = editedItems.map(item => ({
        procedure_id: item.procedure_id,
        item_name: item.item_name,
        amount: item.amount
      }));

      const success = await updateSession(editedSession, itemsForm);
      if (success) {
        isEditing = false;
        // Reload the session data
        const updatedSession = await loadSession(session.id);
        if (updatedSession) {
          session = updatedSession;
          editedSession = { ...session };
          editedItems = session.items ? [...session.items] : [];
        }
      }
    } catch (error) {
      console.error('Error saving session:', error);
    } finally {
      isSaving = false;
    }
  }

  function handleDeleteClick() {
    showDeleteConfirm = true;
  }

  function handleDeleteCancel() {
    showDeleteConfirm = false;
  }

  async function handleDeleteConfirm() {
    isDeleting = true;
    try {
      const success = await deleteSession(session.id);
      if (success) {
        dispatch('close');
      }
    } catch (error) {
      console.error('Error deleting session:', error);
    } finally {
      isDeleting = false;
      showDeleteConfirm = false;
    }
  }

  function handleClose() {
    dispatch('close');
  }

  // Check if invoice exists on mount and when session ID changes
  let lastCheckedSessionId = null;
  
  onMount(async () => {
    await checkInvoiceStatus();
    lastCheckedSessionId = session.id;
  });

  // Re-check invoice status when session ID changes
  $: if (session && session.id && session.id !== lastCheckedSessionId) {
    lastCheckedSessionId = session.id;
    checkInvoiceStatus();
  }

  async function checkInvoiceStatus() {
    const result = await getInvoiceBySession(session.id);
    if (result.success && result.invoice) {
      hasInvoice = true;
      invoiceStatus = result.invoice.status || null;
    } else {
      hasInvoice = false;
      invoiceStatus = null;
    }
  }

  // Check if invoice is active (not cancelled)
  $: isInvoiceActive = hasInvoice && invoiceStatus && invoiceStatus !== 'cancelled';

  function handleCreateInvoice() {
    showInvoiceModal = true;
  }

  function handleInvoiceModalClose() {
    showInvoiceModal = false;
  }

  async function handleInvoiceConfirmed(event) {
    const { invoice } = event.detail;
    hasInvoice = true;
    invoiceStatus = invoice.status || null;
    invoiceSuccess = `Invoice ${invoice.invoice_number} created successfully!`;
    // Clear success message after 3 seconds
    setTimeout(() => {
      invoiceSuccess = null;
    }, 3000);
    // Reload session to get updated data
    const updatedSession = await loadSession(session.id);
    if (updatedSession) {
      session = updatedSession;
    }
    // Re-check invoice status after reload
    await checkInvoiceStatus();
  }

  function handleViewInvoice() {
    // Placeholder for future invoice viewing functionality
    console.log('View invoice clicked - to be implemented in Income section');
  }

  function addItem() {
    editedItems.push({
      procedure_id: null,
      item_name: '',
      amount: 0
    });
  }

  function removeItem(index) {
    editedItems.splice(index, 1);
  }

  function updateItem(index, field, value) {
    if (field === 'amount') {
      editedItems[index].amount = parseInt(value) || 0;
    } else {
      editedItems[index][field] = value;
    }
    editedSession.total_amount = calculateTotal();
  }
</script>

<div class="modal-overlay" on:click={handleClose}>
  <div class="modal-content" on:click|stopPropagation>
    <div class="modal-header">
      <h2>Session Details</h2>
      <button class="close-btn" on:click={handleClose} title="Close">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18"/>
          <line x1="6" y1="6" x2="18" y2="18"/>
        </svg>
      </button>
    </div>

    {#if isEditing}
      <div class="session-form">
        <div class="form-group">
          <label>Patient</label>
          <div class="info-value">{session.patient_name || 'Unknown'}</div>
        </div>

        <div class="form-group">
          <label>Dentist</label>
          <div class="info-value">{session.dentist_name || 'Unknown'}</div>
        </div>

        <div class="form-group">
          <label>Date & Time</label>
          <input 
            type="datetime-local" 
            bind:value={editedSession.session_date}
            class="form-input"
          />
        </div>

        <div class="form-group">
          <label>Status</label>
          <div class="radio-group">
            <label class="radio-label">
              <input 
                type="radio" 
                value="completed" 
                bind:group={editedSession.status}
              />
              <span>Completed</span>
            </label>
            <label class="radio-label">
              <input 
                type="radio" 
                value="in-progress" 
                bind:group={editedSession.status}
              />
              <span>In-progress</span>
            </label>
          </div>
        </div>

        <div class="form-group">
          <label>Procedures</label>
          <div class="items-list">
            {#each editedItems as item, index}
              <div class="item-row">
                <input 
                  type="text" 
                  placeholder="Procedure name"
                  bind:value={item.item_name}
                  class="form-input item-name"
                />
                <input 
                  type="number" 
                  placeholder="Amount (cents)"
                  value={item.amount}
                  on:input={(e) => updateItem(index, 'amount', e.target.value)}
                  class="form-input item-amount"
                />
                <button 
                  class="btn-remove" 
                  on:click={() => removeItem(index)}
                  type="button"
                >
                  ×
                </button>
              </div>
            {/each}
            <button class="btn-add-item" on:click={addItem} type="button">
              + Add Procedure
            </button>
          </div>
        </div>

        <div class="form-group">
          <label>Total Amount</label>
          <div class="total-amount">{formatCurrency(calculateTotal())} SYP</div>
        </div>

        <div class="form-group">
          <label>Notes</label>
          <textarea 
            bind:value={editedSession.notes}
            class="form-textarea"
            rows="4"
            placeholder="Session notes..."
          ></textarea>
        </div>

        <div class="form-actions">
          <button class="btn btn-secondary" on:click={handleCancel} disabled={isSaving}>
            Cancel
          </button>
          <button class="btn btn-primary" on:click={handleSave} disabled={isSaving}>
            {#if isSaving}
              Saving...
            {:else}
              Save Changes
            {/if}
          </button>
        </div>
      </div>
    {:else}
      <div class="session-details">
        <div class="detail-section">
          <div class="detail-row">
            <span class="label">Patient</span>
            <span class="value">{session.patient_name || 'Unknown'}</span>
          </div>
          <div class="detail-row">
            <span class="label">Dentist</span>
            <span class="value">{session.dentist_name || 'Unknown'}</span>
          </div>
          <div class="detail-row">
            <span class="label">Date & Time</span>
            <span class="value">{formatDate(session.session_date)}</span>
          </div>
          <div class="detail-row">
            <span class="label">Status</span>
            <span class="status-badge {getStatusClass(session.status)}">
              {getStatusText(session.status)}
            </span>
          </div>
        </div>

        <div class="detail-section">
          <h3>Procedures</h3>
          {#if session.items && session.items.length > 0}
            <div class="items-list">
              {#each session.items as item}
                <div class="item-row">
                  <span class="item-name">{item.item_name}</span>
                  <span class="item-amount">{formatCurrency(item.amount)} SYP</span>
                </div>
              {/each}
            </div>
          {:else}
            <p class="no-items">No procedures recorded</p>
          {/if}
        </div>

        <div class="detail-section">
          <div class="total-row">
            <span class="label">Total Amount</span>
            <span class="total-value">{formatCurrency(session.total_amount)} SYP</span>
          </div>
        </div>

        {#if session.notes}
          <div class="detail-section">
            <h3>Notes</h3>
            <p class="notes-text">{session.notes}</p>
          </div>
        {/if}

        {#if invoiceSuccess}
          <div class="invoice-success-message">
            {invoiceSuccess}
          </div>
        {/if}

        <div class="detail-actions">
          <button class="btn btn-secondary" on:click={handleEdit} disabled={isInvoiceActive}>
            Edit Session
          </button>
          {#if hasInvoice}
            <button class="btn btn-primary" on:click={handleViewInvoice}>
              View Invoice
            </button>
          {:else}
            <button class="btn btn-primary" on:click={handleCreateInvoice}>
              Create Invoice
            </button>
          {/if}
          <button class="btn btn-danger" on:click={handleDeleteClick} disabled={isDeleting || isInvoiceActive}>
            Delete Session
          </button>
          <button class="btn btn-secondary" on:click={handleClose}>
            Close
          </button>
        </div>
      </div>
    {/if}

    {#if showDeleteConfirm}
      <div class="delete-confirm-overlay">
        <div class="delete-confirm">
          <h3>Delete Session?</h3>
          <p>Are you sure you want to delete this session? This action cannot be undone.</p>
          <div class="confirm-actions">
            <button class="btn btn-secondary" on:click={handleDeleteCancel} disabled={isDeleting}>
              Cancel
            </button>
            <button class="btn btn-danger" on:click={handleDeleteConfirm} disabled={isDeleting}>
              {#if isDeleting}
                Deleting...
              {:else}
                Delete
              {/if}
            </button>
          </div>
        </div>
      </div>
    {/if}

    <InvoiceConfirmationModal 
      sessionId={session.id}
      bind:open={showInvoiceModal}
      on:close={handleInvoiceModalClose}
      on:confirmed={handleInvoiceConfirmed}
    />
  </div>
</div>

<style>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 1000;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
  }

  .modal-content {
    background: var(--color-card);
    color: var(--color-text);
    border-radius: 16px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
    max-width: 700px;
    width: 100%;
    max-height: 90vh;
    overflow-y: auto;
    position: relative;
    border: 1px solid var(--color-border);
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1.5rem;
    border-bottom: 1px solid var(--color-border);
  }

  .modal-header h2 {
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

  .session-details,
  .session-form {
    padding: 1.5rem;
  }

  .detail-section {
    margin-bottom: 2rem;
  }

  .detail-section h3 {
    margin: 0 0 1rem 0;
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--color-text);
  }

  .detail-row {
    display: flex;
    justify-content: space-between;
    padding: 0.75rem 0;
    border-bottom: 1px solid var(--color-border);
  }

  .detail-row:last-child {
    border-bottom: none;
  }

  .label {
    font-weight: 500;
    color: var(--color-text);
    opacity: 0.7;
  }

  .value {
    font-weight: 500;
    color: var(--color-text);
  }

  .status-badge {
    display: inline-block;
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

  .items-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .item-row {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.75rem;
    background: var(--color-panel);
    border-radius: 8px;
    border: 1px solid var(--color-border);
  }

  .item-name {
    flex: 1;
  }

  .item-amount {
    font-weight: 600;
    color: var(--color-accent);
    min-width: 100px;
    text-align: right;
  }

  .no-items {
    color: var(--color-text);
    opacity: 0.6;
    font-style: italic;
  }

  .total-row {
    display: flex;
    justify-content: space-between;
    padding: 1rem;
    background: var(--color-panel);
    border-radius: 8px;
    border: 2px solid var(--color-accent);
  }

  .total-value {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--color-accent);
  }

  .notes-text {
    color: var(--color-text);
    line-height: 1.6;
    white-space: pre-wrap;
  }

  .invoice-success-message {
    padding: 0.75rem 1rem;
    background: rgba(76, 175, 80, 0.15);
    color: #2e7d32;
    border-radius: 8px;
    margin-bottom: 1rem;
    font-size: 0.9rem;
    text-align: center;
  }

  .detail-actions,
  .form-actions {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
    margin-top: 2rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--color-border);
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

  .info-value {
    padding: 0.75rem;
    background: var(--color-panel);
    border-radius: 8px;
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

  .item-name {
    flex: 1;
  }

  .item-amount {
    width: 150px;
  }

  .btn-remove {
    background: var(--color-danger);
    color: white;
    border: none;
    border-radius: 4px;
    width: 32px;
    height: 32px;
    cursor: pointer;
    font-size: 1.25rem;
    line-height: 1;
  }

  .btn-remove:hover {
    opacity: 0.8;
  }

  .btn-add-item {
    padding: 0.75rem;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px dashed var(--color-border);
    border-radius: 8px;
    cursor: pointer;
    width: 100%;
    font-size: 1rem;
  }

  .btn-add-item:hover {
    background: var(--color-border);
  }

  .total-amount {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--color-accent);
    padding: 0.75rem;
    background: var(--color-panel);
    border-radius: 8px;
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

  .btn-danger {
    background: #ef4444;
    color: white;
  }

  .btn-danger:hover:not(:disabled) {
    background: #dc2626;
  }

  .delete-confirm-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.7);
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 16px;
  }

  .delete-confirm {
    background: var(--color-card);
    padding: 2rem;
    border-radius: 12px;
    border: 1px solid var(--color-border);
    max-width: 400px;
    width: 90%;
  }

  .delete-confirm h3 {
    margin: 0 0 1rem 0;
    color: var(--color-text);
  }

  .delete-confirm p {
    margin: 0 0 1.5rem 0;
    color: var(--color-text);
    opacity: 0.8;
  }

  .confirm-actions {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
  }
</style>

