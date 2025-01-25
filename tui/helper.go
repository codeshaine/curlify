package tui

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/codeshaine/curlify/utils/request"
	"github.com/codeshaine/curlify/utils/response"
)

func (m *Model) makeRequest() {
	url, err := request.ValidateUrl(m.URLInput.Value())
	if err != nil {
		m.Result.SetContent(err.Error())
		return
	}

	if strings.EqualFold(m.MethodInput.Value(), "GET") {
		header, err := request.ParseHeaders(m.HeaderValues)
		if err != nil {
			m.Result.SetContent(err.Error())
			return
		}
		req := request.NewGet(url, header)
		res, err := req.Do()
		if err != nil {
			m.Result.SetContent(err.Error())
			return
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			m.Result.SetContent(fmt.Sprintf("Error reading response body: %v", err))
			return
		}
		//json or html (That's all we support now)
		if json.Valid(body) {
			m.Result.SetContent(response.FormatJSONResponse(body))
		} else {

			m.Result.SetContent(response.FormatHTMLResponse(body))
		}
		m.Result.GotoTop()
		return
	} else if strings.EqualFold(m.MethodInput.Value(), "POST") {
		header, err := request.ParseHeaders(m.HeaderValues)
		if err != nil {
			m.Result.SetContent(err.Error())
			return
		}
		header["Content-Type"] = []string{"application/json"}
		if !json.Valid([]byte(m.BodyValues)) {
			m.Result.SetContent("Invalid JSON body")
			return
		}

		bodyReader := strings.NewReader(m.BodyValues)

		req := request.NewPost(url, header, io.NopCloser(bodyReader))
		res, err := req.Do()
		if err != nil {
			m.Result.SetContent(err.Error())
			return
		}
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			m.Result.SetContent(fmt.Sprintf("Error reading response body: %v", err))
			return
		}

		// Format and set the response
		if json.Valid(resBody) {
			m.Result.SetContent(response.FormatJSONResponse(resBody))
		} else {
			m.Result.SetContent(response.FormatHTMLResponse(resBody))
		}
		m.Result.GotoTop()
		return
	} else {
		m.Result.SetContent("Invalid method")
		return
	}
}
