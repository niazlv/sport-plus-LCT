package main

import (
	"flag"
	"log"

	"github.com/niazlv/sport-plus-LCT/internal/importer"
)

func main() {
	// Путь до папки с JSON файлом и изображениями
	var folderPath string
	flag.StringVar(&folderPath, "path", ".", "Path to the folder containing main_images.json and images")
	flag.Parse()

	jsonFilePath := folderPath + "/main_images.json"
	// "http://sport-plus.sorewa.ru:8080/v1"
	err := importer.ImportExercisesFromJSON(jsonFilePath, "http://localhost:8080/v1")
	if err != nil {
		log.Fatal("Error importing courses: ", err)
	}

	log.Println("Data import completed successfully")
}
