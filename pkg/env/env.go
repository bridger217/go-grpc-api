package env

import "os"

// Environment variable names
const (
	ENV_FIREBASE_JSON = "FIREBASE_JSON" // holds a base64 encoded json string
	ENV_DB_USER       = "DB_USER"
	ENV_DB_PASSWORD   = "DB_PASSWORD"
	ENV_DB_IP_ADDR    = "DB_IP_ADDR"
)

type EnvManager struct{}

func (m *EnvManager) GetFirebaseJson() string {
	return os.Getenv(ENV_FIREBASE_JSON)
}

func (m *EnvManager) GetDbUser() string {
	return os.Getenv(ENV_DB_USER)
}

func (m *EnvManager) GetDbPassword() string {
	return os.Getenv(ENV_DB_PASSWORD)
}

func (m *EnvManager) GetDbIpAddr() string {
	return os.Getenv(ENV_DB_IP_ADDR)
}
