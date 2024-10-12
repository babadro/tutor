// Code generated by go-swagger; DO NOT EDIT.

package swagger

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Chat chat
//
// swagger:model Chat
type Chat struct {

	// The current message in the chat. 0 based index.
	CurrentMessageIDx int32 `json:"cur_m,omitempty"`

	// prepared messages
	PreparedMessages []string `json:"prep_msgs"`

	// chat Id
	ChatID string `json:"chatId,omitempty"`

	// time
	Time int64 `json:"time,omitempty"`

	// title
	Title string `json:"title,omitempty"`

	// typ
	Typ ChatType `json:"typ,omitempty"`
}

// Validate validates this chat
func (m *Chat) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTyp(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Chat) validateTyp(formats strfmt.Registry) error {
	if swag.IsZero(m.Typ) { // not required
		return nil
	}

	if err := m.Typ.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("typ")
		}
		return err
	}

	return nil
}

// ContextValidate validate this chat based on the context it is used
func (m *Chat) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateTyp(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Chat) contextValidateTyp(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Typ.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("typ")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Chat) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Chat) UnmarshalBinary(b []byte) error {
	var res Chat
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
