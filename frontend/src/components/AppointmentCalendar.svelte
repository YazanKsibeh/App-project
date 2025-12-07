<script>
import FullCalendar from 'svelte-fullcalendar';
import { onMount } from 'svelte';
import { appointments, loadAppointments } from '../stores/appointmentStore.js';
import { patients, loadPatients } from '../stores/patientStore.js';
import dayGridPlugin from '@fullcalendar/daygrid';
import timeGridPlugin from '@fullcalendar/timegrid';
import interactionPlugin from '@fullcalendar/interaction';
import EditAppointmentModal from './EditAppointmentModal.svelte';
import AddAppointmentModal from './AddAppointmentModal.svelte';

let calendarOptions = {
    plugins: [dayGridPlugin, timeGridPlugin, interactionPlugin],
    initialView: 'dayGridMonth',
    headerToolbar: {
        left: 'prev,next today',
        center: 'title',
        right: 'dayGridMonth,timeGridWeek'
    },
    events: [],
    height: 'auto',
    eventDisplay: 'block',
    initialDate: new Date(),
    eventTimeFormat: {
        hour: 'numeric',
        minute: '2-digit',
        meridiem: 'short'
    },
    slotMinTime: '08:00:00',
    slotMaxTime: '20:00:00',
    slotDuration: '00:30:00',
    slotLabelInterval: '01:00',
    slotLabelFormat: {
        hour: 'numeric',
        minute: '2-digit',
        hour12: true
    },
    slotEventOverlap: true,
    dayMaxEvents: true,
    dayMaxEventRows: 4,
    views: {
        timeGridWeek: {
            slotMinTime: '08:00:00',
            slotMaxTime: '20:00:00',
            slotDuration: '00:30:00',
            slotLabelInterval: '01:00',
            slotLabelFormat: {
                hour: 'numeric',
                minute: '2-digit',
                hour12: true
            },
            eventOverlap: true,
            eventMaxStack: 4,
            slotHeight: 80
        }
    }
};

onMount(() => {
    // Load both patients and appointments
    loadPatients();
    loadAppointments();
});

let showEditModal = false;
let selectedAppointment = null;
let showAddModal = false;

function handleEventClick(info) {
    // info.event.id is the appointment id as string
    const apptId = parseInt(info.event.id);
    const appt = $appointments.find(a => a.id === apptId);
    if (appt) {
        selectedAppointment = appt;
        showEditModal = true;
    }
}

// Add eventContent callback for week view to show time and title on one line
function eventContent(arg) {
    // Only show the event title (patient name)
    const titleText = arg.event.title;
    return { html: `<div style="white-space:nowrap;overflow:hidden;text-overflow:ellipsis;">${titleText}</div>` };
}

// Reactive statement to update calendar events when data changes
$: if ($appointments && $patients) {
    const events = $appointments.map(appt => {
        const patient = $patients.find(p => p.id === appt.patient_id);
        let start = appt.datetime;
        if (start && start.length === 16) start = start + ':00';
        let end = start;
        try {
            const startDateObj = new Date(start);
            if (isNaN(startDateObj.getTime())) return null;
            // Calculate end time using duration
            const endDateObj = new Date(startDateObj.getTime() + (appt.duration || 30) * 60000);
            end = endDateObj.toISOString();
            start = startDateObj.toISOString();
        } catch {
            return null;
        }
        return {
            id: appt.id.toString(),
            title: (patient ? patient.name : 'Unknown Patient') + (appt.notes ? ': ' + appt.notes : ''),
            start,
            end,
            allDay: false,
            backgroundColor: '#667eea',
            borderColor: '#667eea',
            textColor: '#ffffff'
        };
    }).filter(event => event !== null);

    calendarOptions = {
        ...calendarOptions,
        events,
        eventClick: handleEventClick,
        eventContent: eventContent
    };
}

// Function to navigate to appointment date
function goToAppointmentDate() {
    if ($appointments && $appointments.length > 0) {
        const firstAppt = $appointments[0];
        let dateStr = firstAppt.datetime;
        if (dateStr && dateStr.length === 16) {
            dateStr = dateStr + ':00';
        }
        const date = new Date(dateStr);
        if (!isNaN(date.getTime())) {
            calendarOptions = {
                ...calendarOptions,
                initialDate: date
            };
        }
    }
}

function closeEditModal() {
    showEditModal = false;
    selectedAppointment = null;
}
</script>

<div class="calendar-container">
    <div class="calendar-header">
        <h2>Appointment Calendar</h2>
        <div class="calendar-stats">
            <span>Appointments: {$appointments ? $appointments.length : 0}</span>
            <span>Patients: {$patients ? $patients.length : 0}</span>
            <span>Current Date: {new Date().toLocaleDateString()}</span>
        </div>
        <button class="add-btn" on:click={() => showAddModal = true}>+ Add Appointment</button>
    </div>
    <FullCalendar options={calendarOptions} />
</div>

{#if showEditModal}
    <EditAppointmentModal appointment={selectedAppointment} on:close={closeEditModal} />
{/if}

{#if showAddModal}
    <AddAppointmentModal on:close={() => showAddModal = false} />
{/if}

<style>
.calendar-container {
    background: white;
    border-radius: 12px;
    padding: 2rem;
    box-shadow: 0 4px 12px rgba(0,0,0,0.08);
    margin: 2rem auto;
    width: 90vw;
    max-width: 1200px;
}

.calendar-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid #eee;
    flex-wrap: wrap;
    gap: 1rem;
}

.calendar-header h2 {
    color: #333;
    margin: 0;
    font-size: 1.5rem;
    font-weight: 600;
}

.calendar-stats {
    display: flex;
    gap: 1rem;
    font-size: 0.9rem;
    color: #666;
}

.calendar-stats span {
    background: #f5f5f5;
    padding: 0.25rem 0.75rem;
    border-radius: 15px;
}

.nav-btn {
    background: #667eea;
    color: white;
    border: none;
    border-radius: 6px;
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
    cursor: pointer;
    transition: background 0.2s;
}

.nav-btn:hover {
    background: #5a6fd8;
}

:global(.fc) {
    color: #222;
    background: white;
}

:global(.fc-toolbar-title),
:global(.fc-button),
:global(.fc-col-header-cell-cushion),
:global(.fc-daygrid-day-number),
:global(.fc-timegrid-slot-label),
:global(.fc-timegrid-axis-cushion),
:global(.fc-event-title),
:global(.fc-event-time) {
    color: #222 !important;
}

:global(.fc-col-header-cell),
:global(.fc-daygrid-day),
:global(.fc-timegrid-slot) {
    background: #fff !important;
}

:global(.fc-event) {
    background: #667eea !important;
    color: #fff !important;
    border: none !important;
    border-radius: 6px !important;
    font-weight: 500;
    cursor: pointer;
}

:global(.fc-event:hover) {
    background: #5a6fd8 !important;
    transform: translateY(-1px);
    box-shadow: 0 2px 8px rgba(0,0,0,0.15);
}

:global(.fc-toolbar) {
    color: #222;
}

:global(.fc-button) {
    background: #667eea !important;
    border-color: #667eea !important;
}

:global(.fc-button:hover) {
    background: #5a6fd8 !important;
    border-color: #5a6fd8 !important;
}

:global(.fc-button-active) {
    background: #4a5bb8 !important;
    border-color: #4a5bb8 !important;
}

.add-btn {
    background: #667eea;
    color: white;
    border: none;
    border-radius: 6px;
    padding: 0.5rem 1.2rem;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.2s;
}

.add-btn:hover {
    background: #5a67d8;
}
</style> 