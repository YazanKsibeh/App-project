# DentistApp - Comprehensive Functionality Report

## Executive Summary

**DentistApp** is a desktop application designed for dental practice management. Built using **Wails v2** framework, it combines a **Go (Golang)** backend with a **Svelte.js** frontend. The application provides comprehensive patient management, appointment scheduling, payment tracking, and file management capabilities with a sophisticated license-based access control system.

**Technology Stack:**
- **Backend**: Go 1.23.0 with SQLite database
- **Frontend**: Svelte 3.49.0 with Vite build system
- **Framework**: Wails v2 (Desktop application framework)
- **Database**: SQLite with foreign key constraints
- **License System**: Custom validation using Fibonacci sequences and Base64 encoding

**Platform Support**: Windows, macOS, Linux

---

## Module 1: License Management System

### Overview
The application implements a mandatory license validation system that gates all functionality. Users must enter a valid license key before accessing any features.

### Functionalities

#### 1.1 License Gate Screen
- **Purpose**: Blocks application access until valid license is provided
- **Features**:
  - Displays branded screen with DentistApp logo and tagline
  - License key input field with validation
  - Real-time license validation feedback
  - Automatic license check on application startup
  - License key persistence in browser localStorage
  - Error messages for invalid or expired licenses
  - Support for Enter key to submit license

#### 1.2 License Validation Process
- **Algorithm**: Custom license validation using:
  - Fibonacci sequence positions (3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987)
  - Base64 encoding/decoding for date extraction
  - Offset-based character extraction from license key
  - Date parsing and expiration checking

- **Validation Features**:
  - Extracts expiry date from license key
  - Compares expiry date with current date
  - Returns validation status (valid/invalid)
  - Provides expiry date information
  - Shows days remaining/expired messages
  - Warning for licenses expiring within 7 days

#### 1.3 License Service Integration
- **API Protection**: All backend API endpoints require valid license key
- **Methods**:
  - `ValidateLicense(licenseKey)`: Returns detailed license information
  - `IsLicenseValid(licenseKey)`: Returns boolean validation status
  - `checkLicense(licenseKey)`: Internal validation helper used by all API methods

#### 1.4 License Storage
- **LocalStorage**: License key and expiry date stored in browser localStorage
- **Auto-validation**: License automatically validated on application startup
- **Session Management**: License status maintained throughout application session
- **Clear License**: Ability to clear invalid licenses from storage

---

## Module 2: Patient Management

### Overview
Comprehensive patient record management system with full CRUD operations, search functionality, file system integration, and data validation.

### Functionalities

#### 2.1 Patient List View
- **Display**:
  - Grid layout showing patient cards (4 columns)
  - Patient cards display: Name, Phone (formatted as XXX-XXX-XXXX), Age, Gender
  - Pagination (12 patients per page)
  - Sorting by most recent (highest ID first)
  - Empty state message when no patients exist
  - Loading spinner during data fetch
  - Error messages with retry functionality

- **Search Functionality**:
  - Real-time search bar at top of page
  - Searches by patient name (partial match)
  - Searches by phone number (partial match)
  - Case-insensitive search
  - Automatic list refresh when search term cleared
  - Search results displayed in same grid format

- **Actions Available**:
  - **Add Patient**: Button to open add patient modal
  - **View Patient**: Click on patient card to view details
  - **Edit Patient**: Edit button on patient card
  - **Delete Patient**: Delete button with confirmation dialog

#### 2.2 Add Patient
- **Form Fields**:
  - **Name**: Text input (required)
  - **Phone**: Text input (required, validated for uniqueness)
  - **Age**: Number input (required, validated range 6-100 years)
  - **Gender**: Dropdown/Selection (required)
  - **Occupation**: Text input (optional)

- **Validation Rules**:
  - Phone number must be unique across all patients
  - Age must be between 6 and 100 years
  - All required fields must be filled
  - Error messages displayed for validation failures

- **On Success**:
  - Patient record created in database
  - Patient folder created in file system: `patient_data/{patient_id}/{patient_name}/`
  - Patient name sanitized for filesystem compatibility
  - Patient list automatically refreshed
  - New patient ID returned

#### 2.3 Edit Patient
- **Features**:
  - Same form as Add Patient, pre-populated with existing data
  - Same validation rules as Add Patient
  - Phone uniqueness check excludes current patient
  - Updates all patient fields including total_required amount
  - Patient folder structure preserved

#### 2.4 Patient Detail View
- **Display Information**:
  - Full patient information (all fields)
  - Patient's appointments (linked)
  - Patient's payment history
  - Patient balance information (total required, total paid, remaining)
  - Quick actions for appointments and payments

- **Actions**:
  - **Edit Patient**: Modify patient information
  - **Delete Patient**: Remove patient with confirmation
  - **Open Patient Folder**: Opens file explorer to patient's directory
  - **View Appointments**: Navigate to appointments for this patient
  - **View Payments**: Navigate to payments for this patient
  - **Back to List**: Return to patient list view

#### 2.5 Delete Patient
- **Safety Features**:
  - Confirmation dialog before deletion
  - Transaction-based deletion (all-or-nothing)
  - Foreign key cascade deletion:
    - All appointments for patient automatically deleted
    - All payments for patient automatically deleted
  - Patient folder and all contents deleted from file system
  - Verification that related records are deleted

#### 2.6 Delete All Patients
- **Functionality**:
  - Bulk deletion of all patients
  - Transaction-based operation
  - Deletes all appointments
  - Deletes all payments
  - Removes entire patient_data directory
  - Requires license validation

#### 2.7 Patient File System Integration
- **Folder Structure**:
  - Base directory: `patient_data/`
  - Patient folder: `patient_data/{patient_id}/`
  - Named subfolder: `patient_data/{patient_id}/{sanitized_patient_name}/`

- **Name Sanitization**:
  - Removes/replaces unsafe filesystem characters (/ \ : * ? " < > | .)
  - Replaces multiple spaces/dashes with single dash
  - Limits name length to 50 characters
  - Trims leading/trailing spaces
  - Handles empty names (defaults to "Unknown-Patient")

- **Open Patient Folder**:
  - Cross-platform file explorer integration
  - Windows: Uses `explorer` command
  - macOS: Uses `open` command
  - Linux: Uses `xdg-open` command
  - Creates folder if it doesn't exist
  - Non-blocking execution (runs in goroutine)

#### 2.8 Patient Data Model
- **Fields**:
  - **ID**: Auto-incrementing integer (primary key)
  - **Name**: String (required, not null)
  - **Phone**: String (required, not null, unique)
  - **Age**: Integer (required, validated 6-100)
  - **Gender**: String (required, not null)
  - **Occupation**: String (optional)
  - **TotalRequired**: Integer (default 0, for payment tracking)

---

## Module 3: Appointment Management

### Overview
Complete appointment scheduling system with calendar integration, filtering, and full CRUD operations.

### Functionalities

#### 3.1 Appointment List View
- **Display**:
  - Table/list format showing all appointments
  - Pagination (10 appointments per page)
  - Sorted by datetime (most recent first)
  - Shows: Patient name, Date/Time, Duration, Notes
  - Empty state when no appointments exist
  - Loading states during data fetch

- **Filtering**:
  - **Date Filter**: Filter appointments by specific date
  - **Search Filter**: Search by patient name
  - **Combined Filters**: Both filters work together
  - **Clear Filters**: Button to reset all filters
  - Real-time filtering as user types

- **Actions**:
  - **Add Appointment**: Button to open add appointment modal
  - **Edit Appointment**: Click to modify appointment
  - **Delete Appointment**: Delete with confirmation
  - **View Patient**: Navigate to patient detail page

#### 3.2 Add Appointment
- **Form Fields**:
  - **Patient Selection**: 
    - Searchable dropdown/autocomplete
    - Search by patient name
    - Displays patient name and phone
    - Required field
  - **Date & Time**: 
    - Date-time picker (Flatpickr component)
    - RFC3339 format (e.g., "2024-06-01T14:00:00Z")
    - Required field
  - **Duration**: 
    - Number input (in minutes)
    - Minimum value: 1
    - Optional field
  - **Notes**: 
    - Textarea for additional information
    - Optional field

- **Validation**:
  - Patient must be selected
  - Date and time must be provided
  - Duration must be positive if provided

- **On Success**:
  - Appointment saved to database
  - Linked to selected patient via foreign key
  - Appointment list automatically refreshed
  - New appointment ID returned

#### 3.3 Edit Appointment
- **Features**:
  - Pre-populated form with existing appointment data
  - Can change patient, date/time, duration, and notes
  - Same validation rules as Add Appointment
  - Updates existing appointment record

#### 3.4 Delete Appointment
- **Features**:
  - Confirmation dialog before deletion
  - Permanent deletion from database
  - Appointment list automatically refreshed
  - No cascade effect (appointment deletion doesn't affect patient)

#### 3.5 Appointment Data Model
- **Fields**:
  - **ID**: Auto-incrementing integer (primary key)
  - **PatientID**: Integer (foreign key to patients table, required)
  - **DateTime**: String (RFC3339 format, required)
  - **Duration**: Integer (minutes, optional)
  - **Notes**: String (optional)

- **Relationships**:
  - Many-to-one relationship with Patient
  - Cascade delete: If patient deleted, all appointments deleted

---

## Module 4: Calendar View

### Overview
Interactive calendar interface using FullCalendar library for visual appointment scheduling and management.

### Functionalities

#### 4.1 Calendar Display
- **Views Available**:
  - **Month View**: Grid calendar showing all days of month
  - **Week View**: Time-based week view (8:00 AM - 8:00 PM)
  - **Day View**: Single day time grid

- **Calendar Features**:
  - Navigation buttons (Previous, Next, Today)
  - View switcher (Month/Week)
  - Current date highlighted
  - Appointment events displayed as blocks
  - Event colors and styling
  - Time slots: 30-minute intervals
  - Business hours: 8:00 AM to 8:00 PM

#### 4.2 Appointment Events on Calendar
- **Event Display**:
  - Each appointment shown as colored block
  - Event title shows patient name
  - Time displayed for each appointment
  - Duration represented by block height (week view)
  - Multiple appointments on same day shown stacked (month view)
  - Maximum 4 event rows per day (month view)

#### 4.3 Calendar Interactions
- **Click Event**: 
  - Click on appointment to view/edit
  - Opens edit appointment modal
  - Pre-populated with appointment data

- **Add Appointment**:
  - Button to add new appointment
  - Opens add appointment modal
  - Can select patient and set date/time

#### 4.4 Calendar Data Loading
- **Automatic Loading**:
  - Loads all appointments on calendar mount
  - Loads all patients for selection
  - Converts appointment data to calendar event format
  - Updates calendar when appointments change

---

## Module 5: Payment Management

### Overview
Comprehensive payment tracking system with balance calculations, payment history, and financial summaries.

### Functionalities

#### 5.1 Payment Overview Page
- **Display**:
  - List of all patients with payment summaries
  - Shows for each patient:
    - Patient name and phone
    - Total Required amount
    - Total Paid amount
    - Remaining Balance
    - Last Payment Date
  - Pagination (10 patients per page)
  - Search/filter by patient name

- **Actions**:
  - **View Patient Payments**: Click to see detailed payment history
  - **Add Payment**: Quick add payment for patient
  - **Edit Total Required**: Modify total required amount for patient

#### 5.2 Patient Payment Detail Page
- **Display**:
  - Patient information header
  - Financial Summary:
    - Total Required (editable)
    - Total Paid (calculated)
    - Remaining Balance (calculated)
  - Payment History Table:
    - All payments for patient
    - Sorted by date (most recent first)
    - Shows: Date, Amount, Note
    - Actions: Edit, Delete

- **Actions**:
  - **Add Payment**: Add new payment record
  - **Edit Payment**: Modify existing payment
  - **Delete Payment**: Remove payment with confirmation
  - **Edit Total Required**: Update total required amount
  - **Back to Overview**: Return to payment overview

#### 5.3 Add Payment
- **Form Fields**:
  - **Patient**: Pre-selected (if coming from patient detail)
  - **Amount**: Number input (required, positive integer)
  - **Date**: Date picker (required)
  - **Note**: Textarea (optional, for payment description)

- **Validation**:
  - Amount must be positive
  - Date must be provided
  - Patient must be selected

- **On Success**:
  - Payment saved to database
  - Linked to patient via foreign key
  - Balance automatically recalculated
  - Payment list refreshed

#### 5.4 Edit Payment
- **Features**:
  - Pre-populated form with existing payment data
  - Can modify amount, date, and note
  - Same validation as Add Payment
  - Updates existing payment record
  - Balance recalculated after update

#### 5.5 Delete Payment
- **Features**:
  - Confirmation dialog before deletion
  - Permanent deletion from database
  - Balance automatically recalculated
  - Payment list refreshed

#### 5.6 Balance Calculation
- **Automatic Calculation**:
  - **Total Required**: Stored in patient record (editable)
  - **Total Paid**: Sum of all payment amounts for patient
  - **Remaining**: Total Required - Total Paid
  - Real-time updates when payments added/edited/deleted

#### 5.7 Update Total Required
- **Functionality**:
  - Allows clinic to set/update total treatment cost
  - Updates patient's total_required field
  - Balance automatically recalculated
  - Can be done from patient detail or payment page

#### 5.8 Payment Data Model
- **Fields**:
  - **ID**: Auto-incrementing integer (primary key)
  - **PatientID**: Integer (foreign key to patients table, required)
  - **Amount**: Integer (required, stored in smallest currency unit)
  - **Date**: String (date format, required)
  - **Note**: String (optional, payment description)

- **Relationships**:
  - Many-to-one relationship with Patient
  - Cascade delete: If patient deleted, all payments deleted

---

## Module 6: Settings Management

### Overview
Application settings and configuration management through a sidebar interface.

### Functionalities

#### 6.1 Settings Sidebar
- **Access**: 
  - Settings button in top navigation bar (gear icon)
  - Opens from right side of screen
  - Overlay backdrop (click to close)
  - Close button (X) in top corner

#### 6.2 License Management in Settings
- **Display**:
  - Current license key (masked/displayed)
  - License expiry date
  - License validation status
  - Days remaining/expired message

- **Actions**:
  - **Update License**: Enter new license key
  - **Validate License**: Re-validate current license
  - **Clear License**: Remove license from storage

#### 6.3 Data Management
- **Delete All Patients**:
  - Button to delete all patient data
  - Confirmation dialog (double confirmation)
  - Deletes:
    - All patients
    - All appointments
    - All payments
    - All patient folders
  - Irreversible operation

#### 6.4 Theme Management
- **Features**:
  - Theme toggle (Light/Dark mode)
  - Theme preference saved in localStorage
  - Theme applied globally across application
  - Persistent across sessions

---

## Module 7: File System Integration

### Overview
Automatic file system management for patient-related documents and files.

### Functionalities

#### 7.1 Patient Folder Creation
- **Automatic Creation**:
  - When patient added, folder structure created automatically
  - Base directory: `patient_data/`
  - Patient folder: `patient_data/{patient_id}/`
  - Named subfolder: `patient_data/{patient_id}/{sanitized_name}/`

#### 7.2 Folder Structure
- **Purpose**: 
  - Store patient-specific files (X-rays, documents, photos, etc.)
  - Organized by patient ID and name
  - Easy file management

#### 7.3 Open Patient Folder
- **Functionality**:
  - Button in patient detail view
  - Opens system file explorer to patient's folder
  - Cross-platform support:
    - Windows: Uses `explorer` command
    - macOS: Uses `open` command
    - Linux: Uses `xdg-open` command
  - Creates folder if doesn't exist
  - Non-blocking execution

#### 7.4 Folder Cleanup
- **On Patient Deletion**:
  - Patient folder and all contents deleted
  - Subdirectories recursively removed
  - Transaction-safe (if deletion fails, rollback)

#### 7.5 Name Sanitization
- **Process**:
  - Removes filesystem-unsafe characters
  - Replaces with dashes
  - Limits length to 50 characters
  - Handles special cases (empty names, etc.)

---

## Module 8: Database Architecture

### Overview
SQLite database with proper relational structure and foreign key constraints.

### Database Schema

#### 8.1 Patients Table
```sql
CREATE TABLE patients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    age INTEGER NOT NULL,
    gender TEXT NOT NULL,
    occupation TEXT,
    total_required INTEGER DEFAULT 0
);
```
- **Primary Key**: id
- **Unique Constraint**: phone (enforced in application logic)
- **Indexes**: name (for search performance)

#### 8.2 Appointments Table
```sql
CREATE TABLE appointments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    patient_id INTEGER NOT NULL,
    datetime TEXT NOT NULL,
    duration INTEGER,
    notes TEXT,
    FOREIGN KEY(patient_id) REFERENCES patients(id) ON DELETE CASCADE
);
```
- **Primary Key**: id
- **Foreign Key**: patient_id → patients.id
- **Cascade Delete**: When patient deleted, appointments deleted

#### 8.3 Payments Table
```sql
CREATE TABLE payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    patient_id INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    date TEXT NOT NULL,
    note TEXT,
    FOREIGN KEY(patient_id) REFERENCES patients(id) ON DELETE CASCADE
);
```
- **Primary Key**: id
- **Foreign Key**: patient_id → patients.id
- **Cascade Delete**: When patient deleted, payments deleted

#### 8.4 Foreign Key Constraints
- **Enabled**: Foreign keys enabled at database connection level
- **Verification**: Foreign key status verified on connection
- **Cascade Behavior**: Automatic deletion of related records
- **Transaction Safety**: All operations use transactions

---

## Module 9: User Interface Components

### Overview
Modern, responsive UI built with Svelte components and custom styling.

### Components

#### 9.1 Navigation Bar
- **Features**:
  - Top navigation with module buttons
  - Active state highlighting
  - Settings button (gear icon)
  - Responsive design
  - Full-width layout

#### 9.2 Search Bar
- **Features**:
  - Real-time search input
  - Placeholder text
  - Search icon
  - Clear button
  - Used in Patient Management and Appointments

#### 9.3 Patient Card
- **Display**:
  - Patient name (prominent)
  - Phone number (formatted)
  - Age and gender
  - Action buttons (View, Edit, Delete)
  - Hover effects
  - Consistent sizing (260x260px)

#### 9.4 Modals
- **Add Patient Modal**:
  - Form with all patient fields
  - Validation feedback
  - Submit/Cancel buttons
  - Keyboard support (Enter to submit, Escape to close)

- **Edit Appointment Modal**:
  - Pre-populated form
  - Patient selection dropdown
  - Date-time picker
  - Duration and notes fields

- **Add Appointment Modal**:
  - Patient search/autocomplete
  - Date-time picker
  - Duration and notes fields

#### 9.5 Company Footer
- **Location**: Bottom of Patient Management page
- **Content**:
  - NovaTech Labs logo
  - Company tagline: "Engineering Ideas Into Reality"
  - Contact phone: 0941-414-122
  - Social media links (GitHub, LinkedIn, Facebook, X/Twitter)
  - Product links
  - Company links
  - Support links
  - Copyright information

#### 9.6 Pagination
- **Features**:
  - Previous/Next buttons
  - Page number display
  - Disabled state for first/last page
  - Used in Patient List and Appointments List

#### 9.7 Loading States
- **Spinners**: 
  - Circular loading indicators
  - Shown during data fetch operations
  - Centered in content area

#### 9.8 Error Handling
- **Error Messages**:
  - Displayed in red/error styling
  - Retry buttons where applicable
  - Clear error descriptions
  - License validation errors

---

## Module 10: Data Validation and Business Rules

### Overview
Comprehensive validation rules and business logic enforced throughout the application.

### Validation Rules

#### 10.1 Patient Validation
- **Name**: Required, any string
- **Phone**: Required, must be unique across all patients
- **Age**: Required, must be between 6 and 100 years
- **Gender**: Required, must be selected
- **Occupation**: Optional, any string

#### 10.2 Appointment Validation
- **Patient**: Required, must be existing patient
- **DateTime**: Required, must be valid date-time format (RFC3339)
- **Duration**: Optional, must be positive integer if provided
- **Notes**: Optional, any string

#### 10.3 Payment Validation
- **Patient**: Required, must be existing patient
- **Amount**: Required, must be positive integer
- **Date**: Required, must be valid date
- **Note**: Optional, any string

#### 10.4 License Validation
- **License Key**: Required, must match expected format
- **Expiry Date**: Must be valid date in future (or current date)
- **Format Validation**: Must pass Fibonacci sequence extraction
- **Base64 Decoding**: Must decode to valid date string

---

## Module 11: State Management

### Overview
Reactive state management using Svelte stores for data synchronization.

### Stores

#### 11.1 Patient Store
- **State**:
  - `patients`: Array of all patients
  - `loading`: Boolean loading state
  - `error`: Error message string
  - `searchTerm`: Current search term
  - `selectedPatient`: Currently selected patient object

- **Functions**:
  - `loadPatients()`: Fetch all patients
  - `searchPatients(term)`: Search patients by name/phone
  - `addPatient(data)`: Add new patient
  - `updatePatient(data)`: Update existing patient
  - `deletePatient(id)`: Delete patient
  - `deleteAllPatients()`: Delete all patients
  - `openPatientFolder(id)`: Open patient folder in file explorer

#### 11.2 Appointment Store
- **State**:
  - `appointments`: Array of all appointments
  - `loadingAppointments`: Boolean loading state
  - `appointmentError`: Error message string

- **Functions**:
  - `loadAppointments()`: Fetch all appointments
  - `addAppointment(data)`: Add new appointment
  - `updateAppointment(data)`: Update existing appointment
  - `deleteAppointment(id)`: Delete appointment

#### 11.3 Settings Store
- **State**:
  - `theme`: Current theme (light/dark)
  - `account`: Account/license information object
  - `licenseValidationStatus`: License validation state
  - `currentLicenseKey`: Current valid license key
  - `licenseValid`: Derived boolean (is license valid)

- **Functions**:
  - `validateCurrentLicense()`: Validate stored license
  - `setLicense(key, expiry)`: Set and validate new license
  - `clearLicense()`: Remove license from storage

---

## Module 12: API Architecture

### Overview
RESTful API structure through Wails bindings connecting frontend to Go backend.

### API Endpoints

#### 12.1 License Endpoints
- `ValidateLicense(licenseKey)`: Returns license validation information
- `IsLicenseValid(licenseKey)`: Returns boolean validation status

#### 12.2 Patient Endpoints
- `AddPatient(patient, licenseKey)`: Create new patient
- `GetPatients(licenseKey)`: Retrieve all patients
- `GetPatient(id, licenseKey)`: Retrieve single patient
- `UpdatePatient(patient, licenseKey)`: Update patient
- `DeletePatient(id, licenseKey)`: Delete patient
- `SearchPatients(searchTerm, licenseKey)`: Search patients
- `OpenPatientFolder(id, licenseKey)`: Open patient folder
- `DeleteAllPatients(licenseKey)`: Delete all patients

#### 12.3 Appointment Endpoints
- `AddAppointment(appointment, licenseKey)`: Create appointment
- `GetAppointments(licenseKey)`: Retrieve all appointments
- `GetAppointment(id, licenseKey)`: Retrieve single appointment
- `UpdateAppointment(appointment, licenseKey)`: Update appointment
- `DeleteAppointment(id, licenseKey)`: Delete appointment

#### 12.4 Payment Endpoints
- `AddPayment(payment, licenseKey)`: Create payment
- `GetPaymentsForPatient(patientID, licenseKey)`: Get patient payments
- `GetLastPaymentForPatient(patientID, licenseKey)`: Get most recent payment
- `UpdateTotalRequired(patientID, total, licenseKey)`: Update total required
- `GetPatientBalance(patientID, licenseKey)`: Get balance information
- `DeletePayment(paymentID, licenseKey)`: Delete payment
- `UpdatePayment(payment, licenseKey)`: Update payment

#### 12.5 Security
- **All Endpoints**: Require valid license key parameter
- **License Check**: Executed before every operation
- **Error Handling**: Returns appropriate error messages for invalid licenses
- **Session Management**: License validated per request

---

## Module 13: Application Workflow

### Overview
End-to-end user workflows and application lifecycle.

### Application Startup
1. **License Check**: Application checks for saved license in localStorage
2. **License Validation**: If license found, validates with backend
3. **License Gate**: If no valid license, shows license gate screen
4. **Main Application**: If license valid, shows main application interface

### Patient Management Workflow
1. **View Patients**: User navigates to Patient Management tab
2. **Search (Optional)**: User can search for specific patient
3. **View Details**: Click patient card to see full details
4. **Actions**: 
   - Add new patient
   - Edit existing patient
   - Delete patient
   - Open patient folder
   - View appointments/payments

### Appointment Workflow
1. **View Appointments**: Navigate to Appointments tab
2. **Filter (Optional)**: Filter by date or patient name
3. **Add Appointment**: Click to add, select patient, date/time, duration, notes
4. **Edit Appointment**: Click appointment to modify
5. **Delete Appointment**: Click delete with confirmation

### Calendar Workflow
1. **View Calendar**: Navigate to Calendar tab
2. **Switch Views**: Toggle between Month and Week views
3. **Navigate**: Use Previous/Next/Today buttons
4. **Interact**: Click appointment to edit, click Add to create new

### Payment Workflow
1. **View Payments**: Navigate to Payments tab
2. **Select Patient**: Click patient to view payment details
3. **Add Payment**: Enter amount, date, note
4. **Edit Payment**: Modify existing payment
5. **Update Total**: Set total required amount for patient
6. **View Balance**: See calculated remaining balance

---

## Technical Architecture Details

### Backend Architecture
- **Language**: Go 1.23.0
- **Database**: SQLite (modernc.org/sqlite driver)
- **Pattern**: Handler-based architecture
- **Transactions**: Used for data integrity
- **Error Handling**: Comprehensive error messages
- **Logging**: Debug logging for troubleshooting

### Frontend Architecture
- **Framework**: Svelte 3.49.0
- **Build Tool**: Vite 3.0.7
- **State Management**: Svelte stores (writable, derived)
- **Styling**: CSS variables for theming
- **Components**: Modular component architecture
- **Routing**: Component-based navigation

### Integration
- **Wails v2**: Connects Go backend to Svelte frontend
- **Type Safety**: TypeScript definitions for Go methods
- **Asset Embedding**: Frontend assets embedded in binary
- **Cross-Platform**: Single codebase for Windows, macOS, Linux

### Data Flow
1. **User Action**: User interacts with UI component
2. **Store Function**: Svelte store function called
3. **API Call**: Wails binding calls Go method
4. **License Check**: License validated
5. **Handler**: Go handler processes request
6. **Database**: SQLite database queried/updated
7. **Response**: Data returned to frontend
8. **Store Update**: Svelte store updated
9. **UI Update**: Component reactively updates

---

## Conclusion

DentistApp is a comprehensive dental practice management system providing:

- **Complete Patient Management**: Full CRUD operations with file system integration
- **Appointment Scheduling**: Calendar-based scheduling with filtering
- **Payment Tracking**: Financial management with balance calculations
- **Security**: License-based access control
- **User Experience**: Modern, responsive interface
- **Data Integrity**: Foreign key constraints and transactions
- **Cross-Platform**: Works on Windows, macOS, and Linux

The application is production-ready with proper error handling, validation, and data management capabilities.

---

**Report Generated**: 2025-01-07
**Application Version**: DentistApp (Wails v2)
**Documentation Version**: 1.0

