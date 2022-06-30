package module

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	client "github.com/ariary/coco-modules/pkg/ipc"
	"github.com/ariary/coco/pkg/agent"
	ipc "github.com/james-barrow/golang-ipc"
)

type Module interface {
	//Work: the payload of your module
	Work(params []agent.Params, output chan string)
	//GetName: return the module's name
	GetName() string
	//GetPrefix: return prefix for output purpose
	GetPrefix() string
}

//Work: the payload of your module
func Work(module Module, params []agent.Params, output chan string) {
	module.Work(params, output)
}

//ConnectToAgent: Retrive socket name, and connect to it, then send message to announce the connection
func ConnectToAgent(socket string) (cc *ipc.Client, err error) {
	//Get socket name
	// flag.StringVar(&socket, "socket", "", "provide socket name for ipc communication")
	// flag.Parse()

	// if socket == "" {
	// 	//try envvar
	socket = os.Getenv(agent.COCO_SOCKET_ENVVAR)
	if socket == "" {
		//try from stdin
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		fmt.Println("input:", input)
		if input != "" {
			socket = input[:len(input)-1]
		}
	}
	// }
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
		client.CheckSendMessage(cc, agent.CONNECTION_KEYWORD+":"+module.GetName())
		msg, err := cc.Read()
		if err == nil {
			msgStr := string(msg.Data)
			if strings.HasPrefix(msgStr, agent.LOADED_KEYWORD) {
				fmt.Println(module.GetPrefix(), "⏳ connected to agent, wait instruction...")
				return
			}
		} else {
			fmt.Println(module.GetPrefix(), "Error while receiving ipc message:", err)
		}
	}
}

//WaitInstruction: wait indefinitely instructions from agent
func WaitInstruction(cc *ipc.Client, module Module) {
	for {
		msg, err := cc.Read()
		if err == nil {
			instr := agent.Instruction{}
			if err := json.Unmarshal([]byte(msg.Data), &instr); err != nil {
				fmt.Println(module.GetPrefix(), err, "\nclient do not received instruction but:\"", string(msg.Data), "\"")
			}
			switch instr.Type {
			case agent.Run:
				output := make(chan string)
				// client.CheckSendMessage(cc, agent.INSTR_OK)
				go Work(module, instr.Params, output)
				client.CheckSendMessage(cc, <-output)
			case agent.Kill:
				client.CheckSendMessage(cc, agent.INSTR_OK)
				os.Exit(0)
			}

		} else {
			fmt.Println(module.GetPrefix(), "Error while receiving ipc message:", err)
		}
	}
}
