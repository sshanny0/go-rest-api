package utils

import (
	"fmt"
	"log"
)

// SendResetEmail is a minimal, pluggable mailer for password resets.
// Currently it just logs the reset link. Replace this with SMTP or external
// provider integration as needed (config via ENV).
func SendResetEmail(toEmail, resetLink string) error {
	// TODO: integrate with SMTP or third-party email provider
	log.Printf("[mailer] Sending password reset email to %s: %s", toEmail, resetLink)
	// For debugging/development return nil. In production, return error on failure.
	fmt.Printf("Password reset link for %s: %s\n", toEmail, resetLink)
	return nil
}
