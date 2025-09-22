let pkgs = import <nixpkgs> { };
in pkgs.buildGoModule rec {
  pname = "releasefeed";
  version = "0.0.5";
  src = pkgs.lib.cleanSource ./.;
  vendorHash = "sha256-Er+A+yACCFuSPsm+mQzQFpbUCKNu0khX724lpBHIi4Q=";

  # network access is required for tests
  doCheck = false;
}
