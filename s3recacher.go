package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/AdRoll/goamz/aws"
	"github.com/AdRoll/goamz/s3"
)

var (
	bucketName            string
	prefix                string
	delimiter             string
	fileName              string
	results               int = 0
	lastMarker            string
	maxObjs               int
	stopMarker            string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	cacheAge              string
	doStop                bool = false
)

func init() {
	flag.StringVar(&bucketName, "bucket", "", "Bucket Name")
	flag.StringVar(&prefix, "prefix", "", "path within bucket")
	flag.StringVar(&delimiter, "delimiter", "/", "delimeter to use")
	flag.StringVar(&lastMarker, "startid", "", "Object to start with")
	flag.StringVar(&stopMarker, "stopid", "", "Object to stop at")
	flag.IntVar(&maxObjs, "maxobjs", 0, "Maximum number of objects to perform against")
	flag.StringVar(&cacheAge, "age", "", "Cache-Age")
	flag.StringVar(&AWS_ACCESS_KEY_ID, "AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID"), "AWS_ACCESS_KEY_ID")
	flag.StringVar(&AWS_SECRET_ACCESS_KEY, "AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY"), "AWS_SECRET_ACCESS_KEY")

}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	if AWS_ACCESS_KEY_ID == "" || AWS_SECRET_ACCESS_KEY == "" {
		log.Fatal("AWS Credentials Required")
	}

	os.Setenv("AWS_ACCESS_KEY_ID", AWS_ACCESS_KEY_ID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", AWS_SECRET_ACCESS_KEY)

	// Since we're not messing with cacheAge, it's just easier to treat is as a string.
	if cacheAge == "" {
		cacheAge = "0"
	}

	if maxObjs != 0 || stopMarker != "" {
		// Set the conditional bit to check to stop
		doStop = true
	}

	log.Println("Starting Cache Alterations:")

	//  Connect to AWS using goamz
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Panic(err.Error())
	}

	// Instantiate S3 Object
	s := s3.New(auth, aws.USEast)

	// Set the Bucket
	Bucket := s.Bucket(bucketName)

	doUpdate(prefix, Bucket)

	log.Println("Wrote to", results, " S3 Objects. Last object was:", lastMarker)
}

func doUpdate(newprefix string, Bucket *s3.Bucket) {
	// Initial Request - Outside Loop
	Response, err := Bucket.List(newprefix, delimiter, lastMarker, 1000)
	if err != nil {
		log.Panic(err.Error())
	}

	for _, v := range Response.CommonPrefixes {
		fmt.Printf("recursing into: %s\n", v)
		doUpdate(v, Bucket)
	}

	// Set up the header for iterating.
	opts := s3.CopyOptions{}
	opts.CacheControl = "max-age=" + cacheAge
	opts.MetadataDirective = "REPLACE"

	log.Println("-> 0 START")

	// Loop Results
	for _, v := range Response.Contents {
		if v.Key == newprefix { // this is the directory itself
			continue
		}
		fmt.Printf("updating: %s\n", bucketName+"/"+v.Key)
		_, err := Bucket.PutCopy(v.Key, s3.PublicRead, opts, bucketName+"/"+v.Key)
		if err != nil {
			log.Panic(err.Error())
		}
		// We generate our own lastMarker.  This allows us to perform our own resume.
		lastMarker = v.Key
		results++

		if doStop == true {
			if results == maxObjs || lastMarker == stopMarker {
				return // End here.
			}
		}
	}

	fmt.Printf("\n")
	log.Println("->", results, " ", lastMarker)

	// Did Amazon say there was more?  If so, keep going.
	if Response.IsTruncated == true {
		for {
			// Issue List Command
			Response, err := Bucket.List(newprefix, delimiter, lastMarker, 1000)
			if err != nil {
				panic(err.Error())
			}

			// Loop through Response and dump it to the console.
			for _, v := range Response.Contents {
				fmt.Printf(".") // Indicator that something is happening
				_, err := Bucket.PutCopy(v.Key, s3.PublicRead, opts, bucketName+"/"+v.Key)
				if err != nil {
					log.Panic(err.Error())
				}
				lastMarker = v.Key
				results++

				if doStop == true {
					if results == maxObjs || lastMarker == stopMarker {
						return // End here.
					}
				}
			}

			if Response.IsTruncated == false {
				return // End loop
			} else {
				fmt.Printf("\n")
				log.Println("->", results, " ", lastMarker)
			}
		}
	}

}
