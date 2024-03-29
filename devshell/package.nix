{ name
, version
, lib
, buildGoModule
, installShellFiles
}:

buildGoModule rec {
  pname = name;
  inherit version;

  src = lib.cleanSource ./..;

  vendorHash = "sha256-g1xHsz8tMZVGE17wEmye587QP5FQvdy4ztp68HV9Up8=";

  subPackages = [ "cmd/norden" ];

  ldflags = [
    "-s"
    "-w"
    "-X github.com/xrelkd/norden/pkg/version.AppName=${pname}"
    "-X github.com/xrelkd/norden/pkg/version.Version=${version}"
  ];

  nativeBuildInputs = [ installShellFiles ];

  postInstall = ''
    installShellCompletion --cmd norden \
      --bash <($out/bin/norden completion bash) \
      --fish <($out/bin/norden completion fish) \
      --zsh  <($out/bin/norden completion zsh)
  '';
}
