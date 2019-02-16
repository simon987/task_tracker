package api

import (
	"crypto"
	"crypto/hmac"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/simon987/task_tracker/config"
	"github.com/simon987/task_tracker/storage"
	"github.com/valyala/fasthttp"
	"hash"
	"strings"
)

type GitPayload struct {
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	After      string `json:"after"`
	Repository struct {
		Id    int64 `json:"id"`
		Owner struct {
			Id       int64  `json:"id"`
			Username string `json:"username"`
			Login    string `json:"login"`
			FullName string `json:"full_name"`
			Email    string `json:"email"`
		} `json:"owner"`
		Name          string `json:"name"`
		FullName      string `json:"full_name"`
		Private       bool   `json:"private"`
		Fork          bool   `json:"fork"`
		Size          int64  `json:"size"`
		HtmlUrl       string `json:"html_url"`
		SshUrl        string `json:"ssh_url"`
		CloneUrl      string `json:"clone_url"`
		DefaultBranch string `json:"default_branch"`
	} `json:"repository"`
}

func (g GitPayload) String() string {
	jsonBytes, _ := json.Marshal(g)
	return string(jsonBytes)
}

func (api *WebAPI) ReceiveGitWebHook(r *Request) {

	if !signatureValid(r) {
		logrus.Error("WebHook signature does not match!")
		r.Ctx.SetStatusCode(403)
		_, _ = fmt.Fprintf(r.Ctx, "Signature does not match")
		return
	}

	payload := &GitPayload{}
	err := json.Unmarshal(r.Ctx.Request.Body(), payload)
	if err != nil {
		r.Ctx.SetStatusCode(400)
		return
	}

	logrus.WithFields(logrus.Fields{
		"payload": payload,
	}).Info("Received git WebHook")

	if !isProductionBranch(payload) {
		return
	}

	project := api.getAssociatedProject(payload)
	if project == nil {
		return
	}

	version := getVersion(payload)

	project.Version = version
	err = api.Database.UpdateProject(project)
	handleErr(err, r)
}

func signatureValid(r *Request) (matches bool) {

	signature := parseSignatureFromRequest(r.Ctx)

	if signature == "" {
		return false
	}

	body := r.Ctx.PostBody()

	mac := hmac.New(getHashFuncFromConfig(), config.Cfg.WebHookSecret)
	mac.Write(body)

	expectedMac := hex.EncodeToString(mac.Sum(nil))
	matches = strings.Compare(expectedMac, signature) == 0

	logrus.WithFields(logrus.Fields{
		"expected":  expectedMac,
		"signature": signature,
		"matches":   matches,
	}).Trace("Validating WebHook signature")

	return
}

func getHashFuncFromConfig() func() hash.Hash {

	if config.Cfg.WebHookHash == "sha1" {
		return crypto.SHA1.New
	} else if config.Cfg.WebHookHash == "sha256" {
		return crypto.SHA256.New
	}

	logrus.WithFields(logrus.Fields{
		"hash": config.Cfg.WebHookHash,
	}).Error("Invalid hash function from config")

	return nil
}

func parseSignatureFromRequest(ctx *fasthttp.RequestCtx) string {

	signature := string(ctx.Request.Header.Peek(config.Cfg.WebHookSigHeader))
	sigParts := strings.Split(signature, "=")
	signature = sigParts[len(sigParts)-1]

	return signature
}

func (api *WebAPI) getAssociatedProject(payload *GitPayload) *storage.Project {

	project := api.Database.GetProjectWithRepoName(payload.Repository.FullName)

	logrus.WithFields(logrus.Fields{
		"project": project,
	}).Trace("Found project associated with WebHook")

	return project
}

func isProductionBranch(payload *GitPayload) (isProd bool) {

	isProd = strings.HasSuffix(payload.Ref, "master")

	logrus.WithFields(logrus.Fields{
		"isProd": isProd,
	}).Trace("Identified if push event occured in production branch")

	return
}

func getVersion(payload *GitPayload) (version string) {

	version = payload.After

	logrus.WithFields(logrus.Fields{
		"version": version,
	}).Trace("Got new version")

	return
}
