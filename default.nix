let pkgs = import <nixpkgs> { };
in pkgs.buildGoModule rec {
  pname = "releasefeed";
  version = "0.0.3-snapshot";
  src = pkgs.lib.cleanSource ./.;
  vendorHash = "sha256-/OzNsgU3VNnkL9sXDoZahJ7fMqoYCEmstnNnGvmF03A=";

  # network access is required for tests
  doCheck = false;
}
