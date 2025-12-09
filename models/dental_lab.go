package models

// DentalLab represents a dental lab vendor
type DentalLab struct {
	ID             int    `json:"id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	ContactPerson  string `json:"contact_person"`
	PhonePrimary   string `json:"phone_primary"`
	PhoneSecondary string `json:"phone_secondary"`
	Email          string `json:"email"`
	Specialties    string `json:"specialties"`
	IsActive       bool   `json:"is_active"`
	Notes          string `json:"notes"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// DentalLabForm represents data needed to create/update a dental lab
type DentalLabForm struct {
	Name           string `json:"name"`
	ContactPerson  string `json:"contact_person"`
	PhonePrimary   string `json:"phone_primary"`
	PhoneSecondary string `json:"phone_secondary"`
	Email          string `json:"email"`
	Specialties    string `json:"specialties"`
	IsActive       bool   `json:"is_active"`
	Notes          string `json:"notes"`
}

// DentalLabsResponse represents paginated dental labs response
type DentalLabsResponse struct {
	Labs        []DentalLab `json:"labs"`
	CurrentPage int         `json:"current_page"`
	TotalPages  int         `json:"total_pages"`
	TotalCount  int         `json:"total_count"`
	PageSize    int         `json:"page_size"`
}

