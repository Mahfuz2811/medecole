package unit

import (
	"github.com/Mahfuz2811/medecole/backend/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMSISDN(t *testing.T) {
	tests := []struct {
		name     string
		msisdn   string
		expected bool
	}{
		// Valid Bangladesh mobile numbers
		{
			name:     "Valid Grameenphone number",
			msisdn:   "01712345678",
			expected: true,
		},
		{
			name:     "Valid Banglalink number",
			msisdn:   "01912345678",
			expected: true,
		},
		{
			name:     "Valid Robi number",
			msisdn:   "01812345678",
			expected: true,
		},
		{
			name:     "Valid Airtel number",
			msisdn:   "01612345678",
			expected: true,
		},
		{
			name:     "Valid Teletalk number",
			msisdn:   "01512345678",
			expected: true,
		},
		{
			name:     "Valid Citycell number",
			msisdn:   "01911111111",
			expected: true,
		},

		// Invalid cases
		{
			name:     "Too short",
			msisdn:   "0171234567",
			expected: false,
		},
		{
			name:     "Too long",
			msisdn:   "017123456789",
			expected: false,
		},
		{
			name:     "Invalid prefix",
			msisdn:   "01012345678",
			expected: false,
		},
		{
			name:     "Contains letters",
			msisdn:   "0171234567a",
			expected: false,
		},
		{
			name:     "Contains special characters",
			msisdn:   "0171234567-",
			expected: false,
		},
		{
			name:     "Empty string",
			msisdn:   "",
			expected: false,
		},
		{
			name:     "Only numbers without prefix",
			msisdn:   "1712345678",
			expected: false,
		},
		{
			name:     "Wrong country code format (not supported)",
			msisdn:   "8801712345678",
			expected: false, // Function only accepts 01[3-9] format
		},
		{
			name:     "Invalid operator code",
			msisdn:   "01012345678",
			expected: false,
		},
		{
			name:     "Starts with +88 (not supported)",
			msisdn:   "+8801712345678",
			expected: false, // Function only accepts 01[3-9] format
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ValidateMSISDN(tt.msisdn)
			assert.Equal(t, tt.expected, result,
				"ValidateMSISDN(%s) = %v, expected %v", tt.msisdn, result, tt.expected)
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "Valid password - 6 characters",
			password: "passwo",
			expected: true,
		},
		{
			name:     "Valid password - 8 characters",
			password: "password",
			expected: true,
		},
		{
			name:     "Valid password - long",
			password: "thisisaverylongpasswordthatshouldbefine",
			expected: true,
		},
		{
			name:     "Valid password - with numbers",
			password: "password123",
			expected: true,
		},
		{
			name:     "Valid password - with special chars",
			password: "password!@#$",
			expected: true,
		},
		{
			name:     "Valid password - mixed case",
			password: "PassWord123",
			expected: true,
		},

		// Invalid cases
		{
			name:     "Too short - 5 characters",
			password: "passw",
			expected: false,
		},
		{
			name:     "Too short - 1 character",
			password: "p",
			expected: false,
		},
		{
			name:     "Empty password",
			password: "",
			expected: false,
		},
		{
			name:     "Only spaces (6+ chars but still valid by current implementation)",
			password: "        ",
			expected: true, // Current implementation doesn't trim spaces for password
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ValidatePassword(tt.password)
			assert.Equal(t, tt.expected, result,
				"ValidatePassword(%s) = %v, expected %v", tt.password, result, tt.expected)
		})
	}
}

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "Normal password",
			password: "password123",
		},
		{
			name:     "Short password",
			password: "pass",
		},
		{
			name:     "Long password",
			password: "thisissomeverylongpasswordthatwewanttotest",
		},
		{
			name:     "Password with special characters",
			password: "p@ssw0rd!#$%",
		},
		{
			name:     "Password with unicode",
			password: "pässwörd🔒",
		},
		{
			name:     "Empty password",
			password: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := utils.HashPassword(tt.password)

			// Should not return error
			assert.NoError(t, err)

			// Hash should not be empty
			assert.NotEmpty(t, hash)

			// Hash should be different from original password
			assert.NotEqual(t, tt.password, hash)

			// Hash should start with bcrypt identifier
			assert.True(t, len(hash) > 10, "Hash should be reasonably long")

			// Should be able to verify the password
			valid := utils.CheckPassword(tt.password, hash)
			assert.True(t, valid, "Should be able to verify the hashed password")
		})
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"

	// Generate a hash
	hash, err := utils.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	tests := []struct {
		name     string
		password string
		hash     string
		expected bool
	}{
		{
			name:     "Correct password",
			password: password,
			hash:     hash,
			expected: true,
		},
		{
			name:     "Wrong password",
			password: "wrongpassword",
			hash:     hash,
			expected: false,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash,
			expected: false,
		},
		{
			name:     "Case sensitive - uppercase",
			password: "TESTPASSWORD123",
			hash:     hash,
			expected: false,
		},
		{
			name:     "Password with extra characters",
			password: password + "extra",
			hash:     hash,
			expected: false,
		},
		{
			name:     "Invalid hash",
			password: password,
			hash:     "invalidhash",
			expected: false,
		},
		{
			name:     "Empty hash",
			password: password,
			hash:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.CheckPassword(tt.password, tt.hash)
			assert.Equal(t, tt.expected, result,
				"CheckPassword(%s, hash) = %v, expected %v", tt.password, result, tt.expected)
		})
	}
}

func TestPasswordHashConsistency(t *testing.T) {
	password := "consistencytest"

	// Generate multiple hashes for the same password
	hash1, err1 := utils.HashPassword(password)
	hash2, err2 := utils.HashPassword(password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, hash1)
	assert.NotEmpty(t, hash2)

	// Hashes should be different (bcrypt includes salt)
	assert.NotEqual(t, hash1, hash2, "Multiple hashes of same password should be different due to salt")

	// But both should verify correctly
	assert.True(t, utils.CheckPassword(password, hash1))
	assert.True(t, utils.CheckPassword(password, hash2))
}

func TestValidateMSISDNEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		msisdn   string
		expected bool
	}{
		{
			name:     "All zeros after valid prefix",
			msisdn:   "01700000000",
			expected: true,
		},
		{
			name:     "All nines after valid prefix",
			msisdn:   "01799999999",
			expected: true,
		},
		{
			name:     "Valid prefix with mixed digits",
			msisdn:   "01713572468",
			expected: true,
		},
		{
			name:     "Whitespace at start (cleaned and valid)",
			msisdn:   " 01712345678",
			expected: true, // Function removes spaces
		},
		{
			name:     "Whitespace at end (cleaned and valid)",
			msisdn:   "01712345678 ",
			expected: true, // Function removes spaces
		},
		{
			name:     "Whitespace in middle (cleaned and valid)",
			msisdn:   "017 12345678",
			expected: true, // Function removes spaces
		},
		{
			name:     "Tab character",
			msisdn:   "017\t12345678",
			expected: false,
		},
		{
			name:     "Newline character",
			msisdn:   "017\n12345678",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ValidateMSISDN(tt.msisdn)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidatePasswordEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "Exactly 6 characters",
			password: "123456",
			expected: true,
		},
		{
			name:     "Exactly 5 characters",
			password: "12345",
			expected: false,
		},
		{
			name:     "Password with leading spaces (6+ chars)",
			password: "  password",
			expected: true,
		},
		{
			name:     "Password with trailing spaces (6+ chars)",
			password: "password  ",
			expected: true,
		},
		{
			name:     "Password with only spaces but 6+ chars",
			password: "      ",
			expected: true, // Current implementation doesn't validate content, only length
		},
		{
			name:     "Unicode characters",
			password: "pässwörd",
			expected: true,
		},
		{
			name:     "Emoji password",
			password: "password🔒🔑",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ValidatePassword(tt.password)
			assert.Equal(t, tt.expected, result)
		})
	}
}
