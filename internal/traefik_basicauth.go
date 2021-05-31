package tbam

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

// Traefik Basic Auth struct
type TraefikBasicAuth struct {
	Http struct {
		Middlewares struct {
			BagetAuth struct {
				BasicAuth struct {
					Users []string `yaml:"users"`
				} `yaml:"basicAuth"`
			} `yaml:"baget-auth"`
		} `yaml:"middlewares"`
	} `yaml:"http"`
}

func RemoveCredentials(user string) error {
	t := TraefikBasicAuth{}

	yamlFile, err := ioutil.ReadFile(config.AuthFile)
	if err != nil {
		log.Debugf("yamlFile.Get err   #%v ", err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, &t)
	if err != nil {
		log.Debugf("error: %v", err)
	}
	//fmt.Printf("--- t:\n%v\n\n", t)

	exists := false
	id := 0
	for i, hashString := range t.Http.Middlewares.BagetAuth.BasicAuth.Users {
		tmpUsername := strings.Split(hashString, ":")[0]
		if tmpUsername == user {
			exists = true
			id = i
			break
		}
	}

	if exists {
		t.Http.Middlewares.BagetAuth.BasicAuth.Users[id] = t.Http.Middlewares.BagetAuth.BasicAuth.Users[len(t.Http.Middlewares.BagetAuth.BasicAuth.Users)-1]
		t.Http.Middlewares.BagetAuth.BasicAuth.Users = t.Http.Middlewares.BagetAuth.BasicAuth.Users[:len(t.Http.Middlewares.BagetAuth.BasicAuth.Users)-1]
	}

	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Debugf("error: %v", err)
		return err
	}

	//fmt.Printf("--- t dump:\n%s\n\n", string(d))

	ioutil.WriteFile(config.AuthFile, d, 0644)

	return nil
}

func UpdateCredentials(user string) error {
	t := TraefikBasicAuth{}

	yamlFile, err := ioutil.ReadFile(config.AuthFile)
	if err != nil {
		log.Debugf("yamlFile.Get err   #%v ", err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, &t)
	if err != nil {
		log.Debugf("error: %v", err)
		return err
	}
	//fmt.Printf("--- t:\n%v\n\n", t)

	notExists := true
	for i, hashString := range t.Http.Middlewares.BagetAuth.BasicAuth.Users {
		tmpUsername := strings.Split(hashString, ":")[0]
		if tmpUsername == user {
			notExists = false
			hash, err := GetValue("users", user)
			if err != nil {
				log.Debugf("error: %v", err)
				return err
			}
			t.Http.Middlewares.BagetAuth.BasicAuth.Users[i] = user + ":" + hash
		}
	}

	if notExists {
		hash, err := GetValue("users", user)
		if err != nil {
			log.Debugf("error: %v", err)
			return err
		}
		t.Http.Middlewares.BagetAuth.BasicAuth.Users = append(t.Http.Middlewares.BagetAuth.BasicAuth.Users, user+":"+hash)
	}

	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Debugf("error: %v", err)
		return err
	}

	//fmt.Printf("--- t dump:\n%s\n\n", string(d))

	ioutil.WriteFile(config.AuthFile, d, 0644)

	return nil
}
