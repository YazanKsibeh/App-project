<script>
  import { onMount, onDestroy } from 'svelte';
  import { invoiceOverview, invoiceOverviewLoading, invoiceOverviewError, loadInvoiceOverview } from '../stores/financialsStore.js';
  import { invoices as invoiceList, invoicesLoading, invoicesError, invoicesCurrentPage, invoicesTotalPages, loadInvoices, refreshInvoices } from '../stores/invoiceListStore.js';
  import {
    payments as paymentList,
    paymentsLoading,
    paymentsError,
    paymentsCurrentPage,
    paymentsTotalPages,
    loadPayments as loadPaymentList,
    refreshPayments as refreshPaymentList
  } from '../stores/paymentListStore.js';
  import {
    filteredExpenseCategories,
    expenseCategoriesLoading,
    expenseCategoriesError,
    expenseCategoriesSuccess,
    expenseCategorySearch,
    expenseCategoriesCurrentPage,
    expenseCategoriesTotalPages,
    loadExpenseCategoriesPaginated,
    createExpenseCategory,
    updateExpenseCategory,
    deleteExpenseCategory,
    permanentlyDeleteExpenseCategory,
    clearExpenseCategorySuccess,
    clearExpenseCategoryError
  } from '../stores/expenseCategoryStore.js';
  import PaymentModal from './PaymentModal.svelte';

  const financialSections = [
    {
      id: 'income',
      label: 'Income',
      icon: 'income',
      children: [
        {
          id: 'income_overview',
          label: 'Income Overview',
          description: 'Track revenue performance at a glance.'
        },
        {
          id: 'income_invoices',
          label: 'Invoices',
          description: 'Review issued invoices and their statuses.'
        },
        {
          id: 'income_payments',
          label: 'Payments',
          description: 'Upcoming payments module.',
          emptyMessage: 'Payments tracking is coming soon.'
        },
        {
          id: 'income_receivables',
          label: 'Outstanding Receivables',
          description: 'Monitor pending balances.',
          emptyMessage: 'Receivables insights will appear here soon.'
        },
        {
          id: 'income_reports',
          label: 'Revenue Reports',
          description: 'Historical performance and exports.',
          emptyMessage: 'Revenue reports are under construction.'
        }
      ]
    },
    {
      id: 'expenses',
      label: 'Expenses',
      icon: 'expenses',
      children: [
        {
          id: 'expenses_overview',
          label: 'Expenses Overview',
          description: 'Snapshot of outgoing costs.',
          emptyMessage: 'Expense overview will arrive in a future update.'
        },
        {
          id: 'expenses_bills',
          label: 'Bills & Payables',
          description: 'Track outstanding bills.',
          emptyMessage: 'Bills & payables tracking is coming soon.'
        },
        {
          id: 'expenses_categories',
          label: 'Expense Categories',
          description: 'Organize expenses by category.',
          emptyMessage: 'Expense categories management will be added later.'
        },
        {
          id: 'expenses_payments',
          label: 'Payment Tracking',
          description: 'Log outgoing payments.',
          emptyMessage: 'Expense payment tracking is on the roadmap.'
        },
        {
          id: 'expenses_reports',
          label: 'Expense Reports',
          description: 'Analyze spending trends.',
          emptyMessage: 'Expense reports are planned for a future release.'
        }
      ]
    }
  ];

  let expandedSections = [];
  let activeItem = 'income_overview';
  let refreshInterval = null;
  let paymentModalOpen = false;
  let selectedInvoice = null;
  let paymentToast = '';
  let paymentToastTimeout = null;

  // Expense Category modal state
  let showExpenseCategoryModal = false;
  let expenseCategoryModalMode = 'create';
  let expenseCategoryForm = {
    name: '',
    description: '',
    color: '#3498db',
    expense_type: 'operational',
    budget_amount: 0,
    budget_period: 'monthly',
    is_tax_deductible: true,
    cost_center: 'main',
    account_code: '',
    parent_category_id: null,
    is_active: true,
    requires_approval: false,
    approval_threshold: 0,
    reporting_group: '',
    sort_order: 0
  };
  let expenseCategoryFormError = '';
  let expenseCategoryFormLoading = false;
  let showExpenseCategoryDeleteConfirm = false;
  let expenseCategoryToDelete = null;

  onMount(() => {
    loadInvoiceOverview();
    loadInvoices(1);
    loadPaymentList(1);
    refreshInterval = setInterval(() => {
      loadInvoiceOverview();
      refreshInvoices();
      refreshPaymentList();
    }, 60 * 1000);
  });

  onDestroy(() => {
    if (refreshInterval) {
      clearInterval(refreshInterval);
    }
    if (paymentToastTimeout) {
      clearTimeout(paymentToastTimeout);
    }
  });


  function toggleSection(sectionId) {
    console.log('[Financials] toggleSection', sectionId, 'current:', expandedSections);
    if (expandedSections.includes(sectionId)) {
      expandedSections = expandedSections.filter((id) => id !== sectionId);
    } else {
      expandedSections = [...expandedSections, sectionId];
    }
    console.log('[Financials] expandedSections after toggle', expandedSections);
  }

  function handleSelectItem(sectionId, itemId) {
    console.log('[Financials] select item', itemId, 'under', sectionId);
    if (!expandedSections.includes(sectionId)) {
      expandedSections = [...expandedSections, sectionId];
      console.log('[Financials] auto-expanding section for selection', expandedSections);
    }
    activeItem = itemId;
    
    // Load expense categories when section is selected
    if (itemId === 'expenses_categories') {
      loadExpenseCategoriesPaginated(1);
    }
  }

  // Expense Category handlers
  function clearExpenseCategoryFieldErrors() {
    expenseCategoryFormError = '';
  }

  function openCreateExpenseCategoryModal() {
    expenseCategoryModalMode = 'create';
    expenseCategoryForm = {
      name: '',
      description: '',
      color: '#3498db',
      expense_type: 'operational',
      budget_amount: 0,
      budget_period: 'monthly',
      is_tax_deductible: true,
      cost_center: 'main',
      account_code: '',
      parent_category_id: null,
      is_active: true,
      requires_approval: false,
      approval_threshold: 0,
      reporting_group: '',
      sort_order: 0
    };
    clearExpenseCategoryFieldErrors();
    showExpenseCategoryModal = true;
  }

  let expenseCategoryToEdit = null;

  function openEditExpenseCategoryModal(category) {
    expenseCategoryModalMode = 'edit';
    expenseCategoryToEdit = category;
    expenseCategoryForm = {
      name: category.name || '',
      description: category.description || '',
      color: category.color || '#3498db',
      expense_type: category.expense_type || 'operational',
      // Set defaults for fields not shown in form (will be sent to backend)
      budget_amount: 0,
      budget_period: 'monthly',
      is_tax_deductible: true,
      cost_center: 'main',
      account_code: '',
      parent_category_id: null,
      is_active: category.is_active !== undefined ? category.is_active : true,
      requires_approval: false,
      approval_threshold: 0,
      reporting_group: '',
      sort_order: 0
    };
    clearExpenseCategoryFieldErrors();
    showExpenseCategoryModal = true;
  }

  async function handleExpenseCategorySave() {
    clearExpenseCategoryFieldErrors();
    
    // Validation - only validate required fields
    if (!expenseCategoryForm.name || expenseCategoryForm.name.trim() === '') {
      expenseCategoryFormError = 'Category name is required';
      return;
    }

    if (!expenseCategoryForm.expense_type) {
      expenseCategoryFormError = 'Expense type is required';
      return;
    }

    expenseCategoryFormLoading = true;
    
    try {
      // Set defaults for fields not shown in form
      const formData = {
        name: expenseCategoryForm.name.trim(),
        description: expenseCategoryForm.description.trim() || '',
        color: expenseCategoryForm.color || '#3498db',
        expense_type: expenseCategoryForm.expense_type,
        budget_amount: 0,
        budget_period: 'monthly',
        is_tax_deductible: true,
        cost_center: 'main',
        account_code: '',
        parent_category_id: null,
        is_active: expenseCategoryForm.is_active,
        requires_approval: false,
        approval_threshold: 0,
        reporting_group: '',
        sort_order: 0
      };

      if (expenseCategoryModalMode === 'create') {
        await createExpenseCategory(formData);
      } else if (expenseCategoryToEdit) {
        await updateExpenseCategory(expenseCategoryToEdit.id, formData);
      }
      
      showExpenseCategoryModal = false;
      clearExpenseCategoryError();
    } catch (error) {
      expenseCategoryFormError = error.message || 'Failed to save expense category';
    } finally {
      expenseCategoryFormLoading = false;
    }
  }

  function confirmPermanentDeleteExpenseCategory(category) {
    expenseCategoryToDelete = category;
    expenseCategoryFormError = '';
    showExpenseCategoryDeleteConfirm = true;
  }

  async function handleDeleteExpenseCategory() {
    if (!expenseCategoryToDelete) return;
    
    expenseCategoryFormLoading = true;
    expenseCategoryFormError = '';
    
    try {
      await permanentlyDeleteExpenseCategory(expenseCategoryToDelete.id);
      showExpenseCategoryDeleteConfirm = false;
      expenseCategoryToDelete = null;
      clearExpenseCategoryError();
    } catch (error) {
      expenseCategoryFormError = error.message || 'Failed to delete expense category';
      // Don't close modal on error - user needs to see the error message
    } finally {
      expenseCategoryFormLoading = false;
    }
  }

  function formatExpenseType(type) {
    const types = {
      'operational': 'Operational',
      'capital': 'Capital',
      'personnel': 'Personnel',
      'marketing': 'Marketing',
      'administrative': 'Administrative'
    };
    return types[type] || type;
  }

  function previousExpenseCategories() {
    if ($expenseCategoriesCurrentPage > 1) {
      loadExpenseCategoriesPaginated($expenseCategoriesCurrentPage - 1);
    }
  }

  function nextExpenseCategories() {
    if ($expenseCategoriesCurrentPage < $expenseCategoriesTotalPages) {
      loadExpenseCategoriesPaginated($expenseCategoriesCurrentPage + 1);
    }
  }

  // Clear success/error messages after a delay
  $: if ($expenseCategoriesSuccess) {
    setTimeout(() => clearExpenseCategorySuccess(), 3000);
  }
  $: if ($expenseCategoriesError) {
    setTimeout(() => clearExpenseCategoryError(), 5000);
  }

  function handleInvoiceRowKeydown(event, invoice) {
    const actionableKeys = ['Enter', ' '];
    if (actionableKeys.includes(event.key)) {
      event.preventDefault();
      openPaymentModal(invoice);
    }
  }

  function openPaymentModal(invoice) {
    if (!invoice || $invoicesLoading) return;
    selectedInvoice = invoice;
    paymentModalOpen = true;
  }

  function handlePaymentModalClose() {
    paymentModalOpen = false;
    selectedInvoice = null;
  }

  function handlePaymentSuccess(event) {
    paymentToast = 'Payment recorded successfully.';
    const updatedDetails = event?.detail?.details;
    if (selectedInvoice && updatedDetails?.status) {
      selectedInvoice = { ...selectedInvoice, status: updatedDetails.status };
    }
    refreshInvoices();
    loadInvoiceOverview();
    refreshPaymentList();
    if (paymentToastTimeout) {
      clearTimeout(paymentToastTimeout);
    }
    paymentToastTimeout = setTimeout(() => {
      paymentToast = '';
    }, 3000);
  }

  function formatCurrency(value = 0) {
    const amount = Number(value) || 0;
    return amount.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
  }

  function formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
      return dateString;
    }
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }

  function getInvoiceStatusClass(status) {
    if (!status) return 'status-in-progress';
    const normalized = status.toLowerCase();
    if (normalized === 'paid') return 'status-completed';
    if (normalized === 'cancelled') return 'status-cancelled';
    return 'status-in-progress';
  }

  function formatStatusText(status) {
    if (!status) return 'Issued';
    return status.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase());
  }

  function goToInvoicePage(page) {
    if (!$invoicesLoading) {
      loadInvoices(page);
    }
  }

  function previousInvoices() {
    if ($invoicesCurrentPage > 1) {
      goToInvoicePage($invoicesCurrentPage - 1);
    }
  }

  function nextInvoices() {
    if ($invoicesCurrentPage < $invoicesTotalPages) {
      goToInvoicePage($invoicesCurrentPage + 1);
    }
  }

  function goToPaymentsPage(page) {
    if (!$paymentsLoading) {
      loadPaymentList(page);
    }
  }

  function previousPayments() {
    if ($paymentsCurrentPage > 1) {
      goToPaymentsPage($paymentsCurrentPage - 1);
    }
  }

  function nextPayments() {
    if ($paymentsCurrentPage < $paymentsTotalPages) {
      goToPaymentsPage($paymentsCurrentPage + 1);
    }
  }

  function getSectionIcon(sectionId) {
    if (sectionId === 'income') {
      return (
        '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="1" x2="12" y2="23"/><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>'
      );
    }
    return (
      '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="23" x2="12" y2="1"/><path d="M17 19H9.5a3.5 3.5 0 0 1 0-7h5a3.5 3.5 0 0 0 0-7H6"/></svg>'
    );
  }

  $: activeItemMeta = (() => {
    for (const section of financialSections) {
      const child = section.children.find((child) => child.id === activeItem);
      if (child) {
        return { ...child, parent: section };
      }
    }
    return null;
  })();

  $: overviewData = $invoiceOverview || {};
  $: overviewCards = [
    { key: 'today', label: 'Today', amount: overviewData.today_total || 0, count: overviewData.today_count || 0 },
    { key: 'week', label: 'This Week', amount: overviewData.week_total || 0, count: overviewData.week_count || 0 },
    { key: 'month', label: 'This Month', amount: overviewData.month_total || 0, count: overviewData.month_count || 0 }
  ];
</script>

<div class="financials">
  <div class="financials-container">
    <!-- Left Sidebar Navigation -->
    <aside class="sidebar">
      <h2 class="sidebar-title">Financials</h2>
      <nav class="sidebar-nav">
        {#each financialSections as section}
          <div class="section-group">
            <button
              class="section-toggle"
              class:expanded={expandedSections.includes(section.id)}
              on:click={() => toggleSection(section.id)}
              aria-expanded={expandedSections.includes(section.id)}
            >
              <span class="section-icon" aria-hidden="true">{@html getSectionIcon(section.id)}</span>
              <span class="section-label">{section.label}</span>
              <span class="chevron" aria-hidden="true">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class:is-open={expandedSections.includes(section.id)}>
                  <polyline points="6 9 12 15 18 9"></polyline>
                </svg>
              </span>
            </button>
            {#if expandedSections.includes(section.id)}
              <div class="section-children">
                {#each section.children as child}
                  <button
                    class="sub-nav-item"
                    class:active={activeItem === child.id}
                    on:click={() => handleSelectItem(section.id, child.id)}
                  >
                    <span class="sub-indicator"></span>
                    <span>{child.label}</span>
                  </button>
                {/each}
              </div>
            {/if}
          </div>
        {/each}
      </nav>
    </aside>

    <!-- Right Content Panel -->
    <main class="content-panel">
      <div class="section-content">
        <div class="section-header">
          {#if activeItemMeta?.parent}
            <p class="section-breadcrumb">{activeItemMeta.parent.label}</p>
          {/if}
          <h1>{activeItemMeta?.label || 'Financials'}</h1>
          {#if activeItemMeta?.description}
            <p class="section-subtitle">{activeItemMeta.description}</p>
          {/if}
        </div>

        {#if activeItem === 'income_overview'}
          <div class="overview-section">
            <div class="overview-heading">
              <div class="overview-icon">📊</div>
              <div>
                <p class="overview-label">Invoices Overview</p>
                <h2>Revenue snapshot</h2>
              </div>
            </div>

            {#if $invoiceOverviewLoading}
              <div class="overview-loading">
                <div class="spinner"></div>
                <p>Loading overview...</p>
              </div>
            {:else if $invoiceOverviewError}
              <div class="overview-error">
                <p>{$invoiceOverviewError}</p>
                <button class="btn btn-secondary" on:click={loadInvoiceOverview}>Retry</button>
              </div>
            {:else}
              <div class="overview-cards">
                {#each overviewCards as card}
                  <div class="overview-card" data-period={card.key}>
                    <div class="card-label">{card.label}</div>
                    <div class="card-amount">
                      {formatCurrency(card.amount)}
                      <span class="currency">SYP</span>
                    </div>
                    <div class="card-count">{card.count} {card.count === 1 ? 'invoice' : 'invoices'}</div>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {:else if activeItem === 'income_invoices'}
          <div class="invoices-section">
            <div class="section-subheader">
              <div>
                <p class="overview-label">Invoices List</p>
                <h2>Detailed records</h2>
              </div>
            </div>

            {#if $invoicesLoading}
              <div class="table-state loading-state">
                <div class="spinner"></div>
                <p>Loading invoices...</p>
              </div>
            {:else if $invoicesError}
              <div class="table-state error-state">
                <p>{$invoicesError}</p>
                <button class="btn btn-secondary" on:click={() => loadInvoices($invoicesCurrentPage)}>Retry</button>
              </div>
            {:else if !$invoiceList || $invoiceList.length === 0}
              <div class="table-state empty-message">
                <p>No invoices found.</p>
              </div>
            {:else}
              <div class="sessions-table-container invoices-table-container">
                <table class="sessions-table invoices-table">
                  <thead>
                    <tr>
                      <th>Invoice Number</th>
                      <th>Patient Name</th>
                      <th>Session ID</th>
                      <th>Invoice Date</th>
                      <th>Total Amount</th>
                      <th>Status</th>
                    </tr>
                  </thead>
                  <tbody>
                    {#each $invoiceList as invoice}
                      <tr
                        class="invoice-row"
                        role="button"
                        tabindex="0"
                        on:click={() => openPaymentModal(invoice)}
                        on:keydown={(event) => handleInvoiceRowKeydown(event, invoice)}
                      >
                        <td class="invoice-number">{invoice.invoice_number}</td>
                        <td class="patient-name">{invoice.patient_name || 'Unknown'}</td>
                        <td class="session-id">#{invoice.session_id}</td>
                        <td class="date">{formatDate(invoice.invoice_date)}</td>
                        <td class="total-cost">{formatCurrency(invoice.total_amount)} SYP</td>
                        <td class="status">
                          <span class="status-badge {getInvoiceStatusClass(invoice.status)}">
                            <span class="status-icon">●</span>
                            <span class="status-text">{formatStatusText(invoice.status)}</span>
                          </span>
                        </td>
                      </tr>
                    {/each}
                  </tbody>
                </table>
              </div>

              {#if $invoicesTotalPages > 1}
                <div class="pagination">
                  <button class="page-btn" disabled={$invoicesCurrentPage === 1 || $invoicesLoading} on:click={previousInvoices}>
                    Previous
                  </button>
                  <span class="page-info">
                    Page {$invoicesCurrentPage} of {$invoicesTotalPages}
                  </span>
                  <button class="page-btn" disabled={$invoicesCurrentPage >= $invoicesTotalPages || $invoicesLoading} on:click={nextInvoices}>
                    Next
                  </button>
                </div>
              {/if}
            {/if}
          </div>
        {:else if activeItem === 'income_payments'}
          <div class="payments-section">
            <div class="section-subheader">
              <div>
                <p class="overview-label">Payments List</p>
                <h2>Recorded cash payments</h2>
              </div>
            </div>

            {#if $paymentsLoading}
              <div class="table-state loading-state">
                <div class="spinner"></div>
                <p>Loading payments...</p>
              </div>
            {:else if $paymentsError}
              <div class="table-state error-state">
                <p>{$paymentsError}</p>
                <button class="btn btn-secondary" on:click={() => loadPaymentList($paymentsCurrentPage)}>Retry</button>
              </div>
            {:else if !$paymentList || $paymentList.length === 0}
              <div class="table-state empty-message">
                <p>No payments found.</p>
              </div>
            {:else}
              <div class="sessions-table-container payments-table-container">
                <table class="sessions-table payments-table">
                  <thead>
                    <tr>
                      <th>Patient Name</th>
                      <th>Payment Code</th>
                      <th>Payment Amount</th>
                      <th>Invoice Number</th>
                      <th>Invoice Amount</th>
                      <th>Invoice Status</th>
                      <th>Invoice Date</th>
                      <th>Payment Date</th>
                      <th>Note</th>
                    </tr>
                  </thead>
                  <tbody>
                    {#each $paymentList as payment}
                      <tr>
                        <td class="patient-name">{payment.patient_name || 'Unknown'}</td>
                        <td class="payment-code">{payment.payment_code}</td>
                        <td class="payment-amount">{formatCurrency(payment.payment_amount)} SYP</td>
                        <td class="invoice-number">{payment.invoice_number || '—'}</td>
                        <td class="total-cost">{formatCurrency(payment.invoice_amount)} SYP</td>
                        <td class="status">
                          <span class="status-badge {getInvoiceStatusClass(payment.invoice_status)}">
                            <span class="status-icon">●</span>
                            <span class="status-text">{formatStatusText(payment.invoice_status)}</span>
                          </span>
                        </td>
                        <td class="date">{formatDate(payment.invoice_date) || '—'}</td>
                        <td class="date">{formatDate(payment.payment_date) || '—'}</td>
                        <td class="note-column">{payment.note || '—'}</td>
                      </tr>
                    {/each}
                  </tbody>
                </table>
              </div>

              {#if $paymentsTotalPages > 1}
                <div class="pagination">
                  <button class="page-btn" disabled={$paymentsCurrentPage === 1 || $paymentsLoading} on:click={previousPayments}>
                    Previous
                  </button>
                  <span class="page-info">
                    Page {$paymentsCurrentPage} of {$paymentsTotalPages}
                  </span>
                  <button class="page-btn" disabled={$paymentsCurrentPage >= $paymentsTotalPages || $paymentsLoading} on:click={nextPayments}>
                    Next
                  </button>
                </div>
              {/if}
            {/if}
          </div>
        {:else if activeItem === 'expenses_categories'}
          <div class="expense-categories-section">
            <div class="content-cards">
                <div class="action-card">
                  <div class="procedures-toolbar">
                    <div class="search-group">
                      <input
                        type="text"
                        class="form-input search-input"
                        placeholder="Search categories..."
                        bind:value={$expenseCategorySearch}
                      />
                    </div>
                    <button class="btn btn-primary" on:click={openCreateExpenseCategoryModal}>
                      ➕ Create New Category
                    </button>
                  </div>

                  {#if $expenseCategoriesError}
                    <div class="alert alert-error">
                      {$expenseCategoriesError}
                    </div>
                  {/if}
                  {#if $expenseCategoriesSuccess}
                    <div class="alert alert-success">
                      {$expenseCategoriesSuccess}
                    </div>
                  {/if}

                  <div class="procedure-list">
                    {#if $expenseCategoriesLoading}
                      <div class="loading-state">
                        <div class="spinner"></div>
                        <p>Loading categories...</p>
                      </div>
                    {:else if $filteredExpenseCategories.length === 0}
                      <div class="empty-state">
                        <p>No expense categories found.</p>
                      </div>
                    {:else}
                      <table>
                        <thead>
                          <tr>
                            <th>Category Name</th>
                            <th>Description</th>
                            <th>Status</th>
                            <th class="actions-col">Actions</th>
                          </tr>
                        </thead>
                        <tbody>
                          {#each $filteredExpenseCategories as category}
                            <tr>
                              <td>{category.name}</td>
                              <td>{category.description || '-'}</td>
                              <td>
                                <span class="status-badge {category.is_active ? 'active' : 'inactive'}">
                                  <span class="status-icon">●</span>
                                  <span class="status-text">{category.is_active ? 'Active' : 'Inactive'}</span>
                                </span>
                              </td>
                              <td class="actions-col">
                                <button class="icon-btn" on:click={() => openEditExpenseCategoryModal(category)} title="Edit">
                                  ✏️
                                </button>
                                <button class="icon-btn danger" on:click={() => confirmPermanentDeleteExpenseCategory(category)} title="Delete">
                                  🗑️
                                </button>
                              </td>
                            </tr>
                          {/each}
                        </tbody>
                      </table>
                      
                      {#if $expenseCategoriesTotalPages > 1}
                        <div class="pagination">
                          <button 
                            class="page-btn" 
                            disabled={$expenseCategoriesCurrentPage === 1 || $expenseCategoriesLoading} 
                            on:click={previousExpenseCategories}
                          >
                            Previous
                          </button>
                          <span class="page-info">
                            Page {$expenseCategoriesCurrentPage} of {$expenseCategoriesTotalPages}
                          </span>
                          <button 
                            class="page-btn" 
                            disabled={$expenseCategoriesCurrentPage >= $expenseCategoriesTotalPages || $expenseCategoriesLoading} 
                            on:click={nextExpenseCategories}
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
        {:else}
          <div class="placeholder-card">
            <div class="placeholder-icon">🛠️</div>
            <h3>{activeItemMeta?.label}</h3>
            <p>{activeItemMeta?.emptyMessage || 'This area is being prepared.'}</p>
          </div>
        {/if}
      </div>
    </main>
  </div>

  {#if paymentToast}
    <div class="payment-toast">
      {paymentToast}
    </div>
  {/if}

  <PaymentModal
    bind:open={paymentModalOpen}
    invoice={selectedInvoice}
    on:close={handlePaymentModalClose}
    on:paymentSuccess={handlePaymentSuccess}
  />

  <!-- Expense Category Modal -->
  {#if showExpenseCategoryModal}
    <div 
      class="modal-overlay" 
      role="button"
      tabindex="0"
      on:click={() => { clearExpenseCategoryFieldErrors(); showExpenseCategoryModal = false; }}
      on:keydown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          clearExpenseCategoryFieldErrors();
          showExpenseCategoryModal = false;
        }
      }}
    >
      <div class="modal-content expense-category-modal" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
        <div class="modal-header">
          <h3>{expenseCategoryModalMode === 'create' ? 'Create New Expense Category' : 'Edit Expense Category'}</h3>
          <button class="close-btn" on:click={() => { clearExpenseCategoryFieldErrors(); showExpenseCategoryModal = false; }}>×</button>
        </div>
        <div class="modal-body">
          {#if expenseCategoryFormError}
            <div class="alert alert-error">{expenseCategoryFormError}</div>
          {/if}
          
          <div class="form-group">
            <label for="categoryName">Category Name *</label>
            <input 
              type="text" 
              class="form-input {expenseCategoryFormError && !expenseCategoryForm.name ? 'error' : ''}" 
              placeholder="Enter category name"
              id="categoryName"
              bind:value={expenseCategoryForm.name}
              disabled={expenseCategoryFormLoading}
              required
            />
          </div>
          
          <div class="form-group">
            <label for="categoryDescription">Description</label>
            <textarea 
              class="form-textarea" 
              placeholder="Enter description (optional)"
              id="categoryDescription"
              bind:value={expenseCategoryForm.description}
              disabled={expenseCategoryFormLoading}
              rows="3"
            ></textarea>
          </div>
          
          <div class="form-group">
            <label for="expenseType">Expense Type *</label>
            <select 
              class="form-input" 
              id="expenseType"
              bind:value={expenseCategoryForm.expense_type}
              disabled={expenseCategoryFormLoading}
              required
            >
              <option value="operational">Operational</option>
              <option value="capital">Capital</option>
              <option value="personnel">Personnel</option>
              <option value="marketing">Marketing</option>
              <option value="administrative">Administrative</option>
            </select>
          </div>
          
          <div class="form-group">
            <label for="color">Color</label>
            <input 
              type="color" 
              class="form-input color-input" 
              id="color"
              bind:value={expenseCategoryForm.color}
              disabled={expenseCategoryFormLoading}
            />
          </div>
          
          <div class="form-group">
            <label for="status">Status</label>
            <div class="toggle-group">
              <div class="toggle-item">
                <label for="status" class="toggle-label">
                  <span>{expenseCategoryForm.is_active ? 'Active' : 'Inactive'}</span>
                  <div class="toggle-switch">
                    <input
                      type="checkbox"
                      id="status"
                      bind:checked={expenseCategoryForm.is_active}
                      class="toggle-input"
                      disabled={expenseCategoryFormLoading}
                    />
                    <span class="toggle-slider"></span>
                  </div>
                </label>
              </div>
            </div>
          </div>
        </div>
        <div class="form-actions">
          <button type="button" class="btn btn-secondary" on:click={() => { clearExpenseCategoryFieldErrors(); showExpenseCategoryModal = false; }} disabled={expenseCategoryFormLoading}>
            Cancel
          </button>
          <button type="button" class="btn btn-primary" disabled={expenseCategoryFormLoading} on:click={handleExpenseCategorySave}>
            {#if expenseCategoryFormLoading}
              <div class="spinner"></div>
              {expenseCategoryModalMode === 'create' ? 'Creating...' : 'Saving...'}
            {:else}
              {expenseCategoryModalMode === 'create' ? 'Create Category' : 'Save Changes'}
            {/if}
          </button>
        </div>
      </div>
    </div>
  {/if}

  <!-- Delete Expense Category Confirmation -->
  {#if showExpenseCategoryDeleteConfirm}
    <div 
      class="modal-overlay" 
      role="button"
      tabindex="0"
      on:click={() => { if (!expenseCategoryFormError) { showExpenseCategoryDeleteConfirm = false; expenseCategoryToDelete = null; } }}
      on:keydown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          if (!expenseCategoryFormError) {
            showExpenseCategoryDeleteConfirm = false;
          }
        }
      }}
    >
      <div class="confirmation-modal" tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
        <div class="modal-header danger-header">
          <h3>Permanently Delete Expense Category</h3>
        </div>
        <div class="modal-content">
          {#if expenseCategoryFormError}
            <div class="alert alert-error">{expenseCategoryFormError}</div>
          {/if}
          <p>Are you sure you want to <strong>permanently delete</strong> the expense category <strong>{expenseCategoryToDelete?.name}</strong>?</p>
          <p class="warning-text">⚠️ This action <strong>cannot be undone</strong>. The category will be permanently removed from the database.</p>
        </div>
        <div class="modal-actions">
          <button class="btn btn-danger" on:click={handleDeleteExpenseCategory} disabled={expenseCategoryFormLoading || expenseCategoryFormError}>
            {expenseCategoryFormLoading ? 'Deleting...' : 'Permanently Delete'}
          </button>
          <button class="btn btn-secondary" on:click={() => { showExpenseCategoryDeleteConfirm = false; expenseCategoryToDelete = null; expenseCategoryFormError = ''; }} disabled={expenseCategoryFormLoading}>Cancel</button>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .financials {
    height: calc(100vh - 80px);
    background: var(--color-bg);
    color: var(--color-text);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .financials-container {
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
    padding: 0 0.75rem;
    gap: 0.5rem;
  }

  .section-group {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  .section-toggle {
    display: flex;
    align-items: center;
    width: 100%;
    gap: 0.75rem;
    padding: 0.65rem 0.75rem;
    border: none;
    background: transparent;
    color: var(--color-text);
    font-size: 0.95rem;
    font-weight: 600;
    cursor: pointer;
    border-radius: 10px;
    transition: background 0.2s ease, color 0.2s ease;
    text-align: left;
  }

  .section-toggle:hover {
    background: var(--color-panel);
  }

  .section-icon svg {
    flex-shrink: 0;
  }

  .section-label {
    flex: 1;
  }

  .chevron svg {
    transition: transform 0.2s ease;
  }

  .chevron svg.is-open {
    transform: rotate(180deg);
  }

  .section-children {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    padding-left: 2.5rem;
  }

  .sub-nav-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.45rem 0.5rem;
    border: none;
    background: transparent;
    color: var(--color-text);
    font-size: 0.9rem;
    font-weight: 500;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.2s ease, color 0.2s ease;
    text-align: left;
  }

  .sub-nav-item:hover {
    background: var(--color-panel);
  }

  .sub-nav-item.active {
    background: rgba(102, 126, 234, 0.15);
    color: var(--color-text);
  }

  .sub-indicator {
    width: 6px;
    height: 6px;
    border-radius: 999px;
    background: var(--color-border);
    transition: transform 0.2s ease, background 0.2s ease;
  }

  .sub-nav-item.active .sub-indicator {
    background: var(--color-accent);
    transform: scale(1.4);
  }

  /* Content Panel Styles */
  .content-panel {
    flex: 1;
    overflow-y: auto;
    padding: 2rem;
    background: var(--color-bg);
  }

  .section-content {
    max-width: 1000px;
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
    margin-bottom: 1.5rem;
  }

  .section-breadcrumb {
    margin: 0;
    text-transform: uppercase;
    letter-spacing: 1px;
    font-size: 0.75rem;
    color: var(--color-text);
    opacity: 0.6;
  }

  .section-header h1 {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--color-text);
    margin: 0.35rem 0 0.25rem 0;
  }

  .section-subtitle {
    margin: 0;
    color: var(--color-text);
    opacity: 0.7;
  }

  .empty-state {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 400px;
    text-align: center;
  }

  .placeholder-card {
    background: var(--color-card);
    border: 1px solid var(--color-border);
    border-radius: 16px;
    padding: 3rem 2rem;
    text-align: center;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.75rem;
  }

  .placeholder-icon {
    font-size: 2rem;
  }

  .empty-state.muted {
    min-height: 200px;
    opacity: 0.6;
  }

  .empty-state p {
    font-size: 1.125rem;
    color: var(--color-text);
    opacity: 0.6;
    margin: 0;
  }

  /* Overview Cards */
  .overview-section {
    background: var(--color-card);
    border: 1px solid var(--color-border);
    border-radius: 16px;
    padding: 1.75rem;
    box-shadow: var(--color-shadow);
    margin-bottom: 2rem;
  }

  .overview-heading {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin-bottom: 1.5rem;
  }

  .overview-icon {
    width: 48px;
    height: 48px;
    border-radius: 12px;
    background: var(--color-accent);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.5rem;
    color: white;
  }

  .overview-label {
    margin: 0;
    font-size: 0.85rem;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: var(--color-text);
    opacity: 0.6;
  }

  .overview-heading h2 {
    margin: 0.25rem 0 0 0;
    font-size: 1.4rem;
    color: var(--color-text);
  }

  .overview-cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: 1rem;
  }

  .overview-card {
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: 1.25rem;
    background: var(--color-panel);
    transition: transform 0.2s ease, box-shadow 0.2s ease;
  }

  .overview-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.08);
  }

  .card-label {
    font-size: 0.95rem;
    font-weight: 600;
    color: var(--color-text);
    opacity: 0.8;
    margin-bottom: 0.75rem;
  }

  .card-amount {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--color-text);
    margin-bottom: 0.5rem;
  }

  .currency {
    font-size: 0.95rem;
    font-weight: 600;
    margin-left: 0.35rem;
    color: var(--color-text);
    opacity: 0.7;
  }

  .card-count {
    font-size: 0.95rem;
    color: var(--color-text);
    opacity: 0.7;
  }

  .overview-loading,
  .overview-error {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.75rem;
    padding: 1rem;
    border-radius: 12px;
    border: 1px dashed var(--color-border);
    color: var(--color-text);
  }

  .overview-error button {
    margin-left: 0.5rem;
  }

  .spinner {
    width: 20px;
    height: 20px;
    border: 3px solid rgba(255, 255, 255, 0.2);
    border-top-color: var(--color-accent);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .invoices-section {
    margin-top: 2.5rem;
  }

  .section-subheader {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1.25rem;
  }

  .section-subheader h2 {
    margin: 0.2rem 0 0 0;
    font-size: 1.35rem;
    color: var(--color-text);
  }

  .table-state {
    background: var(--color-card);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: 2rem;
    text-align: center;
    color: var(--color-text);
  }

  .table-state p {
    margin: 0;
  }

  .error-state {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    align-items: center;
  }

  .loading-state {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.75rem;
  }

  .empty-message {
    opacity: 0.7;
  }

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
    padding: 1rem;
    text-align: left;
    font-weight: 600;
    color: var(--color-text);
    font-size: 0.9rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .sessions-table tbody tr {
    border-bottom: 1px solid var(--color-border);
    transition: background 0.2s ease, box-shadow 0.2s ease;
  }

  .sessions-table tbody tr:hover {
    background: var(--color-panel);
  }

  .sessions-table tbody tr.invoice-row {
    cursor: pointer;
  }

  .sessions-table tbody tr.invoice-row:hover {
    box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.08);
  }

  .sessions-table tbody tr.invoice-row:focus-visible {
    outline: 2px solid var(--color-accent);
    outline-offset: -2px;
  }

  .sessions-table tbody tr:last-child {
    border-bottom: none;
  }

  .sessions-table td {
    padding: 1rem;
    color: var(--color-text);
  }
  .note-column {
    max-width: 220px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .total-cost {
    font-weight: 600;
    color: var(--color-accent);
  }

  .session-id {
    font-weight: 600;
    color: var(--color-text);
  }

  .status-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.375rem 0.75rem;
    border-radius: 6px;
    font-size: 0.875rem;
    font-weight: 500;
  }

  .status-completed {
    background: rgba(34, 197, 94, 0.1);
    color: #22c55e;
  }

  .status-in-progress {
    background: rgba(251, 191, 36, 0.1);
    color: #fbbf24;
  }

  .status-cancelled {
    background: rgba(239, 68, 68, 0.1);
    color: #ef4444;
  }

  .status-icon {
    font-size: 0.75rem;
  }

  .pagination {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 1rem;
    margin-top: 2rem;
  }

  .page-btn {
    padding: 0.5rem 1rem;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 6px;
    cursor: pointer;
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
  }

  .payment-toast {
    position: fixed;
    bottom: 2rem;
    right: 2rem;
    background: rgba(34, 197, 94, 0.95);
    color: #fff;
    padding: 0.85rem 1.5rem;
    border-radius: 999px;
    box-shadow: 0 15px 35px rgba(34, 197, 94, 0.35);
    font-weight: 600;
    z-index: 2200;
  }

  /* Expense Categories Section Styles */
  .expense-categories-section {
    width: 100%;
  }

  .section-description {
    margin: 0.5rem 0 0 0;
    color: var(--color-text);
    opacity: 0.7;
    font-size: 0.95rem;
  }

  .content-cards {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .action-card {
    background: var(--color-card);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: 1.5rem;
    box-shadow: var(--color-shadow);
  }

  .procedures-toolbar {
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

  .procedure-list {
    margin-top: 1rem;
  }

  .procedure-list table {
    width: 100%;
    border-collapse: collapse;
    border: 1px solid var(--color-border);
    border-radius: 10px;
    overflow: hidden;
  }

  .procedure-list th,
  .procedure-list td {
    padding: 0.75rem 1rem;
    border-bottom: 1px solid var(--color-border);
    text-align: left;
  }

  .procedure-list th {
    background: var(--color-panel);
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--color-text);
  }

  .procedure-list tr:last-child td {
    border-bottom: none;
  }

  .procedure-list tr:hover td {
    background: rgba(0, 0, 0, 0.02);
  }

  .procedure-list .actions-col {
    text-align: right;
    width: 160px;
    white-space: nowrap;
  }

  .procedure-list .actions-col .icon-btn {
    display: inline-block;
    vertical-align: middle;
    margin-left: 0.35rem;
  }

  .procedure-list .actions-col .icon-btn:first-child {
    margin-left: 0;
  }

  .procedure-list .icon-btn {
    border: none;
    background: var(--color-panel);
    border-radius: 8px;
    padding: 0.35rem 0.65rem;
    cursor: pointer;
    font-size: 0.95rem;
    transition: all 0.2s ease;
    margin-left: 0.35rem;
  }

  .procedure-list .icon-btn:hover {
    background: var(--color-accent);
    color: #fff;
  }

  .procedure-list .icon-btn.danger:hover {
    background: var(--color-danger);
    color: #fff;
  }

  .procedure-list .icon-btn.permanent-delete {
    background: rgba(239, 68, 68, 0.1);
    color: #ef4444;
    border: 1px solid rgba(239, 68, 68, 0.3);
  }

  .procedure-list .icon-btn.permanent-delete:hover {
    background: #ef4444;
    color: #fff;
    border-color: #ef4444;
  }

  .expense-type-badge {
    display: inline-block;
    padding: 0.25rem 0.75rem;
    border-radius: 6px;
    font-size: 0.8125rem;
    font-weight: 500;
    text-transform: capitalize;
  }

  .expense-type-badge.expense-type-operational {
    background: rgba(59, 130, 246, 0.1);
    color: #3b82f6;
  }

  .expense-type-badge.expense-type-capital {
    background: rgba(139, 92, 246, 0.1);
    color: #8b5cf6;
  }

  .expense-type-badge.expense-type-personnel {
    background: rgba(236, 72, 153, 0.1);
    color: #ec4899;
  }

  .expense-type-badge.expense-type-marketing {
    background: rgba(251, 146, 60, 0.1);
    color: #fb923c;
  }

  .expense-type-badge.expense-type-administrative {
    background: rgba(107, 114, 128, 0.1);
    color: #6b7280;
  }

  .budget-period-badge {
    display: inline-block;
    padding: 0.25rem 0.75rem;
    border-radius: 6px;
    font-size: 0.8125rem;
    font-weight: 500;
    background: rgba(102, 126, 234, 0.1);
    color: #667eea;
    text-transform: capitalize;
  }

  .warning-text {
    color: #ef4444;
    font-size: 0.9rem;
    margin-top: 0.75rem;
    padding: 0.75rem;
    background: rgba(239, 68, 68, 0.1);
    border-left: 3px solid #ef4444;
    border-radius: 4px;
  }

  .info-text {
    color: var(--color-text);
    opacity: 0.7;
    font-size: 0.9rem;
    margin-top: 0.75rem;
    padding: 0.75rem;
    background: rgba(59, 130, 246, 0.1);
    border-left: 3px solid #3b82f6;
    border-radius: 4px;
  }

  .procedure-list .loading-state,
  .procedure-list .empty-state {
    padding: 2rem 1rem;
    text-align: center;
    color: var(--color-text);
    opacity: 0.75;
  }

  .procedure-list .pagination {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 1rem;
    margin-top: 2rem;
    padding-top: 1rem;
  }

  .procedure-list .page-btn {
    padding: 0.5rem 1rem;
    background: var(--color-panel);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: 6px;
    cursor: pointer;
    font-size: 0.875rem;
    transition: all 0.2s ease;
  }

  .procedure-list .page-btn:hover:not(:disabled) {
    background: var(--color-border);
  }

  .procedure-list .page-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .procedure-list .page-info {
    color: var(--color-text);
    font-weight: 500;
    font-size: 0.875rem;
  }

  .status-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.375rem;
    padding: 0.25rem 0.75rem;
    border-radius: 6px;
    font-size: 0.8125rem;
    font-weight: 500;
  }

  .status-badge.active {
    background: rgba(34, 197, 94, 0.1);
    color: #22c55e;
  }

  .status-badge.inactive {
    background: rgba(107, 114, 128, 0.1);
    color: #6b7280;
  }

  .status-icon {
    font-size: 0.625rem;
  }

  .status-text {
    font-size: 0.8125rem;
  }

  /* Modal styles for expense categories */
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
    max-width: 800px;
    max-height: 90vh;
    overflow: hidden;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
    display: flex;
    flex-direction: column;
  }

  .modal-content.expense-category-modal {
    max-width: 500px;
    padding: 0;
    overflow-x: hidden;
    box-sizing: border-box;
  }

  .expense-category-modal .modal-header {
    padding: 1.5rem 2rem;
    margin: 0;
  }

  .expense-category-modal .modal-body {
    padding: 0 1.5rem 1.5rem 1.5rem;
    overflow-x: hidden;
    box-sizing: border-box;
  }

  .expense-category-modal .form-actions {
    padding: 1.25rem 2rem;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.25rem 1.5rem;
    border-bottom: 1px solid var(--color-border);
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
  }

  .modal-header h3 {
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

  .modal-actions {
    display: flex;
    gap: 0.75rem;
    padding: 1.25rem 1.5rem;
    border-top: 1px solid var(--color-border);
  }

  .modal-actions .btn {
    flex: 1;
  }

  .modal-body {
    padding: 0 1.5rem 1.5rem 1.5rem;
    overflow-y: auto;
    max-height: calc(90vh - 140px);
  }

  .expense-category-form {
    padding: 0;
  }

  .form-actions {
    display: flex;
    gap: 0.75rem;
    padding: 1.25rem 1.5rem;
    border-top: 1px solid var(--color-border);
    justify-content: flex-end;
  }

  .form-actions .btn {
    min-width: 120px;
  }

  .form-actions .spinner {
    display: inline-block;
    width: 14px;
    height: 14px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-top-color: white;
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
    margin-right: 0.5rem;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  /* Advanced Fields Section */
  .advanced-fields-section {
    margin-top: 1.5rem;
    margin-bottom: 1rem;
  }

  .advanced-fields-header {
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
  }

  .advanced-fields-header:hover {
    background: var(--color-border);
    border-color: var(--color-accent);
  }

  .advanced-fields-header .section-title {
    color: var(--color-text);
  }

  .advanced-fields-header .chevron {
    width: 18px;
    height: 18px;
    color: var(--color-text);
    transition: transform 0.3s ease;
  }

  .advanced-fields-header .chevron.expanded {
    transform: rotate(180deg);
  }

  .advanced-fields-content {
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

  .color-input {
    height: 40px;
    cursor: pointer;
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
    font-size: 0.95rem;
  }

  .toggle-status-text {
    margin-left: 0.5rem;
    font-size: 0.875rem;
    color: var(--color-text);
    opacity: 0.8;
  }

  .toggle-switch {
    position: relative;
    width: 44px;
    height: 24px;
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
    border-radius: 24px;
  }

  .toggle-slider:before {
    position: absolute;
    content: "";
    height: 18px;
    width: 18px;
    left: 3px;
    bottom: 3px;
    background-color: white;
    transition: 0.3s;
    border-radius: 50%;
  }

  .toggle-input:checked + .toggle-slider {
    background-color: #667eea;
  }

  .toggle-input:checked + .toggle-slider:before {
    transform: translateX(20px);
  }

  .toggle-input:disabled + .toggle-slider {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .form-group {
    margin-bottom: 1.5rem;
    width: 100%;
    box-sizing: border-box;
  }

  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-text);
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

  .form-input:disabled {
    opacity: 0.6;
    cursor: not-allowed;
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
    transition: border-color 0.2s, box-shadow 0.2s;
  }

  .form-input select {
    box-sizing: border-box;
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

  .form-textarea:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .error-message {
    color: var(--color-danger);
    font-size: 0.95rem;
  }

  .confirmation-modal {
    background: var(--color-card);
    border-radius: 12px;
    width: 90%;
    max-width: 500px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  }

  .danger-header {
    background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  }

  .alert-error {
    background: rgba(239, 68, 68, 0.1);
    color: #ef4444;
    border: 1px solid rgba(239, 68, 68, 0.2);
  }

  .alert-success {
    background: rgba(34, 197, 94, 0.1);
    color: #22c55e;
    border: 1px solid rgba(34, 197, 94, 0.2);
  }

  .field-error {
    color: #ef4444;
    font-size: 0.875rem;
    margin-bottom: 0.5rem;
    margin-top: -0.25rem;
    display: block;
  }
</style>

