<script>
  import { onMount } from 'svelte';
  import { CreateUser, GetAllUsers } from '../../wailsjs/go/main/App.js';
  import { currentLicenseKey, account } from '../stores/settingsStore.js';
  import { currentUser, isAdmin } from '../stores/authStore.js';
  import { get } from 'svelte/store';

  let users = [];
  let loading = false;
  let error = '';
  let showAddModal = false;
  
  let newUser = {
    username: '',
    password: '',
    role: 'Dentist'
  };

  const roles = ['Admin', 'Dentist', 'Assistant'];

  // Helper function to get current license key
  function getLicenseKey() {
    let licenseKey = '';
    try {
      licenseKey = get(currentLicenseKey);
    } catch (e) {
      licenseKey = localStorage.getItem('dentist_license_key') || '';
    }
    return licenseKey;
  }

  onMount(() => {
    loadUsers();
  });

  async function loadUsers() {
    loading = true;
    error = '';
    
    try {
      const licenseKey = getLicenseKey();
      users = await GetAllUsers(licenseKey);
    } catch (err) {
      error = err.message || 'Failed to load users';
    } finally {
      loading = false;
    }
  }

  async function handleCreateUser() {
    if (!newUser.username.trim() || !newUser.password.trim()) {
      error = 'Please fill in all required fields';
      return;
    }

    if (newUser.password.length < 6) {
      error = 'Password must be at least 6 characters';
      return;
    }

    loading = true;
    error = '';

    try {
      const licenseKey = getLicenseKey();
      const currentUserData = get(currentUser);
      
      if (!currentUserData || !isAdmin()) {
        throw new Error('Only admins can create users');
      }

      await CreateUser(newUser, currentUserData.id, licenseKey);
      
      // Reset form
      newUser = {
        username: '',
        password: '',
        role: 'Dentist'
      };
      showAddModal = false;
      
      // Reload users
      await loadUsers();
    } catch (err) {
      error = err.message || 'Failed to create user';
    } finally {
      loading = false;
    }
  }

  function openAddModal() {
    showAddModal = true;
    error = '';
    newUser = {
      username: '',
      password: '',
      role: 'Dentist'
    };
  }

  function closeAddModal() {
    showAddModal = false;
    error = '';
  }
</script>

<div class="user-management">
  <div class="header">
    <h2>User Management</h2>
    {#if isAdmin()}
      <button class="add-user-btn" on:click={openAddModal} disabled={loading}>
        + Add User
      </button>
    {/if}
  </div>

  {#if error}
    <div class="error-message">
      ⚠️ {error}
    </div>
  {/if}

  {#if loading}
    <div class="loading">
      <div class="spinner"></div>
      <p>Loading users...</p>
    </div>
  {:else if users.length === 0}
    <div class="empty-state">
      <p>No users found</p>
    </div>
  {:else}
    <div class="users-table">
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Username</th>
            <th>Role</th>
          </tr>
        </thead>
        <tbody>
          {#each users as user}
            <tr>
              <td>{user.id}</td>
              <td>{user.username}</td>
              <td>
                <span class="role-badge role-{user.role.toLowerCase()}">
                  {user.role}
                </span>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

{#if showAddModal}
  <div class="modal-overlay" on:click={closeAddModal}>
    <div class="modal-content" on:click|stopPropagation>
      <div class="modal-header">
        <h3>Add New User</h3>
        <button class="close-btn" on:click={closeAddModal}>×</button>
      </div>

      <div class="modal-body">
        {#if error}
          <div class="error-message">{error}</div>
        {/if}

        <div class="form-group">
          <label for="username">Username:</label>
          <input
            id="username"
            type="text"
            bind:value={newUser.username}
            placeholder="Enter username"
            disabled={loading}
          />
        </div>

        <div class="form-group">
          <label for="password">Password:</label>
          <input
            id="password"
            type="password"
            bind:value={newUser.password}
            placeholder="Enter password (min 6 characters)"
            disabled={loading}
          />
        </div>

        <div class="form-group">
          <label for="role">Role:</label>
          <select id="role" bind:value={newUser.role} disabled={loading}>
            {#each roles as role}
              <option value={role}>{role}</option>
            {/each}
          </select>
        </div>
      </div>

      <div class="modal-footer">
        <button class="btn-cancel" on:click={closeAddModal} disabled={loading}>
          Cancel
        </button>
        <button class="btn-submit" on:click={handleCreateUser} disabled={loading}>
          {loading ? 'Creating...' : 'Create User'}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .user-management {
    padding: 2rem;
    max-width: 1200px;
    margin: 0 auto;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }

  .header h2 {
    margin: 0;
    color: var(--color-text, #fff);
  }

  .add-user-btn {
    background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    padding: 0.75rem 1.5rem;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .add-user-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  }

  .add-user-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .error-message {
    background: #fee;
    border: 1px solid #fcc;
    border-radius: 8px;
    padding: 1rem;
    margin-bottom: 1.5rem;
    color: #c33;
  }

  .loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 4rem;
    color: var(--color-text, #fff);
  }

  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid rgba(255, 255, 255, 0.3);
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
    color: var(--color-text, #fff);
  }

  .users-table {
    background: var(--color-card, #fff);
    border-radius: 12px;
    overflow: hidden;
    box-shadow: var(--color-shadow, 0 2px 8px rgba(0,0,0,0.1));
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  thead {
    background: var(--color-navbar, #f1f3fa);
  }

  th {
    padding: 1rem;
    text-align: left;
    font-weight: 600;
    color: var(--color-text, #333);
    border-bottom: 2px solid var(--color-border, #e1e5e9);
  }

  td {
    padding: 1rem;
    border-bottom: 1px solid var(--color-border, #e1e5e9);
    color: var(--color-text, #333);
  }

  tbody tr:hover {
    background: var(--color-panel, #f8f9fa);
  }

  .role-badge {
    display: inline-block;
    padding: 0.25rem 0.75rem;
    border-radius: 12px;
    font-size: 0.875rem;
    font-weight: 600;
  }

  .role-badge.role-admin {
    background: #fee;
    color: #c53030;
  }

  .role-badge.role-dentist {
    background: #e6f3ff;
    color: #1e40af;
  }

  .role-badge.role-assistant {
    background: #f0fdf4;
    color: #166534;
  }

  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modal-content {
    background: white;
    border-radius: 12px;
    width: 90%;
    max-width: 500px;
    max-height: 90vh;
    overflow-y: auto;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid #e1e5e9;
  }

  .modal-header h3 {
    margin: 0;
    color: #333;
  }

  .close-btn {
    background: none;
    border: none;
    font-size: 2rem;
    color: #666;
    cursor: pointer;
    padding: 0;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
  }

  .close-btn:hover {
    background: #f1f3fa;
  }

  .modal-body {
    padding: 1.5rem;
  }

  .form-group {
    margin-bottom: 1.5rem;
  }

  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    color: #333;
    font-weight: 600;
  }

  .form-group input,
  .form-group select {
    width: 100%;
    padding: 0.75rem;
    border: 2px solid #e1e5e9;
    border-radius: 8px;
    font-size: 1rem;
    box-sizing: border-box;
  }

  .form-group input:focus,
  .form-group select:focus {
    outline: none;
    border-color: #667eea;
  }

  .form-group input:disabled,
  .form-group select:disabled {
    background-color: #f5f5f5;
    cursor: not-allowed;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    padding: 1.5rem;
    border-top: 1px solid #e1e5e9;
  }

  .btn-cancel,
  .btn-submit {
    padding: 0.75rem 1.5rem;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    border: none;
    transition: all 0.2s;
  }

  .btn-cancel {
    background: #e1e5e9;
    color: #333;
  }

  .btn-cancel:hover:not(:disabled) {
    background: #cbd5e1;
  }

  .btn-submit {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
  }

  .btn-submit:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  }

  .btn-cancel:disabled,
  .btn-submit:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>

