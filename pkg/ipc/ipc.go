package ipc

import (
	"fmt"

	ipc "github.com/james-barrow/golang-ipc"
)

//CheckSendMessage: send message and check errors
func CheckSendMessage(cc *ipc.Client, msg string) {
	if err := cc.Write(4, []byte(msg)); err != nil {
		fmt.Println("Error while sending client message:", err)
	}
}
