package utils

import (
	"fmt"
	"strings"
)

// IsFilterSafe performs basic sanitization check on filter string
// Rejects filters containing dangerous SQL keywords and characters
func IsFilterSafe(filter string) bool {
	if filter == "" || filter == "1=1" {
		return true
	}

	// List of dangerous SQL keywords and characters
	dangerous := []string{
		";",      // Statement separator
		"--",     // SQL comment
		"/*",     // Block comment start
		"*/",     // Block comment end
		"DROP",   // DDL
		"DELETE", // DML
		"UPDATE", // DML (when not part of updated_at column)
		"INSERT", // DML
		"EXEC",   // Execute
		"EXECUTE",
		"XP_",   // Extended procedures
		"SP_",   // System procedures
		"UNION", // SQL injection technique
		"CONCAT",
		"CHAR(",
		"NCHAR(",
		"VARCHAR(",
		"NVARCHAR(",
	}

	filterUpper := strings.ToUpper(filter)

	for _, d := range dangerous {
		if strings.Contains(filterUpper, d) {
			return false
		}
	}

	return true
}

// IsSortSafe performs basic sanitization check on sort column
// Rejects sort containing dangerous characters
func IsSortSafe(sort string) bool {
	if sort == "" {
		return true
	}

	// Reject if contains dangerous characters
	dangerous := []string{";", "--", "/*", "*/", "(", ")", "=", "'", "\""}

	for _, d := range dangerous {
		if strings.Contains(sort, d) {
			return false
		}
	}

	return true
}

// SanitizeFilter returns error if filter is unsafe
func SanitizeFilter(filter string) error {
	if !IsFilterSafe(filter) {
		return fmt.Errorf("invalid filter: potentially dangerous SQL detected")
	}
	return nil
}

// SanitizeSort returns error if sort is unsafe
func SanitizeSort(sort string) error {
	if !IsSortSafe(sort) {
		return fmt.Errorf("invalid sort: potentially dangerous characters detected")
	}
	return nil
}
