package database

import "database/sql"

func ClearCustomerTokens(customerID string) (bool, error) {
	stmt, err := db.Prepare("UPDATE customer_token SET refresh_token='', access_token='' WHERE customerId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(customerID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func ClearTradespersonTokens(tradespersonID string) (bool, error) {
	stmt, err := db.Prepare("UPDATE tradesperson_token SET refresh_token='', access_token='' WHERE tradespersonId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(tradespersonID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func UpdateCustomerTokens(customerID, refreshToken, accessToken string) (bool, error) {
	stmt, err := db.Prepare("UPDATE customer_token SET refresh_token=?, access_token=? WHERE customerId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(refreshToken, accessToken, customerID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func UpdateTradespersonTokens(tradespersonID, refreshToken, accessToken string) (bool, error) {
	stmt, err := db.Prepare("UPDATE tradesperson_token SET refresh_token=?, access_token=? WHERE tradespersonId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(refreshToken, accessToken, tradespersonID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func SaveCustomerTokens(customerID, refreshToken, accessToken string) (bool, error) {
	saved := false

	stmt, err := db.Prepare("SELECT customerId FROM customer_token WHERE customerId=?")
	if err != nil {
		return saved, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(customerID)
	switch err = row.Scan(&customerID); err {
	case sql.ErrNoRows:
		stmt, err := db.Prepare("INSERT INTO customer_token (customerId, refresh_token, access_token) VALUES (?, ?, ?)")
		if err != nil {
			return saved, err
		}
		defer stmt.Close()

		results, err := stmt.Exec(customerID, refreshToken, accessToken)
		if err != nil {
			return saved, err
		}

		rowsAffected, err := results.RowsAffected()
		if err != nil {
			return saved, err
		}

		if rowsAffected == 1 {
			saved = true
		}
	case nil:
		stmt, err := db.Prepare("UPDATE customer_token SET refresh_token=?, access_token=? WHERE customerId=?")
		if err != nil {
			return saved, err
		}
		defer stmt.Close()

		results, err := stmt.Exec(refreshToken, accessToken, customerID)
		if err != nil {
			return saved, err
		}

		rowsAffected, err := results.RowsAffected()
		if err != nil {
			return saved, err
		}

		if rowsAffected == 1 {
			saved = true
		}
	default:
		return saved, err
	}

	return saved, nil
}

func SaveTradespersonTokens(tradespersonID, refreshToken, accessToken string) (bool, error) {
	saved := false

	stmt, err := db.Prepare("SELECT tradespersonId FROM tradesperson_token WHERE tradespersonId=?")
	if err != nil {
		return saved, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonID)

	switch err = row.Scan(&tradespersonID); err {
	case sql.ErrNoRows:
		stmt, err := db.Prepare("INSERT INTO tradesperson_token (tradespersonId, refresh_token, access_token) VALUES (?, ?, ?)")
		if err != nil {
			return saved, err
		}
		defer stmt.Close()

		results, err := stmt.Exec(tradespersonID, refreshToken, accessToken)
		if err != nil {
			return saved, err
		}

		rowsAffected, err := results.RowsAffected()
		if err != nil {
			return saved, err
		}

		if rowsAffected == 1 {
			saved = true
		}
	case nil:
		stmt, err := db.Prepare("UPDATE tradesperson_token SET refresh_token=?, access_token=? WHERE tradespersonId=?")
		if err != nil {
			return saved, err
		}
		defer stmt.Close()

		results, err := stmt.Exec(refreshToken, accessToken, tradespersonID)
		if err != nil {
			return saved, err
		}

		rowsAffected, err := results.RowsAffected()
		if err != nil {
			return saved, err
		}

		if rowsAffected == 1 {
			saved = true
		}
	default:
		return saved, err
	}

	return saved, nil
}

func CheckCustomerAccessToken(customerID, token string) (bool, error) {
	valid := false

	stmt, err := db.Prepare("SELECT access_token FROM customer_token WHERE customerId=?")
	if err != nil {
		return valid, err
	}
	defer stmt.Close()

	accessToken := ""
	if err := stmt.QueryRow(customerID).Scan(&accessToken); err != nil {
		return valid, err
	}

	if accessToken != "" {
		if accessToken == token {
			valid = true
		}
	}

	return valid, nil
}

func CheckCustomerRefreshToken(customerID, token string) (bool, error) {
	valid := false

	stmt, err := db.Prepare("SELECT refresh_token FROM customer_token WHERE customerId=?")
	if err != nil {
		return valid, err
	}
	defer stmt.Close()

	refreshToken := ""
	if err := stmt.QueryRow(customerID).Scan(&refreshToken); err != nil {
		return valid, err
	}

	if refreshToken != "" {
		if refreshToken == token {
			valid = true
		}
	}

	return valid, nil
}

func CheckTradespersonAccessToken(tradespersonID, token string) (bool, error) {
	valid := false

	stmt, err := db.Prepare("SELECT access_token FROM tradesperson_token WHERE tradespersonId=?")
	if err != nil {
		return valid, err
	}
	defer stmt.Close()

	accessToken := ""
	if err := stmt.QueryRow(tradespersonID).Scan(&accessToken); err != nil {
		return valid, err
	}

	if accessToken != "" {
		if accessToken == token {
			valid = true
		}
	}

	return valid, nil
}

func CheckTradespersonRefreshToken(tradespersonID, token string) (bool, error) {
	valid := false

	stmt, err := db.Prepare("SELECT refresh_token FROM tradesperson_token WHERE tradespersonId=?")
	if err != nil {
		return valid, err
	}
	defer stmt.Close()

	refreshToken := ""
	if err := stmt.QueryRow(tradespersonID).Scan(&refreshToken); err != nil {
		return valid, err
	}

	if refreshToken != "" {
		if refreshToken == token {
			valid = true
		}
	}

	return valid, nil
}
