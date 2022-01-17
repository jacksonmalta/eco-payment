package app

import (
	"context"
	"fmt"
	"strconv"
)

type accreditation struct {
	log        Logger
	repository Persistence
}

func validateDocumentNumber(v string) *CreateAccountOutput {
	_, err := strconv.Atoi(v)

	if err == nil && (len([]rune(v)) == 11 || len([]rune(v)) == 14) {
		return nil
	} else {
		return &CreateAccountOutput{
			Error:  true,
			Code:   "document-invalid",
			Detail: "invalid document number",
		}
	}
}

func (a *accreditation) CreateAccountWithContext(ctx context.Context, input *CreateAccountInput) (*CreateAccountOutput, error) {
	createAccountOuput := &CreateAccountOutput{
		Error: false,
	}

	v := validateDocumentNumber(input.DocumentNumber)
	if v != nil {
		return v, nil
	}

	i := &InsertInput{
		DocumentNumber: input.DocumentNumber,
		ExternalKey:    input.ExternalKey,
	}

	res, err := a.repository.InsertWithContext(ctx, i)

	if err != nil {
		a.log.Error(fmt.Sprintf("Repository insert error %s", err.Error()))
		return nil, err
	}

	if res != nil && res.AlreadyExists {
		return &CreateAccountOutput{
			Error:  true,
			Code:   "item-already-exists",
			Detail: "item already exists",
		}, nil
	}

	return createAccountOuput, nil
}

func (a *accreditation) GetAccountWithContext(ctx context.Context, input *GetAccountInput) (*GetAccountOutput, error) {
	i := &GetInput{
		ExternalKey: input.ExternalKey,
	}

	o, err := a.repository.GetWithContext(ctx, i)
	if err != nil {
		return nil, err
	}

	if o != nil {
		return &GetAccountOutput{
			ExternalKey:    o.ExternalKey,
			DocumentNumber: o.DocumentNumber,
		}, nil
	}

	return nil, nil
}

func New(r Persistence, log Logger) Accreditation {
	return &accreditation{
		repository: r,
		log:        log,
	}
}
