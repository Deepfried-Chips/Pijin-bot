FROM golang
RUN apt install g++ gcc make
WORKDIR /app
COPY . .
CMD ./pijin -token ${TOKEN} -unsplash-token ${UNSPLASH_TOKEN}
