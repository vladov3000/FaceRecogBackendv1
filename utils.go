package main

import (
	"os"
	"fmt"
	"strings"
    "net/http"
	"log"
	"os/exec"
	"path/filepath"
	"bytes"
	"errors"
)

func findPythonFolder(path ...string) (string, error) {
	envvar := os.Getenv(pfrenvvar)
	if envvar != "" {
		log.Printf("Using Python Folder from $%s", pfrenvvar)
		return envvar, nil
	}

	log.Printf("Finding Python Folder")
	pythonFolder, err := filepath.Abs(filepath.Join(path...))
	if err != nil {
		log.Printf("Couldn't Find Python Folder at %v: %s", path, err)
		return "", err
	}
	os.Setenv(pfrenvvar, pythonFolder)
	log.Printf("Found Python Folder at %s, set $%s to it", pythonFolder, pfrenvvar)
	return pythonFolder, nil
}

func runPyScript(name string, args ...string) (string, error) {
	log.Printf("Running Python Script")

	file2 := "";
	if len(args) < 1 {
		s := "Not Enough Files Provided to Run Python Script"
		log.Printf(s)
		return "", errors.New(s)
	} else if len(args) > 1 { 
		file2 = args[1] 
	}

	pythonFolder, _ := findPythonFolder("..", "FaceRecogPy")
	scriptPath := filepath.Join(pythonFolder, name + ".py")
	pythonExePath := os.Getenv("VIRTUAL_ENV") + "/bin/python"
	if p := os.Getenv(pythonexevar); p != "" {
		pythonExePath = p
	}

	var out bytes.Buffer
	var cmdErr bytes.Buffer

	cmd := exec.Command(pythonExePath, scriptPath, args[0], file2)
	cmd.Stdout = &out
	cmd.Stderr = &cmdErr
	err := cmd.Run()

	if err != nil {
		log.Printf("Failed to Run Command: %s", cmdErr.String())
		return "", err
	}
	log.Print(out.String())
	return out.String(), nil
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
	  name = strings.ToLower(name)
	  for _, h := range headers {
		request = append(request, fmt.Sprintf("%v: %v", name, h))
	  }
	}

	 // Return the request as a string
	 return strings.Join(request, "\n")
   }