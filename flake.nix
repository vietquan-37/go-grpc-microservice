{
  description = "Custom Go 1.11.1 package for Devbox";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-21.11";
    flake-utils.url = "github:numtide/flake-utils";
    go-src = {
      url = "https://dl.google.com/go/go1.11.1.src.tar.gz";
      flake = false;
      narHash = "sha256-+5LJQV1+lwj8sJt3aCIY6QipYbJ0D2cD5vPqMefDuTk=";
    };
  };

  outputs = { self, nixpkgs, flake-utils, go-src }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        goPackage = pkgs.stdenv.mkDerivation {
          name = "go-1.11.1";
          src = go-src;
          buildInputs = [ pkgs.go ];
          buildPhase = ''
            export GOROOT_BOOTSTRAP=${pkgs.go}/share/go
            cd src
            ./make.bash
          '';
          installPhase = ''
            mkdir -p $out/bin
            cp -r bin/* $out/bin/
          '';
        };
      in
      {
        packages.default = goPackage;
        packages.go = goPackage;
      }
    );
}