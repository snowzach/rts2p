## RTS2P - RTSP Stream Proxy

This is a simple RTSP stream proxy based off of the live 555 library. http://www.live555.com/liveMedia/
I needed something simple to proxy IP Camera feeds that's easy to use and supports docker. 

This is another foray for me into learning how to use CGO and wrap C++ libraries. 

## Config
The config is very simple, by default the docker image looks for a config file `/opt/rts2p/rts2p.yaml` but the
config library supports yaml, toml or json. You can specify `-c` to the config file you want. 

Every config option is shown:
```
server:
  port: 5554
  max_out_packet_size: 2000000
  username: myusername
  password: mypassword
  http_port: 8080

streams:
  - url: "rtsp://wowzaec2demo.streamlock.net/vod/mp4:BigBuckBunny_115k.mov"
    name: "mytestfeed"
    username: feedusername
    password: feedpassword
    verbosity: 0
    still: "https://raw.githubusercontent.com/libgit2/libgit2sharp/master/square-logo.png"
```

This would create the feed `rtsp://server:5554/mytestfeed` and require a login of `myusername` with password `mypassword`
Omit the username and password fields if you do not want to require a login.

## Still image serving
If you include server.http_port it will also serve still images on an http server with the same credentials (using basic auth)
from the RTSP stream. 

You can proxy a still image url by putting an http url into the still parameter.

You can also capture from the proxied RTSP stream and serve frames from the video using a couple options
 * stream - This will start a streaming client and serve frames
 * stream_rpi - This will do the same as stream but attempt to offload decoding to the Raspberry Pi GPU
 * once - This will create a client and capture a single frame and shutdown the client
 * false - do not serve frames for this stream

The stream options can use a decent amount of CPU but they will serve frames very fast as it's continuously decoding frames
The once option shuts down the client in between but will take a while to start and serve a frame on the next invocation

If you do not specify the still option, it defaults to `stream` if you have `server.http_port` listed. 

## Docker

Example docker command: `docker run -it -p 5554:5554 -v /my/config.yaml:/opt/rts2p/rts2p.yaml snowzach/rts2p:latest`
