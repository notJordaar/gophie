package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

// Props : The scraping engine Properties and description about the engine (e.g NetNaijaEngine)
type Props struct {
	Name        string
	BaseURL     *url.URL // The Base URL for the engine
	SearchURL   *url.URL // URL for searching
	ListURL     *url.URL // URL to return movie lists
	Description string
}

// PropsJSON : JSON structure of all downloadable movies
type PropsJSON struct {
	Props
	BaseURL   string
	SearchURL string
	ListURL   string
}

// MarshalJSON Props structure to return from api
func (p *Props) MarshalJSON() ([]byte, error) {
	props := PropsJSON{
		Props:     *p,
		BaseURL:   p.BaseURL.String(),
		SearchURL: p.SearchURL.String(),
		ListURL:   p.ListURL.String(),
	}

	return json.Marshal(props)
}

// Engine : interface for all engines
type Engine interface {
	Search(query string) SearchResult
	Scrape(mode string) ([]Movie, error)
	List(page int) SearchResult
	String() string
}

// Movie : the structure of all downloadable movies
type Movie struct {
	Index          int
	Title          string
	CoverPhotoLink string
	Description    string
	Size           string
	DownloadLink   *url.URL
	Year           int
	IsSeries       bool
	SDownloadLink  []*url.URL // Other links for downloads if movies is series
	UploadDate     string
	Source         string // The Engine From which it is gotten from
}

// MovieJSON : JSON structure of all downloadable movies
type MovieJSON struct {
	Movie
	DownloadLink  string
	SDownloadLink []string
}

func (m *Movie) String() string {
	return fmt.Sprintf("%s (%v)", m.Title, m.Year)
}

// MarshalJSON Json structure to return from api
func (m *Movie) MarshalJSON() ([]byte, error) {
	var sDownloadLink []string
	for _, link := range m.SDownloadLink {
		sDownloadLink = append(sDownloadLink, link.String())
	}

	movie := MovieJSON{
		Movie:         *m,
		DownloadLink:  m.DownloadLink.String(),
		SDownloadLink: sDownloadLink,
	}

	return json.Marshal(movie)

}

// SearchResult : the results of search from engine
type SearchResult struct {
	Query  string
	Movies []Movie
}

// Titles : Get a slice of the titles of movies
func (s *SearchResult) Titles() []string {
	var titles []string
	for _, movie := range s.Movies {
		titles = append(titles, movie.Title)
	}
	return titles
}

// GetMovieByTitle : Return a movie object from title passed
func (s *SearchResult) GetMovieByTitle(title string) (Movie, error) {
	for _, movie := range s.Movies {
		if movie.Title == title {
			return movie, nil
		}
	}
	return Movie{}, errors.New("Movie not Found")
}

// GetIndexFromTitle : return movieIndex from title
func (s *SearchResult) GetIndexFromTitle(title string) (int, error) {
	for index, movie := range s.Movies {
		if movie.Title == title {
			return index, nil
		}
	}
	return 0, errors.New("Movie not Found")
}

// GetEngines : Returns all the usable engines in the application
func GetEngines() map[string]Engine {
	engines := make(map[string]Engine)
	engines["netnaija"] = NewNetNaijaEngine()
	engines["fzmovies"] = NewFzEngine()
	return engines
}

// GetEngine : Return an engine
func GetEngine(engine string) (Engine, error) {
	e := GetEngines()[strings.ToLower(engine)]
	if e == nil {
		return nil, fmt.Errorf("Engine %s Does not exist", engine)
	}
	return e, nil
}

// Get the movie index context stored in Request
func getMovieIndexFromCtx(r *colly.Request) int {
	movieIndex, err := strconv.Atoi(r.Ctx.Get("movieIndex"))
	if err != nil {
		log.Fatal(err)
	}
	return movieIndex
}
