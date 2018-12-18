package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/yeqown/gateway/config/rule"
	"github.com/yeqown/gateway/logger"
)

var (
	// ErrTimeout ...
	ErrTimeout = errors.New("combineReq timeout error")
)

type responseChan struct {
	Err   error
	Field string
	Data  map[string]interface{}
}

func combineReq(ctx context.Context, serverHost string, body io.Reader,
	comb rule.Combiner, rc chan<- responseChan) {
	var (
		err error
		r   = responseChan{
			Err:   nil,
			Field: comb.Field(),
			Data:  nil,
		}
	)

	// logger.Logger.Info("combineReq calling")

	defer func() {
		if v := recover(); v != nil {
			logger.Logger.Errorf("plugin.Proxy combineReq panic: %s", debug.Stack())
			err = v.(error)
			r.Err = err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Logger.Info("plugin.Proxy combineReq timeout!")
		r.Err = ErrTimeout
		break
	default:
		url := fmt.Sprintf("%s%s", serverHost, comb.Path())
		req, err := http.NewRequest(comb.Method(), url, body)
		if err != nil {
			r.Err = err
			logger.Logger.Errorf("could not finish NewRequest: %v", err)
			break
		}
		client := http.Client{
			Timeout: 5 * time.Second, // TODO: support config
		}
		resp, err := client.Do(req)
		if err != nil {
			r.Err = err
			logger.Logger.Errorf("could not finish client.Do: %v", err)
			break
		}

		byts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			r.Err = err
			logger.Logger.Errorf("could not ReadAll resp.Body: %v", err)
			break
		}
		resp.Body.Close()
		json.Unmarshal(byts, &r.Data)
		// logger.Logger.Info("combineReq done")
	}

	// put into channel
	rc <- r
}
