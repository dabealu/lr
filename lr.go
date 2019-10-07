package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func PrintHelp(errcode int) {
	fmt.Printf("Usage:\n" +
		"  ls-images                 list registry images (short: li)\n" +
		"  ls-tags image             list tags of an image (short: lt)\n" +
		"  rm-image image            remove all image tags (short: ri)\n" +
		"  rm-tags image:tag1,tag2   remove some image tag(s) (short: rt)\n" +
		"  help                      print help\n")
	os.Exit(errcode)
}

func GetArgs() map[string]string {
	// try to get connection info from file: ~/.lr.json
	var addr, user, password string
	type Config struct {
		Addr, User, Password string
	}
	file, _ := os.Open(os.Getenv("HOME") + "/.lr.json")
	decoder := json.NewDecoder(file)
	config := Config{}
	errJson := decoder.Decode(&config)

	if errJson == nil {
		addr = config.Addr
		user = config.User
		password = config.Password
	} else {
		// fallback to env vars if config json parsing failed
		addr = os.Getenv("REGISTRY_ADDRESS")
		user = os.Getenv("REGISTRY_USER")
		password = os.Getenv("REGISTRY_PASSWORD")
	}

	if addr == "" || user == "" || password == "" {
		fmt.Printf("Set registry connection information inside config file ~/.lr.json:\n" +
			`  {"addr":"https://registry.example.com","user":"myuser","password":"mypassword"}` +
			"\nor with env variables:\n" +
			"  export REGISTRY_ADDRESS=https://registry.example.com\n" +
			"  export REGISTRY_USER=myuser\n" +
			"  export REGISTRY_PASSWORD=mypassword\n" +
			"note that config file takes precedence over env vars\n")
		os.Exit(1)
	}

	var action string
	if len(os.Args) > 1 {
		action = os.Args[1]
	} else {
		fmt.Printf("Specify action as a first arg, i.e.: %s help\n", os.Args[0])
		os.Exit(1)
	}

	validaction := false
	for _, a := range []string{"ls-images", "li", "ls-tags", "lt", "rm-image", "ri", "rm-tags", "rt", "help"} {
		if action == a {
			validaction = true
		}
	}
	if validaction == false {
		PrintHelp(1)
	}

	var actionarg string
	if len(os.Args) > 2 {
		actionarg = os.Args[2]
	}

	var argsMap = map[string]string{
		"Action":    action,
		"ActionArg": actionarg,
		"Address":   addr,
		"User":      user,
		"Password":  password,
	}

	return argsMap
}

func GetBody(addr, user, password, action, actionarg string) []byte {
	switch action {
	case "ls-images", "li":
		addr = addr + "/v2/_catalog"
	case "ls-tags", "lt":
		addr = addr + "/v2/" + actionarg + "/tags/list"
	}

	regClient := http.Client{
		Timeout: time.Second * 15,
	}

	req, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(user, password)
	req.Header.Set("User-Agent", "registry-client")

	res, getErr := regClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	return body
}

func GetImages(body []byte) {
	type ImageList struct {
		Field []string `json:"repositories"`
	}
	getlist := ImageList{}
	jsonErr := json.Unmarshal(body, &getlist)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	for i := 0; i < len(getlist.Field); i++ {
		fmt.Println(getlist.Field[i])
	}
}

func GetTags(body []byte) []string {
	type TagList struct {
		Field []string `json:"tags"`
	}
	getlist := TagList{}
	jsonErr := json.Unmarshal(body, &getlist)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return getlist.Field
}

func RmTag(addr, user, password, action, actionarg string) {
	split := strings.Split(actionarg, ":") // actionarg ex.: image:tag1,tag2
	image := split[0]

	var tagslist []string
	switch action {
	case "rm-image", "ri":
		// get all tags of image and then delete them
		tagslist = GetTags(GetBody(addr, user, password, "ls-tags", image))
	case "rm-tags", "rt":
		tags := split[1]
		tagslist = strings.Split(tags, ",")
	}

	for _, tag := range tagslist {
		// get tag digest
		digestaddr := addr + "/v2/" + image + "/manifests/" + tag

		regClient := http.Client{
			Timeout: time.Second * 5,
		}

		req, err := http.NewRequest("GET", digestaddr, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(user, password)
		req.Header.Set("User-Agent", "registry-client")
		req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

		res, getErr := regClient.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}

		digest := res.Header.Get("Docker-Content-Digest")

		// rm tag by digest
		deleteaddr := addr + "/v2/" + image + "/manifests/" + digest

		req, err = http.NewRequest("DELETE", deleteaddr, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(user, password)
		req.Header.Set("User-Agent", "registry-client")

		res, getErr = regClient.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}

		status := res.Status
		fmt.Printf("%s\t%s:%s\t(digest: %s)\n", status, image, tag, digest)
	}
}

func main() {
	args := GetArgs()
	addr := args["Address"]
	user := args["User"]
	password := args["Password"]
	action := args["Action"]
	actionarg := args["ActionArg"]

	switch action {
	case "ls-images", "li":
		GetImages(GetBody(addr, user, password, action, actionarg))
	case "ls-tags", "lt":
		tags := GetTags(GetBody(addr, user, password, action, actionarg))
		for i := 0; i < len(tags); i++ {
			fmt.Println(tags[i])
		}
	case "rm-tags", "rt", "rm-image", "ri":
		RmTag(addr, user, password, action, actionarg)
	case "help":
		PrintHelp(0)
	}
}
