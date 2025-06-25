package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

func ValidateJSONBody(r *http.Request, v interface{}, log *zap.Logger, op string) error {
	if r.Body == nil || r.Body == http.NoBody {
		log.Error("body is nil", zap.String("op", op))
		return fmt.Errorf("body is nil")
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("failed to read request body", zap.String("op", op), zap.Error(err))
		return fmt.Errorf("invalid request body")
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		var detailedErr error
		switch e := err.(type) {
		case *json.UnmarshalTypeError:
			detailedErr = fmt.Errorf("field %s must be like %s", e.Field, e.Type)
		case *json.SyntaxError:
			detailedErr = fmt.Errorf("json syntax error")
		default:
			if strings.Contains(err.Error(), "unknown field") {
				detailedErr = fmt.Errorf("unknown field JSON")
			} else {
				detailedErr = fmt.Errorf("JSON parsing error: %v", err)
			}
		}
		log.Error("invalid JSON",
			zap.String("op", op),
			zap.Error(detailedErr))
		return detailedErr
	}
	return nil
}
