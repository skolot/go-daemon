package daemon

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"os"
	"testing"
)

var (
	filename                = os.TempDir() + "/test.lock"
	fileperm    os.FileMode = 0644
	invalidname             = "/x/y/unknown"
)

func TestCreatePidFile(test *testing.T) {
	if _, err := CreatePidFile(invalidname, fileperm); err == nil {
		test.Fatal("CreatePidFile(): Error was not detected on invalid name")
	}

	lock, err := CreatePidFile(filename, fileperm)
	if err != nil {
		test.Fatal(err)
	}
	defer lock.Remove()

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		test.Fatal(err)
	}
	if string(data) != fmt.Sprint(os.Getpid()) {
		test.Fatal("pids not equal")
	}

	file, err := os.OpenFile(filename, os.O_RDONLY, fileperm)
	if err != nil {
		test.Fatal(err)
	}
	if err = NewLockFile(file).WritePid(); err == nil {
		test.Fatal("WritePid(): Error was not detected on invalid permissions")
	}
}

func TestNewLockFile(test *testing.T) {
	lock := NewLockFile(os.NewFile(1001, ""))
	err := lock.Remove()
	if err == nil {
		test.Fatal("Remove(): Error was not detected on invalid fd")
	}
	err = lock.WritePid()
	if err == nil {
		test.Fatal("WritePid(): Error was not detected on invalid fd")
	}
}

func TestGetFdName(test *testing.T) {
	name, err := GetFdName(0)
	if err != nil {
		test.Error(err)
	} else {
		devnull := "/dev/null"

		if runtime.GOOS == "solaris" {
			devnull = "/devices/pseudo/mm@0:null"
		}

		if name != devnull {
			test.Errorf("Filename of fd 0: `%s'", name)
		}
	}

	name, err = GetFdName(1011)
	if err == nil {
		test.Errorf("GetFdName(): Error was not detected on invalid fd, name: `%s'", name)
	}
}

func TestReadPid(test *testing.T) {
	lock, err := CreatePidFile(filename, fileperm)
	if err != nil {
		test.Fatal(err)
	}
	defer lock.Remove()

	pid, err := lock.ReadPid()
	if err != nil {
		test.Fatal("ReadPid(): Unable read pid from file:", err)
	}

	if pid != os.Getpid() {
		test.Fatal("Pid not equal real pid")
	}
}

func TestLockFileLock(test *testing.T) {
	if runtime.GOOS == "solaris" {
		// sensless under solaris, exclusive lock is always unlock for same process 
		return
	}

	lock1, err := OpenLockFile(filename, fileperm)
	if err != nil {
		test.Fatal(err)
	}
	fmt.Println("lock1")
	if err := lock1.Lock(); err != nil {
		test.Fatal(err)
	}
	defer lock1.Remove()

	lock2, err := OpenLockFile(filename, fileperm)
	if err != nil {
		test.Fatal(err)
	}
	fmt.Println("lock2")
	if err := lock2.Lock(); err != ErrWouldBlock {
		test.Fatal("To lock file more than once must be unavailable.")
	}
}
