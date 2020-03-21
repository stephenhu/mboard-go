package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/madsportslab/glbs"
)

const (

	EXT_DOT     = "."
	EXT_AVI     = "avi"
	EXT_JPG			= "jpg"
	EXT_JPEG    = "jpeg"
	EXT_M4A     = "m4a"
	EXT_MP3     = "mp3"
	EXT_MP4			= "mp4"
	EXT_MOV     = "mov"
	EXT_PNG			= "png"
	EXT_AAC     = "aac"

)

const (
	MEDIA_AUDIO 	= "AUDIO"
	MEDIA_PHOTO   = "PHOTO"
	MEDIA_VIDEO		= "VIDEO"
	MEDIA_UNKNOWN = "UNKNOWN"
)

const (

	MediaCreate = "INSERT into media(key, meta, tags) " +
	  "VALUES($1, $2, $3)"

  MediaDelete = "DELETE from media WHERE id=?"
	
	MediaGet = "SELECT " +
	  "id, game_id, key, meta, tags, created, updated " +
		"FROM media " +
		"WHERE id=?"

	MediaGetAll = "SELECT " +
	  "id, game_id, key, meta, tags, created, updated " + 
		"FROM media " +
		"ORDER BY created DESC"

)

type MediaMeta struct {
	Size		int64			`json:"size"`
	Name    string    `json:"name"`
	Ext     string    `json:"ext"`
}

type Media struct {
	ID					int   				`json:"id"`
	GameID      sql.NullInt64 `json:"gameId"`
	Key					string				`json:"key"`
	Meta  			*MediaMeta		`json:"meta"`
	Tag         string        `json:"tag"`
	Created     string				`json:"created"`
	Updated     string				`json:"updated"`
}

var SUPPORTED_VIDEO = []string{EXT_AVI, EXT_MP4, EXT_MOV}
var SUPPORTED_AUDIO = []string{EXT_MP3, EXT_AAC}
var SUPPORTED_PHOTO = []string{EXT_JPEG, EXT_JPG, EXT_PNG}


func findExt(exts []string, target string) bool {

	for _, ext := range exts {
		if EXT_DOT + ext == strings.ToLower(target) {
			return true
		}
	}

	return false

} // findExt

func getTag(ext string) string {
	
	if findExt(SUPPORTED_VIDEO, ext) {
		return MEDIA_VIDEO
	} else if findExt(SUPPORTED_AUDIO, ext) {
		return MEDIA_AUDIO
	} else if findExt(SUPPORTED_PHOTO, ext) {
		return MEDIA_PHOTO
	} else {
		return MEDIA_UNKNOWN
	}

} // getTag

func getMeta(filename string, size int64) (*MediaMeta, error) {

	m := MediaMeta{}

	m.Size 	= size
	m.Name 	= filename
	m.Ext		= filepath.Ext(filename)

	return &m, nil
	
} // getMeta

func createMedia(key string, filename string, size int64) {

	meta, err := getMeta(filename, size)

	if err != nil {
		log.Println(err)
	} else {

		tag := getTag(meta.Ext)

		j, err := json.Marshal(meta)

		if err != nil {
			log.Println(err)
		} else {
	
			_, err := data.Exec(
				MediaCreate, key, j, tag,
			)
		
			if err != nil {
				log.Println(err)
			}
		
		}

	}

} // createMedia


func removeMedia() {

} // removeMedia


func getMediaList() []Media {

	rows, err := data.Query(MediaGetAll)
	
	if err != nil {

		log.Println("getMediaList(): ", err)
		return nil

	}

	defer rows.Close()

	all := []Media{}

	for rows.Next() {

		m := Media{}

		jstr := ""

		err := rows.Scan(&m.ID, &m.GameID, &m.Key, &jstr, &m.Tag, &m.Created,
			&m.Updated)

		if err == sql.ErrNoRows || err != nil {
			
			log.Println("getMediaList(): ", err)
			return nil

		}

		mm := MediaMeta{}

		errJson := json.Unmarshal([]byte(jstr), &mm)

		if errJson != nil {
			log.Println("getMediaList(): ", errJson)
		} else {

			m.Meta = &mm
	
			all = append(all, m)
	
		}


	}

	return all

} // getMediaList


func mediaHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
	case http.MethodPost:

		err := r.ParseMultipartForm(200000)

		if err != nil {
			log.Println(err)
		}

    form := r.MultipartForm

		media := form.File["media"]

		for _, m := range media {

			file, err := m.Open()

			defer file.Close()

			if err != nil {
				log.Println(err)
			} else {

				glbs.SetNamespace("blobs")

				k := glbs.Put(file)

				log.Printf("%s uploaded successfully, %s", m.Filename , *k)

				createMedia(*k, m.Filename, m.Size)

			}

		}

	case http.MethodGet:
		
		// search database for all media

		all := getMediaList()

		j, err := json.Marshal(all)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.Write(j)
		}
		
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // mediaHandler
