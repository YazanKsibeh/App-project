<script>
    import { createEventDispatcher } from 'svelte';
    import { openPatientFolder } from '../stores/patientStore.js';
    
    const dispatch = createEventDispatcher();
    
    export let patient;
    
    function handleEdit() {
        dispatch('edit', patient);
    }
    
    function handleDelete() {
        dispatch('delete');
    }
    
    function handleView() {
        dispatch('view', patient);
    }
    
    async function handleOpenFolder() {
        console.log('[PATIENT-CARD] Open folder button clicked for patient:', patient.id);
        try {
            await openPatientFolder(patient.id);
            console.log('[PATIENT-CARD] Open folder completed successfully for patient:', patient.id);
        } catch (error) {
            console.error('[PATIENT-CARD] Open folder failed for patient:', patient.id, 'Error:', error);
        }
    }
    
    function getGenderIcon(gender) {
        return gender.toLowerCase() === 'male' ? '👨' : gender.toLowerCase() === 'female' ? '👩' : '👤';
    }
    
    function formatPhone(phone) {
        return phone.replace(/(\d{3})(\d{3})(\d{4})/, '($1) $2-$3');
    }
</script>

<div class="patient-card">
    <div class="card-header">
        <div class="patient-info">
            <h3 class="patient-name">{patient.name}</h3>
            <p class="patient-id">ID: {patient.id}</p>
        </div>
    </div>
    
    <div class="card-content">
        <div class="info-grid">
            <div class="info-row">
                <span class="label">📞 Phone:</span>
                <span class="value">{formatPhone(patient.phone)}</span>
            </div>
            <div class="info-row">
                <span class="label">🎂 Age:</span>
                <span class="value">{patient.age} years</span>
            </div>
            <div class="info-row">
                <span class="label">⚧ Gender:</span>
                <span class="value">{patient.gender}</span>
            </div>
        </div>
    </div>
    
    <div class="card-actions">
        <button on:click={handleView} class="action-btn view-btn" title="View Details">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                <circle cx="12" cy="12" r="3"/>
            </svg>
        </button>
        <button on:click={handleEdit} class="action-btn edit-btn" title="Edit Patient">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                <path d="m18.5 2.5 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
            </svg>
        </button>
        <button on:click={handleDelete} class="action-btn delete-btn" title="Delete Patient">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3,6 5,6 21,6"/>
                <path d="m19,6v14a2,2 0 0,1 -2,2H7a2,2 0 0,1 -2,-2V6m3,0V4a2,2 0 0,1 2,-2h4a2,2 0 0,1 2,2v2"/>
            </svg>
        </button>
        <button on:click={handleOpenFolder} class="action-btn folder-btn" title="Open Patient Folder">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M3 7a2 2 0 0 1 2-2h4l2 3h7a2 2 0 0 1 2 2v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
            </svg>
        </button>
    </div>
</div>

<style>
    .patient-card {
        background: var(--color-card);
        border-radius: 12px;
        padding: 1.5rem;
        box-shadow: var(--color-shadow);
        transition: all 0.3s ease;
        border: 1px solid var(--color-border);
        color: var(--color-text);
        display: flex;
        flex-direction: column;
        min-height: 280px;
        height: 100%;
    }
    
    .patient-card:hover {
        transform: translateY(-4px);
        box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
    }
    
    .card-header {
        margin-bottom: 1rem;
    }
    
    .patient-name {
        font-size: 1.15rem;
        font-weight: 700;
        color: var(--color-text);
        word-break: break-word;
        margin-bottom: 0.25rem;
    }
    
    .patient-id {
        font-size: 0.875rem;
        color: var(--color-text);
        margin: 0;
    }
    
    .card-content {
        flex: 1;
    }
    
    .info-grid {
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
    }
    
    .info-row {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0.25rem 0;
    }
    
    .card-actions {
        display: flex;
        gap: 0.5rem;
        justify-content: flex-end;
        padding-top: 1rem;
        margin-top: auto;
        border-top: 1px solid var(--color-border);
    }
    
    .action-btn {
        background: none;
        border: none;
        padding: 0.5rem;
        border-radius: 6px;
        cursor: pointer;
        transition: all 0.2s ease;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    
    .action-btn svg {
        width: 18px;
        height: 18px;
    }
    
    /* Button color styles remain the same */
    .view-btn { color: var(--color-accent); }
    .view-btn:hover { background: rgba(102, 126, 234, 0.1); }
    .edit-btn { color: #f39c12; }
    .edit-btn:hover { background: rgba(243, 156, 18, 0.1); }
    .delete-btn { color: var(--color-danger); }
    .delete-btn:hover { background: rgba(231, 76, 60, 0.1); }
    .folder-btn { color: #764ba2; }
    .folder-btn:hover { background: rgba(118, 75, 162, 0.1); }
    
    .label {
        font-weight: 500;
        color: var(--color-text);
        font-size: 0.9rem;
    }
    
    .value {
        color: var(--color-text);
        font-weight: 500;
        text-align: right;
        max-width: 60%;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }
    
    @media (max-width: 768px) {
        .patient-card {
            padding: 1rem;
            min-height: auto;
        }
        
        .info-row {
            flex-direction: column;
            align-items: flex-start;
            gap: 0.25rem;
        }
        
        .value {
            text-align: left;
            max-width: 100%;
        }
    }
</style>