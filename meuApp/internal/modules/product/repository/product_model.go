package repository

import (
	"time"

	"meuApp/pkg/contracts"
)

// ProductModel representa a estrutura da tabela products no banco
type ProductModel struct {
	ID          string    `gorm:"primaryKey;size:36"`
	Name        string    `gorm:"size:100;not null"`
	Description string    `gorm:"size:500"`
	Price       float64   `gorm:"type:decimal(10,2);not null"`
	Stock       int       `gorm:"default:0;not null"`
	CategoryID  string    `gorm:"size:36;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// TableName especifica o nome da tabela
func (ProductModel) TableName() string {
	return "products"
}

// ToContract converte ProductModel para contracts.Product
func (p *ProductModel) ToContract() *contracts.Product {
	return &contracts.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CategoryID:  p.CategoryID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// FromContract converte contracts.Product para ProductModel
func (p *ProductModel) FromContract(product *contracts.Product) {
	p.ID = product.ID
	p.Name = product.Name
	p.Description = product.Description
	p.Price = product.Price
	p.Stock = product.Stock
	p.CategoryID = product.CategoryID
	p.CreatedAt = product.CreatedAt
	p.UpdatedAt = product.UpdatedAt
}
