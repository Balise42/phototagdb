package main
import (
	"fmt"
	"log"
	"os"
	vision "cloud.google.com/go/vision/apiv1"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
	"golang.org/x/net/context"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
)

func main() {
	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	filename := "/home/isa/Photos/voyages/Californie/yosemite/p1070504.jpg"

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	defer file.Close()
	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	res, err := client.AnnotateImage(ctx, &pb.AnnotateImageRequest {
		Image: image,
		Features: []*pb.Feature{
			{Type: pb.Feature_LANDMARK_DETECTION, MaxResults: 5},
			{Type: pb.Feature_LABEL_DETECTION, MaxResults: 10},
			{Type: pb.Feature_IMAGE_PROPERTIES},
		},
	})

	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
	}


	data, err := proto.Marshal(res)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	ioutil.WriteFile("/home/isa/tmp/protobuf", data, 0644)


	fmt.Println(res)

}
