package wallpaper

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Wallpaper struct {
	Id    string `json:"id"`
	Urls  Urls   `json:"urls"`
	Links Links  `json:"links"`
}

type Urls struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
}

type Links struct {
	Download string `json:"download"`
}

func HandleErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func CreateFolder() {
	// making the directory to store the wallpapers, if not already exists
	if err := os.Mkdir("/home/sufyaan/wallpaper_app_images", os.ModePerm); err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}
}

func FetchWallpaperIntoFolder() {
	// getting the wallpaper object from unsplash API
	var wallpaper Wallpaper
	wallpaper, err := getWallpaperObj()
	HandleErr(err)

	// getting the image and storing it in the image directory
	res, err := http.Get(wallpaper.Links.Download)
	HandleErr(err)

	defer res.Body.Close()

	// fname := wallpaper.Id
	fname := "current"
	f, err := os.Create("/home/sufyaan/wallpaper_app_images/" + fname)
	HandleErr(err)

	defer f.Close()
	_, err = f.ReadFrom(res.Body)
	HandleErr(err)
}

func getWallpaperObj() (Wallpaper, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", "https://api.unsplash.com/photos/random?orientation=landscape&topics=bo8jQKTaE0Y,CDwuwXJAbEw", nil)
	if err != nil {
		return Wallpaper{}, err
	}
	authKey := "Client-ID " + os.Getenv("ACCESS_KEY")
	req.Header.Set("Authorization", authKey)
	res, err := client.Do(req)
	if err != nil {
		return Wallpaper{}, err
	}

	var wallpaper Wallpaper

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&wallpaper)
	if err != nil {
		return Wallpaper{}, err
	}

	return wallpaper, nil
}

func SetWallpaper() error {
	var out bytes.Buffer
	var file string = "/home/sufyaan/wallpaper_app_images/current"
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "color-scheme")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	mode := out.String()
	if strings.Compare(mode, "'prefer-dark'\n") == 0 {
		return exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", strconv.Quote("file://"+file)).Run()
	} else {
		return exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", strconv.Quote("file://"+file)).Run()
	}
}

func GetCurrentWallpaper() (string, error) {
	colorScheme, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "color-scheme").Output()
	if err != nil {
		panic(err)
	}
	mode := string(colorScheme)
	if strings.Compare(mode, "'prefer-dark'\n") == 0 {
		output, err := exec.Command("gsettings", "get", "org.gnome.desktop.background", "picture-uri-dark").Output()
		outputStr := string(output)
		outputStr = outputStr[8 : len(outputStr)-2]
		return outputStr, err
	} else {
		output, err := exec.Command("gsettings", "get", "org.gnome.desktop.background", "picture-uri").Output()
		outputStr := string(output)
		outputStr = outputStr[8 : len(outputStr)-2]
		return outputStr, err
	}
}
