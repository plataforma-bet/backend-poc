package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

type HTTPClient func(r *http.Request) (*http.Response, error)

type Request[T Seriazable[T]] struct {
	URL     URL
	Method  string
	Headers map[string]string
	Body    T
}

type Response[T Seriazable[T]] struct {
	StatusCode int
	Headers    map[string][]string
	Body       T
	EncodeBody string
}

type ResponseFunc[Req Seriazable[Req], Res Seriazable[Res]] func(context.Context, HTTPClient) (Response[Res], error)

func SendRequest[Res Seriazable[Res], Req Seriazable[Req]](req Request[Res]) ResponseFunc[Req, Res] {
	return func(ctx context.Context, client HTTPClient) (Response[Res], error) {
		var encodedBody []byte

		if req.Body.Codec() != nil {
			var err error
			encodedBody, err = req.Body.Codec().Encode(req.Body)
			if err != nil {
				return Response[Res]{}, fmt.Errorf("error encoding body: %w", err)
			}
		}

		url, err := req.URL.Build()
		if err != nil {
			return Response[Res]{}, fmt.Errorf("error creating request url: %w", err)
		}

		request, err := http.NewRequestWithContext(ctx, req.Method, url, bytes.NewReader(encodedBody))
		if err != nil {
			return Response[Res]{}, fmt.Errorf("error creating request: %w", err)
		}

		for key, value := range req.Headers {
			request.Header.Add(key, value)
		}

		response, err := client(request)
		if err != nil {
			return Response[Res]{}, fmt.Errorf("error sending request: %w", err)
		}
		defer response.Body.Close()

		var empty Res

		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			return Response[Res]{}, fmt.Errorf("error reading response body: %w", err)
		}

		decodeBody, err := empty.Codec().Decode(responseBody)
		if err != nil {
			return Response[Res]{
				StatusCode: response.StatusCode,
				Headers:    response.Header,
				EncodeBody: string(responseBody),
			}, nil
		}

		return Response[Res]{
			StatusCode: response.StatusCode,
			Headers:    response.Header,
			Body:       decodeBody,
			EncodeBody: string(responseBody),
		}, nil
	}
}
