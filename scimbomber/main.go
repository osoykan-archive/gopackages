package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Result struct {
	Results []Profile `json:"results"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Md5      string `json:"md5"`
}

type Profile struct {
	Name  Name   `json:"name"`
	Email string `json:"email"`
	Login Login  `json:"login"`
	Phone string `json:"phone"`
	Nat   string `json:"nat"`
}

type Name struct {
	FamilyName      string `json:"familyName"`
	GivenName       string `json:"givenName"`
	HonorificPrefix string `json:"honorificPrefix"`
}

type TokenResponse struct {
	Token string `json:"access_token"`
}

/*SCIM*/
type Email struct {
	Value string `json:"value"`
}

type PhoneNumber struct {
	Value string `json:"value"`
}
type ScimRequest struct {
	Schemas      []string      `json:"schemas"`
	UserName     string        `json:"userName"`
	Name         Name          `json:"name"`
	Active       bool          `json:"active"`
	Emails       []Email       `json:"emails"`
	PhoneNumbers []PhoneNumber `json:"phoneNumbers"`
}

var (
	scimProfileSchema = "urn:ietf:params:scim:schemas:core:2.0:Profile"
	randomUserApiUrl  = "https://randomuser.me/api"
	iamUrl            = "https://localhost"
	tokenUrl          = "https://localhost/id/connect/token"
	scimUrl           = fmt.Sprintf("%s/admin/api/scim/v2/Profiles", iamUrl)
)

func main() {

	profileCount := flag.Int("profile", 100, "Profiles count to be created")
	flag.Parse()

	if *profileCount == 0 {
		flag.PrintDefaults()
		os.Exit(0)
		return
	}
	log.Printf("INFORMATION: %d profiles will be created", *profileCount)
	start := time.Now()
	var wg sync.WaitGroup
	CreateProfiles(&wg, *profileCount)
	wg.Wait()
	elapsed := time.Since(start)
	log.Printf("Took %s", elapsed)
}

func CreateProfiles(wg *sync.WaitGroup, profileCount int) {
	processorCount := runtime.NumCPU()
	iterationCount, remainder := divMod(int64(profileCount), int64(processorCount))
	log.Printf("Processor Count: %d, Iteration Count: %d", processorCount, iterationCount)
	for i := 1; i <= processorCount; i++ {
		wg.Add(1)
		if i == processorCount {
			iterationCount += remainder
		}
		go func(iteration int64) {
			defer wg.Done()
			for j := int64(0); j < iteration; j++ {
				profile, err := getRandomProfile()
				if err != nil {
					log.Printf("Random api had problem: %s", err.Error())
				}
				err = createProfile(profile)
				if err != nil {
					log.Printf("ERROR: Create profile %s", err.Error())
				}
			}
		}(iterationCount)
	}
}

func divMod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}

func createProfile(profile Profile) error {
	request := ScimRequest{
		Schemas:  []string{scimProfileSchema},
		UserName: profile.Login.Username,
		Name: Name{
			FamilyName:      profile.Name.FamilyName,
			GivenName:       profile.Name.GivenName,
			HonorificPrefix: profile.Name.HonorificPrefix,
		},
		Active:       true,
		Emails:       []Email{{Value: profile.Email}},
		PhoneNumbers: []PhoneNumber{{Value: profile.Phone}},
	}
	requestBytes, err := json.Marshal(request)
	requestBuffer := bytes.NewBuffer(requestBytes)
	if err != nil {
		log.Println(err.Error())
	}

	accessToken, err := GetAccessToken()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	req, _ := http.NewRequest(http.MethodPost, scimUrl, requestBuffer)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if response != nil && response.StatusCode == 200 {
		log.Printf("Profile created!")
	}

	return nil
}

func GetAccessToken() (token string, err error) {
	var tokenResponse TokenResponse
	data := url.Values{}
	data.Set("client_id", "user-generator")
	data.Set("client_secret", "password")
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "iam.scim.admin")

	response, err := http.Post(tokenUrl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	if response.StatusCode != 200 {
		return "", err
	}
	_ = json.NewDecoder(response.Body).Decode(&tokenResponse)
	return tokenResponse.Token, nil
}

func getRandomProfile() (Profile, error) {
	var result Result
	client := retryablehttp.NewClient()
	client.Logger = nil
	response, err := client.Get(randomUserApiUrl)
	if err != nil {
		log.Printf("There was a problem while getting the random user")
		return Profile{}, err
	}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return Profile{}, err
	}

	return result.Results[0], nil
}
