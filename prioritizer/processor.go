package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

const dateLayout = "2006-01-02T15:04:05Z"

var processorLock = &sync.Mutex{}

func prioritizePods(interval int, done chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case <-time.After(time.Duration(interval) * time.Second):
			err := prioritizeLoop()
			if err != nil {
				log.Println(err)
			}
		case <-done:
			wg.Done()
			log.Println("Stopped reconciliation loop.")
			return
		}
	}
}

func prioritizeLoop() error {
	processorLock.Lock()
	defer processorLock.Unlock()
	pods, err := getUnscheduledPods()
	if err != nil {
		return err
	}
	for _, pod := range pods {
		//gather variables
		labels := pod.Metadata.Labels

		expectedDuration, _ := labels["expectedDuration"]
		ed, _ := strconv.Atoi(expectedDuration)
		edf := float64(ed)

		valueMult, _ := labels["valueMult"]
		vm, _ := strconv.Atoi(valueMult)
		vmf := float64(vm)

		t, err := time.Parse(dateLayout, pod.Metadata.CreationTimestamp)
		if err != nil {
			fmt.Println(err)
		}
		waitTime := time.Since(t).Seconds()

		//calculate priority
		priority := int(vmf * (1 + waitTime/edf))

		//update pod with priority
		fmt.Printf("updating pod %s with priority %s\n", pod.Metadata.Name, strconv.Itoa(priority))
		patchErr := patchPriority(&pod, strconv.Itoa(priority))
		if patchErr != nil {
			return patchErr
		}
	}
	return nil
}
