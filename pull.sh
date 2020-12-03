#!/bin/bash -ex

if [ "$1" == "--force" ]; then
  force="yes"
fi

srcsha="--root"
dstsha="HEAD"
commitrange="${srcsha} ${dstsha}"

while read srcpath dstpath; do
  if [ "${srcpath::1}" == "#" ] ; then
    continue
  fi
  
  relpath="${srcpath//\*}"
  
  rm -rf patches/
 
  dstsha=$(git rev-parse HEAD)
  if [ -f "../${dstpath}/.synced" ]; then
    srcsha=$(cat "../${dstpath}/.synced" | tr -d '\n')
    commitrange="${srcsha}..${dstsha}"
  fi

  git format-patch --find-copies --break-rewrites --find-renames=100% --relative="${relpath}" --no-stat --minimal --minimal --no-cover-letter --no-signature "${commitrange}" -o patches/ -- "${srcpath}"
 
  for p in $(ls patches/); do
    grep -q 'From: Vasiliy Tolstov <v.tolstov' "patches/${p}" || sed -i '/Signed-off-by: Vasiliy Tolstov/d' "patches/${p}"
  done
  
  if [ "x$(find patches/ -type f -name '*.patch' | wc -l)" != "x0" ]; then
    pushd ../${dstpath} >/dev/null
    git am --rerere-autoupdate --3way ../micro/patches/*.patch
    popd >/dev/null
  fi

  echo -n "${dstsha}" > ../${dstpath}/.synced

done < mapping.txt

