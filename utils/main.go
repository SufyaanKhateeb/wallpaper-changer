package util

import (
	"fmt"

	Wallpaper "github.com/SufyaanKhateeb/wallpaper_app/wallpaper"
	"github.com/joho/godotenv"
)

func Mainfunc() {
	// load the env variables
	err := godotenv.Load(".env")
	Wallpaper.HandleErr(err)

	Wallpaper.CreateFolder()             // create wallpaper_app_images folder if not exists
	Wallpaper.FetchWallpaperIntoFolder() // fetch the wallpaper from the api into the above folder
	err = Wallpaper.SetWallpaper()       // set the wallpaper
	Wallpaper.HandleErr(err)

	// printing the wallpaper location
	background, err := Wallpaper.GetCurrentWallpaper()
	Wallpaper.HandleErr(err)
	fmt.Println("current wallpaper location:", background)
}
