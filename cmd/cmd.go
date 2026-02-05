package cmd

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/netbill/awsx"
	"github.com/netbill/evebox/box/inbox"
	"github.com/netbill/evebox/box/outbox"
	"github.com/netbill/logium"
	"github.com/netbill/pgdbx"
	"github.com/netbill/places-svc/internal/bucket"
	"github.com/netbill/places-svc/internal/core/modules/organization"
	"github.com/netbill/places-svc/internal/core/modules/pclass"
	"github.com/netbill/places-svc/internal/core/modules/place"
	"github.com/netbill/places-svc/internal/geogueser"
	"github.com/netbill/places-svc/internal/messenger"
	"github.com/netbill/places-svc/internal/messenger/inbound"
	"github.com/netbill/places-svc/internal/messenger/outbound"
	"github.com/netbill/places-svc/internal/repository"
	"github.com/netbill/places-svc/internal/repository/pg"
	"github.com/netbill/places-svc/internal/rest"
	"github.com/netbill/places-svc/internal/rest/controller"
	"github.com/netbill/places-svc/internal/rest/middlewares"
	"github.com/netbill/places-svc/internal/tokenmanager"
	"github.com/netbill/restkit"
)

func StartServices(ctx context.Context, cfg Config, log *logium.Logger, wg *sync.WaitGroup) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	pool, err := pgxpool.New(ctx, cfg.Database.SQL.URL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}
	db := pgdbx.NewDB(pool)

	awsCfg := aws.Config{
		Region: cfg.S3.AWS.Region,
		Credentials: credentials.NewStaticCredentialsProvider(
			cfg.S3.AWS.AccessKeyID,
			cfg.S3.AWS.SecretAccessKey,
			"",
		),
	}

	s3Client := s3.NewFromConfig(awsCfg)
	presignClient := s3.NewPresignClient(s3Client)

	awsS3 := awsx.New(
		cfg.S3.AWS.BucketName,
		s3Client,
		presignClient,
	)

	placeIconValidator := &awsx.ImgObjectValidator{
		AllowedContentTypes: cfg.S3.Upload.Place.Icon.AllowedContentTypes,
		AllowedFormats:      cfg.S3.Upload.Place.Icon.AllowedFormats,
		MaxWidth:            cfg.S3.Upload.Place.Icon.MaxWidth,
		MaxHeight:           cfg.S3.Upload.Place.Icon.MaxHeight,
		ContentLengthMax:    cfg.S3.Upload.Place.Icon.ContentLengthMax,
	}

	placeBannerValidator := &awsx.ImgObjectValidator{
		AllowedContentTypes: cfg.S3.Upload.Place.Banner.AllowedContentTypes,
		AllowedFormats:      cfg.S3.Upload.Place.Banner.AllowedFormats,
		MaxWidth:            cfg.S3.Upload.Place.Banner.MaxWidth,
		MaxHeight:           cfg.S3.Upload.Place.Banner.MaxHeight,
		ContentLengthMax:    cfg.S3.Upload.Place.Banner.ContentLengthMax,
	}

	placeCLassIconValidator := &awsx.ImgObjectValidator{
		AllowedContentTypes: cfg.S3.Upload.PlaceClass.Icon.AllowedContentTypes,
		AllowedFormats:      cfg.S3.Upload.PlaceClass.Icon.AllowedFormats,
		MaxWidth:            cfg.S3.Upload.PlaceClass.Icon.MaxWidth,
		MaxHeight:           cfg.S3.Upload.PlaceClass.Icon.MaxHeight,
		ContentLengthMax:    cfg.S3.Upload.PlaceClass.Icon.ContentLengthMax,
	}

	s3Bucket := bucket.New(bucket.Config{
		S3:                      awsS3,
		PlaceIconValidator:      placeIconValidator,
		PlaceBannerValidator:    placeBannerValidator,
		PlaceClassIconValidator: placeCLassIconValidator,
		UploadTokensTTL: bucket.UploadTokensTTL{
			Place:      cfg.S3.Upload.Token.TTL.Place,
			PlaceClass: cfg.S3.Upload.Token.TTL.PlaceClass,
		},
	})

	tokenManager := tokenmanager.New(
		cfg.S3.Upload.Token.SecretKey,
		cfg.S3.Upload.Token.TTL.Place,
		cfg.S3.Upload.Token.TTL.PlaceClass,
	)

	outBox := outbox.New(db)
	inBox := inbox.New(db)

	kafkaOutbound := outbound.New(log, outBox)

	orgMemberRolesSql := pg.NewOrgMemberRolesQ(db)
	orgMembersSql := pg.NewOrgMembersQ(db)
	organizationsSql := pg.NewOrganizationsQ(db)
	orgRolePermLinksSql := pg.NewOrgRolePermissionLinksQ(db)
	orgRolesSql := pg.NewOrgRolesQ(db)
	placesSql := pg.NewPlacesQ(db)
	placeClassesSql := pg.NewPlaceClassesQ(db)
	transactioner := pg.NewTransaction(db)

	repo := &repository.Repository{
		OrganizationsQ:          organizationsSql,
		OrgMembersQ:             orgMembersSql,
		OrgMemberRolesQ:         orgMemberRolesSql,
		OrgRolePermissionLinksQ: orgRolePermLinksSql,
		OrgRolesQ:               orgRolesSql,
		PlacesQ:                 placesSql,
		PlaceClassesQ:           placeClassesSql,
		Transactioner:           transactioner,
	}

	geogusser, err := geogueser.New()
	if err != nil {
		log.Fatal("failed to initialize geogueser", "error", err)
	}

	orgSvc := organization.New(repo)
	pclasessSvc := pclass.New(repo, kafkaOutbound, s3Bucket, tokenManager)
	placesSvc := place.New(repo, kafkaOutbound, s3Bucket, geogusser, tokenManager)

	kafkaInbound := inbound.New(log, orgSvc)

	responser := restkit.NewResponser()
	ctrl := controller.New(placesSvc, pclasessSvc, log)
	mdll := middlewares.New(log, middlewares.Config{
		AccountAccessSK: cfg.Auth.Account.Token.Access.SecretKey,
		UploadFilesSK:   cfg.S3.Upload.Token.SecretKey,
	}, responser)
	router := rest.New(log, mdll, ctrl)

	kafkaConsumer := messenger.NewConsumer(log, inBox, kafkaInbound, cfg.Kafka.Brokers...)

	kafkaProducer := messenger.NewProducer(log, outBox, cfg.Kafka.Brokers...)

	log.Infof("starting kafka brokers %s", cfg.Kafka.Brokers)

	run(func() { router.Run(ctx, cfg) })

	run(func() { kafkaConsumer.Run(ctx) })

	run(func() { kafkaProducer.Run(ctx) })
}
