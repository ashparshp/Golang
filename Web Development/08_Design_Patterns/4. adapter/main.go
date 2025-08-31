package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

// ToDo is the type for the data we are working with.
type ToDo struct {
	UserID    int    `json:"userId" xml:"userId"`
	ID        int    `json:"id" xml:"id"`
	Title     string `json:"title" xml:"title"`
	Completed bool   `json:"completed" xml:"completed"`
}

// DataInterface is simply our target interface, which defines all the methods that
// any type which implements this interface must have. In our case, we only have
// one, but you might have one for each of Create, Read, Update, Delete, and more.
type DataInterface interface {
	GetData() (*ToDo, error)
}

// RemoteService is the Adaptor type. It embeds a DataInterface interface
// (which is critical to the pattern). It is a simple wrapper for this interface.
type RemoteService struct {
	Remote DataInterface
}

// CallRemoteService is the function on RemoteService which lets us
// call any adaptor which implements the DataInterface type.
func (rs *RemoteService) CallRemoteService() (*ToDo, error) {
	return rs.Remote.GetData()
}

// JSONBackend is the JSON adaptee, which needs to satisfy the DataInterface by
// have a GetData method.
type JSONBackend struct{}

// GetData is necessary so that JSONBackend satisifies the DataInterface requirements.
func (jb *JSONBackend) GetData() (*ToDo, error) {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var todo ToDo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// XMLBackend is the XML adaptee, which needs to satisfy the DataInterface by
// have a GetData method.
type XMLBackend struct{}

// GetData is necessary so that JSONBackend satisifies the DataInterface requirements.
func (xb *XMLBackend) GetData() (*ToDo, error) {
	xmlFile := `
<?xml version="1.0" encoding="UTF-8" ?>
<root>
	<userId>1</userId>
	<id>1</id>
	<title>delectus aut autem</title>
	<completed>false</completed>
</root>
`

	var todo ToDo
	err := xml.Unmarshal([]byte(xmlFile), &todo)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func main() {
	// No adapter
	todo := getRemoteData()
	fmt.Println("TODO without adapter:\t", todo.ID, todo.Title)

	// With adapter, using JSON
	jsonBackend := &JSONBackend{}
	jsonAdapter := &RemoteService{Remote: jsonBackend}
	tdFromJSON, _ := jsonAdapter.CallRemoteService()
	fmt.Println("From JSON Adapter:\t", tdFromJSON.ID, tdFromJSON.Title)

	// With adapter, using XML
	xmlBackend := &XMLBackend{}
	xmlAdapter := &RemoteService{Remote: xmlBackend}
	tdFromXML, err := xmlAdapter.CallRemoteService()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("From XML Adapter:\t", tdFromXML.ID, tdFromXML.Title)
}

func getRemoteData() *ToDo {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var todo ToDo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		log.Fatalln(err)
	}

	return &todo
}
