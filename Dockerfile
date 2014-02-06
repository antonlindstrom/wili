# Git Dockerfile
FROM base
MAINTAINER Anton <anton@alley.se>

ENV DEBIAN_FRONTEND noninteractive
RUN echo "deb http://archive.ubuntu.com/ubuntu quantal main universe multiverse" > /etc/apt/sources.list
RUN apt-get update
RUN apt-get install -y openssh-server git-core curl

RUN useradd git -d /home/git -m -s /usr/bin/git-shell
RUN mkdir /home/git/.ssh
RUN mkdir /home/git/git-shell-commands
RUN mkdir /var/run/sshd

ADD no-interactive-login /home/git/git-shell-commands/no-interactive-login
RUN chmod +x /home/git/git-shell-commands/no-interactive-login

ADD sshkey.pub /home/git/.ssh/authorized_keys

RUN chown -R git:git /home/git

RUN git init --bare /build.git
ADD post-receive-hook /build.git/hooks/post-receive
RUN chmod +x /build.git/hooks/post-receive

RUN chown -R git:git /build.git

EXPOSE 22
ENTRYPOINT ["/usr/sbin/sshd"]
CMD ["-D"]
