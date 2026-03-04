package config

type Config struct {
	// Add your configuration fields here
	// For example:
	// Port int
	// DBHost string
	// DBPort int
	// DBUser string
	// DBPassword string
	// DBName string
}

func New() *Config {
	return &Config{
		// Initialize your configuration fields here
		// For example:
		// Port: 8080,
		// DBHost: "localhost",
		// DBPort: 5432,
		// DBUser: "user",
		// DBPassword: "password",
		// DBName: "mydb",
	}
}
