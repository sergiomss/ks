package user

import (
	"gopkg.in/AlecAivazis/survey.v1"
)

func Prompt(p survey.Prompt) (string, error) {
	var qs = []*survey.Question{
		{
			Name:   "answer",
			Prompt: p,
		},
	}

	answers := struct {
		Answer string `survey:"answer"`
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return "", err
	}
	return answers.Answer, nil
}
