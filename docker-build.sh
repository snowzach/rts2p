#/bin/sh

echo "Building amd64"
docker build -t snowzach/rts2p:amd64 .
docker push snowzach/rts2p:amd64

echo "Building arm32v7"
docker buildx build --platform linux/arm/v7 -t snowzach/rts2p:arm32v7 --push -f Dockerfile .

echo "Building arm64"
docker buildx build --platform linux/arm64 -t snowzach/rts2p:arm64 --push -f Dockerfile .

echo "Creating latest manifest"
docker manifest push --purge snowzach/rts2p:latest
docker manifest create snowzach/rts2p:latest snowzach/rts2p:amd64 snowzach/rts2p:arm32v7 snowzach/rts2p:arm64
docker manifest push --purge snowzach/rts2p:latest
