package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
)

// config is the command's local configuration.
type config struct {
	KeyRef string `toml:"key_ref"`
}

// loadConfig reads and checks a config file.
func loadConfig(path string) (config, error) {
	var c config
	md, err := toml.DecodeFile(path, &c)
	if err != nil {
		return config{}, fmt.Errorf("config %s: %w", path, err)
	}
	if undecoded := md.Undecoded(); len(undecoded) > 0 {
		keys := make([]string, len(undecoded))
		for i, k := range undecoded {
			keys[i] = k.String()
		}
		return config{}, fmt.Errorf("config %s: unknown keys %s", path, strings.Join(keys, ", "))
	}
	if c.KeyRef == "" {
		return config{}, fmt.Errorf("config %s: missing key_ref", path)
	}
	return c, nil
}

// readKey returns the secret at ref through the 1Password CLI op. A failure reports op's
// stderr, never its stdout, which may hold the secret.
func readKey(ctx context.Context, op, ref string) (string, error) {
	out, err := exec.CommandContext(ctx, op, "read", ref).Output()
	if err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			if msg := strings.TrimSpace(string(exit.Stderr)); msg != "" {
				return "", fmt.Errorf("op read %s: %w: %s", ref, err, msg)
			}
		}
		return "", fmt.Errorf("op read %s: %w", ref, err)
	}
	key := strings.TrimSpace(string(out))
	if key == "" {
		return "", errors.New("op read " + ref + ": empty secret")
	}
	return key, nil
}
