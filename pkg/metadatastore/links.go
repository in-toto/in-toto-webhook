package metadatastore

import (
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	// utility that should go away
	"net/http/httputil"
)

func HandleLinkSerialization(w http.ResponseWriter,
	r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}
	dump, err := httputil.DumpRequest(r, true)

	if err != nil {
		fmt.Printf("proper error handling is in order here")
	}
	// Should handle by returning a 200 ok or so
	fmt.Fprintf(w, "%q\n", dump)
	//fmt.Printf("%q\n", dump)

	linkName, namespace := parseURI(r.URL)
	body, _ := ioutil.ReadAll(r.Body)
	serializeLink(linkName, namespace, body)
}

func parseURI (location *url.URL) (string, string) {

	linkName := ""
	namespace := ""

	requesturi := location.RequestURI()

	fields := strings.Split(requesturi, "/")
	if len(fields) != 4 {
		fmt.Printf("proper error handling is in order here%v", fields)
	}
	linkName = fields[3]
	namespace = fields[2]

	return linkName, namespace

}

func serializeLink(link string, namespace string, data []byte) {

	err := os.MkdirAll(namespace, 0755)
	if err != nil {
		fmt.Printf("couldn't create namespace dir!")
	}

	ioutil.WriteFile(fmt.Sprintf("%s/%s", namespace, link), data, 0644)

}
