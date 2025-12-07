import './style.css'
import App from './App.svelte'
import { debugLicenseStatus, testOpenPatientFolder } from './stores/patientStore.js'

const app = new App({
  target: document.getElementById('app')
})

// Add debug functions to global window for console testing
window.debugLicenseStatus = debugLicenseStatus;
window.testOpenPatientFolder = testOpenPatientFolder;

export default app
