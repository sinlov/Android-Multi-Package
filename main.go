package main

import (
	"github.com/mkideal/cli"
	"os"
	"fmt"
	"regexp"
	"errors"
	"github.com/gsdocker/gsos/fs"
	"github.com/klauspost/compress/zip"
	"io/ioutil"
)

const (
	versionName string = "1.0.0"
	//TODO set comm info
	commInfo string = "Android mutil utils"
)

var ExcPathIsNotApk = errors.New("you input path is not apk, please check!")
var ExcPathIsFolder = errors.New("you input path is folder, please check!")

type filterCLI struct {
	cli.Helper
	Version  bool `cli:"version" usage:"version"`
	Verbose  bool `cli:"verbose" usage:"see Verbose of utils"`
	Channel  string `cli:"c,channel" usage:"channel name input"`
	Resource string `cli:"r,resource" usage:"resource apk path"`
	Output   string `cli:"o,output" usage:"output apk path"`
}

var isPrintVerbose = false

func verbosePrint(format string, a ...interface{}) {
	if isPrintVerbose {
		fmt.Printf(format, a)
	}
}

func isFilePathApk(apkPath string) (bool, error) {
	return regexp.MatchString(`(\.apk$)`, apkPath)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Warning you input is error pleae use -h to see help")
		os.Exit(-1)
	}
	cli.Run(new(filterCLI), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*filterCLI)
		if argv.Version {
			ctx.String(commInfo + "\n\tversion: " + versionName)
			os.Exit(0)
		}
		if argv.Verbose {
			isPrintVerbose = true
		}
		if argv.Channel != "" && argv.Resource != "" && argv.Output != "" {
			channelName := argv.Channel
			if channelName == "" {
				ctx.String("Channel code error")
				os.Exit(1)
			}
			resource := argv.Resource
			outPutPath := argv.Output
			isApk, err := isFilePathApk(resource)
			isApk, err = isFilePathApk(outPutPath)
			if !isApk || err != nil {
				ctx.String("Error %v, sys %v\n", ExcPathIsNotApk, err)
				os.Exit(1)
			}
			err = insertApkChannelInfo(channelName, resource, outPutPath)
			if err != nil {
				fmt.Printf("insert error %v\n", err)
				os.Exit(1)
			}
			ctx.String("insert channel info into APK success!" +
				"\n\tChannel: %v" +
				"\n\tResource path: %v" +
				"\n\tOutpath: %v", channelName, resource, outPutPath)
			os.Exit(0)
		} else {
			ctx.String("Your params is error please check" +
				"\n\tChannel: %v" +
				"\n\tResource path: %v" +
				"\n\tOutpath: %v" +
				"\nAll this must be has!", argv.Channel, argv.Resource, argv.Output)
		}

		return nil
	})
}

func insertApkChannelInfo(channel string, resource string, outPath string) error {
	if fs.IsDir(resource) {
		return ExcPathIsFolder
	}
	reader, err := zip.OpenReader(resource)
	if err != nil {
		return ExcPathIsNotApk
	}
	defer reader.Close()
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}

	writer := zip.NewWriter(outFile)

	var insertFiles = []struct {
		Name, Body string
	}{
		{"META-INF/pl_channel_" + channel, "channel = " + channel},
	}

	for _, file := range reader.File {
		f, err := writer.Create(file.Name)
		if err != nil {
			fmt.Printf("write err %v", err)
			return err
		}
		eachFile, err := file.Open()
		if err != nil {
			fmt.Printf("write err %v", err)
			return err
		}
		eachData, err := ioutil.ReadAll(eachFile)
		if err != nil {
			fmt.Printf("write err %v", err)
			return err
		}
		_, err = f.Write(eachData)
		if err != nil {
			fmt.Printf("write err %v", err)
			return err
		}
	}

	for _, file := range insertFiles {
		f, err := writer.Create(file.Name)
		if err != nil {
			fmt.Printf("write err %v", err)
			return err
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			fmt.Printf("write err %v", err)
			return err
		}
	}

	err = writer.Close()
	if err != nil {
		fmt.Printf("close zip file err %v", err)
		return err
	}
	return nil
}
