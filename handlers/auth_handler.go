package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"DentistApp/models"

	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication-related operations
type AuthHandler struct {
	db *sql.DB
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// InitializeAdmin creates the default admin user if it doesn't exist
func (h *AuthHandler) InitializeAdmin() error {
	// Check if admin user exists
	var count int
	err := h.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", "admin").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check admin user: %v", err)
	}

	if count > 0 {
		// Admin already exists
		return nil
	}

	// Create default admin user
	// Default credentials: admin / admin123
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %v", err)
	}

	query := `INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)`
	_, err = h.db.Exec(query, "admin", string(passwordHash), "Admin")
	if err != nil {
		return fmt.Errorf("failed to create admin user: %v", err)
	}

	return nil
}

// Login validates user credentials and returns a session token
func (h *AuthHandler) Login(username, password string) (*models.LoginResponse, error) {
	var user models.User
	query := `SELECT id, username, password_hash, role FROM users WHERE username = ?`
	err := h.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.LoginResponse{
				Success: false,
				Message: "Invalid username or password",
			}, nil
		}
		return nil, fmt.Errorf("failed to query user: %v", err)
	}

	// Compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return &models.LoginResponse{
			Success: false,
			Message: "Invalid username or password",
		}, nil
	}

	// Generate session token
	token, err := generateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %v", err)
	}

	// Clear password hash from response
	user.PasswordHash = ""

	return &models.LoginResponse{
		Success: true,
		Token:   token,
		User:    user,
		Message: "Login successful",
	}, nil
}

// ValidateSession validates a session token (simple implementation)
// In production, you might want to store tokens in database or use JWT
func (h *AuthHandler) ValidateSession(token string) (*models.User, error) {
	// For now, we'll use a simple approach where token contains user info
	// In production, implement proper token validation (JWT, database lookup, etc.)
	// This is a basic implementation - you may want to enhance it
	
	// For simplicity, we'll extract username from token format: username_timestamp_random
	// In production, use JWT or store tokens in database
	
	// For now, return nil to indicate we need to re-authenticate
	// This is a placeholder - implement proper token validation
	return nil, fmt.Errorf("session validation not fully implemented - please re-login")
}

// CreateUser creates a new user (admin only)
func (h *AuthHandler) CreateUser(userForm models.UserForm, createdByID int) (int64, error) {
	// Verify that the creator is an admin
	var creatorRole string
	err := h.db.QueryRow("SELECT role FROM users WHERE id = ?", createdByID).Scan(&creatorRole)
	if err != nil {
		return 0, fmt.Errorf("failed to verify creator: %v", err)
	}

	if creatorRole != "Admin" {
		return 0, fmt.Errorf("only admins can create users")
	}

	// Check if username already exists
	var existingCount int
	err = h.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", userForm.Username).Scan(&existingCount)
	if err != nil {
		return 0, fmt.Errorf("failed to check username uniqueness: %v", err)
	}
	if existingCount > 0 {
		return 0, fmt.Errorf("username already exists")
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userForm.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %v", err)
	}

	// Set default role if not provided
	role := userForm.Role
	if role == "" {
		role = "Dentist"
	}

	// Insert user
	query := `INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)`
	result, err := h.db.Exec(query, userForm.Username, string(passwordHash), role)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get user ID: %v", err)
	}

	return id, nil
}

// GetUserByID retrieves a user by ID
func (h *AuthHandler) GetUserByID(id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, password_hash, role FROM users WHERE id = ?`
	err := h.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %v", err)
	}

	// Clear password hash
	user.PasswordHash = ""
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (h *AuthHandler) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, password_hash, role FROM users WHERE username = ?`
	err := h.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %v", err)
	}

	// Clear password hash
	user.PasswordHash = ""
	return &user, nil
}

// GetAllUsers returns all users (admin only)
func (h *AuthHandler) GetAllUsers() ([]models.User, error) {
	query := `SELECT id, username, role FROM users ORDER BY username`
	rows, err := h.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %v", err)
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Role)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// generateSessionToken generates a simple session token
// In production, consider using JWT or storing tokens in database
func generateSessionToken() (string, error) {
	// Generate a random token
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	// Combine with timestamp for uniqueness
	timestamp := time.Now().Unix()
	token := fmt.Sprintf("%d_%s", timestamp, hex.EncodeToString(bytes))
	
	return token, nil
}

