package sockets

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/salasberryfin/dockhomer/client"
)

// ImageResponse is the Docker API response to /images
type ImageResponse struct {
	ID          string   `json:"Id"`
	ParentID    string   `json:"ParentId"`
	RepoTags    []string `json:"RepoTags"`
	RepoDigests []string `json:"RepoDigests"`
	Created     int      `json:"Created"`
	Size        int      `json:"Size"`
	VirtualSize int      `json:"VirtualSize"`
	SharedSize  int      `json:"SharedSize"`
	//Labels map `json:"Labels"`
	Containers int `json:"Containers"`
}

func parse(r []byte) {
	var unmarsh []ImageResponse
	if err := json.Unmarshal(r, &unmarsh); err != nil {
		log.Fatal("Something went wrong parsing the response: ", err)
	}
	log.Println(unmarsh[0])
	log.Println(unmarsh[1])
}

// Get sends a "GET HTTP" request to the Docker unix socket
func Get(c *client.APIClient) {
	// HTTP over unix socket
	resp, err := c.Client.Get("http://foo/images/json")
	if err != nil {
		log.Fatal("Failed to send HTTP request:", err)
	}
	log.Println("HTTP response:", resp.Body)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read response:", err)
	}

	parse(body)
}
