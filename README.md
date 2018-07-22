# C3 Image Recoginition Example

> An image recognition example in Go that runs on [C3](https://github.com/c3systems/c3)

## About
This library was forked from [tinrab/go-tensorflow-image-recognition](https://github.com/tinrab/go-tensorflow-image-recognition)

## Installation
If you don't want to run the docker file, you'll first need to install tensorflow:
```bash
$ sudo curl -L \
   "https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-1.3.0.tar.gz" | \
   tar -C "/usr/local" -xz

$ sudo ldconfig
```

## Usage
This container accepts images and keeps track of how many of each type of image is has received. It has three methods:

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

## License
[MIT](LICENSE)
