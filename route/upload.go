package route

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"soci-html-cdn/util"

	"github.com/google/uuid"
)

// UploadFile takes the form upload and delegates to the encoders
func UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		util.SendResponse(w, "", 200)
		return
	}
	// Parse our multipart form, set a 1GB max upload size
	r.ParseMultipartForm(1 << 30)

	// Get the user's email if we're authorized
	bearerToken := r.Header.Get("Authorization")
	fmt.Println(bearerToken)
	user, err := util.GetUserEmail(bearerToken)
	fmt.Println(user)
	if err != nil {
		util.SendError(w, fmt.Sprintf("User is not authorized. Token: %v", bearerToken), 400)
		return
	}

	// Parse our url, and check if the url is available
	url := r.FormValue("url")
	if url != "" {
		urlIsAvailable, err := util.CheckIfURLIsAvailable(url)
		if err != nil {
			util.SendError(w, "Error checking requested url.", 500)
			return
		}
		if urlIsAvailable == false {
			util.SendError(w, "Url is taken.", 400)
			return
		}
	} else {
		url = uuid.New().String()
	}

	formdata := r.MultipartForm

	files := formdata.File["files"]
	// Parse our file and assign it to the proper handlers depending on the type
	if err != nil {
		util.SendError(w, "Error: no file was found in the \"files\" field, or they could not be parsed.", 400)
		return
	}

	// Let's create a temp folder to store the files in
	tempFolder, err := ioutil.TempDir("files/temp-html", "html-*")
	if err != nil {
		fmt.Println(err)
		return
	}

	os.Chmod(tempFolder, 0644)

	// Loop through the files and write them to the temp folder
	for i, _ := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			util.SendError(w, "Error opening file.", 500)
			return
		}

		out, err := os.Create(tempFolder + "/" + files[i].Filename)
		defer out.Close()
		if err != nil {
			util.SendError(w, "Error creating file.", 500)
			return
		}

		// Create a temp file
		_, err = io.Copy(out, file)
		if err != nil {
			util.SendError(w, "Error creating temp file.", 500)
			return
		}
	}

	// If all is good, let's log what the hell is going on
	fmt.Printf("%v is uploading a file of to %v\n", user, url)

	util.SendResponse(w, filepath.Base(tempFolder), 200)
}
