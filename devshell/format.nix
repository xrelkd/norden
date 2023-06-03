{ pkgs, }:

pkgs.runCommand "check-format"
{
  buildInputs = with pkgs; [
    treefmt
    gofumpt
    fd
    nixpkgs-fmt
    nodePackages.prettier
    shellcheck
    shfmt
  ];
} ''
  treefmt \
    --fail-on-change \
    --no-cache \
    -C ${./..}

  # it worked!
  touch $out
''
