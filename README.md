# gocv-sample

Code used within LT of golang.tokyo #19.

## Install

```
go get -u github.com/usk81/gocv-sample
```

## Setup

```bash
# your cascade file
export DEFAULTCASCADEFILE="data/haarcascade_frontalface_default.xml"
# directory to output detected file
export GOCVOUTPUTDIR="/tmp"
```

## Usage

```bash
$ cd $GOPATH/src/github.com/usk81/gocv-sample

# face detect from webcam
go run main.go stream

# face detect from a image file
go run main.go detect xxx.jpg

# crop faces from a image file
go run main.go crop xxx.jpg

# crop faces from a image file on web
go run main.go online https://xxxxx.yyy/aaaaa.png
```

## Author

[Yusuke Komatsu](https://github.com/usk81)
