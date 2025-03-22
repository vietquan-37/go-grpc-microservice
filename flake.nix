{
  description = "Local Flake for Go 1.11.1 and protoc-gen-go";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-20.03";

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in {
      packages = {
        x86_64-linux = {
          go_1_11_1 = pkgs.stdenv.mkDerivation rec {
            pname = "go";
            version = "1.11.1";

            src = pkgs.fetchurl {
              url = "https://dl.google.com/go/go1.11.1.linux-amd64.tar.gz";
              sha256 = "sha256-KHEnDY/wyMafFhqq5C+fKHOYVf9cUgR1Ko2Socn2OZM=";
            };

            installPhase = ''
              mkdir -p $out
              tar -C $out -xzf $src --strip-components=1
            '';
          };
        };
      };
    };
}
