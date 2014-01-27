package mpdata

import (
	"encoding/json"
	"time"
)

func (mp *MealPlan) MarshalJSON() (text []byte, err error) {
	mpj := mealPlanJson{
		ID:        mp.ID,
		Notes:     mp.Notes,
		StartDate: mp.StartDate.Format(JsonDateFormat),
		EndDate:   mp.EndDate.Format(JsonDateFormat),
	}

	return json.Marshal(mpj)
}

func (mp *MealPlan) UnmarshalJSON(text []byte) (err error) {
	var mpj mealPlanJson
	err = json.Unmarshal(text, &mpj)
	if err != nil {
		return err
	}

	mp.ID = mpj.ID
	mp.Notes = mpj.Notes

	mp.StartDate, err = time.Parse(JsonDateFormat, mpj.StartDate)
	if err != nil {
		return err
	}

	mp.EndDate, err = time.Parse(JsonDateFormat, mpj.EndDate)
	if err != nil {
		return err
	}

	return nil
}

func (s *Serving) MarshalJSON() (text []byte, err error) {
	sj := servingJson{
		MealPlanID: s.MealPlanID,
		Date:       s.Date.Format(JsonDateFormat),
		MealID:     s.MealID,
	}

	return json.Marshal(sj)
}

func (s *Serving) UnmarshalJSON(text []byte) (err error) {
	var sj servingJson
	err = json.Unmarshal(text, &sj)
	if err != nil {
		return err
	}

	s.MealPlanID = sj.MealPlanID
	s.MealID = sj.MealID

	s.Date, err = time.Parse(JsonDateFormat, sj.Date)
	if err != nil {
		return err
	}

	return nil
}
