#!/bin/bash
echo ============================================
echo Golang Build Script
echo ============================================
export GIT_USERNAME=FoxBoxsnet
export GIT_PROJECT_REPONAME=letsencrypt.sh-dns-route53

echo .
go get github.com/${GIT_USERNAME}/${GIT_PROJECT_REPONAME}
go get github.com/mitchellh/gox
go get github.com/tcnksm/ghr

git clone https://github.com/${GIT_USERNAME}/${GIT_PROJECT_REPONAME}
cd ${GIT_PROJECT_REPONAME}

gox -ldflags "-s" -output "dist/route53_{{.OS}}_{{.Arch}}"
ghr -t ${GITHUB_TOKEN} -u ${GIT_USERNAME} -r ${GIT_PROJECT_REPONAME} `git tag` dist/

exec "$@"