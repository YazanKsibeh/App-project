import { 
    CreateInvoice, 
    GetInvoiceBySession,
    PreviewInvoice
} from '../../wailsjs/go/main/App.js';
import { currentLicenseKey } from './settingsStore.js';
import { get } from 'svelte/store';

// Helper to get license key
function getLicenseKey() {
    try {
        const key = get(currentLicenseKey);
        if (key) {
            return key;
        }
    } catch (e) {
        // Fallback to localStorage
    }
    return localStorage.getItem('dentist_license_key') || '';
}

// Create invoice from session
export async function createInvoice(sessionID) {
    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        
        const invoice = await CreateInvoice(sessionID, licenseKey);
        return { success: true, invoice };
    } catch (error) {
        console.error('Error creating invoice:', error);
        return { success: false, error: error.message || 'Failed to create invoice' };
    }
}

// Get invoice by session ID
export async function getInvoiceBySession(sessionID) {
    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        
        const invoice = await GetInvoiceBySession(sessionID, licenseKey);
        return { success: true, invoice };
    } catch (error) {
        console.error('Error getting invoice:', error);
        return { success: false, error: error.message || 'Failed to get invoice' };
    }
}

// Preview invoice (get preview data without creating)
export async function previewInvoice(sessionID) {
    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        
        const preview = await PreviewInvoice(sessionID, licenseKey);
        return { success: true, preview };
    } catch (error) {
        console.error('Error previewing invoice:', error);
        return { success: false, error: error.message || 'Failed to preview invoice' };
    }
}

