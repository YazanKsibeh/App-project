import { writable, get } from 'svelte/store';
import { currentLicenseKey } from './settingsStore.js';
import { GetInvoicePaymentDetails, CreateInvoicePayment } from '../../wailsjs/go/main/App.js';

const paymentDetailsStore = writable(null);
const paymentLoading = writable(false);
const paymentError = writable(null);

function getLicenseKey() {
  try {
    const key = get(currentLicenseKey);
    if (key) {
      return key;
    }
  } catch (err) {
    // Fallback to localStorage in case store is unavailable
  }
  return localStorage.getItem('dentist_license_key') || '';
}

export async function loadInvoicePayments(invoiceId) {
  if (!invoiceId) {
    paymentDetailsStore.set(null);
    return { success: false, error: 'Missing invoice id' };
  }

  paymentDetailsStore.set(null);
  paymentLoading.set(true);
  paymentError.set(null);

  try {
    const licenseKey = getLicenseKey();
    if (!licenseKey) {
      throw new Error('License key required. Please validate your license.');
    }

    const details = await GetInvoicePaymentDetails(invoiceId, licenseKey);
    paymentDetailsStore.set(details);
    return { success: true, details };
  } catch (error) {
    console.error('[PaymentStore] loadInvoicePayments error:', error);
    const message = error?.message || 'Failed to load payment details';
    paymentError.set(message);
    return { success: false, error: message };
  } finally {
    paymentLoading.set(false);
  }
}

export async function addInvoicePayment({ invoiceId, amount, paymentDate, note }) {
  if (!invoiceId) {
    return { success: false, error: 'Missing invoice id' };
  }

  try {
    const licenseKey = getLicenseKey();
    if (!licenseKey) {
      throw new Error('License key required. Please validate your license.');
    }

    const details = await CreateInvoicePayment(invoiceId, amount, paymentDate, note || '', licenseKey);
    paymentDetailsStore.set(details);
    return { success: true, details };
  } catch (error) {
    console.error('[PaymentStore] addInvoicePayment error:', error);
    return { success: false, error: error?.message || 'Failed to add payment' };
  }
}

export const paymentDetails = paymentDetailsStore;
export const paymentDetailsLoading = paymentLoading;
export const paymentDetailsError = paymentError;

