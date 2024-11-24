FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm

RUN apt-get update && \
    apt-get install -y openjdk-17-jdk graphviz gnuplot && \
    apt-get clean && rm -rf /var/lib/apt/lists/*
