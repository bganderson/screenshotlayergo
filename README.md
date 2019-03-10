# ScreenshotLayer Go
ScreenshotLayer Go is a [Go](https://golang.org) package for accessing the [screenshotlayer.com API](https://www.screenshotlayer.com/). This package 
implements all available query parameters using the REST API.

## Notice
Please note that is package is not complete and is considered pre-release, meaning it
is subject to change. Safety not guaranteed.

## Basic Usage

Assuming you have Go already setup and working, grab the latest version of the package 
from master.

```sh
go get github.com/bganderson/screenshotlayergo
```

Import the package into your project.

```go
import "github.com/bganderson/screenshotlayergo"
```

Create a screenshotlayer API client using your API Access Key. Optionally, you can specify if you want to use HTTPS or a different API endpoint.

```go
screenshotlayer := screenshotlayergo.Client{
    AccessKey:  "<YOUR_API_ACCESS_KEY>",
    HTTPS:      true,
}
```

Make an API request and save the image to disk. The only required parameter is `URL`. For all available parameters, refer to the [API Documention](https://screenshotlayer.com/documentation).

```go
r, err := screenshotlayer.Screenshot(&screenshotlayergo.APIRequest{
    URL: "https://www.bganderson.com",
})
if err != nil {
    log.Fatalln(err)
}

if err := ioutil.WriteFile("screenshot.png", r.Bytes, 0644); err != nil {
    log.Fatalln(err)
}
```
