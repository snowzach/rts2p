#ifndef GO_LIVEMEDIA_H
#define GO_LIVEMEDIA_H

#ifdef __cplusplus
extern "C"
{
#endif

    typedef void *_UsageEnvironment;
    typedef void *_RTSPServer;
    typedef void *_ServerMediaSession;

    _UsageEnvironment NewEnvironment(void);
    void RunEventLoop(_UsageEnvironment env);
    _RTSPServer CreateRTSPServer(_UsageEnvironment env, int port, char *username, char *password, int max_out_packet_size);
    int RTSPServerAddProxyStream(_UsageEnvironment env, _RTSPServer rtsp_server, char *url, char *name, char *username, char *password, int http_port, int verbosity);

#ifdef __cplusplus
}
#endif

#endif
