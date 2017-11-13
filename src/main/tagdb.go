package main
import (
	"fmt"
	"log"
	"os"
	vision "cloud.google.com/go/vision/apiv1"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
	"golang.org/x/net/context"
	"io/ioutil"
	"strings"
	"gopkg.in/h2non/bimg.v1"
)

func main() {
	path := os.Args[1]
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal("can't read directory", err)
	}

	toTreat := make([]string, 0, len(files))

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".JPG") {
			toTreat = append(toTreat, path + "/" + file.Name())
		}
	}

	InitDB()
	defer CloseDB()

	tagFiles(toTreat)
}

func tagFiles(files []string) {
	ctx := context.Background()
        client, err := vision.NewImageAnnotatorClient(ctx)
        if err != nil {
              log.Fatalf("Failed to create client: %v", err)
        }
	totalImg := len(files)
	numImg := 0
	for numImg < totalImg {
		reqs := make([]*pb.AnnotateImageRequest, 0, 16)
		for i:= 0; i<16 && i < totalImg ; i++ {
			reqs = append(reqs, imageToRequest(files[numImg]))
			numImg++
		}
		res, err := client.BatchAnnotateImages(ctx, &pb.BatchAnnotateImagesRequest {
			Requests: reqs,
		})
		if err != nil {
			log.Fatal("Failed to annotate images", err)
		}
		storeResults(res, files)
	}
}

func storeResults(res *pb.BatchAnnotateImagesResponse, files []string) {
	for i, ann := range(res.GetResponses()) {
		StoreAnnotations(files[i], ann)
	}
}

func imageToRequest(filename string) *pb.AnnotateImageRequest {
	req := &pb.AnnotateImageRequest {
		Image: &pb.Image { Content: getImageBytes(filename) },
		 Features: []*pb.Feature{
                        {Type: pb.Feature_LANDMARK_DETECTION, MaxResults: 5},
                        {Type: pb.Feature_LABEL_DETECTION, MaxResults: 10},
                        {Type: pb.Feature_IMAGE_PROPERTIES},
                },
	}
	return req
}

func getImageBytes(filename string) []byte {
	fmt.Println(filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Failed to read file", err)
	}
	oldImage := bimg.NewImage(data)
	oldImageSize, _ := oldImage.Size()
	var newImage []byte
	if oldImageSize.Width > 1000 || oldImageSize.Height > 1000 {
		newImage, err = oldImage.Resize(oldImageSize.Width/4, oldImageSize.Height/4)
		if err != nil {
			log.Fatal("Failed to resize", err)
		}
	} else {
		newImage = data
	}
	return newImage
}
