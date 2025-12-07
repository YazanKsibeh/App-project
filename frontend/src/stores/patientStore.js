import { writable, get } from 'svelte/store';
import { 
    GetPatients, 
    SearchPatients, 
    AddPatient, 
    UpdatePatient, 
    DeletePatient, 
    DeleteAllPatients,
    OpenPatientFolder
} from '../../wailsjs/go/main/App.js';
import { currentLicenseKey, account } from './settingsStore.js';

// Create stores
export const patients = writable([]);
export const loading = writable(false);
export const error = writable(null);
export const searchTerm = writable('');
export const selectedPatient = writable(null);

// Helper function to get current license key
function getLicenseKey() {
    // Try multiple sources to get the license key
    let licenseKey = '';
    
    // 1. Try currentLicenseKey store
    try {
        licenseKey = get(currentLicenseKey);
        console.log('getLicenseKey - From currentLicenseKey store:', licenseKey ? licenseKey.substring(0, 20) + '...' : 'empty');
    } catch (e) {
        console.error('getLicenseKey - Error getting from currentLicenseKey store:', e);
    }
    
    // 2. If empty, try account store
    if (!licenseKey) {
        try {
            const accountData = get(account);
            licenseKey = accountData.licenseKey || '';
            console.log('getLicenseKey - From account store:', licenseKey ? licenseKey.substring(0, 20) + '...' : 'empty');
        } catch (e) {
            console.error('getLicenseKey - Error getting from account store:', e);
        }
    }
    
    // 3. If still empty, try localStorage directly
    if (!licenseKey) {
        try {
            licenseKey = localStorage.getItem('dentist_license_key') || '';
            console.log('getLicenseKey - From localStorage:', licenseKey ? licenseKey.substring(0, 20) + '...' : 'empty');
        } catch (e) {
            console.error('getLicenseKey - Error getting from localStorage:', e);
        }
    }
    
    console.log('getLicenseKey - Final result length:', licenseKey.length);
    return licenseKey;
}

// Load all patients
export async function loadPatients() {
    loading.set(true);
    error.set(null);
    
    try {
        const licenseKey = getLicenseKey();
        const patientList = await GetPatients(licenseKey);
        patients.set(patientList);
    } catch (err) {
        const errorMessage = err.message || 'Failed to load patients';
        error.set(errorMessage);
        console.error('Error loading patients:', err);
        
        // If license error, clear patients
        if (errorMessage.includes('license')) {
            patients.set([]);
        }
    } finally {
        loading.set(false);
    }
}

// Search patients
export async function searchPatients(term) {
    if (!term.trim()) {
        await loadPatients();
        return;
    }
    
    loading.set(true);
    error.set(null);
    
    try {
        const licenseKey = getLicenseKey();
        const results = await SearchPatients(term, licenseKey);
        patients.set(results);
    } catch (err) {
        const errorMessage = err.message || 'Failed to search patients';
        error.set(errorMessage);
        console.error('Error searching patients:', err);
    } finally {
        loading.set(false);
    }
}

// Add new patient
export async function addPatient(patientData) {
    loading.set(true);
    error.set(null);
    
    try {
        const licenseKey = getLicenseKey();
        const newId = await AddPatient(patientData, licenseKey);
        await loadPatients(); // Reload the list
        return newId;
    } catch (err) {
        const errorMessage = err.message || 'Failed to add patient';
        error.set(errorMessage);
        console.error('Error adding patient:', err);
        throw err;
    } finally {
        loading.set(false);
    }
}

// Update patient
export async function updatePatient(patientData) {
    loading.set(true);
    error.set(null);
    
    try {
        const licenseKey = getLicenseKey();
        await UpdatePatient(patientData, licenseKey);
        await loadPatients(); // Reload the list
    } catch (err) {
        const errorMessage = err.message || 'Failed to update patient';
        error.set(errorMessage);
        console.error('Error updating patient:', err);
        throw err;
    } finally {
        loading.set(false);
    }
}

// Delete patient
export async function deletePatient(id) {
    loading.set(true);
    error.set(null);
    
    try {
        const licenseKey = getLicenseKey();
        await DeletePatient(id, licenseKey);
        await loadPatients(); // Reload the list
    } catch (err) {
        const errorMessage = err.message || 'Failed to delete patient';
        error.set(errorMessage);
        console.error('Error deleting patient:', err);
        throw err;
    } finally {
        loading.set(false);
    }
}

// Delete all patients
export async function deleteAllPatients() {
    loading.set(true);
    error.set(null);
    try {
        const licenseKey = getLicenseKey();
        await DeleteAllPatients(licenseKey);
        await loadPatients(); // Reload the list
    } catch (err) {
        const errorMessage = err.message || 'Failed to delete all patients';
        error.set(errorMessage);
        console.error('Error deleting all patients:', err);
        throw err;
    } finally {
        loading.set(false);
    }
}

// Debug function to check license status
export function debugLicenseStatus() {
    console.log('=== LICENSE DEBUG INFO ===');
    
    // Check localStorage
    const localStorageKey = localStorage.getItem('dentist_license_key');
    const localStorageExpiry = localStorage.getItem('dentist_license_expiry');
    console.log('localStorage - Key:', localStorageKey ? localStorageKey.substring(0, 20) + '...' : 'empty');
    console.log('localStorage - Expiry:', localStorageExpiry);
    
    // Check stores
    try {
        const currentKey = get(currentLicenseKey);
        const accountData = get(account);
        console.log('currentLicenseKey store:', currentKey ? currentKey.substring(0, 20) + '...' : 'empty');
        console.log('account store - Key:', accountData.licenseKey ? accountData.licenseKey.substring(0, 20) + '...' : 'empty');
        console.log('account store - Expiry:', accountData.licenseExpiry);
    } catch (e) {
        console.error('Error reading stores:', e);
    }
    
    // Test getLicenseKey function
    const finalKey = getLicenseKey();
    console.log('getLicenseKey() result:', finalKey ? finalKey.substring(0, 20) + '...' : 'empty');
    console.log('========================');
    
    return finalKey;
}

// Open patient folder in file explorer
export async function openPatientFolder(patientId) {
    error.set(null);
    
    try {
        const licenseKey = getLicenseKey();
        console.log('OpenPatientFolder - Patient ID:', patientId);
        console.log('OpenPatientFolder - License key length:', licenseKey ? licenseKey.length : 0);
        console.log('OpenPatientFolder - License key (first 20 chars):', licenseKey ? licenseKey.substring(0, 20) + '...' : 'empty');
        
        if (!licenseKey) {
            throw new Error('No license key found. Please check your license settings.');
        }
        
        await OpenPatientFolder(patientId, licenseKey);
        console.log('OpenPatientFolder - Success: Folder opened for patient', patientId);
    } catch (err) {
        const errorMessage = err.message || 'Failed to open patient folder';
        error.set(errorMessage);
        console.error('Error opening patient folder:', err);
        console.error('Error details - Patient ID:', patientId);
        console.error('Error details - License key available:', !!getLicenseKey());
        throw err;
    }
}

// Test function for debugging - can be called from browser console
export function testOpenPatientFolder(patientId = 1) {
    console.log('=== TESTING OPEN PATIENT FOLDER ===');
    debugLicenseStatus();
    
    return openPatientFolder(patientId)
        .then(() => {
            console.log('✅ OpenPatientFolder SUCCESS for patient', patientId);
        })
        .catch((error) => {
            console.error('❌ OpenPatientFolder FAILED for patient', patientId, ':', error);
            throw error;
        });
} 