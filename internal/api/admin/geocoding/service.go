package geocoding

type AdminGeocodingServiceInterface interface {
	ReverseGeocodeService(lat, lng string) (map[string]interface{}, error)
}

type AdminGeocodingServiceStruct struct {
	repository AdminGeocodingRepositoryInterface
}

func NewAdminGeocodingService(repository AdminGeocodingRepositoryInterface) AdminGeocodingServiceInterface {
	return &AdminGeocodingServiceStruct{repository: repository}
}

func (s AdminGeocodingServiceStruct) ReverseGeocodeService(lat, lng string) (map[string]interface{}, error) {
	return s.repository.ReverseGeocodeRepository(lat, lng)
}
