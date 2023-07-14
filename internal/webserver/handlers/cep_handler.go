package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/gsouza97/go-multithreading-api/internal/dto"
	"github.com/gsouza97/go-multithreading-api/internal/entity"
)

type CepHandler struct{}

type Error struct {
	Message string `json:"message"`
}

func NewCepHandler() *CepHandler {
	return &CepHandler{}
}

func (handler *CepHandler) GetCep(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if cep == "" || !validateCep(cep) {
		e := Error{Message: "cep is empty or invalid"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(e)
		return
	}

	ch1 := make(chan interface{})
	ch2 := make(chan interface{})

	go findByViaCep(cep, ch1)
	go findByCdnCep(cep, ch2)

	select {
	case viaCep := <-ch1:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(viaCep)
		fmt.Println("ViaCep")
	case cdnCep := <-ch2:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cdnCep)
		fmt.Println("ViaCep")
	case <-time.After(1 * time.Second):
		e := Error{Message: "Exceeded time limit"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(e)
	}
}

func validateCep(cep string) bool {
	regexPatter := regexp.MustCompile(`^[0-9]{5}-?[0-9]{3}$`)
	return regexPatter.MatchString(cep)
}

func findByViaCep(cep string, ch chan<- interface{}) {
	url := "https://viacep.com.br/ws/" + cep + "/json/"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		ch <- entity.Cep{}
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		ch <- entity.Cep{}
		return
	}
	var c dto.ViaCepResponse
	err = json.Unmarshal(body, &c)
	if err != nil {
		fmt.Println(err)
		ch <- entity.Cep{}
		return
	}
	// fmt.Println(c)
	ch <- entity.MapViaCepResponseToCep(c)
}

func findByCdnCep(cep string, ch chan<- interface{}) {
	url := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		ch <- entity.Cep{}
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		ch <- entity.Cep{}
		return
	}
	var c dto.CdnCepResponse
	err = json.Unmarshal(body, &c)
	if err != nil {
		fmt.Println(err)
		ch <- entity.Cep{}
		return
	}

	// fmt.Println(c)

	ch <- entity.MapCdnCepResponseToCep(c)
}
