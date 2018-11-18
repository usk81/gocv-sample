package commands

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var (
	onlineCmd = &cobra.Command{
		Use:   "online",
		Short: "Download and crop image data",
		Long:  "Download and crop image data",
		Run:   onlineCommand,
	}
	cascadeFile string
)

func init() {
	onlineCmd.PersistentFlags().StringVarP(&cascadeFile, "cascadeFile", "c", "", "custom cascade file path")
	RootCmd.AddCommand(onlineCmd)
}

func onlineCommand(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		Exit(fmt.Errorf("How to run:\n\tonline [imgfile]"), 1)
	}

	u := args[0]

	res, err := http.Get(u)
	if err != nil {
		Exit(err, 1)
	}
	defer res.Body.Close()

	var bs []byte
	bs, err = ioutil.ReadAll(res.Body)
	if err != nil {
		Exit(err, 1)
	}
	res.Body = ioutil.NopCloser(bytes.NewBuffer(bs))

	_, format, err := image.DecodeConfig(res.Body)
	if err != nil {
		// 画像フォーマットではない場合はエラーが発生する
		Exit(err, 1)
	}

	img, err := gocv.IMDecode(bs, gocv.IMReadColor)
	if err != nil {
		Exit(err, 1)
	}
	defer img.Close()
	if img.Empty() {
		fmt.Printf("Error reading image from: %s\n", u)
		return
	}

	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

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

	now := time.Now().Unix()

	// detect faces
	rects := cf.DetectMultiScaleWithParams(gray, 1.1, 12, 0, image.Point{50, 50}, image.Point{500, 500})

	// draw a rectangle around each face on the original image
	for i, r := range rects {
		result := gocv.NewMatWithSize(200, 200, gocv.MatTypeCV8U)
		gocv.Resize(img.Region(r), &result, image.Pt(result.Rows(), result.Cols()), 0, 0, gocv.InterpolationCubic)
		gocv.IMWrite(fmt.Sprintf("%s/%d.%d.%s", d, now, i, format), result)
	}
}
