package basic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/stack-labs/stack-rpc-plugins/service/stackweb/plugins/basic/tools"
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/util/errors"
)

func (a *api) rpcCall(w http.ResponseWriter, ctx context.Context, rpcReq *rpcRequest) {
	if len(rpcReq.Service) == 0 {
		tools.WriteError(w, fmt.Errorf("Service Is Not found "))
		return
	}

	if len(rpcReq.Endpoint) == 0 {
		tools.WriteError(w, fmt.Errorf("Endpoint Is Not found err "))
		return
	}

	// decode rpc request param body
	if req, ok := rpcReq.Request.(string); ok {
		d := json.NewDecoder(strings.NewReader(req))
		d.UseNumber()

		if err := d.Decode(&rpcReq.Request); err != nil {
			tools.WriteError(w, fmt.Errorf("error decoding request string err: %s", err))
			return
		}
	}

	// create request/response
	var response json.RawMessage
	var err error
	req := a.rpcClient.NewRequest(rpcReq.Service, rpcReq.Endpoint, rpcReq.Request, client.WithContentType("application/json"))

	var opts []client.CallOption

	// set timeout
	if rpcReq.timeout > 0 {
		opts = append(opts, client.WithRequestTimeout(time.Duration(rpcReq.timeout)*time.Second))
	}

	// remote call
	if len(rpcReq.Address) > 0 {
		opts = append(opts, client.WithAddress(rpcReq.Address))
	}

	// remote call
	err = a.rpcClient.Call(ctx, req, &response, opts...)
	if err != nil {
		ce := errors.Parse(err.Error())
		switch ce.Code {
		case 0:
			// assuming it's totally screwed
			ce.Code = 500
			ce.Id = "stack.rpc"
			ce.Status = http.StatusText(500)
			ce.Detail = "error during request: " + ce.Detail
			w.WriteHeader(500)
		default:
			w.WriteHeader(int(ce.Code))
		}
		w.Write([]byte(ce.Error()))
		return
	}

	if strings.Contains(rpcReq.URL, "/v1/rpc") {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		tools.WriteJsonData(w, response)
	}
}
