package route

import (
	"fmt"
	"net/http"
	"os"
	"soci-html-cdn/util"
)

// MoveFile takes the temp file and renames it to match the url
func MoveFolder(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		util.SendResponse(w, "", 200)
		return
	}

	fmt.Println("moving a folder")
	r.ParseMultipartForm(1 << 30)

	// Get the user's email if we're authorized
	bearerToken := r.Header.Get("Authorization")
	fmt.Println(bearerToken)
	user, err := util.GetUserEmail(bearerToken)
	fmt.Println(user)
	if err != nil {
		util.SendError(w, "User is not authorized.", 400)
		return
	}

	// Parse our url, and check if the url is available
	url := r.FormValue("url")
	urlIsAvailable, err := util.CheckIfURLIsAvailable(url)
	if err != nil {
		util.SendError(w, fmt.Sprintf("Error checking requested url: %v", url), 500)
		fmt.Println(err)
		return
	}
	if urlIsAvailable == false {
		util.SendError(w, fmt.Sprintf("Url \"%v\" is taken.", url), 400)
		return
	}

	// Check if the file we're moving exists
	tempFolder := r.FormValue("oldUrl")
	if _, err := os.Stat(fmt.Sprintf("files/temp-html/%v", tempFolder)); os.IsNotExist(err) {
		util.SendError(w, "No temp folder exists with that name.", 400)
		fmt.Println(err)
		return
	}

	// If everything else looks good, lets move the folder
	err = os.Rename(fmt.Sprintf("files/temp-html/%v", tempFolder), fmt.Sprintf("files/html/%v", url))
	if err != nil {
		util.SendError(w, "Error renaming folder.", 500)
		return
	}

	// Send back a response.
	util.SendResponse(w, url, 200)
}
