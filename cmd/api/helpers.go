package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]any

func (app *Application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	// Encode the data to JSON, returning the error if there was one.
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}
func (app *Application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	//disallow unknown fields, doesn't ignore fields that are not meant to be in json request body
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}
	return nil
}

func (app *Application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}
func (app *Application) readIntParam(r *http.Request, intName string) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName(intName), 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("invalid %s parameter", intName)
	}
	return id, nil
}

func (app *Application) readStringParam(r *http.Request, paramName string) (string, error) {
	params := httprouter.ParamsFromContext(r.Context())
	value := params.ByName(paramName)
	if value == "" {
		return "", fmt.Errorf("missing %s parameter", paramName)
	}
	return value, nil
}

func (app *Application) GetDate(r *http.Request) (time.Time, error) {
	dateString, err := app.readStringParam(r, "date")
	if err != nil {
		return time.Now(), err
	}
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return time.Now(), errors.New("invalid date format")
	}
	return date, nil
}

func (app *Application) readBoolParam(r *http.Request, paramName string) (bool, error) {
	params := httprouter.ParamsFromContext(r.Context())
	value := params.ByName(paramName)
	if value == "" {
		return false, fmt.Errorf("missing %s parameter", paramName)
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("invalid boolean value for %s parameter", paramName)
	}
	return boolValue, nil
}

//For flexible update querys

func (app *Application) readString(qs url.Values, key string, defaultValue string) string {

	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	// Otherwise return the string.
	return s
}

func (app *Application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return defaultValue
	}
	return strings.Split(csv, ",")
}

// func (app *Application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
// 	s := qs.Get(key)
// 	if s == "" {
// 		return defaultValue
// 	}
// 	i, err := strconv.Atoi(s)
// 	if err != nil {
// 		v.AddError(key, "must be an integer value")
// 		return defaultValue
// 	}
// 	return i
// }

func (app *Application) Background(fn func()) {
	app.Wg.Add(1)
	go func() {
		defer app.Wg.Done()
		defer func() {
			if err := recover(); err != nil {
				app.Logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()
		fn()
		app.Wg.Wait()
	}()
}
