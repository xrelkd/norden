{ pkgs }:

pkgs.runCommand "check-format"
  {
    buildInputs = with pkgs; [
      treefmt
      gofumpt
      fd
      taplo
      nixfmt-rfc-style
      nodePackages.prettier
      shellcheck
      shfmt
    ];
  }

  ''
    treefmt \
      --fail-on-change \
      --no-cache \
      --allow-missing-formatter \
      --formatters go \
      --formatters prettier \
      --formatters nix \
      --formatters toml \
      --formatters shell \
      -C ${./..}

    # it worked!
    touch $out
  ''
