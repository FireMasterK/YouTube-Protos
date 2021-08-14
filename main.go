package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/FireMasterK/youtube-protos/v2/compiled/github.com/FireMasterK/youtube-protos/youtubeprotos"
	"google.golang.org/protobuf/proto"
)

// http/2 client
var h2client = &http.Client{
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 20 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		IdleConnTimeout:       30 * time.Second,
		ReadBufferSize:        16 * 1024,
		ForceAttemptHTTP2:     true,
		MaxConnsPerHost:       0,
		MaxIdleConnsPerHost:   10,
		MaxIdleConns:          0,
	},
}

// user agent to use
var ua = "com.google.android.youtube/16.20.35(Linux; U; Android 10; en_US; Pixel 4 XL Build/QQ3A.200805.001) gzip"

func main() {

	ObjectRequest := youtubeprotos.PlayerRequest{
		VideoId: "dQw4w9WgXcQ",
		Context: &youtubeprotos.Context{
			Client: &youtubeprotos.ClientInfo{
				ClientName:    youtubeprotos.Client_ANDROID,
				ClientVersion: "16.20.35",
			},
		},
	}

	out, err := proto.Marshal(&ObjectRequest)

	if err != nil {
		log.Fatalln("Failed to encode Object: ", err)
	}

	request, err := http.NewRequest("POST", "https://youtubei.googleapis.com/youtubei/v1/player?key=AIzaSyA8eiZmM1FaDVjRy-df2KTyQ_vz_yYM39w", bytes.NewReader(out))

	request.Header.Set("User-Agent", ua)
	request.Header.Set("Content-Type", "application/x-protobuf")
	request.Header.Set("x-goog-api-format-version", "1")

	if err != nil {
		log.Fatalln("Failed to create Request: ", err)
	}

	resp, err := h2client.Do(request)

	if err != nil {
		log.Fatalln("Failed to perform Request: ", err)
	}

	respByte, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln("Failed to read Response: ", err)
	}

	respObj := &youtubeprotos.PlayerResponse{}

	proto.Unmarshal(respByte, respObj)

	for _, form := range respObj.StreamingData.AdaptiveFormats {
		fmt.Println(form.Itag) // or form.Url for example
	}

	fmt.Println(respObj.VideoDetails.Title)
	fmt.Println(respObj.VideoDetails.Thumbnail[0])
	fmt.Println(respObj.StreamingData.ExpiresInSeconds)
	fmt.Println(respObj.VideoDetails.ViewCount)
	fmt.Println(respObj.VideoDetails.IsPrivate)
	fmt.Println(respObj.VideoDetails.Keywords[0])
	fmt.Println(respObj.VideoDetails.AllowRatings)
	// fmt.Println(respObj.VideoDetails.ShortDescription)

	// ioutil.WriteFile("response.dump", respByte, 0644)

	// fmt.Println(string(respByte))

}
