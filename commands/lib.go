package commands

import (
	"fmt"
	"os"

	"gocv.io/x/gocv"
)

func loadCascadeClassifier(path string) (gocv.CascadeClassifier, error) {
	if path == "" {
		df := os.Getenv("DEFAULTCASCADEFILE")
		if df == "" {
			return gocv.CascadeClassifier{}, fmt.Errorf("environment variable (DEFAULTCASCADEFILE) is empty")
		}
		path = df
	}
	cf := gocv.NewCascadeClassifier()
	if !cf.Load(path) {
		cf.Close()
		return gocv.CascadeClassifier{}, fmt.Errorf("Error reading cascade file: %s", path)
	}
	return cf, nil
}

func getOutputDir() (string, error) {
	d := os.Getenv("GOCVOUTPUTDIR")
	if d == "" {
		return "", fmt.Errorf("environment variable (GOCVOUTPUTDIR) is empty")
	}
	return d, nil
}
