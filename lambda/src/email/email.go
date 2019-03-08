package email

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"lambda/src/jwt"
	"net/http"
	"os"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/sdeoras/api"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	ProjectName = "lambda"
	Name        = "email"
)

// Send sends email via sendgrid api
func Send(w http.ResponseWriter, r *http.Request) {
	// validate input request
	err := jwt.Manager.Validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error readhing http request body in SendEMail:%v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	sendRequest := new(api.EmailRequest)
	if err := proto.Unmarshal(b, sendRequest); err != nil {
		http.Error(w, fmt.Sprintf("could not unmarshal email send request:%v", err), http.StatusBadRequest)
		return
	}

	from := mail.NewEmail(sendRequest.FromName, sendRequest.FromEmail)
	subject := sendRequest.Subject
	to := mail.NewEmail(sendRequest.ToName, sendRequest.ToEmail)
	plainTextContent := "no plain text"
	htmlContent := string(sendRequest.Body)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not send email request:%v", err), http.StatusInternalServerError)
		return
	} else {
		sendResponse := new(api.EmailResponse)

		sendResponse.StatusCode = int64(response.StatusCode)
		sendResponse.Body = response.Body
		sendResponse.Headers = make(map[string]*api.ListOfString)

		for key, val := range response.Headers {
			key, val := key, val
			listOfString := new(api.ListOfString)
			listOfString.Value = val
			sendResponse.Headers[key] = listOfString
		}

		b, err := proto.Marshal(sendResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("error serializing email send response:%v", err), http.StatusInternalServerError)
			return
		}

		var contentType string
		var contentSize string
		contentSize = strconv.FormatInt(int64(len(b)), 10)
		if len(b) >= 512 {
			contentType = http.DetectContentType(b[:512])
		} else {
			contentType = http.DetectContentType(b)
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", contentSize)

		n, err := io.Copy(w, bytes.NewReader(b))
		if err != nil {
			http.Error(w, fmt.Sprintf("error copying bytes to http response writer:%v", err), http.StatusInternalServerError)
			return
		}

		if n != int64(len(b)) {
			http.Error(w, fmt.Sprintf("error: not all bytes transferred"), http.StatusInternalServerError)
			return
		}

		return
	}
}
