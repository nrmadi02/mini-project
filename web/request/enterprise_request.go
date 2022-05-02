package request

type CreateEnterpriseRequest struct {
	Name        string   `json:"name"`
	NumberPhone string   `json:"number_phone"`
	Address     string   `json:"address"`
	Postcode    int      `json:"postcode"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Latitude    string   `json:"latitude"`
	Longitude   string   `json:"longitude"`
}
