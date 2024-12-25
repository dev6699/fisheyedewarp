# FisheyeDewarp 🐟👁️

[![GoDoc](https://pkg.go.dev/badge/github.com/dev6699/fisheyedewarp)](https://pkg.go.dev/github.com/dev6699/fisheyedewarp)
[![Go Report Card](https://goreportcard.com/badge/github.com/dev6699/fisheyedewarp)](https://goreportcard.com/report/github.com/dev6699/fisheyedewarp)
[![License](https://img.shields.io/github/license/dev6699/fisheyedewarp)](LICENSE)

FisheyeDewarp dewarps fisheye distortion in images.

## Installation

Use `go get` to install this package:

```bash
go get github.com/dev6699/fisheyedewarp
```

### Sample usage
Check [cmd/fisheyedewarp/main.go](cmd/fisheyedewarp/main.go) for more details.

- For help
```bash
./fisheyedewarp --help
```
```bash
Usage of ./fisheyedewarp:
  -fov float
    	Input fisheye field of view in degrees. Must be between 0 and 180. (default 180)
  -img string
    	Path to the input fisheye image file. (default "images/fisheye.jpg")
  -pfov float
    	Output perspective field of view in degrees. Must be between 0 and 180. (default 120)
  -ptype string
    	Type of projection to apply. Options: 'Linear', 'EqualArea', 'Orthographic', 'Stereographic' (default "Linear")

```

### Examples

| ![Image 1](/images/fisheye.jpg) | ![Image 2](/images/fisheye_Linear_180_120.jpg) |
|:-------------------------------:|:-------------------------------:|
| **Original Fisheye**                     | **Linear projection, 120 PFOV**                     |

| ![Image 3](/images/fisheye_Linear_180_130.jpg) | ![Image 4](/images/fisheye_Linear_180_140.jpg) |
|:-------------------------------:|:-------------------------------:|
| **Linear projection, 130 PFOV**                     | **Linear projection, 140 PFOV**                     |
