import { writable, derived, get } from 'svelte/store';
import { 
    CreateExpenseCategory, 
    GetExpenseCategories, 
    GetExpenseCategoriesPaginated,
    UpdateExpenseCategory, 
    DeleteExpenseCategory,
    PermanentlyDeleteExpenseCategory 
} from '../../wailsjs/go/main/App.js';
import { currentLicenseKey } from './settingsStore.js';
import { currentUser } from './authStore.js';

// Store state
export const expenseCategories = writable([]);
export const expenseCategoriesLoading = writable(false);
export const expenseCategoriesError = writable(null);
export const expenseCategoriesSuccess = writable(null);
export const expenseCategorySearch = writable('');

// Pagination state
export const expenseCategoriesCurrentPage = writable(1);
export const expenseCategoriesTotalPages = writable(1);
export const expenseCategoriesTotalCount = writable(0);
export const expenseCategoriesPageSize = writable(10);

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

// Derived store for filtered expense categories
export const filteredExpenseCategories = derived(
    [expenseCategories, expenseCategorySearch],
    ([$expenseCategories, $expenseCategorySearch]) => {
        const term = $expenseCategorySearch.trim().toLowerCase();
        if (!term) return $expenseCategories;
        return $expenseCategories.filter(cat => 
            cat.name.toLowerCase().includes(term) ||
            (cat.description && cat.description.toLowerCase().includes(term))
        );
    }
);

// Load all expense categories
export async function loadExpenseCategories() {
    expenseCategoriesLoading.set(true);
    expenseCategoriesError.set(null);
    
    try {
        const licenseKey = getLicenseKey();
        const categories = await GetExpenseCategories(licenseKey);
        expenseCategories.set(categories || []);
        expenseCategoriesLoading.set(false);
    } catch (error) {
        console.error('[expenseCategoryStore] Error loading expense categories:', error);
        expenseCategoriesError.set(error.message || 'Failed to load expense categories');
        expenseCategoriesLoading.set(false);
    }
}

// Load paginated expense categories
export async function loadExpenseCategoriesPaginated(page = 1) {
    expenseCategoriesLoading.set(true);
    expenseCategoriesError.set(null);
    
    try {
        const licenseKey = getLicenseKey();
        const pageSize = get(expenseCategoriesPageSize);
        const response = await GetExpenseCategoriesPaginated(page, pageSize, licenseKey);
        
        if (response) {
            // Handle both possible response structures
            const categories = response.categories || response.Categories || [];
            expenseCategories.set(categories);
            expenseCategoriesCurrentPage.set(response.current_page || response.CurrentPage || 1);
            expenseCategoriesTotalPages.set(response.total_pages || response.TotalPages || 1);
            expenseCategoriesTotalCount.set(response.total_count || response.TotalCount || 0);
        }
        expenseCategoriesLoading.set(false);
    } catch (error) {
        console.error('[expenseCategoryStore] Error loading paginated expense categories:', error);
        expenseCategoriesError.set(error.message || 'Failed to load expense categories');
        expenseCategoriesLoading.set(false);
    }
}

// Create expense category
export async function createExpenseCategory(category) {
    expenseCategoriesLoading.set(true);
    expenseCategoriesError.set(null);
    expenseCategoriesSuccess.set(null);
    
    try {
        const licenseKey = getLicenseKey();
        const user = get(currentUser);
        const userID = user ? user.id : parseInt(localStorage.getItem('dentist_user_id') || '0');
        
        if (!userID || userID === 0) {
            throw new Error('User not authenticated');
        }
        
        const id = await CreateExpenseCategory(category, userID, licenseKey);
        
        // Reload categories
        await loadExpenseCategoriesPaginated(get(expenseCategoriesCurrentPage));
        
        expenseCategoriesSuccess.set('Expense category created successfully');
        expenseCategoriesLoading.set(false);
        return id;
    } catch (error) {
        console.error('[expenseCategoryStore] Error creating expense category:', error);
        expenseCategoriesError.set(error.message || 'Failed to create expense category');
        expenseCategoriesLoading.set(false);
        throw error;
    }
}

// Update expense category
export async function updateExpenseCategory(id, category) {
    expenseCategoriesLoading.set(true);
    expenseCategoriesError.set(null);
    expenseCategoriesSuccess.set(null);
    
    try {
        const licenseKey = getLicenseKey();
        const user = get(currentUser);
        const userID = user ? user.id : parseInt(localStorage.getItem('dentist_user_id') || '0');
        
        if (!userID || userID === 0) {
            throw new Error('User not authenticated');
        }
        
        await UpdateExpenseCategory(id, category, userID, licenseKey);
        
        // Reload categories
        await loadExpenseCategoriesPaginated(get(expenseCategoriesCurrentPage));
        
        expenseCategoriesSuccess.set('Expense category updated successfully');
        expenseCategoriesLoading.set(false);
    } catch (error) {
        console.error('[expenseCategoryStore] Error updating expense category:', error);
        expenseCategoriesError.set(error.message || 'Failed to update expense category');
        expenseCategoriesLoading.set(false);
        throw error;
    }
}

// Delete expense category (soft delete)
export async function deleteExpenseCategory(id) {
    expenseCategoriesLoading.set(true);
    expenseCategoriesError.set(null);
    expenseCategoriesSuccess.set(null);
    
    try {
        const licenseKey = getLicenseKey();
        await DeleteExpenseCategory(id, licenseKey);
        
        // Reload categories
        await loadExpenseCategoriesPaginated(get(expenseCategoriesCurrentPage));
        
        expenseCategoriesSuccess.set('Expense category deleted successfully');
        expenseCategoriesLoading.set(false);
    } catch (error) {
        console.error('[expenseCategoryStore] Error deleting expense category:', error);
        expenseCategoriesError.set(error.message || 'Failed to delete expense category');
        expenseCategoriesLoading.set(false);
        throw error;
    }
}

// Permanently delete expense category (hard delete)
export async function permanentlyDeleteExpenseCategory(id) {
    expenseCategoriesLoading.set(true);
    expenseCategoriesError.set(null);
    expenseCategoriesSuccess.set(null);
    
    try {
        const licenseKey = getLicenseKey();
        await PermanentlyDeleteExpenseCategory(id, licenseKey);
        
        // Reload categories
        await loadExpenseCategoriesPaginated(get(expenseCategoriesCurrentPage));
        
        expenseCategoriesSuccess.set('Expense category permanently deleted');
        expenseCategoriesLoading.set(false);
    } catch (error) {
        console.error('[expenseCategoryStore] Error permanently deleting expense category:', error);
        expenseCategoriesError.set(error.message || 'Failed to permanently delete expense category');
        expenseCategoriesLoading.set(false);
        throw error;
    }
}

// Clear success message
export function clearExpenseCategorySuccess() {
    expenseCategoriesSuccess.set(null);
}

// Clear error message
export function clearExpenseCategoryError() {
    expenseCategoriesError.set(null);
}

