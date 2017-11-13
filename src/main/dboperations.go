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

// Initialize the database. Required before any other operation.
func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "../../resources/imgtag.db")
	if err != nil {
		log.Fatal("can't open db", err)
	}
}

// Closes the database. Required once all the operations are done
func CloseDB() {
	db.Close()
}

// Reads a protobuf marshaled into "protopath" (file path) and stores its info in the
// database with the provided key (typically corresponding to the photo full path)
func StoreProtobuf(protopath string, key string) {
	data, err := ioutil.ReadFile(protopath)
	if err != nil {
		log.Fatal("can't read file", err)
	}

	response := &pb.AnnotateImageResponse{}
        err = proto.Unmarshal(data, response)
	fmt.Println(response)
        if err != nil {
                log.Fatal("can't unmarshall", err)
        }
	StoreAnnotations(key, response)
}

// Stores the annotations given in AnnotateImageResponse with the provided key (typically
// corresponding to the photo full path). Stores both the raw data (for later processing
// if desired) as protobuf, and the values of the attributes. 
func StoreAnnotations(key string, data *pb.AnnotateImageResponse) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("can't start transaction", err)
	}
	StoreRaw(key, data, tx)
	StoreResponseValues(key, data, tx)
	tx.Commit()
}

// Stores the raw data as marshaled protobuf.
func StoreRaw(filename string, data *pb.AnnotateImageResponse, tx *sql.Tx) {
	raw, err := proto.Marshal(data)
	if err != nil {
		log.Fatal("can't marshal", err)
	}

	stmt, err := db.Prepare("INSERT OR REPLACE INTO protobufs values(?, ?)");
	if err != nil {
		tx.Rollback()
                log.Fatal("can't create statement", err)
        }

	_, err = tx.Stmt(stmt).Exec(filename, raw)
	if err != nil {
		tx.Rollback()
		log.Fatal("can't insert", err)
	}
}

// Stores the response values. Supports: labels, landmarks, colors.
func StoreResponseValues(filename string, data *pb.AnnotateImageResponse, tx *sql.Tx) {
	StoreLabels(filename, data, tx)
	StoreLandmarks(filename, data, tx)
	StoreColors(filename, data, tx)
}

func StoreLabels(filename string, data *pb.AnnotateImageResponse, tx *sql.Tx) {
	for _, label := range(data.GetLabelAnnotations()) {
		storeLabel(label.Mid, label.Description, tx)
		storeImageLabel(filename, label, tx)
	}
}

// Stores the identified landmarks. There's no difference between "landmark" and "label"
// as far as the DB is concerned, we store everything as labels.
func StoreLandmarks(filename string, data *pb.AnnotateImageResponse, tx *sql.Tx) {
	for _, landmark := range(data.GetLandmarkAnnotations()) {
		storeLabel(landmark.Mid, landmark.Description, tx)
		storeImageLabel(filename, landmark, tx)
	}
}

// Transforms the RGB value of the dominant colors into a descriptive color string and
// store the aggregated color values.
func StoreColors(filename string, data *pb.AnnotateImageResponse, tx *sql.Tx) {
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

