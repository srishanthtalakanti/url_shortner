package models

type ShortenUrl struct {
	URL string `json:"url"`
}
type EditUrl struct {
	Short_code string `json:"short_code"`
	Long_url   string `json:"long_url"`
}
type Credentials struct {
	User_Id  int    `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Jwt struct {
	Token string `json:"jwt"`
}
