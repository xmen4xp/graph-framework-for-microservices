package echo_server

import (
	"api-gw/pkg/authn"
	"api-gw/pkg/openapi"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"api-gw/pkg/config"
	"api-gw/pkg/model"
	"api-gw/pkg/utils"
	openMiddleware "github.com/go-openapi/runtime/middleware"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
)

type EchoServer struct {
	Echo   *echo.Echo
	Config *config.Config
}

func InitEcho(stopCh chan struct{}, conf *config.Config) {
	log.Infoln("Init Echo")
	e := NewEchoServer(conf)
	e.RegisterRoutes()
	e.Start(stopCh)
}

func (s *EchoServer) StartHTTPServer() {
	if err := s.Echo.Start(":80"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error %v", err)
	}
}

func (s *EchoServer) Start(stopCh chan struct{}) {
	// Start watching URI notification
	go func() {
		log.Debug("NodeUpdateNotifications")
		if err := s.NodeUpdateNotifications(stopCh); err != nil {
			s.StopServer()
			InitEcho(stopCh, s.Config)
		}
	}()

	// Start Server
	go func() {
		log.Info("Start Echo Server")
		if utils.IsServerConfigValid(s.Config) && utils.IsFileExists(s.Config.Server.CertPath) && utils.IsFileExists(s.Config.Server.KeyPath) {
			log.Infof("Server Config %v", s.Config.Server)
			log.Info("Start TLS Server")
			if err := s.Echo.StartTLS(s.Config.Server.Address, s.Config.Server.CertPath, s.Config.Server.KeyPath); err != nil && err != http.ErrServerClosed {
				log.Fatalf("TLS Server error %v", err)
			}
		} else {
			log.Info("Certificates or TLS port not configured correctly, hence starting the HTTP Server")
			s.StartHTTPServer()
		}
	}()
}

type NexusContext struct {
	echo.Context
	NexusURI string
	Codes    nexus.HTTPCodesResponse

	// Kube
	CrdType   string
	GroupName string
	Resource  string
}

func (s *EchoServer) RegisterRoutes() {
	// OpenAPI route
	s.Echo.GET("/openapi.json", func(c echo.Context) error {
		return c.JSON(http.StatusOK, openapi.Schema)
	})

	// Swagger-UI
	opts := openMiddleware.SwaggerUIOpts{
		SpecURL: "/openapi.json",
		Title:   "API Gateway Documentation",
	}
	s.Echo.GET("/docs", echo.WrapHandler(openMiddleware.SwaggerUI(opts, nil)))

	err := authn.RegisterCallbackHandler(s.Echo)
	if err != nil {
		log.Errorln("Error registering the OIDC callback path")
		// should we panic?
	}

	err = authn.RegisterRefreshAccessTokenEndpoint(s.Echo)
	if err != nil {
		log.Errorln("Error registering the OIDC refresh access token path")
		// should we panic?
	}

	err = authn.RegisterLogoutEndpoint(s.Echo)
	if err != nil {
		log.Errorln("Error registering the logout path")
		// should we panic?
	}
}

func (s *EchoServer) RegisterRouter(restURI nexus.RestURIs) {
	urlPattern := model.ConstructEchoPathParamURL(restURI.Uri)
	for method, codes := range restURI.Methods {
		log.Infof("Registered Router Path %s Method %s\n", urlPattern, method)
		switch method {
		case http.MethodGet:
			s.Echo.GET(urlPattern, getHandler, authn.VerifyAuthenticationMiddleware, func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					nc := &NexusContext{
						Context:  c,
						NexusURI: restURI.Uri,
						Codes:    codes,
					}
					return next(nc)
				}
			})
		case http.MethodPut:
			s.Echo.PUT(urlPattern, putHandler, authn.VerifyAuthenticationMiddleware, func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					nc := &NexusContext{
						Context:  c,
						NexusURI: restURI.Uri,
						Codes:    codes,
					}
					return next(nc)
				}
			})
		case http.MethodDelete:
			s.Echo.DELETE(urlPattern, deleteHandler, authn.VerifyAuthenticationMiddleware, func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					nc := &NexusContext{
						Context:  c,
						NexusURI: restURI.Uri,
						Codes:    codes,
					}
					return next(nc)
				}
			})
		}
	}
}

func (s *EchoServer) RegisterCrdRouter(crdType string) {
	crdParts := strings.Split(crdType, ".")
	groupName := strings.Join(crdParts[1:], ".")
	resourcePattern := fmt.Sprintf("/apis/%s/v1/%s", groupName, crdParts[0])
	resourceNamePattern := resourcePattern + "/:name"

	// TODO NPT-313 support authentication for kubectl proxy requests
	s.Echo.GET(resourceNamePattern, kubeGetByNameHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			nc := &NexusContext{
				Context:   c,
				CrdType:   crdType,
				GroupName: groupName,
				Resource:  crdParts[0],
			}
			return next(nc)
		}
	})
	s.Echo.GET(resourcePattern, kubeGetHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			nc := &NexusContext{
				Context:   c,
				CrdType:   crdType,
				GroupName: groupName,
				Resource:  crdParts[0],
			}
			return next(nc)
		}
	})
	s.Echo.POST(resourcePattern, kubePostHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			nc := &NexusContext{
				Context:   c,
				CrdType:   crdType,
				GroupName: groupName,
				Resource:  crdParts[0],
			}
			return next(nc)
		}
	})
	s.Echo.DELETE(resourceNamePattern, kubeDeleteHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			nc := &NexusContext{
				Context:   c,
				CrdType:   crdType,
				GroupName: groupName,
				Resource:  crdParts[0],
			}
			return next(nc)
		}
	})
}

func (s *EchoServer) NodeUpdateNotifications(stopCh chan struct{}) error {
	for {
		select {
		case <-stopCh:
			return fmt.Errorf("stop signal received")
		case restURIs := <-model.RestURIChan:
			log.Debugln("Rest route notification received")
			for _, v := range restURIs {
				s.RegisterRouter(v)
			}
		case crdType := <-model.CrdTypeChan:
			log.Debugln("CRD route notification received")
			s.RegisterCrdRouter(crdType)
		case oidcNodeEvent := <-model.OidcChan:
			log.Debugln("OIDC notification received")
			err := authn.HandleOidcNodeUpdate(&oidcNodeEvent, s.Echo)
			if err != nil {
				log.Errorf("error occurred while handling OIDC node update notification: %s", err)
			}
		}
	}
}

func (s *EchoServer) StopServer() {
	if err := s.Echo.Shutdown(context.Background()); err != nil {
		log.Fatalf("Shutdown signal received")
	} else {
		log.Debugln("Server exiting")
	}
}

func NewEchoServer(conf *config.Config) *EchoServer {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())

	// Setup proxy to api server
	kubeSetupProxy(e)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "ACCESS[${time_rfc3339}] method=${method}, uri=${uri}, status=${status}\n",
	}))

	return &EchoServer{
		// create a new echo_server instance
		Echo:   e,
		Config: conf,
	}
}
