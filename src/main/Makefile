all: tagdirectory querylabel querytext

tagdirectory: tagdirectory.go tagfetching.go dboperations.go dominantcolor.go
	go build $^

querylabel: querylabel.go dboperations.go dominantcolor.go
	go build $^

querytext: querytext.go dboperations.go dominantcolor.go
	go build $^

clean:
	rm -f tagdirectory querylabel querytext
