package livemedia

/*
#include <stdlib.h>
#include "livemedia.go.h"
#cgo CFLAGS: -I /usr/include/UsageEnvironment -I /usr/include/groupsock -I /usr/include/liveMedia -I /usr/include/BasicUsageEnvironment
#cgo CXXFLAGS: -std=c++11 -I /usr/include/UsageEnvironment -I /usr/include/groupsock -I /usr/include/liveMedia -I /usr/include/BasicUsageEnvironment
#cgo LDFLAGS: -lliveMedia -lgroupsock -lBasicUsageEnvironment -lUsageEnvironment
*/
import "C"
import (
	"fmt"
	// "unsafe"
)

// https://github.com/andreymal/live555ProxyServerEx/blob/master/live555ProxyServerEx.cpp

type RTSPServer struct {
	env    C._UsageEnvironment
	server C._RTSPServer

	// options
	port             int
	maxOutPacketSize int
	username         *C.char
	password         *C.char
}

type RTSPServerOption func(*RTSPServer) error

func MaxOutPacketSize(maxOutPacketSize int) RTSPServerOption {
	return func(rtsps *RTSPServer) error {
		rtsps.maxOutPacketSize = maxOutPacketSize
		return nil
	}
}

func Port(port int) RTSPServerOption {
	return func(rtsps *RTSPServer) error {
		rtsps.port = port
		return nil
	}
}

func Login(username string, password string) RTSPServerOption {
	return func(rtsps *RTSPServer) error {
		rtsps.username = C.CString(username)
		rtsps.password = C.CString(password)
		return nil
	}
}

func NewRTSPServer(optionFuncs ...RTSPServerOption) (*RTSPServer, error) {

	env := C.NewEnvironment()
	rtsps := &RTSPServer{
		env:              env,
		port:             554,
		maxOutPacketSize: 0,
	}

	// Apply the options
	for _, f := range optionFuncs {
		err := f(rtsps)
		if err != nil {
			return nil, err
		}
	}

	rtsps.server = C.CreateRTSPServer(env, C.int(rtsps.port), rtsps.username, rtsps.password, C.int(rtsps.maxOutPacketSize))

	return rtsps, nil

}

func (rtsps *RTSPServer) Run() {

	C.RunEventLoop(rtsps.env)

}

type RTSPStream struct {
	url       *C.char
	name      *C.char
	username  *C.char
	password  *C.char
	httpPort  int
	verbosity int
}

type RTSPStreamOption func(*RTSPStream) error

func Credentials(username string, password string) RTSPStreamOption {
	return func(s *RTSPStream) error {
		s.username = C.CString(username)
		s.password = C.CString(password)
		return nil
	}
}

func Verbosity(level int) RTSPStreamOption {
	return func(s *RTSPStream) error {
		s.verbosity = level
		return nil
	}
}

func HTTPPort(httpPort int) RTSPStreamOption {
	return func(s *RTSPStream) error {
		s.httpPort = httpPort
		return nil
	}
}

func (rtsps *RTSPServer) AddProxyStream(url string, name string, optionFuncs ...RTSPStreamOption) error {

	s := &RTSPStream{
		url:  C.CString(url),
		name: C.CString(name),
	}

	// Apply the options
	for _, f := range optionFuncs {
		err := f(s)
		if err != nil {
			return err
		}
	}

	ret := C.RTSPServerAddProxyStream(rtsps.env, rtsps.server, s.url, s.name, s.username, s.password, C.int(s.httpPort), C.int(s.verbosity))
	if ret != 0 {
		return fmt.Errorf("could not create proxy stream")
	}

	return nil
}
