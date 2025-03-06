package models

const (
	ErrNoFilmsFound = "no films found"
)

//	type Person struct {
//		ID        int    `json:"id"`
//		FullName  string `json:"full_name"`
//		BirthDate string `json:"birth_date"`
//		FilmsID   []int  `json:"films"`
//	}
//
//	type Actor struct {
//		Person
//	}
//
//	type CrewMember struct {
//		Person
//	}
//
//	type Position struct {
//		ID       int    `json:"id"`
//		Name     string `json:"name"`
//		PersonID int    `json:"person_id"`
//		FilmID   int    `json:"film_id"`
//	}
//
//	type Role struct {
//		ID          int `json:"id"`
//		ActorID     int `json:"actor_id"`
//		CharacterID int `json:"character_id"`
//		FilmID      int `json:"film_id"`
//	}
//
//	type Character struct {
//		ID   int    `json:"id"`
//		Name string `json:"name"`
//	}
//
//	type Genre struct {
//		ID   int    `json:"id"`
//		Name string `json:"name"`
//	}
//
//	type Film struct {
//		ID          string          `json:"id"`
//		Title       string       `json:"title"`
//		Description string       `json:"description"`
//		ReleaseYear int          `json:"release_year"`
//		Genres      []Genre      `json:"genres"`
//		Country     string       `json:"country"`
//		Duration    int          `json:"duration"`
//		Budget      int          `json:"budget"`
//		BoxOffice   int          `json:"box_office"`
//		Actors      []Actor      `json:"actors"`
//		CrewMembers []CrewMember `json:"crew_members"`
//	}
type Film struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseYear int    `json:"release_year"`
	Country     string `json:"country"`
	Duration    int    `json:"duration"`
	Budget      int    `json:"budget"`
	BoxOffice   int    `json:"box_office"`
}

type FilmInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseYear int    `json:"release_year"`
	Country     string `json:"country"`
	Duration    int    `json:"duration"`
	Budget      int    `json:"budget"`
	BoxOffice   int    `json:"box_office"`
}
