package types

type StorageProvider string

const (
	GoogleCloudStorage StorageProvider = "google_cloud"
	LocalStorage       StorageProvider = "local"
)

func (s StorageProvider) IsValid() bool {
	switch s {
	case GoogleCloudStorage, LocalStorage:
		return true
	default:
		return false
	}
}
