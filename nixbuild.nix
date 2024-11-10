{ pkgs ? import <nixpkgs> { } }:

pkgs.dockerTools.buildImage {
  name = "tcpforwarder";
  tag = "latest";
  fromImage = "alpine"; # Specify the base image
  copyToRoot = ./.; # Copy the current directory content to the root of the image

  buildInputs = [ pkgs.buildPackages.docker pkgs.buildPackages.qemu ];

  runAsRoot = ''
    apk update
    apk add --no-cache bash
  '';

  platforms = [
    "x86_64-linux"
    "aarch64-linux"
  ];
}