package wps //WifiPositioningSystem

const uRlGeo string = "https://www.googleapis.com/geolocation/v1/geolocate?key="
const aPIGeo string = "AIzaSyCpnL8k0TZJ-BUxKyRlVobneWNuMQtcEpI"
const aPIGeo1 string = "AIzaSyAUJ5IKYiDXj16H2uYbahPWhLQwZ_ih_AU"

var count int = 0

func GetGeoAPI() string {
	count++
	if count > 1 {
		count = 0
	}
	if 0 == count {

		return aPIGeo
	} else {
		return aPIGeo1
	}
}
