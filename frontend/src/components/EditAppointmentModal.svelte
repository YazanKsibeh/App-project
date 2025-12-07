<script>
import { createEventDispatcher, onMount } from 'svelte';
import { patients } from '../stores/patientStore.js';
import { updateAppointment, deleteAppointment, loadAppointments } from '../stores/appointmentStore.js';

export let appointment = null;
const dispatch = createEventDispatcher();

let patientId = '';
let datetime = '';
let duration = 30;
let notes = '';
let error = '';
let showDeleteConfirm = false;

onMount(() => {
    if (appointment) {
        patientId = appointment.patient_id;
        datetime = appointment.datetime;
        duration = appointment.duration;
        notes = appointment.notes;
    }
});

async function handleSave() {
    error = '';
    if (!patientId || !datetime) {
        error = 'Patient and date/time are required.';
        return;
    }
    try {
        await updateAppointment({
            id: appointment.id,
            patient_id: parseInt(patientId),
            datetime,
            duration: parseInt(duration),
            notes
        });
        await loadAppointments();
        dispatch('close');
    } catch (err) {
        error = err.message || 'Failed to update appointment';
    }
}

function handleDelete() {
    showDeleteConfirm = true;
}

async function confirmDelete() {
    await deleteAppointment(appointment.id);
    await loadAppointments();
    showDeleteConfirm = false;
    dispatch('close');
}

function cancelDelete() {
    showDeleteConfirm = false;
}

function handleCancel() {
    dispatch('close');
}
</script>

{#if showDeleteConfirm}
    <div class="modal-backdrop"></div>
    <div class="modal confirm-modal">
        <p>Confirm deleting appointment</p>
        <div class="actions">
            <button on:click={confirmDelete}>Yes</button>
            <button on:click={cancelDelete}>Cancel</button>
        </div>
    </div>
{:else}
    <div class="modal-backdrop" on:click={handleCancel}></div>
    <div class="modal">
        <h3>Edit Appointment</h3>
        {#if error}
            <p class="error">{error}</p>
        {/if}
        <form on:submit|preventDefault={handleSave}>
            <label>Patient</label>
            <select bind:value={patientId} required>
                <option value="" disabled>Select patient</option>
                {#each $patients as p}
                    <option value={p.id}>{p.name}</option>
                {/each}
            </select>
            <label>Date & Time</label>
            <input type="datetime-local" bind:value={datetime} required />
            <label>Duration (minutes)</label>
            <input type="number" min="1" bind:value={duration} />
            <label>Notes</label>
            <textarea bind:value={notes} placeholder="Optional"></textarea>
            <div class="actions">
                <button type="submit">Save</button>
                <button type="button" on:click={handleDelete}>Delete</button>
                <button type="button" on:click={handleCancel}>Cancel</button>
            </div>
        </form>
    </div>
{/if}

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
    background: white;
    padding: 2rem;
    border-radius: 12px;
    box-shadow: 0 8px 32px rgba(0,0,0,0.18);
    z-index: 1001;
    min-width: 320px;
    max-width: 90vw;
}
form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}
label {
    font-weight: 500;
}
input, select, textarea {
    padding: 0.5rem;
    border-radius: 6px;
    border: 1px solid #ddd;
    font-size: 1rem;
}
.actions {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
}
button[type="submit"], .actions button {
    background: #667eea;
    color: white;
    border: none;
    border-radius: 6px;
    padding: 0.5rem 1.2rem;
    font-size: 1rem;
    cursor: pointer;
    transition: background 0.2s;
}
button[type="submit"]:hover, .actions button:hover {
    background: #5a67d8;
}
button[type="button"] {
    background: #eee;
    color: #333;
    border: none;
    border-radius: 6px;
    padding: 0.5rem 1.2rem;
    font-size: 1rem;
    cursor: pointer;
}
.error {
    color: #e74c3c;
    font-weight: bold;
}
.confirm-modal {
    min-width: 250px;
    text-align: center;
}
</style> 