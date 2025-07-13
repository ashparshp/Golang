package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
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

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, values := range headers[0] {
			w.Header()[key] = values
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
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

// invalidCredentials sends a JSON response with an error message for unauthorized access
func (app *application) invalidCredentials(w http.ResponseWriter) error {
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = true
	payload.Message = "Invalid credentials"

	err := app.writeJSON(w, http.StatusUnauthorized, payload)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	return nil
}

// passwordMatch checks if the provided password matches the stored hash
func (app *application) passwordMatch(storedHash, providedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(providedPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			app.errorLog.Println("Error comparing password:", err)
			return false, err
		}
	}
	return true, nil
}
