package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type StarModel struct {
	TokenId     string       `json:"token_id"`
	Coordinates string       `json:"coordinates"`
	Name        string       `json:"name"`
	Price       string       `json:"price"`
	IsForSale   bool         `json:"is_for_sale"`
	Date        time.Time    `json:"date"`
	Wallet      *WalletModel `json:"wallet"`
}

func (s StarModel) Validate() error {
	if err := s.ValidateTokenId(); err != nil {
		return err
	}
	if err := s.Wallet.Validate(); err != nil {
		return err
	}
	if err := s.ValidateCoordinates(); err != nil {
		return err
	}
	if err := s.ValidateName(); err != nil {
		return err
	}
	if s.IsForSale {
		if err := s.ValidatePrice(); err != nil {
			return err
		}
	}
	return nil
}

func (s StarModel) ValidateTokenId() error {
	pattern := "^[1-9](?:[0-9]+)?$"

	match, err := regexp.MatchString(pattern, s.TokenId)
	if err != nil {
		return err
	}

	if match {
		return nil
	}

	return errors.New("invalid token id")
}

func (s StarModel) ValidateCoordinates() error {
	pattern := "^(?:[0-1][0-9]|2[0-4])(?:[0-5][0-9]|60){2}\\.(?:[0-9][0-9])[+-](?:[0-8][0-9]|90)(?:[0-5][0-9]|60){2}\\.(?:[0-9][0-9])$"

	match, err := regexp.MatchString(pattern, s.Coordinates)
	if err != nil {
		return err
	}

	if match {
		return nil
	}

	return errors.New("invalid coordinates")
}

func (s StarModel) ValidateName() error {
	errorMsg := "invalid name"

	nameLength := len(s.Name)
	if nameLength > 32 || nameLength < 4 {
		return errors.New(errorMsg)
	}

	pattern := "^(?:[a-z]+)(?: [a-z]+)*$"

	match, err := regexp.MatchString(pattern, s.Name)
	if err != nil {
		return err
	}

	if match {
		return nil
	}

	return errors.New(errorMsg)
}

func (s StarModel) ValidatePrice() error {
	errorMsg := "invalid price"

	integerPattern := "^[1-9][0-9]*$"
	fractionPattern := "^[0-9]*[1-9]$"

	priceSlice := strings.Split(s.Price, ".")

	switch len(priceSlice) {
	case 1:
		if len(priceSlice[0]) > 12 {
			return errors.New(errorMsg)
		}

		match, err := regexp.MatchString(integerPattern, priceSlice[0])
		if err != nil {
			return err
		}

		if match {
			return nil
		}

		return errors.New(errorMsg)
	case 2:
		if len(priceSlice[0]) > 12 {
			return errors.New(errorMsg)
		}

		if len(priceSlice[1]) > 18 {
			return errors.New(errorMsg)
		}

		match, err := regexp.MatchString(integerPattern, priceSlice[0])
		if err != nil {
			return err
		}

		if !match {
			return errors.New(errorMsg)
		}

		match, err = regexp.MatchString(fractionPattern, priceSlice[1])
		if err != nil {
			return err
		}

		if !match {
			return errors.New(errorMsg)
		}

		return nil
	default:
		return errors.New(errorMsg)
	}
}

type StarRangeModel struct {
	Start       string
	End         string
	OldestFirst bool
}

func (s StarRangeModel) Validate() error {
	errorMsg := "invalid range"

	pattern := "^[1-9](?:[0-9]+)?$"

	match, err := regexp.MatchString(pattern, s.Start)
	if err != nil {
		return err
	}

	if !match {
		return errors.New(errorMsg)
	}

	match, err = regexp.MatchString(pattern, s.End)
	if err != nil {
		return err
	}

	if !match {
		return errors.New(errorMsg)
	}

	return nil
}