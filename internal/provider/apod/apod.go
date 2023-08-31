package apod

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/diianpro/stellar-scope/internal/domain"
	"io"
	"net/http"
	"net/url"
	"time"
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

func (a Apod) GetMetaData(ctx context.Context, date time.Time) (apod *domain.ApodData, image []byte, err error) {
	values := url.Values{}
	values.Add("api_key", a.apiKey)
	values.Add("date", date.String())

	uri := fmt.Sprintf("%s?%s", a.apodAddr, values.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := a.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	var apodMetaData metaData
	err = json.NewDecoder(res.Body).Decode(&apodMetaData)
	if err != nil {
		return nil, nil, err
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, apodMetaData.URL, nil)
	if err != nil {
		return nil, nil, err
	}

	resImage, err := a.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resImage.Body.Close()

	image, err = io.ReadAll(resImage.Body)
	if err != nil {
		return nil, nil, err
	}

	return &domain.ApodData{
		Title:       apodMetaData.Title,
		Explanation: apodMetaData.Explanation,
		Date:        apodMetaData.Date,
		Copyright:   apodMetaData.Copyright,
	}, image, nil
}
