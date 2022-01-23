package user

const (
	// CreateUserQuery - declare user create query.
	CreateUserQuery = `INSERT INTO users
		(user_id, name, email, password_hash, roles, date_created, date_updated)
	VALUES
		(:user_id, :name, :email, :password_hash, :roles, :date_created, :date_updated)`

	// UpdateUserQuery - declare user update query.
	UpdateUserQuery = `UPDATE 
							users 
						SET
							"name" = :name,
							"email" = :email,
							"roles" = :roles,
							"password_hash" = :password_hash,
							"date_updated" = :date_updated
						WHERE
							user_id = :user_id`

	// DeleteUserQuery - declare user delete query.
	DeleteUserQuery = `
	DELETE FROM
		users
	WHERE
		user_id = :user_id`

	// ListUserQuery - declare user list query.
	ListUserQuery = `
	SELECT
		*
	FROM
		users
	ORDER BY
		user_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	// IDUserQuery - declare user ID query.
	IDUserQuery = `
	SELECT
		*
	FROM
		users
	WHERE 
		user_id = :user_id`

	// EmailUserQuery - declare user Email query.
	EmailUserQuery = `
	SELECT
		*
	FROM
		users
	WHERE
		email = :email`
)
