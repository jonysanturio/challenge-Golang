# Database Schema Diagram

## Entity Relationship Diagram (Text Representation)

```
+----------------+      +---------------------+      +----------------+
|    Products    |      | Product_Categories  |      |  Categories    |
+----------------+      +---------------------+      +----------------+
| id (PK)        |<>----| product_id (FK, PK) |>----<| id (PK)        |
| name           |      | category_id (FK, PK)|      | name           |
| description    |      +---------------------+      | description    |
| price          |                                  created_at    |
| stock          |                                  updated_at    |
+----------------+
| Product_History|
+----------------+
| id (PK)        |
| product_id (FK)|
| price          |
| stock          |
| changed_at     |
+----------------+

## Table Details

### products
- **id**: Primary key, auto-increment
- **name**: Product name (VARCHAR, required)
- **description**: Product description (TEXT)
- **price**: Product price (DECIMAL(10,2), required)
- **stock**: Available quantity (INTEGER, required)
- **created_at**: Timestamp when created
- **updated_at**: Timestamp when last updated

### categories
- **id**: Primary key, auto-increment
- **name**: Category name (VARCHAR, required, unique)
- **description**: Category description (TEXT)
- **created_at**: Timestamp when created
- **updated_at**: Timestamp when last updated

### product_categories (Join Table)
- **product_id**: Foreign key to products.id (part of composite PK)
- **category_id**: Foreign key to categories.id (part of composite PK)

### product_history
- **id**: Primary key, auto-increment
- **product_id**: Foreign key to products.id
- **price**: Price at time of change (DECIMAL(10,2))
- **stock**: Stock level at time of change (INTEGER)
- **changed_at**: Timestamp when change occurred

## Indexes
- Primary keys on all id fields
- Foreign key indexes for performance
- Composite primary key on (product_id, category_id) in product_categories
- Index on product_id in product_history for efficient history queries

## Constraints
- NOT NULL on required fields
- UNIQUE constraint on categories.name
- FOREIGN KEY constraints with CASCADE DELETE where appropriate
- CHECK constraints for price >= 0 and stock >= 0 (implemented via application validation)