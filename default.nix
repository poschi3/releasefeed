let pkgs = import <nixpkgs> { };
in pkgs.buildGoModule rec {
  pname = "releasefeed";
  version = "0.0.1-snapshot";
  src = pkgs.lib.cleanSource ./.;
  vendorHash = null;

  # network access is required for tests
  doCheck = false;
}
