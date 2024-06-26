# Setting Up LiveKit Server for OpenMeeting

OpenMeeting uses the LiveKit server as the media server to support video calls and video meeting services.

## About LiveKit

[LiveKit](https://github.com/livekit/livekit-server) is an open-source WebRTC SFU written in Go, built on top of the excellent [Pion](https://github.com/pion) project. For more information, visit the [LiveKit website](https://livekit.io/).

## Quick Start

To self-host LiveKit, start the server with the following Docker command:

```bash
docker run -d \
    -p 7880:7880 \
    -p 7881:7881 \
    -p 7882:7882/udp \
    -v $PWD/config/livekit.yml:/livekit.yaml \
    livekit/livekit-server \
    --config /livekit.yaml \
    --bind 0.0.0.0
```

## Viewing Logs

To check the server logs and ensure everything is running correctly, use the following command:

```bash
docker logs livekit/livekit-server
```

## Configuring the LiveKit Address in OpenIM Chat

Update the `config/live.yml` file to configure the LiveKit server address:

```yaml
liveKit:
  url: "ws://127.0.0.1:7880" # LIVEKIT_URL, LiveKit server address and port
```

By following these steps, you can set up and configure the LiveKit server for use with OpenMeeting.

## More about Deploying LiveKi

For detailed instructions on deploying LiveKit, refer to the self-hosting [deployment documentation](https://docs.livekit.io/realtime/self-hosting/deployment/).
