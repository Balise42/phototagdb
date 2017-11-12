package main

import (
	"io/ioutil"
	"log"
	"github.com/golang/protobuf/proto"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
)

var (
	db *sql.DB
)

func main() {
	var err error
	db, err = sql.Open("sqlite3", "../../resources/imgtag.db")
	if err != nil {
		log.Fatal("can't open db", err)
	}

	data, err := ioutil.ReadFile("/home/isa/tmp/protobuf")
	if err != nil {
		log.Fatal("can't read file", err)
	}

	response := &pb.AnnotateImageResponse{}
        err = proto.Unmarshal(data, response)
	fmt.Println(response)
        if err != nil {
                log.Fatal("can't unmarshall", err)
        }
	Store("/home/isa/Photos/originaux/2012/2012-01-12/IMGP5425.JPG", response)

}

func Store(filename string, data *pb.AnnotateImageResponse) {
	fmt.Println(data)
	raw, err := proto.Marshal(data)
	if err != nil {
		log.Fatal("can't marshal", err)
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("can't start transaction", err)
	}
	StoreRaw(filename, raw, tx)
	StoreResponseValues(filename, data, tx)
	tx.Commit()
}

func StoreRaw(filename string, bytes []byte, tx *sql.Tx) {
	stmt, err := db.Prepare("INSERT OR REPLACE INTO protobufs values(?, ?)");

	if err != nil {
		tx.Rollback()
                log.Fatal("can't create statement", err)
        }

	_, err = tx.Stmt(stmt).Exec(filename, bytes)

	if err != nil {
		tx.Rollback()
		log.Fatal("can't insert", err)
	}
}

func StoreResponseValues(filename string, data *pb.AnnotateImageResponse, tx *sql.Tx) {
	storeLabels(filename, data, tx)
	storeLandmarks(filename, data, tx)
	storeColors(filename, data, tx)
}

func storeLabels(filename string, data *pb.AnnotateImageResponse, tx *sql.Tx) {
	for _, label := range(data.GetLabelAnnotations()) {
		storeLabel(label.Mid, label.Description, tx)
		storeImageLabel(filename, label, tx)
	}
}

func storeLandmarks(filename string, data *pb.AnnotateImageResponse, tx *sql.Tx) {
	for _, landmark := range(data.GetLandmarkAnnotations()) {
		storeLabel(landmark.Mid, landmark.Description, tx)
		storeImageLabel(filename, landmark, tx)
	}
}

func storeColors(filename string, data *pb.AnnotateImageResponse, tx *sql.Tx) {
	colors := ComputeColors(data.GetImagePropertiesAnnotation().GetDominantColors())
	for color, amount := range(colors) {
		storeColor(filename, color, amount, tx)
	}
}

func storeColor(filename string, color string, amount float32, tx *sql.Tx) {
	stmt, err := db.Prepare("INSERT OR REPLACE INTO colors values(?, ?, ?)")
	if err != nil {
		tx.Rollback()
		log.Fatal("can't create statement", err)
	}
	_, err = tx.Stmt(stmt).Exec(filename, color, amount)
	if err != nil {
		tx.Rollback()
		log.Fatal("can't insert", err)
	}
}

func storeLabel(mid string, description string, tx *sql.Tx) {
        stmt, err := db.Prepare("INSERT OR IGNORE INTO labels values(?, ?)");

        if err != nil {
		tx.Rollback()
                log.Fatal("can't create statement", err)
        }

        _, err = tx.Stmt(stmt).Exec(mid, description)

        if err != nil {
		tx.Rollback()
                log.Fatal("can't insert", err)
        }
}

func storeImageLabel(filename string, label *pb.EntityAnnotation, tx *sql.Tx) {
        stmt, err := db.Prepare("INSERT OR REPLACE INTO imagelabels values(?, ?, ?)");

        if err != nil {
		tx.Rollback()
                log.Fatal("can't create statement", err)
        }

        _, err = tx.Stmt(stmt).Exec(filename, label.Mid, label.Score)

        if err != nil {
		tx.Rollback()
                log.Fatal("can't insert", err)
        }
}

