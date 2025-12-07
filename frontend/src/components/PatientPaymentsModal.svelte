<script>
import { createEventDispatcher, onMount } from 'svelte';
import { GetPaymentsForPatient, GetPatientBalance, UpdateTotalRequired, AddPayment } from '../../wailsjs/go/main/App';
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

onMount(async () => {
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
        [totalRequired, totalPaid, remaining] = Array.isArray(balanceResult) ? balanceResult : [0, 0, 0];
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
            patient_id: patient.id,
            amount: parseInt(newAmount),
            payment_date: newDate,
            note: newNote
        }, licenseKey);
        newAmount = '';
        newDate = '';
        newNote = '';
        await loadPayments();
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
    } catch (e) {
        error = e.message || 'Failed to update required amount';
    }
}
</script>

<div class="modal-backdrop" on:click={() => dispatch('close')}></div>
<div class="modal">
    <h2>Payments for {patient.name}</h2>
    <div class="summary-row">
        <div><b>Phone:</b> {patient.phone}</div>
        <div><b>Total Required:</b> {totalRequired} <button class="edit-btn" on:click={() => updatingRequired = !updatingRequired}>Edit</button></div>
        <div><b>Total Paid:</b> {totalPaid}</div>
        <div><b>Remaining:</b> {remaining}</div>
    </div>
    {#if updatingRequired}
        <div class="update-required-row">
            <input type="number" bind:value={newRequired} min="0" />
            <button class="save-btn" on:click={handleUpdateRequired}>Save</button>
            <button class="cancel-btn" on:click={() => { updatingRequired = false; newRequired = totalRequired; }}>Cancel</button>
        </div>
    {/if}
    <h3>Payment History</h3>
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
                </tr>
            </thead>
            <tbody>
                {#each payments as p}
                    <tr>
                        <td>{p.amount}</td>
                        <td>{p.payment_date ? new Date(p.payment_date).toLocaleDateString() : '—'}</td>
                        <td>{p.note}</td>
                    </tr>
                {/each}
            </tbody>
        </table>
    {/if}
    <h3>Add Payment</h3>
    <form class="add-payment-form" on:submit|preventDefault={handleAddPayment}>
        <input type="number" placeholder="Amount" bind:value={newAmount} min="1" required />
        <input type="date" bind:value={newDate} required />
        <input type="text" placeholder="Note (optional)" bind:value={newNote} />
        <button type="submit" class="save-btn">Add</button>
    </form>
    <button class="close-btn" on:click={() => dispatch('close')}>Close</button>
</div>

<style>
.modal-backdrop {
    position: fixed;
    top: 0; left: 0; right: 0; bottom: 0;
    background: rgba(0,0,0,0.3);
    z-index: 1000;
}
.modal {
    position: fixed;
    top: 50%; left: 50%;
    transform: translate(-50%, -50%);
    background: #fff;
    color: #222b45;
    padding: 2rem;
    border-radius: 12px;
    box-shadow: 0 8px 32px rgba(0,0,0,0.22);
    z-index: 1001;
    min-width: 340px;
    max-width: 480px;
    width: 100%;
    box-sizing: border-box;
    overflow-x: hidden;
}
.summary-row {
    display: flex;
    flex-wrap: wrap;
    gap: 1.2rem;
    margin-bottom: 1.2rem;
    font-size: 1.08rem;
    color: #222b45;
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
.cancel-btn, .close-btn {
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
    background: #f9fafe;
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
</style> 