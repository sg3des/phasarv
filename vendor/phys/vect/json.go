package vect 

import (
	"encoding/json"
	"log"
)

func (v Vect) MarshalJSON() ([]byte, error) {
	return json.Marshal(&[2]float32{v.X, v.Y})
}

func (v *Vect) UnmarshalJSON(data []byte) error {
	vectData := [2]float32{}

	//try unmarshalling array form
	err := json.Unmarshal(data, &vectData)
	if err != nil {
		//try other form
		vectData := struct {
			X, Y float32
		}{}

		err := json.Unmarshal(data, &vectData)

		if err != nil {
			log.Printf("Error decoding Vect")
			return err
		}
		v.X = vectData.X
		v.Y = vectData.Y
		return nil
	}

	v.X = vectData[0]
	v.Y = vectData[1]

	return nil
}
