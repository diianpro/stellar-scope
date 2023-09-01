package apod

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/diianpro/stellar-scope/internal/domain/apod"
	"github.com/diianpro/stellar-scope/internal/domain/image"
)

type Apod struct {
	client   *http.Client
	apodAddr string
	apiKey   string
}

func NewApodClient(cfg Config, client *http.Client) *Apod {
	return &Apod{
		client:   client,
		apodAddr: cfg.Address,
		apiKey:   cfg.ApiKey,
	}
}

func (a Apod) GetMetadata(ctx context.Context, date time.Time) (*apod.Data, error) {
	values := url.Values{}
	values.Add("api_key", a.apiKey)
	values.Add("date", date.String())

	uri := fmt.Sprintf("%s?%s", a.apodAddr, values.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var apodMetaData metaData
	err = json.NewDecoder(res.Body).Decode(&apodMetaData)
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, apodMetaData.URL, nil)
	if err != nil {
		return nil, err
	}

	resImage, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resImage.Body.Close()

	img, err := io.ReadAll(resImage.Body)
	if err != nil {
		return nil, err
	}

	return &apod.Data{
		Explanation: apodMetaData.Explanation,
		Date:        apodMetaData.Date,
		Image: &image.Image{
			Title:     apodMetaData.Title,
			Extension: filepath.Ext(apodMetaData.URL),
			Data:      img,
		},
		Copyright: apodMetaData.Copyright,
	}, nil
}
