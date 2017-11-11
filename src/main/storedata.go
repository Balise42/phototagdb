package main

import (
	"io/ioutil"
	"log"
	"github.com/golang/protobuf/proto"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
	"fmt"
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
	fmt.Println(response)
	StoreRaw("plop", response)
}

func Store(filename string, data *pb.AnnotateImageResponse) {
	bytes, err := proto.Marshal(data)
	if err != nil {
		log.Fatal("can't marshal data", err)
	}

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