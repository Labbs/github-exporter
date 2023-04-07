package bootstrap

import (
	"context"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/die-net/lrucache"
	"github.com/google/go-github/v51/github"
	"github.com/gregjones/httpcache"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"

	"github.com/labbs/github-exporter/config"
	"github.com/labbs/github-exporter/internal"
)

func NewGHClient(logger zerolog.Logger) *github.Client {
	var httpClient *http.Client
	var client *github.Client

	cache := lrucache.New(config.CacheHTTPSize, 0)
	cachedTransport := httpcache.NewTransport(cache)

	if config.Github.Token != "" {
		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, cachedTransport.Client())
		httpClient = oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.Github.Token}))
	} else {
		transport, err := ghinstallation.NewKeyFromFile(
			cachedTransport,
			config.Github.ApplicationID,
			config.Github.ApplicationInstallationId,
			config.Github.ApplicationPrivateKey,
		)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error to create the authentication transport for the GitHub App")
		}
		if config.Github.EnterpriseURL != "" {
			apiUrl, err := internal.FormatEnterpriseAPIURL(config.Github.EnterpriseURL)
			if err != nil {
				logger.Fatal().Err(err).Msg("Error to format enterprise url")
			}
			transport.BaseURL = apiUrl
		}
	}

	if config.Github.EnterpriseURL != "" {
		var err error
		client, err = github.NewEnterpriseClient(config.Github.EnterpriseURL, config.Github.EnterpriseURL, httpClient)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error to create the GitHub client")
		}
	} else {
		client = github.NewClient(httpClient)
	}

	return client
}
