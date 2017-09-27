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
	"strings"
)

const (
	versionName string = "1.0.0"
	//TODO set comm info
	commInfo string = "Android multi utils"
)

var ExcPathIsNotApk = errors.New("you input path is not apk, please check!")
var ExcPathIsFolder = errors.New("you input path is folder, please check!")
var ExcOutPathIsExist = errors.New("you out path is exist, stop program, please check!")

type filterCLI struct {
	cli.Helper
	Version    bool `cli:"version" usage:"version"`
	Verbose    bool `cli:"verbose" usage:"see Verbose of utils"`
	Channel    string `cli:"c,channel" usage:"channel name input"`
	Properties string `cli:"p,properties" usage:"channel properties file input"`
	Resource   string `cli:"r,resource" usage:"resource apk path"`
	Output     string `cli:"o,output" usage:"output apk path"`
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

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)
}

func ReadFileContentAsStringLines(filePath string) ([]string, error) {
	result := []string{}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return result, err
	}
	s := string(b)
	for _, lineStr := range strings.Split(s, "\n") {
		lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			continue
		}
		result = append(result, lineStr)
	}
	return result, nil
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
		var channelName string
		var resource string
		var outPutPath string
		var propertiesPath string
		var properties string
		var err error
		if argv.Channel != "" && argv.Resource != "" && argv.Output != "" {
			channelName = argv.Channel
			if channelName == "" {
				ctx.String("Channel code error")
				os.Exit(1)
			}
			resource = argv.Resource
			outPutPath = argv.Output
			if FileExist(outPutPath) {
				ctx.String("Out path is exist\n-> %v\nExit 1!", outPutPath)
				os.Exit(1)
			}
			isApk := false
			isApk, err = isFilePathApk(resource)
			isApk, err = isFilePathApk(outPutPath)
			if !isApk || err != nil {
				ctx.String("Error %v, sys %v\n", ExcPathIsNotApk, err)
				os.Exit(1)
			}
			properties = fmt.Sprintf("channel = %v", channelName)
			err = nil
		} else {
			ctx.String("Your params is error please check"+
				"\n\tChannel: %v"+
				"\n\tResource path: %v"+
				"\n\tOutPath: %v"+
				"\n\tProperties: %v"+
				"\nAll this must be has!", argv.Channel, argv.Resource, argv.Output, argv.Properties)
		}
		if argv.Properties != "" {
			propertiesPath = argv.Properties
			if ! FileExist(propertiesPath) {
				ctx.String("Properties file path is error\n-> %v\nExit 1", propertiesPath)
				os.Exit(1)
			}
			lines, err := ReadFileContentAsStringLines(propertiesPath)
			if err != nil {
				ctx.String("Read Properties file error\n-> Path: %v\n-> error: %v\n\nExit 1", propertiesPath, err)
				os.Exit(1)
			}
			for _, line := range lines {
				properties = fmt.Sprintf("%s\n%s", properties, line)
			}
		}

		err = insertApkChannelInfo(channelName, properties, resource, outPutPath)
		if err != nil {
			fmt.Printf("insert error %v\n", err)
			os.Exit(1)
		}
		ctx.String("insert channel info into APK success!"+
			"\n\tChannel: %v"+
			"\n\tResource path: %v"+
			"\n\tOutpath: %v", channelName, resource, outPutPath)
		os.Exit(0)

		return nil
	})
}

func insertApkChannelInfo(channel string, body string, resource string, outPath string) error {
	if fs.IsDir(resource) {
		return ExcPathIsFolder
	}
	reader, err := zip.OpenReader(resource)
	if err != nil {
		return ExcPathIsNotApk
	}
	defer reader.Close()
	if FileExist(outPath) {
		return ExcOutPathIsExist
	}
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}

	writer := zip.NewWriter(outFile)

	var insertFiles = []struct {
		Name, Body string
	}{
		{"META-INF/pl_channel_" + channel, body},
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
