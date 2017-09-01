package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gogate/gate/types"
	"io/ioutil"
	"net/http"
	"strings"
)

// find key value for types.State
func FindKey(body types.State, key string, contents types.Content) string {
	if key == "ip" {
		var ip string
		var last_ipv4 string
		for _, v := range body.Interfaces {
			// if there is already a same key in content
			for _, _ip := range v.Ips {
				_, exists := contents[_ip]
				if exists {
					return _ip
				}
				if !strings.Contains(_ip, "::") {
					last_ipv4 = _ip
				}
				ip = _ip
			}
		}
		if last_ipv4 != "" {
			return last_ipv4
		}
		return ip
	}

	if key == "mac" {
		var mac string
		for _, v := range body.Interfaces {
			// if there is already a same key in content
			_, exists := contents[v.Mac]
			if exists {
				return v.Mac
			}
			mac = v.Mac
		}
		return mac
	}

	//hostname as default
	return body.Hostname
}

func ForwardToMaster(master string, data types.State) (string, error) {
	content, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	res, err := http.Post(master, "text/plain", bytes.NewReader(content))
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	if res.StatusCode != 201 {
		return "", errors.New(fmt.Sprintf("server returned %d, not 201/Created", res.StatusCode))
	}
	if body, err := ioutil.ReadAll(res.Body); err != nil {
		return "", err
	} else {
		return string(body), nil
	}
}