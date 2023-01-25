package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

// We will use these 3 following functions to extract and parse values from the query string (return the values passed as query parameter strings),

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

// it reads a string value from the query string and then splits it into a slice on the comma character
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

// it reads a string value from the query string and converts it to an integer before returning
func (app *application) readInt(qs url.Values, key string, defaultValue int) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("It must be an integer value")
		return defaultValue
	}

	return i
}
