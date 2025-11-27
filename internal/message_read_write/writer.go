package messagereadwrite

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

func createRequestJson(data string, req_type string) string {
	return fmt.Sprintf("{\"request\": \"%s\", \"status\": \"200\", \"data\": %s}", req_type, data)
}

func createErrorRequestJson(data string, req_type string) string {
	return fmt.Sprintf("{\"request\": \"%s\", \"status\": \"406\", \"error\": %s}", req_type, data)
}

func (rw *ReadWriter) Write(req_type string, data string) error {
	if rw.conn == nil {
		log.Printf("conn is nill could not write message %s\n", req_type)
		return errors.New("Couldn't access the websocket connection")
	}
	log.Printf("\t%s: %s\n", req_type, data)

	req_json := createRequestJson("\""+data+"\"", req_type)
	rw.writer_jobs <- req_json
	return nil
}

func (rw *ReadWriter) WriteError(req_type string, err string) error {
	if rw.conn == nil {
		log.Printf("conn is nill could not write message %s\n", req_type)
		return errors.New("Couldn't access the websocket connection")
	}
	log.Printf("\terr[%s]: %s\n", req_type, err)

	req_json := createErrorRequestJson("\""+err+"\"", req_type)
	rw.writer_jobs <- req_json
	return nil
}

func (rw *ReadWriter) WriteAny(req_type string, data any) error {
	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if rw.conn == nil {
		log.Printf("conn is nill could not write message %s\n", req_type)
		return errors.New("Couldn't access the websocket connection")
	}

	log.Printf("\t%s: %s\n", req_type, string(dataJson))

	req_json := createRequestJson(string(dataJson), req_type)
	rw.writer_jobs <- req_json
	return err
}
