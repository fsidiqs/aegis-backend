package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"

	"github.com/fsidiqs/aegis-backend/handler/authhandler"
	"github.com/fsidiqs/aegis-backend/handler/middleware"
	"github.com/fsidiqs/aegis-backend/handler/userhandler"
	"github.com/fsidiqs/aegis-backend/service/tokenservice"

	"github.com/fsidiqs/aegis-backend/repository"
	"github.com/fsidiqs/aegis-backend/service"

	"github.com/fsidiqs/aegis-backend/db"
	"github.com/gin-gonic/gin"
)

// inject will initialize a handler starting from database
// which inject into repository layer
// which inject into service layer
// which inject into handler layer

func inject(d *db.DataSources, m *mailer) (*gin.Engine, error) {
	// baseLogger := logservice.NewLogger(&logservice.LogWriter{}, log.Ldate|log.Ltime)

	log.Println("Injecting data sources")
	//
	// repository layer
	//
	userRepository := repository.NewUserRepository(d.DB)
	// redisRepository := repository.NewRedisRepository(d.RedisClient)
	// programRepository := repository.NewProgramRepository(d.DB)
	// subscriptionRepository := repository.NewSubscriptionRepository(d.DB)
	// paymentRepository := repository.NewPaymentRepository(d.DB)
	// moodTrackerRepository := repository.NewMoodTrackerRepository(d.DB)
	// activityLogRepository := repository.NewActivityLogRepository(d.DB)
	// safeRoomRepository := repository.NewSafeRoomRepository(d.DB)
	// promoRepository := repository.NewPromoRepository(d.DB)
	// userPointRepository := repository.NewUserPointHistoryRepository(d.DB)
	//
	// mail client
	//

	//
	// service layer
	//
	// emailVerifTokSecret := os.Getenv("EMAIL_VERIFICATION_TOKEN_SECRET")
	// emailVerifTokExp := os.Getenv("EMAIL_VERIFICATION_TOKEN_EXP")

	otpValueFromStr := os.Getenv("OTP_VALUE_FROM")
	otpMaxValueStr := os.Getenv("OTP_VALUE_LESS_THAN")
	otpExpirationSecsStr := os.Getenv("OTP_EXP")

	// emailVerifExp, err := strconv.ParseInt(emailVerifTokExp, 0, 64)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not parse EMAIL_VERIFICATION_TOKEN_EXP: %v", err)
	// }

	otpMaxValue, err := strconv.Atoi(otpMaxValueStr)
	if err != nil {
		return nil, fmt.Errorf("could not parse OTP_MAX_VALUE : %v", err)
	}

	otpValueFrom, err := strconv.Atoi(otpValueFromStr)
	if err != nil {
		return nil, fmt.Errorf("could not parse OTP_VALUE_FROM: %v", err)
	}

	otpExpirationSecs, err := strconv.ParseInt(otpExpirationSecsStr, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse OTP_EXP as int: %v", err)
	}

	//---------- userService
	userService := service.NewUserService(
		&service.USConfig{
			UserRepository: userRepository,
			// RedisRepository: redisRepository,
			MailClient: m.client,
			USVerificationConfig: service.USVerificationConfig{
				// EmailVerificationTokenSecret: emailVerifTokSecret,
				// EmailTokenExpirationSecs:     emailVerifExp,

				OTPValueFrom:      otpValueFrom,
				OTPMaxValue:       otpMaxValue,
				OTPExpirationSecs: otpExpirationSecs,
			},
		},
	)
	//---------- activityLogService
	// activityLogService := service.NewActivityLogService(&service.ActivityLogServiceConfig{
	// 	ActivityLogRepository: activityLogRepository,
	// })

	privateKeyFile := os.Getenv("PRIVATE_KEY_FILE")
	priv, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read priavte key pem file: %v", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %v", err)
	}

	pubKeyFile := os.Getenv("PUB_KEY_FILE")
	pub, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read public key pem file: %v", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		return nil, fmt.Errorf("could not parse the public key: %v", err)
	}

	// load refresh token secret from env variable
	refreshSecret := os.Getenv("REFRESH_SECRET")
	// load expiration lengths from env variables and parse as int
	authTokenExp := os.Getenv("AUTH_TOKEN_EXP")
	refreshTokenExp := os.Getenv("REFRESH_TOKEN_EXP")

	idExp, err := strconv.ParseInt(authTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse AUTH_TOKEN_EXP as int: %v", err)
	}

	refreshExp, err := strconv.ParseInt(refreshTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse REFRESH_TOKEN_EXP as int: %v", err)
	}

	publicTokenSecret := os.Getenv("PUBLIC_TOKEN_SECRET")
	publicTokenExpStr := os.Getenv("PUBLIC_TOKEN_EXP")
	publicTokenExp, err := strconv.ParseInt(publicTokenExpStr, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse PUBLIC_TOKEN_EXP as int: %v", err)
	}

	//------------- tokenService
	tokenService := service.NewTokenService(&service.TSConfig{
		// RedisRepository:       redisRepository,
		PrivKey:               privKey,
		PubKey:                pubKey,
		RefreshSecret:         refreshSecret,
		IDExpirationSecs:      idExp,
		RefreshExpirationSecs: refreshExp,
	})

	//-------------- publicTokenService
	publicTokenService := tokenservice.NewPublicTokenService(&tokenservice.PublicTSConfig{
		Secret:               publicTokenSecret,
		SecretExpirationSecs: publicTokenExp,
	})

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	// read in Base_URL
	baseURL := os.Getenv("API_URL")
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	//--------------- firebaseService
	// firebaseServiceAccountKeyPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY")
	// fAuthClient, err := firebaseapp.NewFirebaseAuth(firebaseServiceAccountKeyPath)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to initiate firebase auth client:%v\n", err)
	// }

	// amplitudeApiKey := os.Getenv("AMPLITUDE_API_KEY")
	// amplitudeEventTracker := amplitude.NewAmplitudeClient(amplitudeApiKey)

	if err := authhandler.NewHandler(&authhandler.Config{
		R:               router,
		PubTokenService: publicTokenService,
		BaseURL:         baseURL,
		// Log:             baseLogger,
	}); err != nil {
		return nil, err
	}

	//---------------- promoService
	// promoService := promoservice.NewPromoService(&promoservice.PromoSvcConfig{
	// 	PromoRepository: promoRepository,
	// })

	//---------------- userpointhistoryservice
	// userPointHistoryService := userpointhistorysvc.NewUserPointHistory(&userpointhistorysvc.UserPointHistoryCfg{
	// 	PromoRepository:     promoRepository,
	// 	UserPointHistoryRep: userPointRepository,
	// })

	userhandler.NewHandler(&userhandler.Config{
		R:            router,
		UserService:  userService,
		TokenService: tokenService,
		// PromoService:            promoService,
		// UserPointHistoryService: userPointHistoryService,
		PublicTokenService: publicTokenService,
		MailClient:         m.client,
		// ActivityLogService:      activityLogService,
		BaseURL: baseURL,
		// FirebaseAuth:            fAuthClient,
		// Log:                     baseLogger,
	})

	//------------- Program Service
	// vidStorage, err := initVideoStorage()
	// if err != nil {
	// 	log.Printf("error initVideoStorage: %v\n", err)

	// 	return nil, err
	// }
	// generalStorage, err := initStorage()
	// if err != nil {
	// 	log.Printf("error initStorage: %v\n", err)
	// 	return nil, err
	// }
	// maxEnrCoachingStr := os.Getenv("MAX_ENROLLMENT_COACHING")
	// maxEnrCoaching, err := strconv.Atoi(maxEnrCoachingStr)
	// if err != nil {
	// 	log.Printf("error parsing max coaching enrollments per user:%+v\n", err)
	// 	return nil, err
	// }

	// maxEnrFitStr := os.Getenv("MAX_ENROLLMENT_FITNESS")
	// maxEnrFit, err := strconv.Atoi(maxEnrFitStr)
	// if err != nil {
	// 	log.Printf("error parsing max fitness enrollments per user:%+v\n", err)
	// 	return nil, err
	// }

	// programService := service.NewProgramService(&service.PSConfig{
	// 	RedisRepository:   redisRepository,
	// 	ProgramRepository: programRepository,
	// 	VideoStorage:      vidStorage.client,
	// 	GeneralStorage:    generalStorage.client,
	// 	MaxEnrollConfig: service.MaxEnrollConfig{
	// 		MaxEnrollCoaching: maxEnrCoaching,
	// 		MaxEnrollFitness:  maxEnrFit,
	// 	},
	// })

	// envDraftProgramWhitelist := os.Getenv("DRAFT_PROGRAM_WHITELIST")
	// draftProgramWhitelist := strings.Split(envDraftProgramWhitelist, ",")
	// programhandler.NewHandler(&programhandler.Config{
	// 	R:                  router,
	// 	UserService:        userService,
	// 	TokenService:       tokenService,
	// 	PublicTokenService: publicTokenService,
	// 	// ProgramService:     programService,
	// 	// ActivityLogService: activityLogService,
	// 	BaseURL: baseURL,
	// 	Log:     baseLogger,

	// 	EmailWhitelist: draftProgramWhitelist,
	// })

	// Subscription Service
	// subscriptionService := service.NewSubscriptionService(&service.SSConfig{
	// 	SubscriptionRepository: subscriptionRepository,
	// })

	//---------------- paymentService
	// xenditApiPrivKey := os.Getenv("XENDIT_PRIVATE_API_KEY")
	// xenditCallbackToken := os.Getenv("XENDIT_CALLBACK_VERIFICATION_TOKEN")
	// paymentClient, err := xendit.NewXenditClient(xenditApiPrivKey)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to creating paymentService:%v\n", err)
	// }
	// paymentService := paymentservice.NewPaymentService(&paymentservice.PSConfig{
	// 	// PaymentRepository: paymentRepository,
	// 	Client: paymentClient,
	// })

	// subscriptionhandler.NewHandler(&subscriptionhandler.Config{
	// 	R:              router,
	// 	UserService:    userService,
	// 	TokenService:   tokenService,
	// PaymentService: paymentService,
	// SubscriptionService: subscriptionService,
	// XenditCallbackToken: xendit.XenditCallbackToken(xenditCallbackToken),
	// ActivityLogService:  activityLogService,
	// 	EventTracker: amplitudeEventTracker,

	// 	BaseURL: baseURL,
	// 	Log:     baseLogger,
	// })

	// MoodTracker Service
	// moodTrackerService := service.NewMoodTrackerService(&service.MTSConfig{
	// 	MoodTrackerRepository: moodTrackerRepository,
	// })

	// moodtrackerhandler.NewHandler(&moodtrackerhandler.Config{
	// 	R:                  router,
	// 	UserService:        userService,
	// 	TokenService:       tokenService,
	// 	MoodTrackerService: moodTrackerService,
	// 	ActivityLogService: activityLogService,
	// 	BaseURL:            baseURL,
	// })

	// activityloghandler.NewHandler(&activityloghandler.Config{
	// 	R:                  router,
	// 	TokenService:       tokenService,
	// 	ActivityLogService: activityLogService,
	// 	BaseURL:            baseURL,
	// })

	// safeRoomService := service.NewSafeRoomService(&service.SafeRoomServiceConfig{
	// 	SafeRoomRepository: safeRoomRepository,
	// })
	// saferoomhandler.NewHandler(&saferoomhandler.Config{
	// 	R:               router,
	// 	TokenService:    tokenService,
	// 	SafeRoomService: safeRoomService,
	// 	BaseURL:         baseURL,
	// })

	return router, nil
}
