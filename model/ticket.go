package model

type Ticket struct {
	ID           string `json:"id" bson:"id,omitempty"`
	TournamentID string `json:"tournament_id" bson:"tournament_id"`
	Price        int    `json:"price" bson:"price"`
	Category     string `json:"category" bson:"category"` // "VIP", "Regular", etc.
	Quantity     int    `json:"quantity" bson:"quantity"`
	Available    int    `json:"available" bson:"available"` // jumlah tiket yang masih tersedia
	Name         string `json:"name" bson:"name"`
	Description  string `json:"description" bson:"description"`
}

type Order struct {
	ID           string `json:"id" bson:"id,omitempty"`
	TicketID     string `json:"ticket_id" bson:"ticket_id"`
	TournamentID string `json:"tournament_id" bson:"tournament_id"`
	BuyerID      string `json:"buyer_id" bson:"buyer_id"` // ID pembeli
	Quantity     int    `json:"quantity" bson:"quantity"`
	TotalAmount  int    `json:"total_amount" bson:"total_amount"`
}
