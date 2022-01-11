package configuration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadToStructs(t *testing.T) {
	tests := []struct {
		name    string
		wantOut Env
	}{
		{
			name: "test",
			wantOut: Env{
				DataBase: DataBase{
					Type:     "postgres",
					Host:     "localhost",
					Port:     5432,
					Name:     "gochat",
					User:     "postgres",
					Password: "pass",
				},
				Redis: Redis{
					Host: "localhost",
					Port: 6379,
				},
			},
		},
	}

	actual, err := InitEnv("../../tests/test_configs/test_config.yml")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, err)
			require.NotNil(t, actual)
			assert.Equal(t, &tt.wantOut, actual)
		})
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := InitEnv("../../configs/test_configs/missing_config.yml")
	if err == nil {
		t.Fatal("Loaded non-existent file")
	}
}

func ExampleInitEnv() {
	env, _ := InitEnv("../../tests/test_configs/test_config.yml")
	fmt.Println(env.DataBase.Host)
	fmt.Println(env.DataBase.Port)
	// Output:
	// localhost
	// 5432
}
