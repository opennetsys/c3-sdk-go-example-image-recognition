[![Automated Release Notes by gren](https://img.shields.io/badge/%F0%9F%A4%96-release%20notes-00B2EE.svg)](https://github-tools.github.io/github-release-notes/)

# C3 Go Image Recoginition Example

> An image recognition example in Go that runs on [C3](https://github.com/c3systems/c3-go)

## About
This library was forked from [tinrab/go-tensorflow-image-recognition](https://github.com/tinrab/go-tensorflow-image-recognition)

## Installation
If you don't want to run the docker file, you'll first need to [install tensorflow](https://www.tensorflow.org/install/install_go).

## Usage
This container accepts images and keeps track of how many of each type of image is has received. To run:

```bash
$ make build/docker
$ make run/sandbox IMAGEID=<docker_image_id>
```

The "smart container" has three methods:

### Process Image
* **Method Name:** "processImage"
* **Method Payload:** [bytesBuffer string, imageType (one of "png" or "jpg")]

### Get All Results
* **Method Name:** "getAllResults"
* **Method Payload:** nil
* **Returns:** A json object with a count of all of the image types received where the image type is the key and the count is the value.

e.g.
```json
{
  "Egyptian cat": 5,
  "Arctic fox": 30,
  "Weasel": 1,
}
```

### Get Results For Type
* **Method Name:** "getResultsForType"
* **Method Payload:** [type name string]
* **Returns:** the count of that type, e.g. "Arctic fox" would return "30"

### Demo

In terminal 1

```bash
make run/node
```

In terminal 2

Deploy

```bash
make test/node IMAGEID="QmPd91JZkGyb9nE11XETxY5NL9QUh7iR1aqwPYgaNq
so7J" PEERID="/ip4/127.0.0.1/tcp/9005/ipfs/QmSHxagveKLvxU2TMEPRwY5fRzrMU
hfjrXvqWFbPrAwdyr" METHOD="deploy"
```

Invoke Method

```bash
make test/node IMAGEID="QmPd91JZkGyb9nE11XETxY5NL9QUh7iR1aqwPYgaNq
so7J" PEERID="/ip4/127.0.0.1/tcp/9005/ipfs/QmSHxagveKLvxU2TMEPRwY5fRzrMU
hfjrXvqWFbPrAwdyr" METHOD="invokeMethod"
```

Replace `IMAGEID` and `PEERID` with your own values.

You get `IMAGEID` from `c3 deploy <docker_image_id>`

You get `PEERID` from `make run/node`

## License
[MIT](LICENSE)
