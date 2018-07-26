all: build

# Build application
.PHONY: build
build:
	@CGO_ENABLED=1 go build -o app .

# Build Dockerfile
.PHONY: build/docker
build/docker:
	@docker build .

# Download tensorflow models
.PHONY: model
model:
	@ mkdir -p model && \
  	wget "https://storage.googleapis.com/download.tensorflow.org/models/inception5h.zip" -O model/inception.zip && \
  unzip model/inception.zip -d model && \
  chmod -R 777 model && rm model/inception.zip

# Convert image to hex
# example:
# $ make img2hex IMAGE=images/cat/cat.jpg
.PHONY: img2hex
img2hex:
	@go run scripts/img2hex.go $(IMAGE)

# Run application
.PHONY: run
run:
	@go run main.go image_tensor.go model.go

# Send a cat photo the application
.PHONY: send/cat
send/cat:
	@echo "[\"processImage\", \"$$($(MAKE) img2hex IMAGE=images/cat/cat.jpg)\", \"jpg\"]" | nc localhost 3333

# Send a dog photo the application
.PHONY: send/dog
send/dog:
	@echo "[\"processImage\", \"$$($(MAKE) img2hex IMAGE=images/dog/dog.jpg)\", \"jpg\"]" | nc localhost 3333

# Run the docker container given the image ID
# example:
# $ make run/docker IMAGEID=<docker_image_id>
.PHONY: run/docker
run/docker:
	@docker run -p 3333:3333 $(IMAGEID)

# Run the docker container using last image built
.PHONY: run/docker/last
run/docker/last:
	@$(MAKE) run/docker IMAGEID=$$(docker images -q | grep -m1 "")

# Run sandbox test given image ID
# example:
# $ make run/sandbox IMAGEID=<docker_image_id>
.PHONY: run/sandbox
run/sandbox:
	@go run scripts/play.go $(IMAGEID) '{}' './images/cat/cat.jpg' 'jpg'

# Run sandbox test using last image built
.PHONY: run/sandbox/last
run/sandbox/last:
	@$(MAKE) run/sandbox IMAGEID=$$(docker images -q | grep -m1 "")

# Run demo
# example:
# $ make demo IMAGEID=<image_id> PEERID=<peer_id> method="deploy"
# note: method can be "deploy" or "invokeMethod"
.PHONY: demo
demo:
	@IMAGEID="$(IMAGEID)" PEERID="$(PEERID)" METHOD="$(METHOD)" go run demo/main.go $(ARGS)

.PHONY: run/node
run/node:
	@c3-go node start --pem=demo/priv1.pem --uri /ip4/0.0.0.0/tcp/9005 --data-dir ~/.c3-1 --difficulty 5
