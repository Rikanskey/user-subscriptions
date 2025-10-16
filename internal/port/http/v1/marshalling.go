package v1

import (
	"github.com/go-chi/render"
	"net/http"
	"time"
	"user-subscriptions/internal/app"
)

func marshalUserSubscription(w http.ResponseWriter, r *http.Request, sub app.UserSubscription) {
	response := marshalUserSubscriptionToUserSubscriptionResponse(sub)

	render.Respond(w, r, response)
}

func marshalUserSubscriptionResponses(w http.ResponseWriter, r *http.Request, subs []app.UserSubscription) {
	sr := make([]UserSubscriptionResponse, len(subs))
	for i, s := range subs {
		sr[i] = marshalUserSubscriptionToUserSubscriptionResponse(s)
	}

	render.Respond(w, r, sr)
}

func marshalUserSubscriptionToUserSubscriptionResponse(sub app.UserSubscription) UserSubscriptionResponse {
	return UserSubscriptionResponse{
		Id:          sub.Id,
		UserId:      sub.UserId,
		ServiceName: sub.Service,
		Price:       int(sub.Price),
		StartDate:   sub.StartDate.Format("01-2006"),
		EndDate: func(ed *time.Time) *string {
			if ed != nil {
				s := ed.Format("01-2006")
				return &s
			}
			return nil
		}(sub.EndDate),
	}
}

func marshalUserSubscriptionPriceResponse(w http.ResponseWriter, r *http.Request, price uint) {
	render.Respond(w, r, UserSubscriptionSumPriceResponse(price))
}
