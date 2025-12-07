<script>
  import { createEventDispatcher } from 'svelte';
  import {
    paymentDetails,
    paymentDetailsLoading,
    paymentDetailsError,
    loadInvoicePayments,
    addInvoicePayment
  } from '../stores/paymentStore.js';

  const dispatch = createEventDispatcher();

  export let open = false;
  export let invoice = null;

  let amount = '';
  let paymentDate = new Date().toISOString().slice(0, 10);
  let note = '';
  let submitError = null;
  let isSubmitting = false;
  let successMessage = '';
  let lastInvoiceId = null;

  $: if (open && invoice?.id && invoice.id !== lastInvoiceId) {
    lastInvoiceId = invoice.id;
    initializeModal();
  }

  async function initializeModal() {
    amount = '';
    note = '';
    submitError = null;
    successMessage = '';
    paymentDate = new Date().toISOString().slice(0, 10);
    await loadInvoicePayments(invoice.id);
  }

  function handleClose() {
    open = false;
    submitError = null;
    successMessage = '';
    lastInvoiceId = null;
    dispatch('close');
  }

  function formatCurrency(value = 0) {
    const num = Number(value) || 0;
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
  }

  function formatDateTime(value) {
    if (!value) return 'Unknown date';
    const parsed = new Date(value);
    if (isNaN(parsed.getTime())) {
      return value;
    }
    return parsed.toLocaleString();
  }

  $: paymentInfo = $paymentDetails;
  $: loadingDetails = $paymentDetailsLoading;
  $: loadError = $paymentDetailsError;
  $: isPaidView = paymentInfo && paymentInfo.status === 'paid';
  $: allowPayments = paymentInfo?.allow_payments;
  $: totalAmount = paymentInfo?.total_amount ?? invoice?.total_amount ?? 0;
  $: totalPaid = paymentInfo?.total_paid ?? 0;
  $: remaining = paymentInfo?.remaining ?? Math.max(totalAmount - totalPaid, 0);
  $: hasPreviousPayments = totalPaid > 0;

  function validatePayment() {
    if (!allowPayments) {
      return 'Payments are not allowed for this invoice.';
    }

    const numericAmount = parseInt(amount, 10);
    if (isNaN(numericAmount) || numericAmount <= 0) {
      return 'Payment amount must be greater than zero.';
    }

    if (numericAmount > remaining) {
      return 'Payment amount cannot exceed remaining balance.';
    }

    if (!paymentDate) {
      return 'Payment date is required.';
    }

    const selectedDate = new Date(paymentDate);
    const today = new Date();
    selectedDate.setHours(0, 0, 0, 0);
    today.setHours(0, 0, 0, 0);

    if (selectedDate.getTime() > today.getTime()) {
      return 'Payment date cannot be in the future.';
    }

    return null;
  }

  async function handleAddPayment() {
    submitError = null;
    successMessage = '';

    const validationError = validatePayment();
    if (validationError) {
      submitError = validationError;
      return;
    }

    isSubmitting = true;

    const response = await addInvoicePayment({
      invoiceId: invoice.id,
      amount: parseInt(amount, 10),
      paymentDate,
      note
    });

    isSubmitting = false;

    if (!response.success) {
      submitError = response.error;
      return;
    }

    successMessage = 'Payment recorded successfully.';
    dispatch('paymentSuccess', { details: response.details });
    setTimeout(() => {
      handleClose();
    }, 800);
  }
</script>

{#if open}
  <div class="modal-overlay" on:click={handleClose}>
    <div class="payment-modal" on:click|stopPropagation>
      <header class="modal-header">
        <div>
          <p class="modal-eyebrow">{invoice?.invoice_number || 'Invoice'}</p>
          <h2>{isPaidView ? 'Payment Details' : 'Add Payment'}</h2>
        </div>
        <button class="close-btn" on:click={handleClose} aria-label="Close modal">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18" />
            <line x1="6" y1="6" x2="18" y2="18" />
          </svg>
        </button>
      </header>

      <section class="modal-body">
        {#if loadingDetails}
          <div class="state-card">
            <div class="spinner"></div>
            <p>Loading payment information...</p>
          </div>
        {:else if loadError}
          <div class="state-card error">
            <p>{loadError}</p>
            <button class="btn btn-secondary" on:click={initializeModal}>Retry</button>
          </div>
        {:else if paymentInfo}
          <div class="info-grid">
            <div class="info-card">
              <p class="label">Patient</p>
              <p class="value">{paymentInfo.patient_name || invoice?.patient_name || 'Unknown'}</p>
            </div>
            <div class="info-card">
              <p class="label">Invoice</p>
              <p class="value">{paymentInfo.invoice?.invoice_number || invoice?.invoice_number}</p>
            </div>
            <div class="info-card">
              <p class="label">Total Amount</p>
              <p class="value accent">{formatCurrency(totalAmount)} SYP</p>
            </div>
            <div class="info-card">
              <p class="label">Status</p>
              <span class={`status-badge ${paymentInfo.status}`}>
                <span class="dot"></span>
                {paymentInfo.status.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase())}
              </span>
            </div>
            {#if hasPreviousPayments}
              <div class="info-card">
                <p class="label">Previous Payments</p>
                <p class="value">{formatCurrency(totalPaid)} SYP</p>
              </div>
            {/if}
            <div class="info-card">
              <p class="label">Remaining</p>
              <p class="value {remaining === 0 ? 'success' : ''}">{formatCurrency(remaining)} SYP</p>
            </div>
          </div>

          {#if paymentInfo.payments && paymentInfo.payments.length > 0}
            <div class="history-section">
              <div class="history-header">
                <h3>Payment History</h3>
                <span>{paymentInfo.payments.length} {paymentInfo.payments.length === 1 ? 'entry' : 'entries'}</span>
              </div>
              <div class="history-list">
                {#each paymentInfo.payments as payment}
                  <div class="history-row">
                    <div>
                      <p class="history-amount">{formatCurrency(payment.amount)} SYP</p>
                      <p class="history-meta">Payment #{payment.payment_code || payment.id}</p>
                    </div>
                    <div class="history-details">
                      <p>{formatDateTime(payment.payment_date)}</p>
                      {#if payment.note}
                        <p class="history-note">{payment.note}</p>
                      {/if}
                    </div>
                  </div>
                {/each}
              </div>
            </div>
          {/if}

          {#if !isPaidView && allowPayments}
            <div class="form-section">
              <h3>Record Payment</h3>

              {#if submitError}
                <div class="alert error">
                  {submitError}
                </div>
              {/if}

              {#if successMessage}
                <div class="alert success">
                  {successMessage}
                </div>
              {/if}

              <div class="form-grid">
                <div class="form-field">
                  <label for="paymentAmount">Payment Amount</label>
                  <input
                    id="paymentAmount"
                    type="number"
                    min="1"
                    bind:value={amount}
                    placeholder="Enter payment amount"
                  />
                </div>

                <div class="form-field">
                  <label for="paymentDate">Payment Date</label>
                  <input
                    id="paymentDate"
                    type="date"
                    max={new Date().toISOString().slice(0, 10)}
                    bind:value={paymentDate}
                  />
                </div>
              </div>

              <div class="form-field">
                <label for="paymentNote">Notes (optional)</label>
                <textarea
                  id="paymentNote"
                  rows="3"
                  placeholder="Additional details about this payment"
                  bind:value={note}
                ></textarea>
              </div>
            </div>
          {/if}
        {/if}
      </section>

      <footer class="modal-footer">
        <button class="btn btn-secondary" on:click={handleClose}>Close</button>

        {#if !isPaidView && allowPayments}
          <button class="btn btn-primary" on:click={handleAddPayment} disabled={isSubmitting || loadingDetails}>
            {#if isSubmitting}
              Processing...
            {:else}
              Add Payment
            {/if}
          </button>
        {/if}
      </footer>
    </div>
  </div>
{/if}

<style>
  .modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.55);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    z-index: 2100;
  }

  .payment-modal {
    width: min(720px, 100%);
    max-height: 95vh;
    background: var(--color-card);
    border-radius: 20px;
    border: 1px solid var(--color-border);
    box-shadow: 0 25px 80px rgba(15, 23, 42, 0.4);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.75rem 2rem 1.25rem;
    border-bottom: 1px solid var(--color-border);
  }

  .modal-eyebrow {
    margin: 0;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    font-size: 0.75rem;
    color: var(--color-muted);
  }

  .modal-header h2 {
    margin: 0.25rem 0 0;
    font-size: 1.5rem;
  }

  .close-btn {
    border: none;
    background: transparent;
    color: var(--color-text);
    padding: 0.35rem;
    border-radius: 50%;
    cursor: pointer;
  }

  .close-btn:hover {
    background: rgba(255, 255, 255, 0.05);
  }

  .modal-body {
    padding: 1.5rem 2rem;
    overflow-y: auto;
  }

  .state-card {
    border: 1px dashed var(--color-border);
    border-radius: 16px;
    padding: 2rem;
    text-align: center;
    color: var(--color-text);
    display: flex;
    flex-direction: column;
    gap: 1rem;
    align-items: center;
  }

  .state-card.error {
    border-color: rgba(239, 68, 68, 0.4);
    color: #fca5a5;
  }

  .spinner {
    width: 32px;
    height: 32px;
    border: 3px solid rgba(255, 255, 255, 0.2);
    border-top-color: var(--color-accent);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .info-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
    gap: 1rem;
    margin-bottom: 1.5rem;
  }

  .info-card {
    padding: 1rem;
    border: 1px solid var(--color-border);
    border-radius: 14px;
    background: var(--color-panel);
  }

  .label {
    font-size: 0.8rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--color-muted);
    margin-bottom: 0.35rem;
  }

  .value {
    margin: 0;
    font-size: 1rem;
    font-weight: 600;
  }

  .value.accent {
    color: var(--color-accent);
    font-size: 1.1rem;
  }

  .value.success {
    color: #22c55e;
  }

  .status-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    padding: 0.4rem 0.8rem;
    border-radius: 999px;
    font-size: 0.85rem;
    font-weight: 600;
    text-transform: capitalize;
    background: rgba(148, 163, 184, 0.1);
    color: var(--color-text);
  }

  .status-badge.paid {
    background: rgba(34, 197, 94, 0.15);
    color: #4ade80;
  }

  .status-badge.partially_paid,
  .status-badge.issued {
    background: rgba(250, 204, 21, 0.15);
    color: #fde68a;
  }

  .dot {
    display: inline-block;
    width: 8px;
    height: 8px;
    border-radius: 999px;
    background: currentColor;
  }

  .history-section {
    border: 1px solid var(--color-border);
    border-radius: 16px;
    padding: 1.25rem;
    margin-bottom: 1.5rem;
    background: var(--color-panel);
  }

  .history-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .history-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    max-height: 220px;
    overflow-y: auto;
  }

  .history-row {
    padding: 0.85rem;
    border: 1px solid var(--color-border);
    border-radius: 12px;
    display: flex;
    justify-content: space-between;
    gap: 1rem;
  }

  .history-amount {
    margin: 0;
    font-weight: 600;
  }

  .history-meta {
    margin: 0.25rem 0 0;
    font-size: 0.85rem;
    color: var(--color-muted);
  }

  .history-details p {
    margin: 0;
    font-size: 0.9rem;
  }

  .history-note {
    margin-top: 0.25rem;
    font-size: 0.85rem;
    font-style: italic;
    opacity: 0.8;
  }

  .form-section {
    border: 1px solid var(--color-border);
    border-radius: 16px;
    padding: 1.25rem;
    background: var(--color-panel);
  }

  .form-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .form-field {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
  }

  label {
    font-size: 0.9rem;
    color: var(--color-text);
    opacity: 0.8;
  }

  input,
  textarea {
    width: 100%;
    border-radius: 10px;
    border: 1px solid var(--color-border);
    background: var(--color-input, rgba(15, 23, 42, 0.6));
    color: var(--color-text);
    padding: 0.7rem 0.85rem;
    font-size: 1rem;
    transition: border 0.2s ease, box-shadow 0.2s ease;
  }

  input:focus,
  textarea:focus {
    border-color: var(--color-accent);
    box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
    outline: none;
  }

  textarea {
    resize: vertical;
  }

  .alert {
    padding: 0.75rem 1rem;
    border-radius: 10px;
    margin-bottom: 1rem;
    font-size: 0.95rem;
  }

  .alert.error {
    background: rgba(239, 68, 68, 0.1);
    color: #f87171;
    border: 1px solid rgba(239, 68, 68, 0.4);
  }

  .alert.success {
    background: rgba(34, 197, 94, 0.1);
    color: #4ade80;
    border: 1px solid rgba(34, 197, 94, 0.4);
  }

  .modal-footer {
    padding: 1.25rem 2rem 1.75rem;
    border-top: 1px solid var(--color-border);
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
  }

  .btn {
    border: none;
    border-radius: 10px;
    padding: 0.75rem 1.5rem;
    font-size: 1rem;
    cursor: pointer;
    font-weight: 600;
    transition: transform 0.2s ease, opacity 0.2s ease;
  }

  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .btn-secondary {
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
  }

  .btn-primary {
    background: linear-gradient(135deg, #667eea 0%, #7f9cf5 100%);
    color: white;
  }
</style>

