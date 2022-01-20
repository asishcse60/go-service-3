package db

const (
	// UserCreateQuery - declare user create query.
	UserCreateQuery = `INSERT INTO users
		(user_id, name, email, password_hash, roles, date_created, date_updated)
	VALUES
		(:user_id, :name, :email, :password_hash, :roles, :date_created, :date_updated)`

	// UserUpdateQuery - declare user update query.
	UserUpdateQuery = `UPDATE 
							users 
						SET
							"name" = :name,
							"email" = :email,
							"roles" = :roles,
							"password_hash" = :password_hash,
							"date_updated" = :date_updated
						WHERE
							user_id = :user_id`

	// UserDeleteQuery - declare user delete query.
	UserDeleteQuery = `
	DELETE FROM
		users
	WHERE
		user_id = :user_id`

	// UserListQuery - declare user list query.
	UserListQuery = `
	SELECT
		*
	FROM
		users
	ORDER BY
		user_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	// UserIDQuery - declare user ID query.
	UserIDQuery = `
	SELECT
		*
	FROM
		users
	WHERE 
		user_id = :user_id`

	// UserEmailQuery - declare user Email query.
	UserEmailQuery = `
	SELECT
		*
	FROM
		users
	WHERE
		email = :email`
)
