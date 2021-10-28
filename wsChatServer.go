package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func main() {
	http.ListenAndServe(":8443", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errMsg = make(chan error)
		var queue = make(chan int, 10)
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		// defer conn.Close()
		if err != nil {
			conn.Close()
			return
		}
		fmt.Println(r.RequestURI) //http path loc
		fmt.Println(r.Header)
		var client = conn.RemoteAddr().String()
		fmt.Println(client)
		// go wsRecv(conn, errMsg)
		// go wsSend(conn, errMsg)
		go func() {
			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					fmt.Println("ws connect close")
					errMsg <- err
					return
				}
				switch op {
				case ws.OpText:
					fmt.Println("server Recv=>", string(msg))
				}
			}
		}()
		// go func() {
		// 	for {
		// 		reader := bufio.NewReader(os.Stdin)
		// 		buf, _, err := reader.ReadLine()
		// 		err = wsutil.WriteServerMessage(conn, 0x01, buf)
		// 		if err != nil {
		// 			fmt.Println("server send err", err)
		// 			return
		// 		}
		// 		fmt.Println("server Send<=", string(buf))
		// 	}
		// }()
		go wsSend(conn, errMsg)
		go func() {
			defer func() {
				fmt.Println("error receive quit")
			}()
			for {
				select {
				case msg := <-queue:
					fmt.Println(msg)
				case <-errMsg:
					return
				}

			}
		}()
	}))
}
func wsRecv(conn net.Conn) {

}
func wsSend(conn net.Conn, errMsg chan error) {
	for {
		reader := bufio.NewReader(os.Stdin)
		buf, _, err := reader.ReadLine()
		err = wsutil.WriteServerMessage(conn, 0x01, buf)
		if err != nil {
			fmt.Println("server send err", err)
			errMsg <- err
			return
		}
		fmt.Println("server Send<=", string(buf))
	}
}
