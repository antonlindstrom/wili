package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/build", buildHandler).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}/log", logHandler).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}/run", runHandler).Methods("GET")
	http.Handle("/", r)

	log.Printf("Starting on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Build a container
func buildHandler(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()

	if err != nil {
		log.Printf("Error, could not read multipart: %s\n", err)
		return
	}

	err = ReadMultipart(reader)

	if err != nil {
		log.Printf("Error, could not read multipart: %s\n", err)
		return
	}

	// Run build
	build(w)
}

func build(w http.ResponseWriter) {
	cmd := exec.Command("./deploy.sh")
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Printf("Failed to execute deployment script.\n")
		return
	}

	process := bufio.NewReader(stdout)
	err = cmd.Start()

	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}

	for {
		line, _, err := process.ReadLine()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error: %s\n")
			return
		}

		fmt.Fprintf(w, "- %s\n", line)
	}
	cmd.Wait()
}

// DEPRECATED - Follow the log of a build
func logHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Fprintf(w, "Hello %s\n", id)
}

// DEPRECATED - Start the container
func runHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Fprintf(w, "Not implemented - %s\n", id)
}

// Read multipart
func ReadMultipart(r *multipart.Reader) error {
	for {
		part, err := r.NextPart()

		fileName := "build.tar.gz"

		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("%s\n", err)
			return err
		}

		file, ferr := os.OpenFile("/tmp/"+fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)

		if ferr != nil {
			log.Printf("Could not open file /tmp/%s\n", fileName)
			return ferr
		}

		defer file.Close()
		err = WriteToFile(file, part)

		if err != nil {
			log.Printf("%s\n", err)
			return err
		}
	}
	return nil
}

// Write multipart to file
func WriteToFile(file *os.File, part *multipart.Part) error {
	for {
		buffer := make([]byte, 1024)
		bufbytes, err := part.Read(buffer)

		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("%s\n", err)
			return err
		}

		file.Write(buffer[:bufbytes])
	}
	return nil
}
