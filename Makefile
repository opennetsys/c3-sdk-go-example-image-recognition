all: build

.PHONY: build
build:
	@CGO_ENABLED=1 go build -o app .

.PHONY: build/docker
build/docker:
	@docker build .

.PHONY: model
model:
	@ mkdir -p model && \
  	wget "https://storage.googleapis.com/download.tensorflow.org/models/inception5h.zip" -O model/inception.zip && \
  unzip model/inception.zip -d model && \
  chmod -R 777 model && rm model/inception.zip

# example:
# make img2hex IMAGE=images/cat/cat.jpg
.PHONY: img2hex
img2hex:
	@go run scripts/img2hex.go $(IMAGE)

.PHONY: run
run:
	@go run main.go image_tensor.go model.go

.PHONY: send/cat
send/cat:
	@echo "[\"processImage\", \"$$($(MAKE) img2hex IMAGE=images/cat/cat.jpg)\", \"jpg\"]" | nc localhost 3333

.PHONY: send/dog
send/dog:
	@echo "[\"processImage\", \"$$($(MAKE) img2hex IMAGE=images/dog/dog.jpg)\", \"jpg\"]" | nc localhost 3333

# example:
# make run/docker IMAGEID=<docker_image_id>
.PHONY: run/docker
run/docker:
	@docker run -p 3333:3333 $(IMAGEID)

.PHONY: run/docker/last
run/docker/last:
	@$(MAKE) run/docker IMAGEID=$$(docker images -q | grep -m1 "")

# example:
# make run/sandbox IMAGEID=<docker_image_id>
.PHONY: run/sandbox
run/sandbox:
	@go run scripts/play.go $(IMAGEID) '{}' './images/cat/cat.jpg' 'jpg'

.PHONY: run/sandbox/last
run/sandbox/last:
	@$(MAKE) run/sandbox IMAGEID=$$(docker images -q | grep -m1 "")

# make build/docker
# make run/docker IMAGEID=52b74c0ae0df
# (in another new terminal)
# make send/cat
# make send/dog
