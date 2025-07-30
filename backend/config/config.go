package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var Cfg *Config

type Config struct {
	Port         string
	SMTPEmail    string
	SMTPPassword string
	SMTPHost     string
	SMTPPort     string

	SupabaseHost     string
	SupabasePort     string
	SupabaseUser     string
	SupabasePassword string
	SupabaseDB       string

	SecretKey string

	GcloudPrintFilesBucket string

	EmailWhitelistEnabled bool

	AdminEmail     string
	AdminFirstName string
	AdminLastName  string
}

func LoadConfig() error {
	_ = godotenv.Load()

	requiredKeys := []string{
		"PORT", "SMTP_EMAIL", "SMTP_PASSWORD", "SMTP_HOST", "SMTP_PORT",
		"SUPABASE_HOST", "SUPABASE_PORT", "SUPABASE_USER", "SUPABASE_PASSWORD", "SUPABASE_DB", "SECRET_KEY",
		"GCLOUD_PRINT_FILES_BUCKET", "EMAIL_WHITELIST_ENABLED",
		"ADMIN_EMAIL", "ADMIN_FIRST_NAME", "ADMIN_LAST_NAME",
	}

	for _, key := range requiredKeys {
		if os.Getenv(key) == "" {
			return fmt.Errorf("missing key %s in environment variables", key)
		}
	}

	Cfg = &Config{
		Port:             os.Getenv("PORT"),
		SMTPEmail:        os.Getenv("SMTP_EMAIL"),
		SMTPPassword:     os.Getenv("SMTP_PASSWORD"),
		SMTPHost:         os.Getenv("SMTP_HOST"),
		SMTPPort:         os.Getenv("SMTP_PORT"),
		SupabaseHost:     os.Getenv("SUPABASE_HOST"),
		SupabasePort:     os.Getenv("SUPABASE_PORT"),
		SupabaseUser:     os.Getenv("SUPABASE_USER"),
		SupabasePassword: os.Getenv("SUPABASE_PASSWORD"),
		SupabaseDB:       os.Getenv("SUPABASE_DB"),

		SecretKey:              os.Getenv("SECRET_KEY"),
		GcloudPrintFilesBucket: os.Getenv("GCLOUD_PRINT_FILES_BUCKET"),
		EmailWhitelistEnabled:  os.Getenv("EMAIL_WHITELIST_ENABLED") == "true",

		AdminEmail:     os.Getenv("ADMIN_EMAIL"),
		AdminFirstName: os.Getenv("ADMIN_FIRST_NAME"),
		AdminLastName:  os.Getenv("ADMIN_LAST_NAME"),
	}

	return nil
}
