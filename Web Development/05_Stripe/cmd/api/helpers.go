package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1_048_576 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	err := dec.Decode(data)
	if err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must contain only a single JSON value")
	}
	return nil
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) error {
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = true
	payload.Message = err.Error()
	out, err := json.MarshalIndent(payload, "", "   ")
	if err != nil {
		app.errorLog.Println(err)
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(out)
	if err != nil {
		app.errorLog.Println(err)
	}
	return nil
}
