# Dockerfile
FROM ubuntu:xenial

RUN apt-get update && apt-get install -y \
    ruby \
    ruby-dev \
    rubygems \
    build-essential

RUN gem install --no-ri --no-rdoc fpm
RUN fpm --version
