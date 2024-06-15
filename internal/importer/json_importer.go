package importer

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Exercise struct {
	OriginalUri       string   `json:"original_uri"`
	OriginalUrii      string   `json:"originalUri"`
	Name              string   `json:"name"`
	Muscle            string   `json:"muscle"`
	AdditionalMuscle  string   `json:"additional_muscle"`
	AdditionalMusclee string   `json:"additionalMuscle"`
	Type              string   `json:"type"`
	Equipment         string   `json:"equipment"`
	Difficulty        string   `json:"difficulty"`
	Photos            []string `json:"photos"`
}

type CreateExerciseResponse struct {
	ID int `json:"id"`
}

type Photo struct {
	URL string `json:"url"`
}

func ImportExercisesFromJSON(filePath string, apiBaseURL string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var exercises []Exercise
	if err := json.Unmarshal(byteValue, &exercises); err != nil {
		return err
	}

	imagesFolder := filepath.Dir(filePath)

	// Аутентификация
	token, err := Authenticate(apiBaseURL, "aboba", "passwd")
	if err != nil {
		return err
	}

	for _, exercise := range exercises {
		// Убедитесь, что все обязательные поля заполнены
		if exercise.OriginalUri == "" {
			if exercise.OriginalUrii == "" {
				exercise.OriginalUri = "default_uri" // Замените на реальное значение по умолчанию, если необходимо
			} else {
				exercise.OriginalUri = exercise.OriginalUrii
			}
		}
		if exercise.AdditionalMuscle == "" {
			if exercise.AdditionalMusclee == "" {
				exercise.AdditionalMuscle = "default_muscle" // Замените на реальное значение по умолчанию, если необходимо
			} else {
				exercise.AdditionalMuscle = exercise.AdditionalMusclee
			}
		}

		// Загружаем фотографии и добавляем их URL в упражнение
		var photos []Photo
		var urls []string
		for _, photo := range exercise.Photos {
			photoPath := filepath.Join(imagesFolder, photo)
			uploadedImageURL, err := uploadImage(apiBaseURL, photoPath, token)
			if err != nil {
				return err
			}

			photos = append(photos, Photo{URL: uploadedImageURL})
			urls = append(urls, uploadedImageURL)
			log.Printf("Uploaded image URL: %s\n", uploadedImageURL)
		}
		exercise.Photos = urls

		// Создаем упражнение с загруженными фотографиями
		_, err := createExercise(apiBaseURL, exercise, token)
		if err != nil {
			return err
		}
	}

	return nil
}

func createExercise(apiBaseURL string, exercise Exercise, token string) (int, error) {
	exerciseURL := apiBaseURL + "/exercise"
	exerciseData, err := json.Marshal(exercise)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", exerciseURL, bytes.NewBuffer(exerciseData))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		log.Printf("Failed to create exercise: %s\nResponse Body: %s", resp.Status, bodyString)
		return 0, errors.New("failed to create exercise: " + resp.Status)
	}

	var createExerciseResponse CreateExerciseResponse
	err = json.NewDecoder(resp.Body).Decode(&createExerciseResponse)
	if err != nil {
		return 0, err
	}

	return createExerciseResponse.ID, nil
}

func uploadImage(apiBaseURL, imagePath, token string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(imagePath))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}
	writer.Close()

	uploadURL := apiBaseURL + "/upload"
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		log.Printf("Failed to upload image: %s\nResponse Body: %s", resp.Status, bodyString)
		return "", errors.New("failed to upload image: " + resp.Status)
	}

	var uploadImageResponse struct {
		URL string `json:"url"`
	}
	err = json.NewDecoder(resp.Body).Decode(&uploadImageResponse)
	if err != nil {
		return "", err
	}

	return uploadImageResponse.URL, nil
}

func addPhotoToExercise(apiBaseURL string, exerciseID int, photo Photo, token string) error {
	photoURL := apiBaseURL + "/exercise/" + strconv.Itoa(exerciseID) + "/photos"
	photoData, err := json.Marshal(photo)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", photoURL, bytes.NewBuffer(photoData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		log.Printf("Failed to add photo to exercise: %s\nResponse Body: %s", resp.Status, bodyString)
		return errors.New("failed to add photo to exercise: " + resp.Status)
	}

	return nil
}
