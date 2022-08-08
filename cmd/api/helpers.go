package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
)

type envelope map[string]interface{}

// This is used as a helper to read ID parameters
func readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64) // the parameter passed to the url will be integer
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// Converting data to JSON format
func (app *application) writeJSON(w http.ResponseWriter, data envelope) error {
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	jsonData = append(jsonData, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	return nil
}

// Reading JSON data
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&data) // it reads the json
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character &d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("incorrect JSON type for field &q\n", unmarshalTypeError.Field)
			}
			return fmt.Errorf("incorrect JSON type for type &q\n", unmarshalTypeError.Field)

		case errors.Is(err, io.EOF):
			return errors.New("body empty")

		case errors.As(err, invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	return nil
}
