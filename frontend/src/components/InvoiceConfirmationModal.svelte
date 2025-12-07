<script>
  import { createEventDispatcher } from 'svelte';
  import { previewInvoice, createInvoice } from '../stores/invoiceStore.js';

  export let sessionId = null;
  export let open = false;

  const dispatch = createEventDispatcher();

  let preview = null;
  let loading = false;
  let error = null;
  let isCreating = false;

  // Watch for open changes and load preview
  $: if (open && sessionId) {
    loadPreview();
  }

  async function loadPreview() {
    loading = true;
    error = null;
    try {
      const result = await previewInvoice(sessionId);
      if (result.success) {
        preview = result.preview;
      } else {
        error = result.error || 'Failed to load invoice preview';
      }
    } catch (err) {
      console.error('Error loading preview:', err);
      error = err.message || 'Failed to load invoice preview';
    } finally {
      loading = false;
    }
  }

  function handleClose() {
    open = false;
    preview = null;
    error = null;
    dispatch('close');
  }

  async function handleConfirm() {
    if (!sessionId) return;
    
    isCreating = true;
    error = null;
    try {
      const result = await createInvoice(sessionId);
      if (result.success) {
        dispatch('confirmed', { invoice: result.invoice });
        handleClose();
      } else {
        error = result.error || 'Failed to create invoice';
      }
    } catch (err) {
      console.error('Error creating invoice:', err);
      error = err.message || 'Failed to create invoice';
    } finally {
      isCreating = false;
    }
  }

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
</script>

{#if open}
  <div class="modal-overlay" on:click={handleClose}>
    <div class="modal-content" on:click|stopPropagation>
      <div class="modal-header">
        <h2>Create Invoice Confirmation</h2>
        <button class="close-btn" on:click={handleClose} title="Close">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>

      <div class="modal-body">
        {#if loading}
          <div class="loading-state">
            <p>Loading invoice preview...</p>
          </div>
        {:else if error}
          <div class="error-message">
            <p>{error}</p>
            <button class="btn btn-secondary" on:click={loadPreview}>Retry</button>
          </div>
        {:else if preview}
          <div class="invoice-preview">
            <div class="preview-section">
              <div class="preview-row">
                <span class="label">Patient:</span>
                <span class="value">{preview.patient_name || 'Unknown'}</span>
              </div>
              <div class="preview-row">
                <span class="label">Session Date:</span>
                <span class="value">{formatDate(preview.session_date)}</span>
              </div>
              <div class="preview-row">
                <span class="label">Invoice Number:</span>
                <span class="value invoice-number">{preview.invoice_number}</span>
              </div>
            </div>

            <div class="preview-section">
              <h3>Treatment Procedures:</h3>
              {#if preview.procedures && preview.procedures.length > 0}
                <div class="procedures-list">
                  {#each preview.procedures as procedure}
                    <div class="procedure-item">
                      <span class="procedure-name">• {procedure.item_name}</span>
                      <span class="procedure-amount">{formatCurrency(procedure.amount)} SYP</span>
                    </div>
                  {/each}
                </div>
              {:else}
                <p class="no-procedures">No procedures recorded</p>
              {/if}
            </div>

            <div class="preview-section total-section">
              <div class="total-row">
                <span class="label">Total:</span>
                <span class="total-value">{formatCurrency(preview.total_amount)} SYP</span>
              </div>
            </div>
          </div>
        {/if}
      </div>

      {#if preview && !loading}
        <div class="modal-actions">
          <button class="btn btn-secondary" on:click={handleClose} disabled={isCreating}>
            Cancel
          </button>
          <button class="btn btn-primary" on:click={handleConfirm} disabled={isCreating}>
            {#if isCreating}
              Creating...
            {:else}
              Confirm & Create
            {/if}
          </button>
        </div>
      {/if}
    </div>
  </div>
{/if}

<style>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 2000;
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
    max-width: 600px;
    width: 100%;
    max-height: 90vh;
    overflow-y: auto;
    position: relative;
    border: 1px solid var(--color-border);
    display: flex;
    flex-direction: column;
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
    color: var(--color-text);
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

  .modal-body {
    padding: 1.5rem;
    flex: 1;
  }

  .loading-state {
    text-align: center;
    padding: 2rem;
    color: var(--color-text);
    opacity: 0.7;
  }

  .error-message {
    padding: 1rem;
    background: rgba(239, 68, 68, 0.1);
    border: 1px solid rgba(239, 68, 68, 0.3);
    border-radius: 8px;
    color: #ef4444;
  }

  .error-message p {
    margin: 0 0 1rem 0;
  }

  .invoice-preview {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .preview-section {
    padding: 1rem 0;
  }

  .preview-section h3 {
    margin: 0 0 1rem 0;
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--color-text);
  }

  .preview-row {
    display: flex;
    justify-content: space-between;
    padding: 0.75rem 0;
    border-bottom: 1px solid var(--color-border);
  }

  .preview-row:last-child {
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

  .invoice-number {
    font-weight: 600;
    color: var(--color-accent);
  }

  .procedures-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .procedure-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem;
    background: var(--color-panel);
    border-radius: 8px;
    border: 1px solid var(--color-border);
  }

  .procedure-name {
    flex: 1;
    color: var(--color-text);
  }

  .procedure-amount {
    font-weight: 600;
    color: var(--color-accent);
    min-width: 120px;
    text-align: right;
  }

  .no-procedures {
    color: var(--color-text);
    opacity: 0.6;
    font-style: italic;
    padding: 1rem;
    text-align: center;
  }

  .total-section {
    border-top: 2px solid var(--color-border);
    padding-top: 1.5rem;
    margin-top: 0.5rem;
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

  .modal-actions {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
    padding: 1.5rem;
    border-top: 1px solid var(--color-border);
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

