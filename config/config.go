package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HomeDir string `yaml:"home" config:"desc=Base directory for all compositions,default=~/fibber/compositions"`
	Author  string `yaml:"author" config:"desc=Default author name for new compositions"`
	Genre   string `yaml:"default_genre" config:"desc=Fallback genre if none specified,default=general"`
	LogFmt  string `yaml:"log_fmt" config:"desc=Log output format (json, formatted, text),default=json"`
}

var defaultConfig = Config{
	HomeDir: "~/fibber/compositions",
	Genre:   "general",
	LogFmt:  "json",
	Author:  fallbackAuthor(),
}

func fallbackAuthor() string {
	if u, err := user.Current(); err == nil {
		if u.Username != "" {
			return u.Username
		}
	}
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Base(home)
	}
	return "unknown"
}

func Path() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".fibber/config.yaml"
	}
	return filepath.Join(home, ".fibber", "config.yaml")
}

func Load() (*Config, error) {
	path := Path()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := defaultConfig
			return &cfg, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	path := Path()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func Set(key, value string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	rv := reflect.ValueOf(cfg).Elem()
	t := rv.Type()
	found := false
	for i := 0; i < rv.NumField(); i++ {
		yamlTag := t.Field(i).Tag.Get("yaml")
		if yamlTag == key {
			if yamlTag == "home" {
				absPath, err := filepath.Abs(value)
				if err == nil {
					value = absPath
				}
			}
			rv.Field(i).SetString(value)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("unknown config key: %s", key)
	}
	return Save(cfg)
}

func Get(key string) (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", err
	}
	rv := reflect.ValueOf(cfg).Elem()
	t := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		yamlTag := t.Field(i).Tag.Get("yaml")
		if yamlTag == key {
			return rv.Field(i).String(), nil
		}
	}
	return "", fmt.Errorf("unknown config key: %s", key)
}

func Describe() ([]string, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	rv := reflect.ValueOf(cfg).Elem()
	t := rv.Type()
	var output []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		yamlTag := field.Tag.Get("yaml")
		descTag := field.Tag.Get("config")
		parts := parseTag(descTag)
		value := rv.Field(i).String()
		if strings.TrimSpace(value) == "" && parts["default"] != "" {
			value = parts["default"]
		}
		desc := parts["desc"]
		output = append(output, fmt.Sprintf("  %s = %s\n    → %s", yamlTag, value, desc))
	}
	return output, nil
}

func DescribeJSON() ([]byte, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	rv := reflect.ValueOf(cfg).Elem()
	t := rv.Type()
	out := map[string]map[string]string{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		yamlKey := field.Tag.Get("yaml")
		tags := parseTag(field.Tag.Get("config"))
		val := rv.Field(i).String()
		if val == "" {
			val = tags["default"]
		}
		out[yamlKey] = map[string]string{
			"value":   val,
			"desc":    tags["desc"],
			"default": tags["default"],
		}
	}
	return json.MarshalIndent(out, "", "  ")
}

func parseTag(tag string) map[string]string {
	out := make(map[string]string)
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) == 2 {
			out[kv[0]] = kv[1]
		}
	}
	return out
}
