package wps //WifiPositioningSystem

import (
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

func responseStringToData(strRes string) (float64, float64, float64, error) {
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

	fmt.Println()
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
	lat, lon, accRange, err := responseStringToData(string(body))
	if err != nil {
		return 0, 0, 0, err
	}

	return lat, lon, accRange, nil
}
