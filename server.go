package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"gocv.io/x/gocv"
)

const (
	StillOnce      = "once"
	StillOnceRPI   = "once_rpi"
	StillStream    = "stream"
	StillStreamRPI = "stream_rpi"
	StillFalse     = "false"
	StillDefault   = ""
)

// ServeStillUrl returns a proxied frame from a still url
func ServeStillUrl(stillURL string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		response, err := http.DefaultClient.Get(stillURL)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not fetch frame: %v", err.Error()), http.StatusBadGateway)
			return
		}
		defer response.Body.Close()

		for key, values := range response.Header {
			w.Header().Add(key, strings.Join(values, " "))
		}
		io.Copy(w, response.Body)

	}
}

// ServeStillStream scrapes a frame from RTSP stream
func ServeStillStream(rtspURL string, method string) http.HandlerFunc {

	var err error
	var capture *gocv.VideoCapture
	var m sync.RWMutex

	// The OpenCV Buffer
	img := gocv.NewMat()

	needFrame := make(chan struct{})
	frame := make(chan []byte)

	go func() {

	needFrameLoop:
		for {
			<-needFrame
			if capture == nil {
				if method == StillOnce {
					capture, err = gocv.OpenVideoCapture(fmt.Sprintf("%s", rtspURL)) // gstreamer doesn't like to be opened and closed
				} else if method == StillStreamRPI {
					capture, err = gocv.OpenVideoCapture(fmt.Sprintf("rtspsrc location=%s ! rtph264depay ! h264parse ! omxh264dec ! appsink max-buffers=1 drop=true", rtspURL))
				} else if method == StillStream || method == StillDefault {
					capture, err = gocv.OpenVideoCapture(fmt.Sprintf("rtspsrc location=%s ! decodebin ! videoconvert ! appsink max-buffers=1 drop=true", rtspURL))
					// capture, err = gocv.OpenVideoCapture(fmt.Sprintf("rtspsrc location=%s latency=10 ! rtph264depay ! h264parse ! avdec_h264 ! videoconvert ! appsink max-buffers=1 drop=true", rtspURL))
					// capture, err = gocv.OpenVideoCapture(fmt.Sprintf("rtspsrc location=%s ! rtph264depay ! h264parse ! avdec_h264 ! videoconvert ! appsink max-buffers=1 drop=true", rtspURL))
				} else {
					log.Printf("unsupported still method: %s", method)
					frame <- nil
					continue
				}

				if err != nil {
					log.Printf("could not open rtsp: %s error: %v", rtspURL, err)
					frame <- nil
					continue
				}
			}

			m.RLock()

			for {
				// Read an image
				if ok := capture.Read(&img); !ok {
					log.Printf("could not read frame: %v", err)
					frame <- nil
					continue needFrameLoop
				}
				if img.Empty() {
					continue
				}
				break
			}

			// Re-Encode to jpg
			data, err := gocv.IMEncode(".jpg", img)
			if err != nil {
				log.Printf("could not create frame: %v", err)
				frame <- nil
				continue
			}

			frame <- data

			if method == StillOnce || method == StillOnceRPI {
				capture.Close()
				capture = nil
			}

		}

	}()

	return func(w http.ResponseWriter, r *http.Request) {

		needFrame <- struct{}{}
		data := <-frame
		if data == nil {
			http.Error(w, fmt.Sprintf("could not get frame"), http.StatusBadGateway)
			return
		}
		w.Header().Add("Content-Type", "image/jpg")
		w.Write(data)

	}
}
