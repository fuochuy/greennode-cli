package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// ConfigFileWriter creates/updates INI config files.
type ConfigFileWriter struct {
	configDir string
}

// NewConfigFileWriter creates a new writer targeting the default config directory.
func NewConfigFileWriter() *ConfigFileWriter {
	return &ConfigFileWriter{configDir: DefaultConfigDir()}
}

// ensureDir creates the config directory with proper permissions.
func (w *ConfigFileWriter) ensureDir() error {
	return os.MkdirAll(w.configDir, 0700)
}

// WriteCredentials writes client_id and client_secret for the given profile.
func (w *ConfigFileWriter) WriteCredentials(profile, clientID, clientSecret string) error {
	if err := w.ensureDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	filePath := filepath.Join(w.configDir, "credentials")
	cfg, err := w.loadOrCreate(filePath)
	if err != nil {
		return err
	}

	section, err := cfg.NewSection(profile)
	if err != nil {
		return fmt.Errorf("failed to create section '%s': %w", profile, err)
	}
	section.Key("client_id").SetValue(clientID)
	section.Key("client_secret").SetValue(clientSecret)

	return w.save(cfg, filePath)
}

// WriteConfig writes region, output, and project_id for the given profile.
// An empty projectID is written as an empty key to explicitly clear any
// previously-saved value.
func (w *ConfigFileWriter) WriteConfig(profile, region, output, projectID string) error {
	if err := w.ensureDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	filePath := filepath.Join(w.configDir, "config")
	cfg, err := w.loadOrCreate(filePath)
	if err != nil {
		return err
	}

	sectionName := profile
	if profile != "default" {
		sectionName = "profile " + profile
	}

	section, err := cfg.NewSection(sectionName)
	if err != nil {
		return fmt.Errorf("failed to create section '%s': %w", sectionName, err)
	}
	section.Key("region").SetValue(region)
	section.Key("output").SetValue(output)
	section.Key("project_id").SetValue(projectID)

	return w.save(cfg, filePath)
}

func (w *ConfigFileWriter) loadOrCreate(filePath string) (*ini.File, error) {
	if _, err := os.Stat(filePath); err == nil {
		return ini.Load(filePath)
	}
	return ini.Empty(), nil
}

func (w *ConfigFileWriter) save(cfg *ini.File, filePath string) error {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to write %s: %w", filePath, err)
	}
	defer f.Close()

	_, err = cfg.WriteTo(f)
	return err
}
