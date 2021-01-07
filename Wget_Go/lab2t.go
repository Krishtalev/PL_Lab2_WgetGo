package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type PassThru struct {
	total int64 // Total # of bytes transferred
}

// Write 'overrides' the underlying io.Writer's Write method.
func (pt *PassThru) Write(p []byte) (n int, err error) {
	b := len(p)
	pt.total += int64(b)
	return b, nil
}

// Calling main
func main() {
	//https://upload.wikimedia.org/wikipedia/commons/f/ff/Pizigani_1367_Chart_10MB.jpg
	//http://www.calprog.com/Sounds/CP2008_IZZ_Montage.mp3
	var url, fileName string
	fmt.Println("Examples of url:")
	fmt.Println("https://upload.wikimedia.org/wikipedia/commons/f/ff/Pizigani_1367_Chart_10MB.jpg")
	fmt.Println("http://www.calprog.com/Sounds/CP2008_IZZ_Montage.mp3")

	fmt.Println("Enter url: ")
	fmt.Scan(&url)

	fileName = strings.Split(url,"/")[len(strings.Split(url,"/"))-1]
	fmt.Println("Downloading", url, "to", fileName)

	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	traf := PassThru{0}
	quit := make(chan bool)
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <- ticker.C:
				fmt.Println("downloaded",traf.total,"bytes")
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()

	n, err := io.Copy(output, io.TeeReader(response.Body, &traf))
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	quit<-true

	fmt.Println(n, "bytes downloaded.")
}

