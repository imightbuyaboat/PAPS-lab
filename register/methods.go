package register

import (
	"PAPS-LAB/studiodb"
	"strconv"
)

func NewRegister(db *studiodb.DB) *Register {
	return &Register{db}
}

func (r *Register) Insert(i Item) error {
	query := "INSERT INTO register (organization, city, phone) VALUES ($1, $2, $3)"

	_, err := r.Exec(query, i.Organization, i.City, i.Phone)
	return err
}

func (r *Register) SelectAll() ([]Item, error) {
	query := "SELECT * FROM register"
	rows, err := r.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Items := []Item{}
	for rows.Next() {
		i := Item{}
		err := rows.Scan(&i.Id, &i.Organization, &i.City, &i.Phone)
		if err != nil {
			return nil, err
		}
		Items = append(Items, i)
	}
	return Items, nil
}

func (r *Register) SelectAny(i Item) ([]Item, error) {
	query := "SELECT * FROM register where 1=1"
	var args []interface{}

	if i.Organization != "" {
		query += " and organization = $"
		query += strconv.Itoa(len(args) + 1)
		args = append(args, i.Organization)
	}
	if i.City != "" {
		query += " and city = $"
		query += strconv.Itoa(len(args) + 1)
		args = append(args, i.City)
	}
	if i.Phone != "" {
		query += " and phone = $"
		query += strconv.Itoa(len(args) + 1)
		args = append(args, i.Phone)
	}

	rows, err := r.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Items := []Item{}
	for rows.Next() {
		i := Item{}
		err := rows.Scan(&i.Id, &i.Organization, &i.City, &i.Phone)
		if err != nil {
			return nil, err
		}
		Items = append(Items, i)
	}
	return Items, nil
}

func (r *Register) Delete(id int) error {
	_, err := r.Exec("DELETE FROM register WHERE id = $1", id)
	return err
}
