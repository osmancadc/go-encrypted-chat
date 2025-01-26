package websocket

type Handler interface {
	Run()
	readPump()
	writePump()
	handleMessage(message []byte) error
}
