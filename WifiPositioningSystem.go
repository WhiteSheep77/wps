package wps //WifiPositioningSystem

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WifiInfo struct {
	Mac  string
	Rssi int
}

func wifiInfoToStringForOpenData(ArrrayWifiInfo []WifiInfo) (string, error) {
	if len(ArrrayWifiInfo) <= 0 {
		return "", errors.New("wifiInfoToStringForOpenData Error:WIfi Arrarylen is zero")
	}

	var strtmp string
	if 1 == len(ArrrayWifiInfo) {
		strtmp = "http://api.mylnikov.org/geolocation/wifi?v=1.2&bssid="
		strtmp = strtmp + ArrrayWifiInfo[0].Mac + "," + fmt.Sprintf("%d", ArrrayWifiInfo[0].Rssi) + ";"
	} else {
		strtmp = "http://api.mylnikov.org/geolocation/wifi?v=1.2&search="
		var strMacTmp string
		for i := 0; i < len(ArrrayWifiInfo); i++ {
			strMacTmp = strMacTmp + ArrrayWifiInfo[i].Mac + "," + fmt.Sprintf("%d", ArrrayWifiInfo[i].Rssi) + ";"
		}
		strMacBase64Tmp := base64.StdEncoding.EncodeToString([]byte(strMacTmp))
		strtmp = strtmp + strMacBase64Tmp
	}
	return strtmp, nil
}

func responseStringToDataForOpenData(strRes string) (float64, float64, float64, error) {
	type StructData struct {
		Lat   float64
		Range float64
		Lon   float64
		Time  int
	}

	type StructHttpRes struct {
		Result int
		Data   StructData
	}

	tmphttpRes := StructHttpRes{}
	json.Unmarshal([]byte(strRes), &tmphttpRes)

	if 200 != tmphttpRes.Result {
		return 0, 0, 0, errors.New("responseStringToData Error:result not 200")
	}

	return tmphttpRes.Data.Lat, tmphttpRes.Data.Lon, tmphttpRes.Data.Range, nil
}

func GetPositionByOpenData(ArrrayWifiInfo []WifiInfo) (float64, float64, float64, error) {
	strtmp, err := wifiInfoToStringForOpenData(ArrrayWifiInfo)
	if err != nil {
		return 0, 0, 0, err
	}

	response, err := http.Get(strtmp)
	if err != nil {
		return 0, 0, 0, err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	lat, lon, accRange, err := responseStringToDataForOpenData(string(body))
	if err != nil {
		return 0, 0, 0, err
	}

	return lat, lon, accRange, nil
}

func wifiInfoToJasonForGeolocation(ArrrayWifiInfo []WifiInfo) (string, error) {
	type StructWifiPoint struct {
		MacAddress     string `json:"macAddress"`
		SignalStrength int    `json:"signalStrength"`
	}
	type StructJasonData struct {
		ConsiderIp       string            `json:"considerIp"`
		WifiAccessPoints []StructWifiPoint `json:"wifiAccessPoints"`
	}

	if len(ArrrayWifiInfo) <= 0 {
		return "", errors.New("wifiInfoToJasonForGeolocation Error:WIfi Arrarylen is zero")
	}

	tmpWifiAccessPoints := make([]StructWifiPoint, len(ArrrayWifiInfo))
	for i := 0; i < len(ArrrayWifiInfo); i++ {
		tmpWifiAccessPoints[i].MacAddress = ArrrayWifiInfo[i].Mac
		tmpWifiAccessPoints[i].SignalStrength = ArrrayWifiInfo[i].Rssi
	}
	tmpMarshal := StructJasonData{"false", tmpWifiAccessPoints}
	tmpByte, err := json.Marshal(tmpMarshal)
	if err != nil {
		return "", err
	}

	return string(tmpByte), nil
}

func responseStringToDataForGeolocation(strRes string) (float64, float64, float64, error) {
	type StructLocation struct {
		Lat float64
		Lng float64
	}

	type StructData struct {
		Location StructLocation
		Accuracy float64
	}

	tmphttpRes := StructData{}
	json.Unmarshal([]byte(strRes), &tmphttpRes)

	if 0 == tmphttpRes.Location.Lat {
		type StructDataError struct {
			Errors  string
			Code    int `json:"code"`
			Message string
		}
		type StructError struct {
			Error StructDataError
		}
		tmpError := StructError{}
		json.Unmarshal([]byte(strRes), &tmpError)
		var ErrString = "code:" + fmt.Sprintf("%d", tmpError.Error.Code)
		return 0, 0, 0, errors.New(ErrString)
	}
	return tmphttpRes.Location.Lat, tmphttpRes.Location.Lng, tmphttpRes.Accuracy, nil
}

func GetPositionByGeolocation(ArrrayWifiInfo []WifiInfo) (float64, float64, float64, error) {

	wifiJason, err := wifiInfoToJasonForGeolocation(ArrrayWifiInfo)
	if err != nil {
		return 0, 0, 0, err
	}

	req := bytes.NewBuffer([]byte(wifiJason))
	body_type := "application/json;charset=utf-8"
	response, err := http.Post(uRlGeo+aPIGeo, body_type, req)
	if err != nil {
		return 0, 0, 0, err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	lat, lon, accRange, err := responseStringToDataForGeolocation(string(body))
	if err != nil {
		return 0, 0, 0, err
	}

	return lat, lon, accRange, nil
}
