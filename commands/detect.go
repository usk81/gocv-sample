package commands

import (
	"fmt"
	"image/color"

	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var (
	detectCmd = &cobra.Command{
		Use:   "detect",
		Short: "detects faces from a photo file",
		Long:  "detects faces from a photo file",
		Run:   detectCommand,
	}
)

func init() {
	detectCmd.PersistentFlags().StringVarP(&cascadeFile, "cascadeFile", "c", "", "custom cascade file path")
	RootCmd.AddCommand(detectCmd)
}

func detectCommand(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		Exit(fmt.Errorf("How to run:\n\tshowimage [imgfile]"), 1)
	}

	fn := args[0]

	img := gocv.IMRead(fn, gocv.IMReadColor)
	defer img.Close()
	if img.Empty() {
		fmt.Printf("Error reading image from: %s\n", fn)
		return
	}

	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	cf, err := loadCascadeClassifier(cascadeFile)
	if err != nil {
		Exit(err, 1)
	}
	defer cf.Close()

	blue := color.RGBA{0, 0, 255, 0}

	// detect faces
	rects := cf.DetectMultiScale(img)
	fmt.Printf("found %d faces\n", len(rects))

	// draw a rectangle around each face on the original image
	for _, r := range rects {
		gocv.Rectangle(&img, r, blue, 3)
		window.IMShow(img)
	}
	window.WaitKey(0)
}
