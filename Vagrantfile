# -*- mode: ruby -*-
# vi: set ft=ruby :
##############################################################################
# Copyright (c) 2020
# All rights reserved. This program and the accompanying materials
# are made available under the terms of the Apache License, Version 2.0
# which accompanies this distribution, and is available at
# http://www.apache.org/licenses/LICENSE-2.0
##############################################################################

$no_proxy = ENV['NO_PROXY'] || ENV['no_proxy'] || "127.0.0.1,localhost"
# NOTE: This range is based on vagrant-libvirt network definition CIDR 192.168.121.0/24
(1..254).each do |i|
  $no_proxy += ",192.168.121.#{i}"
end
$no_proxy += ",10.0.2.15"

File.exists?("/usr/share/qemu/OVMF.fd") ? loader = "/usr/share/qemu/OVMF.fd" : loader = File.join(File.dirname(__FILE__), "OVMF.fd")
if not File.exists?(loader)
  system('curl -O https://download.clearlinux.org/image/OVMF.fd')
end

distros = YAML.load_file(File.dirname(__FILE__) + '/distros_supported.yml')

Vagrant.configure("2") do |config|
  config.vm.provider :libvirt
  config.vm.provider :virtualbox

  config.vm.synced_folder './', '/vagrant', type: "rsync",
    rsync__args: ["--verbose", "--archive", "--delete", "-z"]
  distros["linux"].each do |distro|
    config.vm.define distro["alias"] do |node|
      node.vm.box = distro["name"]
      node.vm.box_check_update = false
      if distro["alias"] == "clearlinux"
        node.vm.provider 'libvirt' do |v|
          v.loader = loader
        end
      end
    end
  end

   # Install requirements
  config.vm.provision 'shell', privileged: false, inline: <<-SHELL
    source /etc/os-release || source /usr/lib/os-release
    case ${ID,,} in
        ubuntu|debian)
            sudo apt-get update
            sudo apt-get install -y -qq -o=Dpkg::Use-Pty=0 curl
        ;;
    esac
    # NOTE: Shorten link -> https://github.com/electrocucaracha/pkg-mgr_scripts
    curl -fsSL http://bit.ly/install_pkg | PKG="docker docker-compose make git" bash
  SHELL

  # Deploy docker compose services
  config.vm.provision 'shell', privileged: false, inline: <<-SHELL
    cd /vagrant
    make deploy-dev
  SHELL

  # Validate bash retrieval function
  config.vm.provision 'shell', privileged: false, inline: <<-SHELL
    set -o errexit
    attempt_counter=0
    max_attempts=5
    until $(curl --output /dev/null --silent --fail http://localhost:3000/install_pkg?pkg=vagrant); do
       if [ ${attempt_counter} -eq ${max_attempts} ];then
           echo "Max attempts reached"
           exit 1
       fi
       attempt_counter=$(($attempt_counter+1))
       sleep 5
    done
    curl -fsSL http://localhost:3000/install_pkg?pkg=vagrant | bash
    if ! command -v vagrant; then
        echo "Vagrant command line wasn't installed properly"
        exit 1
    fi
  SHELL

  [:virtualbox, :libvirt].each do |provider|
  config.vm.provider provider do |p|
      p.cpus = 1
      p.memory = ENV['MEMORY'] || 1024
    end
  end

  config.vm.provider "virtualbox" do |v|
    v.gui = false
  end

  config.vm.provider :libvirt do |v|
    v.random_hostname = true
    v.management_network_address = "192.168.121.0/24"
  end

  if ENV['http_proxy'] != nil and ENV['https_proxy'] != nil
    if Vagrant.has_plugin?('vagrant-proxyconf')
      config.proxy.http     = ENV['http_proxy'] || ENV['HTTP_PROXY'] || ""
      config.proxy.https    = ENV['https_proxy'] || ENV['HTTPS_PROXY'] || ""
      config.proxy.no_proxy = $no_proxy
      config.proxy.enabled = { docker: false, git: false }
    end
  end
end
