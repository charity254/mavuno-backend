package api

import (
	"database/sql"
	"net/http"


	"github.com/gorilla/mux"
	"github.com/mavuno/mavuno-backend/internal/config"
	"github.com/mavuno/mavuno-backend/internal/middleware"
	"github.com/mavuno/mavuno-backend/internal/services"
    "github.com/mavuno/mavuno-backend/internal/storage"
    "github.com/mavuno/mavuno-backend/internal/utils"
)

func NewRouter(db *sql.DB, cfg *config.Config) *mux.Router {
	router := mux.NewRouter()

	// ── Health Check ────────────────────────────────────────
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.JSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	}).Methods("GET")

	// ── Auth Routes ─────────────────────────────────────────
	userRepo := storage.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := NewAuthHandler(authService)

	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")
	

	 // ── Product Routes ───────────────────────────────────────
	 productRepo := storage.NewProductRepository(db)
	 productService :=services.NewProductService(productRepo)
	 productHandler := NewProductHandler(productService)


	 //AuthMiddleware enforces the required role rule
	 router.Handle("/api/products",
		middleware.AuthMiddleware(cfg.JWTSecret)(
			middleware.RequiredRole("farmer")(
				http.HandlerFunc(productHandler.CreateProduct),
			),
		),
	).Methods("POST")
	router.Handle("/api/products",
        middleware.AuthMiddleware(cfg.JWTSecret)(
            middleware.RequiredRole("farmer")(
                http.HandlerFunc(productHandler.GetProducts),
            ),
        ),
    ).Methods("GET")

    router.Handle("/api/products/{id}",
        middleware.AuthMiddleware(cfg.JWTSecret)(
            middleware.RequiredRole("farmer")(
                http.HandlerFunc(productHandler.GetProduct),
            ),
        ),
    ).Methods("GET")

    router.Handle("/api/products/{id}",
        middleware.AuthMiddleware(cfg.JWTSecret)(
            middleware.RequiredRole("farmer")(
                http.HandlerFunc(productHandler.UpdateProduct),
            ),
        ),
    ).Methods("PUT")

    router.Handle("/api/products/{id}",
        middleware.AuthMiddleware(cfg.JWTSecret)(
            middleware.RequiredRole("farmer")(
                http.HandlerFunc(productHandler.DeleteProduct),
            ),
        ),
    ).Methods("DELETE")
	// ── Entry Routes ─────────────────────────────────────────────
	// Set up the layers: repository → service → handler
	entryRepo := storage.NewProduceEntryRepository(db)
	entryService := services.NewProduceEntryService(entryRepo, productRepo)
	entryHandler := NewProduceEntryHandler(entryService)

	// All entry routes require a valid JWT token and farmer role
	router.Handle("/api/entries",
		middleware.AuthMiddleware(cfg.JWTSecret)(
			middleware.RequiredRole("farmer")(
				http.HandlerFunc(entryHandler.CreateEntry),
			),
		),
	).Methods("POST")

	router.Handle("/api/entries",
		middleware.AuthMiddleware(cfg.JWTSecret)(
			middleware.RequiredRole("farmer")(
				http.HandlerFunc(entryHandler.GetEntries),
			),
		),
	).Methods("GET")

	router.Handle("/api/entries/{id}",
		middleware.AuthMiddleware(cfg.JWTSecret)(
			middleware.RequiredRole("farmer")(
				http.HandlerFunc(entryHandler.GetEntry),
			),
		),
	).Methods("GET")

	router.Handle("/api/entries/{id}",
		middleware.AuthMiddleware(cfg.JWTSecret)(
			middleware.RequiredRole("farmer")(
				http.HandlerFunc(entryHandler.UpdateEntry),
			),
		),
	).Methods("PUT")

	router.Handle("/api/entries/{id}",
		middleware.AuthMiddleware(cfg.JWTSecret)(
			middleware.RequiredRole("farmer")(
				http.HandlerFunc(entryHandler.DeleteEntry),
			),
		),
	).Methods("DELETE")

	// ── Supply Location Routes ────────────────────────────────────
	// Set up the layers: repository → service → handler
	locRepo := storage.NewSupplyLocationRepository(db)
	locService := services.NewSupplyLocationService(locRepo)
	locHandler := NewSupplyLocationHandler(locService)

	// All supply location routes require a valid JWT token and farmer role
	router.Handle("/api/supply-locations",
		middleware.AuthMiddleware(cfg.JWTSecret)(
			middleware.RequiredRole("farmer")(
				http.HandlerFunc(locHandler.CreateSupplyLocation),
			),
		),
	).Methods("POST")

	router.Handle("/api/supply-locations",
		middleware.AuthMiddleware(cfg.JWTSecret)(
			middleware.RequiredRole("farmer")(
				http.HandlerFunc(locHandler.GetSupplyLocations),
			),
		),
	).Methods("GET")

	router.Handle("/api/supply-locations/{id}",
		middleware.AuthMiddleware(cfg.JWTSecret)(
			middleware.RequiredRole("farmer")(
				http.HandlerFunc(locHandler.GetSupplyLocation),
			),
		),
	).Methods("GET")

	router.Handle("/api/supply-locations/{id}",
		middleware.AuthMiddleware(cfg.JWTSecret)(
			middleware.RequiredRole("farmer")(
				http.HandlerFunc(locHandler.UpdateSupplyLocation),
			),
		),
	).Methods("PUT")

	router.Handle("/api/supply-locations/{id}",
		middleware.AuthMiddleware(cfg.JWTSecret)(
			middleware.RequiredRole("farmer")(
				http.HandlerFunc(locHandler.DeleteSupplyLocation),
			),
		),
	).Methods("DELETE")

	// ── Supply Agreement Routes ───────────────────────────────────
// Set up the layers: repository → service → handler
agreementRepo := storage.NewSupplyAgreementRepository(db)
agreementService := services.NewSupplyAgreementService(agreementRepo, locRepo, productRepo)
agreementHandler := NewSupplyAgreementHandler(agreementService)

// All supply agreement routes require a valid JWT token and farmer role
router.Handle("/api/supply-agreements",
    middleware.AuthMiddleware(cfg.JWTSecret)(
        middleware.RequiredRole("farmer")(
            http.HandlerFunc(agreementHandler.CreateSupplyAgreement),
        ),
    ),
).Methods("POST")

router.Handle("/api/supply-agreements",
    middleware.AuthMiddleware(cfg.JWTSecret)(
        middleware.RequiredRole("farmer")(
            http.HandlerFunc(agreementHandler.GetSupplyAgreements),
        ),
    ),
).Methods("GET")

// Active agreements endpoint — used by frontend for dashboard reminders
router.Handle("/api/supply-agreements/active",
    middleware.AuthMiddleware(cfg.JWTSecret)(
        middleware.RequiredRole("farmer")(
            http.HandlerFunc(agreementHandler.GetActiveSupplyAgreements),
        ),
    ),
).Methods("GET")

router.Handle("/api/supply-agreements/{id}",
    middleware.AuthMiddleware(cfg.JWTSecret)(
        middleware.RequiredRole("farmer")(
            http.HandlerFunc(agreementHandler.GetSupplyAgreement),
        ),
    ),
).Methods("GET")

router.Handle("/api/supply-agreements/{id}",
    middleware.AuthMiddleware(cfg.JWTSecret)(
        middleware.RequiredRole("farmer")(
            http.HandlerFunc(agreementHandler.UpdateSupplyAgreement),
        ),
    ),
).Methods("PUT")

router.Handle("/api/supply-agreements/{id}",
    middleware.AuthMiddleware(cfg.JWTSecret)(
        middleware.RequiredRole("farmer")(
            http.HandlerFunc(agreementHandler.DeleteSupplyAgreement),
        ),
    ),
).Methods("DELETE")

    return router
}