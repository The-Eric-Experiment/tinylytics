package db

import (
	"fmt"
	"sync"
	"tinylytics/helpers"
)

// DatabaseManager manages singleton database instances per domain
// The mutex is only needed during initialization to prevent race conditions
// when creating new database instances. Once created, each Database has its
// own mutex for operations.
type DatabaseManager struct {
	databases map[string]*Database
	mu        sync.Mutex // Only for initialization, not operations
}

var manager = &DatabaseManager{
	databases: make(map[string]*Database),
}

// GetDatabase returns a singleton database instance for the given file path
// It creates and initializes the database if it doesn't exist
func GetDatabase(file string) (*Database, error) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	// Check if database already exists
	if db, exists := manager.databases[file]; exists {
		return db, nil
	}

	// Create new database instance
	db := &Database{}
	db.Connect(file)
	db.Initialize()

	// Store in map
	manager.databases[file] = db

	return db, nil
}

// GetDatabaseByDomain returns a database instance for a domain
func GetDatabaseByDomain(domain string) (*Database, error) {
	dbFile, err := helpers.GetDatabaseFileName(domain)
	if err != nil {
		return nil, err
	}
	return GetDatabase(dbFile)
}

// InitializeAllDatabases initializes all databases for configured websites
func InitializeAllDatabases(domains []string) error {
	for _, domain := range domains {
		_, err := GetDatabaseByDomain(domain)
		if err != nil {
			return fmt.Errorf("failed to initialize database for domain %s: %w", domain, err)
		}
	}
	return nil
}

// CloseAll closes all database connections gracefully
func CloseAll() {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	for file, db := range manager.databases {
		db.Close()
		delete(manager.databases, file)
	}
}

