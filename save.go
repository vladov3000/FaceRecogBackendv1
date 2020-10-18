package main

import (
	"strings"
	"net/http"
	"log"
	"io/ioutil"
	"errors"
	"encoding/base64"
	"image/png"
	"fmt"
	"os"
	"image"
)

func saveReqFile(r *http.Request, filename string) (string, error) {
	contentType := r.Header.Get("Content-type")
	log.Printf("Request Headers:\n%s", formatRequest(r))
	if (strings.HasPrefix(contentType, "application/x-www-form-urlencoded")){
		log.Print("Header type detected as application/x-www-form-urlencoded")
		return saveURL(r)
	} else if strings.HasPrefix(contentType, "multipart/form-data") {
		log.Print("Header type detected as multipart/form-data")
		return saveFormDataFile(r, filename)
	}
	log.Print("Cannot Save File Due to Invalid Content Type for Request")
	return "", errors.New("Invalid Content Type")
}

func saveURL(r *http.Request) (string, error) {
	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(tempimgfolder, "upload-*.png")
	if err != nil {
		log.Printf("Error Creating Temp File for URL from Request: %s", err)
		return "", err
	}

	/*
	// Decode base 64
	rBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error Reading Bytes from Request: %s", err)
		return "", err
	}

	var unbased []byte;
	_, err = base64.StdEncoding.Decode(unbased, rBytes)
	if err != nil {
		log.Printf("Error Decoding Base64 from URL: %s", err)
		return "", err
	}

	reader := bytes.NewReader(unbased)
	im, err := png.Decode(reader)
	if err != nil {
		log.Printf("Error Decoding PNG Bytes from URL: %s", err)
		return "", err
	}
	// write this png to our temporary file
	png.Encode(tempFile, im)*/

	// dataBytes, _ := ioutil.ReadAll(r.Body)
	// data := string(dataBytes)

	reader := base64.NewDecoder(base64.StdEncoding, r.Body)
    m, formatString, err := image.Decode(reader)
    if err != nil {
        log.Fatal(err)
    }
    bounds := m.Bounds()
    fmt.Println(bounds, formatString)

    //Encode from image format to writer
    pngFilename := tempFile.Name()
    f, err := os.OpenFile(pngFilename, os.O_WRONLY|os.O_CREATE, 0777)
    if err != nil {
        log.Fatal(err)
        return "", err
    }

    err = png.Encode(f, m)
    if err != nil {
        log.Fatal(err)
        return "", err
    }
    fmt.Println("Png file", pngFilename, "created")
	
	log.Printf("Successfully uploaded file to %s...", tempimgfolder)
	return tempFile.Name(), nil
}

func saveFormDataFile(r *http.Request, filename string) (string, error) {
	log.Printf("Uploading file to %s...", tempimgfolder)

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile(filename)
	if err != nil {
		log.Printf("Error Retrieving File '%s' from Request %s: %s", filename, err)
		return "", err
	}
	//defer file.Close()
	log.Printf("Uploaded File: %+v", handler.Filename)
	log.Printf("File Size: %+v", handler.Size)
	log.Printf("MIME Header: %+v", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(tempimgfolder, "upload-*.png")
	if err != nil {
		log.Printf("Error Creating Temp File '%s' from Request: %s", filename, err)
		return "", err
	}

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Error Reading Bytes from File '%s' from Request: %s", filename, err)
		return "", err
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	
	log.Printf("Successfully uploaded file to %s...", tempimgfolder)
	return tempFile.Name(), nil
}