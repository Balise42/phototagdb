# phototagdb
A small collection of utilities to tag pictures in a directory thanks to the Google Cloud Vision API, store and search through tags.

# How to build to a useable state
* checkout code
* Install the following libraries:
	* "github.com/golang/protobuf/proto"
	* "github.com/mattn/go-sqlite3"
	* "google.golang.org/genproto/googleapis/cloud/vision/v1"
	* "github.com/jkl1337/go-chromath"
	* "github.com/jkl1337/go-chromath/deltae"
	* "google.golang.org/genproto/googleapis/type/color"
	* "golang.org/x/net/context"
	* "gopkg.in/h2non/bimg.v1" (also requires the installation of the libvips C library: https://jcupitt.github.io/libvips/)
* install sqlite3
* create the sqlite3 database as $PHOTOTAGDB\_FOLDER/resources/imgtag.db, create the tables given in $PHOTOTAGDB\_FOLDER/resources/create\_db.sql
* make (should build all executables)
* voilà!

# How to use
* You need a Google Cloud Vision API key - see https://cloud.google.com/vision/docs/auth - and export the GOOGLE\_APPLICATION\_CREDENTIALS variable as indicated in the documentation. NOTE: to get an API key, you need to provide billing information.
* `./tagdirectory <directory>`: tags all the JPG images of a directory. NOTE: this uploads and uses the Google API (and counts towards quotas and monthly limits for said API)
* `./querylabel <label1> <label2>...`: returns all the images references from the DB that contain all the labels given as argument. Does not access the Google API.
* `./querytext <text>`: returns all the images containing the provided text as characters in the image (as stored in the database). Does not access the Google API.
* `./querycolor <color> <amount>`: returns all the images containing more than the specified amount (between 0.0 and 1.0) of a color, as stored in the database. Does not access the Google API.

# TODO
* add proper tests (okay, any tests)
* make things better when it comes to error handling (PROBABLY it's not necessary to come to a screeching halt every time ;) )
* make the DB location not hardcoded, make the DB initialization less manual
* add more query tools (tags associated to a given picture, all the labels...)
* make labels case insensitive (it's an issue on locations)
* add a "recompute stored data from existing API results" to be able to modify the treatment of the raw data without accessing the API
* the color detection can probably be improved by adding more colors to the predefined ones. could also maybe add qualifiers ("dark green" vs "light green" for instance)
* make the query tool all-in-one, so that it can also for instance be used to search for yellow ducks
* a visualization of the images would be very cool instead of just a list of files
