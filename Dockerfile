FROM debian:buster-slim

ADD gofe /

ENTRYPOINT ["sh", "-c", "/gofe -backend $BACKEND"]
