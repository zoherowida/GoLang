package controllers

import (
	"api/database"
	"api/models"
	"api/repository"
	"api/repository/crud"
	"api/responses"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/mux"
)

func SearchTweets(w http.ResponseWriter, r *http.Request) {

	//twitters := models.Twitter{}

	r.ParseForm()
	search_field := r.FormValue("search")
	if search_field == "" {
		responses.JSON(w, http.StatusOK, struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
			Data    string `json:"data"`
		}{
			Status:  http.StatusUnprocessableEntity,
			Message: "error",
			Data:    "Search parameter is required",
		})

		return
	}

	config := oauth1.NewConfig("ahCowaVjdVvO5lmJGMpigchRR", "o7dy5avGMI5AJutFPWjYHAeTAw5kVR8YBMt00OKXKvnUzBUvak")
	token := oauth1.NewToken("630413671-toB4t1du9TkC3NVwfdxFzZpZKJ6kDqiIGcBIG2Bx", "fnTfmLeB1uiEheUPUHCEcFchcylJbBxLmvQHdrFF5wZl7")
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	search, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: search_field,
		Count: 50,
	})
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}
	db, err := database.Connect()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	for _, s := range search.Statuses {

		tweet := string(s.Text)

		query := &models.Twitter{Tweet: tweet}

		db.Create(&query)
	}

	responses.JSON(w, http.StatusOK, struct {
		Status  int             `json:"status"`
		Message string          `json:"message"`
		Data    *twitter.Search `json:"data"`
	}{
		Status:  http.StatusOK,
		Message: "success",
		Data:    search,
	})
}

func GetAllTweet(w http.ResponseWriter, r *http.Request) {

	db, err := database.Connect()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		fmt.Println(4)
		return
	}
	defer db.Close()

	// Start Get Number of Page And convert To Uint64
	r.ParseForm()
	page := r.FormValue("page")

	u64, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	wd := uint(u64)
	// End Get Number of Page And convert To Uint64

	repo := crud.NewRepositoryTwittersCRUD(db)
	func(twitterRepository repository.TwitterRepository) {
		offset := wd
		twitters, err := twitterRepository.FindAll(uint32(offset))

		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		responses.JSON(w, http.StatusOK, struct {
			Status  int              `json:"status"`
			Message string           `json:"message"`
			Data    []models.Twitter `json:"data"`
		}{
			Status:  http.StatusOK,
			Message: "success",
			Data:    twitters,
		})
	}(repo)

}

func CreateTweet(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	twitter := models.Twitter{}
	err = json.Unmarshal(body, &twitter)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = twitter.Validate("create")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		fmt.Println(4)

		return
	}
	defer db.Close()

	repo := crud.NewRepositoryTwittersCRUD(db)

	func(twitterRepository repository.TwitterRepository) {
		twitter, err := twitterRepository.Save(twitter)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, twitter.ID))
		responses.JSON(w, http.StatusOK, struct {
			Status  int            `json:"status"`
			Message string         `json:"message"`
			Data    models.Twitter `json:"data"`
		}{
			Status:  http.StatusOK,
			Message: "success",
			Data:    twitter,
		})

	}(repo)

}

func GetTweet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := crud.NewRepositoryTwittersCRUD(db)

	func(twitterRepository repository.TwitterRepository) {
		twitter, err := twitterRepository.FindByID(uint32(uid))
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
		responses.JSON(w, http.StatusOK, struct {
			Status  int            `json:"status"`
			Message string         `json:"message"`
			Data    models.Twitter `json:"data"`
		}{
			Status:  http.StatusOK,
			Message: "success",
			Data:    twitter,
		})

	}(repo)
}

func UpdateTweet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tweet := models.Twitter{}
	err = json.NewDecoder(r.Body).Decode(&tweet)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = tweet.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := crud.NewRepositoryTwittersCRUD(db)

	func(twitterRepository repository.TwitterRepository) {
		_, err := twitterRepository.Update(uint32(uid), tweet)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
		responses.JSON(w, http.StatusOK, struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
			Data    string `json:"data"`
		}{
			Status:  http.StatusOK,
			Message: "success",
			Data:    "",
		})

	}(repo)
}

func DeleteTweet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := crud.NewRepositoryTwittersCRUD(db)

	func(twitterRepository repository.TwitterRepository) {
		_, err = twitterRepository.Delete(uint32(uid))
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}

		w.Header().Set("Entity", fmt.Sprintf("%d", uid))
		responses.JSON(w, http.StatusOK, struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
			Data    string `json:"data"`
		}{
			Status:  http.StatusOK,
			Message: "success",
			Data:    "",
		})

	}(repo)
}
