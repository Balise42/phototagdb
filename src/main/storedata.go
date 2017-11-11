package main

import (
	"io/ioutil"
	"log"
	"github.com/golang/protobuf/proto"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
//	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	data, err := ioutil.ReadFile("/home/isa/tmp/protobuf")
	if err != nil {
		log.Fatal("can't read file", err)
	}

	response := &pb.AnnotateImageResponse{}
        err = proto.Unmarshal(data, response)

        if err != nil {
                log.Fatal("can't unmarshall", err)
        }
	Store("plop", response)

}

func Store(filename string, data *pb.AnnotateImageResponse) {
	raw, err := proto.Marshal(data)
	if err != nil {
		log.Fatal("can't marshal", err)
	}
	StoreRaw(filename, raw)
	StoreResponseValues(filename, data)
}

func StoreRaw(filename string, bytes []byte) {
	db, err := sql.Open("sqlite3", "../../resources/imgtag.db")

	if err != nil {
		log.Fatal("can't open db", err)
	}

	stmt, err := db.Prepare("INSERT OR REPLACE INTO protobufs values(?, ?)");

	if err != nil {
                log.Fatal("can't create statement", err)
        }

	_, err = stmt.Exec(filename, bytes)

	if err != nil {
		log.Fatal("can't insert", err)
	}
}

func StoreResponseValues(filename string, data *pb.AnnotateImageResponse) {
	storeLabels(filename, data)
}

func storeLabels(filename string, data *pb.AnnotateImageResponse) {
	for _, label := range(data.GetLabelAnnotations()) {
		storeLabel(label.Mid, label.Description)
		storeImageLabel(filename, label)
	}
}

func storeLabel(mid string, description string) {
	db, err := sql.Open("sqlite3", "../../resources/imgtag.db")

        if err != nil {
                log.Fatal("can't open db", err)
        }

        stmt, err := db.Prepare("INSERT OR IGNORE INTO labels values(?, ?)");

        if err != nil {
                log.Fatal("can't create statement", err)
        }

        _, err = stmt.Exec(mid, description)

        if err != nil {
                log.Fatal("can't insert", err)
        }
}

func storeImageLabel(filename string, label *pb.EntityAnnotation) {
	db, err := sql.Open("sqlite3", "../../resources/imgtag.db")
	if err != nil {
                log.Fatal("can't open db", err)
        }

        stmt, err := db.Prepare("INSERT OR REPLACE INTO imagelabels values(?, ?, ?)");

        if err != nil {
                log.Fatal("can't create statement", err)
        }

        _, err = stmt.Exec(filename, label.Mid, label.Score)

        if err != nil {
                log.Fatal("can't insert", err)
        }
}

