package zgrab2

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// RunScans prepares a scan Monitor, and run all the registered scans.
func RunScans() (err error) {

	wg := sync.WaitGroup{}
	monitor := MakeMonitor(1, &wg)
	monitor.Callback = func(_ string) {
		dumpHeapProfile()
	}
	start := time.Now()
	// log.Infof("started grab at %s", start.Format(time.RFC3339))
	Process(monitor)
	end := time.Now()
	// log.Infof("finished grab at %s", end.Format(time.RFC3339))
	monitor.Stop()
	wg.Wait()
	s := Summary{
		StatusesPerModule: monitor.GetStatuses(),
		StartTime:         start.Format(time.RFC3339),
		EndTime:           end.Format(time.RFC3339),
		Duration:          end.Sub(start).String(),
	}
	enc := json.NewEncoder(GetMetaFile())
	if err := enc.Encode(&s); err != nil {
		// log.Fatalf("unable to write summary: %s", err.Error())
	}
	return
}

// StartStats - Calls all functions need to profile the scan process
func StartStats() {
	startCPUProfile()
}

// DumpStats - To be called in a defer statement, will
// dump all data related to scan process profiling
func DumpStats() {
	stopCPUProfile()
	dumpHeapProfile()
}

// Get the value of the ZGRAB2_MEMPROFILE variable (or the empty string).
// This may include {TIMESTAMP} or {NANOS}, which should be replaced using
// getFormattedFile().
func getMemProfileFile() string {
	return os.Getenv("ZGRAB2_MEMPROFILE")
}

// Get the value of the ZGRAB2_CPUPROFILE variable (or the empty string).
// This may include {TIMESTAMP} or {NANOS}, which should be replaced using
// getFormattedFile().
func getCPUProfileFile() string {
	return os.Getenv("ZGRAB2_CPUPROFILE")
}

// Replace instances in formatString of {TIMESTAMP} with when formatted as
// YYYYMMDDhhmmss, and {NANOS} as the decimal nanosecond offset.
func getFormattedFile(formatString string, when time.Time) string {
	timestamp := when.Format("20060102150405")
	nanos := fmt.Sprintf("%d", when.Nanosecond())
	ret := strings.Replace(formatString, "{TIMESTAMP}", timestamp, -1)
	ret = strings.Replace(ret, "{NANOS}", nanos, -1)
	return ret
}

// If memory profiling is enabled (ZGRAB2_MEMPROFILE is not empty), perform a GC
// then write the heap profile to the profile file.
func dumpHeapProfile() {
	if file := getMemProfileFile(); file != "" {
		now := time.Now()
		fullFile := getFormattedFile(file, now)
		f, err := os.Create(fullFile)
		if err != nil {
			log.Fatal("could not create heap profile: ", err)
		}
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write heap profile: ", err)
		}
		f.Close()
	}
}

// If CPU profiling is enabled (ZGRAB2_CPUPROFILE is not empty), start tracking
// CPU profiling in the configured file. Caller is responsible for invoking
// stopCPUProfile() when finished.
func startCPUProfile() {
	if file := getCPUProfileFile(); file != "" {
		now := time.Now()
		fullFile := getFormattedFile(file, now)
		f, err := os.Create(fullFile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
	}
}

// If CPU profiling is enabled (ZGRAB2_CPUPROFILE is not empty), stop profiling
// CPU usage.
func stopCPUProfile() {
	if getCPUProfileFile() != "" {
		pprof.StopCPUProfile()
	}
}
