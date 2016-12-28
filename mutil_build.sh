#!/usr/bin/env bash

buildPath="build"
repoPath="repo"
ReNameCode="android_multi_package"
VERSION_MAJOR=0
VERSION_MINOR=0
VERSION_PATCH=1
VERSION_BUILD=0

VersionCode=$[$[VERSION_MAJOR * 100000000] + $[VERSION_MINOR * 100000] + $[VERSION_PATCH * 100] + $[VERSION_BUILD]]
VersionName="${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_PATCH}.${VERSION_BUILD}"
packageReName="${ReNameCode}_${VersionName}"

shell_running_path=$(cd `dirname $0`; pwd)

if [ -d "${buildPath}" ]; then
    rm -rf ${buildPath}
    sleep 1
fi

echo -e "============\nPrint build info start"
go version
which go
echo -e "Your settings is
\tVersion Name -> ${ReNameCode}
\tVersion code -> ${VersionCode}
\tVersion name -> ${VersionName}
\tPackage rename -> ${packageReName}
\tOut Path -> ${shell_running_path}/${buildPath}
"
echo -e "Print build info end\n============"

mkdir -p ${buildPath}
echo "start build OSX 64"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
mv main "${buildPath}/${packageReName}_osx_64"
echo "build OSX 64 finish"

echo "start build OSX 32"
CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build main.go
mv main "${buildPath}/${packageReName}_osx_86"
echo "build OSX 32 finish"

echo "start build Linux 64"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
mv main "${buildPath}/${packageReName}_linux_64"
echo "build linux 64 finish"

echo "start build Linux 32"
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build main.go
mv main "${buildPath}/${packageReName}_linux_86"
echo "build linux 32 finish"

echo "start build windows 64"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
mv main.exe "${buildPath}/${packageReName}_win_64.exe"
echo "build windows 64 finish"

echo "start build windows 32"
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build main.go
mv main.exe "${buildPath}/${packageReName}_win_86.exe"
echo "build windows 32 finish"

read -p "Do you want repo:(nothing is not need) " word
if [ -n "$word" ] ;then
    mkdir -p "${repoPath}/${ReNameCode}/${VERSION_MAJOR}/${VERSION_MINOR}/${VERSION_PATCH}/"
    cp -r "${buildPath}/" "${repoPath}/${ReNameCode}/${VERSION_MAJOR}/${VERSION_MINOR}/${VERSION_PATCH}/"
    echo -e "\tRepo path: ${repoPath}/${ReNameCode}/${VERSION_MAJOR}/${VERSION_MINOR}/${VERSION_PATCH}/"
    cat > "${repoPath}/${ReNameCode}/${VERSION_MAJOR}/${VERSION_MINOR}/${VERSION_PATCH}/Watch.bat" << EOF
@echo off
@echo. ==== start watch info ====
${packageReName}_win_86.exe -s "%~nx1"
pause
EOF
    cat > "${repoPath}/${ReNameCode}/${VERSION_MAJOR}/${VERSION_MINOR}/${VERSION_PATCH}/Watch_64.bat" << EOF
@echo off
@echo. ==== start watch info ====
${packageReName}_win_64.exe -s "%~nx1"
pause
EOF
    cat > "${repoPath}/${ReNameCode}/${VERSION_MAJOR}/${VERSION_MINOR}/${VERSION_PATCH}/MD5.bat" << EOF
@echo off
@echo. ==== start watch info ====
${packageReName}_win_86.exe -m "%~nx1"
pause
EOF
    cat > "${repoPath}/${ReNameCode}/${VERSION_MAJOR}/${VERSION_MINOR}/${VERSION_PATCH}/MD5_64.bat" << EOF
@echo off
@echo. ==== start watch info ====
${packageReName}_win_64.exe -m "%~nx1"
pause
EOF
fi

echo -e "============\nAll the build is finish! at Path\n${shell_running_path}/${buildPath}"
