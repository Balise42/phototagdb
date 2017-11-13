package main

import (
	"fmt"
	cm "github.com/jkl1337/go-chromath"
	"github.com/jkl1337/go-chromath/deltae"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
	"google.golang.org/genproto/googleapis/type/color"
)

var (
	targetIlluminant = &cm.IlluminantRefD50
	rgb2xyz          = cm.NewRGBTransformer(&cm.SpaceSRGB, &cm.AdaptationBradford, targetIlluminant, &cm.Scaler8bClamping, 1.0, nil)
	lab2xyz          = cm.NewLabTransformer(targetIlluminant)
	// table of predefined colors. Probably needs some expanding/refining.
	PredefinedColors = map[cm.Lab]string{
		RGBtoLab(cm.RGB{255, 0, 0}):     "red",
		RGBtoLab(cm.RGB{0, 255, 0}):     "green",
		RGBtoLab(cm.RGB{0, 0, 255}):     "blue",
		RGBtoLab(cm.RGB{255, 255, 255}): "white",
		RGBtoLab(cm.RGB{0, 0, 0}):       "black",
		RGBtoLab(cm.RGB{127, 127, 127}): "grey",
		RGBtoLab(cm.RGB{255, 192, 203}): "pink",
		RGBtoLab(cm.RGB{34, 139, 24}):   "green",
		RGBtoLab(cm.RGB{0, 191, 255}):   "blue",
		RGBtoLab(cm.RGB{255, 255, 0}):   "yellow",
		RGBtoLab(cm.RGB{0, 255, 255}):   "blue",
		RGBtoLab(cm.RGB{255, 0, 255}):   "purple",
		RGBtoLab(cm.RGB{148, 0, 211}):   "purple",
		RGBtoLab(cm.RGB{80, 80, 80}):    "grey",
		RGBtoLab(cm.RGB{130, 130, 100}): "brown",
		RGBtoLab(cm.RGB{60, 90, 60}):    "green",
		RGBtoLab(cm.RGB{128, 0, 0}):     "brown",
		RGBtoLab(cm.RGB{0, 128, 0}):     "green",
		RGBtoLab(cm.RGB{0, 0, 128}):     "blue",
		RGBtoLab(cm.RGB{255, 165, 0}):   "orange",
		RGBtoLab(cm.RGB{255, 215, 0}):   "yellow",
		RGBtoLab(cm.RGB{124, 252, 0}):   "green",
		RGBtoLab(cm.RGB{220, 20, 60}):   "red",
		RGBtoLab(cm.RGB{255, 228, 196}): "beige",
		RGBtoLab(cm.RGB{245, 245, 220}): "beige",
		RGBtoLab(cm.RGB{255, 248, 220}): "beige",
		RGBtoLab(cm.RGB{139, 69, 19}):   "brown",
		RGBtoLab(cm.RGB{210, 105, 30}):  "brown",
		RGBtoLab(cm.RGB{70, 90, 110}):   "blue",
		RGBtoLab(cm.RGB{160, 150, 130}): "beige",
		RGBtoLab(cm.RGB{70, 80, 110}):   "blue",
		RGBtoLab(cm.RGB{170, 170, 130}): "beige",
		RGBtoLab(cm.RGB{130, 120, 80}):  "khaki",
		RGBtoLab(cm.RGB{90, 80, 50}):    "brown",
		RGBtoLab(cm.RGB{150, 160, 60}):  "green",
		RGBtoLab(cm.RGB{70, 50, 30}):    "brown",
		RGBtoLab(cm.RGB{115, 120, 60}):  "green",
	}
)

// Given the DominantColors of an image, transforms them into "labeled colors"
// and aggregates everything with the same label (e.g. "red"). Returns a map of
// color labels/proportion of the image.
func ComputeColors(data *pb.DominantColorsAnnotation) map[string]float32 {
	res := make(map[string]float32)

	for _, colorInfo := range data.GetColors() {
		name := computeColorName(colorInfo.GetColor())
		val := res[name]
		res[name] = val + colorInfo.GetPixelFraction()
	}

	return res
}

func computeColorName(c *color.Color) string {
	var dist float64
	dist = 255*255*3 + 1
	res := "undef"
	for predefined, name := range PredefinedColors {
		tmpDist := distance(colorToLab(c), predefined)
		if tmpDist < dist {
			dist = tmpDist
			res = name
		}
	}
	fmt.Println(c.Red, c.Green, c.Blue, res)
	return res
}

func colorToLab(c *color.Color) cm.Lab {
	return RGBtoLab(cm.RGB{float64(c.Red), float64(c.Green), float64(c.Blue)})
}

func RGBtoLab(rgb cm.RGB) cm.Lab {
	xyz := rgb2xyz.Convert(rgb)
	return lab2xyz.Invert(xyz)
}

func distance(c1 cm.Lab, c2 cm.Lab) float64 {
	return deltae.CIE2000(c1, c2, &deltae.KLChDefault)
}
