package api

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectToDB establishes a connection to the PostgreSQL database
func PoolConnectToDB(connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return pool, nil
}

// checkAndCreateTable checks if a table exists and creates it if it does not
func checkAndCreateTable(db *pgxpool.Pool, tableName, createTableSQL string) error {
	var exists bool
	// Check if the table exists
	err := db.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1);", tableName).Scan(&exists)
	if err != nil {
		return err
	}

	// If the table does not exist, create it
	if !exists {
		_, err = db.Exec(context.Background(), createTableSQL)
		if err != nil {
			return err
		}
		fmt.Printf("Table '%s' created.\n", tableName)
	}

	return nil
}

// checkAndCreateTableChat checks and creates the chat table
func CheckAndCreateTableChat(db *pgxpool.Pool) error {
	// Check and create the chat table
	chatTableSQL := `
		-- 創建主分區表（父表），包含 order_date 作為唯一約束的一部分
		CREATE TABLE orders (
			order_id SERIAL,
			customer_id INT,
			order_date DATE,
			total_amount NUMERIC(10, 2),
			PRIMARY KEY (order_id, order_date)  -- 添加 order_date 至主鍵
		) PARTITION BY RANGE (order_date);

		CREATE OR REPLACE FUNCTION create_yearly_partitions(year INT) 
		RETURNS void AS $$
		DECLARE
			start_date DATE;
			end_date DATE;
			partition_name TEXT;
		BEGIN
			-- 設置分區的開始和結束日期
			start_date := TO_DATE(year || '-01-01', 'YYYY-MM-DD');
			end_date := TO_DATE((year + 1) || '-01-01', 'YYYY-MM-DD');

			-- 設置分區名稱
			partition_name := 'orders_' || year;

			-- 檢查分區表是否已經存在，如果不存在則創建
			IF NOT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = partition_name) THEN
				EXECUTE format('
					CREATE TABLE %I PARTITION OF orders
					FOR VALUES FROM (%L) TO (%L);',
					partition_name, start_date, end_date);
			END IF;
		END;
		$$ LANGUAGE plpgsql;

		-- 創建 2023 年的分區
		SELECT create_yearly_partitions(2023);

		-- 創建 2024 年的分區
		SELECT create_yearly_partitions(2024);

		DO $$ 
		DECLARE
			i INT;
		BEGIN
			-- 確保分區存在
			PERFORM create_yearly_partitions(EXTRACT(YEAR FROM NOW())::INT);

			-- 開始插入資料
			FOR i IN 1..1000000 LOOP
				INSERT INTO orders (customer_id, order_date, total_amount)
				VALUES (
					FLOOR(RANDOM() * 1000),  -- 隨機生成 customer_id
					NOW() - INTERVAL '1 day' * (i % 365),  -- 隨機生成 order_date
					FLOOR(RANDOM() * 1000 * 100) / 100.0  -- 隨機生成 total_amount，保留兩位小數
				);
			END LOOP;
		END $$;
		`
	if err := checkAndCreateTable(db, "orders", chatTableSQL); err != nil {
		return err
	}

	return nil
}

// CreateIndex creates a new index on the orders table
func CreateIndex(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), "CREATE INDEX idx_customer_order_date ON orders (customer_id, order_date)")
	if err != nil {
		return err
	}
	fmt.Println("Created index idx_customer_order_date")
	return nil
}

// DropIndex drops the specified index
func DropIndex(db *pgxpool.Pool, indexName string) error {
	_, err := db.Exec(context.Background(), fmt.Sprintf("DROP INDEX IF EXISTS %s", indexName))
	if err != nil {
		return err
	}
	fmt.Println("Dropped index", indexName)
	return nil
}

// ReIndex recreates the index
func ReIndex(db *pgxpool.Pool, indexName string) error {
	_, err := db.Exec(context.Background(), fmt.Sprintf("REINDEX INDEX %s", indexName))
	if err != nil {
		return err
	}
	fmt.Println("Re-indexed", indexName)
	return nil
}

// ExplainAnalyze executes EXPLAIN ANALYZE for a given query and prints the query plan
func ExplainAnalyze(db *pgxpool.Pool, query string) error {
	rows, err := db.Query(context.Background(), "EXPLAIN ANALYZE "+query)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println("EXPLAIN ANALYZE (using index):")
	for rows.Next() {
		var plan string
		if err := rows.Scan(&plan); err != nil {
			return err
		}
		fmt.Println(plan)
	}

	return nil
}

func checkIndexExists(db *pgxpool.Pool, indexName string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM pg_indexes
			WHERE indexname = $1
		);
	`
	err := db.QueryRow(context.Background(), query, indexName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if index exists: %v", err)
	}
	return exists, nil
}

func CreateIndexIfNotExists(db *pgxpool.Pool, indexName, tableName string) error {
	exists, err := checkIndexExists(db, indexName)
	if err != nil {
		return err
	}

	if exists {
		fmt.Printf("Index '%s' already exists.\n", indexName)
		return nil
	}

	// 如果索引不存在，則創建它
	_, err = db.Exec(context.Background(), fmt.Sprintf("CREATE INDEX %s ON %s (customer_id, order_date)", indexName, tableName))
	if err != nil {
		return fmt.Errorf("failed to create index: %v", err)
	}
	fmt.Printf("Index '%s' created successfully.\n", indexName)
	return nil
}
