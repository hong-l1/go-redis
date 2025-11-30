package tcp

import (
	"awesomeProject3/go-redis/ineterface/tcp"
	"awesomeProject3/go-redis/pkg/wait"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	Addr string
}

func ListenAndServerWithSignal(config2 Config, handler tcp.Handler) error {
	closeCh := make(chan struct{})
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	go func() {
		sign := <-signalCh
		switch sign {
		case syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT:
			closeCh <- struct{}{}
		}
	}()
	listen, err := net.Listen("tcp", config2.Addr)
	log.Println("listen:", listen.Addr())
	if err != nil {
		log.Fatal(err)
	}
	ListenAndServer(listen, handler, closeCh)
	fmt.Println("exit")
	return nil
}
func ListenAndServer(conn net.Listener, handler tcp.Handler, ech <-chan struct{}) {
	go func() {
		<-ech
		log.Println("Shutting down server...")
		conn.Close()
		handler.Close()
	}()
	defer func() {
		conn.Close()
		handler.Close()
	}()
	ctx := context.Background()
	w := wait.Wait{}
	for {
		conn, err := conn.Accept()
		if err != nil {
			break
		}
		w.Add(1)
		go func() {
			defer w.Done()
			handler.Handle(ctx, conn)
		}()
	}
	w.Wait()
}
