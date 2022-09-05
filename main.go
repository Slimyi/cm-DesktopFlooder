package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func openFile(fname string) {
	// Gives the directory of the file with the filename to the command prompt to open it
	open := exec.Command("cmd.exe", "/c", fname)
	err := open.Run()
	// Errors are mostly because of faulty filetypes, they are ignored and logged.
	if err != nil {
		fmt.Println(err)
		fmt.Println("error with file" + fname)
	}
}

type catalog []struct {
	Page    int `json:"page"`
	Threads []struct {
		No            int    `json:"no"`
		Sticky        int    `json:"sticky,omitempty"`
		Closed        int    `json:"closed,omitempty"`
		Now           string `json:"now"`
		Name          string `json:"name"`
		Sub           string `json:"sub,omitempty"`
		Com           string `json:"com"`
		Filename      string `json:"filename"`
		Ext           string `json:"ext"`
		W             int    `json:"w"`
		H             int    `json:"h"`
		TnW           int    `json:"tn_w"`
		TnH           int    `json:"tn_h"`
		Tim           int64  `json:"tim"`
		Time          int    `json:"time"`
		Md5           string `json:"md5"`
		Fsize         int    `json:"fsize"`
		Resto         int    `json:"resto"`
		Capcode       string `json:"capcode,omitempty"`
		SemanticURL   string `json:"semantic_url"`
		Replies       int    `json:"replies"`
		Images        int    `json:"images"`
		OmittedPosts  int    `json:"omitted_posts"`
		OmittedImages int    `json:"omitted_images"`
		LastReplies   []struct {
			No       int    `json:"no"`
			Now      string `json:"now"`
			Name     string `json:"name"`
			Com      string `json:"com"`
			Filename string `json:"filename"`
			Ext      string `json:"ext"`
			W        int    `json:"w"`
			H        int    `json:"h"`
			TnW      int    `json:"tn_w"`
			TnH      int    `json:"tn_h"`
			Tim      int64  `json:"tim"`
			Time     int    `json:"time"`
			Md5      string `json:"md5"`
			Fsize    int    `json:"fsize"`
			Resto    int    `json:"resto"`
			Capcode  string `json:"capcode"`
		} `json:"last_replies"`
		LastModified int `json:"last_modified"`
		Bumplimit    int `json:"bumplimit,omitempty"`
		Imagelimit   int `json:"imagelimit,omitempty"`
	} `json:"threads"`
}

func main() {
	var stdout bytes.Buffer
	userCom := exec.Command("cmd.exe", "/c", "echo", "%username%")
	userCom.Stdout = &stdout
	err := userCom.Run()
	if err != nil {
		panic(err)
	}
	outStr := string(stdout.Bytes()[:len(stdout.Bytes())-2])
	// This is the directory of the user's desktop
	dir := "/Users/" + outStr + "/Desktop/"

	resp, err := http.Get("https://a.4cdn.org/cm/catalog.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.StatusCode)
	var cat catalog
	json.NewDecoder(resp.Body).Decode(&cat)
	for i := 0; i < len(cat); i++ {
		for j := 0; j < len(cat[i].Threads); j++ {
			respCm, err := http.Get("https://i.4cdn.org/cm/" + fmt.Sprintf("%d", cat[i].Threads[j].Tim) + cat[i].Threads[j].Ext)
			if err != nil {
				panic(err)
			}
			fmt.Println(cat[i].Threads[j].Ext)
			f, _ := os.Create(dir + fmt.Sprintf("%d", cat[i].Threads[j].Tim) + cat[i].Threads[j].Ext)
			io.Copy(f, respCm.Body)
			respCm.Body.Close()
			f.Close()
			openFile(dir + fmt.Sprintf("%d", cat[i].Threads[j].Tim) + cat[i].Threads[j].Ext)
		}
	}
	fmt.Println(resp.StatusCode)
	resp.Body.Close()
}
