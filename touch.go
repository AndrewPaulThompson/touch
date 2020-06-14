package main

import (
	"flag"
	"fmt"
)

type options struct {
	changeAccessTime       bool
	noCreate               bool
	date                   string
	changeModificationTime bool
	referenceFile          string
	timestamp              string
	time                   string
}

func main() {

	options := getFlags()
	fmt.Println(options)

}

func getFlags() options {
	opts := options{}

	flag.BoolVar(&opts.changeAccessTime, "a", false, "Change only the access time")
	flag.BoolVar(&opts.noCreate, "c", false, "Do not create any files")
	flag.BoolVar(&opts.noCreate, "no-create", false, "Do not create any files")
	flag.StringVar(&opts.date, "d", "", "Parse STRING and use it instead of current time")
	flag.StringVar(&opts.date, "date", "", "Parse STRING and use it instead of current time")
	flag.BoolVar(&opts.changeModificationTime, "m", false, "Change only the modification time")
	flag.StringVar(&opts.referenceFile, "r", "", "Use this file's times instead of current time")
	flag.StringVar(&opts.referenceFile, "reference", "", "Use this file's times instead of current time")
	flag.StringVar(&opts.timestamp, "t", "", "Use [[CC]YY]MMDDhhmm[.ss] instead of current time")
	flag.StringVar(&opts.time, "time", "", "Change the specified time: WORD is access, atime, or use: equivalent to -a WORD is modify or mtime: equivalent to -m")

	flag.Parse()

	return opts
}
