module register

go 1.23.1

require PAPS-LAB/studiodb v0.0.0

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
)

replace PAPS-LAB/studiodb => ../studiodb
