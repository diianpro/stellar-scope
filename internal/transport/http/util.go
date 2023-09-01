package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func encodeJSONResponse[T any](w http.ResponseWriter, code int, data T) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if code == http.StatusNoContent {
		return nil
	}

	return json.NewEncoder(w).Encode(data)
}

func getIntQueryParam(r *http.Request, key string) (int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return 0, fmt.Errorf("param %s missed", key)
	}

	res, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}
	return res, nil
}
