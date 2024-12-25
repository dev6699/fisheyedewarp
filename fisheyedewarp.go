package fisheyedewarp

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"unsafe"

	"gocv.io/x/gocv"
)

// ProjectionType defines the type of projection for the dewarping operation.
type ProjectionType string

const (
	Linear        ProjectionType = "Linear"
	EqualArea     ProjectionType = "EqualArea"
	Orthographic  ProjectionType = "Orthographic"
	Stereographic ProjectionType = "Stereographic"
)

// Dewarp transforms a fisheye image into a perspective image based on the provided field of view (FOV) values and projection type.
//
//   - fov (float64): The field of view of the input fisheye image in degrees. A value of 180 represents a full hemispherical fisheye image.
//     The valid range is 0 < fov <= 180. Default value is 180 for a full hemisphere.
//   - pfov (float64): The output perspective field of view in degrees. Values must be 0 < pfov < 180.
//     The default value is 120 degrees both vertically and horizontally for a circular fisheye image and diagonally for a full frame fisheye.
//     The pFOV relative to the input fov determines the proportional area of the fisheye image that will be transformed.
func Dewarp(img gocv.Mat, fov float64, pfov float64, projectionType ProjectionType) (*gocv.Mat, error) {
	if fov <= 0 || fov > 180 {
		return nil, fmt.Errorf("invalid FOV: %f, must be in the range 0 < fov <= 180", fov)
	}
	if pfov <= 0 || pfov >= 180 {
		return nil, fmt.Errorf("invalid PFOV: %f, must be in the range 0 < pfov < 180", pfov)
	}

	rawWidth := img.Cols()
	rawHeight := img.Rows()
	rawXCenter := rawWidth / 2
	rawYCenter := rawHeight / 2

	dim := int(math.Min(float64(rawWidth), float64(rawHeight)))
	x0 := rawXCenter - dim/2
	xf := rawXCenter + dim/2
	y0 := rawYCenter - dim/2
	yf := rawYCenter + dim/2

	roi := img.Region(image.Rect(x0, y0, xf, yf))
	defer roi.Close()

	width := roi.Cols()
	height := roi.Rows()

	dimF := math.Sqrt(math.Pow(float64(width), 2) + math.Pow(float64(height), 2))
	ofoc := dimF / (2.0 * math.Tan(pfov*math.Pi/360.0))
	ofocinv := 1.0 / ofoc

	i := arange(width)
	j := arange(height)
	iGrid, jGrid := meshgrid(i, j)

	xCenter := float64(width-1) / 2.0
	yCenter := float64(height-1) / 2.0
	xs, ys := remap(iGrid, jGrid, ofocinv, dimF, fov, xCenter, yCenter, projectionType)

	xsMat, err := convertToMat(xs)
	if err != nil {
		return nil, err
	}
	defer xsMat.Close()

	ysMat, err := convertToMat(ys)
	if err != nil {
		return nil, err
	}
	defer ysMat.Close()

	remapped := gocv.NewMat()
	gocv.Remap(roi, &remapped, &xsMat, &ysMat, gocv.InterpolationLinear, gocv.BorderConstant, color.RGBA{})
	return &remapped, nil
}

func arange(n int) []float64 {
	arr := make([]float64, n)
	for k := 0; k < n; k++ {
		arr[k] = float64(k)
	}
	return arr
}

func meshgrid(x, y []float64) ([][]float64, [][]float64) {
	m := len(x)
	n := len(y)
	xGrid := make([][]float64, n)
	yGrid := make([][]float64, n)

	for i := 0; i < n; i++ {
		xGrid[i] = make([]float64, m)
		yGrid[i] = make([]float64, m)
		for j := 0; j < m; j++ {
			xGrid[i][j] = x[j]
			yGrid[i][j] = y[i]
		}
	}

	return xGrid, yGrid
}

func remap(i, j [][]float64, ofocinv, dim, fov, xcenter, ycenter float64, projectionType ProjectionType) ([][]float32, [][]float32) {
	rows := len(i)
	cols := len(i[0])

	xd := make([][]float64, rows)
	yd := make([][]float64, rows)
	rd := make([][]float64, rows)
	phiang := make([][]float64, rows)
	rr := make([][]float64, rows)

	for row := 0; row < rows; row++ {
		xd[row] = make([]float64, cols)
		yd[row] = make([]float64, cols)
		rd[row] = make([]float64, cols)
		phiang[row] = make([]float64, cols)
		rr[row] = make([]float64, cols)

		for col := 0; col < cols; col++ {
			xdVal := i[row][col] - xcenter
			ydVal := j[row][col] - ycenter
			xd[row][col] = xdVal
			yd[row][col] = ydVal
			rdVal := math.Hypot(xdVal, ydVal)
			rd[row][col] = rdVal
			phiang[row][col] = math.Atan(ofocinv * rdVal)
		}
	}

	var ifoc float64
	switch projectionType {
	case Linear:
		ifoc = dim * 180.0 / (fov * math.Pi)

	case EqualArea:
		ifoc = dim / (2.0 * math.Sin(fov*math.Pi/720.0))

	case Orthographic:
		ifoc = dim / (2.0 * math.Sin(fov*math.Pi/360.0))

	case Stereographic:
		ifoc = dim / (2.0 * math.Tan(fov*math.Pi/720.0))
	}

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			var rrVal float64
			switch projectionType {
			case Linear:
				rrVal = ifoc * phiang[row][col]

			case EqualArea:
				rrVal = ifoc * math.Sin(phiang[row][col]/2.0)

			case Orthographic:
				rrVal = ifoc * math.Sin(phiang[row][col])

			case Stereographic:
				rrVal = ifoc * math.Tan(phiang[row][col]/2.0)
			}

			rr[row][col] = rrVal
		}
	}

	xs := make([][]float32, rows)
	ys := make([][]float32, rows)

	for row := 0; row < rows; row++ {
		xs[row] = make([]float32, cols)
		ys[row] = make([]float32, cols)

		for col := 0; col < cols; col++ {
			if rd[row][col] != 0 {
				scale := rr[row][col] / rd[row][col]
				xs[row][col] = float32(scale*xd[row][col] + xcenter)
				ys[row][col] = float32(scale*yd[row][col] + ycenter)
			} else {
				xs[row][col] = 0
				ys[row][col] = 0
			}
		}
	}
	return xs, ys
}

func convertToMat(data [][]float32) (gocv.Mat, error) {
	rows := len(data)
	cols := len(data[0])
	totalElements := rows * cols

	flatData := make([]float32, 0, totalElements)

	for _, row := range data {
		flatData = append(flatData, row...)
	}

	byteSlice := unsafe.Slice((*byte)(unsafe.Pointer(&flatData[0])), totalElements*4)
	return gocv.NewMatFromBytes(rows, cols, gocv.MatTypeCV32F, byteSlice)
}
