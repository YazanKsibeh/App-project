import { writable, derived, get } from 'svelte/store';
import { 
    CreateColorShade, 
    GetColorShadesPaginated,
    UpdateColorShade, 
    DeleteColorShade 
} from '../../wailsjs/go/main/App.js';
import { currentLicenseKey } from './settingsStore.js';
import { currentUser } from './authStore.js';

// Store state
export const colorShades = writable([]);
export const colorShadesLoading = writable(false);
export const colorShadesError = writable(null);
export const colorShadesSuccess = writable(null);
export const colorShadeSearch = writable('');

// Pagination state
export const colorShadesCurrentPage = writable(1);
export const colorShadesTotalPages = writable(1);
export const colorShadesTotalCount = writable(0);
export const colorShadesPageSize = writable(10);

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

// Derived store for filtered color shades
export const filteredColorShades = derived(
    [colorShades, colorShadeSearch],
    ([$colorShades, $colorShadeSearch]) => {
        const term = $colorShadeSearch.trim().toLowerCase();
        if (!term) return $colorShades;
        return $colorShades.filter(cs => 
            cs.name.toLowerCase().includes(term) ||
            (cs.description && cs.description.toLowerCase().includes(term)) ||
            (cs.hex_color && cs.hex_color.toLowerCase().includes(term))
        );
    }
);

// Load paginated color shades from backend
export async function loadColorShadesPaginated(page = 1) {
    colorShadesLoading.set(true);
    colorShadesError.set(null);
    colorShadesSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }

        const pageSize = get(colorShadesPageSize) || 10;
        const response = await GetColorShadesPaginated(page, pageSize, licenseKey);
        
        if (response) {
            colorShades.set(response.color_shades || []);
            colorShadesCurrentPage.set(response.current_page || 1);
            colorShadesTotalPages.set(response.total_pages || 1);
            colorShadesTotalCount.set(response.total_count || 0);
            colorShadesPageSize.set(response.page_size || 10);
        } else {
            colorShades.set([]);
            colorShadesCurrentPage.set(1);
            colorShadesTotalPages.set(1);
            colorShadesTotalCount.set(0);
        }
    } catch (err) {
        console.error('[colorShadeStore] loadColorShadesPaginated error:', err);
        colorShadesError.set(err.message || 'Failed to load color shades');
        colorShades.set([]);
        colorShadesCurrentPage.set(1);
        colorShadesTotalPages.set(1);
        colorShadesTotalCount.set(0);
    } finally {
        colorShadesLoading.set(false);
    }
}

// Create color shade
export async function createColorShade(form) {
    colorShadesLoading.set(true);
    colorShadesError.set(null);
    colorShadesSuccess.set(null);

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
        
        await CreateColorShade(form, userID, licenseKey);
        colorShadesSuccess.set('Color shade created successfully');
        
        // Reload current page
        const currentPage = get(colorShadesCurrentPage);
        if (currentPage && currentPage > 0) {
            await loadColorShadesPaginated(currentPage);
        } else {
            await loadColorShadesPaginated(1);
        }
        return true;
    } catch (err) {
        console.error('[colorShadeStore] createColorShade error:', err);
        colorShadesError.set(err.message || 'Failed to create color shade');
        return false;
    } finally {
        colorShadesLoading.set(false);
    }
}

// Update color shade
export async function updateColorShade(id, form) {
    colorShadesLoading.set(true);
    colorShadesError.set(null);
    colorShadesSuccess.set(null);

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
        
        await UpdateColorShade(id, form, userID, licenseKey);
        colorShadesSuccess.set('Color shade updated successfully');
        
        // Reload current page
        const currentPage = get(colorShadesCurrentPage);
        if (currentPage && currentPage > 0) {
            await loadColorShadesPaginated(currentPage);
        } else {
            await loadColorShadesPaginated(1);
        }
        return true;
    } catch (err) {
        colorShadesError.set(err.message || 'Failed to update color shade');
        console.error('updateColorShade error:', err);
        return false;
    } finally {
        colorShadesLoading.set(false);
    }
}

// Delete color shade
export async function deleteColorShade(id) {
    colorShadesLoading.set(true);
    colorShadesError.set(null);
    colorShadesSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        await DeleteColorShade(id, licenseKey);
        colorShadesSuccess.set('Color shade deleted successfully');
        
        // Reload current page
        const currentPage = get(colorShadesCurrentPage);
        if (currentPage && currentPage > 0) {
            await loadColorShadesPaginated(currentPage);
        } else {
            await loadColorShadesPaginated(1);
        }
        return true;
    } catch (err) {
        colorShadesError.set(err.message || 'Failed to delete color shade');
        console.error('deleteColorShade error:', err);
        return false;
    } finally {
        colorShadesLoading.set(false);
    }
}

// Refresh color shades on current page
export function refreshColorShades() {
    const page = get(colorShadesCurrentPage) || 1;
    return loadColorShadesPaginated(page);
}

