#!/bin/bash
set -euo pipefail

DIR=$(dirname "${BASH_SOURCE[0]}" )
version=$(cat "$DIR/VERSION")
binary=test-vmbackend
kbdir="/keybase/team/bitnami.publicservices/${binary}"
archs="linux darwin windows"

t=$(mktemp -d)
function cleanup {
  rm -rf "${t}"
}
trap cleanup EXIT

suffix() {
  os=$1
  case $os in
    windows) echo .exe ;;
    *) ;;
  esac
}

releasedir="${kbdir}/releases/${binary}-${version}"
if [ -d "${releasedir}" ]; then
  echo "Release exists already at ${releasedir}"
  echo "You shoudn't re-release the same version, please do ./next-release.sh first"
  exit 1
fi

go_ld_flags="-s -w -X 'main.Version=${version}'"

echo "Building release binaries"
for os in ${archs}; do
  GOOS=${os} go build -o "${t}/${os}" -ldflags="${go_ld_flags}"
done

echo "Copying binaries to keybase"
keybase fs mkdir "${releasedir}"
for os in ${archs}; do
  outdir=${releasedir}/${os}-amd64
  keybase fs mkdir "${outdir}"

  out="${outdir}/${binary}$(suffix "$os")"
  keybase fs cp -f "${t}/${os}" "${out}"
  chmod +x "${out}"
  echo "${out}"
done

echo "Zip binaries & sources"
pushd "${kbdir}/releases/"
zip -r "${binary}-${version}.zip" "${binary}-${version}"/*
popd
pushd ..
zip -r "${kbdir}/releases/${binary}-${version}.zip" \
  "./${binary}" -x '*.git*' -x "*release.sh" -x "*next-version.sh" \
  -x "*vms.json" -x "*test-vmbackend" -x "*.vscode*"
popd

echo "Pointing latest to ${version}"
rm -f "${kbdir}/releases/latest"
ln -snf "${releasedir}" "${kbdir}/releases/latest"
