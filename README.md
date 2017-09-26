# This is Android Multi Package utils

For android use this library

https://github.com/MDL-Sinlov/MDL_Android-Multi-Package

# OSX or Linux build

```sh
./mutil_build.sh
```

# Use

```sh
go build -o main main.go
# see help
./main -h
# test pakcage
./main -c test_Channel -r app-debug.apk -o out.apk
```

you can download apk-to test

https://github.com/sinlov/Android-Mulit-Package/raw/master/app-debug.apk

# Principle

apk will not check `META-INF` at package

you can new file at path `./META-INF`

like `pl_channel_` + `channel_name`

> Inner file format is properties!

Android API will read this file

default properties is

|key|value|
|---|-----|
|channel|channel_name|

# Warning v1.x

if use [﻿APK Signature Scheme v2](https://source.android.com/security/apksigning/v2.html)
must return old way for sign

```gradle
android {
    ﻿signingConfigs{
        ﻿v2SigningEnabled false
    }
}
```
