package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const (
	helpText = `
----------------------
\help: print help
\join <roomname>: enter the name of the room to join, or create a new one
\list: list all available rooms
\name <name>: change the user name
\quit: quit
----------------------
`
	commandHelp = `\help`
	commandJoin = `\join`
	commandList = `\list`
	commandName = `\name`
	commandQuit = `\quit`
)

// Client represents a connection for a client.
type Client struct {
	Name        string
	Receive     chan string
	conn        net.Conn
	currentRoom *Room
	rooms       *Rooms
	sendChan    chan string
}

// formatMessage is the message clients will see with a timestamp.
func formatMessage(msg string) string {
	return time.Now().Format(time.RFC3339) + "] " + msg + "\n"
}

// NewClient initializes a client for a connection.
func NewClient(conn net.Conn, rooms *Rooms) (*Client, error) {
	c := &Client{
		Receive:     make(chan string),
		conn:        conn,
		currentRoom: nil,
		rooms:       rooms,
		sendChan:    make(chan string),
	}
	name, err := c.readInput(helpText + "$ Enter name: ")
	if err != nil {
		return nil, err
	}
	c.Name = name
	err = c.WriteMessage(fmt.Sprintf("Hi %s, join a room first with command `%s`, then enter text", name, commandJoin), false)
	if err != nil {
		return nil, err
	}
	go c.prompt()
	return c, nil
}

// WriteMessage will write back a message to a client connection
// If forRoom a room formatted message is written, otherwise a server
// formatted message is written.
func (c *Client) WriteMessage(msg string, forRoom bool) error {
	if forRoom {
		_, err := c.conn.Write([]byte(formatMessage(msg)))
		return err
	}
	_, err := c.conn.Write([]byte("$ " + msg + "\n"))
	return err
}

func (c *Client) readInput(msg string) (string, error) {
	c.conn.Write([]byte(msg))
	s, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	s = strings.Trim(s, "\r\n")
	return s, nil
}

func (c *Client) prompt() {
PROMPT_LOOP:
	for {
		msg, err := c.readInput("")
		if err != nil {
			log.Println("ERRO:", c.conn.RemoteAddr().String(), err)
			continue
		}
		switch {
		case strings.HasPrefix(msg, commandJoin+" "):
			roomName := strings.SplitN(msg, " ", 2)
			if len(roomName[1]) == 0 {
				c.WriteMessage("invalid room name "+roomName[1], false)
				continue PROMPT_LOOP
			}
			room, err := c.rooms.AddOrCreateRoom(c, roomName[1])
			if err != nil {
				c.WriteMessage(err.Error(), false)
				continue PROMPT_LOOP
			}
			c.currentRoom = room
		case msg == commandQuit:
			if c.currentRoom != nil {
				c.currentRoom.Broadcast(fmt.Sprintf("[%s] %s has left..", c.currentRoom.Name, c.Name))
				c.currentRoom.RemoveClient(c)
			}
			c.Close()
			return
		case msg == commandHelp:
			c.WriteMessage(helpText, false)
		case msg == commandList:
			c.WriteMessage("Current rooms\n"+strings.Join(c.rooms.List(), "\n"), false)
		case strings.HasPrefix(msg, commandName+" "):
			name := strings.SplitN(msg, " ", 2)
			if len(name[1]) == 0 {
				c.WriteMessage("invalid user name "+name[1], false)
				continue
			}
			if c.currentRoom != nil {
				c.currentRoom.Broadcast(fmt.Sprintf("[%s] %s has left..", c.currentRoom.Name, c.Name))
				c.currentRoom.RemoveClient(c)
				c.currentRoom = nil
			}
			c.Name = name[1]
			c.WriteMessage("Rejoin a room to send messages "+name[1], false)
		default:
			if c.currentRoom == nil {
				c.WriteMessage("Join a room to send a message", false)
				continue
			}
			c.currentRoom.Broadcast(fmt.Sprintf("[%s] (%s) %s", c.currentRoom.Name, c.Name, msg))
		}
	}
}

// Close will close the client connection
func (c *Client) Close() {
	log.Println("INFO: closing client", c.conn.RemoteAddr().String(), c.Name)
	c.conn.Close()
}
