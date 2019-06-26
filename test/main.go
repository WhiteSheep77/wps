package main

import (
	"fmt"

	"github.com/WhiteSheep77/wps"
)

func main() {
	var ArrrayWifiInfo = make([]wps.WifiInfo, 3)
	ArrrayWifiInfo[0].Mac = "00:1D:AA:83:E4:60"
	ArrrayWifiInfo[0].Rssi = -56
	ArrrayWifiInfo[1].Mac = "9C:D6:43:F2:CD:74"
	ArrrayWifiInfo[1].Rssi = -73
	ArrrayWifiInfo[2].Mac = "6C:19:8F:E3:BD:CF"
	ArrrayWifiInfo[2].Rssi = -78
	lat, lon, accRange, err := wps.GetPositionByOpenData(ArrrayWifiInfo)
	fmt.Println(lat, lon, accRange, err)

	ArrrayWifiInfo = make([]wps.WifiInfo, 1)
	ArrrayWifiInfo[0].Mac = "00:1D:AA:83:E4:60"
	ArrrayWifiInfo[0].Rssi = -56
	lat, lon, accRange, err = wps.GetPositionByOpenData(ArrrayWifiInfo)
	fmt.Println(lat, lon, accRange, err)

	ArrrayWifiInfo = make([]wps.WifiInfo, 1)
	ArrrayWifiInfo[0].Mac = "00:1D:55:83:E4:60"
	ArrrayWifiInfo[0].Rssi = -56
	lat, lon, accRange, err = wps.GetPositionByOpenData(ArrrayWifiInfo)
	fmt.Println(lat, lon, accRange, err)

	ArrrayWifiInfo = make([]wps.WifiInfo, 3)
	ArrrayWifiInfo[0].Mac = "00:1D:AA:83:E4:60"
	ArrrayWifiInfo[0].Rssi = -56
	ArrrayWifiInfo[1].Mac = "9C:D6:43:F2:CD:74"
	ArrrayWifiInfo[1].Rssi = -73
	ArrrayWifiInfo[2].Mac = "6C:19:8F:E3:BD:CF"
	ArrrayWifiInfo[2].Rssi = -78
	//	lat, lon, accRange, err = wps.GetPositionByGeolocation(ArrrayWifiInfo)
	fmt.Println(lat, lon, accRange, err)

	ArrrayWifiInfo = make([]wps.WifiInfo, 1)
	ArrrayWifiInfo[0].Mac = "00:1D:AA:83:E4:60"
	ArrrayWifiInfo[0].Rssi = -56
	//	lat, lon, accRange, err = wps.GetPositionByGeolocation(ArrrayWifiInfo)
	fmt.Println(lat, lon, accRange, err)

	var ArrayCellID [3]string
	var ArrayLAC [3]string
	var ArrayMNC [3]string
	var ArrayMCC [3]string

	ArrayCellID[0] = "80676995"
	ArrayLAC[0] = "380"
	ArrayMNC[0] = "92"
	ArrayMCC[0] = "466"
	lat, lon, accRange, err = wps.GetPositionByGeolocationLBS(ArrayCellID, ArrayLAC, ArrayMNC, ArrayMCC)
	fmt.Println(lat, lon, accRange, err)

	ArrayCellID[0] = "21532831"
	ArrayLAC[0] = "2862"
	ArrayMNC[0] = "7"
	ArrayMCC[0] = "214"
	lat, lon, accRange, err = wps.GetPositionByGeolocationLBS(ArrayCellID, ArrayLAC, ArrayMNC, ArrayMCC)
	fmt.Println(lat, lon, accRange, err)

}
