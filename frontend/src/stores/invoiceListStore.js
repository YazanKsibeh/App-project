import { writable, get } from 'svelte/store';
import { GetInvoices } from '../../wailsjs/go/main/App.js';
import { currentLicenseKey } from './settingsStore.js';

const PAGE_SIZE = 5;

export const invoices = writable([]);
export const invoicesLoading = writable(false);
export const invoicesError = writable(null);
export const invoicesCurrentPage = writable(1);
export const invoicesTotalPages = writable(1);

function getLicenseKey() {
  try {
    const key = get(currentLicenseKey);
    if (key) {
      return key;
    }
  } catch (error) {
    console.warn('[invoiceListStore] Failed to read license key from store', error);
  }
  return localStorage.getItem('dentist_license_key') || '';
}

export async function loadInvoices(page = 1) {
  invoicesLoading.set(true);
  invoicesError.set(null);

  try {
    const licenseKey = getLicenseKey();
    if (!licenseKey) {
      throw new Error('License key required. Please validate your license.');
    }

    const result = await GetInvoices(page, PAGE_SIZE, licenseKey);
    invoices.set(result?.invoices || []);
    invoicesCurrentPage.set(result?.current_page || 1);
    invoicesTotalPages.set(result?.total_pages || 1);
  } catch (error) {
    console.error('[invoiceListStore] Failed to load invoices', error);
    invoicesError.set(error.message || 'Failed to load invoices');
    invoices.set([]);
    invoicesCurrentPage.set(1);
    invoicesTotalPages.set(1);
  } finally {
    invoicesLoading.set(false);
  }
}

export function refreshInvoices() {
  const page = get(invoicesCurrentPage) || 1;
  return loadInvoices(page);
}

export const invoicesPageSize = PAGE_SIZE;

