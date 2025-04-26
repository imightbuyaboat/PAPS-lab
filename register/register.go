package register

import (
	bt "papslab/basic_types"
	"papslab/studiodb"
	"strconv"
)

type Register struct {
	*studiodb.DB
}

func NewRegister(db *studiodb.DB) *Register {
	return &Register{db}
}

func (r *Register) Insert(i bt.Item) error {
	var maxVirtualID int
	err := r.QueryRow("SELECT COALESCE(MAX(virtual_id), -1) FROM register").Scan(&maxVirtualID)
	if err != nil {
		return err
	}
	newVirtualID := maxVirtualID + 1

	query := "INSERT INTO register (organization, city, phone, virtual_id) VALUES ($1, $2, $3, $4)"
	_, err = r.Exec(query, i.Organization, i.City, i.Phone, newVirtualID)
	return err
}

func (r *Register) SelectAll() ([]bt.Item, error) {
	query := "SELECT organization, city, phone, virtual_id FROM register ORDER BY virtual_id"
	rows, err := r.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Items := []bt.Item{}
	for rows.Next() {
		i := bt.Item{}
		err := rows.Scan(&i.Organization, &i.City, &i.Phone, &i.Id)
		if err != nil {
			return nil, err
		}
		Items = append(Items, i)
	}
	return Items, nil
}

func (r *Register) SelectAny(i bt.Item) ([]bt.Item, error) {
	query := "SELECT organization, city, phone, virtual_id FROM register where 1=1"
	var args []interface{}

	if i.Organization != "" {
		query += " and organization = $" + strconv.Itoa(len(args)+1)
		args = append(args, i.Organization)
	}
	if i.City != "" {
		query += " and city = $" + strconv.Itoa(len(args)+1)
		args = append(args, i.City)
	}
	if i.Phone != "" {
		query += " and phone = $" + strconv.Itoa(len(args)+1)
		args = append(args, i.Phone)
	}
	query += " ORDER BY virtual_id"

	rows, err := r.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Items := []bt.Item{}
	for rows.Next() {
		i := bt.Item{}
		err := rows.Scan(&i.Organization, &i.City, &i.Phone, &i.Id)
		if err != nil {
			return nil, err
		}
		Items = append(Items, i)
	}
	return Items, nil
}

func (r *Register) Delete(virtualID int) error {
	var realID int
	err := r.QueryRow("SELECT id FROM register WHERE virtual_id = $1", virtualID).Scan(&realID)
	if err != nil {
		return err
	}

	_, err = r.Exec("DELETE FROM register WHERE id = $1", realID)
	if err != nil {
		return err
	}

	rows, err := r.Query("SELECT id FROM register ORDER BY virtual_id")
	if err != nil {
		return err
	}
	defer rows.Close()

	newVirtualID := 0
	for rows.Next() {
		var currentRealID int
		if err := rows.Scan(&currentRealID); err != nil {
			return err
		}
		_, err := r.Exec("UPDATE register SET virtual_id = $1 WHERE id = $2", newVirtualID, currentRealID)
		if err != nil {
			return err
		}
		newVirtualID++
	}
	return nil
}
