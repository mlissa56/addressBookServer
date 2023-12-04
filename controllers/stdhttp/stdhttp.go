package stdhttp

import (
    "addressBookServer/pkg"
	"addressBookServer/gate/psg"
	"addressBookServer/models/dto"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// Controller обрабатывает HTTP запросы для адресной книги.
type Controller struct {
	DB  *psg.Psg
}

// NewController создает новый Controller.
func NewController(db *psg.Psg) *Controller {
	return &Controller{db}
}

func WriteNotFound(w http.ResponseWriter) {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte(`{"code":404,"msg":"Not Found"}`))
}

func WriteInternalServerError(w http.ResponseWriter, err error) {
    w.WriteHeader(http.StatusInternalServerError)
    panic(err)
    w.Write([]byte(`{"code":500,"msg":"` + err.Error() + `"}`))
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log.Printf("%s %s", r.Method, r.URL)

    w.Header().Set("content-type", "application/json") 

    if r.Method == http.MethodGet {
        c.RecordsGet(w, r)
    } else if r.Method == http.MethodPost {
        c.RecordAdd(w, r)
    } else if r.Method == "PUT" {
        c.RecordUpdate(w, r)
    } else if r.Method == http.MethodDelete {
        c.RecordDeleteByPhone(w, r)
    } else {
        WriteNotFound(w) 
    } 
}

// RecordAdd обрабатывает HTTP запрос для добавления новой записи.
func (c *Controller) RecordAdd(w http.ResponseWriter, r *http.Request) {
    var rec dto.Record
    err := json.NewDecoder(r.Body).Decode(&rec)
    if err != nil {
        log.Println(err)
        WriteInternalServerError(w, err) 
        return
    }

    rec.Phone, err = pkg.PhoneNormalize(rec.Phone)
    if err != nil {
        log.Println(err)
        WriteInternalServerError(w, err) 
        return
    }
    
    id, err := c.DB.RecordAdd(rec)
    if err != nil {
        log.Println(err)
        WriteInternalServerError(w, err) 
        return
    }

    w.Write([]byte(`{"id": "` + string(id) + `"}`))
}

// RecordsGet обрабатывает HTTP запрос для получения записей на основе предоставленных полей Record.
func (c *Controller) RecordsGet(w http.ResponseWriter, r *http.Request) {
    var rec dto.Record
    err := json.NewDecoder(r.Body).Decode(&rec)
    if err != nil {
        log.Println(err)
        WriteInternalServerError(w, err) 
        return
    }
    
    matchedRecs, err := c.DB.RecordsGet(rec)
    if err != nil {
        log.Println(err)
        WriteInternalServerError(w, err) 
        return
    }

    var jsonb []byte
    jsonb = append(jsonb, []byte("[")...)
    for i, matchedRec := range(matchedRecs) {
        jsonbMatchedRec, err := json.Marshal(matchedRec)
        if err != nil {
            log.Println(err)
            WriteInternalServerError(w, err) 
            return
        }

        jsonb = append(jsonb, jsonbMatchedRec...)
        if i != len(matchedRecs) - 1 {
            jsonb = append(jsonb, []byte(",")...)
        }
    }
    jsonb = append(jsonb, []byte("]")...)

    w.Write(jsonb)  
}

// RecordUpdate обрабатывает HTTP запрос для обновления записи.
func (c *Controller) RecordUpdate(w http.ResponseWriter, r *http.Request) {
    var rec dto.Record
    err := json.NewDecoder(r.Body).Decode(&rec)
    if err != nil {
        log.Println(err)
        WriteInternalServerError(w, err) 
        return
    }

    rec.Phone, err = pkg.PhoneNormalize(rec.Phone)
    if err != nil {
        log.Println(err)
        WriteInternalServerError(w, err) 
        return
    }
    
    qerr := c.DB.RecordUpdate(rec)
    if qerr != nil {
        log.Println(qerr)
        WriteInternalServerError(w, qerr) 
        return
    }
}

// RecordDeleteByPhone обрабатывает HTTP запрос для удаления записи по номеру телефона.
func (c *Controller) RecordDeleteByPhone(w http.ResponseWriter, r *http.Request) {
    phone := strings.Split(r.URL.Path, "/")[1:][1] 

    err := c.DB.RecordDeleteByPhone(phone)
    if err != nil {
        log.Println(err) 
        WriteInternalServerError(w, err)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("{}"))
}
