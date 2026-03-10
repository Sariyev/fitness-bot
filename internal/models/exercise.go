package models

import "time"

type Exercise struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	Technique          string    `json:"technique"`
	CommonMistakes     string    `json:"common_mistakes"`
	EasierModification string    `json:"easier_modification"`
	HarderModification string    `json:"harder_modification"`
	RestSeconds        int       `json:"rest_seconds"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type WorkoutExercise struct {
	ID              int    `json:"id"`
	WorkoutID       int    `json:"workout_id"`
	ExerciseID      int    `json:"exercise_id"`
	Sets            int    `json:"sets"`
	Reps            string `json:"reps"`
	DurationSeconds int    `json:"duration_seconds"`
	SortOrder       int    `json:"sort_order"`
}
