package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ticket represents a ticket type for a tournament
type Ticket struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" example:"64f123abc456def789012345" description:"ID unik tiket"`
	TournamentID      primitive.ObjectID `bson:"tournament_id" json:"tournament_id" example:"64f123abc456def789012345" description:"ID tournament terkait"`
	Price             int                `bson:"price" json:"price" example:"50000" description:"Harga tiket dalam Rupiah"`
	QuantityAvailable int                `bson:"quantity_available" json:"quantity_available" example:"1000" description:"Jumlah tiket yang tersedia"`
	Description       string             `bson:"description,omitempty" json:"description,omitempty" example:"Tiket Reguler - Tribun A (Standing)" description:"Deskripsi jenis tiket"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at" example:"2025-07-16T07:28:37.016Z" description:"Waktu pembuatan tiket"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at" example:"2025-07-16T07:28:37.016Z" description:"Waktu terakhir diupdate"`
}

// TicketRequest represents request body for creating/updating ticket
type TicketRequest struct {
	TournamentID      string `json:"tournament_id" validate:"required" example:"64f123abc456def789012345" swaggertype:"string" description:"ID tournament yang akan dibuat tiketnya (ObjectID format)"`
	Price             int    `json:"price" validate:"required,min=0" example:"50000" description:"Harga tiket dalam Rupiah"`
	QuantityAvailable int    `json:"quantity_available" validate:"required,min=0" example:"1000" description:"Jumlah tiket yang tersedia"`
	Description       string `json:"description,omitempty" example:"Tiket Reguler - Tribun A (Standing)" description:"Deskripsi jenis tiket (opsional)"`
}

// TicketResponse represents response for ticket operations
type TicketResponse struct {
	Message  string `json:"message" example:"Ticket type created successfully" description:"Pesan konfirmasi operasi"`
	TicketID string `json:"ticket_id,omitempty" example:"64f123abc456def789012345" description:"ID tiket yang dibuat (untuk operasi create)"`
}

// TicketWithTournament represents ticket with tournament details
type TicketWithTournament struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" example:"64f123abc456def789012345" description:"ID unik tiket"`
	TournamentID      primitive.ObjectID `bson:"tournament_id" json:"tournament_id" example:"64f123abc456def789012345" description:"ID tournament terkait"`
	Price             int                `bson:"price" json:"price" example:"50000" description:"Harga tiket dalam Rupiah"`
	QuantityAvailable int                `bson:"quantity_available" json:"quantity_available" example:"1000" description:"Jumlah tiket yang tersedia"`
	Description       string             `bson:"description,omitempty" json:"description,omitempty" example:"Tiket Reguler - Tribun A (Standing)" description:"Deskripsi jenis tiket"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at" example:"2025-07-16T07:28:37.016Z" description:"Waktu pembuatan tiket"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at" example:"2025-07-16T07:28:37.016Z" description:"Waktu terakhir diupdate"`
	Tournament        *TournamentBasic   `json:"tournament,omitempty" description:"Detail tournament (tersedia jika populate=true)"`
}

// TournamentBasic represents minimal tournament info for ticket details
type TournamentBasic struct {
	ID     primitive.ObjectID `json:"_id" example:"64f123abc456def789012345" description:"ID unik tournament"`
	Name   string             `json:"name" example:"Mobile Legends Championship 2025" description:"Nama tournament"`
	Status string             `json:"status" example:"upcoming" description:"Status tournament (upcoming/ongoing/completed)"`
}
