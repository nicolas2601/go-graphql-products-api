// Package postgres provee una implementacion de domain.ProductRepository sobre PostgreSQL.
package postgres

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nicolas2601/go-graphql-products-api/internal/domain"
)

//go:embed schema.sql
var schemaSQL string

// codigo SQLSTATE de violacion de unique/primary key en PostgreSQL.
const uniqueViolationCode = "23505"

// Repository persiste productos en PostgreSQL.
type Repository struct {
	pool *pgxpool.Pool
}

// New crea un repositorio PostgreSQL sobre el pool dado.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Verificacion en compilacion de que *Repository satisface el contrato del dominio.
var _ domain.ProductRepository = (*Repository)(nil)

// Migrate aplica el esquema de la tabla products de forma idempotente.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, schemaSQL); err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}
	return nil
}

// Create inserta un producto nuevo. Devuelve ErrProductAlreadyExists ante id duplicado.
func (r *Repository) Create(ctx context.Context, product domain.Product) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO products (id, name, price, stock, created_at) VALUES ($1, $2, $3, $4, $5)`,
		product.ID, product.Name, product.Price, product.Stock, product.CreatedAt)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
		return domain.ErrProductAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("insert product: %w", err)
	}
	return nil
}

// GetByID devuelve un producto por id, o ErrProductNotFound si no existe.
func (r *Repository) GetByID(ctx context.Context, id string) (domain.Product, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, name, price, stock, created_at FROM products WHERE id = $1`, id)

	product, err := scanProduct(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Product{}, domain.ErrProductNotFound
	}
	if err != nil {
		return domain.Product{}, fmt.Errorf("query product: %w", err)
	}
	return product, nil
}

// List devuelve todos los productos ordenados por id.
func (r *Repository) List(ctx context.Context) ([]domain.Product, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, price, stock, created_at FROM products ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("query products: %w", err)
	}
	defer rows.Close()

	products := make([]domain.Product, 0)
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate products: %w", err)
	}
	return products, nil
}

// Update reemplaza nombre, precio y stock de un producto, o ErrProductNotFound si no existe.
func (r *Repository) Update(ctx context.Context, product domain.Product) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE products SET name = $2, price = $3, stock = $4 WHERE id = $1`,
		product.ID, product.Name, product.Price, product.Stock)
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

// Delete elimina un producto por id, o ErrProductNotFound si no existe.
func (r *Repository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

// rowScanner abstrae QueryRow y Rows para reutilizar el escaneo de un producto.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanProduct(row rowScanner) (domain.Product, error) {
	var product domain.Product
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CreatedAt)
	return product, err
}
