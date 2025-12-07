<script>
import { theme, account, licenseValid, setLicense, licenseValidationStatus, validateCurrentLicense } from '../stores/settingsStore.js';
import { deleteAllPatients } from '../stores/patientStore.js';
import { currentUser, isAdmin, logout } from '../stores/authStore.js';
import UserManagement from './UserManagement.svelte';
import { onMount } from 'svelte';

export let open = false;
export let onClose = () => {};

let showLicenseInput = false;
let newKey = '';
let validatingLicense = false;
let showDeleteAllConfirm = false;
let deletingAllPatients = false;
let showUserManagement = false;

onMount(() => {
    document.body.setAttribute('data-theme', $theme);
});

$: document.body.setAttribute('data-theme', $theme);

function switchTheme() {
    theme.update(t => t === 'dark' ? 'light' : 'dark');
}

async function handleLicenseSave() {
    if (!newKey.trim()) {
        return;
    }
    
    validatingLicense = true;
    try {
        const isValid = await setLicense(newKey);
        if (isValid) {
            showLicenseInput = false;
            newKey = '';
        }
    } catch (error) {
        console.error('License save error:', error);
    } finally {
        validatingLicense = false;
    }
}

async function refreshLicense() {
    await validateCurrentLicense();
}

async function handleDeleteAllPatients() {
    deletingAllPatients = true;
    try {
        await deleteAllPatients();
        showDeleteAllConfirm = false;
    } catch (error) {
        console.error('Error deleting all patients:', error);
        alert('Failed to delete all patients. Please try again.');
    } finally {
        deletingAllPatients = false;
    }
}

function handleLogout() {
    if (confirm('Are you sure you want to logout?')) {
        logout();
        onClose();
        // The App component will handle showing the login screen
        window.location.reload(); // Simple way to reset the app state
    }
}
</script>

{#if open}
    <div class="sidebar-backdrop" on:click={onClose}></div>
    <aside class="settings-sidebar">
        <button class="close-btn" on:click={onClose} title="Close">&times;</button>
        <h2>Settings</h2>
        <div class="section">
            <h3>Color Theme</h3>
            <button class="theme-toggle" on:click={switchTheme}>
                {#if $theme === 'dark'}🌙 Dark{:else}☀️ Light{/if}
            </button>
        </div>
        <div class="section">
            <h3>User Account</h3>
            <div class="account-info">
                {#if $currentUser}
                    <div><b>Username:</b> {$currentUser.username}</div>
                    <div><b>Role:</b> {$currentUser.role}</div>
                {/if}
            </div>
            <button class="logout-btn" on:click={handleLogout}>
                🚪 Logout
            </button>
        </div>
        
        {#if isAdmin()}
        <div class="section">
            <h3>User Management</h3>
            <button class="user-mgmt-btn" on:click={() => showUserManagement = true}>
                👥 Manage Users
            </button>
        </div>
        {/if}
        
        <div class="section">
            <h3>License & Account</h3>
            <div class="account-info">
                <div><b>Email:</b> {$account.email}</div>
                <div><b>License Key:</b> {$account.licenseKey ? `${$account.licenseKey.substring(0, 20)}...` : 'Not set'}</div>
                <div><b>License Expiry:</b> {$account.licenseExpiry || 'Unknown'}</div>
                <div><b>Status:</b> 
                    {#if $licenseValidationStatus.isChecking}
                        <span class="checking">Checking...</span>
                    {:else if $licenseValid}
                        <span class="valid">Valid</span>
                    {:else}
                        <span class="expired">Invalid/Expired</span>
                    {/if}
                </div>
                {#if $licenseValidationStatus.message}
                    <div class="license-message"><b>Message:</b> {$licenseValidationStatus.message}</div>
                {/if}
            </div>
            
            <div class="license-actions">
                <button class="refresh-license-btn" on:click={refreshLicense} disabled={$licenseValidationStatus.isChecking}>
                    {#if $licenseValidationStatus.isChecking}🔄{:else}🔍{/if} Check License
                </button>
                <button class="edit-license-btn" on:click={() => showLicenseInput = !showLicenseInput}>Update License</button>
            </div>
            
            {#if showLicenseInput}
                <div class="license-edit-form">
                    <input 
                        type="text" 
                        placeholder="Enter new license key..." 
                        bind:value={newKey}
                        disabled={validatingLicense}
                    />
                    <div class="form-buttons">
                        <button 
                            class="save-btn" 
                            on:click={handleLicenseSave}
                            disabled={validatingLicense || !newKey.trim()}
                        >
                            {#if validatingLicense}🔄 Validating...{:else}💾 Save & Validate{/if}
                        </button>
                        <button class="cancel-btn" on:click={() => { showLicenseInput = false; newKey = ''; }}>
                            ❌ Cancel
                        </button>
                    </div>
                </div>
            {/if}
        </div>
        
        <div class="section danger-section">
            <h3>⚠️ Attention Required</h3>
            <p class="danger-warning">These actions cannot be undone. Please be careful!</p>
            <button class="delete-all-btn" on:click={() => showDeleteAllConfirm = true}>
                🗑️ Delete All Patients
            </button>
        </div>
    </aside>
{/if}

{#if showDeleteAllConfirm}
    <div class="modal-backdrop" on:click={() => showDeleteAllConfirm = false}></div>
    <div class="confirmation-modal">
        <div class="modal-header">
            <h3>⚠️ Delete All Patients</h3>
        </div>
        <div class="modal-content">
            <p><strong>This action will permanently delete:</strong></p>
            <ul class="deletion-list">
                <li>All patient records</li>
                <li>All patient appointments</li>
                <li>All patient payment records</li>
                <li>All patient folder data</li>
            </ul>
            <p class="final-warning"><strong>This action cannot be undone!</strong></p>
        </div>
        <div class="modal-actions">
            <button 
                class="confirm-delete-btn" 
                on:click={handleDeleteAllPatients}
                disabled={deletingAllPatients}
            >
                {#if deletingAllPatients}🔄 Deleting...{:else}✅ Yes, Delete All{/if}
            </button>
            <button 
                class="cancel-btn" 
                on:click={() => showDeleteAllConfirm = false}
                disabled={deletingAllPatients}
            >
                ❌ Cancel
            </button>
        </div>
    </div>
{/if}

{#if showUserManagement}
    <div class="user-mgmt-modal-overlay" on:click={() => showUserManagement = false}>
        <div class="user-mgmt-modal-content" on:click|stopPropagation>
            <div class="user-mgmt-modal-header">
                <h3>User Management</h3>
                <button class="close-btn" on:click={() => showUserManagement = false}>×</button>
            </div>
            <div class="user-mgmt-modal-body">
                <UserManagement />
            </div>
        </div>
    </div>
{/if}

<style>
.sidebar-backdrop {
    position: fixed;
    top: 0; left: 0; right: 0; bottom: 0;
    background: rgba(0,0,0,0.18);
    z-index: 2000;
}
.settings-sidebar {
    position: fixed;
    top: 0; right: 0;
    width: 340px;
    height: 100vh;
    background: var(--color-card);
    color: var(--color-text);
    box-shadow: -4px 0 24px rgba(0,0,0,0.18);
    z-index: 2001;
    padding: 2.2rem 2rem 2rem 2rem;
    display: flex;
    flex-direction: column;
    gap: 2rem;
    border-radius: 18px 0 0 18px;
    animation: slideIn 0.25s cubic-bezier(.4,1.7,.6,1) 1;
}
@keyframes slideIn {
    from { right: -400px; opacity: 0; }
    to { right: 0; opacity: 1; }
}
.close-btn {
    position: absolute;
    top: 1.2rem;
    right: 1.2rem;
    background: none;
    border: none;
    color: #fff;
    font-size: 2rem;
    cursor: pointer;
}
.section {
    margin-bottom: 1.5rem;
}
.theme-toggle {
    background: var(--color-accent-gradient);
    color: #fff;
    border: none;
    border-radius: 8px;
    padding: 0.7rem 1.5rem;
    font-size: 1.1rem;
    font-weight: 600;
    cursor: pointer;
    margin-top: 0.7rem;
}
.account-info {
    margin-bottom: 1rem;
    font-size: 1.05rem;
}
.valid { color: #4caf50; font-weight: 600; }
.expired { color: var(--color-danger); font-weight: 600; }
.checking { color: #ff9800; font-weight: 600; }

.license-message {
    margin-top: 0.5rem;
    font-size: 0.9rem;
    color: #666;
}

.license-actions {
    display: flex;
    gap: 0.5rem;
    margin-top: 1rem;
}

.refresh-license-btn {
    background: #e3f2fd;
    color: #1976d2;
    border: 1px solid #1976d2;
    border-radius: 8px;
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    flex: 1;
}

.refresh-license-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.edit-license-btn {
    background: #fff0f0;
    color: var(--color-danger);
    border: 1px solid var(--color-danger);
    border-radius: 8px;
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    flex: 1;
}

.license-edit-form {
    display: flex;
    flex-direction: column;
    gap: 0.7rem;
    margin-top: 0.7rem;
}

.license-edit-form input {
    padding: 0.7rem;
    border-radius: 6px;
    border: 1px solid var(--color-border);
    font-size: 0.9rem;
    font-family: 'Courier New', monospace;
    background: var(--color-card);
    color: var(--color-text);
}

.form-buttons {
    display: flex;
    gap: 0.5rem;
}

.save-btn {
    background: var(--color-accent);
    color: white;
    border: none;
    border-radius: 6px;
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    flex: 1;
}

.save-btn:hover:not(:disabled) {
    background: #5a67d8;
}

.save-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.cancel-btn {
    background: #f5f5f5;
    color: #666;
    border: 1px solid #ddd;
    border-radius: 6px;
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    flex: 1;
}

.danger-section {
    border: 2px solid #ff6b6b;
    border-radius: 8px;
    padding: 1rem;
    background: #fff5f5;
}

.danger-section h3 {
    color: #d63031;
    margin: 0 0 0.5rem 0;
    font-size: 1.1rem;
}

.danger-warning {
    color: #666;
    font-size: 0.9rem;
    margin: 0 0 1rem 0;
    font-style: italic;
}

.delete-all-btn {
    background: linear-gradient(135deg, #ff6b6b, #ee5a52);
    color: white;
    border: none;
    border-radius: 8px;
    padding: 0.75rem 1.5rem;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    width: 100%;
    transition: all 0.3s ease;
}

.delete-all-btn:hover {
    background: linear-gradient(135deg, #ff5252, #d32f2f);
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(255, 107, 107, 0.3);
}

.confirmation-modal {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: white;
    border-radius: 12px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    z-index: 2002;
    min-width: 400px;
    max-width: 500px;
    width: 90vw;
}

.confirmation-modal .modal-header {
    background: #ff6b6b;
    color: white;
    padding: 1rem 1.5rem;
    border-radius: 12px 12px 0 0;
    margin: 0;
}

.confirmation-modal .modal-header h3 {
    margin: 0;
    font-size: 1.2rem;
}

.confirmation-modal .modal-content {
    padding: 1.5rem;
    color: #333;
}

.deletion-list {
    margin: 1rem 0;
    padding-left: 1.5rem;
}

.deletion-list li {
    margin: 0.5rem 0;
    color: #d63031;
    font-weight: 500;
}

.final-warning {
    color: #d63031;
    font-size: 1.1rem;
    text-align: center;
    margin-top: 1rem;
}

.modal-actions {
    display: flex;
    gap: 1rem;
    padding: 0 1.5rem 1.5rem 1.5rem;
}

.confirm-delete-btn {
    background: linear-gradient(135deg, #ff6b6b, #ee5a52);
    color: white;
    border: none;
    border-radius: 8px;
    padding: 0.75rem 1.5rem;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    flex: 1;
    transition: all 0.3s ease;
}

.confirm-delete-btn:hover:not(:disabled) {
    background: linear-gradient(135deg, #ff5252, #d32f2f);
}

.confirm-delete-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.modal-actions .cancel-btn {
    background: #f8f9fa;
    color: #666;
    border: 1px solid #ddd;
    flex: 1;
}

.logout-btn {
    background: linear-gradient(135deg, #ff6b6b 0%, #ee5a6f 100%);
    color: white;
    border: none;
    padding: 0.75rem 1.5rem;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    width: 100%;
    margin-top: 1rem;
    transition: all 0.2s;
}

.logout-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(255, 107, 107, 0.3);
}

.user-mgmt-btn {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    padding: 0.75rem 1.5rem;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    width: 100%;
    transition: all 0.2s;
}

.user-mgmt-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.user-mgmt-modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 3000;
}

.user-mgmt-modal-content {
    background: white;
    border-radius: 12px;
    width: 90%;
    max-width: 800px;
    max-height: 90vh;
    overflow: hidden;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
    display: flex;
    flex-direction: column;
}

.user-mgmt-modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid #e1e5e9;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
}

.user-mgmt-modal-header h3 {
    margin: 0;
}

.user-mgmt-modal-header .close-btn {
    background: rgba(255, 255, 255, 0.2);
    color: white;
    border: none;
    font-size: 2rem;
    cursor: pointer;
    padding: 0;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
}

.user-mgmt-modal-header .close-btn:hover {
    background: rgba(255, 255, 255, 0.3);
}

.user-mgmt-modal-body {
    padding: 0;
    overflow-y: auto;
    flex: 1;
}
</style> 