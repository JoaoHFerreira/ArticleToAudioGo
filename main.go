package main

import (
	"articletoaudio/database"
	"articletoaudio/models"
	"bytes"

	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func main() {
	database.Migrate()

	url := "https://medium.com/@arnesh07/how-golang-can-save-you-days-of-web-scraping-72f019a6de87"
	articleSliceMapArray, articleTitle, err := getMediumArticle(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	articleID, _ := database.InsertArticle(articleTitle, url)
	database.InsertParagraphs(articleSliceMapArray, articleID)

	for _, p := range articleSliceMapArray {
		paragraphIndex := p["index"].(int)
		paragraphContent := p["content"].(string)

		audioData, err := convertToAudio(paragraphContent)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Save the audio data in the AudioParagraph model
		audioParagraph := models.AudioParagraph{
			ArticleID:    articleID,
			ParagraphID:  uint(paragraphIndex),
			AudioContent: audioData,
		}
		database.GetDBInstance().Create(&audioParagraph)
	}
}

func getMediumArticle(url string) ([]map[string]interface{}, string, error) {
	urlExists(url)

	res, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, "", fmt.Errorf("failed to fetch the Medium article. Status: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, "", err
	}

	articleTitle := doc.Find("h1").First().Text()

	paragraphs := make([]map[string]interface{}, 0)
	doc.Find("p").Each(func(index int, s *goquery.Selection) {
		paragraph := s.Text()
		paragraphs = append(paragraphs, map[string]interface{}{
			"index":   index,
			"content": paragraph,
		})
	})

	return paragraphs, articleTitle, nil
}

func urlExists(url string) {
	var count int64

	db := database.GetDBInstance()
	db.Model(&models.Article{}).Where("url = ?", url).Count(&count)
	if count > 0 {
		// Raise an error and panic with the error message
		panic(fmt.Errorf("url %s already exists", url))
	}
}

import (
	"github.com/go-audio/audio"
	"github.com/go-audio/oto"
	"github.com/hajimehoshi/oto/encoding/wav"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/samples"
)

func convertToAudio(content string) ([]byte, error) {
	buf := &bytes.Buffer{}

	player, err := oto.NewPlayer(44100, 2, 2, 8192)
	if err != nil {
		return nil, err
	}
	defer player.Close()

	enc := wav.NewEncoder(buf, 44100, 16, 2, 1)
	decoder := mp3.NewDecoder(strings.NewReader(content))

	pcmBuffer := &audio.IntBuffer{}
	for {
		err := decoder.Decode(pcmBuffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		player.Write(pcmBuffer.Data)
		enc.Write(pcmBuffer)
	}

	if err := enc.Close(); err != nil {
		return nil, err
	}

	audioData := buf.Bytes()
	return audioData, nil
}
