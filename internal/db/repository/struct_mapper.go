package repository

import (
	"fmt"
	"reflect"

	"github.com/thd-spatial-ai/ignis/internal/models"
)

// initializeTabulaData creates a fully initialized TabulaBuildingParameters with all nested structs
func initializeTabulaData() *models.TabulaBuildingParameters {
	data := &models.TabulaBuildingParameters{
		BasicParameters: &models.BasicParameters{
			BuildingAppearance: &models.BuildingThematic{},
			Envelope:           &models.Envelope{},
		},
		AdvancedParameters: &models.AdvancedParameters{
			AirInfiltration:       &models.AirInfiltration{},
			ClimateConditions:     &models.ClimateConditions{},
			Uvalues:               &models.Uvalues{},
			Insulation:            &models.InsulationThicknesses{},
			SolarGains:            &models.SolarGains{},
			ThermalBridges:        &models.ThermalBridgeParameters{},
			HeatLosses:            &models.TransmissionHeatLoss{},
			ThermalResistances:    &models.ThermalResistances{},
			InsulationMeasures:    &models.InsulationPredefinedMeasures{},
			ActualInsulation:      &models.ActualInsulationThicknesses{},
			HeatTransfer:          &models.HeatTransferCoefficients{},
			PredefinedCodes:       &models.PredefinedCodes{},
			MeasureTypes:          &models.MeasureTypeCodes{},
			SolarTransmittance:    &models.SolarEnergyTransmittance{},
			MeasureFractions:      &models.MeasureAreaFractions{},
			AdditionalResistances: &models.AdditionalThermalResistance{},
		},
	}

	// Set constant default values that aren't in the database
	data.AdvancedParameters.PredefinedCodes.F_Corr_CeilingHeight = 1.0

	return data
}

// populateStructFromMap walks through all nested structs and populates fields using JSON tags
func populateStructFromMap(target interface{}, dataMap map[string]interface{}) {
	populateStruct(reflect.ValueOf(target), dataMap)
}

func populateStruct(val reflect.Value, dataMap map[string]interface{}) {
	// Dereference pointers
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Get JSON tag
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			// No JSON tag, check if it's a nested struct
			if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Struct {
				populateStruct(field, dataMap)
			} else if field.Kind() == reflect.Struct {
				populateStruct(field, dataMap)
			}
			continue
		}

		// Handle nested structs (pointers to structs)
		if field.Kind() == reflect.Ptr {
			if field.Type().Elem().Kind() == reflect.Struct {
				populateStruct(field, dataMap)
				continue
			}
		} else if field.Kind() == reflect.Struct {
			populateStruct(field, dataMap)
			continue
		}

		// Get value from map using JSON tag as key
		dbValue, ok := dataMap[jsonTag]
		if !ok || dbValue == nil {
			continue
		}

		// Set the field value
		setFieldValue(field, dbValue)
	}
}

func setFieldValue(field reflect.Value, value interface{}) {
	if !field.CanSet() {
		return
	}

	switch field.Kind() {
	case reflect.String:
		if v, ok := value.(string); ok {
			field.SetString(v)
		} else {
			field.SetString(fmt.Sprintf("%v", value))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := value.(type) {
		case int:
			field.SetInt(int64(v))
		case int64:
			field.SetInt(v)
		case int32:
			field.SetInt(int64(v))
		case float64:
			field.SetInt(int64(v))
		case float32:
			field.SetInt(int64(v))
		}
	case reflect.Float32, reflect.Float64:
		switch v := value.(type) {
		case float64:
			field.SetFloat(v)
		case float32:
			field.SetFloat(float64(v))
		case int:
			field.SetFloat(float64(v))
		case int64:
			field.SetFloat(float64(v))
		}
	}
}
