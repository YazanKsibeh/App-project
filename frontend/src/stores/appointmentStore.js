import { writable } from 'svelte/store';
import {
    GetAppointments,
    AddAppointment,
    UpdateAppointment,
    DeleteAppointment
} from '../../wailsjs/go/main/App.js';
import { currentLicenseKey } from './settingsStore.js';

export const appointments = writable([]);
export const loadingAppointments = writable(false);
export const appointmentError = writable(null);

// Helper function to get current license key
function getLicenseKey() {
    let licenseKey = '';
    currentLicenseKey.subscribe(key => licenseKey = key)();
    return licenseKey;
}

export async function loadAppointments() {
    loadingAppointments.set(true);
    appointmentError.set(null);
    try {
        const licenseKey = getLicenseKey();
        const data = await GetAppointments(licenseKey);
        appointments.set(data);
    } catch (err) {
        const errorMessage = err.message || 'Failed to load appointments';
        appointmentError.set(errorMessage);
        
        // If license error, clear appointments
        if (errorMessage.includes('license')) {
            appointments.set([]);
        }
    } finally {
        loadingAppointments.set(false);
    }
}

export async function addAppointment(appt) {
    try {
        const licenseKey = getLicenseKey();
        await AddAppointment(appt, licenseKey);
        await loadAppointments();
    } catch (err) {
        appointmentError.set(err.message || 'Failed to add appointment');
    }
}

export async function updateAppointment(appt) {
    try {
        const licenseKey = getLicenseKey();
        await UpdateAppointment(appt, licenseKey);
        await loadAppointments();
    } catch (err) {
        appointmentError.set(err.message || 'Failed to update appointment');
    }
}

export async function deleteAppointment(id) {
    try {
        const licenseKey = getLicenseKey();
        await DeleteAppointment(id, licenseKey);
        await loadAppointments();
    } catch (err) {
        appointmentError.set(err.message || 'Failed to delete appointment');
    }
} 