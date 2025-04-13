module PAPS-LAB

go 1.23.1

require (
	PAPS-LAB/passwordmanager v0.0.0-00010101000000-000000000000
	PAPS-LAB/register v0.0.0-00010101000000-000000000000
	PAPS-LAB/sessionmanager v0.0.0-00010101000000-000000000000
	PAPS-LAB/studiodb v0.0.0
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1
)

replace PAPS-LAB/studiodb => ./studiodb

replace PAPS-LAB/passwordmanager => ./passwordmanager

replace PAPS-LAB/sessionmanager => ./sessionmanager

replace PAPS-LAB/register => ./register

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
)
