package module

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	client "github.com/ariary/coco-modules/pkg/ipc"
	"github.com/ariary/coco/pkg/c2c"
	ipc "github.com/james-barrow/golang-ipc"
)

type Module interface {
	//Work: the payload of your module
	Work(params []c2c.Params, output chan string)
	//Getname: return the module's name
	GetName() string
}

//Work: the payload of your module
func Work(module Module, params []c2c.Params, output chan string) {
	module.Work(params, output)
}

//ConnectToAgent: Retrive socket name, and connect to it, then send message to announce the connection
func ConnectToAgent(socket string) (cc *ipc.Client, err error) {
	//Get socket name
	flag.StringVar(&socket, "socket", "", "provide socket name for ipc communication")
	flag.Parse()

	if socket == "" {
		//try from stdin
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		fmt.Println("input:", input)
		if input != "" {
			socket = input[:len(input)-1]
		}
	}
	//start ipc client
	cc, err = ipc.StartClient(socket, nil)
	if err != nil {
		return cc, err
	}
	return cc, nil
}

//WaitConnectionValidation: Wait to receive ACK of connection from agent
func WaitConnectionValidation(cc *ipc.Client, module Module) {
	for {
		//confirm connection
		client.CheckSendMessage(cc, c2c.CONNECTION_KEYWORD+":"+module.GetName())
		msg, err := cc.Read()
		if err == nil {
			msgStr := string(msg.Data)
			if strings.HasPrefix(msgStr, c2c.LOADED_KEYWORD) {
				fmt.Println("⏳ connected to agent, wait instruction...")
				return
			}
		} else {
			fmt.Println("Error while receiving ipc message:", err)
		}
	}
}

//WaitInstruction: wait indefinitely instructions from agent
func WaitInstruction(cc *ipc.Client, module Module) {
	for {
		msg, err := cc.Read()
		if err == nil {
			instr := c2c.Instruction{}
			if err := json.Unmarshal([]byte(msg.Data), &instr); err != nil {
				fmt.Println(err, "\nclient do not received instruction but:\"", string(msg.Data), "\"")
			}
			switch instr.Type {
			case c2c.Run:
				output := make(chan string)
				// client.CheckSendMessage(cc, c2c.INSTR_OK)
				go Work(module, instr.Params, output)
				client.CheckSendMessage(cc, <-output)
			case c2c.Kill:
				client.CheckSendMessage(cc, c2c.INSTR_OK)
				os.Exit(0)
			}

		} else {
			fmt.Println("Error while receiving ipc message:", err)
		}
	}
}
