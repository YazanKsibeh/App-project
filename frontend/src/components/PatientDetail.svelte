<script>
  import { createEventDispatcher } from 'svelte';

  export let patient;

  const dispatch = createEventDispatcher();

  function goBack() {
    dispatch('back');
  }

  // Get gender icon
  function getGenderIcon(gender) {
    return gender.toLowerCase() === 'male' ? '👨' : gender.toLowerCase() === 'female' ? '👩' : '👤';
  }
</script>

<div class="patient-detail-view">
  <button class="back-btn" on:click={goBack}>
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <line x1="19" y1="12" x2="5" y2="12"/>
      <polyline points="12 19 5 12 12 5"/>
    </svg>
    Back to List
  </button>

  <div class="detail-card">
    <div class="detail-header">
      <div class="patient-avatar">{getGenderIcon(patient.gender)}</div>
      <div class="patient-info">
        <h1 class="patient-name">{patient.name}</h1>
        <p class="patient-id">ID: {patient.id}</p>
      </div>
    </div>
    <div class="detail-content">
      <div class="info-item">
        <span class="label">📞 Phone</span>
        <span class="value">{patient.phone}</span>
      </div>
      <div class="info-item">
        <span class="label">🎂 Age</span>
        <span class="value">{patient.age} years</span>
      </div>
      <div class="info-item">
        <span class="label">⚧ Gender</span>
        <span class="value">{patient.gender}</span>
      </div>
    </div>
  </div>
</div>

<style>
  .patient-detail-view {
    padding: 2rem;
    max-width: 800px;
    margin: 2rem auto;
    animation: fadeIn 0.5s ease-out;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .back-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    background: none;
    border: 1px solid #ddd;
    color: var(--color-text);
    padding: 0.5rem 1rem;
    border-radius: 8px;
    font-size: 1rem;
    cursor: pointer;
    margin-bottom: 2rem;
    transition: all 0.2s ease;
  }

  body[data-theme="dark"] .back-btn {
    color: var(--color-text-contrast);
    border-color: var(--color-border);
  }

  .back-btn:hover {
    background: #f5f5f5;
    border-color: #ccc;
  }

  .back-btn svg {
    width: 18px;
    height: 18px;
  }

  .detail-card {
    background: white;
    border-radius: 16px;
    box-shadow: 0 8px 30px rgba(0,0,0,0.12);
  }
  
  .detail-header {
    display: flex;
    align-items: center;
    gap: 1.5rem;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 2rem;
    border-top-left-radius: 16px;
    border-top-right-radius: 16px;
  }

  .patient-avatar {
    font-size: 3rem;
    width: 80px;
    height: 80px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(255,255,255,0.2);
    border-radius: 50%;
    flex-shrink: 0;
  }
  
  .patient-name {
    font-size: 2rem;
    font-weight: 700;
    margin: 0 0 0.25rem 0;
    text-shadow: 0 1px 2px rgba(0,0,0,0.2);
  }

  .patient-id {
    font-size: 1rem;
    opacity: 0.9;
    margin: 0;
  }

  .detail-content {
    padding: 2rem;
    display: grid;
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .info-item {
    display: flex;
    justify-content: space-between;
    padding: 1rem;
    background: #f9f9f9;
    border-radius: 8px;
    border: 1px solid #eee;
  }
  
  .label {
    font-weight: 600;
    color: #555;
  }

  .value {
    color: #333;
  }
</style> 