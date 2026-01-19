package utils

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/mozillazg/go-unidecode"
)

var provinces = map[string]uint{
	"ha_noi":      10,
	"cao_bang":    30,
	"tuyen_quang": 50,
	"dien_bien":   11,
	"lai_chau":    12,
	"son_la":      14,
	"lao_cai":     15,
	"thai_nguyen": 19,
	"lang_son":    20,
	"quang_ninh":  22,
	"bac_ninh":    24,
	"phu_tho":     25,
	"hai_phong":   31,
	"hung_yen":    33,
	"ninh_binh":   37,
	"thanh_hoa":   38,
	"nghe_an":     40,
	"ha_tinh":     42,
	"quang_tri":   44,
	"hue":         46,
	"da_nang":     48,
	"quang_ngai":  51,
	"gia_lai":     52,
	"khanh_hoa":   56,
	"dak_lak":     66,
	"lam_dong":    68,
	"dong_nai":    75,
	"ho_chi_minh": 79,
	"tay_ninh":    80,
	"dong_thap":   82,
	"vinh_long":   86,
	"an_giang":    91,
	"can_tho":     92,
	"ca_mau":      96,
}

var universities = map[string]string{
	"hanoi_university_of_science_and_technology": "HUST",
	"national_economics_university":              "NEU",
	"hanoi_medical_university":                   "HMU",
	"foreign_trade_university":                   "FTU",
	"vietnam_national_university_hanoi":          "VNU",
}

var majors = map[string]string{
	"biological_engineering":                          "BF1",
	"food_engineering":                                "BF2",
	"chemical_engineering":                            "CH1",
	"educational_technology":                          "ED2",
	"education_management":                            "ED3",
	"electrical_engineering":                          "EE1",
	"control_and_automation_engineering":              "EE2",
	"energy_management":                               "EM1",
	"industrial_management":                           "EM2",
	"business_administration":                         "EM3",
	"accounting":                                      "EM4",
	"finance_and_banking":                             "EM5",
	"electronics_and_telecommunications_engineering":  "ET1",
	"biomedical_engineering":                          "ET2",
	"environmental_engineering":                       "EV2",
	"national_resources_and_environmental_management": "EV2",
	"thermal_engineering":                             "HE1",
	"computer_science":                                "IT1",
	"computer_engineering":                            "IT2",
	"data_science_and_artificial_intelligence":        "IT-E10",
	"cyber_security":                                  "IT-E15",
	"mechanical_engineering":                          "ME1",
	"mechatronics_engineering":                        "ME2",
	"mathematics_and_computer_science":                "MI1",
	"materials_engineering":                           "MS1",
	"microelectronics_and_nanotechnology_engineering": "MS2",
	"polymer_and_composite_materials_technology":      "MS3",
	"printing_engineering":                            "MS5",
	"engineering_physics":                             "PH1",
	"nuclear_engineering":                             "PH2",
	"medical_physics":                                 "PH2",
	"automotive_engineering":                          "TE1",
	"dynamics_mechanical_engineering":                 "TE2",
	"aeronautical_engineering":                        "TE3",
	"textile_and_garment_technology":                  "TX1",
	"english_language":                                "7220201",
	"chinese_language":                                "7220202",
	"economics":                                       "7310101_1",
	"mathematical_economics":                          "7310108",
	"marketing":                                       "7340115",
	"international_business":                          "7340120",
	"commercial_business":                             "7340121",
	"auditing":                                        "7340302",
	"law":                                             "7380101",
	"economic_law":                                    "7380107",
	"logistics_and_supply_chain_management":           "7510605",
	"public_management":                               "7340403",
	"psychology":                                      "7310401",
	"general_medicine":                                "7720101",
	"nursing":                                         "7720301",
	"midwifery":                                       "7720302",
	"nutrition":                                       "7720401",
	"dentistry":                                       "7720501",
	"optometry":                                       "7720699",
	"medical_laboratory_technology":                   "7720601",
	"medical_imaging_techniques":                      "7720602",
	"rehabilitation_techniques":                       "7720603",
	"public_health":                                   "7720701",
	"mathematics":                                     "QHT01",
	"physics":                                         "QHT03",
	"chemistry":                                       "QHT06",
	"biology":                                         "QHT08",
	"geography":                                       "QHT10",
	"environment":                                     "QHT13",
	"meteorology_and_climatology":                     "QHT16",
	"oceanography":                                    "QHT17",
	"geology":                                         "QHT18",
	"journalism":                                      "QHX01",
	"political_science":                               "QHX02",
	"asian_studies":                                   "QHX06",
	"europe_studies":                                  "QHX07",
	"religious_studies":                               "QHX23",
	"cultural_studies":                                "QHX25",
	"history":                                         "QHX10",
	"linguistics":                                     "QHX12",
	"anthropology":                                    "QHX13",
	"philosophy":                                      "QHX24",
	"literature":                                      "QHX26",
	"sociology":                                       "QHX28",
}

func RandomDigits(n int) (string, error) {
	result := make([]byte, n)

	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		result[i] = byte('0') + byte(num.Int64())
	}

	return string(result), nil
}

func normalize(s string) string {
	s = strings.TrimSpace(s)
	s = unidecode.Unidecode(s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	return s
}

func GetIdNumber() (string, error) {

	serial, err := RandomDigits(12)

	return serial, err
}

func GetDegreeNumber() (string, error) {

	serial, err := RandomDigits(6)

	return serial, err
}

func GetInsuranceNumber(insuranceType string) (string, error) {
	serial, err := RandomDigits(6)
	return serial, err
}

func GetLicenseNumber(class string) (string, error) {
	serial, err := RandomDigits(11)
	return class + "-" + serial, err
}

func GetPassportNumber(nation string) (string, error) {
	serial, err := RandomDigits(11)
	return nation + "-" + serial, err
}
