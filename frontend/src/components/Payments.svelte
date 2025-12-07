<script>
import { onMount } from 'svelte';
import { GetPatients, GetLastPaymentForPatient, GetPatientBalance } from '../../wailsjs/go/main/App';
import { currentLicenseKey, account } from '../stores/settingsStore.js';
import { get } from 'svelte/store';
import PatientPaymentsPage from './PatientPaymentsPage.svelte';

// Helper function to get current license key
function getLicenseKey() {
    // Try multiple sources to get the license key
    let licenseKey = '';
    
    // 1. Try currentLicenseKey store
    try {
        licenseKey = get(currentLicenseKey);
        console.log('getLicenseKey - From currentLicenseKey store:', licenseKey ? licenseKey.substring(0, 20) + '...' : 'empty');
    } catch (e) {
        console.error('getLicenseKey - Error getting from currentLicenseKey store:', e);
    }
    
    // 2. If empty, try account store
    if (!licenseKey) {
        try {
            const accountData = get(account);
            licenseKey = accountData.licenseKey || '';
            console.log('getLicenseKey - From account store:', licenseKey ? licenseKey.substring(0, 20) + '...' : 'empty');
        } catch (e) {
            console.error('getLicenseKey - Error getting from account store:', e);
        }
    }
    
    // 3. If still empty, try localStorage directly
    if (!licenseKey) {
        try {
            licenseKey = localStorage.getItem('dentist_license_key') || '';
            console.log('getLicenseKey - From localStorage:', licenseKey ? licenseKey.substring(0, 20) + '...' : 'empty');
        } catch (e) {
            console.error('getLicenseKey - Error getting from localStorage:', e);
        }
    }
    
    console.log('getLicenseKey - Final result length:', licenseKey.length);
    return licenseKey;
}

let patients = [];
let paymentSummaries = [];
let loading = true;
let error = '';
let selectedPatient = null;
let filter = '';
let currentPage = 1;
const pageSize = 10;

onMount(async () => {
    loading = true;
    error = '';
    try {
        const licenseKey = getLicenseKey();
        console.log('[PAYMENTS] Loading payments data with license key length:', licenseKey.length);
        
        if (!licenseKey) {
            throw new Error('No license key found. Please check your license settings.');
        }
        
        const allPatients = await GetPatients(licenseKey);
        console.log('[PAYMENTS] Loaded patients count:', allPatients.length);
        
        // For each patient, get last payment and balance
        paymentSummaries = await Promise.all(
            allPatients.map(async (p) => {
                const lastPayment = await GetLastPaymentForPatient(p.id, licenseKey) || {};
                let balanceResult = await GetPatientBalance(p.id, licenseKey);
                if (!Array.isArray(balanceResult)) {
                    if (typeof balanceResult === 'object' && balanceResult !== null) {
                        balanceResult = Object.values(balanceResult);
                    } else {
                        balanceResult = [0, 0, 0];
                    }
                }
                const [totalRequired, totalPaid, remaining] = balanceResult;
                return {
                    id: p.id,
                    name: p.name,
                    phone: p.phone,
                    lastPaymentAmount: lastPayment.amount || 0,
                    lastPaymentDate: lastPayment.payment_date || '',
                    totalRequired,
                    remaining
                };
            })
        );
        patients = allPatients;
    } catch (e) {
        error = e.message || 'Failed to load payments data';
    } finally {
        loading = false;
    }
});

// Filter and paginate paymentSummaries
$: filteredSummaries = paymentSummaries.filter(p =>
    p.name.toLowerCase().includes(filter.toLowerCase()) ||
    p.phone.replace(/\D/g, '').includes(filter.replace(/\D/g, ''))
);
$: totalPages = Math.ceil(filteredSummaries.length / pageSize);
$: paginatedSummaries = filteredSummaries.slice((currentPage - 1) * pageSize, currentPage * pageSize);

function handleRowClick(patientId) {
    const patient = patients.find(p => p.id === patientId);
    if (patient) selectedPatient = patient;
}

function closePatientPage() {
    selectedPatient = null;
}

function goToPage(page) {
    if (page >= 1 && page <= totalPages) {
        currentPage = page;
    }
}

function formatNumber(n) {
    if (typeof n === 'number') return n.toLocaleString('en-US');
    if (!isNaN(Number(n))) return Number(n).toLocaleString('en-US');
    return n;
}

async function reloadPatients() {
    loading = true;
    error = '';
    try {
        const licenseKey = getLicenseKey();
        console.log('[PAYMENTS] Reloading payments data with license key length:', licenseKey.length);
        
        if (!licenseKey) {
            throw new Error('No license key found. Please check your license settings.');
        }
        
        const allPatients = await GetPatients(licenseKey);
        console.log('[PAYMENTS] Reloaded patients count:', allPatients.length);
        
        paymentSummaries = await Promise.all(
            allPatients.map(async (p) => {
                const lastPayment = await GetLastPaymentForPatient(p.id, licenseKey) || {};
                let balanceResult = await GetPatientBalance(p.id, licenseKey);
                if (!Array.isArray(balanceResult)) {
                    if (typeof balanceResult === 'object' && balanceResult !== null) {
                        balanceResult = Object.values(balanceResult);
                    } else {
                        balanceResult = [0, 0, 0];
                    }
                }
                const [totalRequired, totalPaid, remaining] = balanceResult;
                return {
                    id: p.id,
                    name: p.name,
                    phone: p.phone,
                    lastPaymentAmount: lastPayment.amount || 0,
                    lastPaymentDate: lastPayment.payment_date || '',
                    totalRequired,
                    remaining
                };
            })
        );
        patients = allPatients;
    } catch (e) {
        error = e.message || 'Failed to load payments data';
    } finally {
        loading = false;
    }
}
</script>

<div class="payments-split-view">
    <div class="patient-list-panel">
        <h2 class="panel-header">Patients</h2>
        <div class="payments-filter-bar">
            <input class="payments-filter-input" type="text" placeholder="Search by name or phone..." bind:value={filter} on:input={() => { currentPage = 1; }} />
        </div>
        {#if loading}
            <p>Loading patients...</p>
        {:else if error}
            <div class="error-message">{error}</div>
        {:else}
            <table class="payments-table">
                <thead>
                    <tr>
                        <th>Patient Name</th>
                        <th>Phone</th>
                        <th>Total Required</th>
                        <th>Remaining</th>
                    </tr>
                </thead>
                <tbody>
                    {#each paginatedSummaries as p}
                        <tr class="clickable-row {selectedPatient && selectedPatient.id === p.id ? 'selected' : ''}" on:click={() => handleRowClick(p.id)}>
                            <td>{p.name}</td>
                            <td>{p.phone}</td>
                            <td>{formatNumber(p.totalRequired)}</td>
                            <td>{formatNumber(p.remaining)}</td>
                        </tr>
                    {/each}
                </tbody>
            </table>
            <div class="pagination-bar">
                <button class="pagination-btn" on:click={() => goToPage(currentPage - 1)} disabled={currentPage === 1}>Previous</button>
                <span class="pagination-info">Page {currentPage} of {totalPages}</span>
                <button class="pagination-btn" on:click={() => goToPage(currentPage + 1)} disabled={currentPage === totalPages}>Next</button>
            </div>
        {/if}
    </div>
    <div class="divider"></div>
    <div class="payment-details-panel">
        {#if selectedPatient}
            <PatientPaymentsPage
                patient={selectedPatient}
                on:close={closePatientPage}
                on:paymentsChanged={reloadPatients}
            />
        {:else}
            <div class="empty-details">Select a patient to view payment details</div>
        {/if}
    </div>
</div>

<style>
.payments-split-view {
    display: flex;
    flex-direction: row;
    gap: 0;
    width: 100%;
    min-height: 80vh;
    max-width: 2000px;
    margin: 0 auto;
    background: none;
}
.patient-list-panel {
    flex: 0 0 35%;
    max-width: 30vw;
    background: var(--color-card);
    color: var(--color-text);
    border-radius: 18px 0 0 18px;
    box-shadow: 0 2px 12px rgba(0,0,0,0.08);
    padding: 2.5rem 2rem 2.5rem 2.5rem;
    min-height: 80vh;
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
}
.panel-header {
    font-size: 1.3rem;
    font-weight: 700;
    margin-bottom: 1.2rem;
    color: var(--color-text);
    letter-spacing: 0.5px;
}
.payments-table {
    width: 100%;
    border-collapse: collapse;
    background: var(--color-card);
    border-radius: 10px;
    overflow: hidden;
    margin-bottom: 1.2rem;
}
.payments-table th, .payments-table td {
    padding: 0.8rem 0.5rem;
    text-align: left;
    color: var(--color-text);
    font-size: 1.05rem;
}
.payments-table th {
    background: var(--color-panel);
    color: var(--color-text);
    font-weight: 700;
}
.payments-table td {
    color: var(--color-text);
}
.payments-table tr:not(:last-child) td {
    border-bottom: 1px solid var(--color-border);
}
.clickable-row {
    cursor: pointer;
    transition: background 0.15s;
}
.clickable-row.selected, .clickable-row:hover {
    background: var(--color-accent);
    color: #fff;
}
.payments-filter-bar {
    display: flex;
    justify-content: flex-start;
    margin-bottom: 1.5rem;
}
.payments-filter-input {
    padding: 1rem 1.5rem;
    border-radius: 10px;
    border: 1.5px solid var(--color-border);
    font-size: 1.2rem;
    width: 100%;
    max-width: 100%;
    outline: none;
    transition: border 0.2s;
    background: #fff;
    color: #222b45;
    box-sizing: border-box;
}
.payments-filter-input:focus {
    border: 1.5px solid var(--color-accent);
    background: #fff;
}
.divider {
    width: 2px;
    background: linear-gradient(180deg, var(--color-accent) 0%, var(--color-accent) 100%);
    margin: 0 0.5rem;
    border-radius: 2px;
}
.payment-details-panel {
    flex: 1 1 0;
    max-width: 900px;
    width: 100%;
    padding: 3rem 2.5rem 2.5rem 2.5rem;
    background: var(--color-card);
    color: var(--color-text);
    border-radius: 0 18px 18px 0;
    box-shadow: 0 2px 12px rgba(0,0,0,0.08);
    min-width: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: flex-start;
    margin: 0 auto;
}
.empty-details {
    color: var(--color-text);
    font-size: 1.2rem;
    margin-top: 4rem;
    text-align: center;
}
@media (max-width: 900px) {
    .payments-split-view {
        flex-direction: column;
    }
    .patient-list-panel {
        border-radius: 18px 18px 0 0;
        min-height: unset;
        padding: 1.2rem;
    }
    .divider {
        width: 100%;
        height: 2px;
        margin: 0.5rem 0;
    }
    .payment-details-panel {
        padding: 1.2rem;
    }
}
</style> 