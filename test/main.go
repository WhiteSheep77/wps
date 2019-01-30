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

	ArrrayWifiInfo = make([]wps.WifiInfo, 0)
	lat, lon, accRange, err = wps.GetPositionByOpenData(ArrrayWifiInfo)
	fmt.Println(lat, lon, accRange, err)
}
