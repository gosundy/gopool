## go pool usage
```golang
var res=uint64(0)
func main() {
	workerPool:=NewGoPool(0,0)
	start:=time.Now()

	for i:=0;i<10;i++{
		job:=Job{Func: func() {
			var sum=0
			for i:=0;i<100000000000;i++{
				sum+=1
			}
			atomic.AddUint64(&res,uint64(sum))
		}}
		workerPool.Dispatch(&job)
	}
	workerPool.WaitAll()
	fmt.Println(res)
	end:=time.Now()
	fmt.Println(end.Sub(start).Nanoseconds())
}

```
