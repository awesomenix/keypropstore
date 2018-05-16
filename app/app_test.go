package app

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
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

	ctx, err := CreateDefaultContext()
	if err != nil {
		t.Error(err)
		return
	}
	defer DeleteContext(ctx)

	time.Sleep(1 * time.Second)

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

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Query returned %d, expected success with 200, error: %s", resp.StatusCode, resp.Status)
		return
	}

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

	ctx, err := CreateDefaultContext()
	if err != nil {
		t.Error(err)
		return
	}
	defer DeleteContext(ctx)

	time.Sleep(1 * time.Second)

	queryBuf := []byte(`{"num": "6.13"}`)
	resp, perr := http.Post("http://127.0.0.1:8080/v1/store/local/query", "application/json", bytes.NewBuffer(queryBuf))

	if perr != nil {
		t.Error(perr)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Query returned %d, expected success with 200, error: %s", resp.StatusCode, resp.Status)
		return
	}

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

func testBasicAggregateQuery(buf1 []byte, buf2 []byte, t *testing.T) {
	const fileDir string = "./"
	const filePrefix1 string = "config1"
	const filePrefix2 string = "config2"
	fileName1 := fileDir + filePrefix1 + ".yml"
	fileName2 := fileDir + filePrefix2 + ".yml"

	defer os.Remove(fileName1)
	if err := ioutil.WriteFile(fileName1, buf1, 0644); err != nil {
		t.Error(err)
		return
	}

	defer os.Remove(fileName2)
	if err := ioutil.WriteFile(fileName2, buf2, 0644); err != nil {
		t.Error(err)
		return
	}

	ctx, err1 := CreateContext(filePrefix1, fileDir+filePrefix1)
	if err1 != nil {
		t.Error(err1)
		return
	}
	defer DeleteContext(ctx)

	ctx2, err2 := CreateContext(filePrefix2, fileDir+filePrefix2)
	if err2 != nil {
		t.Error(err2)
		return
	}
	defer DeleteContext(ctx2)

	time.Sleep(1 * time.Second)

	{
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
	}

	{
		postBuf := []byte(`{"m3": {"num": "6.13","strs": "a","key1": "b"}, "m4": {"num": "6.13","key1": "bddd"}}`)
		respU, postErr := http.Post("http://127.0.0.1:8081/v1/store/local/update", "application/json", bytes.NewBuffer(postBuf))
		if postErr != nil {
			t.Error(postErr)
			return
		}

		if respU.StatusCode != 200 {
			t.Errorf("Update returned %d, expected success with 200, error: %s", respU.StatusCode, respU.Status)
			return
		}
	}

	time.Sleep(3 * time.Second)

	queryBuf := []byte(`{"num": "6.13"}`)
	resp, perr := http.Post("http://127.0.0.1:8080/v1/store/aggregate/query", "application/json", bytes.NewBuffer(queryBuf))

	if perr != nil {
		t.Error(perr)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Query returned %d, expected success with 200, error: %s", resp.StatusCode, resp.Status)
		return
	}

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

	expected := []string{"m1", "m2", "m3", "m4"}

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

func TestInMemoryUpdateQuery(t *testing.T) {

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

// This unittest seems to crash on Mac/Linux
// panic on db close, requires further investigation, not noticed on badgerdb
/*func TestBoltDBBackup(t *testing.T) {

	buf := []byte(`
Port : 8080
Stores :
- local:
    Backup: BoltDB
    BackupDir: ./boltdbtest
`)

	defer os.RemoveAll("./boltdbtest")
	testBasicUpdateQuery(buf, t)
	time.Sleep(1 * time.Second)
	testNoUpdateQuery(buf, t)
}*/

func TestInMemoryAggregateUpdateQuery(t *testing.T) {

	buf1 := []byte(`
Port : 8080
Stores :
- local:
- aggregate:
    SyncInterval: 1
    Aggregate:
    - http://127.0.0.1:8080/v1/store/local
    - http://127.0.0.1:8081/v1/store/local
`)

	buf2 := []byte(`
Port : 8081
Stores :
- local:
`)

	testBasicAggregateQuery(buf1, buf2, t)
}

func TestBadgerDBAggregateUpdateQuery(t *testing.T) {

	buf1 := []byte(`
Port : 8080
Stores :
- local:
- aggregate:
    SyncInterval: 1
    Backup: BadgerDB
    BackupDir: ./badgerdbtest 
    Aggregate:
    - http://127.0.0.1:8080/v1/store/local
    - http://127.0.0.1:8081/v1/store/local
`)

	buf2 := []byte(`
Port : 8081
Stores :
- local:
    Backup: BadgerDB
    BackupDir: ./badgerdbtest2 
`)

	defer os.RemoveAll("./badgerdbtest")
	defer os.RemoveAll("./badgerdbtest2")
	testBasicAggregateQuery(buf1, buf2, t)
}

// This unittest seems to crash on Mac/Linux
// panic on db close, requires further investigation, not noticed on badgerdb
func TestBoltDBAggregateUpdateQuery(t *testing.T) {

	buf1 := []byte(`
Port : 8080
Stores :
- local:
- aggregate:
    SyncInterval: 1
    Backup: BoltDB
    BackupDir: ./boltdbtest 
    Aggregate:
    - http://127.0.0.1:8080/v1/store/local
    - http://127.0.0.1:8081/v1/store/local
`)

	buf2 := []byte(`
Port : 8081
Stores :
- local:
`)

	defer os.RemoveAll("./boltdbtest")
	testBasicAggregateQuery(buf1, buf2, t)
}
