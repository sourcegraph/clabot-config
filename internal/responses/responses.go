package responses

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"google.golang.org/api/forms/v1"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

// https://docs.google.com/forms/d/18Rd5caKejbk6ATTy-znjZofhKKeGS2AfxbsKoKFaXL8/edit
const formID = "18Rd5caKejbk6ATTy-znjZofhKKeGS2AfxbsKoKFaXL8"

type Response struct {
	CreateTime   string
	Name         string
	Email        string
	GitHubHandle string
}

const (
	fieldName         = "Name"
	fieldEmail        = "Email" // email is a form field, we don't collect RespondentEmail
	fieldGitHubHandle = "GitHub Handle"
)

func ListResponses(ctx context.Context, maxPages int) ([]Response, error) {
	// Impersonate a user
	tokenSource, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
		Subject:         os.Getenv("GOOGLE_IMPERSONATE_USER"),
		TargetPrincipal: os.Getenv("GOOGLE_TARGET_SERVICE_ACCOUNT"),
		Scopes:          []string{forms.FormsBodyReadonlyScope, forms.FormsResponsesReadonlyScope},
	})
	if err != nil {
		return nil, fmt.Errorf("impersonate.CredentialsTokenSource: %w", err)
	}
	formsService, err := forms.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("forms.NewService: %w", err)
	}

	// Get question IDs so that we can get the response
	questionIDs := map[string]string{
		fieldName:         "",
		fieldEmail:        "",
		fieldGitHubHandle: "",
	}
	form, err := formsService.Forms.Get(formID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("forms.Get: %w", err)
	}
	for _, item := range form.Items {
		if _, ok := questionIDs[item.Title]; ok {
			questionIDs[item.Title] = item.QuestionItem.Question.QuestionId
		}
	}

	// Page through responses and collect them
	var (
		responses []Response
		pages     int
	)
	if err := formsService.Forms.Responses.List(formID).Pages(ctx, func(page *forms.ListFormResponsesResponse) error {
		pages += 1
		for _, p := range page.Responses {
			name, err := getTextAnswer(p.Answers, questionIDs[fieldName])
			if err != nil {
				continue // discard record
			}
			email, err := getTextAnswer(p.Answers, questionIDs[fieldEmail])
			if err != nil {
				continue // discard record
			}
			githubHandle, err := getTextAnswer(p.Answers, questionIDs[fieldGitHubHandle])
			if err != nil {
				continue // discard record
			}

			responses = append(responses, Response{
				CreateTime:   p.CreateTime,
				Name:         name,
				Email:        email,
				GitHubHandle: cleanGitHubHandle(githubHandle),
			})
		}
		if pages > maxPages {
			return io.EOF // exit
		}
		return nil
	}); err != nil && err != io.EOF {
		return responses, fmt.Errorf("forms.Responses.List: %w", err)
	}
	return responses, err
}

func getTextAnswer(answers map[string]forms.Answer, questionID string) (string, error) {
	a, ok := answers[questionID]
	if !ok {
		return "", fmt.Errorf("no answer with key %q", questionID)
	}
	ta := a.TextAnswers
	if ta == nil {
		return "", fmt.Errorf("answer %q is not a text answer", questionID)
	}
	if len(ta.Answers) == 0 {
		return "", fmt.Errorf("answer %q has no answers", questionID)
	}
	return ta.Answers[0].Value, nil
}

// cleanGitHubHandle fixes common mistakes in provided GitHub handles
func cleanGitHubHandle(handle string) string {
	// If incorrectly prefixed with '@'
	handle = strings.TrimPrefix(handle, "@")
	// If a full URL is provided
	handle = strings.TrimPrefix(handle, "https://github.com/")
	return handle
}
