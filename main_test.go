package main

import (
	"crypto/sha256"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	LoadLimit = 150

	// Singleport portsize have to be '1'.
	SINGLEWorkerSize = 4
	SINGLEPortSize   = 1
	SINGLEQueueSize  = 10

	// Multiport workersize and Multiport portsize have to be equal!
	MULTIWorkerSize = 4
	MULTIPortSize   = 4
	MULTIQueueSize  = 10

	RANDTimeStart = 200_000
	RANDTimeEnd   = 600_000
	RANDTimeType  = time.Microsecond
)

func TestSinglePort(t *testing.T) {
	singleport := NewQueueDesign(SINGLEWorkerSize, SINGLEPortSize, SINGLEQueueSize)

	wg := &sync.WaitGroup{}
	wg.Add(singleport.PortSize)
	wg.Add(singleport.WorkerSize)

	// Build wires.
	for k := 0; k < singleport.PortSize; k++ {
		ch := make(chan []byte, singleport.OneChannelSize)
		singleport.Channels[k] = ch
	}
	for i := 0; i < singleport.PortSize; i++ {
		go func(i int, wg *sync.WaitGroup, singleport *QueueDesign) {
			fun := NewWorkLoad(func() {
				loadcount := 0
				for {
					// Fugazi data stuff.
					h := sha256.New()
					facelessvoid := time.Now()
					h.Write([]byte(facelessvoid.String()))
					timeislocked := h.Sum([]byte{1})
					latency := time.Now()
					// Produce some data to workers.
					singleport.Channels[i] <- timeislocked
					latencyDuration := time.Since(latency).Milliseconds()
					atomic.AddInt64(&singleport.SumLatency, latencyDuration)
					loadcount++
					if loadcount >= LoadLimit/singleport.PortSize {
						break
					}
				}
				wg.Done()
				close(singleport.Channels[0])
			})
			fun.worker()

		}(i, wg, singleport)
	}
	<-time.After(time.Second * 2)
	t.Log("Singleport is started..")
	start := time.Now()

	for m := 0; m < singleport.WorkerSize; m++ {
		go func(m int, wg *sync.WaitGroup, singleport *QueueDesign) {
			fun := NewWorkLoad(func() {
				countofconsume := 0
			out:
				for {
					select {
					// Consume data is from corresponding port.
					case data, ok := <-singleport.Channels[0]:
						if !ok {
							log.Printf("[SINGLEPORT][TEST][WORKER:%v] Consumer is closed. Count of consumed data: %v", m, countofconsume)
							wg.Done()
							break out
						}
						countofconsume++
						// Do something.
						for i := 0; i < 1; i++ {
							_ = data
							source := rand.NewSource(time.Now().UnixNano())
							r := rand.New(source)
							randomNumber := r.Intn(RANDTimeEnd) + RANDTimeStart
							ch := time.After(time.Duration(randomNumber) * RANDTimeType)
							<-ch
						}

					}
				}

			})
			fun.worker()
		}(m, wg, singleport)
	}
	wg.Wait()
	consumingTime := time.Since(start)
	t.Logf("[SINGLEPORT] It took ==> %s | Latency ==> %v", consumingTime, time.Duration(singleport.SumLatency)*time.Millisecond)
}

func TestMultiPort(t *testing.T) {
	multiport := NewQueueDesign(MULTIWorkerSize, MULTIPortSize, MULTIQueueSize)

	wg := &sync.WaitGroup{}
	wg.Add(multiport.PortSize)
	wg.Add(multiport.WorkerSize)

	// Build wires.
	for k := 0; k < multiport.PortSize; k++ {
		ch := make(chan []byte, multiport.OneChannelSize)
		multiport.Channels[k] = ch
	}
	for i := 0; i < multiport.PortSize; i++ {
		go func(i int, wg *sync.WaitGroup, multiport *QueueDesign) {
			fun := NewWorkLoad(func() {
				loadcount := 0
				for {
					// Fugazi data stuff.
					h := sha256.New()
					facelessvoid := time.Now()
					h.Write([]byte(facelessvoid.String()))
					timeislocked := h.Sum([]byte{1})
					latency := time.Now()
					// Produce some data to workers.
					multiport.Channels[i] <- timeislocked
					latencyDuration := time.Since(latency).Milliseconds()
					atomic.AddInt64(&multiport.SumLatency, latencyDuration)
					loadcount++
					if loadcount >= LoadLimit/multiport.PortSize {
						break
					}
				}
				close(multiport.Channels[i])
				wg.Done()
			})
			fun.worker()
		}(i, wg, multiport)
	}

	<-time.After(time.Second * 2)
	t.Log("Multiport is started...")
	start := time.Now()

	for m := 0; m < multiport.WorkerSize; m++ {
		go func(m int, wg *sync.WaitGroup, multiport *QueueDesign) {
			fun := NewWorkLoad(func() {
				countofconsume := 0
			out:
				for {
					select {
					// Consume data is from corresponding port.
					case data, ok := <-multiport.Channels[m]:
						if !ok {
							log.Printf("[MULTIPORT][TEST][WORKER:%v] Consumer is closed. Count of consumed data: %v", m, countofconsume)
							wg.Done()
							break out
						}
						countofconsume++
						// Do something.
						for i := 0; i < 1; i++ {
							_ = data
							source := rand.NewSource(time.Now().UnixNano())
							r := rand.New(source)
							randomNumber := r.Intn(RANDTimeEnd) + RANDTimeStart
							ch := time.After(time.Duration(randomNumber) * RANDTimeType)
							<-ch
						}
					}
				}

			})
			fun.worker()
		}(m, wg, multiport)
	}
	wg.Wait()
	consumingTime := time.Since(start)
	t.Logf("[MULTIPORT] It took ==> %s | Latency ==> %v", consumingTime, time.Duration(multiport.SumLatency)*time.Millisecond)
}
