package handlers

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"
)

// LicenseService handles license validation
type LicenseService struct{}

// NewLicenseService creates a new license service
func NewLicenseService() *LicenseService {
	return &LicenseService{}
}

// LicenseInfo represents license validation result
type LicenseInfo struct {
	IsValid    bool   `json:"is_valid"`
	ExpiryDate string `json:"expiry_date"`
	Message    string `json:"message"`
}

const layout = "20060102"

// arrayInit returns the fibonacci sequence used in key generation
func arrayInit() []int {
	return []int{3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987}
}

// reverse reverses a string (from gen_key package)
func reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}

// ValidateLicense validates a license key and returns license information
func (ls *LicenseService) ValidateLicense(licenseKey string) (*LicenseInfo, error) {
	if len(licenseKey) == 0 {
		return &LicenseInfo{
			IsValid: false,
			Message: "License key is required",
		}, nil
	}

	// Extract and validate the license key using the same algorithm from gen_key
	array := arrayInit()

	// Get offset from key
	var offset int64
	var str string
	var err error

	// Try to parse offset (numbers from 10 to 99)
	if len(licenseKey) >= 3 {
		offset, err = strconv.ParseInt(licenseKey[1:3], 10, 64)
		if err != nil {
			// Try single digit (1 to 9)
			if len(licenseKey) >= 2 {
				offset, err = strconv.ParseInt(licenseKey[1:2], 10, 64)
				if err != nil {
					return &LicenseInfo{
						IsValid: false,
						Message: "Invalid license key format",
					}, nil
				}
			} else {
				return &LicenseInfo{
					IsValid: false,
					Message: "Invalid license key format",
				}, nil
			}
		}
	} else {
		return &LicenseInfo{
			IsValid: false,
			Message: "License key too short",
		}, nil
	}

	// Extract date from key using fibonacci sequence
	for _, v := range array {
		if len(licenseKey) >= int(offset)+v+1 {
			str += string(licenseKey[int(offset)+v+1])
		}
	}

	if len(str) == 0 {
		return &LicenseInfo{
			IsValid: false,
			Message: "Could not extract date from license key",
		}, nil
	}

	// Handle base64 padding - trim to exactly 11 characters if longer
	if len(str) > 11 {
		for len(str) > 11 {
			str = str[:len(str)-1]
		}
	}

	// Reverse the string and add padding
	str = reverse(str)
	str = str + "="

	// Decode base64
	b64, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return &LicenseInfo{
			IsValid: false,
			Message: "Invalid license key encoding",
		}, nil
	}

	dateStr := string(b64)

	// Parse the date
	expiryDate, err := time.Parse(layout, dateStr)
	if err != nil {
		return &LicenseInfo{
			IsValid: false,
			Message: "Invalid date in license key",
		}, nil
	}

	// Check if license is expired
	now := time.Now()
	isValid := now.Before(expiryDate) || now.Equal(expiryDate)

	var message string
	if isValid {
		daysLeft := int(expiryDate.Sub(now).Hours() / 24)
		if daysLeft <= 7 {
			message = fmt.Sprintf("License expires in %d days", daysLeft)
		} else {
			message = "License is valid"
		}
	} else {
		daysExpired := int(now.Sub(expiryDate).Hours() / 24)
		message = fmt.Sprintf("License expired %d days ago", daysExpired)
	}

	return &LicenseInfo{
		IsValid:    isValid,
		ExpiryDate: expiryDate.Format("2006-01-02"),
		Message:    message,
	}, nil
}

// IsLicenseValid is a simplified method that returns just the validation status
func (ls *LicenseService) IsLicenseValid(licenseKey string) bool {
	info, err := ls.ValidateLicense(licenseKey)
	if err != nil {
		return false
	}
	return info.IsValid
}
