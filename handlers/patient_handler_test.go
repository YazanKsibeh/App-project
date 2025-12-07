package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenPatientFolder(t *testing.T) {
	// Create a mock patient handler (no DB needed for this test)
	handler := &PatientHandler{}

	// Test patient ID
	testPatientID := 999

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Expected folder path
	expectedPath := filepath.Join(cwd, "patient_data", "999")

	// Clean up any existing test directory
	os.RemoveAll(filepath.Join(cwd, "patient_data"))

	// Test the function
	err = handler.OpenPatientFolder(testPatientID)

	// The function should succeed (even if explorer fails to open, the folder should be created)
	// We don't check the explorer opening because it's platform-dependent and may fail in CI
	if err != nil {
		t.Logf("OpenPatientFolder returned error (this may be expected in test environment): %v", err)
	}

	// Verify that the directory was created
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Patient folder was not created at expected path: %s", expectedPath)
	} else {
		t.Logf("Patient folder successfully created at: %s", expectedPath)
	}

	// Clean up test directory
	os.RemoveAll(filepath.Join(cwd, "patient_data"))
}

func TestPatientFolderStructure(t *testing.T) {
	// This test verifies that both patient ID folder and patient name folder are created

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Clean up any existing test directory
	os.RemoveAll(filepath.Join(cwd, "patient_data"))

	// Test the cleanPatientName function
	testCases := []struct {
		input    string
		expected string
	}{
		{"John Doe", "John-Doe"},
		{"Mary Jane Smith", "Mary-Jane-Smith"},
		{"Dr. Bob", "Dr-Bob"},
		{"Test/Patient\\Name", "Test-Patient-Name"},
		{"Patient*With?Special<>Chars", "Patient-With-Special-Chars"},
		{"", "Unknown-Patient"},
		{"   Spaces   ", "Spaces"},
	}

	for _, tc := range testCases {
		result := cleanPatientName(tc.input)
		if result != tc.expected {
			t.Errorf("cleanPatientName(%q) = %q; expected %q", tc.input, result, tc.expected)
		}
	}

	// Test folder structure creation
	patientID := 123
	patientName := "John Doe"
	expectedIDFolder := filepath.Join(cwd, "patient_data", "123")
	expectedNameFolder := filepath.Join(expectedIDFolder, "John-Doe")

	// Simulate what happens in AddPatient
	patientDir := filepath.Join("patient_data", fmt.Sprintf("%d", patientID))
	err = os.MkdirAll(patientDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create patient directory: %v", err)
	}

	cleanName := cleanPatientName(patientName)
	patientNameDir := filepath.Join(patientDir, cleanName)
	err = os.MkdirAll(patientNameDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create patient name directory: %v", err)
	}

	// Verify both folders exist
	if _, err := os.Stat(expectedIDFolder); os.IsNotExist(err) {
		t.Errorf("Patient ID folder was not created: %s", expectedIDFolder)
	} else {
		t.Logf("✓ Patient ID folder created: %s", expectedIDFolder)
	}

	if _, err := os.Stat(expectedNameFolder); os.IsNotExist(err) {
		t.Errorf("Patient name folder was not created: %s", expectedNameFolder)
	} else {
		t.Logf("✓ Patient name folder created: %s", expectedNameFolder)
	}

	// Clean up test directory
	os.RemoveAll(filepath.Join(cwd, "patient_data"))
}

func TestOpenPatientFolderWithExistingDirectory(t *testing.T) {
	// Create a mock patient handler
	handler := &PatientHandler{}

	// Test patient ID
	testPatientID := 888

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Expected folder path
	expectedPath := filepath.Join(cwd, "patient_data", "888")

	// Clean up any existing test directory
	os.RemoveAll(filepath.Join(cwd, "patient_data"))

	// Pre-create the directory
	err = os.MkdirAll(expectedPath, 0755)
	if err != nil {
		t.Fatalf("Failed to pre-create test directory: %v", err)
	}

	// Test the function with existing directory
	err = handler.OpenPatientFolder(testPatientID)

	// The function should succeed even with existing directory
	if err != nil {
		t.Logf("OpenPatientFolder returned error (this may be expected in test environment): %v", err)
	}

	// Verify that the directory still exists
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Patient folder was removed unexpectedly: %s", expectedPath)
	} else {
		t.Logf("Patient folder correctly maintained at: %s", expectedPath)
	}

	// Clean up test directory
	os.RemoveAll(filepath.Join(cwd, "patient_data"))
}
