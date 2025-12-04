/**
 * Formats MSISDN (phone number) to ensure proper format
 * - Removes all non-digits
 * - Ensures it starts with 0
 * - Limits to 11 digits max
 */
export const formatMSISDN = (value: string): string => {
	// Remove all non-digits
	const digits = value.replace(/\D/g, "");

	// Limit to 11 digits and ensure it starts with 0
	if (digits.length === 0) return "";

	// If user types without leading 0, add it
	if (digits.length > 0 && !digits.startsWith("0")) {
		return "0" + digits.slice(0, 10);
	}

	// Limit to 11 digits max
	return digits.slice(0, 11);
};

/**
 * Validates MSISDN format
 */
export const validateMSISDN = (msisdn: string): string | null => {
	if (msisdn.length !== 11 || !msisdn.startsWith("0")) {
		return "Phone number must be 11 digits starting with 0";
	}
	return null;
};

/**
 * Validates password requirements
 */
export const validatePassword = (password: string): string | null => {
	if (password.length < 6) {
		return "Password must be at least 6 characters long";
	}
	return null;
};

/**
 * Validates password confirmation
 */
export const validatePasswordConfirmation = (
	password: string,
	confirmPassword: string
): string | null => {
	if (password !== confirmPassword) {
		return "Passwords don't match!";
	}
	return null;
};
