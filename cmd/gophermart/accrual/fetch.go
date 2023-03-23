package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/structures"
)

func (w *Worker) fetchInfo(number string) (*structures.Accrual, error) {
	response, err := w.client.Get(w.AccrualAddress + "/api/orders/" + number) //nolint:noctx
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusNoContent:
		return nil, fmt.Errorf("order %s not registered in accrual system", number) //nolint:goerr113

	case http.StatusTooManyRequests:
		return nil, errors.New("ratelimited")

	case http.StatusOK:
		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("could not read response body: %w", err)
		}

		var accrualResponse structures.Accrual

		err = json.Unmarshal(bytes, &accrualResponse)
		if err != nil {
			return nil, fmt.Errorf("could not parse response body: %s error: %w", string(bytes), err)
		}

		return &accrualResponse, nil

	default:
		return nil, fmt.Errorf("unknown status: %s", response.Status)
	}
}
