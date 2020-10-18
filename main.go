package main

import (
    "os"
    "fmt"
    "net/http"
    "log"
    "github.com/gorilla/mux"
)

func upload(w http.ResponseWriter, r *http.Request) {
    log.Print("Upload Endpoint Hit")
    w.Header().Set("Content-Type", "application/json")

    filename, err := saveReqFile(r, "myFile")
    if (err != nil) {
        log.Printf("Failed to upload file: /%s\n", err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    defer os.Remove(filename)

    // return that we have successfully uploaded our file!
    w.WriteHeader(http.StatusOK)
    log.Printf("Successfully Uploaded File in %s", filename)
    fmt.Fprintf(w, "Successfully Uploaded File in %s\n", filename)
}

func findBoxes(w http.ResponseWriter, r *http.Request) {
    log.Print("Find Boxes Endpoint Hit")
    w.Header().Set("Content-Type", "application/json")

    inputFilepath, err := saveReqFile(r, "myFile")
    if (err != nil) {
        log.Printf("Failed to upload file: /%s\n", err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    defer os.Remove(inputFilepath)
    log.Printf("Saved Uploaded File Temporarily to %s", inputFilepath)

    res, err := runPyScript("find_face", inputFilepath)
    if (err != nil) {
        log.Printf("Failed to Run Python Script to Find Boxes: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, res)
}

func savePerson(w http.ResponseWriter, r *http.Request) {
    log.Print("Save Person Endpoint Hit")
    w.Header().Set("Content-Type", "application/json")

    inputFilepath, err := saveReqFile(r, "myFile")
    if err != nil {
        log.Printf("Failed to upload file: /%s\n", err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    defer os.Remove(inputFilepath)
    log.Printf("Saved Uploaded File Temporarily to %s", inputFilepath)

    // ROBERT DO YOUR SHIT HERE AND FILL IN VAR ID
    id := "thisismyid"
    // END OF ROBERT'S SHIT

    res, err := runPyScript("save_face", inputFilepath, id)
    if err != nil {
        log.Printf("Failed to Run Python Script to Save Face: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
    }
    if res == "{\"worked\": false}" {
        log.Printf("Could Not Encode Face with id %s", id)
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, res)
}

func matchPerson(w http.ResponseWriter, r *http.Request) {
    log.Print("Save Person Endpoint Hit")
    w.Header().Set("Content-Type", "application/json")

    inputFilepath, err := saveReqFile(r, "myFile")
    if err != nil {
        log.Printf("Failed to upload file: /%s\n", err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    defer os.Remove(inputFilepath)
    log.Printf("Saved Uploaded File Temporarily to %s", inputFilepath)


    id, err := runPyScript("match_face", inputFilepath)
    if err != nil {
        log.Printf("Failed to Run Python Script to Match Face: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
    }

    w.WriteHeader(http.StatusOK)
    // ROBERT DO YOUR SHIT HERE AND FILL IN RES JSON
    fmt.Fprint(w, id)
    // END OF ROBERT'S SHIT
}

func status(w http.ResponseWriter, r *http.Request) {
    log.Print("Status Endpoint Hit")
    fmt.Fprintf(w, "Online\n")
}

func setupRoutes() {
    log.Print("Setting Up Routes...")
    r := mux.NewRouter()
    r.HandleFunc("/status", status).Methods("GET")
    r.HandleFunc("/upload", upload).Methods("POST")
    r.HandleFunc("/boxes", findBoxes).Methods("POST")
    r.HandleFunc("/save", savePerson).Methods("POST")
    r.HandleFunc("/match", matchPerson).Methods("POST")

    log.Printf("Listening on port %s...", port)
    log.Fatal(http.ListenAndServe(":" + port, r))
}

func main() {
    if _, err := os.Stat(tempimgfolder); os.IsNotExist(err) {
        log.Printf("Temp Image Folder Does Not Exist. Making One at %s", tempimgfolder)
        os.Mkdir(tempimgfolder, 0700)
    }
    log.Print("Starting Backend Server...")
    setupRoutes()
}