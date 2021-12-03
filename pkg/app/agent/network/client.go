package network

import (
  "fmt"
  "net"
  "openeluer.org/PilotGo/PilotGo/pkg/utils/os"

  "openeluer.org/PilotGo/PilotGo/pkg/protocol"
)

type AgentMessageHandler func(*SocketClient, *protocol.Message) error

type SocketClient struct {
	conn             net.Conn
	MessageProcesser *protocol.MessageProcesser
}

func (c *SocketClient) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	c.conn = conn

	go func(c *SocketClient) {
		readBuff := []byte{}
		for {
			buff := make([]byte, 1024)
			n, err := conn.Read(buff)
			if err != nil {
				fmt.Println("read error:", err)
				return
			}
			readBuff = append(readBuff, buff[:n]...)

			//切割frame
			i, f := protocol.TlvDecode(&readBuff)

			if i != 0 {
				readBuff = readBuff[i:]

				msg := protocol.ParseMessage(*f)
				c.MessageProcesser.ProcessMessage(c, msg)
        cmd := string(msg.Body)
        fmt.Println("recv message:", cmd)

        switch cmd{
        case "firewall stop" :
          tmp:= os.Stop()
          fmt.Println(tmp)
        case "firewall restart":
          tmp:=os.Restart()
          fmt.Println(tmp)
        case "firewall reload":
          tmp:=os.Reload()
          fmt.Println(tmp)
        case "firewall config":
          tmp:= os.Config()
          fmt.Println(tmp)
        default:
          fmt.Println("消息未执行!")
        }
			}
		}
	}(c)
	return nil
}

func (c *SocketClient) Send(msg *protocol.Message) error {
	data := msg.Encode()
	sendData := protocol.TlvEncode(data)

	data_length := len(sendData)
	send_count := 0
	for {
		n, err := c.conn.Write(sendData[send_count:])
		if err != nil {
			return err
		}
		if n+send_count >= data_length {
			send_count = send_count + n
			break
		}
	}
	return nil
}

func (c *SocketClient) BindHandler(t int, f AgentMessageHandler) {
	// c.MessageProcesser.BindHandler(t, (protocol.MessageHandler)(f))
	c.MessageProcesser.BindHandler(t, func(c protocol.MessageContext, msg *protocol.Message) error {
		return f(c.(*SocketClient), msg)
	})
}

func (c *SocketClient) Close() error {
	return c.conn.Close()
}