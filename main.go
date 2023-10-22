package main

type QueueDesign struct {
	PortSize       int
	WorkerSize     int
	OneChannelSize int
	Channels       []chan []byte
	SumLatency     int64
}

type WorkLoad struct {
	worker func()
}

func NewQueueDesign(workersize, portsize, onechannelsize int) *QueueDesign {
	return &QueueDesign{
		PortSize:       portsize,
		WorkerSize:     workersize,
		OneChannelSize: onechannelsize,
		Channels:       make([]chan []byte, workersize),
		SumLatency:     0,
	}
}

func NewWorkLoad(f func()) *WorkLoad {
	return &WorkLoad{
		worker: f,
	}
}

func main() {

}
