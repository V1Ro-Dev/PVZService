package forms

import (
	"time"

	"github.com/google/uuid"

	"pvz/internal/models"
)

type PvzForm struct {
	Id               uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}

func ToPvzForm(pvz models.Pvz) PvzForm {
	return PvzForm{
		Id:               pvz.Id,
		RegistrationDate: pvz.RegistrationDate,
		City:             pvz.City,
	}
}

type GetPvzInfoForm struct {
	StartDate time.Time
	EndDate   time.Time
	Page      int
	Limit     int
}

type GetPvzInfoResult struct {
	Pvz        PvzForm                    `json:"pvz"`
	Receptions []ReceptionProductsFormOut `json:"receptions"`
}

func ToGetPvzInfoFormOut(result []models.PvzInfo) []GetPvzInfoResult {
	var ans []GetPvzInfoResult
	for _, pvzInfo := range result {
		var receptions []ReceptionProductsFormOut

		for _, reception := range pvzInfo.Receptions {
			var products []ProductFormOut

			for _, product := range reception.Products {
				products = append(products, ToProductFormOut(product))
			}

			receptions = append(receptions, ReceptionProductsFormOut{
				Reception: ToReceptionFormOut(reception.Reception),
				Products:  products,
			})
		}

		ans = append(ans, GetPvzInfoResult{
			Pvz:        ToPvzForm(pvzInfo.Pvz),
			Receptions: receptions,
		})
	}

	return ans
}
