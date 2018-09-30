package main

import (
    "fmt"
    "os"
    "bytes"
    "image"
    "image/png"
    "strconv"

    "github.com/mattsan/emattsan-go/amesh"
    "github.com/mattsan/emattsan-go/idobata"

    "github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
    Message string `json:"message"`
}

type Watch struct {
    Radius int
    Threshold int
}

const commentFormat = `
%d/%02d/%02d %02d:%02d の雨雲の状態 Powered by [Tokyo Amesh](http://tokyo-ame.jwa.or.jp)

- 降雨範囲の割合
    - 近辺: %2d %%
    - 周辺: %2d %%

---
`

func getEnvAsInt(name string, defaultValue int) int {
    env := os.Getenv(name)
    if env == "" { return defaultValue }

    value, err := strconv.ParseInt(env, 10, strconv.IntSize)
    if err != nil { return defaultValue }

    return int(value)
}

func officePoint() image.Point {
    return image.Point{
        X: getEnvAsInt("WATCHING_POINT_X", 490),
        Y: getEnvAsInt("WATCHING_POINT_Y", 230),
    }
}

func far() Watch {
    return Watch{
        Radius: getEnvAsInt("WATCHING_FAR_RADIUS", 80),
        Threshold: getEnvAsInt("WATCHING_FAR_THRESHOLD", 50),
    }
}

func near() Watch {
    return Watch{
        Radius: getEnvAsInt("WATCHING_NEAR_RADIUS", 40),
        Threshold: getEnvAsInt("WATCHING_NEAR_THRESHOLD", 250),
    }
}

func endpointUrl() string {
    return os.Getenv("IDOBATA_HOOK_ENDPOINT_URL")
}

func postAmesh() error {
    image, err := amesh.LatestImage()
    if err != nil { return err }

    buffer := new(bytes.Buffer)
    compositeImage, err := image.Composite()
    if err != nil { return err }

    err = png.Encode(buffer, compositeImage)
    if err != nil { return err }

    timestamp := image.Timestamp

    nearRainingRatio := image.RainingRatio(officePoint(), near().Radius)
    farRainingRatio := image.RainingRatio(officePoint(), far().Radius)

    if nearRainingRatio < near().Threshold && farRainingRatio < far().Threshold { return nil }

    comment :=
      fmt.Sprintf(
        commentFormat,
        timestamp.Year(),
        timestamp.Month(),
        timestamp.Day(),
        timestamp.Hour(),
        timestamp.Minute(),
        nearRainingRatio,
        farRainingRatio,
      )

    endpoint := idobata.NewHook(endpointUrl())
    _, err = endpoint.Post(
      &idobata.Image{Reader: buffer, Filename: "amesh.png"},
      &idobata.Source{Value: comment},
      &idobata.Format{Value: "markdown"},
    )
    if err != nil { return err }

    return nil
}

func Handler() (Response, error) {
    err := postAmesh()

    return Response{Message: fmt.Sprintf("%v\n", err)}, nil
}

func main() {
    lambda.Start(Handler)
}
