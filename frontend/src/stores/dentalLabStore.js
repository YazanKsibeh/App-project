import { writable, derived, get } from 'svelte/store';
import { 
    CreateDentalLab, 
    GetDentalLabsPaginated,
    GetDentalLab,
    UpdateDentalLab, 
    DeleteDentalLab 
} from '../../wailsjs/go/main/App.js';
import { currentLicenseKey } from './settingsStore.js';

// Store state
export const dentalLabs = writable([]);
export const dentalLabsLoading = writable(false);
export const dentalLabsError = writable(null);
export const dentalLabsSuccess = writable(null);
export const dentalLabSearch = writable('');

// Pagination state
export const dentalLabsCurrentPage = writable(1);
export const dentalLabsTotalPages = writable(1);
export const dentalLabsTotalCount = writable(0);
export const dentalLabsPageSize = writable(10);

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

// Helper to extract error message from Wails error objects
function extractErrorMessage(err) {
    if (!err) return 'Unknown error';
    
    // If it's already a string, return it
    if (typeof err === 'string') {
        return err;
    }
    
    // Try common error properties
    if (err.message) {
        return err.message;
    }
    
    if (err.error) {
        if (typeof err.error === 'string') {
            return err.error;
        }
        if (err.error.message) {
            return err.error.message;
        }
    }
    
    if (err.data) {
        if (typeof err.data === 'string') {
            return err.data;
        }
        if (err.data.error) {
            return typeof err.data.error === 'string' ? err.data.error : err.data.error.message || err.data.error;
        }
        if (err.data.message) {
            return err.data.message;
        }
    }
    
    // Try toString
    if (err.toString) {
        const errStr = err.toString();
        if (errStr !== '[object Object]') {
            // Try to extract from "Error: message" format
            if (errStr.includes('Error:')) {
                return errStr.split('Error:')[1].trim();
            }
            return errStr;
        }
    }
    
    // Last resort: JSON stringify (but limit length)
    try {
        const jsonStr = JSON.stringify(err);
        if (jsonStr && jsonStr !== '{}' && jsonStr.length < 500) {
            return jsonStr;
        }
    } catch (e) {
        // Ignore JSON errors
    }
    
    return 'Unknown error';
}

// Derived store for filtered dental labs
export const filteredDentalLabs = derived(
    [dentalLabs, dentalLabSearch],
    ([$dentalLabs, $dentalLabSearch]) => {
        const term = $dentalLabSearch.trim().toLowerCase();
        if (!term) return $dentalLabs;
        return $dentalLabs.filter(lab => 
            lab.name.toLowerCase().includes(term) ||
            (lab.contact_person && lab.contact_person.toLowerCase().includes(term)) ||
            (lab.phone_primary && lab.phone_primary.includes(term)) ||
            (lab.phone_secondary && lab.phone_secondary.includes(term))
        );
    }
);

// Load paginated dental labs from backend
export async function loadDentalLabsPaginated(page = 1) {
    dentalLabsLoading.set(true);
    dentalLabsError.set(null);
    dentalLabsSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }

        const pageSize = get(dentalLabsPageSize) || 10;
        const response = await GetDentalLabsPaginated(page, pageSize, licenseKey);
        
        if (response) {
            dentalLabs.set(response.labs || []);
            dentalLabsCurrentPage.set(response.current_page || 1);
            dentalLabsTotalPages.set(response.total_pages || 1);
            dentalLabsTotalCount.set(response.total_count || 0);
            dentalLabsPageSize.set(response.page_size || 10);
        } else {
            dentalLabs.set([]);
            dentalLabsCurrentPage.set(1);
            dentalLabsTotalPages.set(1);
            dentalLabsTotalCount.set(0);
        }
    } catch (err) {
        console.error('[dentalLabStore] loadDentalLabsPaginated error:', err);
        dentalLabsError.set(err.message || 'Failed to load dental labs');
        dentalLabs.set([]);
        dentalLabsCurrentPage.set(1);
        dentalLabsTotalPages.set(1);
        dentalLabsTotalCount.set(0);
    } finally {
        dentalLabsLoading.set(false);
    }
}

// Get a specific dental lab by id
export async function getDentalLab(id) {
    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        return await GetDentalLab(id, licenseKey);
    } catch (err) {
        console.error('[dentalLabStore] getDentalLab error:', err);
        throw err;
    }
}

// Create dental lab
export async function createDentalLab(form) {
    dentalLabsLoading.set(true);
    dentalLabsError.set(null);
    dentalLabsSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        
        await CreateDentalLab(form, licenseKey);
        dentalLabsSuccess.set('Dental lab created successfully');
        
        // Reload current page
        const currentPage = get(dentalLabsCurrentPage);
        if (currentPage && currentPage > 0) {
            await loadDentalLabsPaginated(currentPage);
        } else {
            await loadDentalLabsPaginated(1);
        }
        return true;
    } catch (err) {
        console.error('[dentalLabStore] createDentalLab error:', err);
        const errorMessage = extractErrorMessage(err) || 'Failed to create dental lab';
        console.error('[dentalLabStore] Extracted error message:', errorMessage);
        dentalLabsError.set(errorMessage);
        return false;
    } finally {
        dentalLabsLoading.set(false);
    }
}

// Update dental lab
export async function updateDentalLab(id, form) {
    dentalLabsLoading.set(true);
    dentalLabsError.set(null);
    dentalLabsSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        
        await UpdateDentalLab(id, form, licenseKey);
        dentalLabsSuccess.set('Dental lab updated successfully');
        
        // Reload current page
        const currentPage = get(dentalLabsCurrentPage);
        if (currentPage && currentPage > 0) {
            await loadDentalLabsPaginated(currentPage);
        } else {
            await loadDentalLabsPaginated(1);
        }
        return true;
    } catch (err) {
        console.error('[dentalLabStore] updateDentalLab error:', err);
        const errorMessage = extractErrorMessage(err) || 'Failed to update dental lab';
        dentalLabsError.set(errorMessage);
        return false;
    } finally {
        dentalLabsLoading.set(false);
    }
}

// Delete dental lab
export async function deleteDentalLab(id) {
    dentalLabsLoading.set(true);
    dentalLabsError.set(null);
    dentalLabsSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        await DeleteDentalLab(id, licenseKey);
        dentalLabsSuccess.set('Dental lab deleted successfully');
        
        // Reload current page
        const currentPage = get(dentalLabsCurrentPage);
        if (currentPage && currentPage > 0) {
            await loadDentalLabsPaginated(currentPage);
        } else {
            await loadDentalLabsPaginated(1);
        }
        return true;
    } catch (err) {
        dentalLabsError.set(err.message || 'Failed to delete dental lab');
        console.error('deleteDentalLab error:', err);
        return false;
    } finally {
        dentalLabsLoading.set(false);
    }
}

// Refresh dental labs on current page
export function refreshDentalLabs() {
    const page = get(dentalLabsCurrentPage) || 1;
    return loadDentalLabsPaginated(page);
}

