package entity

import (
	"fmt"

	"github.com/gsouza97/go-multithreading-api/internal/dto"
)

func MapCdnCepResponseToCep(response dto.CdnCepResponse) Cep {
	return Cep{
		Cep:        response.Code,
		Logradouro: response.Address,
		Bairro:     response.District,
		Localidade: response.City,
		Uf:         response.State,
	}
}

func MapViaCepResponseToCep(response dto.ViaCepResponse) Cep {
	fmt.Println(response)
	return Cep{
		Cep:        response.Cep,
		Logradouro: response.Logradouro,
		Bairro:     response.Bairro,
		Localidade: response.Localidade,
		Uf:         response.Uf,
	}
}
