# OS
FROM ubuntu:16.04

MAINTAINER A. Segura <alberto.segura.delgado@gmail.com>

RUN apt-get update

# Install Git and wget
RUN apt-get install -y git git-core wget --force-yes

# Install compilers, etc
RUN apt-get install -y build-essential --force-yes

# Install Redis
RUN wget http://download.redis.io/redis-stable.tar.gz
RUN tar xvzf redis-stable.tar.gz
WORKDIR "/redis-stable"
RUN make
RUN make install
RUN cd ..

# You can expose port 6379 from the container to the host
# EXPOSE 6379


# Install Golang
RUN wget https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.8.3.linux-amd64.tar.gz
RUN mkdir $HOME/gocode
RUN mkdir $HOME/gocode/src
RUN echo 'export PATH=$PATH:/usr/local/go/bin' >> /root/.bashrc
RUN echo 'export GOPATH=$HOME/gocode' >> /root/.bashrc
RUN echo 'export PATH=$PATH:$HOME/gocode/bin' >> /root/.bashrc
ENV PATH="${PATH}:/usr/local/go/bin"
ENV GOPATH="/root/gocode"
ENV PATH="${PATH}:/root/gocode/bin"

RUN mkdir /root/gocode/src/GoHole
WORKDIR "/root/gocode/src/GoHole"

# Copy GoHole code
ADD . .
# Compile
RUN sh install.sh
RUN make install
# Prepare config file
RUN cp config_example.json /root/gohole_config.json

EXPOSE 53 53/udp
EXPOSE 443 443/udp

RUN chmod +x docker/init.sh
ENTRYPOINT docker/init.sh

