package gopool

import (
	"context"
	"runtime"
	"sync"
)
type Job struct {
	Func func()
}
type GoPool struct {
	jobWg sync.WaitGroup
	workWg sync.WaitGroup
	maxWorker int
	maxJob int
	jobs chan *Job
	ctx context.Context
	cancelFunc context.CancelFunc
	flag bool
}
type Worker struct{
	pool *GoPool
}

func NewGoPool(workers int,queueLen int)*GoPool{
	pool:=&GoPool{}
	if workers==0{
		pool.maxWorker=runtime.NumCPU()
	}else{
		pool.maxWorker=workers
	}
	if queueLen==0{
		pool.maxJob=pool.maxWorker
	}else{
		pool.maxJob=queueLen
	}
	pool.ctx,pool.cancelFunc=context.WithCancel(context.Background())
	pool.jobs=make(chan *Job,pool.maxJob)
	return pool
}
func (pool *GoPool)run(){

	for i:=0;i<pool.maxWorker;i++{
		worker:=Worker{pool:pool}
		pool.workWg.Add(1)
		go worker.start()
	}
}
func(pool *GoPool)Dispatch(job *Job){
	pool.jobWg.Add(1)
	if pool.flag==false{
		pool.flag=true
		go pool.run()
	}
	pool.jobs<-job
}

func(worker *Worker)start(){
	defer func() {
		worker.pool.workWg.Done()
	}()
	for{
		select{
			case job :=<-worker.pool.jobs:
				job.Func()
				worker.pool.jobWg.Done()
			case <-worker.pool.ctx.Done():
				return
		}
	}

}
func (pool *GoPool)WaitAll(){
	pool.jobWg.Wait()
	pool.cancelFunc()
	pool.workWg.Wait()
	close(pool.jobs)
}


