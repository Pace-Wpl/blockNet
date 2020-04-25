package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	var poin string
	fmt.Scanln(&poin)
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:"+poin)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		fmt.Println("Client connect error ! " + err.Error())
		return
	}

	defer func(duration int) {
		time.Sleep(time.Duration(duration) * time.Second)
		conn.Close()
		fmt.Println("connect has closed!\n")
	}(5)

	fmt.Println(conn.LocalAddr().String() + " : Client connected, poin:" + poin)

	go MessageReceived(conn)
	onMessageSend(conn)

}
func MessageReceived(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		fmt.Println(msg)
		if err != nil || err == io.EOF || msg == "EXIT" {
			fmt.Println("MessageReceived goroutin exit!")
			break
		}
	}
}
func onMessageSend(conn *net.TCPConn) {

	// reader := bufio.NewReader(conn)
	for {
		// msg, err := reader.ReadString('\n')
		// fmt.Println(msg)

		// if err != nil || err == io.EOF {
		// 	fmt.Println(err)
		// 	break
		// }
		// time.Sleep(time.Second * 2)

		// inputBytes := make([]byte, 512)
		// _, err = os.Stdin.Read(inputBytes)
		var input string
		var err error
		fmt.Scanln(&input)
		// inputReader := bufio.NewReader(os.Stdin)
		// input, err := inputReader.ReadString('\n')
		fmt.Println("输入：" + input)
		// if err != nil || err == io.EOF {
		// 	fmt.Println(err)
		// 	break
		// }

		if input == "EXIT," {
			b := []byte(input + "\n")
			_, err = conn.Write(b)
			if err != nil {
				fmt.Println(err)
				break
			}

			fmt.Println("正在退出...")
			break
		} else {
			b := []byte(input + "\n")
			_, err = conn.Write(b)
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	}
}
