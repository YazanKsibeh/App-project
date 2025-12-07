<script>
import { onMount } from 'svelte';
import { appointments, loadAppointments, loadingAppointments, appointmentError, deleteAppointment } from '../stores/appointmentStore.js';
import AddAppointmentModal from './AddAppointmentModal.svelte';
import { patients } from '../stores/patientStore.js';
import EditAppointmentModal from './EditAppointmentModal.svelte';
import { derived } from 'svelte/store';
import { DateInput } from 'date-picker-svelte';

let showAddModal = false;
let showEditModal = false;
let selectedAppointment = null;
let confirmDeleteId = null;

// Filter state
let filterDate = '';
let filterSearch = '';
let showDeleteAllConfirm = false;

let currentPage = 1;
let pageSize = 10;

onMount(() => {
    loadAppointments();
});

// Use Svelte's reactive stores
$: $patients;
$: $appointments;

// Build patientMap reactively
$: patientMap = $patients.reduce((acc, p) => { acc[String(p.id)] = p; return acc; }, {});

function getPatientName(patientId) {
    const patient = patientMap[String(patientId)];
    return patient ? patient.name : 'Unknown';
}

function getPatientPhone(patientId) {
    const patient = patientMap[String(patientId)];
    return patient ? patient.phone : '';
}

function handleEdit(appt) {
    selectedAppointment = appt;
    showEditModal = true;
}

function handleDelete(id) {
    confirmDeleteId = id;
}

function confirmDelete() {
    if (confirmDeleteId !== null) {
        deleteAppointment(confirmDeleteId).then(() => {
            loadAppointments();
            confirmDeleteId = null;
        });
    }
}

function cancelDelete() {
    confirmDeleteId = null;
}

function closeEditModal() {
    showEditModal = false;
    selectedAppointment = null;
}

function clearFilters() {
    filterDate = '';
    filterSearch = '';
}

let filterDateObj = filterDate ? new Date(filterDate) : null;
$: filterDate = filterDateObj ? filterDateObj.toISOString().slice(0, 10) : '';

// Combined filter logic with sorting
$: filteredAppointments = $appointments
    .filter(appt => {
        let matchesDate = true;
        let matchesSearch = true;
        if (filterDate) {
            const apptDate = new Date(appt.datetime.length === 16 ? appt.datetime + ':00' : appt.datetime).toISOString().slice(0, 10);
            matchesDate = apptDate === filterDate;
        }
        if (filterSearch.trim()) {
            const patient = patientMap[String(appt.patient_id)] || {};
            const name = (patient.name || '').toLowerCase();
            const phone = (patient.phone || '').replace(/\D/g, '');
            const search = filterSearch.toLowerCase();
            let matchesName = name.includes(search);
            let matchesPhone = false;
            if (/\d/.test(search)) { // Only check phone if search has digits
                matchesPhone = phone.includes(search.replace(/\D/g, ''));
            }
            matchesSearch = matchesName || matchesPhone;
            console.log(
              `Checking: name="${name}", phone="${phone}", search="${search}", matchesName=${matchesName}, matchesPhone=${matchesPhone}, matchesSearch=${matchesSearch}`
            );
        }
        return matchesDate && matchesSearch;
    })
    .sort((a, b) => {
        const dateA = new Date(a.datetime.length === 16 ? a.datetime + ':00' : a.datetime).getTime();
        const dateB = new Date(b.datetime.length === 16 ? b.datetime + ':00' : b.datetime).getTime();
        return dateB - dateA;
    });

$: totalPages = Math.ceil(filteredAppointments.length / pageSize);
$: paginatedAppointments = filteredAppointments.slice(
  (currentPage - 1) * pageSize,
  currentPage * pageSize
);

// Reset to page 1 when filters change
$: if (filterSearch || filterDate) currentPage = 1;

// Delete all appointments for selected day
async function deleteAllForDay() {
    showDeleteAllConfirm = false;
    const toDelete = filteredAppointments.map(a => a.id);
    for (const id of toDelete) {
        await deleteAppointment(id);
    }
    await loadAppointments();
}
</script>

<div class="appointments-layout">
    <aside class="filter-sidebar">
        <h2>Filters</h2>
        <div class="filter-group">
            <label for="search-filter">Search</label>
            <div class="search-input-wrapper">
                <svg class="search-icon" viewBox="0 0 24 24" width="18" height="18"><circle cx="11" cy="11" r="8" stroke="currentColor" stroke-width="2" fill="none"/><line x1="21" y1="21" x2="16.65" y2="16.65" stroke="currentColor" stroke-width="2"/></svg>
                <input id="search-filter" type="text" placeholder="Name or phone..." bind:value={filterSearch} />
            </div>
        </div>
        <div class="filter-group">
            <label for="date-filter">Date</label>
            <DateInput
                id="date-filter"
                bind:value={filterDateObj}
                placeholder="Select date"
                format="yyyy-MM-dd"
                on:change={() => filterDate = filterDateObj ? filterDateObj.toISOString().slice(0, 10) : ''}
            />
        </div>
        <button class="clear-filters-btn" on:click={clearFilters} disabled={!filterDate && !filterSearch}>Clear Filters</button>
    </aside>
    <main class="appointments-main">
        <div class="appointments-header-row">
            <h1 class="appointments-header">Appointments</h1>
            <button class="add-btn" on:click={() => showAddModal = true}>+ Add Appointment</button>
        </div>
        {#if filterDate && filteredAppointments.length > 0}
            <button class="delete-all-btn" on:click={() => showDeleteAllConfirm = true}>Delete All for This Day</button>
        {/if}
        {#if showDeleteAllConfirm}
            <div class="modal-backdrop"></div>
            <div class="modal confirm-modal">
                <p>Delete all appointments for {filterDate}?</p>
                <div class="actions">
                    <button on:click={deleteAllForDay}>Yes</button>
                    <button on:click={() => showDeleteAllConfirm = false}>Cancel</button>
                </div>
            </div>
        {/if}
        {#if $loadingAppointments}
            <p>Loading appointments...</p>
        {:else if $appointmentError}
            <p class="error">{$appointmentError}</p>
        {:else}
            <table class="appointment-table">
                <thead>
                    <tr>
                        <th>Patient</th>
                        <th>Date & Time</th>
                        <th>Duration</th>
                        <th>Notes</th>
                        <th style="min-width: 140px;"></th>
                    </tr>
                </thead>
                <tbody>
                    {#each paginatedAppointments as appt}
                        <tr>
                            <td>{getPatientName(appt.patient_id)}</td>
                            <td>{new Date(appt.datetime.length === 16 ? appt.datetime + ':00' : appt.datetime).toLocaleString()}</td>
                            <td>{appt.duration} min</td>
                            <td>{appt.notes}</td>
                            <td>
                                <div class="actions-inline">
                                    <button class="edit-btn" on:click={() => handleEdit(appt)}>Edit</button>
                                    <button class="delete-btn" on:click={() => handleDelete(appt.id)}>Delete</button>
                                </div>
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
            {#if totalPages > 1}
            <div class="pagination">
                <button on:click={() => currentPage = Math.max(1, currentPage - 1)} disabled={currentPage === 1}>Prev</button>
                <span>Page {currentPage} of {totalPages}</span>
                <button on:click={() => currentPage = Math.min(totalPages, currentPage + 1)} disabled={currentPage === totalPages}>Next</button>
            </div>
            {/if}
        {/if}
        {#if showAddModal}
            <AddAppointmentModal on:close={() => showAddModal = false} />
        {/if}
        {#if showEditModal}
            <EditAppointmentModal appointment={selectedAppointment} on:close={closeEditModal} />
        {/if}
        {#if confirmDeleteId !== null}
            <div class="modal-backdrop"></div>
            <div class="modal confirm-modal">
                <p>Confirm deleting appointment</p>
                <div class="actions">
                    <button on:click={confirmDelete}>Yes</button>
                    <button on:click={cancelDelete}>Cancel</button>
                </div>
            </div>
        {/if}
    </main>
</div>

<style>
.appointments-layout {
    display: flex;
    flex-direction: row;
    gap: 2.5rem;
    width: 100%;
    max-width: 1600px;
    align-items: flex-start;
}

.filter-sidebar {
    background: #fff;
    border-radius: 18px;
    padding: 2rem;
    min-width: 260px;
    max-height: 100%;
    box-shadow: 0 2px 8px rgba(0,0,0,0.04);
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
}

.appointments-main {
    background: #fff;
    border-radius: 18px;
    padding: 2rem;
    margin-right: 2rem;
    box-shadow: 0 2px 8px rgba(0,0,0,0.04);
    flex: 1;
}

.appointment-table {
    border-radius: 12px;
    overflow: hidden;
    width: 100%;
    box-shadow: 0 1px 4px rgba(0,0,0,0.03);
}

.filter-sidebar h2 {
    font-size: 1.2rem;
    font-weight: 700;
    margin-bottom: 0.5rem;
}
.filter-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
}
.search-input-wrapper {
    display: flex;
    align-items: center;
    background: var(--color-panel);
    border-radius: 8px;
    padding: 0.2rem 0.7rem;
    border: 1px solid var(--color-border);
}
.search-input-wrapper input {
    border: none;
    background: transparent;
    color: var(--color-text);
    font-size: 1rem;
    padding: 0.6rem 0.5rem;
    outline: none;
    width: 100%;
}
.search-icon {
    color: var(--color-text);
    margin-right: 0.3rem;
}
.clear-filters-btn {
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 0.6rem 1.2rem;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    margin-top: 1.2rem;
    transition: background 0.18s, color 0.18s;
}
.clear-filters-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}
.clear-filters-btn:hover:not(:disabled) {
    background: var(--color-accent);
    color: #fff;
}
.appointments-header-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1.2rem;
}
.appointments-header {
    color: var(--color-text, #222);
    font-size: 2.2rem;
    font-weight: 700;
    margin-bottom: 1.5rem;
}
.add-btn {
    background: var(--color-accent);
    color: #fff;
    border: none;
    border-radius: 8px;
    padding: 0.7rem 1.5rem;
    font-size: 1.1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.18s;
}
.add-btn:hover {
    background: #5a67d8;
}
.delete-all-btn {
    width: auto;
    display: inline-block;
    margin: 1rem 0 1.5rem 0;
    padding: 0.5rem 1.2rem;
    font-size: 1rem;
    background: #ff4d4f;
    color: #fff;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    transition: background 0.2s;
}
.delete-all-btn:hover {
    background: #d9363e;
}
.appointment-table th, .appointment-table td {
    padding: 0.75rem 1rem;
    border-bottom: 1px solid #f0f0f0;
    text-align: left;
    color: #222;
}
.appointment-table th {
    background: #f8f8ff;
    font-weight: 600;
    color: #222;
}
.error {
    color: #e74c3c;
    font-weight: bold;
}
.modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.5);
}
.modal {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: white;
    padding: 2rem;
    border-radius: 12px;
    box-shadow: 0 4px 12px rgba(0,0,0,0.1);
    max-width: 400px;
    width: 100%;
}
.modal p {
    margin-bottom: 1.5rem;
}
.actions {
    display: flex;
    justify-content: space-between;
}
.actions button {
    background: #667eea;
    color: white;
    border: none;
    border-radius: 6px;
    padding: 0.5rem 1.2rem;
    font-size: 1rem;
    cursor: pointer;
    transition: background 0.2s;
}
.actions button:hover {
    background: #5a67d8;
}
.edit-btn {
    background: #667eea;
    color: white;
    border: none;
    border-radius: 6px;
    padding: 0.3rem 0.8rem;
    margin-right: 0.5rem;
    font-size: 0.95rem;
    cursor: pointer;
    transition: background 0.2s;
}
.edit-btn:hover {
    background: #5a67d8;
}
.delete-btn {
    background: #e74c3c;
    color: white;
    border: none;
    border-radius: 6px;
    padding: 0.3rem 0.8rem;
    font-size: 0.95rem;
    cursor: pointer;
    transition: background 0.2s;
}
.delete-btn:hover {
    background: #c0392b;
}
.actions-inline {
    display: flex;
    gap: 0.5rem;
    align-items: center;
}
.confirm-modal {
    min-width: 250px;
    text-align: center;
    background: #222;
    color: #fff;
    font-weight: bold;
    padding: 2rem 1rem;
    border-radius: 12px;
}
.confirm-modal p {
    margin-bottom: 1.5rem;
    color: #fff;
    font-size: 1.1rem;
    font-weight: bold;
}
.actions {
    display: flex;
    justify-content: center;
    gap: 1.5rem;
}
.actions button {
    background: #667eea;
    color: white;
    border: none;
    border-radius: 6px;
    padding: 0.5rem 1.2rem;
    font-size: 1rem;
    cursor: pointer;
    transition: background 0.2s;
}
.actions button:hover {
    background: #5a67d8;
}
.pagination {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 1rem;
    margin-top: 1.5rem;
}
.pagination button {
    background: #e5e7eb;
    border: none;
    border-radius: 6px;
    padding: 0.5rem 1.2rem;
    font-size: 1rem;
    cursor: pointer;
    transition: background 0.2s;
}
.pagination button:disabled {
    background: #cbd5e1;
    cursor: not-allowed;
}

:root {
    --color-text: #222;
}

[data-theme="dark"] {
    --color-text: #fff;
}

/* Filter Sidebar Headings */
.filter-sidebar h2,
.filter-sidebar label {
    color: #222; /* For light theme */
}

:global(body.dark) .filter-sidebar h2,
:global(body.dark) .filter-sidebar label {
    color: #fff; /* For dark theme */
}

/* Appointments List Heading */
.appointments-header {
    color: #222; /* For light theme */
}

:global(body.dark) .appointments-header {
    color: #fff; /* For dark theme */
}
</style> 