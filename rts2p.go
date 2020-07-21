package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	config "github.com/spf13/viper"

	"github.com/snowzach/rts2p/livemedia"
)

type Stream struct {
	Url       string `json:"url" yaml:"url"`
	Name      string `json:"name", yaml:"name"`
	Username  string `json:"username" yaml:"username"`
	Password  string `json:"password" yaml:"password"`
	Verbosity int    `json:"verbosity" yaml:"verbosity"`
	Still     string `json:"still" yaml:"still"`
}

func main() {

	// Disable TLS Validation on the default http client
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Flags
	configFile := flag.String("c", "rts2p.yaml", "config file")
	flag.Parse()

	// Config
	config.SetConfigFile(*configFile)
	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("Could not read config file: %v\n", err)
	}

	// Server
	serverOptions := []livemedia.RTSPServerOption{
		livemedia.Port(config.GetInt("server.port")),
		livemedia.MaxOutPacketSize(config.GetInt("server.max_out_packet_size")),
	}
	if username := config.GetString("server.username"); username != "" {
		serverOptions = append(serverOptions, livemedia.Login(username, config.GetString("server.password")))
	}

	// Listen rtsp
	r, err := livemedia.NewRTSPServer(serverOptions...)
	if err != nil {
		log.Fatalf("error starting server: %+v\n", err)
	}
	log.Printf("Server listening on :%d\n", config.GetInt("server.port"))

	// Listen http
	var router chi.Router
	if httpPort := config.GetString("server.http_port"); httpPort != "" {
		router = chi.NewRouter()
		go func() { log.Fatalf("could not listen: %v", http.ListenAndServe(":"+httpPort, router)) }()
		log.Printf("HTTP server listening on :%s\n", httpPort)

		// Set basic auth if required
		if username := config.GetString("server.username"); username != "" {
			router.Use(middleware.BasicAuth("RTS2P", map[string]string{username: config.GetString("server.password")}))
		}
	}

	var streams []Stream
	err = config.UnmarshalKey("streams", &streams)
	if err != nil {
		log.Fatalf("error parsing stream: %+v\n", err)
	}

	for _, stream := range streams {

		streamOptions := []livemedia.RTSPStreamOption{}
		if stream.Verbosity > 0 {
			streamOptions = append(streamOptions, livemedia.Verbosity(stream.Verbosity))
		}
		if stream.Username != "" {
			streamOptions = append(streamOptions, livemedia.Credentials(stream.Username, stream.Password))
		}
		r.AddProxyStream(stream.Url, stream.Name, streamOptions...)
		log.Printf("Added stream '/%s' from %s\n", stream.Name, stream.Url)

		// Setup stream handling URLS
		if stream.Still != StillFalse && router != nil {
			if strings.HasPrefix(stream.Still, "http") {
				router.Get("/"+stream.Name, ServeStillUrl(stream.Still))
				log.Printf("Added still proxy: '/%s' from %s\n", stream.Name, stream.Still)
			} else {
				var streamURL string
				if username := config.GetString("server.username"); username != "" {
					streamURL = fmt.Sprintf("rtsp://%s:%s@127.0.0.1:%d/%s", username, config.GetString("server.password"), config.GetInt("server.port"), stream.Name)
				} else {
					streamURL = fmt.Sprintf("rtsp://127.0.0.1:%d/%s", config.GetInt("server.port"), stream.Name)
				}
				router.Get("/"+stream.Name, ServeStillStream(streamURL, stream.Still))
				log.Printf("Added still scraper: '/%s' from %s\n", stream.Name, streamURL)
			}
		}
	}

	// This will block
	r.Run()

}
