package app

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func testBasicUpdateQuery(buf []byte, t *testing.T) {
	const fileDir string = "./"
	const filePrefix string = "config"
	fileName := fileDir + filePrefix + ".yml"

	defer os.Remove(fileName)
	if err := ioutil.WriteFile(fileName, buf, 0644); err != nil {
		t.Error(err)
		return
	}

	ctx := CreateContext()
	defer DeleteContext(ctx)

	postBuf := []byte(`{"m1": {"num": "6.13","strs": "a","key1": "b"}, "m2": {"num": "6.13","key1": "bddd"}}`)
	respU, postErr := http.Post("http://127.0.0.1:8080/v1/store/local/update", "application/json", bytes.NewBuffer(postBuf))
	if postErr != nil {
		t.Error(postErr)
		return
	}

	if respU.StatusCode != 200 {
		t.Errorf("Update returned %d, expected success with 200, error: %s", respU.StatusCode, respU.Status)
		return
	}

	queryBuf := []byte(`{"num": "6.13"}`)
	resp, perr := http.Post("http://127.0.0.1:8080/v1/store/local/query", "application/json", bytes.NewBuffer(queryBuf))

	if perr != nil {
		t.Error(perr)
		return
	}

	if resp.StatusCode != 200 {
		t.Errorf("Update returned %d, expected success with 200, error: %s", resp.StatusCode, resp.Status)
		return
	}

	defer resp.Body.Close()

	bodyBytes, rerr := ioutil.ReadAll(resp.Body)
	if rerr != nil {
		t.Error(rerr)
		return
	}

	var res []string
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		t.Error(err)
		return
	}

	expected := []string{"m1", "m2"}

	t.Log("Store returned", res, "Expect", expected)

	for _, expKey := range expected {
		found := false
		for _, key := range res {
			if key == expKey {
				found = true
			}
		}
		if !found {
			t.Errorf("Expected %v", expKey)
			return
		}
	}
}

func testNoUpdateQuery(buf []byte, t *testing.T) {
	const fileDir string = "./"
	const filePrefix string = "config"
	fileName := fileDir + filePrefix + ".yml"

	defer os.Remove(fileName)
	if err := ioutil.WriteFile(fileName, buf, 0644); err != nil {
		t.Error(err)
		return
	}

	ctx := CreateContext()
	defer DeleteContext(ctx)

	queryBuf := []byte(`{"num": "6.13"}`)
	resp, perr := http.Post("http://127.0.0.1:8080/v1/store/local/query", "application/json", bytes.NewBuffer(queryBuf))

	if perr != nil {
		t.Error(perr)
		return
	}

	if resp.StatusCode != 200 {
		t.Errorf("Update returned %d, expected success with 200, error: %s", resp.StatusCode, resp.Status)
		return
	}

	defer resp.Body.Close()

	bodyBytes, rerr := ioutil.ReadAll(resp.Body)
	if rerr != nil {
		t.Error(rerr)
		return
	}

	var res []string
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		t.Error(err)
		return
	}

	expected := []string{"m1", "m2"}

	t.Log("Store returned", res, "Expect", expected)

	for _, expKey := range expected {
		found := false
		for _, key := range res {
			if key == expKey {
				found = true
			}
		}
		if !found {
			t.Errorf("Expected %v", expKey)
			return
		}
	}
}

func TestInMemoryUpdateQuer(t *testing.T) {

	buf := []byte(`
Port : 8080
Stores :
- local:
`)

	testBasicUpdateQuery(buf, t)
}

func TestBadgerDBUpdateQuery(t *testing.T) {

	buf := []byte(`
Port : 8080
Stores :
- local:
    Backup: BadgerDB
    BackupDir: ./badgerdbtest
`)

	defer os.RemoveAll("./badgerdbtest")
	testBasicUpdateQuery(buf, t)
}

func TestBoltDBUpdateQuery(t *testing.T) {

	buf := []byte(`
Port : 8080
Stores :
- local:
    Backup: BoltDB
    BackupDir: ./boltdbtest
`)

	defer os.RemoveAll("./boltdbtest")
	testBasicUpdateQuery(buf, t)
}

func TestBadgerDBBackup(t *testing.T) {

	buf := []byte(`
Port : 8080
Stores :
- local:
    Backup: BadgerDB
    BackupDir: ./badgerdbtest
`)

	defer os.RemoveAll("./badgerdbtest")
	testBasicUpdateQuery(buf, t)
	testNoUpdateQuery(buf, t)
}

func TestBoltDBBackup(t *testing.T) {

	buf := []byte(`
Port : 8080
Stores :
- local:
    Backup: BoltDB
    BackupDir: ./boltdbtest
`)

	defer os.RemoveAll("./boltdbtest")
	testBasicUpdateQuery(buf, t)
	testNoUpdateQuery(buf, t)
}
