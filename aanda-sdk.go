package aandaSdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const apiUrl = "http://api.aanda.ru/xml_gateway/"

type Api struct {
	BuyerId  string
	UserId   string
	Password string
	Language string
}

func sendReq(data url.Values) []byte {
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBufferString(data.Encode()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	return body

}
func parseError(body []byte) error {

	var answ string
	err := json.Unmarshal(body, &answ)
	if err == nil {
		if strings.Index(answ, "Ошибка авторазации") == -1 {
			return errors.New("Ошибка авторазиции")
		}
	}

	re := regexp.MustCompile(`Note="(.*)"`)
	res := re.FindStringSubmatch(string(body))
	if len(res) != 0 {
		return errors.New(res[1])
	}

	return nil
}

func NewApi(auth Auth) *Api {
	return &Api{
		BuyerId:  auth.BuyerId,
		UserId:   auth.UserId,
		Password: auth.Password,
		Language: auth.Language,
	}
}
func (self *Api) CountryListRequest() ([]CountryListAnswer, error) {
	data := url.Values{}
	data.Set("RequestType", "json")
	data.Add("RequestName", "CountryListRequest")
	data.Add("CompanyId", self.BuyerId)
	data.Add("UserId", self.UserId)
	data.Add("Password", self.Password)
	data.Add("Language", self.Language)

	body := sendReq(data)
	jsonData := []CountryListAnswer{}
	err := json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, parseError(body)
	}

	return jsonData, nil
}
func (self *Api) HotelSearchRequest(searchReq HotelSearchRequest) ([]HotelSearchAnswer, error) {
	searchReq.BuyerId = self.BuyerId
	searchReq.UserId = self.UserId
	searchReq.Password = self.Password
	searchReq.Language = self.Language

	jsonReq, err := json.Marshal(searchReq)
	if err != nil {
		panic(err)
	}

	data := url.Values{}
	data.Set("RequestType", "json")
	data.Add("RequestName", "HotelSearchRequest")
	data.Add("JSON", string(jsonReq))

	body := sendReq(data)
	jsonData := []HotelSearchAnswer{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, parseError(body)
	}

	return jsonData, nil
}
