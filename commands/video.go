package commands

import (
	"fmt"
	"image/color"

	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var (
	streamCmd = &cobra.Command{
		Use:   "stream",
		Short: "Stream",
		Long:  "Stream",
		Run:   streamCommand,
	}
)

func init() {
	RootCmd.AddCommand(streamCmd)
}

func streamCommand(cmd *cobra.Command, args []string) {
	// set to use a video capture device 0
	deviceID := 0

	// open webcam
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	// open display window
	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// color for the rect when faces detected
	blue := color.RGBA{0, 0, 255, 0}

	// load classifier to recognize faces
	cf, err := loadCascadeClassifier(cascadeFile)
	if err != nil {
		Exit(err, 1)
	}
	defer cf.Close()

	fmt.Printf("start reading camera device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := cf.DetectMultiScale(img)
		fmt.Printf("found %d faces\n", len(rects))

		// draw a rectangle around each face on the original image
		for _, r := range rects {
			gocv.Rectangle(&img, r, blue, 3)
		}

		// show the image in the window, and wait 1 millisecond
		window.IMShow(img)
		window.WaitKey(1)
	}
}
