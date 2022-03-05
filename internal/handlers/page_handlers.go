package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/alwindoss/astra"
	"github.com/alwindoss/astra/internal/forms"
	"github.com/alwindoss/astra/internal/service"
)

type PageHandler interface {
	RenderHomePage(w http.ResponseWriter, r *http.Request)
	RenderBucketPage(w http.ResponseWriter, r *http.Request)
	RenderCreateBucketPage(w http.ResponseWriter, r *http.Request)
	CreateBucketHandler(w http.ResponseWriter, r *http.Request)
	ViewBucketHandler(w http.ResponseWriter, r *http.Request)
	DeleteBucketHandler(w http.ResponseWriter, r *http.Request)
	AddItemHandler(w http.ResponseWriter, r *http.Request)
	RenderAddItemPage(w http.ResponseWriter, r *http.Request)
	RenderViewItemPage(w http.ResponseWriter, r *http.Request)
	DeleteItemHandler(w http.ResponseWriter, r *http.Request)
}

func NewPageHandler(cfg *astra.Config, session *scs.SessionManager, service service.Service) PageHandler {
	return &pageHandler{
		Cfg:     cfg,
		SessMgr: session,
		Svc:     service,
	}
}

type pageHandler struct {
	Cfg     *astra.Config
	SessMgr *scs.SessionManager
	Svc     service.Service
}

func (h *pageHandler) RenderHomePage(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	h.SessMgr.Put(r.Context(), "remote-ip", remoteIP)
	d := &TemplateData{
		Title: "Astra | Home",
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
		Title:       "Buckets List",
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
		stringMap := make(map[string]string)
		stringMap["bucket_name"] = bucketDtls.Name
		// stringMap["end_date"] = ed
		renderTemplate(w, r, h.Cfg, "create-bucket.page.tmpl", &TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap, // fixes error after invalid data
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
	// d := &TemplateData{
	// 	Title: "Create Bucket",
	// 	Form:  forms.New(nil),
	// }
	http.Redirect(w, r, "/bucket", http.StatusSeeOther)
	// renderTemplate(w, r, h.Cfg, "bucket.page.tmpl", d)
}

type kvPair struct {
	Key   string
	Value string
}

func (h *pageHandler) ViewBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.FormValue("bucket_name")
	log.Printf("Fetching value for the bucket: %s", bucketName)
	data, err := h.Svc.GetAllData(bucketName)
	if err != nil {
		return
	}
	var svcKVPairs []service.KeyValuePair
	var allKV []kvPair
	ok := false
	if svcKVPairs, ok = data.([]service.KeyValuePair); !ok {
		return
	}
	log.Printf("Count of key value pairs returned: %d", len(svcKVPairs))
	for _, svcKVPair := range svcKVPairs {
		kv := kvPair{
			Key:   svcKVPair.Key,
			Value: svcKVPair.Value,
		}
		allKV = append(allKV, kv)
	}
	respData := make(map[string]interface{})
	respData["kv_list"] = allKV
	respData["bucket_name"] = bucketName
	td := &TemplateData{
		Data: respData,
	}
	renderTemplate(w, r, h.Cfg, "view-bucket.page.tmpl", td)
}

func (h *pageHandler) DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.FormValue("bucket_name")
	fmt.Println("Bucket Name in Delete Handler: ", bucketName)
	err := h.Svc.DeleteBucket(bucketName)
	if err != nil {
		h.SessMgr.Put(r.Context(), "error", "cannot delete bucket")
		http.Redirect(w, r, "/bucket", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/bucket", http.StatusSeeOther)
}

func (h *pageHandler) AddItemHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.FormValue("bucket_name")
	key := r.FormValue("key")
	val := r.FormValue("val")
	log.Printf("Bucket Name in the AddItemHandler is: %s", bucketName)
	log.Printf("Key in the AddItemHandler is: %s", key)
	log.Printf("Val in the AddItemHandler is: %s", val)
	itemDls := kvPair{
		Key:   key,
		Value: val,
	}

	form := forms.New(r.PostForm)
	form.Required("key")
	form.Required("val")
	form.Required("bucket_name")
	form.MaxLength("key", 30)
	data := make(map[string]interface{})
	data["item_details"] = itemDls
	data["bucket_name"] = bucketName
	if !form.Valid() {

		// add these lines to fix bad data error
		// stringMap := make(map[string]string)
		// stringMap["end_date"] = ed
		renderTemplate(w, r, h.Cfg, "add-item.page.tmpl", &TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	err := h.Svc.Set(bucketName, key, val)
	if err != nil {
		h.SessMgr.Put(r.Context(), "error", "cannot create an item in the bucket")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// h.SessMgr.Put(r.Context(), "remote-ip", remoteIP)

	// d := &TemplateData{
	// 	Title: "View Bucket",
	// 	Form:  forms.New(nil),
	// 	Data:  data,
	// }
	http.Redirect(w, r, "/view-bucket?bucket_name="+bucketName, http.StatusSeeOther)
	// renderTemplate(w, r, h.Cfg, "view-bucket.page.tmpl", d)
}

func (h *pageHandler) RenderAddItemPage(w http.ResponseWriter, r *http.Request) {
	// remoteIP := r.RemoteAddr

	// h.SessMgr.Put(r.Context(), "remote-ip", remoteIP)
	bucketName := r.FormValue("bucket_name")
	key := r.FormValue("key")
	val := r.FormValue("val")
	log.Printf("Bucket Name in the RenderAddItemPage is: %s", bucketName)
	log.Printf("Key in the RenderAddItemPage is: %s", key)
	log.Printf("Val in the RenderAddItemPage is: %s", val)
	itemDetails := kvPair{
		Key:   key,
		Value: val,
	}
	data := make(map[string]interface{})
	data["bucket_name"] = bucketName
	data["item_details"] = itemDetails

	d := &TemplateData{
		Title: "Add Item",
		Form:  forms.New(nil),
		Data:  data,
	}
	renderTemplate(w, r, h.Cfg, "add-item.page.tmpl", d)
}

func (h *pageHandler) RenderViewItemPage(w http.ResponseWriter, r *http.Request) {
	// remoteIP := r.RemoteAddr

	// h.SessMgr.Put(r.Context(), "remote-ip", remoteIP)
	bucketName := r.FormValue("bucket_name")
	key := r.FormValue("key")
	val, err := h.Svc.Get(bucketName, key)
	if err != nil {
		h.SessMgr.Put(r.Context(), "error", "unable to get the value from the bucket")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	log.Printf("Bucket Name in the RenderViewItemPage is: %s", bucketName)
	log.Printf("Key in the RenderViewItemPage is: %s", key)
	log.Printf("Val in the RenderViewItemPage is: %s", val)
	itemDetails := kvPair{
		Key:   key,
		Value: val,
	}
	data := make(map[string]interface{})
	data["bucket_name"] = bucketName
	data["item_details"] = itemDetails

	d := &TemplateData{
		Title: "View Item",
		Form:  forms.New(nil),
		Data:  data,
	}
	renderTemplate(w, r, h.Cfg, "view-item.page.tmpl", d)
}

func (h *pageHandler) DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.FormValue("bucket_name")
	key := r.FormValue("key")
	err := h.Svc.Delete(bucketName, key)
	if err != nil {
		h.SessMgr.Put(r.Context(), "error", "unable to delete the item from the bucket")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.Printf("Bucket Name in the DeleteItemHandler is: %s", bucketName)
	log.Printf("Key in the DeleteItemHandler is: %s", key)

	// h.SessMgr.Put(r.Context(), "remote-ip", remoteIP)

	// d := &TemplateData{
	// 	Title: "View Bucket",
	// 	Form:  forms.New(nil),
	// 	Data:  data,
	// }
	http.Redirect(w, r, "/view-bucket?bucket_name="+bucketName, http.StatusSeeOther)
	// renderTemplate(w, r, h.Cfg, "view-bucket.page.tmpl", d)
}
