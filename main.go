package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var storageFolder string = ""
var addr string = ":8080"
var logRequests bool = false
var form string = `<!DOCTYPE html>
<html lang="en">
   <head>
      <meta charset="UTF-8" />
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      <meta http-equiv="X-UA-Compatible" content="ie=edge" />
   </head>
   <body>
      <form enctype="multipart/form-data" action="/p" method="post">
         <input type="file" name="file" />
         <input type="submit" value="Upload" />
      </form>
   </body>
</html>`

func logRequest(req *http.Request, endpoint string) {
	if logRequests {
		log.Printf("[REQUEST] %s %s %s - User-Agent: %s - Remote: %s",
			req.Method, endpoint, req.URL.String(), req.UserAgent(), req.RemoteAddr)
	}
}

func exfilGet(w http.ResponseWriter, req *http.Request) {
	logRequest(req, "/g")
	
	// We could also check for req.Header.Get("X-REAL-IP") or "X-FORWARDED-FOR" (better with proxies)
	host := strings.Split(req.RemoteAddr, ":")

	var filename string = fmt.Sprintf("%s_%s.txt", host[0], time.Now().Format("2006-01-02_15-04-05"))
	var filePath = path.Join(storageFolder, filename)

	out, err := os.Create(filePath)

	if err != nil {
		fmt.Println(err)
	} else {
		req.Write(out)
		fmt.Printf("[*] Request Stored (%s)\n", filename)
	}

}

func uploadForm(w http.ResponseWriter, req *http.Request) {
	logRequest(req, "/")
	fmt.Fprint(w, form)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	logRequest(r, "/p")

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	ioutil.WriteFile(path.Join(storageFolder, handler.Filename), fileBytes, 0644)

	fmt.Fprintf(w, "Done\n")
	fmt.Printf("[+] File Uploaded (%s)\n", handler.Filename)
}

func loggedFileServer(fs http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, "/l")
		fs.ServeHTTP(w, r)
	})
}

func setupRoutes() {
	http.HandleFunc("/", uploadForm)
	http.HandleFunc("/p", uploadFile)
	http.HandleFunc("/g", exfilGet)

	// Directory Listing with logging
	http.Handle("/l/", loggedFileServer(http.StripPrefix("/l", http.FileServer(http.Dir(storageFolder)))))

	//HTTP or HTTPs
	if _, err := os.Stat("HTTPUploadExfil.csr"); err == nil {
		fmt.Println("[+] Starting in HTTPs mode")
		http.ListenAndServeTLS(addr, "HTTPUploadExfil.csr", "HTTPUploadExfil.key", nil)
	} else {
		fmt.Println("[+] Starting in HTTP mode")
		http.ListenAndServe(addr, nil)
	}
}

func ascii_art() {

	var ascii string = `
    _____ _____ _____ _____ _____     _           _ _____     ___ _ _ 
   |  |  |_   _|_   _|  _  |  |  |___| |___ ___ _| |   __|_ _|  _|_| |
   |     | | |   | | |   __|  |  | . | | . | .'| . |   __|_'_|  _| | |
   |__|__| |_|   |_| |__|  |_____|  _|_|___|__,|___|_____|_,_|_| |_|_|
                                 |_|                                  
   `

	fmt.Print(ascii)
	fmt.Printf("\n")
	fmt.Println("Version: 2025-07-27")
	fmt.Println("Usage: ./httpuploadexfil [-log] :8080 /home/kali/exfil")
	fmt.Println("	-log: Enable HTTP request logging")
	fmt.Printf("\n")
}

func main() {
	flag.BoolVar(&logRequests, "log", false, "Enable HTTP request logging")
	flag.Parse()

	ascii_art()

	args := flag.Args()
	if len(args) <= 1 {
		storageFolder, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		fmt.Println("[+] Server Running")
		fmt.Printf("[+] Default Settings: Addr '%s'; Folder '%s'\n", addr, storageFolder)
	} else {
		addr = args[0]
		storageFolder = args[1]

		if _, err := os.Stat(storageFolder); os.IsNotExist(err) {
			err := os.Mkdir(storageFolder, 0755)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		}

		fmt.Println("[+] Server Running")
		fmt.Printf("[+] Settings: Addr '%s'; Folder '%s'\n", addr, storageFolder)
	}

	if logRequests {
		fmt.Println("[+] HTTP request logging enabled")
	}

	setupRoutes()
}
