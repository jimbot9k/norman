package reports

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	dbo "github.com/jimbot9k/norman/internal/core/dbobjects"
)

func TestMermaidReportWriter_WriteInventoryReport(t *testing.T) {
	t.Run("writes mermaid file successfully", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.mmd")

		db := dbo.NewDatabase("testdb", nil)
		schema := dbo.NewSchema("public", "owner", nil)
		db.AddSchema(schema)

		writer := &MermaidReportWriter{}
		err := writer.WriteInventoryReport(filePath, db)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		if !strings.HasPrefix(string(content), "erDiagram") {
			t.Errorf("expected file to start with erDiagram, got %s", string(content))
		}
	})

	t.Run("returns error for invalid path", func(t *testing.T) {
		writer := &MermaidReportWriter{}
		db := dbo.NewDatabase("testdb", nil)

		err := writer.WriteInventoryReport("/nonexistent/path/file.mmd", db)

		if err == nil {
			t.Error("expected error for invalid path")
		}
	})
}

func TestGenerateMermaidERD(t *testing.T) {
	t.Run("empty database", func(t *testing.T) {
		db := dbo.NewDatabase("testdb", nil)

		result := GenerateMermaidERD(db)

		if result != "erDiagram\n" {
			t.Errorf("expected 'erDiagram\\n', got %q", result)
		}
	})

	t.Run("database with schema and tables", func(t *testing.T) {
		db := dbo.NewDatabase("testdb", nil)
		schema := dbo.NewSchema("public", "owner", nil)
		db.AddSchema(schema)

		table := dbo.NewTable("users", nil)
		col := dbo.NewColumn("id", "integer", false)
		table.AddColumn(col)
		schema.AddTable(table)

		result := GenerateMermaidERD(db)

		if !strings.Contains(result, "erDiagram") {
			t.Error("expected erDiagram header")
		}
		if !strings.Contains(result, "users") {
			t.Error("expected users table")
		}
		if !strings.Contains(result, "int id") {
			t.Error("expected id column")
		}
	})

	t.Run("database with multiple schemas", func(t *testing.T) {
		db := dbo.NewDatabase("testdb", nil)

		schema1 := dbo.NewSchema("public", "owner", nil)
		table1 := dbo.NewTable("users", nil)
		table1.AddColumn(dbo.NewColumn("id", "integer", false))
		schema1.AddTable(table1)
		db.AddSchema(schema1)

		schema2 := dbo.NewSchema("admin", "owner", nil)
		table2 := dbo.NewTable("roles", nil)
		table2.AddColumn(dbo.NewColumn("id", "integer", false))
		schema2.AddTable(table2)
		db.AddSchema(schema2)

		result := GenerateMermaidERD(db)

		if !strings.Contains(result, "users") {
			t.Error("expected users table from schema1")
		}
		if !strings.Contains(result, "roles") {
			t.Error("expected roles table from schema2")
		}
	})
}

func TestWriteSchemaEntities(t *testing.T) {
	t.Run("schema with tables and relationships", func(t *testing.T) {
		schema := dbo.NewSchema("public", "owner", nil)

		// Create parent table
		parentTable := dbo.NewTable("departments", nil)
		parentCol := dbo.NewColumn("id", "integer", false)
		parentTable.AddColumn(parentCol)
		parentTable.SetPrimaryKey(dbo.NewPrimaryKey("departments_pkey", parentTable, []*dbo.Column{parentCol}))
		schema.AddTable(parentTable)

		// Create child table with FK
		childTable := dbo.NewTable("employees", nil)
		childCol := dbo.NewColumn("id", "integer", false)
		deptIdCol := dbo.NewColumn("department_id", "integer", false)
		childTable.AddColumn(childCol)
		childTable.AddColumn(deptIdCol)

		fk := dbo.NewForeignKey("fk_dept", "departments")
		fk.SetTable(childTable)
		fk.AddColumn(deptIdCol)
		fk.AddReferencedColumn(parentCol)
		childTable.AddForeignKey(fk)
		schema.AddTable(childTable)

		var sb strings.Builder
		writeSchemaEntities(&sb, schema)
		result := sb.String()

		if !strings.Contains(result, "departments") {
			t.Error("expected departments table")
		}
		if !strings.Contains(result, "employees") {
			t.Error("expected employees table")
		}
		if !strings.Contains(result, "}o--||") {
			t.Error("expected relationship")
		}
	})

	t.Run("empty schema", func(t *testing.T) {
		schema := dbo.NewSchema("empty", "owner", nil)

		var sb strings.Builder
		writeSchemaEntities(&sb, schema)
		result := sb.String()

		if result != "" {
			t.Errorf("expected empty string for empty schema, got %q", result)
		}
	})
}

func TestWriteTableEntity(t *testing.T) {
	t.Run("table with columns", func(t *testing.T) {
		table := dbo.NewTable("users", nil)
		table.AddColumn(dbo.NewColumn("id", "integer", false))
		table.AddColumn(dbo.NewColumn("name", "varchar", true))

		var sb strings.Builder
		writeTableEntity(&sb, table)
		result := sb.String()

		if !strings.Contains(result, "users {") {
			t.Error("expected table declaration")
		}
		if !strings.Contains(result, "int id") {
			t.Error("expected id column")
		}
		if !strings.Contains(result, "varchar name") {
			t.Error("expected name column")
		}
		if !strings.Contains(result, "}") {
			t.Error("expected closing brace")
		}
	})

	t.Run("table with no columns", func(t *testing.T) {
		table := dbo.NewTable("empty_table", nil)

		var sb strings.Builder
		writeTableEntity(&sb, table)
		result := sb.String()

		expected := "    empty_table {\n    }\n"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})
}

func TestWriteColumnDefinition(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() (*dbo.Table, *dbo.Column)
		contains []string
	}{
		{
			name: "regular column",
			setup: func() (*dbo.Table, *dbo.Column) {
				table := dbo.NewTable("users", nil)
				col := dbo.NewColumn("name", "varchar", true)
				table.AddColumn(col)
				return table, col
			},
			contains: []string{"varchar name"},
		},
		{
			name: "primary key column",
			setup: func() (*dbo.Table, *dbo.Column) {
				table := dbo.NewTable("users", nil)
				col := dbo.NewColumn("id", "integer", false)
				table.AddColumn(col)
				table.SetPrimaryKey(dbo.NewPrimaryKey("pk", table, []*dbo.Column{col}))
				return table, col
			},
			contains: []string{"int id PK"},
		},
		{
			name: "foreign key column",
			setup: func() (*dbo.Table, *dbo.Column) {
				table := dbo.NewTable("orders", nil)
				col := dbo.NewColumn("user_id", "integer", false)
				table.AddColumn(col)
				fk := dbo.NewForeignKey("fk_user", "users")
				fk.SetTable(table)
				fk.AddColumn(col)
				table.AddForeignKey(fk)
				return table, col
			},
			contains: []string{"int user_id FK"},
		},
		{
			name: "column that is both PK and FK",
			setup: func() (*dbo.Table, *dbo.Column) {
				table := dbo.NewTable("user_roles", nil)
				col := dbo.NewColumn("user_id", "integer", false)
				table.AddColumn(col)
				table.SetPrimaryKey(dbo.NewPrimaryKey("pk", table, []*dbo.Column{col}))
				fk := dbo.NewForeignKey("fk_user", "users")
				fk.SetTable(table)
				fk.AddColumn(col)
				table.AddForeignKey(fk)
				return table, col
			},
			contains: []string{"int user_id PK \"FK\""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table, col := tt.setup()

			var sb strings.Builder
			writeColumnDefinition(&sb, table, col)
			result := sb.String()

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("expected result to contain %q, got %q", expected, result)
				}
			}
		})
	}
}

func TestCollectRelationships(t *testing.T) {
	t.Run("schema with relationships", func(t *testing.T) {
		schema := dbo.NewSchema("public", "owner", nil)

		parentTable := dbo.NewTable("departments", nil)
		parentCol := dbo.NewColumn("id", "integer", false)
		parentTable.AddColumn(parentCol)
		schema.AddTable(parentTable)

		childTable := dbo.NewTable("employees", nil)
		childCol := dbo.NewColumn("department_id", "integer", false)
		childTable.AddColumn(childCol)
		fk := dbo.NewForeignKey("fk_dept", "departments")
		fk.SetTable(childTable)
		fk.AddColumn(childCol)
		childTable.AddForeignKey(fk)
		schema.AddTable(childTable)

		relationships := collectRelationships(schema)

		if len(relationships) != 1 {
			t.Errorf("expected 1 relationship, got %d", len(relationships))
		}
		if !strings.Contains(relationships[0], "employees") {
			t.Error("expected employees in relationship")
		}
		if !strings.Contains(relationships[0], "departments") {
			t.Error("expected departments in relationship")
		}
	})

	t.Run("duplicate relationships are deduplicated", func(t *testing.T) {
		schema := dbo.NewSchema("public", "owner", nil)

		parentTable := dbo.NewTable("departments", nil)
		schema.AddTable(parentTable)

		childTable := dbo.NewTable("employees", nil)
		col1 := dbo.NewColumn("dept_id", "integer", false)
		childTable.AddColumn(col1)

		// Add two FKs to the same table (simulating same relationship)
		fk1 := dbo.NewForeignKey("fk_dept", "departments")
		fk1.SetTable(childTable)
		fk1.AddColumn(col1)
		childTable.AddForeignKey(fk1)

		// Add same FK again (same name = same formatted output)
		fk2 := dbo.NewForeignKey("fk_dept", "departments")
		fk2.SetTable(childTable)
		fk2.AddColumn(col1)
		childTable.AddForeignKey(fk2)

		schema.AddTable(childTable)

		relationships := collectRelationships(schema)

		if len(relationships) != 1 {
			t.Errorf("expected 1 unique relationship (deduped), got %d", len(relationships))
		}
	})

	t.Run("schema with no relationships", func(t *testing.T) {
		schema := dbo.NewSchema("public", "owner", nil)
		table := dbo.NewTable("standalone", nil)
		table.AddColumn(dbo.NewColumn("id", "integer", false))
		schema.AddTable(table)

		relationships := collectRelationships(schema)

		if len(relationships) != 0 {
			t.Errorf("expected 0 relationships, got %d", len(relationships))
		}
	})
}

func TestFormatMermaidRelationship(t *testing.T) {
	t.Run("formats relationship correctly", func(t *testing.T) {
		table := dbo.NewTable("orders", nil)
		col := dbo.NewColumn("user_id", "integer", false)
		table.AddColumn(col)

		fk := dbo.NewForeignKey("fk_user_orders", "users")
		fk.SetTable(table)
		fk.AddColumn(col)

		result := formatMermaidRelationship(table, fk)

		if !strings.Contains(result, "orders }o--|| users") {
			t.Errorf("expected relationship format, got %q", result)
		}
		if !strings.Contains(result, ": \"fk_user_orders\"") {
			t.Errorf("expected FK name label, got %q", result)
		}
	})

	t.Run("handles names with hyphens", func(t *testing.T) {
		table := dbo.NewTable("order-items", nil)
		col := dbo.NewColumn("order_id", "integer", false)
		table.AddColumn(col)

		fk := dbo.NewForeignKey("fk-order-items", "orders-main")
		fk.SetTable(table)
		fk.AddColumn(col)

		result := formatMermaidRelationship(table, fk)

		// Hyphens in table names and FK names should be replaced
		if !strings.Contains(result, "order_items") {
			t.Error("expected order_items (sanitized)")
		}
		if !strings.Contains(result, "orders_main") {
			t.Error("expected orders_main (sanitized)")
		}
		if !strings.Contains(result, "fk_order_items") {
			t.Error("expected fk_order_items (sanitized)")
		}
	})
}

func TestSanitizeMermaidName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple name", "users", "users"},
		{"name with hyphen", "user-roles", "user_roles"},
		{"name with space", "user roles", "user_roles"},
		{"name with multiple hyphens", "user-role-mapping", "user_role_mapping"},
		{"name with hyphen and space", "user-role mapping", "user_role_mapping"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeMermaidName(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestNormalizeMermaidDataType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// varchar variants
		{"character varying", "character varying", "varchar"},
		{"character varying with length", "character varying(255)", "varchar"},
		{"varchar", "varchar", "varchar"},
		{"VARCHAR uppercase", "VARCHAR", "varchar"},

		// integer variants
		{"integer", "integer", "int"},
		{"int", "int", "int"},
		{"int4", "int4", "int"},
		{"INTEGER uppercase", "INTEGER", "int"},

		// bigint variants
		{"bigint", "bigint", "bigint"},
		{"int8", "int8", "bigint"},
		{"BIGINT uppercase", "BIGINT", "bigint"},

		// smallint variants
		{"smallint", "smallint", "smallint"},
		{"int2", "int2", "smallint"},

		// numeric/decimal
		{"numeric", "numeric", "decimal"},
		{"numeric with precision", "numeric(10,2)", "decimal"},
		{"decimal", "decimal", "decimal"},

		// boolean
		{"boolean", "boolean", "bool"},
		{"bool", "bool", "bool"},

		// timestamp
		{"timestamp", "timestamp", "timestamp"},
		{"timestamp with time zone", "timestamp with time zone", "timestamp"},

		// date/time
		{"date", "date", "date"},
		{"time", "time", "time"},

		// text
		{"text", "text", "text"},

		// uuid
		{"uuid", "uuid", "uuid"},

		// json
		{"json", "json", "json"},
		{"jsonb", "jsonb", "json"},

		// default cases
		{"unknown type", "bytea", "bytea"},
		{"type with parentheses", "bit(8)", "bit"},
		{"custom type", "my_custom_type", "my_custom_type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeMermaidDataType(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestIsPrimaryKeyColumn(t *testing.T) {
	t.Run("column is in primary key", func(t *testing.T) {
		table := dbo.NewTable("users", nil)
		col := dbo.NewColumn("id", "integer", false)
		table.AddColumn(col)
		table.SetPrimaryKey(dbo.NewPrimaryKey("pk", table, []*dbo.Column{col}))

		if !isPrimaryKeyColumn(table, "id") {
			t.Error("expected id to be a primary key column")
		}
	})

	t.Run("column is not in primary key", func(t *testing.T) {
		table := dbo.NewTable("users", nil)
		idCol := dbo.NewColumn("id", "integer", false)
		nameCol := dbo.NewColumn("name", "varchar", false)
		table.AddColumn(idCol)
		table.AddColumn(nameCol)
		table.SetPrimaryKey(dbo.NewPrimaryKey("pk", table, []*dbo.Column{idCol}))

		if isPrimaryKeyColumn(table, "name") {
			t.Error("expected name to not be a primary key column")
		}
	})

	t.Run("table has no primary key", func(t *testing.T) {
		table := dbo.NewTable("users", nil)
		col := dbo.NewColumn("id", "integer", false)
		table.AddColumn(col)

		if isPrimaryKeyColumn(table, "id") {
			t.Error("expected false when table has no primary key")
		}
	})

	t.Run("composite primary key", func(t *testing.T) {
		table := dbo.NewTable("user_roles", nil)
		userIdCol := dbo.NewColumn("user_id", "integer", false)
		roleIdCol := dbo.NewColumn("role_id", "integer", false)
		table.AddColumn(userIdCol)
		table.AddColumn(roleIdCol)
		table.SetPrimaryKey(dbo.NewPrimaryKey("pk", table, []*dbo.Column{userIdCol, roleIdCol}))

		if !isPrimaryKeyColumn(table, "user_id") {
			t.Error("expected user_id to be in composite PK")
		}
		if !isPrimaryKeyColumn(table, "role_id") {
			t.Error("expected role_id to be in composite PK")
		}
	})
}

func TestIsForeignKeyColumn(t *testing.T) {
	t.Run("column is in foreign key", func(t *testing.T) {
		table := dbo.NewTable("orders", nil)
		col := dbo.NewColumn("user_id", "integer", false)
		table.AddColumn(col)

		fk := dbo.NewForeignKey("fk_user", "users")
		fk.SetTable(table)
		fk.AddColumn(col)
		table.AddForeignKey(fk)

		if !isForeignKeyColumn(table, "user_id") {
			t.Error("expected user_id to be a foreign key column")
		}
	})

	t.Run("column is not in any foreign key", func(t *testing.T) {
		table := dbo.NewTable("orders", nil)
		idCol := dbo.NewColumn("id", "integer", false)
		userIdCol := dbo.NewColumn("user_id", "integer", false)
		table.AddColumn(idCol)
		table.AddColumn(userIdCol)

		fk := dbo.NewForeignKey("fk_user", "users")
		fk.SetTable(table)
		fk.AddColumn(userIdCol)
		table.AddForeignKey(fk)

		if isForeignKeyColumn(table, "id") {
			t.Error("expected id to not be a foreign key column")
		}
	})

	t.Run("table has no foreign keys", func(t *testing.T) {
		table := dbo.NewTable("users", nil)
		col := dbo.NewColumn("id", "integer", false)
		table.AddColumn(col)

		if isForeignKeyColumn(table, "id") {
			t.Error("expected false when table has no foreign keys")
		}
	})

	t.Run("composite foreign key", func(t *testing.T) {
		table := dbo.NewTable("order_items", nil)
		orderIdCol := dbo.NewColumn("order_id", "integer", false)
		productIdCol := dbo.NewColumn("product_id", "integer", false)
		table.AddColumn(orderIdCol)
		table.AddColumn(productIdCol)

		fk := dbo.NewForeignKey("fk_order_product", "orders")
		fk.SetTable(table)
		fk.AddColumn(orderIdCol)
		fk.AddColumn(productIdCol)
		table.AddForeignKey(fk)

		if !isForeignKeyColumn(table, "order_id") {
			t.Error("expected order_id to be in composite FK")
		}
		if !isForeignKeyColumn(table, "product_id") {
			t.Error("expected product_id to be in composite FK")
		}
	})

	t.Run("multiple foreign keys", func(t *testing.T) {
		table := dbo.NewTable("orders", nil)
		userIdCol := dbo.NewColumn("user_id", "integer", false)
		shippingIdCol := dbo.NewColumn("shipping_address_id", "integer", false)
		table.AddColumn(userIdCol)
		table.AddColumn(shippingIdCol)

		fk1 := dbo.NewForeignKey("fk_user", "users")
		fk1.SetTable(table)
		fk1.AddColumn(userIdCol)
		table.AddForeignKey(fk1)

		fk2 := dbo.NewForeignKey("fk_shipping", "addresses")
		fk2.SetTable(table)
		fk2.AddColumn(shippingIdCol)
		table.AddForeignKey(fk2)

		if !isForeignKeyColumn(table, "user_id") {
			t.Error("expected user_id to be FK")
		}
		if !isForeignKeyColumn(table, "shipping_address_id") {
			t.Error("expected shipping_address_id to be FK")
		}
	})
}

func TestIntegration_FullERD(t *testing.T) {
	t.Run("generates complete ERD for complex schema", func(t *testing.T) {
		db := dbo.NewDatabase("ecommerce", nil)
		schema := dbo.NewSchema("public", "admin", nil)
		db.AddSchema(schema)

		// Users table
		usersTable := dbo.NewTable("users", nil)
		userIdCol := dbo.NewColumn("id", "integer", false)
		usersTable.AddColumn(userIdCol)
		usersTable.AddColumn(dbo.NewColumn("email", "character varying(255)", false))
		usersTable.AddColumn(dbo.NewColumn("created_at", "timestamp with time zone", false))
		usersTable.SetPrimaryKey(dbo.NewPrimaryKey("users_pkey", usersTable, []*dbo.Column{userIdCol}))
		schema.AddTable(usersTable)

		// Orders table
		ordersTable := dbo.NewTable("orders", nil)
		orderIdCol := dbo.NewColumn("id", "integer", false)
		orderUserIdCol := dbo.NewColumn("user_id", "integer", false)
		ordersTable.AddColumn(orderIdCol)
		ordersTable.AddColumn(orderUserIdCol)
		ordersTable.AddColumn(dbo.NewColumn("total", "numeric(10,2)", false))
		ordersTable.SetPrimaryKey(dbo.NewPrimaryKey("orders_pkey", ordersTable, []*dbo.Column{orderIdCol}))

		fkUser := dbo.NewForeignKey("fk_orders_user", "users")
		fkUser.SetTable(ordersTable)
		fkUser.AddColumn(orderUserIdCol)
		ordersTable.AddForeignKey(fkUser)
		schema.AddTable(ordersTable)

		// Order items (junction table with composite PK)
		itemsTable := dbo.NewTable("order_items", nil)
		itemOrderIdCol := dbo.NewColumn("order_id", "integer", false)
		itemProductIdCol := dbo.NewColumn("product_id", "integer", false)
		itemsTable.AddColumn(itemOrderIdCol)
		itemsTable.AddColumn(itemProductIdCol)
		itemsTable.AddColumn(dbo.NewColumn("quantity", "smallint", false))
		itemsTable.SetPrimaryKey(dbo.NewPrimaryKey("order_items_pkey", itemsTable, []*dbo.Column{itemOrderIdCol, itemProductIdCol}))

		fkOrder := dbo.NewForeignKey("fk_items_order", "orders")
		fkOrder.SetTable(itemsTable)
		fkOrder.AddColumn(itemOrderIdCol)
		itemsTable.AddForeignKey(fkOrder)

		fkProduct := dbo.NewForeignKey("fk_items_product", "products")
		fkProduct.SetTable(itemsTable)
		fkProduct.AddColumn(itemProductIdCol)
		itemsTable.AddForeignKey(fkProduct)
		schema.AddTable(itemsTable)

		// Products table
		productsTable := dbo.NewTable("products", nil)
		productIdCol := dbo.NewColumn("id", "integer", false)
		productsTable.AddColumn(productIdCol)
		productsTable.AddColumn(dbo.NewColumn("name", "text", false))
		productsTable.AddColumn(dbo.NewColumn("price", "decimal", false))
		productsTable.AddColumn(dbo.NewColumn("is_active", "boolean", false))
		productsTable.SetPrimaryKey(dbo.NewPrimaryKey("products_pkey", productsTable, []*dbo.Column{productIdCol}))
		schema.AddTable(productsTable)

		result := GenerateMermaidERD(db)

		// Verify structure
		if !strings.HasPrefix(result, "erDiagram\n") {
			t.Error("expected erDiagram header")
		}

		// Verify all tables present
		tables := []string{"users", "orders", "order_items", "products"}
		for _, tbl := range tables {
			if !strings.Contains(result, tbl+" {") {
				t.Errorf("expected table %s", tbl)
			}
		}

		// Verify data types normalized
		if !strings.Contains(result, "varchar email") {
			t.Error("expected varchar for email")
		}
		if !strings.Contains(result, "timestamp created_at") {
			t.Error("expected timestamp for created_at")
		}
		if !strings.Contains(result, "decimal total") {
			t.Error("expected decimal for total")
		}
		if !strings.Contains(result, "smallint quantity") {
			t.Error("expected smallint for quantity")
		}
		if !strings.Contains(result, "bool is_active") {
			t.Error("expected bool for is_active")
		}

		// Verify PK/FK markers
		if !strings.Contains(result, "int id PK") {
			t.Error("expected PK marker on id columns")
		}
		if !strings.Contains(result, "int user_id FK") {
			t.Error("expected FK marker on user_id")
		}
		if !strings.Contains(result, "PK \"FK\"") {
			t.Error("expected PK+FK marker on junction table columns")
		}

		// Verify relationships
		if !strings.Contains(result, "orders }o--|| users") {
			t.Error("expected orders->users relationship")
		}
		if !strings.Contains(result, "order_items }o--|| orders") {
			t.Error("expected order_items->orders relationship")
		}
		if !strings.Contains(result, "order_items }o--|| products") {
			t.Error("expected order_items->products relationship")
		}
	})
}
