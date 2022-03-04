package handlers

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/alwindoss/astra"
	"github.com/alwindoss/astra/internal/dbase"
	"github.com/alwindoss/astra/internal/forms"
)

type PageHandler interface {
	RenderHomePage(w http.ResponseWriter, r *http.Request)
	RenderBucketPage(w http.ResponseWriter, r *http.Request)
	RenderCreateBucketPage(w http.ResponseWriter, r *http.Request)
	CreateBucketHandler(w http.ResponseWriter, r *http.Request)
}

func NewPageHandler(cfg *astra.Config, session *scs.SessionManager, service dbase.Service) PageHandler {
	return &pageHandler{
		Cfg:     cfg,
		SessMgr: session,
		Svc:     service,
	}
}

type pageHandler struct {
	Cfg     *astra.Config
	SessMgr *scs.SessionManager
	Svc     dbase.Service
}

func (h *pageHandler) RenderHomePage(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	h.SessMgr.Put(r.Context(), "remote-ip", remoteIP)
	d := &TemplateData{
		Title: "Home",
	}
	renderTemplate(w, r, h.Cfg, "home.page.tmpl", d)
}

func (h *pageHandler) RenderBucketPage(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	h.SessMgr.Put(r.Context(), "remote-ip", remoteIP)
	buckets, err := h.Svc.GetBuckets()
	if err != nil {
		return
	}
	d := &TemplateData{
		Title:       "Buckets",
		StringSlice: buckets,
	}
	renderTemplate(w, r, h.Cfg, "bucket.page.tmpl", d)
}

func (h *pageHandler) RenderCreateBucketPage(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	h.SessMgr.Put(r.Context(), "remote-ip", remoteIP)
	d := &TemplateData{
		Title: "Create Bucket",
		Form:  forms.New(nil),
	}
	renderTemplate(w, r, h.Cfg, "create-bucket.page.tmpl", d)
}

type bucketDetails struct {
	Name string
}

func (h *pageHandler) CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.SessMgr.Put(r.Context(), "error", "cannot parse the form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	bucketDtls := bucketDetails{
		Name: r.Form.Get("bucket_name"),
	}

	form := forms.New(r.PostForm)
	form.Required("bucket_name")
	form.MaxLength("bucket_name", 30)
	if !form.Valid() {
		data := make(map[string]interface{})
		data["bucket-details"] = bucketDtls

		// add these lines to fix bad data error
		// stringMap := make(map[string]string)
		// stringMap["start_date"] = sd
		// stringMap["end_date"] = ed
		renderTemplate(w, r, h.Cfg, "create-bucket.page.tmpl", &TemplateData{
			Form: form,
			Data: data,
			// StringMap: stringMap, // fixes error after invalid data
		})
		return
	}

	bucketName := r.Form.Get("bucket_name")
	err = h.Svc.CreateBucket(bucketName)
	if err != nil {
		h.SessMgr.Put(r.Context(), "error", "cannot create bucket")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// h.SessMgr.Put(r.Context(), "remote-ip", remoteIP)
	d := &TemplateData{
		Title: "Create Bucket",
		Form:  forms.New(nil),
	}
	renderTemplate(w, r, h.Cfg, "create-bucket.page.tmpl", d)
}
