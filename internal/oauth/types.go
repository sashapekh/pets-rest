package oauth

type User struct {
	Provider      string
	ProviderID    string
	Email         string
	EmailVerified bool
	Name          string
	AvatarURL     string
}
