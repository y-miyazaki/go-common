//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		setting      *Setting
		name         string
		wantSlackNil bool
	}{
		{
			name: "json logger without slack",
			setting: &Setting{
				LoggerFormatter:       "json",
				LoggerOut:             "stdout",
				LoggerLevel:           "info",
				SlackOauthAccessToken: "",
			},
			wantSlackNil: true,
		},
		{
			name: "text logger with slack token",
			setting: &Setting{
				LoggerFormatter:       "text",
				LoggerOut:             "stdout",
				LoggerLevel:           "debug",
				SlackOauthAccessToken: "test-token",
			},
			wantSlackNil: false,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := NewConfig(tt.setting)
			require.NotNil(t, cfg)
			require.NotNil(t, cfg.Logger)
			if tt.wantSlackNil {
				require.Nil(t, cfg.SlackClient)
			} else {
				require.NotNil(t, cfg.SlackClient)
			}
		})
	}
}

func TestNewConfig_InvalidSetting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		setting *Setting
		name    string
	}{
		{
			name: "invalid formatter panics",
			setting: &Setting{
				LoggerFormatter: "invalid",
				LoggerOut:       "stdout",
				LoggerLevel:     "info",
			},
		},
		{
			name: "invalid output panics",
			setting: &Setting{
				LoggerFormatter: "json",
				LoggerOut:       "invalid",
				LoggerLevel:     "info",
			},
		},
		{
			name: "invalid level panics",
			setting: &Setting{
				LoggerFormatter: "json",
				LoggerOut:       "stdout",
				LoggerLevel:     "invalid",
			},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Panics(t, func() { NewConfig(tt.setting) })
		})
	}
}
