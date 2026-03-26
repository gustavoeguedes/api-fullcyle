package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/gustavoeguedes/api-fullcycle/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestProduct_Create(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Product 1", 10)
	assert.Nil(t, err)
	productDb := NewProduct(db)
	err = productDb.Create(product)
	assert.Nil(t, err)
	assert.NotEmpty(t, product.ID)
	assert.Equal(t, "Product 1", product.Name)
}

func TestProduct_FindAll(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	productDb := NewProduct(db)
	for i := 1; i <= 25; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*100)
		assert.Nil(t, err)
		err = productDb.Create(product)
		assert.Nil(t, err)
	}
	products, err := productDb.FindAll(1, 10, "asc")
	assert.Nil(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = productDb.FindAll(1, 10, "desc")
	assert.Nil(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 25", products[0].Name)
	assert.Equal(t, "Product 16", products[9].Name)
}

func TestProduct_FindByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Product 1", 10)
	assert.Nil(t, err)
	productDb := NewProduct(db)
	err = productDb.Create(product)
	assert.Nil(t, err)

	productFound, err := productDb.FindByID(product.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, product.ID, productFound.ID)
	assert.Equal(t, product.Name, productFound.Name)
}

func TestProduct_Delete(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Product 1", 10)
	assert.Nil(t, err)
	productDb := NewProduct(db)
	err = productDb.Create(product)
	assert.Nil(t, err)

	err = productDb.Delete(product.ID.String())
	assert.Nil(t, err)

	_, err = productDb.FindByID(product.ID.String())
	assert.NotNil(t, err)
}

func TestProduct_Update(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Product 1", 10)
	assert.Nil(t, err)
	productDb := NewProduct(db)
	err = productDb.Create(product)
	assert.Nil(t, err)

	product.Name = "Updated Product"
	err = productDb.Update(product)
	assert.Nil(t, err)

	productFound, err := productDb.FindByID(product.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, "Updated Product", productFound.Name)
}
