<script>
  import { onMount } from 'svelte';
  import { login, authError, authLoading } from '../stores/authStore.js';

  export let onLoginSuccess = () => {};

  let username = '';
  let password = '';
  let error = '';

  onMount(() => {
    // Clear any previous errors
    authError.set(null);
  });

  // Watch for auth errors from store
  $: if ($authError) {
    error = $authError;
  }

  async function handleLogin() {
    if (!username.trim() || !password.trim()) {
      error = 'Please enter both username and password';
      return;
    }

    error = '';
    authError.set(null);

    const result = await login(username, password);
    
    if (result.success) {
      onLoginSuccess(result.user);
    } else {
      error = result.message || 'Login failed';
    }
  }

  function handleKeyPress(event) {
    if (event.key === 'Enter') {
      handleLogin();
    }
  }
</script>

<div class="login-container">
  <div class="login-box">
    <div class="login-header">
      <h1>🦷 DentistApp</h1>
      <p class="subtitle">User Login</p>
    </div>

    <div class="login-form">
      {#if error}
        <div class="error-message">
          ❌ {error}
        </div>
      {/if}

      <div class="form-group">
        <label for="username">Username:</label>
        <input
          id="username"
          type="text"
          bind:value={username}
          on:keypress={handleKeyPress}
          placeholder="Enter your username"
          disabled={$authLoading}
          autocomplete="username"
        />
      </div>

      <div class="form-group">
        <label for="password">Password:</label>
        <input
          id="password"
          type="password"
          bind:value={password}
          on:keypress={handleKeyPress}
          placeholder="Enter your password"
          disabled={$authLoading}
          autocomplete="current-password"
        />
      </div>

      <button
        on:click={handleLogin}
        disabled={$authLoading || !username.trim() || !password.trim()}
        class="login-button"
      >
        {#if $authLoading}
          🔄 Logging in...
        {:else}
          🔑 Login
        {/if}
      </button>

      <div class="login-info">
        <p><strong>Default Admin Credentials:</strong></p>
        <p>Username: <code>admin</code></p>
        <p>Password: <code>admin123</code></p>
      </div>
    </div>
  </div>
</div>

<style>
  .login-container {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    z-index: 10000;
  }

  .login-box {
    background: white;
    border-radius: 16px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
    max-width: 450px;
    width: 100%;
    overflow: hidden;
  }

  .login-header {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    text-align: center;
    padding: 2.5rem 2rem;
  }

  .login-header h1 {
    margin: 0;
    font-size: 2.5rem;
    font-weight: 700;
  }

  .subtitle {
    margin: 0.5rem 0 0 0;
    font-size: 1.1rem;
    opacity: 0.9;
  }

  .login-form {
    padding: 2.5rem;
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

  .form-group input {
    width: 100%;
    padding: 1rem;
    border: 2px solid #e1e5e9;
    border-radius: 8px;
    font-size: 1rem;
    transition: border-color 0.2s;
    box-sizing: border-box;
  }

  .form-group input:focus {
    outline: none;
    border-color: #667eea;
  }

  .form-group input:disabled {
    background-color: #f5f5f5;
    cursor: not-allowed;
  }

  .error-message {
    background: #fee;
    color: #c53030;
    padding: 1rem;
    border-radius: 8px;
    margin-bottom: 1.5rem;
    text-align: center;
    border: 1px solid #fed7d7;
  }

  .login-button {
    width: 100%;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    padding: 1rem 2rem;
    border-radius: 8px;
    font-size: 1.1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    margin-bottom: 1.5rem;
  }

  .login-button:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 10px 30px rgba(102, 126, 234, 0.3);
  }

  .login-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none;
  }

  .login-info {
    border-top: 1px solid #e1e5e9;
    padding-top: 1.5rem;
    text-align: center;
    color: #666;
    font-size: 0.9rem;
  }

  .login-info p {
    margin: 0.5rem 0;
  }

  .login-info code {
    background: #f1f3fa;
    padding: 0.2rem 0.5rem;
    border-radius: 4px;
    font-family: 'Courier New', monospace;
    color: #667eea;
    font-weight: 600;
  }
</style>

