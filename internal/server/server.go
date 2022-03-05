package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alwindoss/astra"
	"github.com/alwindoss/astra/internal/dbase"
	"github.com/alwindoss/astra/internal/handlers"
	"github.com/alwindoss/astra/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.etcd.io/bbolt"
)

func Run(cfg *astra.Config) {
	// Create a template cache
	tc, err := handlers.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	cfg.TemplateCache = tc

	session := scs.New()
	session.Lifetime = time.Hour * 24
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = cfg.InProduction

	db, err := bbolt.Open(cfg.Location+"astra.db", 0600, nil)
	if err != nil {
		log.Printf("unable to open the astra db: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	repo := dbase.NewBoltDBRepository(db)
	service := service.NewService(repo)

	router := createRoutes(cfg, session, service)
	log.Printf("Running Astra server on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}

func createRoutes(cfg *astra.Config, sess *scs.SessionManager, service service.Service) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/ping"))
	n := handlers.NoSurf{
		Cfg: cfg,
	}
	router.Use(n.NoSurfMW)
	router.Use(middleware.Logger)
	router.Use(sess.LoadAndSave)
	hdlrs := handlers.NewPageHandler(cfg, sess, service)
	// handlers := &handler{
	// 	Cfg:     cfg,
	// 	SessMgr: sess,
	// }
	router.Get("/", hdlrs.RenderHomePage)
	router.Get("/bucket", hdlrs.RenderBucketPage)
	router.Post("/bucket", hdlrs.CreateBucketHandler)
	router.Get("/create-bucket", hdlrs.RenderCreateBucketPage)
	router.Get("/view-bucket", hdlrs.ViewBucketHandler)
	router.Get("/delete-bucket", hdlrs.DeleteBucketHandler)
	router.Get("/add-item", hdlrs.RenderAddItemPage)
	router.Post("/add-item", hdlrs.AddItemHandler)
	// router.Get("/about", handlers.About)
	// router.Get("/generals-quarters", handlers.Generals)
	// router.Get("/majors-suite", handlers.Majors)
	// router.Get("/search-availability", handlers.Availability)
	// router.Post("/search-availability", handlers.PostAvailability)
	// router.Get("/contact", handlers.Contact)
	// router.Get("/make-reservation", handlers.Reservation)

	staticFileServer := http.FileServer(http.FS(astra.FS))

	router.Handle("/static/*", staticFileServer)
	return router
}
