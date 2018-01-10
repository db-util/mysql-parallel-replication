package mysql_parallel_replication

import (
	repl "github.com/siddontang/go-mysql/replication"
	"github.com/siddontang/go-mysql/mysql"
	"sync"
	"sync/atomic"
	"context"
	"errors"
	"fmt"
)

type ParallelSyncer struct {
	*repl.BinlogSyncer
	ParallelSyncerConfig
	eventChans 	[]chan *binlogEventWrapper
	stopChans 	[]chan bool

}

type ParallelSyncerConfig struct {
	repl.BinlogSyncerConfig
	WorkerCount uint16
	WorkerQueueSize uint

	WorkerAssigner func (e *repl.BinlogEvent) []int
	EventHandler   func (binlog *repl.BinlogEvent)
	ErrorHandler   func (err error)

}

type binlogEventWrapper struct {
	*repl.BinlogEvent
	depCount 	int32
	depWait 	sync.WaitGroup
}

func NewParallelSyncer(cfg ParallelSyncerConfig) *ParallelSyncer {

	syncer := &ParallelSyncer{
		repl.NewBinlogSyncer(cfg.BinlogSyncerConfig),
		cfg,
		make([]chan *binlogEventWrapper, cfg.WorkerCount),
		make([]chan bool, cfg.WorkerCount),
	}

	for i := range syncer.eventChans {
		syncer.eventChans[i] = make(chan *binlogEventWrapper, cfg.WorkerQueueSize)
		syncer.stopChans[i] = make(chan bool)
	}

	return syncer

}

func (s *ParallelSyncer) StartSync(pos mysql.Position) error {
	stream, err := s.BinlogSyncer.StartSync(pos)

	if err != nil {
		return err
	}

	go s.parallelizeStream(stream)

	// TODO

	return nil
}

func (s *ParallelSyncer) Close() {
	s.closeWithError(nil)
}

func (s *ParallelSyncer) closeWithError(err error) {

	s.BinlogSyncer.Close()

	for _, ch := range s.stopChans {

		select {
		case ch <- true:
		default:
		}
	}

	if err != nil {
		s.ErrorHandler(err)
	}

}

// run with goroutine
func (s *ParallelSyncer) parallelizeStream(stream *repl.BinlogStreamer) {

	for {

		event, err := stream.GetEvent(context.Background())
		if err != nil {
			s.closeWithError(err)
			break
		}

		workerIds := s.WorkerAssigner(event)
		eventWrapper := &binlogEventWrapper{event, int32(len(workerIds)), sync.WaitGroup{}}
		eventWrapper.depWait.Add(1)

		for _, idx := range workerIds {
			s.eventChans[idx] <- eventWrapper
		}

	}
}

// run with goroutine
func (s *ParallelSyncer) handleEvent(chId int) {
	// TODO

	eventCh := s.eventChans[chId]
	stopCh := s.stopChans[chId]
	loop: for {

		select {
		case <- stopCh :
			break loop
		case eventWrapper := <-eventCh:
			err := s.handleEvent0(eventWrapper)
			if err != nil {
				s.closeWithError(err)
				break loop
			}
		}

	}

}

func (s *ParallelSyncer) handleEvent0(wrapper *binlogEventWrapper) (e error) {

	defer func() {
		if r := recover(); r != nil {
			e = errors.New(fmt.Sprint(r))
		}
	}()

	if atomic.AddInt32(&wrapper.depCount, -1) != 0 {
		wrapper.depWait.Wait()
	} else {
		defer wrapper.depWait.Done()
		s.EventHandler(wrapper.BinlogEvent)

	}

	return nil
}