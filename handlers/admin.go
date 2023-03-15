package handlers

import (
	"log"
	"redbudway-api/database"
	"redbudway-api/internal"
	"redbudway-api/restapi/operations"
	"time"

	"github.com/go-openapi/runtime/middleware"
)

func PostAdminHandler(params operations.PostAdminParams) middleware.Responder {
	admin := params.Admin

	payload := operations.PostAdminCreatedBody{Created: false}
	response := operations.NewPostAdminCreated().WithPayload(&payload)

	if admin.MasterPass.String() == "MerCedEsAmgGt22$" {
		adminID, err := internal.GenerateUUID()
		if err != nil {
			return response
		}

		db := database.GetConnection()
		stmt, err := db.Prepare("INSERT INTO admin_account (adminId, user, password) VALUES (?, ?, ?)")
		if err != nil {
			return response
		}
		defer stmt.Close()

		passwordHash, err := internal.HashPassword(admin.Password.String())
		if err != nil {
			return response
		}

		results, err := stmt.Exec(adminID.String(), admin.User, passwordHash)
		if err != nil {
			return response
		}

		rowsAffected, err := results.RowsAffected()
		if err != nil {
			return response
		}

		if rowsAffected == 1 {
			payload.Created = true
			payload.AdminID = adminID.String()
			accessToken, err := internal.GenerateToken(adminID.String(), "admin", "access", time.Minute*15)
			if err != nil {
				log.Printf("Failed to generate JWT, %s", err)
				return response
			}
			payload.AccessToken = accessToken

			refreshToken, err := internal.GenerateToken(adminID.String(), "admin", "refresh", time.Minute*20)
			if err != nil {
				log.Printf("Failed to generate JWT, %s", err)
				return response
			}
			payload.RefreshToken = refreshToken

			saved, err := database.SaveAdminTokens(adminID.String(), refreshToken, accessToken)
			if err != nil {
				log.Printf("Failed to save admin tokens, %s", err)
				return response
			}
			if !saved {
				log.Printf("No issues, but failed to save admin")
			}

			response.SetPayload(&payload)
		}
	}

	return response
}

func GetAdminAdminIDTradespeopleHandler(params operations.GetAdminAdminIDTradespeopleParams, principal interface{}) middleware.Responder {

	tradespeople := []*operations.GetAdminAdminIDTradespeopleOKBodyItems0{}
	response := operations.NewGetAdminAdminIDTradespeopleOK().WithPayload(tradespeople)

	db := database.GetConnection()

	stmt, err := db.Prepare("SELECT tradespersonId, name, email, image FROM tradesperson_account")
	if err != nil {
		return response
	}
	defer stmt.Close()

	var tradespersonID, name, email, image string
	rows, err := stmt.Query()
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response
	}

	for rows.Next() {
		if err := rows.Scan(&tradespersonID, &name, &email, &image); err != nil {
			return response
		}
		tradesperson := operations.GetAdminAdminIDTradespeopleOKBodyItems0{}
		tradesperson.TradespersonID = tradespersonID
		tradesperson.Name = name
		tradesperson.Email = email
		tradesperson.Image = image
		tradespeople = append(tradespeople, &tradesperson)
	}
	response.SetPayload(tradespeople)

	return response
}
