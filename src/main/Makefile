all: tagdirectory querylabel querytext querycolor queryimage

tagdirectory: tagdirectory.go tagfetching.go dboperations.go dominantcolor.go
	go build $^

querylabel: querylabel.go dboperations.go dominantcolor.go
	go build $^

querytext: querytext.go dboperations.go dominantcolor.go
	go build $^

querycolor: querycolor.go dboperations.go dominantcolor.go
	go build $^

queryimage: queryimage.go dboperations.go dominantcolor.go
	go build $^

clean:
	rm -f tagdirectory querylabel querytext querycolor queryimage
