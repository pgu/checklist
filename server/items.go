package server

import (
	"encoding/json"
	"fmt"
	// "log"
	"net/http"
	"strconv"

	"appengine"
	"appengine/datastore"
	"github.com/gorilla/mux"
)

const ItemsPrefix = "/items/"

func addApiItems(r *mux.Router) {
	r.HandleFunc(ItemsPrefix, wrapJsonHandler(ListItems)).Methods("GET")
	r.HandleFunc(ItemsPrefix, wrapJsonHandler(NewItem)).Methods("POST", "OPTIONS")
	r.HandleFunc(ItemsPrefix+"{id}", wrapJsonHandler(GetItem)).Methods("GET")
	r.HandleFunc(ItemsPrefix+"{id}", wrapJsonHandler(UpdateItem)).Methods("PUT")
	r.HandleFunc(ItemsPrefix+"{id}", wrapJsonHandler(DeleteItem)).Methods("DELETE", "OPTIONS")
}

type Item struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	IsDone bool   `json:"isDone"`
}

func newDefaultItem() *Item {
	return &Item{0, "", false}
}

// req: GET /items/
// resp: 200
func ListItems(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	var items []Item
	if _, err := datastore.NewQuery("Item").GetAll(c, &items); err != nil {
		return err
	}
	if items == nil {
		items = []Item{}
	}

	return json.NewEncoder(w).Encode(items)
}

// req: POST /items/ {"Title": "Buy bread"}
// resp: 201
func NewItem(w http.ResponseWriter, r *http.Request) error {
	req := struct{ Title string }{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequest{err}
	}
	if req.Title == "" {
		return fmt.Errorf("Empty title!")
	}

	newItem := newDefaultItem()
	newItem.Title = req.Title

	c := appengine.NewContext(r)
	itemKey, e := datastore.Put(c, datastore.NewIncompleteKey(c, "Item", nil), newItem)
	if e != nil {
		return e
	}

	newItem.ID = itemKey.IntID()
	_, e = datastore.Put(c, itemKey, newItem)
	// log.Println("newItem.ID -> ", newItem.ID)
	if e != nil {
		return e
	}

	newUrl := r.URL.Path + strconv.FormatInt(newItem.ID, 10)
	w.Header().Set("Location", newUrl)
	w.WriteHeader(http.StatusCreated)
	return nil
}

// req: /items/99
// resp: 200
func GetItem(w http.ResponseWriter, r *http.Request) error {
	itemId, err := parseID(r)
	if err != nil {
		return badRequest{err}
	}

	c := appengine.NewContext(r)
	itemKey := datastore.NewKey(c, "Item", "", itemId, nil)

	var item Item
	err = datastore.Get(c, itemKey, &item)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(item)
}

// req: /items/99 {ID:99, Title:"title",...}
// resp: 204
func UpdateItem(w http.ResponseWriter, r *http.Request) error {
	itemId, err := parseID(r)
	// log.Println("itemId -> ", itemId)
	if err != nil {
		return badRequest{err}
	}
	var itemUI Item
	if err := json.NewDecoder(r.Body).Decode(&itemUI); err != nil {
		return badRequest{err}
	}
	if itemUI.ID != itemId {
		return badRequest{fmt.Errorf("Inconsistent item IDs")}
	}

	c := appengine.NewContext(r)
	itemKey := datastore.NewKey(c, "Item", "", itemId, nil)

	newItem := &Item{itemId, itemUI.Title, itemUI.IsDone}
	_, err = datastore.Put(c, itemKey, newItem)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// req: /items/99
// resp: 204
func DeleteItem(w http.ResponseWriter, r *http.Request) error {
	itemId, err := parseID(r)
	if err != nil {
		return badRequest{err}
	}

	c := appengine.NewContext(r)
	itemKey := datastore.NewKey(c, "Item", "", itemId, nil)

	err = datastore.Delete(c, itemKey)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
