import { writable, get } from 'svelte/store';
import { GetInvoiceOverview } from '../../wailsjs/go/main/App.js';
import { currentLicenseKey } from './settingsStore.js';

const defaultOverview = {
    today_total: 0,
    today_count: 0,
    week_total: 0,
    week_count: 0,
    month_total: 0,
    month_count: 0
};

export const invoiceOverview = writable(defaultOverview);
export const invoiceOverviewLoading = writable(false);
export const invoiceOverviewError = writable(null);

function getLicenseKey() {
    try {
        const key = get(currentLicenseKey);
        if (key) {
            return key;
        }
    } catch (error) {
        console.warn('[financialsStore] Failed to read license key from store', error);
    }
    return localStorage.getItem('dentist_license_key') || '';
}

export async function loadInvoiceOverview() {
    invoiceOverviewLoading.set(true);
    invoiceOverviewError.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }

        const data = await GetInvoiceOverview(licenseKey);
        invoiceOverview.set(data || defaultOverview);
    } catch (error) {
        console.error('[financialsStore] Failed to load invoice overview', error);
        invoiceOverviewError.set(error.message || 'Failed to load invoice overview');
        invoiceOverview.set(defaultOverview);
    } finally {
        invoiceOverviewLoading.set(false);
    }
}

export function resetInvoiceOverview() {
    invoiceOverview.set(defaultOverview);
    invoiceOverviewError.set(null);
}

