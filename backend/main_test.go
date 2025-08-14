package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

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

func TestPrint(t *testing.T){
	fmt.Println("hello")
}

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

func Test_timeout(t *testing.T){
	test := make(chan int,5)
	go func (tst chan int)  {
		for{
			select{
			case val := <-tst:
				t.Logf("data %d\n",val)
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