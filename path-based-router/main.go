package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

const (
	ConfigPath        = "/etc/nginx/nginx.conf"
	PathPrefix        = "ROUTER_PATH_"
	DestinationPrefix = "ROUTER_DESTINATION_"
	RewriteHostPrefix = "ROUTER_REWRITE_HOST_"
	RewritePathPrefix = "ROUTER_REWRITE_PATH_"
)

type renderCtx struct {
	Routes       map[string]string
	Rewrites     map[string]bool
	PathRewrites map[string]string
}

func main() {
	setupFlags()
	run()
}

func setupFlags() {
	flag.Parse()
}

func run() {
	envVars := parseEnv()
	routes, rewrites, pathRewrites := parseRoutes(envVars)
	renderConfig(routes, rewrites, pathRewrites)
	log.Printf("Successfully generated routing config")
}

func parseEnv() map[string]string {
	rawEnvVars := os.Environ()
	envVars := make(map[string]string)
	for _, r := range rawEnvVars {
		parts := strings.SplitN(r, "=", 2)
		envVars[parts[0]] = parts[1]
	}
	return envVars
}

func parseRoutes(envVars map[string]string) (map[string]string, map[string]bool, map[string]string) {
	routes := make(map[string]string)
	rewrites := make(map[string]bool)
	pathRewrites := make(map[string]string)
	for k, v := range envVars {
		if strings.HasPrefix(k, PathPrefix) {
			path := v
			entry := strings.Replace(k, PathPrefix, "", 1)

			// Find the destination
			destinationEnvVar := fmt.Sprintf("%s%s", DestinationPrefix, entry)
			destination := envVars[destinationEnvVar]
			if destination == "" {
				log.Fatalf("destination cannot be empty: %s", destinationEnvVar)
			}

			// Check if Host header rewrite is required
			var rewrite bool
			rewriteEnvVar := fmt.Sprintf("%s%s", RewriteHostPrefix, entry)
			rewriteValue := envVars[rewriteEnvVar]
			if rewriteValue == "true" {
				rewrite = true
			}

			// Check if path rewrite is mentioned
			var pathRewriteValue string
			pathRewriteKey := fmt.Sprintf("%s%s", RewritePathPrefix, entry)
			if val, ok := envVars[pathRewriteKey]; ok {
				pathRewriteValue = val
			} else {
				pathRewriteValue = path
			}
			_, ok := routes[path]
			if ok {
				log.Fatalf("path cannot be duplicate: %s", path)
			}
			routes[path] = destination
			pathRewrites[path] = pathRewriteValue
			rewrites[path] = rewrite
		}
	}
	return routes, rewrites, pathRewrites
}

func renderConfig(routes map[string]string, rewrites map[string]bool, pathRewrites map[string]string) {
	tmpl, err := template.New("nginx.conf.tmpl").ParseFiles("nginx.conf.tmpl")
	if err != nil {
		log.Fatalf("could not read template: %s", err)
	}

	file, err := os.Create(ConfigPath)
	if err != nil {
		log.Fatalf("could not create config file: %s", err)
	}

	ctx := renderCtx{
		Routes:       routes,
		Rewrites:     rewrites,
		PathRewrites: pathRewrites,
	}
	err = tmpl.Execute(file, ctx)
	if err != nil {
		file.Close()
		log.Fatalf("could not render template: %s", err)
	}

	err = file.Close()
	if err != nil {
		log.Printf("could not close config file: %s", err)
	}
}
