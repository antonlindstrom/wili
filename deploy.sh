#!/bin/sh

if [ ! -f /tmp/build.tar.gz ]; then
  echo "Slug unavailable"
  exit 128
fi

docker ps > /dev/null

if [ $? != 0 ]; then
  echo "Unable to run docker"
  exit 128
fi

echo "Unpacking"

test ! -d /tmp/build && mkdir /tmp/build
tar xzf /tmp/build.tar.gz -C /tmp/build

if [ ! -f /tmp/build/Dockerfile ]; then
  echo "Dockerfile does not exist"
  rm -r /tmp/build
  exit 55
fi

echo "Building"

docker kill build-pipeline-old
docker tag build-pipeline build-pipeline-old
docker build -t build-pipeline /tmp/build

echo "Booting"

docker run -d build-pipeline

echo "Done"
