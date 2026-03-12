package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Path string

const (
	DotEnv Path = ".env"
)

const (
	port            = "PORT"
	logLevel        = "LOG_LEVEL"
	dbHost          = "DB_HOST"
	dbPort          = "DB_PORT"
	dbName          = "DB_NAME"
	dbUser          = "DB_USER"
	dbPassword      = "DB_PASSWORD"
	dbSSLMode       = "DB_SSLMODE"
	shutdownTimeout = "SHUTDOWN_TIMEOUT"
	jwtSecret       = "JWT_SECRET"
)

func Init(path Path) error {
	if err := godotenv.Load(string(path)); err != nil {
		logrus.WithField("path", path).Debug(".env file not loaded, continuing with system environment")
	} else {
		logrus.WithField("path", path).Info(".env file loaded successfully")
	}

	return checkENV()
}

func checkENV() error {
	vars := []string{
		port, logLevel,
		dbHost, dbPort, dbName, dbUser, dbPassword, dbSSLMode,
		shutdownTimeout, jwtSecret,
	}

	for _, v := range vars {
		if _, exists := os.LookupEnv(v); !exists {
			logrus.WithField("variable", v).Error("Environment variable is missing")
			return fmt.Errorf("environment variable %s is missing", v)
		}
	}

	if _, err := getDuration(shutdownTimeout); err != nil {
		logrus.WithFields(logrus.Fields{
			"variable": shutdownTimeout,
			"error":    err,
		}).Error("Invalid duration format")
		return err
	}

	if _, err := getInt(dbPort); err != nil {
		logrus.WithFields(logrus.Fields{
			"variable": dbPort,
			"error":    err,
		}).Error("Invalid integer format")
		return err
	}

	logrus.Info("All environment variables validated successfully")
	return nil
}

func Port() int {
	val, _ := getInt(port)
	return val
}

func LogLevel() string {
	return os.Getenv(logLevel)
}

func DBHost() string {
	return os.Getenv(dbHost)
}

func DBPort() int {
	val, _ := getInt(dbPort)
	return val
}

func DBName() string {
	return os.Getenv(dbName)
}

func DBUser() string {
	return os.Getenv(dbUser)
}

func DBPassword() string {
	return os.Getenv(dbPassword)
}

func DBSSLMode() string {
	return os.Getenv(dbSSLMode)
}

func ShutdownTimeout() time.Duration {
	val, _ := getDuration(shutdownTimeout)
	return val
}

func JWTSecret() []byte {
	return []byte(os.Getenv(jwtSecret))
}

func DSN() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		DBHost(), DBPort(), DBName(), DBUser(), DBPassword(), DBSSLMode())
}

func getDuration(key string) (time.Duration, error) {
	s := os.Getenv(key)
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("variable %s: %w", key, err)
	}
	return time.Duration(n) * time.Second, nil
}

func getInt(key string) (int, error) {
	s := os.Getenv(key)
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("variable %s: %w", key, err)
	}
	return n, nil
}
