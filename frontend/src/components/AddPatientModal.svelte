<script>
  import { createEventDispatcher, onMount } from 'svelte';
  import { addPatient, updatePatient, patients } from '../stores/patientStore.js';
  import { get } from 'svelte/store';
  
  const dispatch = createEventDispatcher();
  
  export let patientToEdit = null;
  
  let formData = {
    id: null,
    name: '',
    phone: '',
    age: '',
    gender: '',
    allergies: '',
    current_medications: '',
    medical_conditions: '',
    smoking_status: false,
    pregnancy_status: false,
    dental_history: '',
    special_notes: ''
  };
  
  let errors = {};
  let isSubmitting = false;
  let isEditMode = false;
  let medicalHistoryExpanded = false;
  
  onMount(() => {
    if (patientToEdit) {
      isEditMode = true;
      formData = {
        id: patientToEdit.id || null,
        name: patientToEdit.name || '',
        phone: patientToEdit.phone || '',
        age: patientToEdit.age || '',
        gender: patientToEdit.gender || '',
        allergies: patientToEdit.allergies || '',
        current_medications: patientToEdit.current_medications || '',
        medical_conditions: patientToEdit.medical_conditions || '',
        smoking_status: patientToEdit.smoking_status || false,
        pregnancy_status: patientToEdit.pregnancy_status || false,
        dental_history: patientToEdit.dental_history || '',
        special_notes: patientToEdit.special_notes || ''
      };
    }
  });
  
  const genderOptions = [
    { value: '', label: 'Select Gender' },
    { value: 'Male', label: 'Male' },
    { value: 'Female', label: 'Female' },
    { value: 'Other', label: 'Other' }
  ];
  
  function validateForm() {
    errors = {};
    
    if (!formData.name.trim()) {
      errors.name = 'Name is required';
    }
    
    if (!formData.phone.trim()) {
      errors.phone = 'Phone is required';
    } else if (!/^\d{10}$/.test(formData.phone.replace(/\D/g, ''))) {
      errors.phone = 'Please enter a valid 10-digit phone number';
    } else {
      // Check for phone number uniqueness
      const allPatients = get(patients);
      const phoneExists = allPatients.some(patient => 
        patient.phone === formData.phone && 
        (!isEditMode || patient.id !== formData.id)
      );
      if (phoneExists) {
        errors.phone = 'This phone number is already registered to another patient';
      }
    }
    
    if (!formData.age) {
      errors.age = 'Age is required';
    } else if (isNaN(formData.age) || formData.age < 6 || formData.age > 100) {
      errors.age = 'Please enter a valid age between 6 and 100 years';
    }
    
    if (!formData.gender) {
      errors.gender = 'Gender is required';
    }
    
    return Object.keys(errors).length === 0;
  }
  
  function formatPhone(event) {
    let value = event.target.value.replace(/\D/g, '');
    if (value.length > 10) {
      value = value.slice(0, 10);
    }
    formData.phone = value;
  }
  
  async function handleSubmit() {
    if (!validateForm()) {
      return;
    }
    
    isSubmitting = true;
    
    try {
      const patientData = {
        id: formData.id,
        name: formData.name.trim(),
        phone: formData.phone,
        age: parseInt(formData.age),
        gender: formData.gender,
        allergies: formData.allergies.trim(),
        current_medications: formData.current_medications.trim(),
        medical_conditions: formData.medical_conditions.trim(),
        smoking_status: formData.smoking_status,
        pregnancy_status: formData.pregnancy_status,
        dental_history: formData.dental_history.trim(),
        special_notes: formData.special_notes.trim()
      };
      
      if (isEditMode) {
        await updatePatient(patientData);
      } else {
        await addPatient(patientData);
      }
      handleClose();
    } catch (error) {
      console.error(`Error ${isEditMode ? 'updating' : 'adding'} patient:`, error);
    } finally {
      isSubmitting = false;
    }
  }
  
  function handleClose() {
    formData = {
      id: null,
      name: '',
      phone: '',
      age: '',
      gender: '',
      allergies: '',
      current_medications: '',
      medical_conditions: '',
      smoking_status: false,
      pregnancy_status: false,
      dental_history: '',
      special_notes: ''
    };
    errors = {};
    medicalHistoryExpanded = false;
    dispatch('close');
  }
  
  function toggleMedicalHistory() {
    medicalHistoryExpanded = !medicalHistoryExpanded;
  }
  
  function handleKeydown(event) {
    if (event.key === 'Escape') {
      handleClose();
    }
  }
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="modal-overlay" on:click={handleClose}>
  <div class="modal-content" on:click|stopPropagation>
    <div class="modal-header">
      <h2>{isEditMode ? 'Edit Patient' : 'Add New Patient'}</h2>
      <button class="close-btn" on:click={handleClose} title="Close">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18"/>
          <line x1="6" y1="6" x2="18" y2="18"/>
        </svg>
      </button>
    </div>
    
    <form on:submit|preventDefault={handleSubmit} class="patient-form">
      <div class="form-group">
        <label for="name">Full Name *</label>
        <input
          id="name"
          type="text"
          bind:value={formData.name}
          class="form-input {errors.name ? 'error' : ''}"
          placeholder="Enter patient's full name"
          required
        />
        {#if errors.name}
          <span class="error-message">{errors.name}</span>
        {/if}
      </div>
      
      <div class="form-group">
        <label for="phone">Phone Number *</label>
        <input
          id="phone"
          type="tel"
          bind:value={formData.phone}
          on:input={formatPhone}
          class="form-input {errors.phone ? 'error' : ''}"
          placeholder="Enter 10-digit phone number"
          maxlength="10"
          required
        />
        {#if errors.phone}
          <span class="error-message">{errors.phone}</span>
        {/if}
      </div>
      
      <div class="form-row">
        <div class="form-group">
          <label for="age">Age *</label>
          <input
            id="age"
            type="number"
            bind:value={formData.age}
            class="form-input {errors.age ? 'error' : ''}"
            placeholder="Age (6-100)"
            min="6"
            max="100"
            required
          />
          {#if errors.age}
            <span class="error-message">{errors.age}</span>
          {/if}
        </div>
        
        <div class="form-group">
          <label for="gender">Gender *</label>
          <select
            id="gender"
            bind:value={formData.gender}
            class="form-input {errors.gender ? 'error' : ''}"
            required
          >
            <option value="">Select gender</option>
            <option value="Male">Male</option>
            <option value="Female">Female</option>
          </select>
          {#if errors.gender}
            <span class="error-message">{errors.gender}</span>
          {/if}
        </div>
      </div>
      
      <!-- Medical History Section (Collapsible) -->
      <div class="medical-history-section">
        <button 
          type="button" 
          class="medical-history-header"
          on:click={toggleMedicalHistory}
          aria-expanded={medicalHistoryExpanded}
        >
          <span class="section-title">Medical History (Optional)</span>
          <svg 
            class="chevron {medicalHistoryExpanded ? 'expanded' : ''}" 
            viewBox="0 0 24 24" 
            fill="none" 
            stroke="currentColor" 
            stroke-width="2"
          >
            <polyline points="6 9 12 15 18 9"></polyline>
          </svg>
        </button>
        
        {#if medicalHistoryExpanded}
          <div class="medical-history-content">
            <div class="form-group">
              <label for="allergies">Allergies</label>
              <textarea
                id="allergies"
                bind:value={formData.allergies}
                class="form-textarea"
                placeholder="List any known allergies"
                rows="3"
              ></textarea>
            </div>
            
            <div class="form-group">
              <label for="current_medications">Current Medications</label>
              <textarea
                id="current_medications"
                bind:value={formData.current_medications}
                class="form-textarea"
                placeholder="List current medications"
                rows="3"
              ></textarea>
            </div>
            
            <div class="form-group">
              <label for="medical_conditions">Medical Conditions</label>
              <textarea
                id="medical_conditions"
                bind:value={formData.medical_conditions}
                class="form-textarea"
                placeholder="List any medical conditions"
                rows="3"
              ></textarea>
            </div>
            
            <div class="toggle-group">
              <div class="toggle-item">
                <label for="smoking_status" class="toggle-label">
                  <span>Smoking Status</span>
                  <div class="toggle-switch">
                    <input
                      type="checkbox"
                      id="smoking_status"
                      bind:checked={formData.smoking_status}
                      class="toggle-input"
                    />
                    <span class="toggle-slider"></span>
                  </div>
                </label>
              </div>
              
              <div class="toggle-item">
                <label for="pregnancy_status" class="toggle-label">
                  <span>Pregnancy Status</span>
                  <div class="toggle-switch">
                    <input
                      type="checkbox"
                      id="pregnancy_status"
                      bind:checked={formData.pregnancy_status}
                      class="toggle-input"
                    />
                    <span class="toggle-slider"></span>
                  </div>
                </label>
              </div>
            </div>
            
            <div class="form-group">
              <label for="dental_history">Dental History</label>
              <textarea
                id="dental_history"
                bind:value={formData.dental_history}
                class="form-textarea"
                placeholder="Previous dental treatments, surgeries, etc."
                rows="3"
              ></textarea>
            </div>
            
            <div class="form-group">
              <label for="special_notes">Special Notes</label>
              <textarea
                id="special_notes"
                bind:value={formData.special_notes}
                class="form-textarea"
                placeholder="Any additional notes or information"
                rows="3"
              ></textarea>
            </div>
          </div>
        {/if}
      </div>
      
      <div class="form-actions">
        <button type="button" class="btn btn-secondary" on:click={handleClose}>
          Cancel
        </button>
        <button type="submit" class="btn btn-primary" disabled={isSubmitting}>
          {#if isSubmitting}
            <div class="spinner"></div>
            {isEditMode ? 'Saving...' : 'Adding...'}
          {:else}
            {isEditMode ? 'Save Changes' : 'Add Patient'}
          {/if}
        </button>
      </div>
    </form>
  </div>
</div>

<style>
  .modal-overlay {
    position: fixed;
    top: 0; left: 0; right: 0; bottom: 0;
    background: rgba(0,0,0,0.18);
    z-index: 1000;
  }
  
  .modal-content {
    background: var(--color-card);
    color: var(--color-text);
    border-radius: 16px;
    box-shadow: var(--color-shadow);
    padding: 2rem 2.5rem 2rem 2.5rem;
    margin: 4vh auto;
    max-width: 500px;
    width: 100%;
    position: relative;
    border: 1px solid var(--color-border);
    max-height: 90vh;
    overflow-y: auto;
  }
  
  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1.5rem;
  }
  
  .modal-header h2 {
    margin: 0;
    color: var(--color-text);
    font-size: 1.5rem;
    font-weight: 600;
  }
  
  .close-btn {
    background: none;
    border: none;
    color: var(--color-text);
    font-size: 2rem;
    cursor: pointer;
    border-radius: 50%;
    transition: background 0.18s;
  }
  
  .close-btn:hover {
    background: var(--color-panel);
  }
  
  .close-btn svg {
    width: 20px;
    height: 20px;
  }
  
  .patient-form {
    padding: 0 1.5rem 1.5rem 1.5rem;
  }
  
  .form-group {
    margin-bottom: 1.5rem;
  }
  
  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }
  
  .patient-form label {
    color: var(--color-text);
  }
  
  .form-input {
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 0.6rem 1rem;
    font-size: 1rem;
    margin-bottom: 0.3rem;
  }
  
  .form-input:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }
  
  .form-input.error {
    border-color: var(--color-danger);
  }
  
  .form-input.error:focus {
    box-shadow: 0 0 0 3px rgba(231, 76, 60, 0.1);
  }
  
  .error-message {
    color: var(--color-danger);
    font-size: 0.95rem;
  }
  
  /* Medical History Section */
  .medical-history-section {
    margin-top: 1.5rem;
    margin-bottom: 1rem;
  }
  
  .medical-history-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    background: var(--color-panel);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 0.75rem 1rem;
    cursor: pointer;
    transition: all 0.2s ease;
    color: var(--color-text);
    font-size: 1rem;
    font-weight: 500;
  }
  
  .medical-history-header:hover {
    background: var(--color-border);
    border-color: var(--color-accent);
  }
  
  .section-title {
    color: var(--color-text);
  }
  
  .chevron {
    width: 18px;
    height: 18px;
    color: var(--color-text);
    transition: transform 0.3s ease;
  }
  
  .chevron.expanded {
    transform: rotate(180deg);
  }
  
  .medical-history-content {
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid var(--color-border);
    animation: slideDown 0.3s ease;
  }
  
  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
  
  .form-textarea {
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 0.6rem 1rem;
    font-size: 1rem;
    font-family: inherit;
    resize: vertical;
    min-height: 80px;
    width: 100%;
    box-sizing: border-box;
    margin-bottom: 0.3rem;
  }
  
  .form-textarea:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }
  
  .form-textarea::placeholder {
    color: var(--color-text);
    opacity: 0.5;
  }
  
  /* Toggle Switches */
  .toggle-group {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-bottom: 1.5rem;
  }
  
  .toggle-item {
    display: flex;
    align-items: center;
  }
  
  .toggle-label {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    cursor: pointer;
    color: var(--color-text);
    font-size: 1rem;
  }
  
  .toggle-switch {
    position: relative;
    width: 50px;
    height: 26px;
    flex-shrink: 0;
  }
  
  .toggle-input {
    opacity: 0;
    width: 0;
    height: 0;
  }
  
  .toggle-slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: var(--color-border);
    transition: 0.3s;
    border-radius: 26px;
  }
  
  .toggle-slider:before {
    position: absolute;
    content: "";
    height: 20px;
    width: 20px;
    left: 3px;
    bottom: 3px;
    background-color: white;
    transition: 0.3s;
    border-radius: 50%;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }
  
  .toggle-input:checked + .toggle-slider {
    background-color: #667eea;
  }
  
  .toggle-input:checked + .toggle-slider:before {
    transform: translateX(24px);
  }
  
  .toggle-input:focus + .toggle-slider {
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.2);
  }
  
  .form-actions {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
    margin-top: 2rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--color-border);
  }
  
  .btn {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  
  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .btn-secondary {
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
  }
  
  .btn-secondary:hover:not(:disabled) {
    background: var(--color-border);
  }
  
  .btn-primary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
  }
  
  .btn-primary:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  }
  
  .spinner {
    width: 16px;
    height: 16px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-top: 2px solid white;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }
  
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
  
  @media (max-width: 768px) {
    .modal-content {
      margin: 1rem;
      max-height: calc(100vh - 2rem);
    }
    
    .form-row {
      grid-template-columns: 1fr;
    }
    
    .form-actions {
      flex-direction: column-reverse;
    }
    
    .btn {
      width: 100%;
      justify-content: center;
    }
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
    max-width: 420px;
    width: 100%;
    box-sizing: border-box;
    overflow-x: hidden;
  }
  form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    width: 100%;
    box-sizing: border-box;
  }
  input, select, textarea {
    padding: 0.5rem;
    border-radius: 6px;
    border: 1px solid var(--color-border);
    background: var(--color-panel);
    color: var(--color-text);
    font-size: 1rem;
    width: 100%;
    box-sizing: border-box;
    min-width: 0;
  }
</style> 