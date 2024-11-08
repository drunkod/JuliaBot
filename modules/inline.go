package modules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/amarnathcjd/gogram/telegram"
)

func PinterestInlineHandle(i *telegram.InlineQuery) error {
	b := i.Builder()
	button := telegram.Button{}
	if i.Args() == "" {
		b.Article("No query", "Please enter a query to search for", "No query", &telegram.ArticleOptions{
			ReplyMarkup: button.Keyboard(
				button.Row(
					button.SwitchInline("Search!!!", true, "pin "),
				),
			),
		})

		i.Answer(b.Results())
		return nil
	}

	images, err := fetchPinterestImages(i.Query, 10)
	if err != nil {
		return err
	}

	if len(images) == 0 {
		b.Article("No images found", "No images found for the query", "No images found", &telegram.ArticleOptions{
			ReplyMarkup: button.Keyboard(
				button.Row(
					button.SwitchInline("Search again", true, "pin "),
				),
			),
		})

		i.Answer(b.Results())
	} else {
		var photos []telegram.Photo
		wg := sync.WaitGroup{}
		wg.Add(len(images))
		for _, image := range images {
			go func(image string) {
				defer wg.Done()
				uploaded, err := i.Client.MessagesUploadMedia("", &telegram.InputPeerSelf{}, &telegram.InputMediaPhotoExternal{
					URL: image,
				})
				if err != nil {
					fmt.Println(err)
					return
				}

				switch uploaded.(type) {
				case *telegram.MessageMediaPhoto:
					photos = append(photos, uploaded.(*telegram.MessageMediaPhoto).Photo)
				}
			}(image)
		}
		wg.Wait()

		for im := range photos {
			b.Photo(photos[im], &telegram.ArticleOptions{
				ID:    fmt.Sprintf("%d", im),
				Title: fmt.Sprintf("pinterest-image-%d", im+1),
				ReplyMarkup: button.Keyboard(
					button.Row(
						button.SwitchInline("Search again", true, "pin "),
					),
				),
			})
		}

		i.Answer(b.Results(), telegram.InlineSendOptions{
			Gallery: true,
		})
	}

	return nil
}

func fetchPinterestImages(query string, lim int) ([]string, error) {
	headers := map[string]string{
		"Accept":                  "application/json, text/javascript, */*; q=0.01",
		"Accept-Language":         "en-IN,en-GB;q=0.9,en-US;q=0.8,en;q=0.7,ml;q=0.6,bn;q=0.5",
		"Cache-Control":           "no-cache",
		"User-Agent":              "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
		"X-App-Version":           "e5cf318",
		"X-Pinterest-Appstate":    "active",
		"X-Pinterest-Pws-Handler": "www/index.js",
		"X-Pinterest-Source-Url":  "/",
		"X-Requested-With":        "XMLHttpRequest",
	}

	params := url.Values{}
	params.Set("source_url", "/search/pins/?eq=World&etslf=67&len=2&q=world%20map&rs=ac")
	params.Set("data", fmt.Sprintf(`{"options":{"query":"%s","redux_normalize_feed":true,"scope":"pins","source_url":"/search/pins/?eq=World&etslf=67&len=2&q=world%%20map&rs=ac"},"context":{}}`, query))

	baseURL := "https://in.pinterest.com/resource/BaseSearchResource/get/"
	reqURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type ImageData struct {
		Images struct {
			Orig struct {
				URL string `json:"url"`
			} `json:"orig"`
		} `json:"images"`
	}
	var parsedResponse struct {
		ResourceResponse struct {
			Data struct {
				Results []ImageData `json:"results"`
			} `json:"data"`
		} `json:"resource_response"`
	}

	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		return nil, err
	}

	var imageUrls []string
	for _, result := range parsedResponse.ResourceResponse.Data.Results {
		imageUrls = append(imageUrls, result.Images.Orig.URL)
	}

	// Limit the number of images to the limit
	if len(imageUrls) > lim {
		imageUrls = imageUrls[:lim]
	}

	return imageUrls, nil
}
