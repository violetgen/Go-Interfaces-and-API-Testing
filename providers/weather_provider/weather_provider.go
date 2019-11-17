package weather_provider

import (
	"encoding/json"
	"fmt"
	"interface-testing/clients/restclient"
	"interface-testing/domain/weather_domain"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	weatherUrl = "https://api.darksky.net/forecast/%s/%v,%v"
)

func GetWeather(request weather_domain.WeatherRequest) (*weather_domain.Weather, *weather_domain.WeatherError) {
	url := fmt.Sprintf(weatherUrl, request.ApiKey, request.Latitude, request.Longitude)
	response, err := restclient.Get(url)
	if err != nil {
		log.Println(fmt.Sprintf("error when trying to get weather from dark sky api %s", err.Error()))
		return nil, &weather_domain.WeatherError{
			Code: http.StatusBadRequest,
			Error:    err.Error(),
		}
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, &weather_domain.WeatherError{
			Code: http.StatusBadRequest,
			Error:    err.Error(),
		}
	}
	defer response.Body.Close()

	//The api owner can decide to change datatypes, etc. When this happen, it might affect the error format returned
	if response.StatusCode > 299 {
		var errResponse weather_domain.WeatherError
		if err := json.Unmarshal(bytes, &errResponse); err != nil {
			return nil, &weather_domain.WeatherError{
				Code: http.StatusInternalServerError,
				Error: "invalid json response body",
			}
		}
		errResponse.Code = response.StatusCode
		return nil, &errResponse
	}
	var result weather_domain.Weather
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf("error when trying to unmarshal weather  successful response: %s", err.Error()))
		return nil, &weather_domain.WeatherError{Code: http.StatusInternalServerError, Error: "error unmarshaling weather fetch response"}
	}
	return &result, nil
}


