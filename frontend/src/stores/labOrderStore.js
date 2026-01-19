import { writable, derived, get } from 'svelte/store';
import { 
    GetLabOrdersPaginated,
    GetLabOrder,
    CreateLabOrder
} from '../../wailsjs/go/main/App.js';
import { currentLicenseKey } from './settingsStore.js';

// Store state
export const labOrders = writable([]);
export const labOrdersLoading = writable(false);
export const labOrdersError = writable(null);
export const labOrderSearch = writable('');
export const labOrderStatusFilter = writable('all');

// Pagination state
export const labOrdersCurrentPage = writable(1);
export const labOrdersTotalPages = writable(1);
export const labOrdersTotalCount = writable(0);
export const labOrdersPageSize = writable(20);

// Helper to get license key
function getLicenseKey() {
    try {
        const key = get(currentLicenseKey);
        if (key) {
            return key;
        }
    } catch (e) {
        // Ignore
    }
    const key = localStorage.getItem('dentist_license_key') || '';
    return key;
}

// Derived store for filtered lab orders (client-side filtering if needed)
export const filteredLabOrders = derived(
    [labOrders, labOrderSearch],
    ([$labOrders, $labOrderSearch]) => {
        const term = $labOrderSearch.trim().toLowerCase();
        if (!term) return $labOrders;
        return $labOrders.filter(order => 
            order.order_number.toLowerCase().includes(term) ||
            (order.patient_name && order.patient_name.toLowerCase().includes(term)) ||
            (order.lab_name && order.lab_name.toLowerCase().includes(term))
        );
    }
);

// Load paginated lab orders from backend
export async function loadLabOrdersPaginated(page = 1, searchOrderNumber = '', searchPatientName = '', searchLabName = '', statusFilter = 'all') {
    labOrdersLoading.set(true);
    labOrdersError.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }

        const pageSize = get(labOrdersPageSize) || 20;
        const response = await GetLabOrdersPaginated(page, pageSize, searchOrderNumber, searchPatientName, searchLabName, statusFilter, licenseKey);
        
        if (response) {
            labOrders.set(response.orders || []);
            labOrdersCurrentPage.set(response.current_page || 1);
            labOrdersTotalPages.set(response.total_pages || 1);
            labOrdersTotalCount.set(response.total_count || 0);
            labOrdersPageSize.set(response.page_size || 20);
        } else {
            labOrders.set([]);
            labOrdersCurrentPage.set(1);
            labOrdersTotalPages.set(1);
            labOrdersTotalCount.set(0);
        }
    } catch (err) {
        console.error('[labOrderStore] loadLabOrdersPaginated error:', err);
        labOrdersError.set(err.message || 'Failed to load lab orders');
        labOrders.set([]);
        labOrdersCurrentPage.set(1);
        labOrdersTotalPages.set(1);
        labOrdersTotalCount.set(0);
    } finally {
        labOrdersLoading.set(false);
    }
}

// Get a specific lab order by id
export async function getLabOrder(id) {
    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        return await GetLabOrder(id, licenseKey);
    } catch (err) {
        console.error('[labOrderStore] getLabOrder error:', err);
        throw err;
    }
}

// Format currency with commas
export function formatCurrency(value = 0) {
    const amount = Number(value) || 0;
    return amount.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
}

// Format date as "Nov 10, 2025"
export function formatOrderDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
        return dateString;
    }
    return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    });
}

// Format date with time as "Nov 10, 2025, 09:00 AM"
export function formatOrderDateTime(dateString) {
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

// Create a new lab order
export async function createLabOrder(form, userID) {
    labOrdersLoading.set(true);
    labOrdersError.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        
        const result = await CreateLabOrder(form, userID, licenseKey);
        
        if (result && result.id && result.order_number) {
            // Reload orders list to show the new order
            await refreshLabOrders();
            return { success: true, orderNumber: result.order_number, id: result.id };
        }
        return { success: false, error: 'Failed to create order' };
    } catch (err) {
        console.error('[labOrderStore] createLabOrder error:', err);
        const errorMessage = err.message || 'Failed to create lab order';
        labOrdersError.set(errorMessage);
        return { success: false, error: errorMessage };
    } finally {
        labOrdersLoading.set(false);
    }
}

// Refresh lab orders on current page
export function refreshLabOrders() {
    const page = get(labOrdersCurrentPage) || 1;
    const searchOrderNumber = get(labOrderSearch) || '';
    const searchPatientName = '';
    const searchLabName = '';
    const statusFilter = get(labOrderStatusFilter) || 'all';
    return loadLabOrdersPaginated(page, searchOrderNumber, searchPatientName, searchLabName, statusFilter);
}

