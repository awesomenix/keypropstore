package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/awesomenix/keypropstore/app"
)

func compareWithExpected(cfg, expectedcfg *app.Config) error {
	if cfg.Port != "8080" {
		return fmt.Errorf("Port %s doesnt match expected 8080", cfg.Port)
	}

	if len(cfg.Stores) != len(expectedcfg.Stores) {
		return fmt.Errorf("Found %d store(s) doesnt match expected %d store(s)", len(cfg.Stores), len(expectedcfg.Stores))
	}

	for i := 0; i < len(expectedcfg.Stores); i++ {
		if cfg.Stores[i].Name != expectedcfg.Stores[i].Name {
			return fmt.Errorf("Store Name %s doesnt match expected %s", cfg.Stores[i].Name, expectedcfg.Stores[i].Name)
		}
		if len(expectedcfg.Stores[i].Backup) > 0 {
			if cfg.Stores[i].Backup != expectedcfg.Stores[i].Backup {
				return fmt.Errorf("Store Backup %s doesnt match expected %s", cfg.Stores[i].Backup, expectedcfg.Stores[i].Backup)
			}
		}
		if len(expectedcfg.Stores[i].Backupdir) > 0 {
			if cfg.Stores[i].Backupdir != expectedcfg.Stores[i].Backupdir {
				return fmt.Errorf("Store BackupDir %s doesnt match expected %s", cfg.Stores[i].Backupdir, expectedcfg.Stores[i].Backupdir)
			}
		}
		if len(expectedcfg.Stores[i].AggregateURLs) > 0 {
			if len(cfg.Stores[i].AggregateURLs) != len(expectedcfg.Stores[i].AggregateURLs) {
				return fmt.Errorf("Store AggregateURLs size %d doesnt match expected %d", len(cfg.Stores[i].AggregateURLs), len(expectedcfg.Stores[i].AggregateURLs))
			}
			for j := range expectedcfg.Stores[i].AggregateURLs {
				if cfg.Stores[i].AggregateURLs[j] != expectedcfg.Stores[i].AggregateURLs[j] {
					return fmt.Errorf("Store AggregateURL %s doesnt match expected %s", cfg.Stores[i].AggregateURLs[j], expectedcfg.Stores[i].AggregateURLs[j])
				}
			}
		}
	}

	return nil
}

func TestBasicConfig(t *testing.T) {
	const fileDir string = "./"
	const filePrefix string = "testcfg"
	fileName := fileDir + filePrefix + ".yml"

	cfg := &app.Config{}
	buf := []byte(`
Port : 8080
Stores :
- local:
`)
	err := ioutil.WriteFile(fileName, buf, 0644)
	defer os.Remove(fileName)
	if err != nil {
		t.Error(err)
		return
	}

	err = cfg.Initialize(filePrefix, fileDir)
	if err != nil {
		t.Error(err)
		return
	}

	cfg.Log()

	stores := make([]app.Store, 1)

	stores[0].Name = "local"

	expectedCfg := &app.Config{"8080", stores}

	err = compareWithExpected(cfg, expectedCfg)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMultipleStoreConfig(t *testing.T) {
	const fileDir string = "./"
	const filePrefix string = "testcfg"
	fileName := fileDir + filePrefix + ".yml"

	cfg := &app.Config{}
	buf := []byte(`
Port : 8080
Stores :
- local:
- second:
- third:
`)
	err := ioutil.WriteFile(fileName, buf, 0644)
	defer os.Remove(fileName)
	if err != nil {
		t.Error(err)
		return
	}

	err = cfg.Initialize(filePrefix, fileDir)
	if err != nil {
		t.Error(err)
		return
	}

	cfg.Log()

	stores := make([]app.Store, 3)

	stores[0].Name = "local"
	stores[1].Name = "second"
	stores[2].Name = "third"

	expectedCfg := &app.Config{"8080", stores}

	err = compareWithExpected(cfg, expectedCfg)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestFullConfig(t *testing.T) {
	const fileDir string = "./"
	const filePrefix string = "testcfg"
	fileName := fileDir + filePrefix + ".yml"

	cfg := &app.Config{}
	buf := []byte(`
Port : 8080
Stores :
- local:
    Backup: BoltDB
    BackupDir: ./boltdb
- second:
    Backup: BoltDB
    BackupDir: ./boltdb
- third:
    Backup: BoltDB
    BackupDir: ./boltdb
    Aggregate:
        - URL1
        - URL2
`)
	err := ioutil.WriteFile(fileName, buf, 0644)
	defer os.Remove(fileName)
	if err != nil {
		t.Error(err)
		return
	}

	err = cfg.Initialize(filePrefix, fileDir)
	if err != nil {
		t.Error(err)
		return
	}

	cfg.Log()

	stores := make([]app.Store, 3)

	stores[0].Name = "local"
	stores[0].Backup = "BoltDB"
	stores[0].Backupdir = "./boltdb"
	stores[1].Name = "second"
	stores[1].Backup = "BoltDB"
	stores[1].Backupdir = "./boltdb"
	stores[2].Name = "third"
	stores[2].Backup = "BoltDB"
	stores[2].Backupdir = "./boltdb"
	stores[2].AggregateURLs = []string{"URL1", "URL2"}

	expectedCfg := &app.Config{"8080", stores}

	err = compareWithExpected(cfg, expectedCfg)
	if err != nil {
		t.Error(err)
		return
	}
}
