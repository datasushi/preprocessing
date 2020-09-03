package main

import (
	"os"
	"strings"

	"io/ioutil"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
	"gopkg.in/gographics/imagick.v3/imagick"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func main() {
	loggger, _ := zap.NewDevelopment()
	defer loggger.Sync()
	logger = loggger.Sugar()
	items, _ := ioutil.ReadDir(".")

	var fileCount int64 = 0
		
	for _, item := range items {
		if !item.IsDir() {
			var extension = filepath.Ext(item.Name())
			if extension == ".pdf" {
				fileCount++
			}
		}
	}

	bar := progressbar.Default(fileCount)
	
	for _, item := range items {
		if !item.IsDir() {
			var extension = filepath.Ext(item.Name())
			if extension == ".pdf" {
//				logger.Debug("File found",
//					zap.String("Filename :", item.Name()),
//					zap.String("Extension :", extension),
//				)
				ConvertPdfToJpg(item.Name())
				bar.Add(1)
			}
		}
	}
}

func ConvertPdfToJpg(pdfName string) {
	fileName := strings.Split(pdfName, ".")[0]
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	if err := mw.SetResolution(300, 300); err != nil {
		logger.Error(err)
	}
	if err := mw.ReadImage(pdfName); err != nil {
		logger.Error(err)
	}
	if err := mw.SetCompressionQuality(95); err != nil {
		logger.Error(err)
	}
	mw.SetIteratorIndex(0)
	if err := mw.SetFormat("jpg"); err != nil {
		logger.Error(err)
	}
   if err := mw.WriteImage(fileName + ".jpg"); err != nil {
	   logger.Error(err)
   }
	_, err := imagick.ConvertImageCommand([]string{
		"convert", fileName + ".jpg", "-colorspace", "gray", "-threshold", "90%", "-negate", fileName + "_n.jpg",
	})
	if err != nil {
		logger.Error(err)
	}
	_, err = imagick.ConvertImageCommand([]string{
		"convert", fileName + "_n.jpg", "-transparent", "black", "-format", "png", fileName + "_n_t.png",
	})
	if err != nil {
		logger.Error(err)
	}
	_, err = imagick.ConvertImageCommand([]string{
		"convert", fileName + "_n_t.png", "-threshold", "90%", "-negate", fileName + "_n_t_n.png",
	})
	if err != nil {
		logger.Error(err)
	}
	_, err = imagick.ConvertImageCommand([]string{
		"convert", fileName + "_n_t_n.png", "-fill", "black", "+opaque", "white", fileName + "_n_t_n_b.png",
	})
	if err != nil {
		logger.Error(err)
	}
	_, err = imagick.ConvertImageCommand([]string{
		"convert", fileName + "_n_t_n_b.png", "-transparent", "white", fileName + "_n_t_n_b_t.png",
	})
	evaluateError(err)
	err = os.Remove(fileName + ".jpg")
	evaluateError(err)
	err = os.Remove(fileName + "_n.jpg")
	evaluateError(err)
	err = os.Remove(fileName + "_n_t.png")
	evaluateError(err)
	err = os.Remove(fileName + "_n_t_n.png")
	evaluateError(err)
	err = os.Remove(fileName + "_n_t_n_b.png")
	evaluateError(err)
	err = os.Rename(fileName + "_n_t_n_b_t.png", fileName + ".png")
	evaluateError(err)
}

func evaluateError(err error) {
	if err != nil {
		logger.Error(err)
	}	
}