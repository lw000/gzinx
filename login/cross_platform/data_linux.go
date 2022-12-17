//go:build linux
// +build linux

package cross_platform

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func Epoll(fd int) {
	var event syscall.EpollEvent
	//创建epoll实例文件描述符，不使用时需关闭以便内核销毁实例释放资源； size参数为内核fd队列大小，内核2.6.8后已升级为动态队列该参数意义不大，但值需大于0
	epfd, e := syscall.EpollCreate(1)

	if e != nil {
		log.Println("epoll_create: ", e)
		os.Exit(1)
	}
	defer syscall.Close(epfd)
	//设置事件模式
	event.Events = syscall.EPOLLIN
	event.Fd = int32(fd) //设置监听描述符
	//注册监听事件（epfd,事件动作,监听的fd,需监听的事件）
	if e = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, &event); e != nil {
		log.Println("epoll_ctl: ", e)
		os.Exit(1)
	}
	epollWait(fd, epfd, event)
}

func EpollWait(fd, epfd int, epollEvent syscall.EpollEvent) {
	var events [10]syscall.EpollEvent
	connect = &Connect{map[int]string{}}
	for {
		nevents, e := syscall.EpollWait(epfd, events[:], -1) //等待获取就绪事件
		if e != nil {
			log.Println("EpollWait: ", e)
		}
		for ev := 0; ev < nevents; ev++ {
			event := events[ev].Events
			efd := events[ev].Fd
			//处理连接
			if int(efd) == fd && event == syscall.EPOLLIN {
				handConn(fd, epfd, &epollEvent)
			} else if event == syscall.EPOLLIN { //可读
				handMsg(epfd, int(efd))
			}
			//可写
			if events[ev].Events == syscall.EPOLLOUT {
			}
		}
	}
}

func DataInit() {
	fmt.Println("use linux")
}
