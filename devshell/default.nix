{ pkgs
, ...
}:

pkgs.mkShell {
  name = "dev-shell";

  nativeBuildInputs = with pkgs; [
    go
    gofumpt
    gotools
    go-tools
    golangci-lint

    tokei

    treefmt

    jq
    nixpkgs-fmt
    shfmt
    nodePackages.prettier

    shellcheck
  ];

  shellHook = ''
    export NIX_PATH="nixpkgs=${pkgs.path}"
  '';
}
