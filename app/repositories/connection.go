package repositories

import "ticket-sync-mvc/framework/database"

// Open initializes the shared connection and the application auth schema.
func Open(path string) error {
	if err := database.Open(path); err != nil {
		return err
	}
	Migrate()
	return nil
}
