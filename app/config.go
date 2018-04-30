package app

import (
	"fmt"

	"github.com/spf13/viper"
)

// Store details for this application
// primary store is always inmemory with
// optional backup store along with backup directory
// also aggregte urls for aggregating multiple stores into the primary
type Store struct {
	Name          string
	Backup        string
	Backupdir     string
	AggregateURLs []string
}

// Config context for this application
type Config struct {
	Port   string
	Stores []Store
}

// Config will look like this
// Port : 8080
// Stores :
//   - Machines :
//	     Backup : BoltDB
//       BackupDir : ./boltdb
//   - GlobalAggregateMachines :
//	     Backup : BoltDB
//       BackupDir : ./boltdb
//		 Aggregate:
//			- URL1
//			- URL2
//			...

// Initialize AppConfig with sample yaml file provided above
func (cfg *Config) Initialize(name, dir string) error {
	viper.SetConfigName(name) // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir) // path to look for the config file in
	viper.AddConfigPath(".") // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return err
	}

	viper.SetDefault("Port", "8080")
	cfg.Port = viper.GetString("Port")

	stores := viper.Get("Stores")
	for _, istorevalues := range stores.([]interface{}) {
		for storename, istoresettings := range istorevalues.(map[interface{}]interface{}) {
			var store Store

			store.Name = storename.(string)

			if istoresettings != nil {
				setting := istoresettings.(map[interface{}]interface{})
				if backup, ok := setting["Backup"]; ok {
					store.Backup = backup.(string)
				}

				if backupdir, ok := setting["BackupDir"]; ok {
					store.Backupdir = backupdir.(string)
				}

				if aggregate, ok := setting["Aggregate"]; ok {
					for _, aggregateURL := range aggregate.([]interface{}) {
						store.AggregateURLs = append(store.AggregateURLs, aggregateURL.(string))
					}
				}
			}

			cfg.Stores = append(cfg.Stores, store)
		}

	}

	return nil
}

// Log AppConfig mainly used for debugging purposes
func (cfg *Config) Log() {
	fmt.Println("Port:", cfg.Port)

	for _, store := range cfg.Stores {
		fmt.Println("Name:", store.Name)
		fmt.Println("Primary: InMemoryStore")

		if len(store.Backup) > 0 {
			fmt.Println("\tBackup:", store.Backup)
		}
		if len(store.Backupdir) > 0 {
			fmt.Println("\tBackupDir:", store.Backupdir)
		}

		if store.AggregateURLs != nil {
			fmt.Println("\tAggregate:", store.AggregateURLs)
		}
	}
}
