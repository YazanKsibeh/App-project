<script>
  import { currentUser, logout } from '../stores/authStore.js';
  import { get } from 'svelte/store';

  export let open = false;
  export let onClose = () => {};

  function handleLogout() {
    if (confirm('Are you sure you want to logout?')) {
      logout();
      onClose();
      // The App component will handle showing the login screen
      window.location.reload();
    }
  }

  function handleKeydown(event) {
    if (event.key === 'Escape') {
      onClose();
    }
  }
</script>

<svelte:window on:keydown={handleKeydown} />

{#if open}
  <div class="modal-overlay" on:click={onClose} on:keydown={(e) => e.key === 'Escape' && onClose()}>
    <div class="modal-content" on:click|stopPropagation>
      <div class="modal-header">
        <h3>User Profile</h3>
        <button class="close-btn" on:click={onClose} title="Close" aria-label="Close">×</button>
      </div>

      <div class="modal-body">
        {#if $currentUser}
          <div class="user-info">
            <div class="user-avatar">
              <span class="avatar-text">{$currentUser.username.charAt(0).toUpperCase()}</span>
            </div>
            <div class="user-details">
              <div class="detail-row">
                <span class="detail-label">Username:</span>
                <span class="detail-value">{$currentUser.username}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Role:</span>
                <span class="detail-value role-badge role-{$currentUser.role.toLowerCase()}">
                  {$currentUser.role}
                </span>
              </div>
            </div>
          </div>
        {:else}
          <div class="no-user">
            <p>No user information available</p>
          </div>
        {/if}
      </div>

      <div class="modal-footer">
        <button class="logout-btn" on:click={handleLogout}>
          🚪 Logout
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
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
    z-index: 2000;
    animation: fadeIn 0.2s ease;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  .modal-content {
    background: var(--color-card);
    border-radius: 12px;
    width: 90%;
    max-width: 400px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
    display: flex;
    flex-direction: column;
    animation: slideUp 0.3s ease;
    border: 1px solid var(--color-border);
  }

  @keyframes slideUp {
    from {
      transform: translateY(20px);
      opacity: 0;
    }
    to {
      transform: translateY(0);
      opacity: 1;
    }
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid var(--color-border);
    background: var(--color-accent-gradient);
    color: white;
  }

  .modal-header h3 {
    margin: 0;
    font-size: 1.5rem;
    color: white;
  }

  .close-btn {
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
    transition: background 0.2s;
    line-height: 1;
  }

  .close-btn:hover {
    background: rgba(255, 255, 255, 0.3);
  }

  .modal-body {
    padding: 2rem;
    color: var(--color-text);
  }

  .user-info {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.5rem;
  }

  .user-avatar {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    background: var(--color-accent-gradient);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 2rem;
    font-weight: 700;
    color: white;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  }

  .avatar-text {
    color: white;
  }

  .user-details {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .detail-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem;
    background: var(--color-panel);
    border-radius: 8px;
    border: 1px solid var(--color-border);
  }

  .detail-label {
    font-weight: 600;
    color: var(--color-text);
    opacity: 0.8;
  }

  .detail-value {
    color: var(--color-text);
    font-weight: 500;
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

  body[data-theme="dark"] .role-badge.role-admin {
    background: rgba(197, 48, 48, 0.2);
    color: #ff6b6b;
  }

  body[data-theme="dark"] .role-badge.role-dentist {
    background: rgba(30, 64, 175, 0.2);
    color: #60a5fa;
  }

  body[data-theme="dark"] .role-badge.role-assistant {
    background: rgba(22, 101, 52, 0.2);
    color: #4ade80;
  }

  .no-user {
    text-align: center;
    padding: 2rem;
    color: var(--color-text);
    opacity: 0.7;
  }

  .modal-footer {
    padding: 1.5rem;
    border-top: 1px solid var(--color-border);
    display: flex;
    justify-content: center;
  }

  .logout-btn {
    background: linear-gradient(135deg, #ff6b6b 0%, #ee5a6f 100%);
    color: white;
    border: none;
    padding: 0.75rem 2rem;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    font-size: 1rem;
    width: 100%;
  }

  .logout-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(255, 107, 107, 0.3);
    background: linear-gradient(135deg, #ff5252, #d32f2f);
  }
</style>

