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
	OriginalUri      string   `json:"originalUri"`
	Name             string   `json:"name"`
	Muscle           string   `json:"muscle"`
	AdditionalMuscle string   `json:"additionalMuscle"`
	Type             string   `json:"type"`
	Equipment        string   `json:"equipment"`
	Difficulty       string   `json:"difficulty"`
	Photos           []string `json:"photos"`
}

type Course struct {
	ID                int     `json:"id"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	Difficulty        string  `json:"difficulty"`
	DifficultyNumeric int     `json:"difficulty_numeric"`
	Direction         string  `json:"direction"`
	TrainerID         int     `json:"trainer_id"`
	Cost              float64 `json:"cost"`
	ParticipantsCount int     `json:"participants_count"`
	Rating            float64 `json:"rating"`
	RequiredTools     string  `json:"required_tools"`
}

type Class struct {
	ID          int    `json:"id"`
	CourseID    int    `json:"course_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
	Content     string `json:"content"`
	Video       string `json:"video"`
	Tips        string `json:"tips"`
}

type ClassImage struct {
	ID      int    `json:"id"`
	ClassID int    `json:"class_id"`
	Image   string `json:"image"`
}

type CreateCourseResponse struct {
	Course Course `json:"course"`
}

type CreateClassResponse struct {
	Class Class `json:"class"`
}

type CreateClassImageResponse struct {
	ClassImage ClassImage `json:"class_image"`
}

type UploadImageResponse struct {
	URL string `json:"url"`
}

func ImportCoursesFromJSON(filePath string, apiBaseURL string) error {
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
		newCourse := Course{
			Title:             exercise.Name,
			Description:       "Muscle: " + exercise.Muscle + ", Additional Muscle: " + exercise.AdditionalMuscle + ", Type: " + exercise.Type + ", Equipment: " + exercise.Equipment,
			Difficulty:        exercise.Difficulty,
			DifficultyNumeric: 1,         // Задайте соответствующее значение
			Direction:         "General", // Задайте соответствующее значение
			TrainerID:         1,         // Задайте соответствующее значение
			Cost:              0,         // Задайте соответствующее значение
			ParticipantsCount: 0,         // Задайте соответствующее значение
			Rating:            0,         // Задайте соответствующее значение
			RequiredTools:     exercise.Equipment,
		}

		courseID, err := createCourse(apiBaseURL, newCourse, token)
		if err != nil {
			return err
		}

		newClass := Class{
			CourseID:    courseID,
			Title:       exercise.Name,
			Description: exercise.Muscle,
		}

		classID, err := createClass(apiBaseURL, newClass, token)
		if err != nil {
			return err
		}

		for _, photo := range exercise.Photos {
			photoPath := filepath.Join(imagesFolder, photo)
			uploadedImageURL, err := uploadImage(apiBaseURL, photoPath, token)
			if err != nil {
				return err
			}

			newClassImage := ClassImage{
				ClassID: classID,
				Image:   uploadedImageURL,
			}

			log.Printf("Creating class image with data: %+v\n", newClassImage)

			err = createClassImage(apiBaseURL, newClassImage, token)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func createCourse(apiBaseURL string, newCourse Course, token string) (int, error) {
	courseURL := apiBaseURL + "/course"
	courseData, err := json.Marshal(newCourse)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", courseURL, bytes.NewBuffer(courseData))
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
		log.Printf("Failed to create course: %s\nResponse Body: %s", resp.Status, bodyString)
		return 0, errors.New("failed to create course: " + resp.Status)
	}

	var createCourseResponse CreateCourseResponse
	err = json.NewDecoder(resp.Body).Decode(&createCourseResponse)
	if err != nil {
		return 0, err
	}

	return createCourseResponse.Course.ID, nil
}

func createClass(apiBaseURL string, newClass Class, token string) (int, error) {
	classURL := apiBaseURL + "/course/" + strconv.Itoa(newClass.CourseID) + "/classes"
	classData, err := json.Marshal(newClass)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", classURL, bytes.NewBuffer(classData))
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
		log.Printf("Failed to create class: %s\nResponse Body: %s", resp.Status, bodyString)
		return 0, errors.New("failed to create class: " + resp.Status)
	}

	var createClassResponse CreateClassResponse
	err = json.NewDecoder(resp.Body).Decode(&createClassResponse)
	if err != nil {
		return 0, err
	}

	return createClassResponse.Class.ID, nil
}

func createClassImage(apiBaseURL string, newClassImage ClassImage, token string) error {
	classImageURL := apiBaseURL + "/course/" + strconv.Itoa(newClassImage.ClassID) + "/classes/" + strconv.Itoa(newClassImage.ClassID) + "/images"
	classImageData, err := json.Marshal(newClassImage)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", classImageURL, bytes.NewBuffer(classImageData))
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
		log.Printf("Failed to create class image: %s\nResponse Body: %s", resp.Status, bodyString)
		return errors.New("failed to create class image: " + resp.Status)
	}

	return nil
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

	var uploadImageResponse UploadImageResponse
	err = json.NewDecoder(resp.Body).Decode(&uploadImageResponse)
	if err != nil {
		return "", err
	}

	return uploadImageResponse.URL, nil
}
