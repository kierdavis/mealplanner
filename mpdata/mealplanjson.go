package mpdata

import (
	"encoding/json"
	"time"
)

// MarshalJSON encodes a meal plan into its JSON form.
func (mp *MealPlan) MarshalJSON() (text []byte, err error) {
	mpj := mealPlanJSON{
		ID:        mp.ID,
		Notes:     mp.Notes,
		StartDate: mp.StartDate.Format(JSONDateFormat),
		EndDate:   mp.EndDate.Format(JSONDateFormat),
	}

	return json.Marshal(mpj)
}

// UnmarshalJSON populates the fields of the receiver with values parsed from
// the input JSON.
func (mp *MealPlan) UnmarshalJSON(text []byte) (err error) {
	var mpj mealPlanJSON
	err = json.Unmarshal(text, &mpj)
	if err != nil {
		return err
	}

	mp.ID = mpj.ID
	mp.Notes = mpj.Notes

	mp.StartDate, err = time.Parse(JSONDateFormat, mpj.StartDate)
	if err != nil {
		return err
	}

	mp.EndDate, err = time.Parse(JSONDateFormat, mpj.EndDate)
	if err != nil {
		return err
	}

	return nil
}

// MarshalJSON encodes a meal serving into its JSON form.
func (s *Serving) MarshalJSON() (text []byte, err error) {
	sj := servingJSON{
		MealPlanID: s.MealPlanID,
		Date:       s.Date.Format(JSONDateFormat),
		MealID:     s.MealID,
	}

	return json.Marshal(sj)
}

// UnmarshalJSON populates the fields of the receiver with values parsed from
// the input JSON.
func (s *Serving) UnmarshalJSON(text []byte) (err error) {
	var sj servingJSON
	err = json.Unmarshal(text, &sj)
	if err != nil {
		return err
	}

	s.MealPlanID = sj.MealPlanID
	s.MealID = sj.MealID

	s.Date, err = time.Parse(JSONDateFormat, sj.Date)
	if err != nil {
		return err
	}

	return nil
}
