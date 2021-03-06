package profile

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/obj/feed_event"
	"github.com/memocash/memo/app/profile"
	"github.com/memocash/memo/app/res"
	"github.com/memocash/memo/app/util"
	"net/http"
)

const (
	PageAll   = "all"
	PagePosts = "posts"
	PagePolls = "polls"
	PageLikes = "likes"
)

var profilePages = []string{
	PageAll,
	PagePosts,
	PagePolls,
	PageLikes,
}

var viewRoute = web.Route{
	Pattern: res.UrlProfileView + "/" + urlAddress.UrlPart(),
	Handler: func(r *web.Response) {
		addressString := r.Request.GetUrlNamedQueryVariable(urlAddress.Id)
		address := wallet.GetAddressFromString(addressString)
		pkHash := address.GetScriptAddress()
		var userPkHash []byte
		var userId uint
		if auth.IsLoggedIn(r.Session.CookieId) {
			user, err := auth.GetSessionUser(r.Session.CookieId)
			if err != nil {
				r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
				return
			}
			key, err := db.GetKeyForUser(user.Id)
			if err != nil {
				r.Error(jerr.Get("error getting key for user", err), http.StatusInternalServerError)
				return
			}
			userPkHash = key.PkHash
			userId = user.Id
		}

		offset := r.Request.GetUrlParameterInt("offset")
		pageType := r.Request.GetUrlParameter("p")
		if ! util.StringInSlice(pageType, profilePages) {
			pageType = PageAll
		}
		var events []*feed_event.Event
		var err error
		switch pageType {
		case PageAll:
			events, err = feed_event.GetUserEvents(userId, userPkHash, pkHash, uint(offset), nil)
		case PagePosts:
			events, err = feed_event.GetUserEvents(userId, userPkHash, pkHash, uint(offset), db.PostEvents)
		case PagePolls:
			events, err = feed_event.GetUserEvents(userId, userPkHash, pkHash, uint(offset), []db.FeedEventType{
				db.FeedEventCreatePoll,
			})
		case PageLikes:
			events, err = feed_event.GetUserEvents(userId, userPkHash, pkHash, uint(offset), []db.FeedEventType{
				db.FeedEventLike,
			})
		}
		if err != nil {
			r.Error(jerr.Get("error getting user events", err), http.StatusInternalServerError)
			return
		}
		r.Helper["FeedItems"] = events

		pf, err := profile.GetBasicProfile(pkHash, userPkHash)
		if err != nil {
			r.Error(jerr.Get("error getting profile for hash", err), http.StatusInternalServerError)
			return
		}
		err = pf.SetQr()
		if err != nil {
			r.Error(jerr.Get("error creating qr", err), http.StatusInternalServerError)
			return
		}

		r.Helper["Profile"] = pf
		r.Helper["PageType"] = pageType
		r.Helper["OffsetLink"] = fmt.Sprintf("%s/%s?p=%s", res.UrlProfileView, address.GetEncoded(), pageType)
		r.Helper["Title"] = fmt.Sprintf("Memo - %s's Profile", pf.Name)
		res.SetPageAndOffset(r, offset)
		r.RenderTemplate(res.UrlProfileView)
	},
}
