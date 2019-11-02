package main

import (
	"flag"
	"log"

	config "github.com/spf13/viper"

	"github.com/snowzach/rts2p/livemedia"
)

type Stream struct {
	Url       string `json:"url" yaml:"url"`
	Name      string `json:"name", yaml:"name"`
	Username  string `json:"username" yaml:"username"`
	Password  string `json:"password" yaml:"password"`
	Verbosity int    `json:"verbosity" yaml:"verbosity"`
}

func main() {

	// Config file
	configFile := flag.String("c", "rts2p.yaml", "config file")
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

	r, err := livemedia.NewRTSPServer(serverOptions...)
	if err != nil {
		log.Fatalf("error starting server: %+v\n", err)
	}
	log.Printf("Server listening on :%d\n", config.GetInt("server.port"))

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

	}

	r.Run()

}
