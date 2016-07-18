package main

import (
	"flag"
	"fmt"
	"github.com/DAddYE/vips"
	"io/ioutil"
	"os"
	"strings"
)

func resizeImage(srcImage string, dstImage string) {
	options := vips.Options{
		Width:        335,
		Height:       250,
		Crop:         true,
		Extend:       vips.EXTEND_WHITE,
		Interpolator: vips.BILINEAR,
		Gravity:      vips.CENTRE,
		Quality:      80,
	}
	origFile, err := os.Open(srcImage)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer origFile.Close()
	buf, err := ioutil.ReadAll(origFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf, err = vips.Resize(buf, options)
	if err != nil {
		fmt.Println(err)
		return
	}
	cacheFile, err := os.Create(dstImage)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cacheFile.Close()
	_, err = cacheFile.Write(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func syncDir(dstPath string, dstMode os.FileMode) {
	_, err := os.Stat(dstPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(dstPath, dstMode)
			if err != nil {
				fmt.Println("can't create target directory: " + dstPath)
				return
			}
		}
	}
}

func iterateDir(srcPath string, dstPath string) {
	fmt.Println("srcPath: " + srcPath)
	fmt.Println("dstPath: " + dstPath)
	entries, err := ioutil.ReadDir(srcPath)
	if err != nil {
		fmt.Println("can't read entries from: " + srcPath)
		return
	}

	for _, r := range entries {
		fmt.Println("Entry: " + r.Name())
		if r.IsDir() {
			childSrcPath := srcPath + "/" + r.Name()
			childDstPath := dstPath + "/" + r.Name()
			syncDir(childDstPath, r.Mode())
			iterateDir(childSrcPath, childDstPath)
		} else {
			fname := strings.ToLower(r.Name())
			if strings.HasSuffix(fname, ".jpg") || strings.HasSuffix(fname, ".jpeg") {
				childSrcPath := srcPath + "/" + r.Name()
				childDstPath := dstPath + "/" + r.Name()
				childSrcInfo, _ := os.Stat(childSrcPath)
				childDstInfo, err := os.Stat(childDstPath)
				if err != nil || childSrcInfo.ModTime().Unix() > childDstInfo.ModTime().Unix() {
					resizeImage(srcPath+"/"+r.Name(), dstPath+"/"+r.Name())
				}
			}
		}
	}
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		fmt.Println("Usage: " + os.Args[0] + " SRCDIR DSTDIR")
		return
	}

	srcDirInfo, err := os.Stat(flag.Args()[0])
	if err != nil {
		fmt.Println("can't stat source directory")
		return
	}

	if !srcDirInfo.IsDir() {
		fmt.Println("source is not a directory")
		return
	}

	dstDirInfo, err := os.Stat(flag.Args()[1])
	if err != nil {
		fmt.Println("can't stat target directory")
		return
	}

	if !dstDirInfo.IsDir() {
		fmt.Println("target directory is not a file")
		return
	}

	iterateDir(flag.Args()[0], flag.Args()[1])
}
