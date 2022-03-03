package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/fogleman/gg"
	"github.com/gofrs/uuid"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/gorilla/mux"
)

type request struct {
	Url  string `json:"url"`
	Text string `json:"text"`
	Auth bool   `json:"auth"`
}

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	//Init router
	r := mux.NewRouter()

	// Route handling and endpoints
	r.HandleFunc("/image", createInspImageHandler).Methods(http.MethodPost)
	r.HandleFunc("/user", userLoginHandler).Methods(http.MethodPost)

	// image file handler
	fs := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/image/").Handler(http.StripPrefix("/image/", fs))

	// static file handler
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("frontend/dist/")))
	r.PathPrefix("/").HandlerFunc(indexHandler("frontend/dist/index.html"))

	fmt.Println("Server listening on port 3000")

	err := http.ListenAndServe(":3000", r)
	if err != nil && err != http.ErrServerClosed {
		log.Println(err)
		os.Exit(1)
	}
}

func indexHandler(entrypoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}
}

func userLoginHandler(w http.ResponseWriter, r *http.Request) {
	//decode response
	decoder := json.NewDecoder(r.Body)

	var user user
	err := decoder.Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//hard coding a log in for now @TODO: add db + secure pw storing
	if user.Username == "test" && user.Password == "test" {
		//fmt.Println("Login successful!")
		fmt.Fprintf(w, `{"status": "success", "user":"%s", "msg":"Login successful!"}`, user.Username)
	} else {
		//fmt.Println("Login failed, wrong username or password.")
		fmt.Fprintf(w, `{"status": "fail", "msg":"Login failed, wrong username or password"}`)
	}
}

func createInspImageHandler(w http.ResponseWriter, r *http.Request) {
	//decode response
	decoder := json.NewDecoder(r.Body)

	var req request
	err := decoder.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Print("premium access: ")
	log.Println(req.Auth)

	//check the request
	if req.Url == "" || req.Text == "" {
		//fmt.Println("Error: Incomplete request.")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"Error: Incomplete request"}`)
		return
	}

	log.Printf("URL and text received, URL=%v, Text=%v", req.Url, req.Text)

	//get http response from url
	res, err := http.Get(req.Url)
	if err != nil {
		//fmt.Println("Invalid URL.")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"Error: Invalid URL"}`)
		return
	}

	//grab the image from the response body
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//fmt.Println("Invalid URL.")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"Error: Invalid Image data"}`)
		return
	}
	res.Body.Close()

	id, err := uuid.NewV4()
	if err != nil {
		//fmt.Println("Invalid URL.")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"Error: Invalid"}`)
		return
	}

	//place text over img
	err = textOverImg(data, req.Text, req.Auth, id.String())
	if err != nil {
		//no image from the given url
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"Error: Could not get image from URL"}`)
		return
	}

	fmt.Fprintf(w, `{"image": "http://localhost:3000/image/inspirational_image_%s.png", "error":"none"}`, id.String())
}

func textOverImg(imgData []byte, text string, premium bool, id string) error {
	//decode from []byte to image.Image
	img, _, err := image.Decode(bytes.NewReader(imgData))
	//if the url is not an image
	if err != nil {
		//fmt.Println("Could not get image from URL.")
		return err
	}
	//get image size
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	//load in a default font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return err
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 48})

	//create canvas for image & drawing text
	dc := gg.NewContext(imgWidth, imgHeight)
	dc.DrawImage(img, 0, 0)
	dc.SetFontFace(face)
	dc.SetColor(color.White)

	//set x/y position of text
	x := float64(imgWidth / 2)
	y := float64(imgHeight / 2)
	maxWidth := float64(imgWidth - 60) //maximum width text can occupy

	dc.DrawStringWrapped(text, x, y, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)
	//check users access
	if !premium {
		//draw a watermark
		dc.DrawStringAnchored("Inspirationifier: Free Version.", 325, y*2-48, 0.5, 0.5)
	}
	filename := "images/inspirational_image_" + id + ".png"
	err = dc.SavePNG(filename)
	if err != nil {
		return err
	}
	log.Println("Inspirational image created: ", filename)
	return nil
}
