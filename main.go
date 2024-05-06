package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Структура для парсинга ответа API
type DogAPIResponse struct {
	Message interface{} `json:"message"`
	Status  string      `json:"status"`
}

func getRandomImage() (string, error) {
	response, err := http.Get("https://dog.ceo/api/breeds/image/random")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var data DogAPIResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	url, ok := data.Message.(string)
	if !ok {
		return "", fmt.Errorf("unexpected type for random image URL")
	}

	return url, nil
}

func getBreedImages(breed string) ([]string, error) {
	breed = strings.ReplaceAll(breed, "-", "/")
	url := fmt.Sprintf("https://dog.ceo/api/breed/%s/images/random/3", breed)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var data DogAPIResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	images, ok := data.Message.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for breed images")
	}

	var imageUrls []string
	for _, img := range images {
		imgStr, ok := img.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected type for image URL")
		}
		imageUrls = append(imageUrls, imgStr)
	}

	return imageUrls, nil
}

func extractBreedFromURL(url string) string {
	parts := strings.Split(url, "/")
	breed := parts[4]
	return breed
}
func DownloadImage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	filename := "downloaded_image.jpg"
	err = os.WriteFile(filename, body, 0644)
	if err != nil {
		return "", err
	}
	return filename, nil
}
func main() {
	http.HandleFunc("/GetRandomBreeds", func(w http.ResponseWriter, r *http.Request) {
		randomImageURL, err := getRandomImage()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		breed := extractBreedFromURL(randomImageURL)
		file, err := DownloadImage(randomImageURL)
		http.HandleFunc("/"+breed, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, file)
		})
		breedImages, err := getBreedImages(breed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<h1>%s</h1><div style=\"display:flex\">", breed)
		for _, img := range breedImages {
			fmt.Fprintf(w, "<img src=\"%s\" style=\"height:360px;width:480px;\" alt=\"Dog image\"/>", img)
		}
		fmt.Println(breed)
		fmt.Fprintf(w, "</div>")
		fmt.Fprintf(w, "<script>")
		fmt.Fprintf(w, "window.onload = function() {")
		fmt.Fprintf(w, "fetch('/"+breed+"', {")
		fmt.Fprintf(w, "method: 'GET',")
		fmt.Fprintf(w, "headers: {'Content-Type': 'application/json'},")
		fmt.Fprintf(w, "})")
		fmt.Fprintf(w, ".then(response => response.blob())")
		fmt.Fprintf(w, ".then(blob => {")
		fmt.Fprintf(w, "const url = window.URL.createObjectURL(new Blob([blob]));")
		fmt.Fprintf(w, "const link = document.createElement('a');")
		fmt.Fprintf(w, "link.href = url;")
		fmt.Fprintf(w, "link.setAttribute('download', 'downloaded_image.jpg');")
		fmt.Fprintf(w, "document.body.appendChild(link);")
		fmt.Fprintf(w, "link.click();")
		fmt.Fprintf(w, "});")
		fmt.Fprintf(w, "};")
		fmt.Fprintf(w, "</script>")
	})
	fmt.Println("Server is running at http://localhost:8080/GetRandomBreeds")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
