package wps //WifiPositioningSystem

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WifiInfo struct {
	Mac  string
	Rssi int
}

func GetPositionByOpenData(ArrrayWifiInfo []WifiInfo) {
	if len(ArrrayWifiInfo) <= 0 {
		fmt.Printf("Arrarylen=%d ,not support\n", len(ArrrayWifiInfo))
		return
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
	fmt.Println(strtmp)
	response, err := http.Get(strtmp)
	if err != nil {
		// handle error
	}
	//程序在使用完回复后必须关闭回复的主体。
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
}
