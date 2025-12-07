import { writable, get } from 'svelte/store';
import { GetInvoicePayments } from '../../wailsjs/go/main/App.js';
import { currentLicenseKey } from './settingsStore.js';

const PAGE_SIZE = 10;

export const payments = writable([]);
export const paymentsLoading = writable(false);
export const paymentsError = writable(null);
export const paymentsCurrentPage = writable(1);
export const paymentsTotalPages = writable(1);

function getLicenseKey() {
  try {
    const key = get(currentLicenseKey);
    if (key) {
      return key;
    }
  } catch (error) {
    console.warn('[paymentListStore] Failed to read license key from store', error);
  }
  return localStorage.getItem('dentist_license_key') || '';
}

export async function loadPayments(page = 1) {
  paymentsLoading.set(true);
  paymentsError.set(null);

  try {
    const licenseKey = getLicenseKey();
    if (!licenseKey) {
      throw new Error('License key required. Please validate your license.');
    }

    const response = await GetInvoicePayments(page, PAGE_SIZE, licenseKey);
    payments.set(response?.payments || []);
    paymentsCurrentPage.set(response?.current_page || 1);
    paymentsTotalPages.set(response?.total_pages || 1);
  } catch (error) {
    console.error('[paymentListStore] Failed to load payments', error);
    payments.set([]);
    paymentsCurrentPage.set(1);
    paymentsTotalPages.set(1);
    paymentsError.set(error?.message || 'Failed to load payments');
  } finally {
    paymentsLoading.set(false);
  }
}

export function refreshPayments() {
  const page = get(paymentsCurrentPage) || 1;
  return loadPayments(page);
}

export const paymentsPageSize = PAGE_SIZE;

