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

streams:
  - url: "rtsp://wowzaec2demo.streamlock.net/vod/mp4:BigBuckBunny_115k.mov"
    name: "mytestfeed"
    username: feedusername
    password: feedpassword
    verbosity: 0
```

This would create the feed `rtsp://server:5554/mytestfeed` and require a login of `myusername` with password `mypassword`
Omit the username and password fields if you do not want to require a login.

## Docker

Example docker command: `docker run -it -p 5554:5554 -v /my/config.yaml:/opt/rts2p/rts2p.yaml snowzach/rts2p:latest`
