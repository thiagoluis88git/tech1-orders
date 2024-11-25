package bdd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/cucumber/godog"
	"github.com/go-chi/chi"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
	"github.com/thiagoluis88git/tech1-orders/internal/core/handler"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"feature"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	api := &apiFeature{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		api.resetResponse(sc)
		return ctx, nil
	})

	ctx.Step(`^I send "([^"]*)" request to "([^"]*)" with payload:$`, api.iSendRequestToWithPayload)
	ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	ctx.Step(`^the response payload should match json:$`, api.theResponsePayloadShouldMatchJson)
}

type paymentCtxKey struct{}

type apiFeature struct{}

type response struct {
	status int
	body   any
}

func (a *apiFeature) resetResponse(*godog.Scenario) {

}

func (a *apiFeature) iSendRequestToWithPayload(ctx context.Context, method, route string, payloadDoc *godog.DocString) (context.Context, error) {
	var reqBody []byte

	if payloadDoc != nil {
		payloadMap := dto.ProductForm{}
		err := json.Unmarshal([]byte(payloadDoc.Content), &payloadMap)
		if err != nil {
			panic(err)
		}

		reqBody, _ = json.Marshal(payloadMap)
	}

	req := httptest.NewRequest(method, route, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	recorder := httptest.NewRecorder()

	createProductUseCase := new(MockCreateProductUseCase)

	createProductUseCase.On("Execute", req.Context(), dto.ProductForm{
		Name:        "Batata 2",
		Description: "Batata média",
		Category:    "Acompanhamento",
		Price:       7.99,
		Images: []dto.ProducImage{
			{
				ImageUrl: "url1",
			},
		},
	}).Return(uint(3), nil)

	createProductHandler := handler.CreateProductHandler(createProductUseCase)

	createProductHandler.ServeHTTP(recorder, req)

	var paymentResponse dto.ProductCreationResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &paymentResponse)

	if err != nil {
		return nil, err
	}

	actual := response{
		status: recorder.Code,
		body:   paymentResponse,
	}

	return context.WithValue(ctx, paymentCtxKey{}, actual), nil
}

func (a *apiFeature) theResponseCodeShouldBe(ctx context.Context, expectedStatus int) error {
	resp, ok := ctx.Value(paymentCtxKey{}).(response)

	if !ok {
		return errors.New("there are no payment")
	}

	if expectedStatus != resp.status {
		if resp.status >= 400 {
			return fmt.Errorf("expected response code to be: %d, but actual is: %d, response message: %s", expectedStatus, resp.status, resp.body)
		}
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", expectedStatus, resp.status)
	}

	return nil
}

func (a *apiFeature) theResponsePayloadShouldMatchJson(ctx context.Context, expectedBody *godog.DocString) error {
	actualResp, ok := ctx.Value(paymentCtxKey{}).(response)
	if !ok {
		return errors.New("there are no payment")
	}

	var response dto.ProductCreationResponse

	err := json.Unmarshal([]byte(expectedBody.Content), &response)

	if err != nil {
		return err
	}

	if !reflect.DeepEqual(actualResp.body, response) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expectedBody, actualResp.body)
	}

	return nil
}
