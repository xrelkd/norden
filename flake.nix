{
  description = "Norden";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    let
      name = "norden";
      version = "0.2.0";
    in
    (flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = import nixpkgs {
            inherit system;
            overlays = [
              self.overlays.default
            ];
          };

        in
        {
          formatter = pkgs.treefmt;

          devShells.default = pkgs.callPackage ./devshell { };

          packages = rec {
            default = norden;
            norden = pkgs.callPackage ./devshell/package.nix {
              inherit name version;
            };
          };

          checks = {
            format = pkgs.callPackage ./devshell/format.nix { };
          };
        })) // {
      overlays.default = final: prev: { };
    };
}
