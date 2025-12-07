<script>
    import { onMount } from 'svelte';
    import { 
        patients, 
        loading, 
        error, 
        searchPatients, 
        loadPatients, 
        deletePatient,
        selectedPatient
    } from '../stores/patientStore.js';
    import PatientCard from './PatientCard.svelte';
    import SearchBar from './SearchBar.svelte';
    import AddPatientModal from './AddPatientModal.svelte';
    import CompanyFooter from './CompanyFooter.svelte';

    let searchTerm = '';
    let showAddModal = false;
    let patientToEdit = null;
    let currentPage = 1;
    const pageSize = 12;

    onMount(() => {
        loadPatients();
    });

    function handleSearch(event) {
        searchTerm = event.detail;
        if (searchTerm.trim()) {
            searchPatients(searchTerm);
        } else {
            loadPatients();
        }
    }

    function handleDelete(id) {
        if (confirm('Are you sure you want to delete this patient?')) {
            deletePatient(id);
        }
    }

    function handleEdit(event) {
        patientToEdit = event.detail;
        openAddModal();
    }

    function handleView(event) {
        selectedPatient.set(event.detail);
    }

    function openAddModal() {
        showAddModal = true;
    }

    function closeAddModal() {
        showAddModal = false;
        patientToEdit = null;
    }

    // Format phone number as 0941-414-122
    function formatPhone(phone) {
        // Remove non-digits
        const digits = phone.replace(/\D/g, '');
        if (digits.length === 10) {
            return `${digits.slice(0,4)}-${digits.slice(4,7)}-${digits.slice(7)}`;
        }
        return phone;
    }

    // Sort patients by most recent (highest ID first)
    $: sortedPatients = [...$patients].sort((a, b) => b.id - a.id);

    // Paginate
    $: totalPages = Math.ceil(sortedPatients.length / pageSize);
    $: paginatedPatients = sortedPatients.slice((currentPage - 1) * pageSize, currentPage * pageSize);

    function goToPage(page) {
        if (page >= 1 && page <= totalPages) {
            currentPage = page;
        }
    }
</script>

<main class="patient-list">
    <div class="top-bar">
        <SearchBar on:search={handleSearch} />
        <button class="add-patient-btn prominent" on:click={openAddModal}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="5" x2="12" y2="19"/>
                <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
            Add Patient
        </button>
    </div>

    {#if $error}
        <div class="error-message">
            <p>⚠️ {$error}</p>
            <button on:click={() => loadPatients()}>Retry</button>
        </div>
    {/if}

    {#if $loading}
        <div class="loading">
            <div class="spinner"></div>
            <p>Loading patients...</p>
        </div>
    {:else if $patients.length === 0}
        <div class="empty-state">
            <div class="empty-icon">👥</div>
            <h3>No patients found</h3>
            <p>{searchTerm ? 'Try adjusting your search terms' : 'Add your first patient to get started'}</p>
            {#if !searchTerm}
                <button class="add-first-patient-btn" on:click={openAddModal}>
                    Add Your First Patient
                </button>
            {/if}
        </div>
    {:else}
        <div class="patients-grid">
            {#each paginatedPatients as patient (patient.id)}
                <PatientCard 
                    {patient} 
                    on:delete={() => handleDelete(patient.id)}
                    on:edit={handleEdit}
                    on:view={handleView}
                />
            {/each}
        </div>
        <div class="pagination-bar">
            <button class="pagination-btn" on:click={() => goToPage(currentPage - 1)} disabled={currentPage === 1}>Previous</button>
            <span class="pagination-info">Page {currentPage} of {totalPages}</span>
            <button class="pagination-btn" on:click={() => goToPage(currentPage + 1)} disabled={currentPage === totalPages}>Next</button>
        </div>
    {/if}
</main>

<!-- Company Footer -->
<CompanyFooter />

{#if showAddModal}
    <AddPatientModal {patientToEdit} on:close={closeAddModal} />
{/if}

<style>
    .patient-list {
        padding: 2rem;
        width: 100%;
        max-width: none;
        margin: 0;
        min-height: 100vh;
        background: var(--color-bg-gradient, var(--color-bg));
        position: relative;
        box-sizing: border-box;
        overflow-x: hidden;
    }

    .top-bar {
        display: flex;
        align-items: center;
        gap: 1.2rem;
        margin-bottom: 2.2rem;
        justify-content: flex-start;
    }

    .add-patient-btn {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 1.1rem;
        font-weight: 600;
        border-radius: 25px;
        padding: 0.85rem 2rem;
        cursor: pointer;
        border: none;
        transition: background 0.2s, box-shadow 0.2s;
        white-space: nowrap;
        min-width: 0;
        width: auto;
    }

    .add-patient-btn.prominent {
        background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
        color: #fff;
        box-shadow: 0 2px 8px rgba(102,126,234,0.10);
        border: none;
    }

    .add-patient-btn.prominent:hover {
        background: linear-gradient(90deg, #5a67d8 0%, #6b47b6 100%);
        box-shadow: 0 4px 16px rgba(102,126,234,0.18);
    }

    .patients-grid {
        display: grid;
        grid-template-columns: repeat(4, 1fr);
        gap: 1.2rem;
        justify-items: center;
        align-items: center;
    }

    .error-message {
        background: #fee;
        border: 1px solid #fcc;
        border-radius: 8px;
        padding: 1rem;
        margin-bottom: 2rem;
        color: #c33;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .error-message button {
        background: #c33;
        color: white;
        border: none;
        padding: 0.5rem 1rem;
        border-radius: 4px;
        cursor: pointer;
    }

    .loading {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 4rem;
        color: white;
    }

    .spinner {
        width: 40px;
        height: 40px;
        border: 4px solid rgba(255,255,255,0.3);
        border-top: 4px solid white;
        border-radius: 50%;
        animation: spin 1s linear infinite;
        margin-bottom: 1rem;
    }

    @keyframes spin {
        0% { transform: rotate(0deg); }
        100% { transform: rotate(360deg); }
    }

    .empty-state {
        text-align: center;
        padding: 4rem;
        color: white;
    }

    .empty-icon {
        font-size: 4rem;
        margin-bottom: 1rem;
    }

    .empty-state h3 {
        font-size: 1.5rem;
        margin: 0 0 0.5rem 0;
    }

    .empty-state p {
        font-size: 1.1rem;
        opacity: 0.9;
        margin: 0 0 2rem 0;
    }

    .add-first-patient-btn {
        background: rgba(255, 255, 255, 0.2);
        border: 2px solid rgba(255, 255, 255, 0.3);
        color: white;
        padding: 1rem 2rem;
        border-radius: 25px;
        font-size: 1.1rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.3s ease;
        backdrop-filter: blur(10px);
    }

    .add-first-patient-btn:hover {
        background: rgba(255, 255, 255, 0.3);
        border-color: rgba(255, 255, 255, 0.5);
        transform: translateY(-2px);
        box-shadow: 0 8px 25px rgba(0, 0, 0, 0.2);
    }

    :global(.patient-card) {
        width: 260px;
        height: 260px;
        min-width: 260px;
        min-height: 260px;
        max-width: 260px;
        max-height: 260px;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
    }

    .pagination-bar {
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 1.5rem;
        margin-top: 2rem;
    }

    .pagination-btn {
        background: #667eea;
        color: white;
        border: none;
        border-radius: 6px;
        padding: 0.5rem 1.2rem;
        font-size: 1rem;
        font-weight: 600;
        cursor: pointer;
        transition: background 0.2s;
    }

    .pagination-btn:disabled {
        background: #bfc7e6;
        cursor: not-allowed;
    }

    .pagination-btn:hover:not(:disabled) {
        background: #5a67d8;
    }

    .pagination-info {
        color: #fff;
        font-size: 1.1rem;
        font-weight: 500;
    }
</style> 