// (C) 2014 Mathias Dalheimer <md@gonium.net>. Code derived from the
// Gorilla WebSocket Demo, which is licensed as follows:
// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defluxio

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var H = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (H *hub) Run() {
	for {
		select {
		case c := <-H.register:
			H.connections[c] = true
		case c := <-H.unregister:
			delete(H.connections, c)
			close(c.send)
		case m := <-H.broadcast:
			for c := range H.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(H.connections, c)
				}
			}
		}
	}
}
