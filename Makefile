all: build

.PHONY: build
build:
	@vgo build main.go image_tensor.go model.go

.PHONY: run
run:
	# example:
	# make run IMAGEID=<docker_image_id>
	@go run scripts/play.go $(IMAGEID) '{}' './images/cat/cat.jpg' 'jpg'
