# Wili

Because Beno was too general.

### Description

This is an example project of how to make deployments of docker containers a
bit easier.

So, basically what this does is that you build a git receiver via SSH in a
docker container and you'll push this to it and it builds the project in
another container.

WARNING: This is pretty dangerous if you use other people run Dockerfiles on
your server. One mount away from doing nasty stuff. Be safe.

This lacks a way to name your projects and to run more than one simple
project.

### Install

    cat .ssh/your_key.pub > sshkey.pub
    chgrp [YOUR_GROUP] /var/run/docker
    docker build -t wili .

    docker run -d wili

Build the builder:

    go build wili.go

Run the builder:

    ./wili

That should hopefully work.
