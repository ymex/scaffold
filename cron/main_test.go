package cron

import (
	"testing"
	"fmt"
	"time"
)

func TestEvery(t *testing.T) {
	i := 0
	c := New()
	spec := "*/1 * * * * ?"//每秒执行一次
	c.AddFunc(spec, func() {
		i++
		fmt.Println("cron running:", i)
	})

	c.Start()
	fmt.Println(spec)
	time.Sleep(5*time.Second)//测试代码：等待定时器打印
}
