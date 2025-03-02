package config

import (
	"fmt"
	"os"
	"path/filepath"

	"hypr-exiled/pkg/logger"
)

// DefaultConfig creates a default configuration.
func DefaultConfig(log *logger.Logger) (*Config, error) {
	log.Debug("Creating default configuration")

	logPath, err := getDefaultPoeLogPath(log)
	if err != nil {
		log.Error("Failed to get default POE log path", err)
		return nil, fmt.Errorf("failed to get PoE log file, pls create the config file manually: %w", err)
	}

	config := &Config{
		poeLogPath: logPath,
		triggers: map[string]string{
			"incoming_trade": `\[INFO Client \d+\] @From (?:<.*> |)(\w+): Hi, I would like to buy your ([^,]+(?:,[^,]+)*) listed for (\d+(?:\.\d+)?) ([^ ]+) in ([^\(]+) \(stash tab "([^"]+)"; position: left (\d+), top (\d+)\)`,
			"outgoing_trade": `\[INFO Client \d+\] @To (?:<.*> |)(\w+): Hi, I would like to buy your ([^,]+(?:,[^,]+)*) listed for (\d+(?:\.\d+)?) ([^ ]+) in ([^\(]+) \(stash tab "([^"]+)"; position: left (\d+), top (\d+)\)`,
		},
		commands: map[string][]string{
			"party":  {"/invite {player}"},
			"finish": {"/kick {player}", "@{player} thanks!"},
			"trade":  {"/tradewith {player}"},
		},
		notifyCommand: "",
		log:           log,
	}

	log.Info("Created default configuration",
		"log_path", logPath,
		"trigger_count", len(config.triggers),
		"command_count", len(config.commands))

	if err := config.compile(); err != nil {
		log.Error("Failed to compile default config triggers", err)
		return nil, fmt.Errorf("failed to compile default config: %w", err)
	}

	return config, nil
}

// getDefaultPoeLogPath finds the default Path of Exile 2 log file.
func getDefaultPoeLogPath(log *logger.Logger) (string, error) {
	log.Debug("Looking for default POE log path")

	home, err := os.UserHomeDir()
	if err != nil {
		log.Error("Failed to get home directory", err)
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Common POE2 log paths
	possiblePaths := []string{
		filepath.Join(home, ".local", "share", "Steam", "steamapps", "common", "Path of Exile 2", "logs", "Client.txt"),
		filepath.Join(home, "Games", "Path of Exile 2", "logs", "Client.txt"),
	}

	for _, path := range possiblePaths {
		log.Debug("Checking possible log path", "path", path)
		if _, err := os.Stat(path); err == nil {
			log.Info("Found POE log file", "path", path)
			return path, nil
		}
	}

	log.Error("No valid POE log file found", nil, "checked_paths", possiblePaths)
	return "", fmt.Errorf("no valid Path of Exile 2 log file found in common locations")
}
