import { writable, derived, get } from 'svelte/store';
import { 
    CreateWorkType, 
    GetWorkTypesPaginated,
    UpdateWorkType, 
    DeleteWorkType 
} from '../../wailsjs/go/main/App.js';
import { currentLicenseKey } from './settingsStore.js';
import { currentUser } from './authStore.js';

// Store state
export const workTypes = writable([]);
export const workTypesLoading = writable(false);
export const workTypesError = writable(null);
export const workTypesSuccess = writable(null);
export const workTypeSearch = writable('');

// Pagination state
export const workTypesCurrentPage = writable(1);
export const workTypesTotalPages = writable(1);
export const workTypesTotalCount = writable(0);
export const workTypesPageSize = writable(10);

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

// Derived store for filtered work types
export const filteredWorkTypes = derived(
    [workTypes, workTypeSearch],
    ([$workTypes, $workTypeSearch]) => {
        const term = $workTypeSearch.trim().toLowerCase();
        if (!term) return $workTypes;
        return $workTypes.filter(wt => 
            wt.name.toLowerCase().includes(term) ||
            (wt.description && wt.description.toLowerCase().includes(term))
        );
    }
);

// Load paginated work types from backend
export async function loadWorkTypesPaginated(page = 1) {
    workTypesLoading.set(true);
    workTypesError.set(null);
    workTypesSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }

        const pageSize = get(workTypesPageSize) || 10;
        const response = await GetWorkTypesPaginated(page, pageSize, licenseKey);
        
        if (response) {
            workTypes.set(response.work_types || []);
            workTypesCurrentPage.set(response.current_page || 1);
            workTypesTotalPages.set(response.total_pages || 1);
            workTypesTotalCount.set(response.total_count || 0);
            workTypesPageSize.set(response.page_size || 10);
        } else {
            workTypes.set([]);
            workTypesCurrentPage.set(1);
            workTypesTotalPages.set(1);
            workTypesTotalCount.set(0);
        }
    } catch (err) {
        console.error('[workTypeStore] loadWorkTypesPaginated error:', err);
        workTypesError.set(err.message || 'Failed to load work types');
        workTypes.set([]);
        workTypesCurrentPage.set(1);
        workTypesTotalPages.set(1);
        workTypesTotalCount.set(0);
    } finally {
        workTypesLoading.set(false);
    }
}

// Create work type
export async function createWorkType(form) {
    workTypesLoading.set(true);
    workTypesError.set(null);
    workTypesSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        
        const user = get(currentUser);
        const userID = user ? user.id : parseInt(localStorage.getItem('dentist_user_id') || '0');
        
        if (!userID || userID === 0) {
            throw new Error('User not authenticated. Please log in.');
        }
        
        await CreateWorkType(form, userID, licenseKey);
        workTypesSuccess.set('Work type created successfully');
        
        // Reload current page
        const currentPage = get(workTypesCurrentPage);
        if (currentPage && currentPage > 0) {
            await loadWorkTypesPaginated(currentPage);
        } else {
            await loadWorkTypesPaginated(1);
        }
        return true;
    } catch (err) {
        console.error('[workTypeStore] createWorkType error:', err);
        workTypesError.set(err.message || 'Failed to create work type');
        return false;
    } finally {
        workTypesLoading.set(false);
    }
}

// Update work type
export async function updateWorkType(id, form) {
    workTypesLoading.set(true);
    workTypesError.set(null);
    workTypesSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        
        const user = get(currentUser);
        const userID = user ? user.id : parseInt(localStorage.getItem('dentist_user_id') || '0');
        
        if (!userID || userID === 0) {
            throw new Error('User not authenticated. Please log in.');
        }
        
        await UpdateWorkType(id, form, userID, licenseKey);
        workTypesSuccess.set('Work type updated successfully');
        
        // Reload current page
        const currentPage = get(workTypesCurrentPage);
        if (currentPage && currentPage > 0) {
            await loadWorkTypesPaginated(currentPage);
        } else {
            await loadWorkTypesPaginated(1);
        }
        return true;
    } catch (err) {
        workTypesError.set(err.message || 'Failed to update work type');
        console.error('updateWorkType error:', err);
        return false;
    } finally {
        workTypesLoading.set(false);
    }
}

// Delete work type
export async function deleteWorkType(id) {
    workTypesLoading.set(true);
    workTypesError.set(null);
    workTypesSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        await DeleteWorkType(id, licenseKey);
        workTypesSuccess.set('Work type deleted successfully');
        
        // Reload current page
        const currentPage = get(workTypesCurrentPage);
        if (currentPage && currentPage > 0) {
            await loadWorkTypesPaginated(currentPage);
        } else {
            await loadWorkTypesPaginated(1);
        }
        return true;
    } catch (err) {
        workTypesError.set(err.message || 'Failed to delete work type');
        console.error('deleteWorkType error:', err);
        return false;
    } finally {
        workTypesLoading.set(false);
    }
}

// Refresh work types on current page
export function refreshWorkTypes() {
    const page = get(workTypesCurrentPage) || 1;
    return loadWorkTypesPaginated(page);
}

