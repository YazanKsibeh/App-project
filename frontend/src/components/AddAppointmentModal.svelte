<script>
import { createEventDispatcher } from 'svelte';
import { patients } from '../stores/patientStore.js';
import { addAppointment } from '../stores/appointmentStore.js';
import { onMount } from 'svelte';
import Flatpickr from 'svelte-flatpickr';
import 'flatpickr/dist/flatpickr.css';

const dispatch = createEventDispatcher();

let patientId = '';
let datetime = '';
let duration = 30;
let notes = '';
let error = '';
let patientSearch = '';
let showPatientDropdown = false;
let filteredPatients = [];

$: filteredPatients = $patients.filter(p =>
    p.name.toLowerCase().includes(patientSearch.toLowerCase())
);

function selectPatient(p) {
    patientId = p.id;
    patientSearch = p.name;
    showPatientDropdown = false;
}

let flatpickrOptions = {
    enableTime: true,
    dateFormat: 'Y-m-d h:i K',
    time_24hr: false,
    minuteIncrement: 5,
};

onMount(() => {
    // No need to manually subscribe to patients, Svelte will handle $patients
    // Set default datetime to today at 14:00 (2:00 PM) if not already set
    if (!datetime) {
        const now = new Date();
        now.setHours(14, 0, 0, 0);
        // Format as yyyy-MM-ddTHH:mm for datetime-local input
        const yyyy = now.getFullYear();
        const mm = String(now.getMonth() + 1).padStart(2, '0');
        const dd = String(now.getDate()).padStart(2, '0');
        const hh = String(now.getHours()).padStart(2, '0');
        const min = String(now.getMinutes()).padStart(2, '0');
        datetime = `${yyyy}-${mm}-${dd}T${hh}:${min}`;
    }
});

async function handleSubmit() {
    error = '';
    if (!patientId || !datetime) {
        error = 'Patient and date/time are required.';
        return;
    }
    try {
        await addAppointment({
            patient_id: parseInt(patientId),
            datetime,
            duration: parseInt(duration),
            notes
        });
        dispatch('close');
    } catch (err) {
        error = err.message || 'Failed to add appointment';
    }
}
</script>

<div class="modal-backdrop" on:click={() => dispatch('close')}></div>
<div class="modal">
    <h3>Add Appointment</h3>
    {#if error}
        <p class="error">{error}</p>
    {/if}
    <form on:submit|preventDefault={handleSubmit} autocomplete="off">
        <label>Patient</label>
        <div class="patient-search-wrapper">
            <input type="text" placeholder="Type patient name..." bind:value={patientSearch} on:input={() => { showPatientDropdown = true; patientId = ''; }} on:focus={() => showPatientDropdown = true} autocomplete="off" required />
            {#if showPatientDropdown && patientSearch.trim()}
                <ul class="patient-dropdown">
                    {#each filteredPatients as p}
                        <li on:click={() => selectPatient(p)}>{p.name}</li>
                    {/each}
                    {#if filteredPatients.length === 0}
                        <li class="no-match">No matches</li>
                    {/if}
                </ul>
            {/if}
        </div>
        <label>Date & Time</label>
        <Flatpickr
            bind:value={datetime}
            options={flatpickrOptions}
            placeholder="Select date & time"
            required
            class="modal-input"
        />
        <label>Duration (minutes)</label>
        <input type="number" min="1" bind:value={duration} />
        <label>Notes</label>
        <textarea bind:value={notes} placeholder="Optional"></textarea>
        <div class="actions">
            <button type="submit">Add</button>
            <button type="button" on:click={() => dispatch('close')}>Cancel</button>
        </div>
    </form>
</div>

<style>
.modal-backdrop {
    position: fixed;
    top: 0; left: 0; right: 0; bottom: 0;
    background: rgba(0,0,0,0.2);
    z-index: 1000;
}
.modal {
    position: fixed;
    top: 50%; left: 50%;
    transform: translate(-50%, -50%);
    background: var(--color-card);
    color: var(--color-text);
    padding: 2rem;
    border-radius: 12px;
    box-shadow: 0 8px 32px rgba(0,0,0,0.18);
    z-index: 1001;
    min-width: 320px;
    max-width: 90vw;
}
h3 {
    color: var(--color-text);
    font-size: 1.3rem;
    font-weight: 700;
    margin-bottom: 1.2rem;
}
form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}
label {
    font-weight: 500;
    color: var(--color-text);
}
input, select, textarea, .modal-input {
    padding: 0.5rem;
    border-radius: 6px;
    border: 1px solid var(--color-border);
    font-size: 1rem;
    color: var(--color-text);
    background: var(--color-panel, #fff);
    width: 100%;
    box-sizing: border-box;
}
.actions {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
}
button[type="submit"] {
    background: var(--color-accent);
    color: #fff;
    border: none;
    border-radius: 6px;
    padding: 0.5rem 1.2rem;
    font-size: 1rem;
    cursor: pointer;
    transition: background 0.2s;
}
button[type="submit"]:hover {
    background: #5a67d8;
}
button[type="button"] {
    background: var(--color-panel, #eee);
    color: var(--color-text);
    border: none;
    border-radius: 6px;
    padding: 0.5rem 1.2rem;
    font-size: 1rem;
    cursor: pointer;
}
.error {
    color: var(--color-danger);
    font-weight: bold;
}
.patient-search-wrapper {
    position: relative;
}
.patient-search-wrapper input[type="text"] {
    width: 100%;
}
.patient-dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: var(--color-panel, #fff);
    border: 1px solid var(--color-border);
    border-radius: 0 0 8px 8px;
    max-height: 180px;
    overflow-y: auto;
    z-index: 10;
    margin: 0;
    padding: 0;
    list-style: none;
    box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}
.patient-dropdown li {
    padding: 0.5rem 1rem;
    cursor: pointer;
    color: var(--color-text);
}
.patient-dropdown li:hover {
    background: #f3f4f6;
}
.patient-dropdown .no-match {
    color: #888;
    cursor: default;
}
</style> 