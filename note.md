## Context && channel 相关内容
### WriteGroup
`WriteGroup`的重要说明：`Writegroup`是`sync`package中重要的`struct`,用来收集等待执行的`goroutine`.
Three methods:
1. `Add()`: it is applied to set the number of goroutine .
Eg: `Add(2)` is as same as `Add(1)` * 2
   waiting for 2 goroutine to be completed

2. `Done()`: mark for having completed one goroutine
3. `Wait()`: stall current goroutine until the counter return to 0
Eg:
```go
func TestTaskControdl(t *testing.T) {
	taskNum := 5

	wg := sync.WaitGroup{}
	wg.Add(taskNum)
	for i:=0;i<taskNum;i++{
		go func (i int)  {
			fmt.Println("info",i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
```   
### Close案例
```go
func Test(t *testing.T){
 test := make(chan int, 10)

 go func(info chan int){
 for{
    select{
        case val, ok := <- test:
        if !ok{
        t.Logf("Channel Closed!")
        return
        }

        t.Logf("data %d\n", val)
    }
}
}(test)

 go func(){
    test <- 1
    time.Sleep(1 * time.Second)
    test <- 2
    close(test)
}()

 time.Sleep(5 *time.Second)
} 
```

**不共用一个传输通道，exit单独关闭**
```go
func TestA(t *testing.T){
	test := make(chan int,5)
	exit := make(chan struct{})
	go func (info chan int,exit chan struct{})  {
		for{
			select{
			case val := <-info:
				t.Logf("data %d\n",val)
			case <- exit:
				t.Logf("Task exit\n")
				return
			}
		
		}
	}(test,exit)

	go func ()  {
		test <- 1
		time.Sleep(time.Second*1)
		test <- 2
		close(exit)
	}()
	time.Sleep(time.Second*3)
}
```

**超时任务反馈**
```go
func Test_timeout(t *testing.T){
	test := make(chan int,5)
	go func (tst chan int)  {
		for{
			select{
			case val := <-tst:
				t.Logf("data %d\n",val)
            // 如果等待超过两秒设置超时
			case <- time.After(time.Second*2):
				t.Logf("Time out")
				return
			}
		}
	}(test)
	go func(){
		test <- 1
		time.Sleep(time.Second*2)
		test <- 2
	}()
	time.Sleep(time.Second*4)
}
```

### Context的场景
`Context`适用于让多级`Goroutine`实现通信的工具，并发安全
多级嵌套：父任务停止，子任务停止，控制停止顺序（比如abcdefg 可以让顺序未efgbcda）
`context.Context`api define 4 methods
```go
type Context interface{
    Deadline()(deadline time.Time,ok bool)
    Done() <- chan struct{}
    Err() error
    Value(key any)any
}
```

通过`context.WithTimeout()`设置上下文的超时时间，到达超时之后自动关闭，而通过`context.Deadline()`设置上下文的截止时间
Eg：创建一个10ms超时的上下文
```go
ctx,cancel := context.WithTimeout(context.Background(),10*time.Millisecond)
```
`context.WithValue(parent context,key any,value any)`
```go
func TestContextValue(t *testing.T){
	a := context.Background()
	b := context.WithValue(a,"k1","val1")
	c := context.WithValue(b,"key1","val1")
	d := context.WithValue(c,"key2","val2")
	e := context.WithValue(d,"key3","val3")
	f := context.WithValue(e,"key3","val4")
	fmt.Printf("%s\n",f.Value("key3"))
	// fmt.Printf("%s\n",e.Value("key3"))
}
```

**场景描述（多级嵌套）**
我们有个父任务`A`,它启动了三个子任务`B`,`C`,`D`,每个子任务还会进一步启动自己的子任务，实现以下功能：
 1. 多级嵌套控制
 2. `E - F - B - G - C - D - A`
```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 定义执行任务
func task(name string, ctx context.Context, wg *sync.WaitGroup){
	defer wg.Done()
	fmt.Printf("Task %s started:\n",name)
	for {
		select {
		case <- ctx.Done():
			fmt.Printf("Task %s finished\n",name)
			return
		default:
			time.Sleep(500*time.Millisecond)
		}
	}
}
func main() {
	ctx_A , cancel_A := context.WithCancel(context.Background())
	ctx_B , cancel_B := context.WithCancel(ctx_A)
	ctx_C , cancel_C := context.WithCancel(ctx_A)
	ctx_D , _  := context.WithCancel(ctx_A)
	ctx_E , _  := context.WithCancel(ctx_B)
	ctx_F , _  := context.WithCancel(ctx_B)
	ctx_G , _ := context.WithCancel(ctx_C)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go task("A",ctx_A,&wg)
	wg.Add(1)
	go task("B",ctx_B,&wg)
	wg.Add(1)
	go task("C",ctx_C,&wg)
	wg.Add(1)
	go task("D",ctx_D,&wg)
	wg.Add(1)
	go task("E",ctx_E,&wg)
	wg.Add(1)
	go task("F",ctx_F,&wg)
	wg.Add(1)
	go task("G",ctx_G,&wg)

	time.Sleep(time.Second*2)
	cancel_B()
	time.Sleep(time.Second*2)
	cancel_C()
	time.Sleep(time.Second*2)
	cancel_A()
	time.Sleep(time.Second*2)
	wg.Wait()
	fmt.Println("All tasks finished")
}
```



