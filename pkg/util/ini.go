package util

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-ini/ini"
)

func GetIniSection(filePath string, sectionName string) (*ini.Section, error) {
	cfg, err := ini.Load(filePath)
	if err != nil {
		return nil, err
	}

	// 获取 [remote "origin"] 段
	section, err := cfg.GetSection(sectionName)
	if err != nil {
		return nil, err
	}
	return section, nil
}

func ParseGitConfigRemoteOrigin() (domain string, project string, err error) {
	section, err := GetIniSection(".git/config", `remote "origin"`)
	if err != nil {
		return "", "", err
	}
	return ParseGitURL(section.Key("url").String())
}

func ParseGitURL(url string) (domain, project string, err error) {
	// 定义正则表达式，匹配 HTTPS 和 SSH 格式
	// HTTPS: https://(domain)/(user)/(repo).git
	// SSH: git@(domain):(user)/(repo).git
	re := regexp.MustCompile(`^(?:https?://|git@)([^/:]+)[/:]([^/]+/[^/]+?)(?:\.git)?$`)

	matches := re.FindStringSubmatch(url)
	if len(matches) != 3 {
		return "", "", fmt.Errorf("无法解析 URL: %s", url)
	}

	domain = matches[1]
	project = strings.TrimSuffix(matches[2], ".git")

	return domain, project, nil
}
