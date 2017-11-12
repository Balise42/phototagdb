package main

import (
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
	"google.golang.org/genproto/googleapis/type/color"
	"fmt"
	"io/ioutil"
	"log"
	"github.com/golang/protobuf/proto"
	cm "github.com/jkl1337/go-chromath"
	"github.com/jkl1337/go-chromath/deltae"
)

var (
	targetIlluminant = &cm.IlluminantRefD50
	rgb2xyz = cm.NewRGBTransformer(&cm.SpaceSRGB, &cm.AdaptationBradford, targetIlluminant, &cm.Scaler8bClamping, 1.0, nil)
	lab2xyz = cm.NewLabTransformer(targetIlluminant)
	PredefinedColors = map[cm.Lab]string {
		RGBtoLab(cm.RGB{255,0,0}): "red",
		RGBtoLab(cm.RGB{0,255,0}): "green",
		RGBtoLab(cm.RGB{0,0,255}): "blue",
		RGBtoLab(cm.RGB{255,255,255}): "white",
		RGBtoLab(cm.RGB{0,0,0}): "black",
		RGBtoLab(cm.RGB{127,127,127}): "grey",
		RGBtoLab(cm.RGB{255,192,203}): "pink",
		RGBtoLab(cm.RGB{34,139,24}): "green",
		RGBtoLab(cm.RGB{0,191,255}): "blue",
		RGBtoLab(cm.RGB{255,255,0}): "yellow",
		RGBtoLab(cm.RGB{0,255,255}): "blue",
		RGBtoLab(cm.RGB{255,0,255}): "purple",
		RGBtoLab(cm.RGB{148,0,211}): "purple",
		RGBtoLab(cm.RGB{80,80,80}): "grey",
		RGBtoLab(cm.RGB{130,130,100}): "brown",
		RGBtoLab(cm.RGB{60,90,60}): "green",
		RGBtoLab(cm.RGB{128,0,0}): "brown",
		RGBtoLab(cm.RGB{0,128,0}): "green",
		RGBtoLab(cm.RGB{0,0,128}): "blue",
		RGBtoLab(cm.RGB{255,165,0}): "orange",
		RGBtoLab(cm.RGB{255,215,0}): "yellow",
		RGBtoLab(cm.RGB{124,252,0}): "green",
		RGBtoLab(cm.RGB{220,20,60}): "red",
		RGBtoLab(cm.RGB{255,228,196}): "beige",
		RGBtoLab(cm.RGB{245,245,220}): "beige",
		RGBtoLab(cm.RGB{255,248,220}): "beige",
		RGBtoLab(cm.RGB{139,69,19}): "brown",
		RGBtoLab(cm.RGB{210,105,30}): "brown",
	}
)

func test() {
	data, err := ioutil.ReadFile("/home/isa/tmp/protobuf")
        if err != nil {
                log.Fatal("can't read file", err)
        }

        response := &pb.AnnotateImageResponse{}
	err = proto.Unmarshal(data, response)

	if err != nil {
                log.Fatal("can't unmarshall", err)
        }
	fmt.Println(ComputeColors(response.GetImagePropertiesAnnotation().GetDominantColors()))
}

func ComputeColors(data *pb.DominantColorsAnnotation) map[string]float32 {
	res := make(map[string]float32)

	for _, colorInfo := range(data.GetColors()) {
		name := computeColorName(colorInfo.GetColor())
		val := res[name]
		res[name] = val + colorInfo.GetPixelFraction()
	}

	return res
}

func computeColorName(c *color.Color) string {
	var dist float64
	dist = 255*255*3+1
	res := "undef"
	for predefined, name := range(PredefinedColors) {
		tmpDist := distance(colorToLab(c), predefined)
		if (tmpDist < dist) {
			dist = tmpDist
			res = name
		}
	}
	fmt.Println(c, res)
	return res
}

func colorToLab(c *color.Color) cm.Lab {
	return RGBtoLab(cm.RGB{float64(c.Red), float64(c.Green), float64(c.Blue)});
}

func RGBtoLab(rgb cm.RGB) cm.Lab {
	xyz := rgb2xyz.Convert(rgb)
	return lab2xyz.Invert(xyz)
}

func distance (c1 cm.Lab, c2 cm.Lab) float64 {
	return deltae.CIE2000(c1, c2, &deltae.KLChDefault)
}
