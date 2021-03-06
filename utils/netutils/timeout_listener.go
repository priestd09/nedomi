package netutils

import "net"

type timeoutConnListener struct {
	net.Listener
	maxSizeOfTransfer int64
	minSizeOfTransfer int64
}

// DeadlineToTimeoutListenerConstructor returns a function that wraps
// a provided listener in a new one
// whose Accept methods returns a wrapped net.Conn whose Deadlines set
// timeouts for each Read|Write individually.
// Example: a conn.SetReadDeadline(time.Now().Add(time.Second)) will set a timeout
// of one second. With the standard conn|listener this will mean that if you start reading a response
// calling Read multiple times but it all takes more than a second it will timeout. With a connection
// from this listener if each call to Read finishes in less than a second the connection will not timeout.
// The sizeOfTransfer argument has the meaning of the size of transfer for each deadline set not for the whole connection.
func DeadlineToTimeoutListenerConstructor(maxSizeOfTransfer, minSizeOfTransfer int64) func(l net.Listener) net.Listener {
	return func(l net.Listener) net.Listener {
		return &timeoutConnListener{
			Listener:          l,
			maxSizeOfTransfer: maxSizeOfTransfer,
			minSizeOfTransfer: minSizeOfTransfer,
		}
	}
}

// Accept calls the underlying accept and wraps the connection if not nil in timeouting connection
func (t *timeoutConnListener) Accept() (net.Conn, error) {
	conn, err := t.Listener.Accept()
	if conn != nil {
		conn = newTimeoutConn(conn, t.maxSizeOfTransfer, t.minSizeOfTransfer)
	}

	return conn, err
}
