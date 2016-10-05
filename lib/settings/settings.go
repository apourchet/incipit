package settings

type SettingsClient interface {
	Set(userId string, field string, value string) error
	Get(userId string, field string) (value string, err error)
}
