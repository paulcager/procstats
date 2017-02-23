package procstats

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	//	"github.com/shirou/gopsutil/cpu"
	//	"github.com/shirou/gopsutil/mem"
	"github.com/paulcager/rolling"
	"github.com/shirou/gopsutil/process"
)

const (
	sampleInterval = time.Second
)

var _ = fmt.Errorf

// TODO - procstats.CpuMon().Serve("/cpu").Start()

var (
	stopChan     chan struct{} = make(chan struct{}, 1)
	myProc       *process.Process
	cpuValues    *rolling.Window
	sysCpuValues *rolling.Window

	lastMinute atomic.Value // of []Point
	lastHour   atomic.Value
	lastDay    atomic.Value
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

	lastMinute.Store(make([]rolling.Point, 60))
	lastHour.Store(make([]rolling.Point, 60))
	lastDay.Store(make([]rolling.Point, 24))
}

func Start() {
	go sample()
}

func Stop() {
	stopChan <- struct{}{}
}

func LastMinute() []rolling.Point {
	return lastMinute.Load().([]rolling.Point)
}
func LastHour() []rolling.Point {
	return lastHour.Load().([]rolling.Point)
}
func LastDay() []rolling.Point {
	return lastDay.Load().([]rolling.Point)
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
			cpu := int64(100 * (newMyTimes.User - myTimes.User))
			sys := int64(100 * (newMyTimes.System - myTimes.System))
			myTimes = newMyTimes
			cpuValues.Push(cpu)
			sysCpuValues.Push(sys)
			lastMinute.Store(cpuValues.Get(0))
			lastHour.Store(cpuValues.Get(1))
			lastDay.Store(cpuValues.Get(2))
		case <-stopChan:
			return
		}
	}
}
