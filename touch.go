package main

import (
	"flag"
	"log"
	"os"
	"syscall"
	"time"
)

type options struct {
	arguments              []string
	changeAccessTime       bool
	noCreate               bool
	date                   string
	changeModificationTime bool
	referenceFile          string
	timestamp              string
	time                   string
}

func main() {
	// Parse flags & args
	options := getFlags()

	// Initialise default time values
	modTime := time.Now()
	accessTime := time.Now()

	// If we have a date to use
	if options.date != "" {
		t, err := time.Parse(time.RFC3339, options.date)
		if err != nil {
			log.Fatal(err)
		}

		modTime, accessTime = t, t
	}

	if options.timestamp != "" {
		layout := "200601021504.05"
		t, err := time.Parse(layout, options.timestamp)
		if err != nil {
			log.Fatal(err)
		}

		modTime, accessTime = t, t
	}

	// Set times to that of the reference file if needed
	if options.referenceFile != "" {
		modTime, accessTime = getReferenceFile(options.referenceFile)
	}

	// Get FileInfo of the file to be updated
	files := getFiles(options)

	for _, file := range files {
		// If we only want to change access time, reset mod time to the existing mod time
		if options.changeAccessTime {
			modTime = file.ModTime()
		}

		// If we only want to change mod time, reset access time to the existing access time
		if options.changeModificationTime {
			accessTime = time.Unix(0, file.Sys().(*syscall.Win32FileAttributeData).LastAccessTime.Nanoseconds())
		}

		// Update the file
		err := os.Chtimes(file.Name(), accessTime, modTime)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func getFiles(options options) []os.FileInfo {
	var files []os.FileInfo

	for _, arg := range options.arguments {
		// Get current file data (if necessary)
		fileInfo, err := os.Stat(arg)

		if err != nil {
			// If we get an error here the file doesn't exist
			if options.noCreate {
				continue
			}

			// Create the file
			file, err := os.Create(arg)
			if err != nil {
				log.Fatal(err)
			}

			// Get the FileInfo
			fileInfo, err = file.Stat()
			if err != nil {
				log.Fatal(err)
			}
		}
		files = append(files, fileInfo)
	}

	return files
}

func getReferenceFile(file string) (time.Time, time.Time) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}

	modTime := fileInfo.ModTime()
	accessTime := time.Unix(0, fileInfo.Sys().(*syscall.Win32FileAttributeData).LastAccessTime.Nanoseconds())

	return modTime, accessTime
}

func getFlags() options {
	opts := options{}

	// Get each expected flag value
	flag.BoolVar(&opts.changeAccessTime, "a", false, "Change only the access time")
	flag.BoolVar(&opts.noCreate, "c", false, "Do not create any files")
	flag.BoolVar(&opts.noCreate, "no-create", false, "Do not create any files")
	flag.StringVar(&opts.date, "d", "", "Parse STRING and use it instead of current time")
	flag.StringVar(&opts.date, "date", "", "Parse STRING and use it instead of current time")
	flag.BoolVar(&opts.changeModificationTime, "m", false, "Change only the modification time")
	flag.StringVar(&opts.referenceFile, "r", "", "Use this file's times instead of current time")
	flag.StringVar(&opts.referenceFile, "reference", "", "Use this file's times instead of current time")
	flag.StringVar(&opts.timestamp, "t", "", "Use [[CC]YY]MMDDhhmm[.ss] instead of current time")
	flag.StringVar(&opts.time, "time", "", "Change the specified time:\nWORD is access, atime, or use: equivalent to -a\nWORD is modify or mtime: equivalent to -m")
	flag.Parse()

	// Do extra validation on time flag, since it has expected values
	if opts.time == "access" || opts.time == "atime" {
		opts.changeAccessTime = true
	} else if opts.time == "modify" || opts.time == "mtime" {
		opts.changeAccessTime = true
	} else if opts.time != "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// If the argument is empty, bail out
	if len(flag.Args()) < 1 {
		log.Fatal("Expected at least 1 argument")
	}

	// Get the 1st argument passed to the command
	opts.arguments = flag.Args()

	return opts
}
