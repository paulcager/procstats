package procstats

import (
//	"github.com/shirou/gopsutil/cpu"
//	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"github.com/paulcager/rolling"
	"os"
//	"sync/atomic"
	"time"
)

const (
	sampleInterval = time.Second
)

var (
	stopChan chan struct{} = make(chan struct{}, 1)
	myProc   *process.Process
	cpuValues *rolling.Window
	sysCpuValues *rolling.Window
)

func init() {
	var err error
	myProc, err = process.NewProcess(int32(os.Getpid()))
	if err != nil {
		// Should be impossible - the running process must exist.
		panic(err)
	}
	cpuValues = rolling.New(rolling.Max, 60, 60, 24)
	sysCpuValues = rolling.New(rolling.Max, 60, 60, 24)
}

func Start() {
	go sample()
}

func Stop() {
	stopChan <- struct{}{}
}

func sample() {
	t := time.NewTicker(sampleInterval)
	defer t.Stop()

	myTimes, _ := myProc.Times()
	//sysTimes, _ := cpu.Times(false)

	for {
		select {
		case <-t.C:
			newMyTimes, _ := myProc.Times()
			cpu := int64(newMyTimes.User-myTimes.User)
			sys := int64(newMyTimes.System-myTimes.System)
			myTimes = newMyTimes
			cpuValues.Push(cpu)
			sysCpuValues.Push(sys)
		case <-stopChan:
			return
		}
	}
}
