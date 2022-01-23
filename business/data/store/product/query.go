package product

const (
	// CreateProductQuery - declare product create query.
	CreateProductQuery =  `INSERT INTO products (product_id, user_id, name, cost, quantity, date_created, date_updated) VALUES (:product_id, :user_id, :name, :cost, :quantity, :date_created, :date_updated)`

	// UpdateProductQuery - declare product update query.
	UpdateProductQuery =  `
	UPDATE
		products
	SET
		"name" = :name,
		"cost" = :cost,
		"quantity" = :quantity,
		"date_updated" = :date_updated
	WHERE
		product_id = :product_id`

	// DeleteProductQuery - declare product delete query.
	DeleteProductQuery =  `
	DELETE FROM
		products
	WHERE
		product_id = :product_id`

	// ListProductQuery - declare product list query.
	ListProductQuery = `
	SELECT
		p.*,
		COALESCE(SUM(s.quantity) ,0) AS sold,
		COALESCE(SUM(s.paid), 0) AS revenue
	FROM
		products AS p
	LEFT JOIN
		sales AS s ON p.product_id = s.product_id
	GROUP BY
		p.product_id
	ORDER BY
		user_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	// IDProductQuery - declare product ID query.
	IDProductQuery = `
	SELECT
		p.*,
		COALESCE(SUM(s.quantity), 0) AS sold,
		COALESCE(SUM(s.paid), 0) AS revenue
	FROM
		products AS p
	LEFT JOIN
		sales AS s ON p.product_id = s.product_id
	WHERE
		p.product_id = :product_id
	GROUP BY
		p.product_id`

	// UserProductIDQuery - declare product user ID query.
	UserProductIDQuery =  `
	SELECT
		p.*,
		COALESCE(SUM(s.quantity), 0) AS sold,
		COALESCE(SUM(s.paid), 0) AS revenue
	FROM
		products AS p
	LEFT JOIN
		sales AS s ON p.product_id = s.product_id
	WHERE
		p.user_id = :user_id
	GROUP BY
		p.product_id`
)