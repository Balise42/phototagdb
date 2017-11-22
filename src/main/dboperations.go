package main

import (
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/proto"
	_ "github.com/mattn/go-sqlite3"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
	"io/ioutil"
	"log"
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

	stmt, err := db.Prepare("INSERT OR REPLACE INTO protobufs values(?, ?)")
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
	StoreTexts(filename, data, tx)
}

func StoreLabels(filename string, data *pb.AnnotateImageResponse, tx *sql.Tx) {
	for _, label := range data.GetLabelAnnotations() {
		storeLabel(label.Mid, label.Description, tx)
		storeImageLabel(filename, label, tx)
	}
}

// Stores the identified landmarks. There's no difference between "landmark" and "label"
// as far as the DB is concerned, we store everything as labels.
func StoreLandmarks(filename string, data *pb.AnnotateImageResponse, tx *sql.Tx) {
	for _, landmark := range data.GetLandmarkAnnotations() {
		storeLabel(landmark.Mid, landmark.Description, tx)
		storeImageLabel(filename, landmark, tx)
	}
}

// Transforms the RGB value of the dominant colors into a descriptive color string and
// store the aggregated color values.
func StoreColors(filename string, data *pb.AnnotateImageResponse, tx *sql.Tx) {
	colors := ComputeColors(data.GetImagePropertiesAnnotation().GetDominantColors())
	for color, amount := range colors {
		storeColor(filename, color, amount, tx)
	}
}

// Queries for a list of labels. Return the list of keys for which ALL the labels are
// stored. Labels can be partial using "%" as a placeholder character.
func QueryLabels(labels []string) []string {
	query, args := buildQueryLabels(labels)

	res, err := db.Query(query, args...)
	defer res.Close()

	if err != nil {
		log.Fatal(err)
	}

	return getFilenamesFromRes(res)
}

func getFilenamesFromRes(res *sql.Rows) []string {
	filenames := make([]string, 0, 20)

	for res.Next() {
		var filename string
		err := res.Scan(&filename)
		if err != nil {
			log.Fatal(err)
		}
		filenames = append(filenames, filename)
	}

	return filenames
}

func buildQueryLabels(labels []string) (string, []interface{}) {
	query := "SELECT DISTINCT filename FROM imagelabels, labels where imagelabels.mid = labels.mid and labels.description LIKE ?"

	for range labels[1:] {
		query = query + " INTERSECT (SELECT filename FROM imagelabels, labels where imagelabels.mid = labels.mid and labels.description LIKE ?)"
	}

	args := make([]interface{}, 0, len(labels))
	for _, l := range labels {
		args = append(args, l)
	}
	return query, args
}

// Queries for text in pictures. Returns the list of keys of images in which the
// corresponding text (or a superstring of it) is present.
func QueryText(text string) []string {
	res, err := db.Query("SELECT DISTINCT filename FROM texts WHERE text like ?", "%"+text+"%")
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	return getFilenamesFromRes(res)
}

// Queries for images that have more than a certain proportion of a certain color.
func QueryColor(color string, amount float64) []string {
	res, err := db.Query("SELECT DISTINCT filename FROM colors WHERE color = ? and amount >= ?", color, amount)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	return getFilenamesFromRes(res)
}
func StoreTexts(filename string, data *pb.AnnotateImageResponse, tx *sql.Tx) {
	for _, text := range data.GetTextAnnotations() {
		storeText(filename, text.Description, tx)
	}
}

// Queries for tags associated to a file
func QueryImage(filename string) []string {
	res, err := db.Query("SELECT DISTINCT labels.description FROM labels, imagelabels where imagelabels.mid = labels.mid and imagelabels.filename = ?", filename)
	defer res.Close()

	if err != nil {
		log.Fatal(err)
	}
	return getFilenamesFromRes(res)
}

func storeText(filename string, text string, tx *sql.Tx) {
	stmt, err := db.Prepare("INSERT OR IGNORE INTO texts values(?, ?)")
	if err != nil {
		tx.Rollback()
		log.Fatal("can't create statement", err)
	}
	_, err = tx.Stmt(stmt).Exec(filename, text)
	if err != nil {
		tx.Rollback()
		log.Fatal("can't insert", err)
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
	stmt, err := db.Prepare("INSERT OR IGNORE INTO labels values(?, ?)")

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
	stmt, err := db.Prepare("INSERT OR REPLACE INTO imagelabels values(?, ?, ?)")

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
