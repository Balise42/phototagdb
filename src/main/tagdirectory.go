package main

import (
	"log"
	"os"
	"io/ioutil"
	"strings"
)

// Takes as argument a directory containing jpg images, tag them and add the tags to the
// database.
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

        TagFiles(toTreat)
}
