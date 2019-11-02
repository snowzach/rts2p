#include "liveMedia.hh"
#include "BasicUsageEnvironment.hh"
#include "MediaSink.hh"
#include "livemedia.go.h"

extern "C"
{
    _UsageEnvironment NewEnvironment(void)
    {
        TaskScheduler *scheduler = BasicTaskScheduler::createNew();
        return BasicUsageEnvironment::createNew(*scheduler);
    }

    void RunEventLoop(_UsageEnvironment env)
    {
        ((UsageEnvironment *)env)->taskScheduler().doEventLoop();
    }

    _RTSPServer CreateRTSPServer(_UsageEnvironment env, int port, char *username, char *password, int max_out_packet_size)
    {
        if (max_out_packet_size)
        {
            OutPacketBuffer::maxSize = max_out_packet_size;
        }

        UserAuthenticationDatabase *authDB = NULL; // new UserAuthenticationDatabase
        if (username && password) {
            authDB = new UserAuthenticationDatabase;
            authDB->addUserRecord(username, password);
        }
        
        return RTSPServer::createNew(*(UsageEnvironment *)env, port, authDB);
    }

    void RTSPServerAddServerMediaSession(_RTSPServer rtsp_server, _ServerMediaSession sms)
    {
        ((RTSPServer *)rtsp_server)->addServerMediaSession((ServerMediaSession *)sms);
    }

    int RTSPServerAddProxyStream(_UsageEnvironment env, _RTSPServer rtsp_server, char *url, char *name, char *username, char *password, int http_port, int verbosity)
    {
        ServerMediaSession *sms = ProxyServerMediaSession::createNew(
            *(UsageEnvironment *)env,
            (RTSPServer *)rtsp_server,
            url,
            name,
            username,
            password,
            http_port, //tunnelOverHTTPPortNum
            verbosity  //verbosity
        );

        if (sms == NULL)
            return 1;

        ((RTSPServer *)rtsp_server)->addServerMediaSession(sms);

        return 0;
    }
}
