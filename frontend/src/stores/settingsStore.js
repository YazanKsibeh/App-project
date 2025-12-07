import { writable, derived } from 'svelte/store';
import { ValidateLicense } from '../../wailsjs/go/main/App.js';

// Theme: 'light' or 'dark'
const initialTheme = localStorage.getItem('theme') || 'dark';
export const theme = writable(initialTheme);
theme.subscribe(value => localStorage.setItem('theme', value));

// License management
const initialLicenseKey = localStorage.getItem('dentist_license_key') || '';
const initialLicenseExpiry = localStorage.getItem('dentist_license_expiry') || '';

// Account/license info
export const account = writable({
    email: 'Test_User', // This can be updated when you implement user management
    licenseKey: initialLicenseKey,
    licenseExpiry: initialLicenseExpiry,
});

// License validation status
export const licenseValidationStatus = writable({
    isValid: false,
    isChecking: false,
    lastChecked: null,
    message: 'License not validated'
});

// Current valid license key (used for API calls)
export const currentLicenseKey = writable(initialLicenseKey);

// Derived store: is license valid and not expired?
export const licenseValid = derived(
    [account, licenseValidationStatus], 
    ([$account, $licenseValidationStatus]) => {
        return $licenseValidationStatus.isValid && $account.licenseKey.length > 0;
    }
);

// Function to validate license with backend
export async function validateCurrentLicense() {
    const $account = account.get ? account.get() : getCurrentAccount();
    
    if (!$account.licenseKey) {
        licenseValidationStatus.set({
            isValid: false,
            isChecking: false,
            lastChecked: new Date(),
            message: 'No license key provided'
        });
        return false;
    }

    licenseValidationStatus.update(status => ({ ...status, isChecking: true }));

    try {
        const licenseInfo = await ValidateLicense($account.licenseKey);
        
        licenseValidationStatus.set({
            isValid: licenseInfo.is_valid,
            isChecking: false,
            lastChecked: new Date(),
            message: licenseInfo.message
        });

        // Update expiry date if validation was successful
        if (licenseInfo.is_valid && licenseInfo.expiry_date) {
            account.update(acc => ({ ...acc, licenseExpiry: licenseInfo.expiry_date }));
            localStorage.setItem('dentist_license_expiry', licenseInfo.expiry_date);
        }

        return licenseInfo.is_valid;
    } catch (error) {
        console.error('License validation error:', error);
        licenseValidationStatus.set({
            isValid: false,
            isChecking: false,
            lastChecked: new Date(),
            message: `Validation error: ${error.message}`
        });
        return false;
    }
}

// Function to set and validate a new license
export async function setLicense(key, expiry = null) {
    // Update the account store
    account.update(acc => ({ 
        ...acc, 
        licenseKey: key, 
        licenseExpiry: expiry || acc.licenseExpiry 
    }));

    // Save to localStorage
    localStorage.setItem('dentist_license_key', key);
    if (expiry) {
        localStorage.setItem('dentist_license_expiry', expiry);
    }

    // Validate the new license
    const isValid = await validateCurrentLicense();
    
    // Update current license key if valid
    if (isValid) {
        currentLicenseKey.set(key);
    } else {
        currentLicenseKey.set('');
        // Remove invalid license from localStorage
        localStorage.removeItem('dentist_license_key');
        localStorage.removeItem('dentist_license_expiry');
    }

    return isValid;
}

// Function to clear license
export function clearLicense() {
    account.update(acc => ({ ...acc, licenseKey: '', licenseExpiry: '' }));
    currentLicenseKey.set('');
    localStorage.removeItem('dentist_license_key');
    localStorage.removeItem('dentist_license_expiry');
    licenseValidationStatus.set({
        isValid: false,
        isChecking: false,
        lastChecked: new Date(),
        message: 'License cleared'
    });
}

// Helper function to get current account (for compatibility)
function getCurrentAccount() {
    return {
        email: 'user@example.com',
        licenseKey: localStorage.getItem('dentist_license_key') || '',
        licenseExpiry: localStorage.getItem('dentist_license_expiry') || ''
    };
}

// Auto-validate license on store initialization
if (initialLicenseKey) {
    validateCurrentLicense();
} 