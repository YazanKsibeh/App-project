<script>
    import { createEventDispatcher } from 'svelte';
    
    const dispatch = createEventDispatcher();
    
    let searchValue = '';
    let searchTimeout;
    
    function handleInput(event) {
        searchValue = event.target.value;
        
        // Clear previous timeout
        if (searchTimeout) {
            clearTimeout(searchTimeout);
        }
        
        // Debounce search to avoid too many API calls
        searchTimeout = setTimeout(() => {
            dispatch('search', searchValue);
        }, 300);
    }
    
    function handleClear() {
        searchValue = '';
        dispatch('search', '');
    }
</script>

<div class="search-container">
    <div class="search-input-wrapper">
        <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"/>
            <path d="m21 21-4.35-4.35"/>
        </svg>
        <input
            type="text"
            placeholder="Search patients by name or phone..."
            bind:value={searchValue}
            on:input={handleInput}
            class="search-input"
        />
        {#if searchValue}
            <button on:click={handleClear} class="clear-button" title="Clear search">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="18" y1="6" x2="6" y2="18"/>
                    <line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
            </button>
        {/if}
    </div>
</div>

<style>
    .search-container {
        min-width: 300px;
    }
    
    .search-input-wrapper {
        position: relative;
        display: flex;
        align-items: center;
    }
    
    .search-icon {
        position: absolute;
        left: 12px;
        width: 20px;
        height: 20px;
        color: #666;
        z-index: 1;
    }
    
    .search-input {
        width: 100%;
        padding: 12px 40px 12px 40px;
        border: none;
        border-radius: 25px;
        background: rgba(255, 255, 255, 0.95);
        font-size: 16px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        transition: all 0.3s ease;
    }
    
    .search-input:focus {
        outline: none;
        background: white;
        box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
        transform: translateY(-1px);
    }
    
    .search-input::placeholder {
        color: #999;
    }
    
    .clear-button {
        position: absolute;
        right: 8px;
        background: none;
        border: none;
        cursor: pointer;
        padding: 4px;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #666;
        transition: all 0.2s ease;
    }
    
    .clear-button:hover {
        background: rgba(0, 0, 0, 0.1);
        color: #333;
    }
    
    .clear-button svg {
        width: 16px;
        height: 16px;
    }
    
    @media (max-width: 768px) {
        .search-container {
            min-width: auto;
            width: 100%;
        }
        
        .search-input {
            font-size: 14px;
            padding: 10px 36px 10px 36px;
        }
    }
</style> 