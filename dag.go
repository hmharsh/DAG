package dag

import (
	"sync"
)

var (
	status        = map[string]int{}
	ErrorMessages = map[string]error{}
	mutex         = &sync.RWMutex{}
)

type Dependencies struct {
	Function    func() error
	Name        string
	Installed   []string
	Uninstalled []string
	Retry       bool
}

var err error

// Init function will Initialize the status map based on input conditions
func Init(conditions []Dependencies) {
	for _, x := range conditions {
		status[x.Name] = 0
	}
}

// Installer will review the status map and trigger next  worker function goroutine
func Installer(conditions []Dependencies) {
	var wg sync.WaitGroup
	for _, x := range conditions {
		if CheckStatus(x.Installed, x.Uninstalled) {
			wg.Add(1)
			go Generator(x.Function, x.Name, x.Retry, &wg, conditions, mutex)
		}
	}
	wg.Wait()
}

// Generator will wrap the core logic of specific worker function with necessary go statements to handel go routines
func Generator(x func() error, name string, retry bool, wg *sync.WaitGroup, conditions []Dependencies, mutex *sync.RWMutex) {
	defer wg.Done()
	if status[name] != 3 {
		mutex.Lock()
		status[name] = 1
		mutex.Unlock()
	}
	err = x()
	if err != nil {
		if retry && status[name] == 1 {
			mutex.Lock()
			status[name] = 3
			mutex.Unlock()
			ErrorMessages["Error in Function "+name+":"] = err
		} else {
			mutex.Lock()
			status[name] = 4
			mutex.Unlock()
			if retry {
				ErrorMessages["Error in Function "+name+" on retry: "] = err
			} else {
				ErrorMessages["Error in Function "+name] = err
			}
		}
	} else {
		mutex.Lock()
		status[name] = 2
		mutex.Unlock()
	}
	Installer(conditions)
}

// Checkstatus uses 'status' map as source of truth and return true
// if charts in installed array are installed already
func CheckStatus(installed, uninstalled []string) bool {
	for _, chart := range installed {
		if status[chart] != 2 {
			return false
		}
	}
	for _, chart := range uninstalled {
		if status[chart] == 1 || status[chart] == 2 || status[chart] == 4 {
			return false
		}
	}
	return true
}

// VerifyStatus is used to verify the final worker execution status at last
func VerifyStatus() bool {
	for _, stat := range status {
		if stat != 2 {
			return false
		}
	}
	return true
}

// FetchError will return the worker function status along with error message map
func FetchError() (map[string]int, map[string]error) {
	return status, ErrorMessages
}
