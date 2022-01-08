package vision

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type BatchAnnotateRequest struct {
	Requests []AnnotateRequest `json:"requests"`
}

type AnnotateRequest struct {
	Image    Image     `json:"image"`
	Features []Feature `json:"features"`
}

type Image struct {
	Content string `json:"content"`
}

type Feature struct {
	FeatureType string `json:"type"`
	Max         int    `json:"maxResults"`
}

type BatchAnnotateResponse struct {
	Responses []AnnotateResponse `json:"responses"`
}

type AnnotateResponse struct {
	Result WebDetection `json:"webDetection"`
}

type WebDetection struct {
	BestGuessLabels []BestGuessLabel `json:"bestGuessLabels"`
	WebEntities     []WebEntity      `json:"webEntities"`
}

type BestGuessLabel struct {
	Label        string `json:"label"`
	LanguageCode string `json:"languageCode"`
}

type WebEntity struct {
	EntityID    string  `json:"entityId"`
	Score       float32 `json:"score"`
	Description string  `json"description"`
}

// Search takes base64 encoded image data and returns the first web detection result
func Search(imageData string) (string, error) {

	req := BatchAnnotateRequest{
		Requests: []AnnotateRequest{
			{
				Image: Image{
					Content: imageData,
				},
				Features: []Feature{
					{
						FeatureType: "WEB_DETECTION",
						Max:         1,
					},
				},
			},
		},
	}
	// fmt.Println(req)

	reqBuf, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	apiKey := os.Getenv("AR_API_KEY")
	apiURL := fmt.Sprintf("https://vision.googleapis.com/v1/images:annotate?key=%v", apiKey)

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(reqBuf))
	// resp, err := http.Post("http://localhost:8000", "application/json", bytes.NewBuffer(reqBuf))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response BatchAnnotateResponse

	err = json.Unmarshal(respBytes, &response)
	if err != nil {
		return "", err
	}

	if len(response.Responses) == 0 {
		return "", errors.New("no results")
	}

	if len(response.Responses[0].Result.BestGuessLabels) == 0 {
		return "", errors.New("no guesses")
	}
	return response.Responses[0].Result.BestGuessLabels[0].Label, nil
}
