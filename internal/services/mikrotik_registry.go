package services

// Simple shared instance registry so background jobs can use the same MikroTik connection

var sharedMikroTikService *MikroTikService

func SetSharedMikroTikService(svc *MikroTikService) {
	sharedMikroTikService = svc
}

func GetSharedMikroTikService() *MikroTikService {
	return sharedMikroTikService
}
