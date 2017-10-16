package main

import (
	"net/http"
	"time"
)

// Validation
type IValidable interface {
	Validate(r *http.Request) error
}

type ValidationError struct {
	Message  string
	methodId string
}

type ValidatedEntity struct {
	Entity interface{}
	Errors []ValidationError
}

var validators = []func(*Rec) *ValidationError{
	validateId,
	validateLinija,
	validateKm,
	validateM,
	validateSkodas,
	validateSuvirino,
	validateOperatorius,
	validateAparatas,
	validateTData,
	validateKelintas,
	validatePirmasId,
	validateNepirmasId,
	validateIesmeSiule,
	validateNeiesmeSiule,
	validateKelias8Pk,
	validateKeliasNe8Pk}

func (r Rec) Validate() ValidatedEntity {
	var verrs = []ValidationError{}

	for _, v := range validators {
		if verr := v(&r); verr != nil {
			verrs = append(verrs, *verr)
		}
	}

	return ValidatedEntity{Entity: r, Errors: verrs}
}

var (
	linijos     = []string{"01", "17", "22", "23", "24", "46", "47", "48", "49", "50", "51", "52", "86", "87", "94", "95", "96"}
	operatoriai = []string{"402", "407", "410", "419", "421", "422", "425", "426", "427", "428", "432", "435", "436", "437"}
	aparatai    = []string{"806", "807", "830", "807"}
	skodai      = []string{"06.3", "06.4"}
	suvirino    = []string{"IF4", "gamykla", "IF3", "GTC", "ŠP", "VitrasS"}
)

func validateId(r *Rec) *ValidationError {
	var ve ValidationError
	// logic
	// ve = ValidationError{"invalid Id", "ValidateId"}
	return &ve
}

func validateLinija(r *Rec) *ValidationError {
	var ve ValidationError
	if !inSlice(r.Linija, linijos) {
		ve = ValidationError{"invalid linija", "validateLinija"}
	}
	return &ve
}

func validateKm(r *Rec) *ValidationError {
	var ve ValidationError
	if r.Km < 1 {
		ve = ValidationError{"invalid km", "validateKm"}
	}
	return &ve
}

func validateM(r *Rec) *ValidationError {
	var ve ValidationError
	if r.M < 1 {
		ve = ValidationError{"invalid m", "validateM"}
	}
	return &ve
}

func validateSkodas(r *Rec) *ValidationError {
	var ve ValidationError
	if !inSlice(r.Skodas, skodai) {
		ve = ValidationError{"invalid sąlyginis kodas", "validateSkodas"}
	}
	return &ve
}

func validateSuvirino(r *Rec) *ValidationError {
	var ve ValidationError
	if !inSlice(r.Suvirino, suvirino) {
		ve = ValidationError{"invalid suvirinusi įmonė", "validateSuvirino"}
	}
	return &ve
}

func validateOperatorius(r *Rec) *ValidationError {
	var ve ValidationError
	if !inSlice(r.Operatorius, operatoriai) {
		ve = ValidationError{"invalid operatorius", "validateOperatorius"}
	}
	return &ve
}

func validateAparatas(r *Rec) *ValidationError {
	var ve ValidationError
	if !inSlice(r.Aparatas, aparatai) {
		ve = ValidationError{"invalid aparatas", "validateAparatas"}
	}
	return &ve
}

func validateTData(r *Rec) *ValidationError {
	const allowedDays = 15
	var ve ValidationError
	// time negali būti ateityje ir negali būti seniau nei prieš allowed days
	if r.TData.After(time.Now()) || r.TData.Sub(time.Now()).Hours() > allowedDays*24 {
		ve = ValidationError{"invalid tikrinimo data", "validateTData"}
	}
	return &ve
}

func validateKelintas(r *Rec) *ValidationError {
	var ve ValidationError
	if r.Kelintas > 4 {
		ve = ValidationError{"invalid kelintas", "validateKelintas"}
	}
	return &ve
}

// pirmas tikrinimas neturi Id
func validatePirmasId(r *Rec) *ValidationError {
	var ve ValidationError
	if r.Kelintas == 1 && r.Id > 0 {
		ve = ValidationError{"pirmasis tikrinimas negali turėti Id", "validatePirmasId"}
	}
	return &ve
}

// nepirmas tikrinimas turi Id
func validateNepirmasId(r *Rec) *ValidationError {
	var ve ValidationError
	if r.Kelintas <= 4 && r.Kelintas != 1 && r.Id == 0 {
		ve = ValidationError{"nepirmasis tikrinimas turi turėti Id", "validateNepirmasId"}
	}
	return &ve
}

// iešme neturi siūlės
func validateIesmeSiule(r *Rec) *ValidationError {
	var ve ValidationError
	if (r.Kelias == 8 || r.Kelias == 9) && r.Siule.Valid {
		ve = ValidationError{"iešme neturi būti nurodyta siūlė", "validateIesmeSiule"}
	}
	return &ve
}

// neiešme turi būti siūlė
func validateNeiesmeSiule(r *Rec) *ValidationError {
	var ve ValidationError
	if r.Kelias != 8 && r.Kelias != 9 && !r.Siule.Valid {
		ve = ValidationError{"neiešme turi būti nurodyta siūlė", "validateNeiesmeSiule"}
	}
	return &ve
}

// kelias 8 - pk 0
func validateKelias8Pk(r *Rec) *ValidationError {
	var ve ValidationError
	if r.Kelias == 8 && r.Pk > 0 {
		ve = ValidationError{"didelės stoties iešme pk turi būti 0", "validateKelias8Pk"}
	}
	return &ve
}

// kelias 8 - pk 0
func validateKeliasNe8Pk(r *Rec) *ValidationError {
	var ve ValidationError
	if r.Kelias != 8 && r.Pk == 0 {
		ve = ValidationError{"pk turi būti > 0", "validateKeliasNe8Pk"}
	}
	return &ve
}

// kelias 8 - pk 00

func inSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
