package main

import proto "github.com/golang/protobuf/proto"
import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"stresstest/kit"
	"syscall"
)

var (
	concurrentCount int
	targetAddr      string
)

type Header struct {
	Size    uint32
	Flag    uint32
	CmdCat  uint16
	CmdId   uint16
	BizId   uint64
	TransId uint64
	Result  int32
}

func usage() {
	fmt.Println("usage: stresstest")
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	flag.IntVar(&concurrentCount, "concurrent", 10, "send goroutin count")
	flag.StringVar(&targetAddr, "target", "127.0.0.1:6221", "target server listening addr")
	flag.Parse()

	if len(os.Args) == 1 {
		usage()
	}

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// todo do job in here.
	conn, err := net.Dial("tcp", targetAddr)
	checkError(err)
	defer conn.Close()

	stopCh := make(chan bool, 2*concurrentCount)
	for i := 0; i < concurrentCount; i++ {
		go sendPack(conn, genTestCase(), stopCh)
		go recvPack(conn, stopCh)
	}

	<-exitChan // wait signal to exit.
	for i := 0; i < concurrentCount; i++ {
		stopCh <- true
	}
}

func genTestCase() []byte {

	body := new(kit.GetUserInfo)
	var userid uint64 = 40008
	body.UserId = &userid
	var nickname string = "zgwang"
	body.Nickname = &nickname
	var clientaddr string = ""
	body.ClientAddr = &clientaddr
	var areaid uint32 = 0
	body.AreaId = &areaid
	buffer, err := proto.Marshal(body)
	if err != nil {
		fmt.Println("proto.Marshal failed:", err)
	}

	fmt.Println("body len:", len(buffer))

	var header Header
	header.Size = 32 + uint32(len(buffer))
	header.Flag = 0
	header.CmdCat = 0x3003
	header.CmdId = 1

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, header)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	pack_buf := buf.Bytes()
	fmt.Println("header len:", len(pack_buf))

	err = binary.Write(buf, binary.LittleEndian, buffer)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	pack_buf = buf.Bytes()
	fmt.Println("pack len:", len(pack_buf))

	return pack_buf
}

func sendPack(conn net.Conn, buf []byte, stopCh chan bool) {
	for {
		select {
		case <-stopCh:
			return
		default:
			conn.Write(buf)
		}
	}
}

func recvPack(conn net.Conn, stopCh chan bool) {
	buf := make([]byte, 512)
	for {
		select {
		case <-stopCh:
			return
		default:
			conn.Read(buf)
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
