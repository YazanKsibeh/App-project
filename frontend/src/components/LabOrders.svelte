<script>
  import { onMount } from 'svelte';
  import {
    filteredDentalLabs,
    dentalLabsLoading,
    dentalLabsError,
    dentalLabsSuccess,
    dentalLabSearch,
    dentalLabsCurrentPage,
    dentalLabsTotalPages,
    loadDentalLabsPaginated,
    getDentalLab,
    createDentalLab,
    updateDentalLab,
    deleteDentalLab
  } from '../stores/dentalLabStore.js';
  import {
    labOrders,
    labOrdersLoading,
    labOrdersError,
    labOrderSearch,
    labOrderStatusFilter,
    labOrdersCurrentPage,
    labOrdersTotalPages,
    labOrdersTotalCount,
    loadLabOrdersPaginated,
    getLabOrder,
    createLabOrder,
    formatCurrency,
    formatOrderDate,
    formatOrderDateTime
  } from '../stores/labOrderStore.js';
  import { patients, loadPatients } from '../stores/patientStore.js';
  import { currentUser } from '../stores/authStore.js';
  import { currentLicenseKey } from '../stores/settingsStore.js';
  import { get } from 'svelte/store';
  import { 
    GetDentalLabsPaginated,
    GetDentalLab
  } from '../../wailsjs/go/main/App.js';
  import { 
    GetWorkTypesPaginated 
  } from '../../wailsjs/go/main/App.js';
  import { 
    GetColorShadesPaginated 
  } from '../../wailsjs/go/main/App.js';

  let selectedSection = 'labs'; // 'labs', 'new-order', 'orders-list', 'tracking'

  // Lab modal state
  let showLabModal = false;
  let labModalMode = 'create'; // 'create' or 'edit'
  let labForm = {
    id: null,
    name: '',
    contact_person: '',
    phone_primary: '',
    phone_secondary: '',
    email: '',
    specialties: '',
    is_active: true,
    notes: ''
  };
  let labFormErrors = {};
  let labFormLoading = false;
  let showOptionalFields = false;
  let labToEdit = null;

  // Lab details modal state
  let showLabDetailModal = false;
  let selectedLab = null;

  // Delete confirmation state
  let showDeleteConfirm = false;
  let labToDelete = null;

  // Order details modal state
  let showOrderDetailModal = false;
  let selectedOrder = null;

  onMount(() => {
    selectedSection = 'labs';
    loadDentalLabsPaginated(1);
    loadPatients(); // Load patients for order form
  });

  function selectSection(section) {
    selectedSection = section;
    if (section === 'labs') {
      loadDentalLabsPaginated(1);
    } else if (section === 'orders-list') {
      loadLabOrdersPaginated(1, '', '', '', 'all');
    } else if (section === 'new-order') {
      initializeOrderForm();
      loadOrderDropdownData();
    }
  }

  // Phone formatting (same as AddPatientModal)
  function formatPhone(event, field) {
    let value = event.target.value.replace(/\D/g, '');
    if (value.length > 10) {
      value = value.slice(0, 10);
    }
    labForm[field] = value;
    // Clear error when user types
    if (labFormErrors[field]) {
      labFormErrors[field] = '';
    }
  }

  // Validation
  function validateLabForm() {
    labFormErrors = {};
    
    if (!labForm.name.trim()) {
      labFormErrors.name = 'Name is required';
    }
    
    if (!labForm.contact_person.trim()) {
      labFormErrors.contact_person = 'Contact person is required';
    }
    
    if (!labForm.phone_primary.trim()) {
      labFormErrors.phone_primary = 'Phone primary is required';
    } else if (!/^\d{10}$/.test(labForm.phone_primary.replace(/\D/g, ''))) {
      labFormErrors.phone_primary = 'Please enter a valid 10-digit phone number';
    }
    
    if (labForm.phone_secondary && !/^\d{10}$/.test(labForm.phone_secondary.replace(/\D/g, ''))) {
      labFormErrors.phone_secondary = 'Please enter a valid 10-digit phone number';
    }
    
    // Validate email format if provided
    if (labForm.email && labForm.email.trim()) {
      const emailRegex = /^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/;
      if (!emailRegex.test(labForm.email.trim())) {
        labFormErrors.email = 'Please enter a valid email address';
      }
    }
    
    return Object.keys(labFormErrors).length === 0;
  }

  // Open create lab modal
  function openCreateLabModal() {
    labModalMode = 'create';
    labForm = {
      id: null,
      name: '',
      contact_person: '',
      phone_primary: '',
      phone_secondary: '',
      email: '',
      specialties: '',
      is_active: true,
      notes: ''
    };
    labFormErrors = {};
    showOptionalFields = false;
    labToEdit = null;
    showLabModal = true;
  }

  // Open edit lab modal
  function openEditLabModal(lab) {
    labModalMode = 'edit';
    labToEdit = lab;
    labForm = {
      id: lab.id,
      name: lab.name || '',
      contact_person: lab.contact_person || '',
      phone_primary: lab.phone_primary || '',
      phone_secondary: lab.phone_secondary || '',
      email: lab.email || '',
      specialties: lab.specialties || '',
      is_active: lab.is_active !== undefined ? lab.is_active : true,
      notes: lab.notes || ''
    };
    labFormErrors = {};
    showOptionalFields = false;
    showLabModal = true;
  }

  // Handle lab save
  async function handleLabSave() {
    if (!validateLabForm()) {
      return;
    }
    
    labFormLoading = true;
    labFormErrors = {}; // Clear previous errors
    
    try {
      const formData = {
        name: labForm.name.trim(),
        contact_person: labForm.contact_person.trim(),
        phone_primary: labForm.phone_primary,
        phone_secondary: labForm.phone_secondary || '',
        email: labForm.email.trim() || '',
        specialties: labForm.specialties.trim() || '',
        is_active: labForm.is_active,
        notes: labForm.notes.trim() || ''
      };

      let success = false;
      if (labModalMode === 'create') {
        success = await createDentalLab(formData);
      } else {
        success = await updateDentalLab(labForm.id, formData);
      }
      
      // If operation failed, parse error from store
      if (!success) {
        const errorMessage = $dentalLabsError || 'Failed to save lab';
        parseAndDisplayError(errorMessage);
        return; // Don't close modal on error
      }
      
      // Success - close modal and reset form
      showLabModal = false;
      labForm = {
        id: null,
        name: '',
        contact_person: '',
        phone_primary: '',
        phone_secondary: '',
        email: '',
        specialties: '',
        is_active: true,
        notes: ''
      };
      labFormErrors = {};
      showOptionalFields = false;
    } catch (error) {
      console.error('Error saving lab:', error);
      const errorMessage = error.message || error.toString();
      parseAndDisplayError(errorMessage);
    } finally {
      labFormLoading = false;
    }
  }

  // Parse backend error messages and display them in form fields
  function parseAndDisplayError(errorMessage) {
    labFormErrors = {};
    
    // Map backend errors to form fields
    if (errorMessage.includes('lab name already exists')) {
      labFormErrors.name = 'Lab name already exists';
    } else if (errorMessage.includes('phone primary already exists')) {
      labFormErrors.phone_primary = 'Phone primary already exists';
    } else if (errorMessage.includes('phone secondary already exists')) {
      labFormErrors.phone_secondary = 'Phone secondary already exists';
    } else if (errorMessage.includes('email already exists')) {
      labFormErrors.email = 'Email already exists';
    } else if (errorMessage.includes('invalid email format')) {
      labFormErrors.email = 'Invalid email format';
    } else if (errorMessage.includes('phone primary must be exactly 10 digits')) {
      labFormErrors.phone_primary = 'Phone primary must be exactly 10 digits';
    } else if (errorMessage.includes('phone secondary must be exactly 10 digits')) {
      labFormErrors.phone_secondary = 'Phone secondary must be exactly 10 digits';
    } else if (errorMessage.includes('lab name is required')) {
      labFormErrors.name = 'Lab name is required';
    } else if (errorMessage.includes('contact person is required')) {
      labFormErrors.contact_person = 'Contact person is required';
    } else if (errorMessage.includes('phone primary is required')) {
      labFormErrors.phone_primary = 'Phone primary is required';
    } else {
      // Generic error - display at top of form
      labFormErrors.general = errorMessage;
    }
  }

  // Open lab details modal
  async function openLabDetailModal(lab) {
    try {
      const fullLab = await getDentalLab(lab.id);
      selectedLab = fullLab;
      showLabDetailModal = true;
    } catch (error) {
      console.error('Error loading lab details:', error);
    }
  }

  // Close lab details modal
  function closeLabDetailModal() {
    showLabDetailModal = false;
    selectedLab = null;
  }

  // Open edit from details modal
  function editFromDetails() {
    if (selectedLab) {
      closeLabDetailModal();
      openEditLabModal(selectedLab);
    }
  }

  // Confirm delete
  function confirmDeleteLab(lab) {
    labToDelete = lab;
    showDeleteConfirm = true;
  }

  // Handle delete
  async function handleDeleteLab() {
    if (!labToDelete) return;
    const success = await deleteDentalLab(labToDelete.id);
    if (success) {
      showDeleteConfirm = false;
      labToDelete = null;
    }
  }

  // Pagination
  function previousLabs() {
    const current = $dentalLabsCurrentPage;
    if (current > 1) {
      loadDentalLabsPaginated(current - 1);
    }
  }

  function nextLabs() {
    const current = $dentalLabsCurrentPage;
    const total = $dentalLabsTotalPages;
    if (current < total) {
      loadDentalLabsPaginated(current + 1);
    }
  }

  // Order functions
  async function openOrderDetailModal(order) {
    try {
      const orderDetail = await getLabOrder(order.id);
      selectedOrder = orderDetail;
      showOrderDetailModal = true;
    } catch (err) {
      console.error('Failed to load order details:', err);
    }
  }

  function closeOrderDetailModal() {
    showOrderDetailModal = false;
    selectedOrder = null;
  }

  function handleOrderSearch() {
    const searchTerm = $labOrderSearch || '';
    loadLabOrdersPaginated(1, searchTerm, '', '', $labOrderStatusFilter || 'all');
  }

  function previousOrders() {
    const current = $labOrdersCurrentPage;
    if (current > 1) {
      const searchTerm = $labOrderSearch || '';
      loadLabOrdersPaginated(current - 1, searchTerm, '', '', $labOrderStatusFilter || 'all');
    }
  }

  function nextOrders() {
    const current = $labOrdersCurrentPage;
    const total = $labOrdersTotalPages;
    if (current < total) {
      const searchTerm = $labOrderSearch || '';
      loadLabOrdersPaginated(current + 1, searchTerm, '', '', $labOrderStatusFilter || 'all');
    }
  }

  // Status badge helpers
  function getStatusClass(status) {
    if (!status) return 'status-draft';
    const normalized = status.toLowerCase();
    switch (normalized) {
      case 'draft':
        return 'status-draft';
      case 'sent':
        return 'status-sent';
      case 'delivered':
        return 'status-delivered';
      case 'cancelled':
        return 'status-cancelled';
      // Keep backward compatibility for old statuses
      case 'in_progress':
        return 'status-in-progress';
      case 'ready':
        return 'status-ready';
      default:
        return 'status-draft';
    }
  }

  function getStatusText(status) {
    if (!status) return 'Draft';
    return status.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase());
  }

  // Order delete confirmation
  let showDeleteOrderConfirm = false;
  let orderToDelete = null;

  function confirmDeleteOrder(order) {
    orderToDelete = order;
    showDeleteOrderConfirm = true;
  }

  function handleDeleteOrder() {
    // TODO: Implement delete order functionality
    console.log('Delete order:', orderToDelete);
    showDeleteOrderConfirm = false;
    orderToDelete = null;
  }

  // New Order Form State
  let orderForm = {
    order_number: '', // Read-only, auto-generated after save
    patient_id: null,
    patient_search: '',
    lab_id: null,
    work_type_id: null,
    color_shade_id: null,
    description: '',
    upper_left: '',
    upper_right: '',
    lower_left: '',
    lower_right: '',
    quantity: 1,
    lab_cost: '',
    order_date: '',
    status: 'draft',
    notes: ''
  };
  let orderFormErrors = {};
  let orderFormLoading = false;
  let selectedOrderPatient = null;
  let filteredOrderPatients = [];
  let showPatientDropdownOrder = false;
  let dentalLabsList = [];
  let workTypesList = [];
  let colorShadesList = [];
  let orderSuccessMessage = '';

  // Load dropdown data when New Order section is selected
  async function loadOrderDropdownData() {
    try {
      const licenseKey = get(currentLicenseKey) || localStorage.getItem('dentist_license_key') || '';
      if (!licenseKey) {
        console.error('[LabOrders] No license key found');
        return;
      }

      console.log('[LabOrders] Loading dropdown data...');

      // Load patients
      await loadPatients();
      console.log('[LabOrders] Patients loaded:', $patients.length);

      // Load dental labs (all active labs)
      const labsResponse = await GetDentalLabsPaginated(1, 1000, licenseKey);
      console.log('[LabOrders] Labs response:', labsResponse);
      if (labsResponse && labsResponse.labs && Array.isArray(labsResponse.labs)) {
        dentalLabsList = [...labsResponse.labs.filter(lab => lab.is_active)];
        console.log('[LabOrders] Active labs loaded:', dentalLabsList.length, dentalLabsList);
      } else {
        console.warn('[LabOrders] No labs in response or labs property missing. Response:', labsResponse);
        dentalLabsList = [];
      }

      // Load work types
      const workTypesResponse = await GetWorkTypesPaginated(1, 1000, licenseKey);
      console.log('[LabOrders] Work types response:', workTypesResponse);
      if (workTypesResponse && workTypesResponse.work_types && Array.isArray(workTypesResponse.work_types)) {
        workTypesList = [...workTypesResponse.work_types];
        console.log('[LabOrders] Work types loaded:', workTypesList.length, workTypesList);
      } else {
        console.warn('[LabOrders] No work types in response or work_types property missing. Response:', workTypesResponse);
        workTypesList = [];
      }

      // Load color shades
      const colorShadesResponse = await GetColorShadesPaginated(1, 1000, licenseKey);
      console.log('[LabOrders] Color shades response:', colorShadesResponse);
      if (colorShadesResponse && colorShadesResponse.color_shades && Array.isArray(colorShadesResponse.color_shades)) {
        colorShadesList = [...colorShadesResponse.color_shades.filter(shade => shade.is_active)];
        console.log('[LabOrders] Active color shades loaded:', colorShadesList.length, colorShadesList);
      } else {
        console.warn('[LabOrders] No color shades in response or color_shades property missing. Response:', colorShadesResponse);
        colorShadesList = [];
      }
    } catch (err) {
      console.error('[LabOrders] Failed to load dropdown data:', err);
      console.error('[LabOrders] Error details:', err.message, err.stack);
    }
  }

  // Note: Order form initialization is handled in selectSection() function

  function initializeOrderForm() {
    // Set default order date to now
    const now = new Date();
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
    orderForm.order_date = now.toISOString().slice(0, 16);
    
    // Reset form
    orderForm = {
      order_number: '',
      patient_id: null,
      patient_search: '',
      lab_id: null,
      work_type_id: null,
      color_shade_id: null,
      description: '',
      upper_left: '',
      upper_right: '',
      lower_left: '',
      lower_right: '',
      quantity: 1,
      lab_cost: '',
      order_date: now.toISOString().slice(0, 16),
      status: 'draft',
      notes: ''
    };
    orderFormErrors = {};
    selectedOrderPatient = null;
    orderSuccessMessage = '';
  }

  // Patient search for order form
  $: {
    if (orderForm.patient_search && !selectedOrderPatient) {
      const term = orderForm.patient_search.toLowerCase();
      filteredOrderPatients = $patients.filter(p => 
        p.name.toLowerCase().includes(term) || 
        (p.phone && p.phone.includes(term))
      ).slice(0, 5);
    } else {
      filteredOrderPatients = [];
    }
  }

  function handleOrderPatientSearch() {
    showPatientDropdownOrder = true;
  }

  function selectOrderPatient(patient) {
    selectedOrderPatient = patient;
    orderForm.patient_id = patient.id;
    orderForm.patient_search = patient.name;
    showPatientDropdownOrder = false;
  }

  function clearOrderPatient() {
    selectedOrderPatient = null;
    orderForm.patient_id = null;
    orderForm.patient_search = '';
  }

  // Helper function to validate tooth numbers for a quadrant
  function validateToothNumbers(value, allowedNumbers, quadrantName) {
    if (!value || value.trim() === '') {
      return null; // Empty is allowed (optional field)
    }

    // Parse the input - support comma, dash, space, or comma-space separators
    const numbers = value
      .split(/[,\s-]+/)
      .map(n => n.trim())
      .filter(n => n !== '')
      .map(n => parseInt(n))
      .filter(n => !isNaN(n));

    if (numbers.length === 0) {
      return `Invalid format. Please enter valid tooth numbers for ${quadrantName}.`;
    }

    // Check if all numbers are in the allowed list
    const invalidNumbers = numbers.filter(n => !allowedNumbers.includes(n));
    if (invalidNumbers.length > 0) {
      return `${quadrantName} only accepts: ${allowedNumbers.join(', ')}. Invalid: ${invalidNumbers.join(', ')}.`;
    }

    return null; // Valid
  }

  // Form validation
  function validateOrderForm() {
    orderFormErrors = {};

    if (!orderForm.patient_id) {
      orderFormErrors.patient_id = 'Patient is required';
    }
    if (!orderForm.lab_id) {
      orderFormErrors.lab_id = 'Dental lab is required';
    }
    if (!orderForm.work_type_id) {
      orderFormErrors.work_type_id = 'Work type is required';
    }
    if (!orderForm.order_date) {
      orderFormErrors.order_date = 'Order date is required';
    }
    if (!orderForm.lab_cost || orderForm.lab_cost === '' || orderForm.lab_cost === null) {
      orderFormErrors.lab_cost = 'Lab cost is required';
    } else if (parseFloat(orderForm.lab_cost) <= 0) {
      orderFormErrors.lab_cost = 'Lab cost must be greater than 0';
    }
    if (!orderForm.quantity || orderForm.quantity < 1) {
      orderFormErrors.quantity = 'Quantity must be at least 1';
    }

    // Validate quadrant tooth numbers
    const upperLeftError = validateToothNumbers(orderForm.upper_left, [12, 13, 14, 15, 16], 'Upper Left');
    if (upperLeftError) {
      orderFormErrors.upper_left = upperLeftError;
    }

    const upperRightError = validateToothNumbers(orderForm.upper_right, [1, 2, 3, 4, 5], 'Upper Right');
    if (upperRightError) {
      orderFormErrors.upper_right = upperRightError;
    }

    const lowerLeftError = validateToothNumbers(orderForm.lower_left, [17, 18, 19, 20, 21], 'Lower Left');
    if (lowerLeftError) {
      orderFormErrors.lower_left = lowerLeftError;
    }

    const lowerRightError = validateToothNumbers(orderForm.lower_right, [29, 30, 31, 32], 'Lower Right');
    if (lowerRightError) {
      orderFormErrors.lower_right = lowerRightError;
    }

    return Object.keys(orderFormErrors).length === 0;
  }

  // Save handlers
  async function handleSaveDraft() {
    if (!validateOrderForm()) {
      return;
    }

    orderFormLoading = true;
    orderFormErrors = {};
    orderSuccessMessage = '';

    try {
      const user = get(currentUser);
      if (!user || !user.id) {
        throw new Error('User not logged in');
      }

      const formData = {
        patient_id: parseInt(orderForm.patient_id) || 0,
        lab_id: parseInt(orderForm.lab_id) || 0,
        work_type_id: parseInt(orderForm.work_type_id) || 0,
        color_shade_id: orderForm.color_shade_id ? (parseInt(orderForm.color_shade_id) || null) : null,
        description: orderForm.description.trim(),
        upper_left: orderForm.upper_left.trim(),
        upper_right: orderForm.upper_right.trim(),
        lower_left: orderForm.lower_left.trim(),
        lower_right: orderForm.lower_right.trim(),
        quantity: parseInt(orderForm.quantity) || 1,
        lab_cost: orderForm.lab_cost !== '' && orderForm.lab_cost !== null ? parseInt(orderForm.lab_cost) : 0,
        order_date: orderForm.order_date,
        status: 'draft',
        notes: orderForm.notes.trim()
      };

      const result = await createLabOrder(formData, user.id);
      if (result.success) {
        orderForm.order_number = result.orderNumber;
        orderSuccessMessage = `Order ${result.orderNumber} saved successfully!`;
        // Navigate to orders list after 2 seconds
        setTimeout(() => {
          selectSection('orders-list');
          loadLabOrdersPaginated(1, '', '', '', 'all');
        }, 2000);
      } else {
        orderFormErrors.general = result.error || 'Failed to save order';
      }
    } catch (err) {
      console.error('Error saving draft:', err);
      orderFormErrors.general = err.message || 'Failed to save order';
    } finally {
      orderFormLoading = false;
    }
  }


  function handleCancelOrder() {
    initializeOrderForm();
    selectSection('orders-list');
  }
</script>

<div class="lab-orders">
  <div class="lab-orders-container">
    <!-- Left Sidebar Navigation -->
    <aside class="sidebar">
      <h2 class="sidebar-title">Lab Orders</h2>
      <nav class="sidebar-nav">
        <button 
          class="nav-item" 
          class:active={selectedSection === 'labs'}
          on:click={() => selectSection('labs')}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M9.5 2A2.5 2.5 0 0 1 12 4.5v15a2.5 2.5 0 0 1-4.96.44L2.5 19.5"/>
            <path d="M14.5 2A2.5 2.5 0 0 0 12 4.5v15a2.5 2.5 0 0 0 4.96.44L21.5 19.5"/>
          </svg>
          <span>Labs Management</span>
        </button>
        
        <button 
          class="nav-item" 
          class:active={selectedSection === 'new-order'}
          on:click={() => selectSection('new-order')}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          <span>New Order</span>
        </button>
        
        <button 
          class="nav-item" 
          class:active={selectedSection === 'orders-list'}
          on:click={() => selectSection('orders-list')}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="8" y1="6" x2="21" y2="6"/>
            <line x1="8" y1="12" x2="21" y2="12"/>
            <line x1="8" y1="18" x2="21" y2="18"/>
            <line x1="3" y1="6" x2="3.01" y2="6"/>
            <line x1="3" y1="12" x2="3.01" y2="12"/>
            <line x1="3" y1="18" x2="3.01" y2="18"/>
          </svg>
          <span>Orders List</span>
        </button>
        
        <button 
          class="nav-item" 
          class:active={selectedSection === 'tracking'}
          on:click={() => selectSection('tracking')}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"/>
            <polyline points="12 6 12 12 16 14"/>
          </svg>
          <span>Order Tracking</span>
        </button>
      </nav>
    </aside>

    <!-- Right Content Panel -->
    <main class="content-panel">
      <!-- Labs Management Section -->
      {#if selectedSection === 'labs'}
        <div class="section-content">
          <div class="section-header">
            <h1>Labs Management</h1>
            <p class="section-description">Manage your dental lab vendors.</p>
          </div>

          <div class="content-cards">
            <div class="action-card">
              <div class="labs-toolbar">
                <div class="search-group">
                  <input
                    type="text"
                    class="form-input search-input"
                    placeholder="Search labs..."
                    bind:value={$dentalLabSearch}
                  />
                </div>
                <button class="btn btn-primary" on:click={openCreateLabModal}>
                  ➕ Create New Lab
                </button>
              </div>

              {#if $dentalLabsError}
                <div class="alert alert-error">
                  {$dentalLabsError}
                </div>
              {/if}
              {#if $dentalLabsSuccess}
                <div class="alert alert-success">
                  {$dentalLabsSuccess}
                </div>
              {/if}

              <div class="labs-list">
                {#if $dentalLabsLoading}
                  <div class="loading-state">
                    <div class="spinner"></div>
                    <p>Loading labs...</p>
                  </div>
                {:else if $filteredDentalLabs.length === 0}
                  <div class="empty-state">
                    <p>No labs found. Add your first lab to get started.</p>
                  </div>
                {:else}
                  <div class="sessions-table-container">
                    <table class="sessions-table">
                      <thead>
                        <tr>
                          <th>Code</th>
                          <th>Name</th>
                          <th>Contact Person</th>
                          <th>Phone Primary</th>
                          <th>Phone Secondary</th>
                          <th class="actions-col">Actions</th>
                        </tr>
                      </thead>
                      <tbody>
                        {#each $filteredDentalLabs as lab}
                          <tr class="lab-row" on:click={() => openLabDetailModal(lab)}>
                            <td class="lab-code">{lab.code || '-'}</td>
                            <td class="lab-name">{lab.name}</td>
                            <td class="lab-contact-person">{lab.contact_person || '-'}</td>
                            <td>{lab.phone_primary || '-'}</td>
                            <td>{lab.phone_secondary || '-'}</td>
                            <td class="actions-col" on:click|stopPropagation>
                              <button class="icon-btn" on:click={() => openEditLabModal(lab)} title="Edit">
                                ✏️
                              </button>
                              <button class="icon-btn danger" on:click={() => confirmDeleteLab(lab)} title="Delete">
                                🗑️
                              </button>
                            </td>
                          </tr>
                        {/each}
                      </tbody>
                    </table>
                  </div>
                  
                  {#if $dentalLabsTotalPages > 1}
                    <div class="pagination">
                      <button 
                        class="page-btn" 
                        disabled={$dentalLabsCurrentPage === 1 || $dentalLabsLoading} 
                        on:click={previousLabs}
                      >
                        Previous
                      </button>
                      <span class="page-info">
                        Page {$dentalLabsCurrentPage} of {$dentalLabsTotalPages}
                      </span>
                      <button 
                        class="page-btn" 
                        disabled={$dentalLabsCurrentPage >= $dentalLabsTotalPages || $dentalLabsLoading} 
                        on:click={nextLabs}
                      >
                        Next
                      </button>
                    </div>
                  {/if}
                {/if}
              </div>
            </div>
          </div>
        </div>
      {/if}

      <!-- New Order Section -->
      {#if selectedSection === 'new-order'}
        <div class="section-content">
          <div class="section-header">
            <h1>New Lab Order</h1>
            <p class="section-description" style="display: none;">Create a new order for dental lab work.</p>
          </div>

          <div class="content-cards">
            {#if orderSuccessMessage}
              <div class="alert alert-success">
                {orderSuccessMessage}
              </div>
            {/if}

            {#if orderFormErrors.general}
              <div class="alert alert-error">
                {orderFormErrors.general}
              </div>
            {/if}

            <form on:submit|preventDefault class="order-form-grid">
              <!-- Card 1: ORDER DETAILS (Top-left) -->
              <div class="order-card">
                <h3 class="card-title">Order Details</h3>
                
                <div class="form-group form-group-inline">
                  <label for="orderDate">Order Date *</label>
                  <input
                    id="orderDate"
                    type="datetime-local"
                    bind:value={orderForm.order_date}
                    class="form-input {orderFormErrors.order_date ? 'error' : ''}"
                    required
                  />
                  {#if orderFormErrors.order_date}
                    <span class="error-message">{orderFormErrors.order_date}</span>
                  {/if}
                </div>

                <div class="form-group form-group-inline">
                  <label for="orderPatient">Patient *</label>
                  <div class="patient-search-container">
                    <input
                      id="orderPatient"
                      type="text"
                      bind:value={orderForm.patient_search}
                      on:input={handleOrderPatientSearch}
                      on:focus={() => showPatientDropdownOrder = true}
                      class="form-input search-input {orderFormErrors.patient_id ? 'error' : ''}"
                      placeholder="Search by name or phone..."
                      required
                    />
                    {#if showPatientDropdownOrder && orderForm.patient_search && !selectedOrderPatient}
                      <div class="dropdown">
                        {#if filteredOrderPatients.length > 0}
                          {#each filteredOrderPatients as patient}
                            <div 
                              class="dropdown-item"
                              on:click={() => selectOrderPatient(patient)}
                            >
                              <div class="patient-name">{patient.name}</div>
                              <div class="patient-phone">{patient.phone}</div>
                            </div>
                          {/each}
                        {:else}
                          <div class="dropdown-item empty">No patients found</div>
                        {/if}
                      </div>
                    {/if}
                  </div>
                  {#if selectedOrderPatient}
                    <div class="selected-patient-badge">
                      <span>{selectedOrderPatient.name} ({selectedOrderPatient.phone})</span>
                      <button type="button" class="clear-btn" on:click={clearOrderPatient}>×</button>
                    </div>
                  {/if}
                  {#if orderFormErrors.patient_id}
                    <span class="error-message">{orderFormErrors.patient_id}</span>
                  {/if}
                </div>

                <div class="form-group form-group-inline">
                  <label for="orderLab">Dental Lab *</label>
                  <select
                    id="orderLab"
                    bind:value={orderForm.lab_id}
                    class="form-input {orderFormErrors.lab_id ? 'error' : ''}"
                    required
                  >
                    <option value="">▾ Select Lab...</option>
                    {#each dentalLabsList as lab}
                      <option value={String(lab.id)}>{lab.name}</option>
                    {/each}
                  </select>
                  {#if dentalLabsList.length === 0}
                    <div style="font-size: 0.85rem; color: var(--text-secondary); margin-top: 0.25rem;">
                      No active labs available. Create a lab first.
                    </div>
                  {/if}
                  {#if orderFormErrors.lab_id}
                    <span class="error-message">{orderFormErrors.lab_id}</span>
                  {/if}
                </div>
              </div>

              <!-- Card 2: WORK SPECIFICATION (Top-right) -->
              <div class="order-card">
                <h3 class="card-title">Work Specification</h3>
                
                <div class="form-group form-group-inline">
                  <label for="workType">Work Type *</label>
                  <select
                    id="workType"
                    bind:value={orderForm.work_type_id}
                    class="form-input {orderFormErrors.work_type_id ? 'error' : ''}"
                    required
                  >
                    <option value="">▾ Select Work Type...</option>
                    {#each workTypesList as workType}
                      <option value={String(workType.id)}>{workType.name}</option>
                    {/each}
                  </select>
                  {#if workTypesList.length === 0}
                    <div style="font-size: 0.85rem; color: var(--text-secondary); margin-top: 0.25rem;">
                      No work types available. Create a work type first.
                    </div>
                  {/if}
                  {#if orderFormErrors.work_type_id}
                    <span class="error-message">{orderFormErrors.work_type_id}</span>
                  {/if}
                </div>

                <div class="form-group form-group-inline">
                  <label for="colorShade">Color/Shade</label>
                  <select
                    id="colorShade"
                    bind:value={orderForm.color_shade_id}
                    class="form-input"
                  >
                    <option value="">▾ Select Color Shade...</option>
                    {#each colorShadesList as shade}
                      <option value={String(shade.id)}>{shade.name} {shade.description ? `(${shade.description})` : ''}</option>
                    {/each}
                  </select>
                  {#if colorShadesList.length === 0}
                    <div style="font-size: 0.85rem; color: var(--text-secondary); margin-top: 0.25rem;">
                      No active color shades available. Create a color shade first.
                    </div>
                  {/if}
                </div>

                <div class="form-group form-group-inline">
                  <label for="quantity">Quantity</label>
                  <input
                    id="quantity"
                    type="number"
                    bind:value={orderForm.quantity}
                    class="form-input {orderFormErrors.quantity ? 'error' : ''}"
                    min="1"
                    required
                  />
                  {#if orderFormErrors.quantity}
                    <span class="error-message">{orderFormErrors.quantity}</span>
                  {/if}
                </div>
              </div>

              <!-- Card 3: TEETH/QUADRANT SELECTION (Bottom-left) -->
              <div class="order-card">
                <h3 class="card-title">Teeth/Quadrant Selection</h3>
                
                <div class="card-row">
                  <div class="form-group width-50">
                    <label for="upperLeft">Upper Left</label>
                    <input
                      id="upperLeft"
                      type="text"
                      bind:value={orderForm.upper_left}
                      class="form-input {orderFormErrors.upper_left ? 'error' : ''}"
                      placeholder="12-16"
                    />
                    {#if orderFormErrors.upper_left}
                      <span class="error-message">{orderFormErrors.upper_left}</span>
                    {/if}
                  </div>

                  <div class="form-group width-50">
                    <label for="upperRight">Upper Right</label>
                    <input
                      id="upperRight"
                      type="text"
                      bind:value={orderForm.upper_right}
                      class="form-input {orderFormErrors.upper_right ? 'error' : ''}"
                      placeholder="1-5"
                    />
                    {#if orderFormErrors.upper_right}
                      <span class="error-message">{orderFormErrors.upper_right}</span>
                    {/if}
                  </div>
                </div>

                <div class="card-row">
                  <div class="form-group width-50">
                    <label for="lowerLeft">Lower Left</label>
                    <input
                      id="lowerLeft"
                      type="text"
                      bind:value={orderForm.lower_left}
                      class="form-input {orderFormErrors.lower_left ? 'error' : ''}"
                      placeholder="17-21"
                    />
                    {#if orderFormErrors.lower_left}
                      <span class="error-message">{orderFormErrors.lower_left}</span>
                    {/if}
                  </div>

                  <div class="form-group width-50">
                    <label for="lowerRight">Lower Right</label>
                    <input
                      id="lowerRight"
                      type="text"
                      bind:value={orderForm.lower_right}
                      class="form-input {orderFormErrors.lower_right ? 'error' : ''}"
                      placeholder="29-32"
                    />
                    {#if orderFormErrors.lower_right}
                      <span class="error-message">{orderFormErrors.lower_right}</span>
                    {/if}
                  </div>
                </div>

                <div class="form-group">
                  <label for="description">Description</label>
                  <textarea
                    id="description"
                    bind:value={orderForm.description}
                    class="form-textarea"
                    placeholder="Add detailed description..."
                    rows="3"
                  ></textarea>
                </div>
              </div>

              <!-- Card 4: FINANCIALS & STATUS (Bottom-right) -->
              <div class="order-card">
                <h3 class="card-title">Financials & Status</h3>
                
                <div class="form-group">
                  <label>Status</label>
                  <div class="status-radio-grid">
                    <label class="radio-option" data-status="draft">
                      <input type="radio" bind:group={orderForm.status} value="draft" />
                      <span>Draft</span>
                    </label>
                    <label class="radio-option" data-status="sent">
                      <input type="radio" bind:group={orderForm.status} value="sent" />
                      <span>Sent</span>
                    </label>
                    <label class="radio-option" data-status="delivered">
                      <input type="radio" bind:group={orderForm.status} value="delivered" />
                      <span>Delivered</span>
                    </label>
                    <label class="radio-option" data-status="cancelled">
                      <input type="radio" bind:group={orderForm.status} value="cancelled" />
                      <span>Cancelled</span>
                    </label>
                  </div>
                </div>

                <div class="form-group">
                  <label for="labCost">Lab Cost</label>
                  <div class="input-with-suffix">
                    <input
                      id="labCost"
                      type="number"
                      bind:value={orderForm.lab_cost}
                      class="form-input {orderFormErrors.lab_cost ? 'error' : ''}"
                      min="0.01"
                      step="0.01"
                      required
                    />
                    <span class="input-suffix">SYP</span>
                  </div>
                  {#if orderFormErrors.lab_cost}
                    <span class="error-message">{orderFormErrors.lab_cost}</span>
                  {/if}
                </div>

                <div class="form-group">
                  <label for="orderNotes">Notes</label>
                  <textarea
                    id="orderNotes"
                    bind:value={orderForm.notes}
                    class="form-textarea"
                    placeholder="Internal notes about the order"
                    rows="3"
                  ></textarea>
                </div>
              </div>
            </form>

            <!-- Action Buttons (Bottom-right, outside cards) -->
            <div class="form-actions-fixed">
              <button type="button" class="btn btn-secondary" on:click={handleCancelOrder} disabled={orderFormLoading}>
                Cancel
              </button>
              <button type="button" class="btn btn-primary" on:click={handleSaveDraft} disabled={orderFormLoading}>
                {orderFormLoading ? 'Saving...' : 'Save'}
              </button>
            </div>
          </div>
        </div>
      {/if}

      <!-- Orders List Section -->
      {#if selectedSection === 'orders-list'}
        <div class="section-content">
          <div class="section-header">
            <h1>Orders List</h1>
            <p class="section-description">Track all lab orders and their status.</p>
          </div>

          <div class="content-cards">
            <div class="action-card">
              <div class="labs-toolbar">
                <div class="search-group">
                  <input
                    type="text"
                    class="form-input search-input"
                    placeholder="Search orders..."
                    bind:value={$labOrderSearch}
                    on:input={handleOrderSearch}
                  />
                </div>
                <select 
                  class="form-input status-filter"
                  bind:value={$labOrderStatusFilter}
                  on:change={handleOrderSearch}
                >
                  <option value="all">All Status</option>
                  <option value="draft">Draft</option>
                  <option value="sent">Sent</option>
                  <option value="delivered">Delivered</option>
                  <option value="cancelled">Cancelled</option>
                </select>
                <button class="btn btn-primary" on:click={() => selectSection('new-order')}>
                  ➕ Create New Order
                </button>
              </div>

              {#if $labOrdersError}
                <div class="alert alert-error">
                  {$labOrdersError}
                </div>
              {/if}

              <div class="labs-list">
                {#if $labOrdersLoading}
                  <div class="loading-state">
                    <div class="spinner"></div>
                    <p>Loading orders...</p>
                  </div>
                {:else if $labOrders.length === 0}
                  <div class="empty-state">
                    <p>No lab orders found. Create your first order to get started.</p>
                  </div>
                {:else}
                  <div class="sessions-table-container">
                    <table class="sessions-table">
                      <thead>
                        <tr>
                          <th>Order #</th>
                          <th>Patient</th>
                          <th>Lab</th>
                          <th>Work Type</th>
                          <th>Status</th>
                          <th>Cost</th>
                          <th>Date</th>
                          <th class="actions-col">Actions</th>
                        </tr>
                      </thead>
                      <tbody>
                        {#each $labOrders as order}
                          <tr class="lab-row" on:click={() => openOrderDetailModal(order)}>
                            <td class="order-number">{order.order_number || '-'}</td>
                            <td>{order.patient_name || 'Unknown'}</td>
                            <td>{order.lab_name || 'Unknown'}</td>
                            <td>{order.work_type_name || 'Unknown'}</td>
                            <td>
                              <span class="status-badge {getStatusClass(order.status)}">
                                <span class="status-icon">●</span>
                                <span class="status-text">{getStatusText(order.status)}</span>
                              </span>
                            </td>
                            <td class="order-cost">{formatCurrency(order.lab_cost || 0)} SYP</td>
                            <td>{formatOrderDate(order.order_date)}</td>
                            <td class="actions-col" on:click|stopPropagation>
                              <button class="icon-btn" on:click={() => openOrderDetailModal(order)} title="Edit">
                                ✏️
                              </button>
                              <button class="icon-btn danger" on:click={() => confirmDeleteOrder(order)} title="Delete">
                                🗑️
                              </button>
                            </td>
                          </tr>
                        {/each}
                      </tbody>
                    </table>
                  </div>
                  
                  {#if $labOrdersTotalPages > 1}
                    <div class="pagination">
                      <button 
                        class="page-btn" 
                        disabled={$labOrdersCurrentPage === 1 || $labOrdersLoading} 
                        on:click={previousOrders}
                      >
                        Previous
                      </button>
                      <span class="page-info">
                        Page {$labOrdersCurrentPage} of {$labOrdersTotalPages}
                      </span>
                      <button 
                        class="page-btn" 
                        disabled={$labOrdersCurrentPage >= $labOrdersTotalPages || $labOrdersLoading} 
                        on:click={nextOrders}
                      >
                        Next
                      </button>
                    </div>
                  {/if}
                {/if}
              </div>
            </div>
          </div>
        </div>
      {/if}

      <!-- Order Tracking Section -->
      {#if selectedSection === 'tracking'}
        <div class="section-content">
          <div class="section-header">
            <h1>Lab Orders</h1>
            <p class="section-description">This section is under development.</p>
          </div>

          <div class="content-cards">
            <div class="action-card">
              <h3>Order Tracking</h3>
              <p class="card-description">
                Features coming soon:
              </p>
              <ul class="feature-list">
                <li>Order Status Monitoring</li>
              </ul>
            </div>
          </div>
        </div>
      {/if}
    </main>
  </div>
</div>

<!-- Create/Edit Lab Modal -->
{#if showLabModal}
  <div 
    class="modal-overlay" 
    role="button"
    tabindex="0"
    on:click={() => showLabModal = false}
    on:keydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        showLabModal = false;
      }
    }}
  >
    <div class="modal-content" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header">
        <h2>{labModalMode === 'create' ? 'Create New Dental Lab' : 'Edit Dental Lab'}</h2>
        <button class="close-btn" on:click={() => showLabModal = false} title="Close">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
      
      <form on:submit|preventDefault={handleLabSave} class="lab-form">
        <!-- General Error Display -->
        {#if labFormErrors.general}
          <div class="error-banner">
            <span class="error-message">{labFormErrors.general}</span>
          </div>
        {/if}
        
        <!-- Required Fields Section -->
        <div class="form-section">
          <div class="form-group">
            <label for="labName">Name *</label>
            <input
              id="labName"
              type="text"
              bind:value={labForm.name}
              class="form-input {labFormErrors.name ? 'error' : ''}"
              placeholder="Enter lab name"
              required
            />
            {#if labFormErrors.name}
              <span class="error-message">{labFormErrors.name}</span>
            {/if}
          </div>
          
          <div class="form-group">
            <label for="contactPerson">Contact Person *</label>
            <input
              id="contactPerson"
              type="text"
              bind:value={labForm.contact_person}
              class="form-input {labFormErrors.contact_person ? 'error' : ''}"
              placeholder="Enter contact person name"
              required
            />
            {#if labFormErrors.contact_person}
              <span class="error-message">{labFormErrors.contact_person}</span>
            {/if}
          </div>
          
          <div class="form-group">
            <label for="phonePrimary">Phone Primary *</label>
            <input
              id="phonePrimary"
              type="tel"
              bind:value={labForm.phone_primary}
              on:input={(e) => formatPhone(e, 'phone_primary')}
              class="form-input {labFormErrors.phone_primary ? 'error' : ''}"
              placeholder="09******** or 011*******"
              maxlength="10"
              required
            />
            {#if labFormErrors.phone_primary}
              <span class="error-message">{labFormErrors.phone_primary}</span>
            {/if}
          </div>
          
          <div class="form-group">
            <div class="toggle-group">
              <div class="toggle-item">
                <label for="labActive" class="toggle-label">
                  <span>{labForm.is_active ? 'Active' : 'Inactive'}</span>
                  <div class="toggle-switch">
                    <input
                      type="checkbox"
                      id="labActive"
                      bind:checked={labForm.is_active}
                      class="toggle-input"
                      disabled={labFormLoading}
                    />
                    <span class="toggle-slider"></span>
                  </div>
                </label>
              </div>
            </div>
          </div>
        </div>

        <!-- Optional Fields Section (Collapsible) -->
        <div class="optional-fields-section">
          <button 
            type="button" 
            class="optional-fields-header"
            on:click={() => showOptionalFields = !showOptionalFields}
            aria-expanded={showOptionalFields}
          >
            <span class="section-title">Show Optional Fields {showOptionalFields ? '▼' : '▶'}</span>
          </button>
          
          {#if showOptionalFields}
            <div class="optional-fields-content">
              <div class="form-group">
                <label for="phoneSecondary">Phone Secondary</label>
                <input
                  id="phoneSecondary"
                  type="tel"
                  bind:value={labForm.phone_secondary}
                  on:input={(e) => formatPhone(e, 'phone_secondary')}
                  class="form-input {labFormErrors.phone_secondary ? 'error' : ''}"
                  placeholder="09******** or 011*******"
                  maxlength="10"
                />
                {#if labFormErrors.phone_secondary}
                  <span class="error-message">{labFormErrors.phone_secondary}</span>
                {/if}
              </div>
              
              <div class="form-group">
                <label for="labEmail">Email</label>
                <input
                  id="labEmail"
                  type="email"
                  bind:value={labForm.email}
                  on:input={() => {
                    if (labFormErrors.email) {
                      labFormErrors.email = '';
                    }
                  }}
                  class="form-input {labFormErrors.email ? 'error' : ''}"
                  placeholder="info@lab.com"
                />
                {#if labFormErrors.email}
                  <span class="error-message">{labFormErrors.email}</span>
                {/if}
              </div>
              
              <div class="form-group">
                <label for="labSpecialties">Specialties</label>
                <textarea
                  id="labSpecialties"
                  bind:value={labForm.specialties}
                  class="form-textarea"
                  placeholder="Crowns, Bridges, Implants"
                  rows="3"
                ></textarea>
              </div>
              
              <div class="form-group">
                <label for="labNotes">Notes</label>
                <textarea
                  id="labNotes"
                  bind:value={labForm.notes}
                  class="form-textarea"
                  placeholder="Additional notes about the lab"
                  rows="3"
                ></textarea>
              </div>
            </div>
          {/if}
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn-secondary" on:click={() => showLabModal = false} disabled={labFormLoading}>
            Cancel
          </button>
          <button type="submit" class="btn btn-primary" disabled={labFormLoading}>
            {#if labFormLoading}
              Saving...
            {:else if labModalMode === 'create'}
              Create Lab
            {:else}
              Save Changes
            {/if}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- Lab Details Modal -->
{#if showLabDetailModal && selectedLab}
  <div 
    class="modal-overlay" 
    role="button"
    tabindex="0"
    on:click={closeLabDetailModal}
    on:keydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        closeLabDetailModal();
      }
    }}
  >
    <div class="modal-content" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header">
        <h2>Lab Details: {selectedLab.name} ({selectedLab.code})</h2>
        <button class="close-btn" on:click={closeLabDetailModal} title="Close">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
      
      <div class="lab-details-content">
        <div class="detail-row">
          <span class="detail-label">Code</span>
          <span class="detail-value">{selectedLab.code || '-'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Name</span>
          <span class="detail-value">{selectedLab.name || '-'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Contact Person</span>
          <span class="detail-value">{selectedLab.contact_person || '-'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Phone Primary</span>
          <span class="detail-value">{selectedLab.phone_primary || '-'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Phone Secondary</span>
          <span class="detail-value">{selectedLab.phone_secondary || '-'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Email</span>
          <span class="detail-value">{selectedLab.email || '-'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Specialties</span>
          <span class="detail-value">{selectedLab.specialties || '-'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Active</span>
          <span class="detail-value">
            {#if selectedLab.is_active}
              ● Active
            {:else}
              ○ Inactive
            {/if}
          </span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Notes</span>
          <span class="detail-value">{selectedLab.notes || '-'}</span>
        </div>
      </div>

      <div class="form-actions">
        <button type="button" class="btn btn-secondary" on:click={closeLabDetailModal}>
          Close
        </button>
        <button type="button" class="btn btn-primary" on:click={editFromDetails}>
          Edit
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Order Details Modal -->
{#if showOrderDetailModal && selectedOrder}
  <div 
    class="modal-overlay" 
    role="button"
    tabindex="0"
    on:click={closeOrderDetailModal}
    on:keydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        closeOrderDetailModal();
      }
    }}
  >
    <div class="modal-content" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header">
        <h2>Order Details: {selectedOrder.order_number}</h2>
        <button class="close-btn" on:click={closeOrderDetailModal} title="Close">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
      
      <div class="lab-details-content">
        <div class="detail-row">
          <span class="detail-label">Order Number</span>
          <span class="detail-value">{selectedOrder.order_number || '-'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Patient</span>
          <span class="detail-value">{selectedOrder.patient_name || 'Unknown'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Lab</span>
          <span class="detail-value">{selectedOrder.lab_name || 'Unknown'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Dentist</span>
          <span class="detail-value">{selectedOrder.dentist_name || 'Unknown'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Work Type</span>
          <span class="detail-value">{selectedOrder.work_type_name || 'Unknown'}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Status</span>
          <span class="detail-value">
            <span class="status-badge {getStatusClass(selectedOrder.status)}">
              <span class="status-icon">●</span>
              <span class="status-text">{getStatusText(selectedOrder.status)}</span>
            </span>
          </span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Order Date</span>
          <span class="detail-value">{formatOrderDateTime(selectedOrder.order_date) || '-'}</span>
        </div>
        
        <div class="detail-section">
          <h3 class="detail-section-title">Quadrants</h3>
          <div class="detail-row">
            <span class="detail-label">Upper Left</span>
            <span class="detail-value">{selectedOrder.upper_left || '-'}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">Upper Right</span>
            <span class="detail-value">{selectedOrder.upper_right || '-'}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">Lower Left</span>
            <span class="detail-value">{selectedOrder.lower_left || '-'}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">Lower Right</span>
            <span class="detail-value">{selectedOrder.lower_right || '-'}</span>
          </div>
        </div>

        <div class="detail-section">
          <h3 class="detail-section-title">Specifications</h3>
          <div class="detail-row">
            <span class="detail-label">Color</span>
            <span class="detail-value">{selectedOrder.color_shade_name || '-'}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">Quantity</span>
            <span class="detail-value">{selectedOrder.quantity || 1} units</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">Cost</span>
            <span class="detail-value">{formatCurrency(selectedOrder.lab_cost || 0)} SYP</span>
          </div>
        </div>

        {#if selectedOrder.description}
          <div class="detail-row">
            <span class="detail-label">Description</span>
            <span class="detail-value">{selectedOrder.description}</span>
          </div>
        {/if}

        {#if selectedOrder.notes}
          <div class="detail-row">
            <span class="detail-label">Notes</span>
            <span class="detail-value">{selectedOrder.notes}</span>
          </div>
        {/if}
      </div>

      <div class="form-actions">
        <button type="button" class="btn btn-secondary" on:click={closeOrderDetailModal}>
          Close
        </button>
        <button type="button" class="btn btn-primary" disabled>
          Edit
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Delete Order Confirmation Modal -->
{#if showDeleteOrderConfirm && orderToDelete}
  <div 
    class="modal-overlay" 
    role="button"
    tabindex="0"
    on:click={() => showDeleteOrderConfirm = false}
    on:keydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        showDeleteOrderConfirm = false;
      }
    }}
  >
    <div class="confirmation-modal" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header danger-header">
        <h3>Delete Order</h3>
      </div>
      <div class="modal-content">
        <p>Are you sure you want to delete the order <strong>{orderToDelete.order_number}</strong>? This action cannot be undone.</p>
      </div>
      <div class="modal-actions">
        <button class="btn btn-danger" on:click={handleDeleteOrder}>Delete</button>
        <button class="btn btn-secondary" on:click={() => { showDeleteOrderConfirm = false; orderToDelete = null; }}>Cancel</button>
      </div>
    </div>
  </div>
{/if}

<!-- Delete Confirmation Modal -->
{#if showDeleteConfirm}
  <div 
    class="modal-overlay" 
    role="button"
    tabindex="0"
    on:click={() => showDeleteConfirm = false}
    on:keydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        showDeleteConfirm = false;
      }
    }}
  >
    <div class="confirmation-modal" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header danger-header">
        <h3>Delete Lab</h3>
      </div>
      <div class="modal-content">
        <p>Are you sure you want to delete the lab <strong>{labToDelete?.name}</strong>? This action cannot be undone.</p>
      </div>
      <div class="modal-actions">
        <button class="btn btn-danger" on:click={handleDeleteLab}>Delete</button>
        <button class="btn btn-secondary" on:click={() => { showDeleteConfirm = false; labToDelete = null; }}>Cancel</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .lab-orders {
    height: calc(100vh - 80px);
    background: var(--color-bg);
    color: var(--color-text);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .lab-orders-container {
    display: flex;
    height: 100%;
    overflow: hidden;
  }

  /* Sidebar Styles */
  .sidebar {
    width: 240px;
    background: var(--color-card);
    border-right: 1px solid var(--color-border);
    display: flex;
    flex-direction: column;
    padding: 1.5rem 0;
    overflow-y: auto;
  }

  .sidebar-title {
    font-size: 1.25rem;
    font-weight: 700;
    color: var(--color-text);
    padding: 0 1.5rem;
    margin: 0 0 1.5rem 0;
  }

  .sidebar-nav {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    padding: 0 0.75rem;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.625rem 1rem;
    border: none;
    background: transparent;
    color: var(--color-text);
    font-size: 0.9375rem;
    font-weight: 500;
    cursor: pointer;
    border-radius: 8px;
    transition: all 0.2s ease;
    text-align: left;
  }

  .nav-item:hover {
    background: var(--color-panel);
  }

  .nav-item.active {
    background: var(--color-accent);
    color: white;
  }

  .nav-item.danger.active {
    background: var(--color-danger);
    color: white;
  }

  .nav-item svg {
    flex-shrink: 0;
  }

  /* Content Panel Styles */
  .content-panel {
    flex: 1;
    overflow-y: auto;
    padding: 2rem;
    background: var(--color-bg);
  }

  .section-content {
    max-width: 1400px;
    animation: fadeIn 0.3s ease;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .section-header {
    margin-bottom: 2rem;
  }

  .section-header h1 {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--color-text);
    margin: 0 0 0.5rem 0;
  }

  .section-description {
    font-size: 0.9375rem;
    color: var(--color-text);
    opacity: 0.7;
    margin: 0;
  }

  .content-cards {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }

  /* Card Styles */
  .action-card {
    background: var(--color-card);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: 1.5rem;
    box-shadow: var(--color-shadow);
  }

  .action-card h3 {
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--color-text);
    margin: 0 0 1rem 0;
  }

  .card-description {
    font-size: 0.875rem;
    color: var(--color-text);
    opacity: 0.7;
    margin: 0 0 1rem 0;
    line-height: 1.5;
  }

  .feature-list {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .feature-list li {
    padding: 0.5rem 0;
    color: var(--color-text);
    opacity: 0.8;
    font-size: 0.875rem;
    position: relative;
    padding-left: 1.5rem;
  }

  .feature-list li::before {
    content: "•";
    position: absolute;
    left: 0;
    color: var(--color-accent);
    font-weight: bold;
    font-size: 1.2rem;
  }

  /* Labs Management Styles */
  .labs-toolbar {
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1.25rem;
  }

  .search-group {
    flex: 1;
    min-width: 220px;
    max-width: 320px;
  }

  .search-input {
    width: 100%;
  }

  .btn {
    padding: 0.5rem 1rem;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
    border: none;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.375rem;
  }

  .btn-primary {
    background: var(--color-accent-gradient);
    color: white;
  }

  .btn-primary:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  }

  .btn-secondary {
    background: var(--color-panel);
    color: var(--color-accent);
    border: 1px solid var(--color-accent);
  }

  .btn-secondary:hover:not(:disabled) {
    background: var(--color-accent);
    color: white;
  }

  .btn-danger {
    background: linear-gradient(135deg, #ff6b6b, #ee5a52);
    color: white;
  }

  .btn-danger:hover:not(:disabled) {
    background: linear-gradient(135deg, #ff5252, #d32f2f);
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(255, 107, 107, 0.3);
  }

  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none !important;
  }

  .alert {
    padding: 0.75rem 1rem;
    border-radius: 10px;
    font-size: 0.9rem;
    margin-bottom: 1rem;
  }

  .alert-success {
    background: rgba(76, 175, 80, 0.15);
    color: #2e7d32;
  }

  .alert-error {
    background: rgba(229, 115, 115, 0.15);
    color: #e53935;
  }

  .loading-state,
  .empty-state {
    padding: 2rem 1rem;
    text-align: center;
    color: var(--color-text);
    opacity: 0.75;
  }

  .spinner {
    width: 32px;
    height: 32px;
    border: 3px solid rgba(255, 255, 255, 0.3);
    border-top-color: var(--color-accent);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin: 0 auto 1rem;
  }

  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }

  /* Sessions Table Styles (matching Sessions.svelte) */
  .sessions-table-container {
    background: var(--color-card);
    border-radius: 12px;
    border: 1px solid var(--color-border);
    overflow: hidden;
  }

  .sessions-table {
    width: 100%;
    border-collapse: collapse;
  }

  .sessions-table thead {
    background: var(--color-panel);
    border-bottom: 2px solid var(--color-border);
  }

  .sessions-table th {
    padding: 0.65rem 1rem;
    text-align: left;
    font-weight: 600;
    color: var(--color-text);
    font-size: 0.9rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    white-space: nowrap;
  }

  .sessions-table tbody tr {
    border-bottom: 1px solid var(--color-border);
    transition: background 0.2s ease;
  }

  .sessions-table tbody tr:hover {
    background: var(--color-panel);
  }

  .sessions-table tbody tr:last-child {
    border-bottom: none;
  }

  .sessions-table td {
    padding: 0.65rem 1rem;
    color: var(--color-text);
    text-align: left;
  }

  .lab-row {
    cursor: pointer;
  }

  .lab-code {
    font-weight: 600;
    color: var(--color-accent);
    white-space: nowrap;
  }

  .lab-name {
    font-weight: 500;
    white-space: nowrap;
  }

  .lab-contact-person {
    white-space: nowrap;
  }

  .order-number {
    font-weight: 600;
    color: var(--color-accent);
    white-space: nowrap;
  }

  .order-cost {
    font-weight: 600;
    color: var(--color-accent);
    white-space: nowrap;
  }

  .status-filter {
    width: auto;
    min-width: 150px;
    margin-left: 1rem;
  }

  .status-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.375rem 0.75rem;
    border-radius: 6px;
    font-size: 0.875rem;
    font-weight: 500;
    white-space: nowrap;
  }

  .status-icon {
    font-size: 0.75rem;
  }

  .status-draft {
    background: rgba(156, 163, 175, 0.1);
    color: #9ca3af;
  }

  .status-sent {
    background: rgba(59, 130, 246, 0.1);
    color: #3b82f6;
  }

  .status-in-progress {
    background: rgba(251, 191, 36, 0.1);
    color: #fbbf24;
  }

  .status-ready {
    background: rgba(249, 115, 22, 0.1);
    color: #f97316;
  }

  .status-delivered {
    background: rgba(34, 197, 94, 0.1);
    color: #22c55e;
  }

  .status-cancelled {
    background: rgba(239, 68, 68, 0.1);
    color: #ef4444;
  }

  .detail-section {
    margin-top: 1.5rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--color-border);
  }

  .detail-section-title {
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-text);
    margin-bottom: 1rem;
  }

  /* New Order Form Styles */
  .order-form {
    padding: 0;
  }

  /* New Order Form Grid Layout */
  .order-form-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    grid-template-rows: auto auto;
    gap: 1rem;
    padding: 0;
    margin-bottom: 5rem;
  }

  .order-card {
    background: var(--color-card);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: 0.75rem;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    display: flex;
    flex-direction: column;
  }

  .card-title {
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-text);
    margin: 0 0 0.75rem 0;
    padding-bottom: 0.5rem;
    border-bottom: 1px solid var(--color-border);
    text-align: center;
  }

  .card-row {
    display: flex;
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  .card-row .form-group {
    margin-bottom: 0;
  }

  /* Center quadrant inputs in Teeth/Quadrant Selection card */
  .order-card:nth-child(3) .card-row {
    justify-content: center;
    max-width: 100%;
    margin-left: auto;
    margin-right: auto;
  }

  .order-card .form-group:not(.card-row .form-group) {
    margin-bottom: 1rem;
  }

  .order-card .form-group:last-child {
    margin-bottom: 0;
  }

  .width-30 {
    width: 30%;
    flex: 0 0 30%;
  }

  .width-40 {
    width: 40%;
    flex: 0 0 40%;
  }

  .width-50 {
    width:40%;
    flex: 0 0 40%;
  }

  .width-60 {
    width: 60%;
    flex: 0 0 60%;
  }

  .width-70 {
    width: 70%;
    flex: 0 0 70%;
  }

  .status-radio-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 0.75rem;
    margin-top: 0.375rem;
  }

  .order-card .status-radio-grid {
    gap: 0.75rem;
  }

  .form-actions-fixed {
    position: fixed;
    bottom: 2rem;
    right: 2rem;
    display: flex;
    gap: 1rem;
    z-index: 100;
  }

  @media (max-width: 768px) {
    .order-form-grid {
      grid-template-columns: 1fr;
      grid-template-rows: auto;
    }

    .card-row {
      flex-direction: column;
    }

    .card-row .form-group {
      width: 100% !important;
      flex: 1 1 100% !important;
    }

    .status-radio-grid {
      grid-template-columns: 1fr;
    }

    .form-actions-fixed {
      position: relative;
      bottom: auto;
      right: auto;
      margin-top: 2rem;
      justify-content: flex-end;
    }
  }

  .form-section-header {
    margin-top: 2rem;
    margin-bottom: 1.5rem;
  }

  .form-section-header:first-child {
    margin-top: 0;
  }

  .form-section-header h3 {
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--color-text);
    margin: 0 0 0.75rem 0;
  }

  .form-section-header h4 {
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-text);
    margin: 0 0 0.75rem 0;
  }

  .section-divider {
    height: 1px;
    background: var(--color-border);
    margin-top: 0.5rem;
  }

  .form-row {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1.5rem;
    margin-bottom: 1.5rem;
  }

  @media (max-width: 768px) {
    .form-row {
      grid-template-columns: 1fr;
    }
  }

  .field-hint {
    font-size: 0.75rem;
    color: var(--color-text);
    opacity: 0.6;
    margin-top: 0.25rem;
    display: block;
  }

  .patient-search-container {
    position: relative;
    display: flex;
    align-items: center;
  }

  .search-icon {
    position: absolute;
    left: 0.75rem;
    color: var(--color-text);
    opacity: 0.5;
    pointer-events: none;
    z-index: 1;
  }

  .patient-search-container .search-input {
    padding-left: 1rem;
  }

  .dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: var(--color-card);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    z-index: 1000;
    max-height: 200px;
    overflow-y: auto;
    margin-top: 0.25rem;
  }

  .dropdown-item {
    padding: 0.75rem 1rem;
    cursor: pointer;
    transition: background 0.2s ease;
    border-bottom: 1px solid var(--color-border);
  }

  .dropdown-item:last-child {
    border-bottom: none;
  }

  .dropdown-item:hover {
    background: var(--color-panel);
  }

  .dropdown-item.empty {
    color: var(--color-text);
    opacity: 0.6;
    cursor: default;
  }

  .patient-name {
    font-weight: 500;
    color: var(--color-text);
    font-size: 0.9375rem;
  }

  .patient-phone {
    font-size: 0.875rem;
    color: var(--color-text);
    opacity: 0.7;
    margin-top: 0.25rem;
  }

  .selected-patient-badge {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.5rem 0.75rem;
    background: var(--color-panel);
    border: 1px solid var(--color-border);
    border-radius: 6px;
    margin-top: 0.5rem;
    font-size: 0.875rem;
    color: var(--color-text);
  }

  .order-card .selected-patient-badge {
    padding: 0.5rem 0.75rem;
    margin-top: 0.375rem;
    font-size: 0.8125rem;
  }

  .form-group-inline .selected-patient-badge {
    padding: 0.5rem 0.75rem;
  }

  .clear-btn {
    background: transparent;
    border: none;
    color: var(--color-text);
    opacity: 0.6;
    cursor: pointer;
    font-size: 1.25rem;
    line-height: 1;
    padding: 0;
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    transition: all 0.2s ease;
  }

  .clear-btn:hover {
    opacity: 1;
    background: rgba(239, 68, 68, 0.1);
    color: #ef4444;
  }

  .input-with-suffix {
    position: relative;
    display: flex;
    align-items: center;
  }

  .input-suffix {
    position: absolute;
    right: 0.75rem;
    color: var(--color-text);
    opacity: 0.7;
    font-size: 0.875rem;
    pointer-events: none;
  }

  .order-card .input-suffix {
    right: 0.625rem;
    font-size: 0.8125rem;
  }

  .status-radio-group {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    margin-top: 0.5rem;
  }

  .radio-option {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
    padding: 0.5rem 1rem;
    border: 1px solid var(--color-border);
    border-radius: 6px;
    background: var(--color-panel);
    transition: all 0.2s ease;
    font-size: 0.875rem;
  }

  .order-card .radio-option {
    padding: 0.375rem 0.75rem;
    font-size: 0.8125rem;
  }

  .radio-option:hover {
    background: var(--color-card);
    border-color: var(--color-accent);
  }

  .radio-option input[type="radio"] {
    margin: 0;
    cursor: pointer;
    opacity: 0;
    position: absolute;
    width: 0;
    height: 0;
  }

  .radio-option:has(input[type="radio"]:checked) {
    color: white;
  }

  .radio-option:has(input[type="radio"]:checked) span {
    color: white;
  }

  /* Status-specific colors when selected */
  .radio-option[data-status="draft"]:has(input[type="radio"]:checked) {
    background: #6b7280;
    border-color: #6b7280;
  }

  .radio-option[data-status="sent"]:has(input[type="radio"]:checked) {
    background: #3b82f6;
    border-color: #3b82f6;
  }

  .radio-option[data-status="delivered"]:has(input[type="radio"]:checked) {
    background: #22c55e;
    border-color: #22c55e;
  }

  .radio-option[data-status="cancelled"]:has(input[type="radio"]:checked) {
    background: #ef4444;
    border-color: #ef4444;
  }

  .actions-col {
    text-align: right;
    width: 180px;
    white-space: nowrap;
  }

  .icon-btn {
    border: none;
    background: var(--color-panel);
    border-radius: 8px;
    padding: 0.35rem 0.65rem;
    cursor: pointer;
    font-size: 0.95rem;
    transition: all 0.2s ease;
    margin-left: 0.35rem;
  }

  .icon-btn:hover {
    background: var(--color-accent);
    color: #fff;
  }

  .icon-btn.danger:hover {
    background: #e53935;
    color: #fff;
  }

  .pagination {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 1rem;
    margin-top: 2rem;
    padding-top: 1rem;
  }

  .page-btn {
    padding: 0.5rem 1rem;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 6px;
    cursor: pointer;
    font-size: 0.875rem;
    transition: all 0.2s ease;
  }

  .page-btn:hover:not(:disabled) {
    background: var(--color-border);
  }

  .page-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .page-info {
    color: var(--color-text);
    font-weight: 500;
    font-size: 0.875rem;
  }

  /* Modal Styles */
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
    z-index: 3000;
    animation: fadeIn 0.2s ease;
  }

  .modal-content {
    background: var(--color-card);
    border-radius: 12px;
    width: 90%;
    max-width: 600px;
    max-height: 90vh;
    overflow: hidden;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
    display: flex;
    flex-direction: column;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.25rem 1.5rem;
    border-bottom: 1px solid var(--color-border);
    background: var(--color-accent-gradient);
    color: white;
  }

  .modal-header h2 {
    margin: 0;
    font-size: 1.125rem;
    color: white;
  }

  .modal-header .close-btn {
    background: rgba(255, 255, 255, 0.2);
    color: white;
    border: none;
    font-size: 1.5rem;
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

  .modal-header .close-btn:hover {
    background: rgba(255, 255, 255, 0.3);
  }

  .modal-header .close-btn svg {
    width: 18px;
    height: 18px;
  }

  .danger-header {
    background: linear-gradient(135deg, #ff6b6b, #ee5a52) !important;
  }

  /* Form Styles */
  .lab-form {
    padding: 1.5rem;
    overflow-y: auto;
    flex: 1;
    background: var(--color-card);
    color: var(--color-text);
  }

  .form-section {
    margin-bottom: 1.5rem;
  }

  .form-group {
    margin-bottom: 1.5rem;
    width: 100%;
    box-sizing: border-box;
  }

  .order-card .form-group {
    margin-bottom: 0.75rem;
  }

  .order-card .form-group:last-child {
    margin-bottom: 0;
  }

  .form-group label {
    display: block;
    margin-bottom: 0.375rem;
    font-size: 0.9375rem;
    font-weight: 600;
    color: var(--color-text);
  }

  .order-card .form-group label {
    margin-bottom: 0.25rem;
    font-size: 0.875rem;
  }

  .form-group-inline {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex-wrap: wrap;
  }

  .form-group-inline label {
    display: inline-block;
    margin-bottom: 0;
    white-space: nowrap;
    flex-shrink: 0;
    min-width: 110px;
    text-align: left;
  }

  .form-group-inline .form-input,
  .form-group-inline .patient-search-container,
  .form-group-inline select {
    width: 300px;
    flex: 0 0 300px;
    min-width: 0;
  }

  .form-group-inline .error-message,
  .form-group-inline > div[style*="font-size"] {
    flex-basis: 100%;
    margin-left: 0;
    margin-top: 0.2rem;
  }

  .form-group-inline .selected-patient-badge {
    margin-top: 0.375rem;
    width: 300px;
    flex: 0 0 300px;
    margin-left: calc(110px + 0.75rem);
    margin-right: 0;
    box-sizing: border-box;
  }

  .form-input {
    width: 100%;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 0.6rem 1rem;
    font-size: 1rem;
    margin-bottom: 0.3rem;
    transition: border-color 0.2s, box-shadow 0.2s;
    box-sizing: border-box;
  }

  .order-card .form-input {
    padding: 0.5rem 0.75rem;
    font-size: 0.9375rem;
  }

  .form-input:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  .form-input.error {
    border-color: #e53935;
  }

  .form-input.error:focus {
    box-shadow: 0 0 0 3px rgba(229, 57, 53, 0.1);
  }

  .form-input:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  /* Hide number input spinner buttons */
  input[type="number"]::-webkit-inner-spin-button,
  input[type="number"]::-webkit-outer-spin-button {
    -webkit-appearance: none;
    margin: 0;
  }

  input[type="number"] {
    -moz-appearance: textfield;
  }

  .form-textarea {
    width: 100%;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 0.6rem 1rem;
    font-size: 1rem;
    font-family: inherit;
    resize: vertical;
    min-height: 80px;
    box-sizing: border-box;
    margin-bottom: 0.3rem;
  }

  .order-card .form-textarea {
    padding: 0.5rem 0.75rem;
    font-size: 0.9375rem;
    min-height: 60px;
  }

  /* Description textarea in Teeth/Quadrant Selection card */
  .order-card:nth-child(3) .form-textarea {
    width: 80%;
  }

  /* Financials & Status card - control input widths */
  /* Center all inputs in Financials & Status card */
  .order-card:nth-child(4) .form-group {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-bottom: 1.5rem;
  }

  .order-card:nth-child(4) .form-group:first-of-type {
    margin-bottom: 2.1rem;
  }

  .order-card:nth-child(4) .form-group:last-child {
    margin-bottom: 0;
  }

  .order-card:nth-child(4) .form-group label {
    text-align: center;
    width: 100%;
  }

  .order-card:nth-child(4) .form-input,
  .order-card:nth-child(4) .input-with-suffix {
    width: 80%;
    max-width: 400px;
    margin-left: auto;
    margin-right: auto;
  }

  .order-card:nth-child(4) .status-radio-grid {
    width: 90%;
    max-width: 100%;
    margin-left: auto;
    margin-right: auto;
    box-sizing: border-box;
    padding: 0;
  }

  .order-card:nth-child(4) .status-radio-grid .radio-option {
    padding: 0.375rem 0.15rem;
    min-width: 0;
    width: 100%;
    max-width: 100%;
    box-sizing: border-box;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    justify-content: center;
  }

  .order-card:nth-child(4) .form-textarea {
    width: 80%;
    margin-left: auto;
    margin-right: auto;
  }

  .form-textarea:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  .error-banner {
    background: rgba(229, 57, 53, 0.1);
    border: 1px solid #e53935;
    border-radius: 6px;
    padding: 0.75rem 1rem;
    margin-bottom: 1rem;
  }

  .error-banner .error-message {
    color: #e53935;
    font-size: 0.875rem;
    font-weight: 500;
    margin: 0;
  }

  @media (prefers-color-scheme: dark) {
    .error-banner {
      background: rgba(229, 57, 53, 0.15);
      border-color: #ef5350;
    }

    .error-banner .error-message {
      color: #ef5350;
    }
  }

  .error-message {
    color: #e53935;
    font-size: 0.875rem;
    margin-top: 0.25rem;
    display: block;
  }

  .order-card .error-message {
    font-size: 0.8125rem;
    margin-top: 0.2rem;
  }

  /* Optional Fields Section */
  .optional-fields-section {
    margin-top: 1.5rem;
    margin-bottom: 1rem;
  }

  .optional-fields-header {
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
    font-family: inherit;
    margin-bottom: 0;
  }

  .optional-fields-header:hover {
    background: var(--color-border);
    border-color: var(--color-accent);
  }

  .section-title {
    color: var(--color-text);
  }

  .optional-fields-content {
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
    justify-content: space-between;
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

  .toggle-input:disabled + .toggle-slider {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .form-actions {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
    margin-top: 2rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--color-border);
  }

  /* Lab Details Modal */
  .lab-details-content {
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .detail-row {
    display: flex;
    justify-content: space-between;
    padding: 1rem;
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
    text-align: right;
  }

  /* Confirmation Modal */
  .confirmation-modal {
    background: var(--color-card);
    border-radius: 12px;
    width: 90%;
    max-width: 500px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  }

  .confirmation-modal .modal-content {
    padding: 1.5rem;
    color: var(--color-text);
  }

  .modal-actions {
    display: flex;
    gap: 0.75rem;
    padding: 1.25rem 1.5rem;
    border-top: 1px solid var(--color-border);
  }

  .modal-actions .btn {
    flex: 1;
  }
</style>

