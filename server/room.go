package server

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Room is a single channel for messages to members
type Room struct {
	sync.Mutex
	Name     string
	messages chan string
	members  map[string]*Client
	logFile  io.WriteCloser
}

// Rooms is a list of open rooms
type Rooms struct {
	sync.Mutex
	logDir string
	list   map[string]*Room
}

// NewRooms initializes a list of all rooms.
func NewRooms(logDir string) *Rooms {
	return &Rooms{
		list:   map[string]*Room{},
		logDir: logDir,
	}
}

// AddOrCreateRoom will initialize a new chat room.
func (r *Rooms) AddOrCreateRoom(c *Client, name string) (*Room, error) {
	log.Printf("INFO: %s joining room %s\n", c.Name, name)
	r.Lock()
	defer r.Unlock()
	if room, ok := r.list[name]; ok {
		room.Broadcast(fmt.Sprintf("[%s] %s joined...", room.Name, c.Name))
		if err := room.AddClient(c); err != nil {
			return nil, err
		}
		return room, nil
	}
	log.Println("INFO: creating room", name)
	room := &Room{
		Name:     name,
		messages: make(chan string),
		members:  make(map[string]*Client, 0),
	}
	if r.logDir != "" {
		f, err := os.OpenFile(filepath.Join(r.logDir, room.Name+".log"), os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("ERRO:", err)
		} else {
			room.logFile = f
		}
	}
	if err := room.AddClient(c); err != nil {
		return nil, err
	}
	r.list[name] = room
	room.Broadcast(fmt.Sprintf("[%s] %s joined...", room.Name, c.Name))
	return room, nil
}

// Close will close all rooms and clients.
func (r *Rooms) Close() {
	log.Println("INFO:", "closing rooms")
	r.Lock()
	defer r.Unlock()
	if len(r.list) == 0 {
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(len(r.list))
	for _, room := range r.list {
		go func(r *Room) {
			r.Close()
			wg.Done()
		}(room)
	}
	wg.Wait()
}

// List will return a list of all current room names.
func (r *Rooms) List() []string {
	rooms := []string{}
	r.Lock()
	for _, room := range r.list {
		rooms = append(rooms, room.Name)
	}
	r.Unlock()
	return rooms
}

// AddClient will add a client to the members list of a room.
func (r *Room) AddClient(c *Client) error {
	log.Println("INFO: adding client", r.Name, c.conn.RemoteAddr().String(), c.Name)
	r.Lock()
	if _, ok := r.members[c.Name]; ok {
		return fmt.Errorf("%s already exists", c.Name)
	}
	r.members[c.Name] = c
	r.Unlock()
	return nil
}

// RemoveClient will remove a client from the members list.
func (r *Room) RemoveClient(c *Client) {
	log.Println("INFO: removing client", r.Name, c.conn.RemoteAddr().String(), c.Name)
	r.Lock()
	delete(r.members, c.Name)
	r.Unlock()
}

// Broadcast a message to members of a room
func (r *Room) Broadcast(msg string) {
	log.Println("broadcast", r.Name, msg)
	if r.logFile != nil {
		r.logFile.Write([]byte(formatMessage(msg)))
	}
	r.Lock()
	wg := sync.WaitGroup{}
	wg.Add(len(r.members))
	for _, client := range r.members {
		go func(c *Client) {
			c.WriteMessage(msg, true)
			wg.Done()
		}(client)
	}
	r.Unlock()
	wg.Wait()
}

// Close will disconnect all members in a room.
func (r *Room) Close() {
	log.Println("INFO: disconnecting clients in", r.Name)
	r.Broadcast("server closing")
	r.Lock()
	defer r.Unlock()
	wg := sync.WaitGroup{}
	wg.Add(len(r.members))
	for _, client := range r.members {
		go func(c *Client) {
			c.Close()
			wg.Done()
		}(client)
	}
	if r.logFile != nil {
		r.logFile.Close()
	}
	wg.Wait()
}
