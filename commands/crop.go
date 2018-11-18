package commands

import (
	"fmt"
	"image"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var (
	cropCmd = &cobra.Command{
		Use:   "crop",
		Short: "Crop a photo file",
		Long:  "Crop a photo file",
		Run:   cropCommand,
	}
)

func init() {
	cropCmd.PersistentFlags().StringVarP(&cascadeFile, "cascadeFile", "c", "", "custom cascade file path")
	RootCmd.AddCommand(cropCmd)
}

func cropCommand(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		Exit(fmt.Errorf("How to run:\n\tshowimage [imgfile]"), 1)
	}

	fn := args[0]

	f := filepath.Base(fn)
	fs := strings.Split(f, ".")

	img := gocv.IMRead(fn, gocv.IMReadColor)
	defer img.Close()
	if img.Empty() {
		fmt.Printf("Error reading image from: %s\n", fn)
		return
	}

	cf, err := loadCascadeClassifier(cascadeFile)
	if err != nil {
		Exit(err, 1)
	}
	defer cf.Close()

	var d string
	d, err = getOutputDir()
	if err != nil {
		Exit(err, 1)
	}

	// detect faces
	rects := cf.DetectMultiScale(img)
	// fmt.Printf("found %d faces\n", len(rects))

	// draw a rectangle around each face on the original image
	for i, r := range rects {
		result := gocv.NewMatWithSize(200, 200, gocv.MatTypeCV8U)
		gocv.Resize(img.Region(r), &result, image.Pt(result.Rows(), result.Cols()), 0, 0, gocv.InterpolationCubic)

		// gocv.Rectangle(&img, r, blue, 3)
		gocv.IMWrite(fmt.Sprintf("%s/%s.%d.%s", d, fs[0], i, fs[1]), result)
	}
}
