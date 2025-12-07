<script>
import { onMount, createEventDispatcher } from 'svelte';
import { GetPaymentsForPatient, GetPatientBalance, UpdateTotalRequired, AddPayment, DeletePayment, UpdatePayment } from '../../wailsjs/go/main/App';
import { currentLicenseKey, account } from '../stores/settingsStore.js';
import { get } from 'svelte/store';

// Helper function to get current license key
function getLicenseKey() {
    // Try multiple sources to get the license key
    let licenseKey = '';
    
    // 1. Try currentLicenseKey store
    try {
        licenseKey = get(currentLicenseKey);
    } catch (e) {
        console.error('getLicenseKey - Error getting from currentLicenseKey store:', e);
    }
    
    // 2. If empty, try account store
    if (!licenseKey) {
        try {
            const accountData = get(account);
            licenseKey = accountData.licenseKey || '';
        } catch (e) {
            console.error('getLicenseKey - Error getting from account store:', e);
        }
    }
    
    // 3. If still empty, try localStorage directly
    if (!licenseKey) {
        try {
            licenseKey = localStorage.getItem('dentist_license_key') || '';
        } catch (e) {
            console.error('getLicenseKey - Error getting from localStorage:', e);
        }
    }
    
    return licenseKey;
}

export let patient = null;
const dispatch = createEventDispatcher();

let payments = [];
let totalRequired = 0;
let totalPaid = 0;
let remaining = 0;
let loading = true;
let error = '';

// Add payment form
let newAmount = '';
let newDate = '';
let newNote = '';
let updatingRequired = false;
let newRequired = '';

let editingPaymentId = null;
let editAmount = '';
let editDate = '';
let editNote = '';

let showEditTotalModal = false;

let lastPatientId = null;
$: if (patient && patient.id !== lastPatientId) {
    lastPatientId = patient.id;
    loadPayments();
}

// Set default date for new payments to today
function setDefaultDate() {
    const today = new Date();
    const yyyy = today.getFullYear();
    const mm = String(today.getMonth() + 1).padStart(2, '0');
    const dd = String(today.getDate()).padStart(2, '0');
    newDate = `${yyyy}-${mm}-${dd}`;
}

onMount(async () => {
    setDefaultDate();
    await loadPayments();
});

async function loadPayments() {
    loading = true;
    error = '';
    try {
        const licenseKey = getLicenseKey();
        
        if (!licenseKey) {
            throw new Error('No license key found. Please check your license settings.');
        }
        
        payments = (await GetPaymentsForPatient(patient.id, licenseKey)) || [];
        const balanceResult = await GetPatientBalance(patient.id, licenseKey);
        totalRequired = balanceResult?.total_required || 0;
        totalPaid = balanceResult?.total_paid || 0;
        remaining = balanceResult?.remaining || 0;
        newRequired = totalRequired;
    } catch (e) {
        error = e.message || 'Failed to load payments';
    } finally {
        loading = false;
    }
}

async function handleAddPayment() {
    if (!newAmount || !newDate) return;
    try {
        const licenseKey = getLicenseKey();
        
        if (!licenseKey) {
            throw new Error('No license key found. Please check your license settings.');
        }
        
        await AddPayment({
            patient_id: patient?.id,
            amount: parseInt(newAmount),
            payment_date: newDate,
            note: newNote
        }, licenseKey);
        newAmount = '';
        setDefaultDate();
        newNote = '';
        await loadPayments();
        dispatch('paymentsChanged');
    } catch (e) {
        error = e.message || 'Failed to add payment';
    }
}

async function handleUpdateRequired() {
    if (!newRequired) return;
    try {
        const licenseKey = getLicenseKey();
        
        if (!licenseKey) {
            throw new Error('No license key found. Please check your license settings.');
        }
        
        await UpdateTotalRequired(patient.id, parseInt(newRequired), licenseKey);
        updatingRequired = false;
        await loadPayments();
        dispatch('paymentsChanged');
    } catch (e) {
        error = e.message || 'Failed to update required amount';
    }
}

async function handleDeletePayment(paymentId) {
    if (confirm('Are you sure you want to delete this payment?')) {
        try {
            const licenseKey = getLicenseKey();
            
            if (!licenseKey) {
                throw new Error('No license key found. Please check your license settings.');
            }
            
            await DeletePayment(paymentId, licenseKey);
            await loadPayments();
            dispatch('paymentsChanged');
        } catch (e) {
            error = e.message || 'Failed to delete payment';
        }
    }
}

function startEditPayment(payment) {
    editingPaymentId = payment.id;
    editAmount = payment.amount.toString();
    editDate = getDateOnly(payment.payment_date);
    editNote = payment.note;
}

async function handleUpdatePayment(paymentId) {
    try {
        const licenseKey = getLicenseKey();
        
        if (!licenseKey) {
            throw new Error('No license key found. Please check your license settings.');
        }
        
        await UpdatePayment({
            id: paymentId,
            amount: parseInt(editAmount),
            payment_date: editDate,
            note: editNote,
            patient_id: patient?.id
        }, licenseKey);
        editingPaymentId = null;
        await loadPayments();
        dispatch('paymentsChanged');
    } catch (e) {
        error = e.message || 'Failed to update payment';
    }
}

function cancelEditPayment() {
    editingPaymentId = null;
}

function goBack() {
    dispatch('close');
}

function openEditTotalModal() {
    showEditTotalModal = true;
    newRequired = totalRequired;
}

function closeEditTotalModal() {
    showEditTotalModal = false;
    newRequired = totalRequired;
}

async function saveEditTotalModal() {
    await handleUpdateRequired();
    showEditTotalModal = false;
}

function formatNumber(n) {
    if (typeof n === 'number') return n.toLocaleString('en-US');
    if (!isNaN(Number(n))) return Number(n).toLocaleString('en-US');
    return n;
}

function getDateOnly(value) {
    if (!value) return '';
    const normalized = value.replace('T', ' ');
    return normalized.split(' ')[0];
}
</script>

<div class="patient-payments-page">
    <button class="back-btn" on:click={goBack}>&larr; Back to Payments</button>
    <h1 class="header">Payments for {patient.name}</h1>
    <div class="section payment-summary">
        <div class="summary-row">
            <div><b>Total Required:</b> {formatNumber(totalRequired)} <button class="edit-btn" on:click={openEditTotalModal}>Edit</button></div>
            <div><b>Total Paid:</b> {formatNumber(totalPaid)}</div>
            <div><b>Remaining:</b> {formatNumber(remaining)}</div>
        </div>
    </div>
    {#if showEditTotalModal}
        <div class="modal-backdrop" on:click={closeEditTotalModal}></div>
        <div class="modal-popup">
            <h3>Edit Total Required</h3>
            <input type="number" class="edit-total-input" bind:value={newRequired} min="0" />
            <div class="modal-actions">
                <button class="save-btn" on:click={saveEditTotalModal}>Save</button>
                <button class="cancel-btn" on:click={closeEditTotalModal}>Cancel</button>
            </div>
        </div>
    {/if}
    <div class="section payment-history">
        <h2>Payment History</h2>
        {#if loading}
            <p>Loading payments...</p>
        {:else if error}
            <div class="error-message">{error}</div>
        {:else}
            <table class="payments-table">
                <thead>
                    <tr>
                        <th>Amount</th>
                        <th>Date</th>
                        <th>Note</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {#each payments as p}
                        <tr>
                            {#if editingPaymentId === p.id}
                                <td><input type="number" class="edit-input" bind:value={editAmount} min="1" /></td>
                                <td><input type="date" class="edit-input" bind:value={editDate} /></td>
                                <td><input type="text" class="edit-input" bind:value={editNote} /></td>
                                <td>
                                    <button class="save-btn" on:click={() => handleUpdatePayment(p.id)}>Save</button>
                                    <button class="cancel-btn" on:click={cancelEditPayment}>Cancel</button>
                                </td>
                            {:else}
                                <td>{formatNumber(p.amount)}</td>
                                <td>{p.payment_date ? new Date(p.payment_date).toLocaleDateString() : '—'}</td>
                                <td>{p.note}</td>
                                <td>
                                    <button class="edit-btn" on:click={() => startEditPayment(p)}>Edit</button>
                                    <button class="delete-btn" on:click={() => handleDeletePayment(p.id)}>Delete</button>
                                </td>
                            {/if}
                        </tr>
                    {/each}
                </tbody>
            </table>
        {/if}
    </div>
    <div class="section add-payment">
        <h2>Add Payment</h2>
        <form class="add-payment-form" on:submit|preventDefault={handleAddPayment}>
            <input type="number" placeholder="Amount" bind:value={newAmount} min="1" required />
            <input type="date" bind:value={newDate} required />
            <input type="text" placeholder="Note (optional)" bind:value={newNote} />
            <button type="submit" class="save-btn">Add</button>
        </form>
    </div>
</div>

<style>
.patient-payments-page {
    background: #fff;
    color: #222b45;
    border-radius: 14px;
    max-width: 100%;
    width: 100%;
    margin: 2rem 0;
    padding: 2.5rem 2.5rem 2rem 2.5rem;
    box-shadow: 0 8px 32px rgba(0,0,0,0.10);
}
.header {
    text-align: center;
    font-size: 2rem;
    font-weight: 700;
    margin-bottom: 2rem;
    color: #222b45;
}
.back-btn {
    background: none;
    border: none;
    color: #667eea;
    font-size: 1.1rem;
    font-weight: 600;
    cursor: pointer;
    margin-bottom: 1.5rem;
    transition: color 0.2s;
}
.back-btn:hover {
    color: #5a67d8;
    text-decoration: underline;
}
.section {
    margin-bottom: 2.2rem;
    padding: 1.2rem 1rem;
    background: #f9fafe;
    border-radius: 10px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}
.summary-row {
    display: flex;
    gap: 2.5rem;
    align-items: center;
    justify-content: space-between;
}
.edit-total-input {
    width: 120px;
    margin-right: 0.5rem;
    padding: 0.3rem 0.7rem;
    border-radius: 6px;
    border: 1px solid #ddd;
    font-size: 1rem;
}
.edit-input {
    width: 100px;
    padding: 0.3rem 0.7rem;
    border-radius: 6px;
    border: 1px solid #ddd;
    font-size: 1rem;
}
.edit-btn, .save-btn {
    background: #667eea;
    color: white;
    border: none;
    border-radius: 6px;
    padding: 0.3rem 1rem;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    margin-left: 0.5rem;
}
.edit-btn:hover, .save-btn:hover {
    background: #5a67d8;
}
.delete-btn {
    background: #fff0f0;
    color: #c0392b;
    border: 1px solid #e57373;
    border-radius: 6px;
    padding: 0.3rem 1rem;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    margin-left: 0.5rem;
    transition: background 0.2s, color 0.2s;
}
.delete-btn:hover {
    background: #ffeaea;
    color: #fff;
    background: linear-gradient(90deg, #e57373 0%, #c0392b 100%);
}
.cancel-btn {
    background: #eee;
    color: #333;
    border: none;
    border-radius: 6px;
    padding: 0.3rem 1rem;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    margin-left: 0.5rem;
}
.payments-table {
    width: 100%;
    border-collapse: collapse;
    background: #fff;
    border-radius: 8px;
    margin-bottom: 1.5rem;
}
.payments-table th, .payments-table td {
    padding: 0.7rem 0.5rem;
    text-align: left;
    color: #222b45;
    font-size: 1.02rem;
}
.payments-table th {
    background: #e6eaf5;
    color: #222b45;
    font-weight: 700;
    font-size: 1.02rem;
}
.payments-table tr:not(:last-child) td {
    border-bottom: 1px solid #e0e6f0;
}
.add-payment-form {
    display: flex;
    gap: 0.7rem;
    margin-bottom: 1.2rem;
}
.add-payment-form input {
    flex: 1 1 0;
    padding: 0.5rem;
    border-radius: 6px;
    border: 1px solid #ddd;
    font-size: 1rem;
    background: #fff;
    color: #222b45;
}
.error-message {
    color: #e74c3c;
    font-weight: bold;
    margin: 1rem 0;
    text-align: center;
    background: #fff0f0;
    border-radius: 8px;
    padding: 1rem;
}
.modal-backdrop {
    position: fixed;
    top: 0; left: 0; right: 0; bottom: 0;
    background: rgba(0,0,0,0.25);
    z-index: 1000;
}
.modal-popup {
    position: fixed;
    top: 50%; left: 50%;
    transform: translate(-50%, -50%);
    background: #fff;
    color: #222b45;
    padding: 2rem 2.5rem 1.5rem 2.5rem;
    border-radius: 14px;
    box-shadow: 0 8px 32px rgba(0,0,0,0.18);
    z-index: 1001;
    min-width: 320px;
    max-width: 90vw;
    width: 350px;
    text-align: center;
}
.modal-popup h3 {
    margin-bottom: 1.2rem;
    font-size: 1.25rem;
    font-weight: 700;
}
.modal-actions {
    display: flex;
    gap: 1rem;
    justify-content: center;
    margin-top: 1.2rem;
}
</style> 