#!/bin/bash
# SPDX-license-identifier: Apache-2.0
##############################################################################
# Copyright (c) 2019
# All rights reserved. This program and the accompanying materials
# are made available under the terms of the Apache License, Version 2.0
# which accompanies this distribution, and is available at
# http://www.apache.org/licenses/LICENSE-2.0
##############################################################################

set -o nounset
set -o errexit
set -o pipefail

# _vercmp() - Function that compares two versions
function _vercmp {
    local v1=$1
    local op=$2
    local v2=$3
    local result

    # sort the two numbers with sort's "-V" argument.  Based on if v2
    # swapped places with v1, we can determine ordering.
    result=$(echo -e "$v1\n$v2" | sort -V | head -1)

    case $op in
        "==")
            [ "$v1" = "$v2" ]
            return
            ;;
        ">")
            [ "$v1" != "$v2" ] && [ "$result" = "$v2" ]
            return
            ;;
        "<")
            [ "$v1" != "$v2" ] && [ "$result" = "$v1" ]
            return
            ;;
        ">=")
            [ "$result" = "$v2" ]
            return
            ;;
        "<=")
            [ "$result" = "$v1" ]
            return
            ;;
        *)
            die $LINENO "unrecognised op: $op"
            ;;
    esac
}

function main {
    local version=${PKG_VAGRANT_VERSION:-2.2.6}

    if command -v vagrant; then
        if _vercmp "$(vagrant version | awk 'NR==1{print $3}')" '>=' "$version"; then
            return
        fi
    fi

    pushd "$(mktemp -d)"
    vagrant_pkg="vagrant_${version}_x86_64."
    # shellcheck disable=SC1091
    source /etc/os-release || source /usr/lib/os-release
    case ${ID,,} in
        opensuse*)
            vagrant_pgp="pgp_keys.asc"
            vagrant_pkg+="rpm"
            curl -o "$vagrant_pgp" "https://keybase.io/hashicorp/$vagrant_pgp"
            curl -o "$vagrant_pkg" "https://releases.hashicorp.com/vagrant/$version/$vagrant_pkg"
            gpg --quiet --with-fingerprint "$vagrant_pgp"
            sudo rpm --import "$vagrant_pgp"
            sudo rpm --checksig "$vagrant_pkg"
            sudo rpm --install "$vagrant_pkg"
            rm $vagrant_pgp
        ;;
        ubuntu|debian)
            vagrant_pkg+="deb"
            curl -o "$vagrant_pkg" "https://releases.hashicorp.com/vagrant/$version/$vagrant_pkg"
            sudo dpkg -i "$vagrant_pkg"
        ;;
        rhel|centos|fedora)
            vagrant_pkg+="rpm"
            curl -o "$vagrant_pkg" "https://releases.hashicorp.com/vagrant/$version/$vagrant_pkg"
            PKG_MANAGER=$(command -v dnf || command -v yum)
            sudo -H -E "${PKG_MANAGER}" -q -y install "$vagrant_pkg"
        ;;
        clear-linux-os)
            vagrant_pkg="vagrant_${version}_linux_amd64.zip"
            curl -o "$vagrant_pkg" "https://releases.hashicorp.com/vagrant/$version/$vagrant_pkg"
            unzip "$vagrant_pkg"
            sudo -H -E swupd bundle-add --quiet devpkg-compat-fuse-soname2 fuse
            sudo mkdir -p /usr/local/bin
            sudo mv vagrant /usr/local/bin/
        ;;
    esac
    rm "$vagrant_pkg"
    popd
}

main
