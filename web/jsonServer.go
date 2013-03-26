package web

import (
	"code.google.com/p/go.net/websocket"
	"github.com/newthinker/onemap-installer/sys"
)

func JsonServer(ws *websocket.Conn) {
	l.Messagef("jsonServer %#v", ws.Config())

	for {
		// get message struct from chan
		var msg sys.Result

		msg = sys.GetResult()
		err := websocket.JSON.Send(ws, msg)
		if err != nil {
			l.Error(err)
			break
		}
		l.Debugf("send:%#v", msg)
	}
}
