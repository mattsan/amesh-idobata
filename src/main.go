package main

import (
    "fmt"
    "os"
    "bytes"
    "image/png"

    "github.com/mattsan/emattsan-go/amesh"
    "github.com/mattsan/emattsan-go/idobata"

    "github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
    Message string `json:"message"`
}

func endpointUrl() string {
    return os.Getenv("IDOBATA_HOOK_ENDPOINT_URL")
}

func postAmesh() error {
    image, datetime, err := amesh.LatestImage()
    if err != nil {
        return err
    }

    buffer := new(bytes.Buffer)
    err = png.Encode(buffer, image)
    if err != nil {
        return err
    }

    comment := fmt.Sprintf("%d/%02d/%02d-%02d:%02d の雨雲の状態\nPowered by [Tokyo Amesh](http://tokyo-ame.jwa.or.jp)\n\n---\n",
      datetime.Year(), datetime.Month(), datetime.Day(), datetime.Hour(), datetime.Minute())

    endpoint := idobata.NewHook(endpointUrl())
    _, err = endpoint.Post(
      &idobata.Image{Reader: buffer, Filename: "amesh.png"},
      &idobata.Source{Value: comment},
      &idobata.Format{Value: "markdown"},
    )
    if err != nil {
        return err
    }

    return nil
}

func Handler() (Response, error) {
    err := postAmesh()

    return Response{Message: fmt.Sprintf("%v\n", err)}, nil
}

func main() {
    lambda.Start(Handler)
}
