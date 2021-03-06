package wps //WifiPositioningSystem

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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
		fmt.Println("GetPositionByOpenData String err")
		return 0, 0, 0, err
	}

	fmt.Println("OpenData str:")
	fmt.Println(strtmp)

	response, err := http.Get(strtmp)
	if err != nil {
		fmt.Println("GetPositionByOpenData http err")
		return 0, 0, 0, err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	lat, lon, accRange, err := responseStringToDataForOpenData(string(body))
	if err != nil {
		fmt.Println("GetPositionByOpenData response err:")
		fmt.Println(string(body))
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

	fmt.Println(strRes)

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

func GetPositionByGeolocationWPS(ArrrayWifiInfo []WifiInfo) (float64, float64, float64, error) {

	wifiJason, err := wifiInfoToJasonForGeolocation(ArrrayWifiInfo)
	if err != nil {
		return 0, 0, 0, err
	}

	req := bytes.NewBuffer([]byte(wifiJason))
	body_type := "application/json;charset=utf-8"
	response, err := http.Post(uRlGeo+GetGeoAPI(), body_type, req)
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

func GetPositionBybaidu(ArrrayWifiInfo []WifiInfo) (float64, float64, float64, error) {
	return 0, 0, 0, nil
}

func LBSInfoToJasonForGeolocation(ArrayCellID [3]string, ArrayLAC [3]string, ArrayMNC [3]string, ArrayMCC [3]string) (string, error) {
	type StructCellTowers struct {
		CellId            int64 `json:"cellId"`
		LocationAreaCode  int64 `json:"locationAreaCode"`
		MobileCountryCode int64 `json:"mobileCountryCode"`
		MobileNetworkCode int64 `json:"mobileNetworkCode"`
	}
	type StructJasonData struct {
		ConsiderIp string             `json:"considerIp"`
		CellTowers []StructCellTowers `json:"cellTowers"`
	}

	tmpStructCellTowers := make([]StructCellTowers, 1)

	tmpStructCellTowers[0].CellId, _ = strconv.ParseInt(ArrayCellID[0], 10, 64)
	tmpStructCellTowers[0].LocationAreaCode, _ = strconv.ParseInt(ArrayLAC[0], 10, 64)
	tmpStructCellTowers[0].MobileCountryCode, _ = strconv.ParseInt(ArrayMCC[0], 10, 64)
	tmpStructCellTowers[0].MobileNetworkCode, _ = strconv.ParseInt(ArrayMNC[0], 10, 64)

	tmpMarshal := StructJasonData{"false", tmpStructCellTowers}
	tmpByte, err := json.Marshal(tmpMarshal)
	if err != nil {
		return "", err
	}

	fmt.Println(string(tmpByte))

	return string(tmpByte), nil
}

func GetPositionByGeolocationLBS(ArrayCellID [3]string, ArrayLAC [3]string, ArrayMNC [3]string, ArrayMCC [3]string) (float64, float64, float64, error) {
	lbsJason, err := LBSInfoToJasonForGeolocation(ArrayCellID, ArrayLAC, ArrayMNC, ArrayMCC)
	if err != nil {
		return 0, 0, 0, err
	}

	req := bytes.NewBuffer([]byte(lbsJason))
	body_type := "application/json;charset=utf-8"
	response, err := http.Post(uRlGeo+GetGeoAPI(), body_type, req)
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

func wifiInfoAndLBSInfoToStringForCellocation(ArrayCellID [3]string, ArrayLAC [3]string, ArrayMNC [3]string, ArrayMCC [3]string, ArrrayWifiInfo []WifiInfo) (string, error) {
	var strtmp string
	strtmp = "http://vip.cellocation.com/loc/nveg5t0h.php?cl=" + ArrayMCC[0] + "," + ArrayMNC[0] + "," + ArrayLAC[0] + "," + ArrayCellID[0] + ",0"

	if len(ArrrayWifiInfo) > 0 {
		strtmp = strtmp + "&wl="
		for i := 0; i < len(ArrrayWifiInfo); i++ {
			if i != 0 {
				strtmp = strtmp + ";"
			}
			strtmp = strtmp + ArrrayWifiInfo[i].Mac + "," + fmt.Sprintf("%d", ArrrayWifiInfo[i].Rssi)
		}
	}

	strtmp = strtmp + "&output=xml"

	return strtmp, nil
}

func responseStringToDataForCellocation(strRes string) (float64, float64, float64, error) {
	type StructData struct {
		XMLName xml.Name `xml:"response"` // 指定最外层的标签为config
		Errcode int      `xml:"errcode"`
		Lat     float64  `xml:"lat"`
		Lon     float64  `xml:"lon"`
		Radius  float64  `xml:"radius"`
		Address string   `xml:"address"`
	}

	tmphttpRes := StructData{}
	xml.Unmarshal([]byte(strRes), &tmphttpRes)
	fmt.Println(tmphttpRes)
	if 0 != tmphttpRes.Errcode {
		fmt.Printf("tmphttpRes.errcode=%d", tmphttpRes.Errcode)
		return 0, 0, 0, errors.New("responseStringToData Error:result not 200")
	}

	return tmphttpRes.Lat, tmphttpRes.Lon, tmphttpRes.Radius, nil
}

func GetPositionByCellocationMix(ArrayCellID [3]string, ArrayLAC [3]string, ArrayMNC [3]string, ArrayMCC [3]string, ArrrayWifiInfo []WifiInfo) (float64, float64, float64, error) {
	strtmp, _ := wifiInfoAndLBSInfoToStringForCellocation(ArrayCellID, ArrayLAC, ArrayMNC, ArrayMCC, ArrrayWifiInfo)
	fmt.Println(strtmp)

	response, err := http.Get(strtmp)
	if err != nil {
		fmt.Println("GetCellocation http err")
		return 0, 0, 0, err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	lat, lon, accRange, err := responseStringToDataForCellocation(string(body))
	if err != nil {
		fmt.Println("GetPositionByByCellocation response err:")
		fmt.Println(string(body))
		return 0, 0, 0, err
	}

	return lat, lon, accRange, nil
}
