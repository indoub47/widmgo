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
	validateKelias,
	validateLinija,
	validateKm,
	validatePk,
	validateM,
	validateSiule,
	validateSkodas,
	validateSuvirino,
	validateOperatorius,
	validateAparatas,
	validateTData,
	validateKelintas,
	validatePirmasId,
	validateNepirmasId,
	validatePirmasSuvirino,
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
	aparatai    = []string{"806", "807", "830", "831"}
	skodai      = []string{"06.3", "06.4"}
	suvirino    = []string{"IF4", "gamykla", "IF3", "GTC", "ŠP", "VitrasS"}
)

func validateId(r *Rec) *ValidationError {
	if !isIdValid(r) {
		return &ValidationError{"invalid Id", "ValidateId"}
	}
	return nil
}

func validateKelias(r *Rec) *ValidationError {
	if !isKeliasValid(r) {
		return &ValidationError{"invalid Kelias", "ValidateKelias"}
	}
	return nil
}

func validateLinija(r *Rec) *ValidationError {
	if !inSlice(r.Linija, linijos) {
		return &ValidationError{"invalid linija", "validateLinija"}
	}
	return nil
}

func validateKm(r *Rec) *ValidationError {
	if r.Km <= 0 {
		return &ValidationError{"invalid km", "validateKm"}
	}
	return nil
}

func validatePk(r *Rec) *ValidationError {
	if !isPkValid(r) {
		return &ValidationError{"invalid pk", "validatePk"}
	}
	return nil
}

func validateM(r *Rec) *ValidationError {
	if r.M < 1 {
		return &ValidationError{"invalid m", "validateM"}
	}
	return nil
}

func validateSiule(r *Rec) *ValidationError {
	if !isSiuleValid(r) {
		return &ValidationError{"invalid siūlė", "validateSiule"}
	}
	return nil
}

func validateSkodas(r *Rec) *ValidationError {
	if !inSlice(r.Skodas, skodai) {
		return &ValidationError{"invalid sąlyginis kodas", "validateSkodas"}
	}
	return nil
}

func validateSuvirino(r *Rec) *ValidationError {
	// gali būti tuščias arba iš leistinų reikšmių
	if !isSuvirinoValid(r) {
		return &ValidationError{"invalid suvirinusi įmonė", "validateSuvirino"}
	}
	return nil
}

func validateOperatorius(r *Rec) *ValidationError {
	if !inSlice(r.Operatorius, operatoriai) {
		return &ValidationError{"invalid operatorius", "validateOperatorius"}
	}
	return nil
}

func validateAparatas(r *Rec) *ValidationError {
	if !inSlice(r.Aparatas, aparatai) {
		return &ValidationError{"invalid aparatas", "validateAparatas"}
	}
	return nil
}

func validateTData(r *Rec) *ValidationError {
	const allowedDays = 15
	// time negali būti ateityje ir negali būti seniau nei prieš allowed days
	if r.TData.After(time.Now()) || r.TData.Sub(time.Now()).Hours() > allowedDays*24 {
		return &ValidationError{"invalid tikrinimo data", "validateTData"}
	}
	return nil
}

func validateKelintas(r *Rec) *ValidationError {
	if !isKelintasValid(r) {
		return &ValidationError{"invalid kelintas", "validateKelintas"}
	}
	return nil
}

// pirmas tikrinimas neturi Id
func validatePirmasId(r *Rec) *ValidationError {
	if isKelintasValid(r) && isIdValid(r) && r.Kelintas == 1 && r.ID.Valid {
		return &ValidationError{"pirmasis tikrinimas negali turėti Id", "validatePirmasId"}
	}
	return nil
}

// nepirmas tikrinimas turi Id
func validateNepirmasId(r *Rec) *ValidationError {
	if isKelintasValid(r) && isIdValid(r) && r.Kelintas != 1 && !r.ID.Valid {
		return &ValidationError{"nepirmasis tikrinimas turi turėti Id", "validateNepirmasId"}
	}
	return nil
}

// pirmas tikrinimas turi suvirino
func validatePirmasSuvirino(r *Rec) *ValidationError {
	if isKelintasValid(r) && isSuvirinoValid(r) && r.Kelintas == 1 && !r.Suvirino.Valid {
		return &ValidationError{"pirmam tikrinimui turi būti nurodyta, kas virino", "validatePirmasSuvirino"}
	}
	return nil
}

// iešme neturi siūlės
func validateIesmeSiule(r *Rec) *ValidationError {
	if isKeliasValid(r) && isSiuleValid(r) && (r.Kelias == 8 || r.Kelias == 9) && r.Siule.Valid {
		return &ValidationError{"iešme neturi būti nurodyta siūlė", "validateIesmeSiule"}
	}
	return nil
}

// neiešme turi būti siūlė
func validateNeiesmeSiule(r *Rec) *ValidationError {
	if isKeliasValid(r) && isSiuleValid(r) && r.Kelias != 8 && r.Kelias != 9 && !r.Siule.Valid {
		return &ValidationError{"neiešme turi būti nurodyta siūlė", "validateNeiesmeSiule"}
	}
	return nil
}

// kelias 8 - pk 0
func validateKelias8Pk(r *Rec) *ValidationError {
	if isKeliasValid(r) && isPkValid(r) && r.Kelias == 8 && r.Pk > 0 {
		return &ValidationError{"didelės stoties iešme pk turi būti 0", "validateKelias8Pk"}
	}
	return nil
}

// kelias 8 - pk 0
func validateKeliasNe8Pk(r *Rec) *ValidationError {
	if isKeliasValid(r) && isPkValid(r) && r.Kelias != 8 && r.Pk == 0 {
		return &ValidationError{"pk turi būti > 0", "validateKeliasNe8Pk"}
	}
	return nil
}

func isIdValid(r *Rec) bool {
	return !r.ID.Valid || r.ID.Int64 > 0
}

func isKeliasValid(r *Rec) bool {
	return r.Kelias > 0
}

func isPkValid(r *Rec) bool {
	return r.Pk >= 0
}

func isSiuleValid(r *Rec) bool {
	return !r.Siule.Valid || r.Siule.Int64 >= 0
}

func isKelintasValid(r *Rec) bool {
	return r.Kelintas >= 0 && r.Kelintas <= 4
}

func isSuvirinoValid(r *Rec) bool {
	return !r.Suvirino.Valid || inSlice(r.Suvirino.String, suvirino)
}

func inSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
