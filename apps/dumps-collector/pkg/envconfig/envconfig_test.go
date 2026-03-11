package envconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// requiredEnvForDockerfileScriptsDeployment is the set of env var names that may be required
// when deploying the Dockerfile.scripts image. Irrespective of PVMountPath, the scripts image
// must deploy without Postgres params (DIAG_POSTGRES_HOST, DIAG_POSTGRES_PORT,
// DIAG_POSTGRES_USERNAME, DIAG_POSTGRES_PASSWORD, DIAG_DB_NAME).
var requiredEnvForDockerfileScriptsDeployment = []string{} // scripts image needs no mandatory env for PG

// TestDockerfileScriptsDeploysWithoutPostgresParams verifies that when deploying the Dockerfile.scripts
// image, the deployment does not require any of the Postgres params. This is independent of PVMountPath:
// the scripts image must deploy without DIAG_POSTGRES_HOST, DIAG_POSTGRES_PORT, DIAG_POSTGRES_USERNAME,
// DIAG_POSTGRES_PASSWORD, and DIAG_DB_NAME.
func TestDockerfileScriptsDeploysWithoutPostgresParams(t *testing.T) {
	pgParams := PostgresParamNames()
	require.Len(t, pgParams, 5, "PostgresParamNames must define exactly the five PG params")

	expected := []string{"DIAG_POSTGRES_HOST", "DIAG_POSTGRES_PORT", "DIAG_POSTGRES_USERNAME", "DIAG_POSTGRES_PASSWORD", "DIAG_DB_NAME"}
	assert.ElementsMatch(t, expected, pgParams)

	// Dockerfile.scripts deployment must not require any of the Postgres params
	requiredForScripts := make(map[string]bool)
	for _, name := range requiredEnvForDockerfileScriptsDeployment {
		requiredForScripts[name] = true
	}
	for _, pgName := range pgParams {
		assert.False(t, requiredForScripts[pgName],
			"Dockerfile.scripts deployment must deploy without Postgres param %q; do not require it when using the scripts image", pgName)
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
