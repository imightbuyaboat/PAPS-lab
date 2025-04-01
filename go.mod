module PAPS-LAB

go 1.23.1

replace PAPS-LAB/studiodb => ./studiodb

require (
	PAPS-LAB/sessionmanager v0.0.0-00010101000000-000000000000
	PAPS-LAB/studiodb v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.1
)

require github.com/lib/pq v1.10.9 // indirect

replace PAPS-LAB/sessionmanager => ./sessionmanager
