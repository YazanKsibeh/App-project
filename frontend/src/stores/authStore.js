import { writable, get } from 'svelte/store';
import { Login, GetCurrentUser } from '../../wailsjs/go/main/App.js';
import { currentLicenseKey } from './settingsStore.js';

// Session state
export const isAuthenticated = writable(false);
export const currentUser = writable(null);
export const sessionToken = writable(null);
export const authError = writable(null);
export const authLoading = writable(false);

// Helper function to get current license key
function getLicenseKey() {
    let licenseKey = '';
    try {
        licenseKey = get(currentLicenseKey);
    } catch (e) {
        licenseKey = localStorage.getItem('dentist_license_key') || '';
    }
    return licenseKey;
}

// Login function
export async function login(username, password) {
    authLoading.set(true);
    authError.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license first.');
        }

        const response = await Login(username, password, licenseKey);
        
        if (response.success) {
            // Store session information
            sessionToken.set(response.token);
            currentUser.set(response.user);
            isAuthenticated.set(true);
            
            // Save to localStorage
            localStorage.setItem('dentist_session_token', response.token);
            localStorage.setItem('dentist_user_id', response.user.id.toString());
            localStorage.setItem('dentist_username', response.user.username);
            localStorage.setItem('dentist_user_role', response.user.role);
            
            return { success: true, user: response.user };
        } else {
            authError.set(response.message || 'Login failed');
            return { success: false, message: response.message };
        }
    } catch (err) {
        const errorMessage = err.message || 'Login failed';
        authError.set(errorMessage);
        return { success: false, message: errorMessage };
    } finally {
        authLoading.set(false);
    }
}

// Logout function
export function logout() {
    isAuthenticated.set(false);
    currentUser.set(null);
    sessionToken.set(null);
    authError.set(null);
    
    // Clear localStorage
    localStorage.removeItem('dentist_session_token');
    localStorage.removeItem('dentist_user_id');
    localStorage.removeItem('dentist_username');
    localStorage.removeItem('dentist_user_role');
}

// Check if user is logged in (on app startup)
export function checkAuth() {
    const token = localStorage.getItem('dentist_session_token');
    const userId = localStorage.getItem('dentist_user_id');
    const username = localStorage.getItem('dentist_username');
    const role = localStorage.getItem('dentist_user_role');
    
    if (token && userId && username) {
        // Restore session
        sessionToken.set(token);
        currentUser.set({
            id: parseInt(userId),
            username: username,
            role: role || 'Dentist'
        });
        isAuthenticated.set(true);
        return true;
    }
    
    return false;
}

// Check if current user is admin
export function isAdmin() {
    const user = get(currentUser);
    return user && user.role === 'Admin';
}

