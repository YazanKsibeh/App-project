import { writable, get } from 'svelte/store';
import {
  CreateSession,
  GetSessions,
  GetSession,
  UpdateSession,
  DeleteSession
} from '../../wailsjs/go/main/App.js';
import { currentLicenseKey } from './settingsStore.js';

const DEFAULT_PAGE_SIZE = 10;

// Store state
export const sessions = writable([]);
export const currentSession = writable(null);
export const sessionsLoading = writable(false);
export const sessionsError = writable(null);
export const sessionsSuccess = writable(null);
export const currentPage = writable(1);
export const totalPages = writable(1);
export const sessionsPageSize = writable(DEFAULT_PAGE_SIZE);
export const sessionsTotalCount = writable(0);

// Filter state
export const sessionFilters = writable({
  patient_id: null,
  status: null,
  dentist_id: null,
  date_from: null,
  date_to: null,
  procedure_ids: []
});

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

// Load sessions from backend (paginated)
export async function loadSessions(page = 1, filters = null) {
    if (!Number.isInteger(page) || page < 1) {
        page = 1;
    }
    console.log('[sessionStore] loadSessions called with page:', page, 'filters:', filters);
    sessionsLoading.set(true);
    sessionsError.set(null);
    sessionsSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        console.log('[sessionStore] License key length:', licenseKey.length);
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        
        // Build filter object for backend
        const filterObj = filters || get(sessionFilters);
        const backendFilters = {
            patient_id: filterObj.patient_id || null,
            status: filterObj.status || null,
            dentist_id: filterObj.dentist_id || null,
            date_from: filterObj.date_from || null,
            date_to: filterObj.date_to || null,
            procedure_ids: filterObj.procedure_ids || []
        };
        
        console.log('[sessionStore] Calling GetSessions backend function with filters:', backendFilters);
        const result = await GetSessions(page, backendFilters, licenseKey);
        console.log('[sessionStore] GetSessions returned:', result);
        console.log('[sessionStore] Result type:', typeof result);
        
        // Handle SessionsResponse object {sessions: [], total_pages: number}
        let sessionsList, total, totalCount, pageSize;
        if (result && typeof result === 'object') {
            if (result.sessions && Array.isArray(result.sessions)) {
                sessionsList = result.sessions;
                total = result.total_pages || 1;
                totalCount = result.total_count || sessionsList.length || 0;
                pageSize = result.page_size || DEFAULT_PAGE_SIZE;
                console.log('[sessionStore] Parsed as SessionsResponse - sessions:', sessionsList.length, 'totalPages:', total);
            } else if (Array.isArray(result)) {
                // Fallback: if it's an array, treat as sessions list
                sessionsList = result;
                total = 1;
                totalCount = sessionsList.length;
                pageSize = sessionsList.length || DEFAULT_PAGE_SIZE;
                console.log('[sessionStore] Parsed as array - sessions:', sessionsList.length);
            } else {
                console.error('[sessionStore] Unexpected result format:', result);
                sessionsList = [];
                total = 1;
                totalCount = 0;
                pageSize = DEFAULT_PAGE_SIZE;
            }
        } else {
            console.error('[sessionStore] Result is null or unexpected type:', result);
            sessionsList = [];
            total = 1;
            totalCount = 0;
            pageSize = DEFAULT_PAGE_SIZE;
        }
        
        const safeSessions = Array.isArray(sessionsList) ? sessionsList : [];
        console.log('[sessionStore] Setting sessions in store, count:', safeSessions.length);
        sessions.set(safeSessions);
        currentPage.set(page);
        totalPages.set(total || 1);
        sessionsTotalCount.set(totalCount || safeSessions.length || 0);
        sessionsPageSize.set(pageSize || DEFAULT_PAGE_SIZE);
    } catch (err) {
        console.error('[sessionStore] loadSessions error:', err);
        console.error('[sessionStore] Error stack:', err.stack);
        sessionsError.set(err.message || 'Failed to load sessions');
        sessions.set([]);
        sessionsPageSize.set(DEFAULT_PAGE_SIZE);
    } finally {
        sessionsLoading.set(false);
        console.log('[sessionStore] loadSessions completed');
    }
}

// Get a specific session by ID
export async function loadSession(id) {
    sessionsLoading.set(true);
    sessionsError.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        
        const session = await GetSession(id, licenseKey);
        currentSession.set(session);
        return session;
    } catch (err) {
        console.error('loadSession error:', err);
        sessionsError.set(err.message || 'Failed to load session');
        return null;
    } finally {
        sessionsLoading.set(false);
    }
}

// Create session
export async function createSession(sessionForm) {
    console.log('[sessionStore] createSession called with form:', sessionForm);
    sessionsLoading.set(true);
    sessionsError.set(null);
    sessionsSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        console.log('[sessionStore] License key length:', licenseKey.length);
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        
        console.log('[sessionStore] Calling CreateSession backend function...');
        const sessionId = await CreateSession(sessionForm, licenseKey);
        console.log('[sessionStore] CreateSession returned session ID:', sessionId);
        sessionsSuccess.set('Session created successfully');
        
        console.log('[sessionStore] Reloading sessions...');
        await loadSessions(get(currentPage), get(sessionFilters));
        console.log('[sessionStore] Sessions reloaded');
        return true;
    } catch (err) {
        console.error('[sessionStore] createSession error:', err);
        console.error('[sessionStore] Error stack:', err.stack);
        sessionsError.set(err.message || 'Failed to create session');
        return false;
    } finally {
        sessionsLoading.set(false);
        console.log('[sessionStore] createSession completed');
    }
}

// Update session
export async function updateSession(session, items) {
    sessionsLoading.set(true);
    sessionsError.set(null);
    sessionsSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        await UpdateSession(session, items, licenseKey);
        sessionsSuccess.set('Session updated successfully');
        await loadSessions(get(currentPage), get(sessionFilters));
        if (currentSession) {
            await loadSession(session.id);
        }
        return true;
    } catch (err) {
        sessionsError.set(err.message || 'Failed to update session');
        console.error('updateSession error:', err);
        return false;
    } finally {
        sessionsLoading.set(false);
    }
}

// Delete session
export async function deleteSession(id) {
    sessionsLoading.set(true);
    sessionsError.set(null);
    sessionsSuccess.set(null);

    try {
        const licenseKey = getLicenseKey();
        if (!licenseKey) {
            throw new Error('License key required. Please validate your license.');
        }
        await DeleteSession(id, licenseKey);
        sessionsSuccess.set('Session deleted successfully');
        await loadSessions(get(currentPage), get(sessionFilters));
        return true;
    } catch (err) {
        sessionsError.set(err.message || 'Failed to delete session');
        console.error('deleteSession error:', err);
        return false;
    } finally {
        sessionsLoading.set(false);
    }
}

// Apply filters and reload sessions
export async function applyFilters(filters) {
    sessionFilters.set(filters);
    await loadSessions(1, filters); // Reset to page 1 when applying filters
}

// Clear all filters
export async function clearAllFilters() {
    const emptyFilters = {
        patient_id: null,
        status: null,
        dentist_id: null,
        date_from: null,
        date_to: null,
        procedure_ids: []
    };
    sessionFilters.set(emptyFilters);
    await loadSessions(1, emptyFilters);
}

// Clear only procedure filters
export async function clearProcedureFilters() {
    const currentFilters = get(sessionFilters);
    const updatedFilters = {
        ...currentFilters,
        procedure_ids: []
    };
    sessionFilters.set(updatedFilters);
    await loadSessions(1, updatedFilters);
}

// Remove a specific filter
export async function removeFilter(filterType, value = null) {
    const currentFilters = get(sessionFilters);
    const updatedFilters = { ...currentFilters };
    
    switch (filterType) {
        case 'patient':
            updatedFilters.patient_id = null;
            break;
        case 'status':
            updatedFilters.status = null;
            break;
        case 'dentist':
            updatedFilters.dentist_id = null;
            break;
        case 'date':
            updatedFilters.date_from = null;
            updatedFilters.date_to = null;
            break;
        case 'procedure':
            if (value !== null) {
                updatedFilters.procedure_ids = updatedFilters.procedure_ids.filter(id => id !== value);
            } else {
                updatedFilters.procedure_ids = [];
            }
            break;
    }
    
    sessionFilters.set(updatedFilters);
    await loadSessions(1, updatedFilters);
}
