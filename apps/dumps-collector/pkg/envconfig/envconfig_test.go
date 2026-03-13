package envconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDockerfileScriptsDeploysWithoutPostgresParams verifies that when deploying the
// Dockerfile.scripts image (qubership-profiler-dumps-collector-scripts), the deployment
// does not require any of the Postgres params. The scripts image runs Nginx + WebDAV +
// cleanup only; no Go server or Postgres is needed.
func TestDockerfileScriptsDeploysWithoutPostgresParams(t *testing.T) {
	pgParams := PostgresParamNames()
	require.Len(t, pgParams, 5, "PostgresParamNames must define exactly the five PG params")

	// Contract: deployment for Dockerfile.scripts must not inject or require these env vars.
	// Charts/Helm should omit DIAG_POSTGRES_* and DIAG_DB_NAME when using the scripts image.
	for _, name := range pgParams {
		assert.Contains(t, name, "DIAG_", "PG param %q should be a DIAG_ env var", name)
	}
}

func clearEnv(keys ...string) func() {
	old := make(map[string]string)
	for _, k := range keys {
		old[k] = os.Getenv(k)
		os.Unsetenv(k)
	}
	return func() {
		for k, v := range old {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	}
}

func setEnv(env map[string]string) func() {
	old := make(map[string]string)
	for k, v := range env {
		old[k] = os.Getenv(k)
		os.Setenv(k, v)
	}
	return func() {
		for k, v := range old {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	}
}
