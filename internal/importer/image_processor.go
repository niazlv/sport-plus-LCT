package importer

import "log"

// ProcessImage принимает ссылку на изображение и возвращает её же
func ProcessImage(imageURL string) string {
	// Пока просто возвращаем ссылку, позже добавим логику загрузки и изменения ссылки
	log.Printf("Processing image: %s", imageURL)
	return imageURL
}
