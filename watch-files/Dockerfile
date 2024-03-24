from alpine:latest

RUN mkdir /watch-files
COPY ./tests/watch_files_linux_arm64 /watch-files/watch_files_linux_arm64

WORKDIR /watch-files

RUN mkdir /watch-files/tests

CMD ["/watch-files/watch_files_linux_arm64"]
