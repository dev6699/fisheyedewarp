package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/dev6699/fisheyedewarp"
	"gocv.io/x/gocv"
)

var (
	pfov  = flag.Float64("pfov", 120.0, "Output perspective field of view in degrees. Must be between 0 and 180.")
	fov   = flag.Float64("fov", 180.0, "Input fisheye field of view in degrees. Must be between 0 and 180.")
	img   = flag.String("img", "images/fisheye.jpg", "Path to the input fisheye image file.")
	ptype = flag.String("ptype", string(fisheyedewarp.Linear), "Type of projection to apply. Options: 'Linear', 'EqualArea', 'Orthographic', 'Stereographic'")
)

func main() {
	flag.Parse()

	mat := gocv.IMRead(*img, gocv.IMReadColor)
	if mat.Empty() {
		panic("failed to read image")
	}
	defer mat.Close()

	dewarped, err := fisheyedewarp.Dewarp(mat, *fov, *pfov, fisheyedewarp.ProjectionType(*ptype))
	if err != nil {
		panic(err)
	}
	defer dewarped.Close()

	dir := filepath.Dir(*img)
	filename := filepath.Base(*img)
	ext := filepath.Ext(filename)
	baseName := filename[:len(filename)-len(ext)]
	newFilename := fmt.Sprintf("%s/%s_%s_%.0f_%.0f%s", dir, baseName, *ptype, *fov, *pfov, ext)
	gocv.IMWrite(newFilename, *dewarped)
}
