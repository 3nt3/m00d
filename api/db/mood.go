package db

import "time"

func GetMoods() ([]Mood, error) {
	rows, err := db.Query("select id, mood, user_id, created_at from moods")
	if err != nil {
		return nil, err
	}

	moods := []Mood{}
	for rows.Next() {
		var mood Mood

		if err := rows.Scan(&mood.ID, &mood.Mood, &mood.UserID, &mood.CreatedAt); err != nil {
			return nil, err
		}

		moods = append(moods, mood)
	}

	return moods, nil
}

func NewMood(mood Mood) (Mood, error) {
	row := db.QueryRow("insert into moods (mood, user_id, created_at) values ($1, $2, $3) returning id", mood.Mood, mood.UserID, time.Now())

	if row.Err() != nil {
		return mood, row.Err()
	}

	var id int
	var newMood Mood

	err := row.Scan(&id)
	if err != nil {
		return newMood, err
	}

	newMood, err = GetMoodByID(id)
	return newMood, err
}

func GetMoodByID(id int) (Mood, error) {
	row := db.QueryRow("select id, mood, user_id, created_at from moods where id = $1", id)
	var mood Mood

	if row.Err() != nil {
		return mood, row.Err()
	}

	err := row.Scan(&mood.ID, &mood.Mood, &mood.UserID, &mood.CreatedAt)
	return mood, err
}
