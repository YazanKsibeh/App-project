import { writable, derived, get } from 'svelte/store';
import { 
    CreateProcedure, 
    GetProcedures, 
    GetProceduresPaginated,
    UpdateProcedure, 
    DeleteProcedure 
} from '../../wailsjs/go/main/App.js';
import { currentLicenseKey } from './settingsStore.js';

// Store state
export const procedures = writable([]);
export const proceduresLoading = writable(false);
export const proceduresError = writable(null);
export const proceduresSuccess = writable(null);
export const procedureSearch = writable('');

// Pagination state
export const proceduresCurrentPage = writable(1);
export const proceduresTotalPages = writable(1);
export const proceduresTotalCount = writable(0);
export const proceduresPageSize = writable(10);

// Helper to get license key
function getLicenseKey() {
    try {
        const key = get(currentLicenseKey);
        if (key) {
            // console.log('[procedureStore] getLicenseKey: found in store, length:', key.length);
            return key;
        }
    } catch (e) {
        // console.log('[procedureStore] getLicenseKey: error getting from store:', e);
    }
    const key = localStorage.getItem('dentist_license_key') || '';
    // console.log('[procedureStore] getLicenseKey: from localStorage, length:', key.length);
    return key;
}

// Derived store for filtered procedures
export const filteredProcedures = derived(
    [procedures, procedureSearch],
    ([$procedures, $procedureSearch]) => {
        const term = $procedureSearch.trim().toLowerCase();
        if (!term) return $procedures;
        return $procedures.filter(proc => 
            proc.name.toLowerCase().includes(term) ||
            proc.price.toString().includes(term)
        );
    }
);

// Load procedures from backend
export async function loadProcedures() {
    // console.log('[procedureStore] loadProcedures() called');
    // console.log('[procedureStore] Setting loading state to true');
    proceduresLoading.set(true);
    proceduresError.set(null);
    proceduresSuccess.set(null);

    try {
        // console.log('[procedureStore] Getting license key...');
        const licenseKey = getLicenseKey();
        // console.log('[procedureStore] License key retrieved, length:', licenseKey.length);
        
        if (!licenseKey) {
            // console.error('[procedureStore] No license key found!');
            throw new Error('License key required. Please validate your license.');
        }
        
        // console.log('[procedureStore] Calling GetProcedures backend function...');
        // const startTime = Date.now();
        const list = await GetProcedures(licenseKey);
        // const duration = Date.now() - startTime;
        // console.log('[procedureStore] GetProcedures returned after', duration, 'ms');
        // console.log('[procedureStore] Raw response:', list);
        // console.log('[procedureStore] Response type:', typeof list);
        // console.log('[procedureStore] Is array?', Array.isArray(list));
        
        // Handle null/undefined response - ensure we always have an array
        const proceduresList = Array.isArray(list) ? list : (list || []);
        // console.log('[procedureStore] Processed procedures count:', proceduresList.length);
        // console.log('[procedureStore] Procedures data:', proceduresList);
        
        // console.log('[procedureStore] Setting procedures in store...');
        procedures.set(proceduresList);
        // console.log('[procedureStore] Procedures set successfully');
    } catch (err) {
        console.error('[procedureStore] loadProcedures error:', err);
        // console.error('[procedureStore] Error details:', {
        //     message: err.message,
        //     stack: err.stack,
        //     name: err.name
        // });
        proceduresError.set(err.message || 'Failed to load procedures');
    } finally {
        // console.log('[procedureStore] Setting loading state to false');
        proceduresLoading.set(false);
        // console.log('[procedureStore] loadProcedures() completed');
    }
}

// Create procedure
export async function createProcedure(form) {
    console.log('[procedureStore] createProcedure called with form:', form);
    proceduresLoading.set(true);
    proceduresError.set(null);
    proceduresSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        console.log('[procedureStore] License key length:', licenseKey.length);
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        console.log('[procedureStore] Calling CreateProcedure backend function...');
        const result = await CreateProcedure(form, licenseKey);
        console.log('[procedureStore] CreateProcedure returned:', result);
        proceduresSuccess.set('Procedure created successfully');
        console.log('[procedureStore] Reloading procedures...');
        // Reload current page if using pagination, otherwise load all
        const currentPage = get(proceduresCurrentPage);
        if (currentPage && currentPage > 0) {
            await loadProceduresPaginated(currentPage);
        } else {
            await loadProcedures();
        }
        console.log('[procedureStore] Procedures reloaded');
        return true;
    } catch (err) {
        console.error('[procedureStore] createProcedure error:', err);
        console.error('[procedureStore] Error stack:', err.stack);
        proceduresError.set(err.message || 'Failed to create procedure');
        return false;
    } finally {
        proceduresLoading.set(false);
        console.log('[procedureStore] createProcedure completed');
    }
}

// Update procedure
export async function updateProcedure(procedure) {
    proceduresLoading.set(true);
    proceduresError.set(null);
    proceduresSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        await UpdateProcedure(procedure, licenseKey);
        proceduresSuccess.set('Procedure updated successfully');
        // Reload current page if using pagination, otherwise load all
        const currentPage = get(proceduresCurrentPage);
        if (currentPage && currentPage > 0) {
            await loadProceduresPaginated(currentPage);
        } else {
            await loadProcedures();
        }
        return true;
    } catch (err) {
        proceduresError.set(err.message || 'Failed to update procedure');
        console.error('updateProcedure error:', err);
        return false;
    } finally {
        proceduresLoading.set(false);
    }
}

// Delete procedure
export async function deleteProcedure(id) {
    proceduresLoading.set(true);
    proceduresError.set(null);
    proceduresSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        await DeleteProcedure(id, licenseKey);
        proceduresSuccess.set('Procedure deleted successfully');
        // Reload current page if using pagination, otherwise load all
        const currentPage = get(proceduresCurrentPage);
        if (currentPage && currentPage > 0) {
            await loadProceduresPaginated(currentPage);
        } else {
            await loadProcedures();
        }
        return true;
    } catch (err) {
        proceduresError.set(err.message || 'Failed to delete procedure');
        console.error('deleteProcedure error:', err);
        return false;
    } finally {
        proceduresLoading.set(false);
    }
}

// Load paginated procedures from backend
export async function loadProceduresPaginated(page = 1) {
    proceduresLoading.set(true);
    proceduresError.set(null);
    proceduresSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }

        const pageSize = get(proceduresPageSize) || 10;
        const response = await GetProceduresPaginated(page, pageSize, licenseKey);
        
        if (response) {
            procedures.set(response.procedures || []);
            proceduresCurrentPage.set(response.current_page || 1);
            proceduresTotalPages.set(response.total_pages || 1);
            proceduresTotalCount.set(response.total_count || 0);
            proceduresPageSize.set(response.page_size || 10);
        } else {
            procedures.set([]);
            proceduresCurrentPage.set(1);
            proceduresTotalPages.set(1);
            proceduresTotalCount.set(0);
        }
    } catch (err) {
        console.error('[procedureStore] loadProceduresPaginated error:', err);
        proceduresError.set(err.message || 'Failed to load procedures');
        procedures.set([]);
        proceduresCurrentPage.set(1);
        proceduresTotalPages.set(1);
        proceduresTotalCount.set(0);
    } finally {
        proceduresLoading.set(false);
    }
}

// Refresh procedures on current page
export function refreshProcedures() {
    const page = get(proceduresCurrentPage) || 1;
    return loadProceduresPaginated(page);
}

